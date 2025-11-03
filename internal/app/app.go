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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// metrics middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(metrics.HTTPMetricsMiddleware())

	// add middleware for validation (skip validation for metrics endpoint)
	r.Use(func(next http.Handler) http.Handler {
	    oapiValidator := oapi_middleware.OapiRequestValidator(swagger)
	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	        // Skip OpenAPI validation for /metrics endpoint
	        if r.URL.Path == "/metrics" {
	            next.ServeHTTP(w, r)
	            return
	        }
	        oapiValidator(next).ServeHTTP(w, r)
	    })
	})
	// r.Use(oapi_middleware.OapiRequestValidator(swagger))

	r.Handle("/metrics", promhttp.Handler())
	gen.HandlerFromMux(server, r)
	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort(config.Host, config.Port),
	}

	log.Printf("Server started on port %s", config.Port)

	log.Fatal(s.ListenAndServe())
}

