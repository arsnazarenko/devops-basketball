// Package app uses for running application
package app

import (
	"log"
	"net"
	"net/http"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/config"
	v1 "github.com/arsnazarenko/devops-basketball/internal/controller/http/v1"
	"github.com/arsnazarenko/devops-basketball/internal/metrics"
	"github.com/arsnazarenko/devops-basketball/internal/usecase"
	"github.com/arsnazarenko/devops-basketball/internal/usecase/repo"
	"github.com/arsnazarenko/devops-basketball/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	oapi_middleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	swagger, err := gen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %s", err)
	}
	swagger.Servers = nil

	pg, err := postgres.New(config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()
	// create chi router
	r := chi.NewRouter()
	// create PlayerServer
	repo := repo.NewPlayerRepo(pg)
	player := usecase.NewPlayerUsecase(repo)
	serversImpl := v1.NewPlayersServerImpl(player)

	server := gen.NewStrictHandler(serversImpl, []gen.StrictMiddlewareFunc{})

	// cors
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	// cors middleware
	r.Use(corsHandler)
	// openapi validation middleware
	r.Use(oapi_middleware.OapiRequestValidator(swagger))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// metrics middleware
	r.Use(metrics.HTTPMetricsMiddleware())

	gen.HandlerFromMux(server, r)
	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort(config.HTTP.Host, config.HTTP.Port),
	}

	// metrics server
	go func() {
		mr := chi.NewMux()
		mr.Use(corsHandler)
		mr.Use(middleware.Logger)
		mr.Handle("/metrics", promhttp.Handler())

		ms := &http.Server{
			Handler: mr,
			Addr:    net.JoinHostPort(config.Metrics.Host, config.Metrics.Port),
		}
		log.Printf("Metrics server started on port %s", config.Metrics.Port)
		log.Fatal(ms.ListenAndServe())
	}()

	log.Printf("Server started on port %s", config.HTTP.Port)

	log.Fatal(s.ListenAndServe())
}
