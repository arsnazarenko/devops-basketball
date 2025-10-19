package app

import (
	"log"
	"net"
	"net/http"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/config"
	v1 "github.com/arsnazarenko/devops-basketball/internal/controller/http/v1"
	"github.com/arsnazarenko/devops-basketball/internal/usecase"
	"github.com/arsnazarenko/devops-basketball/internal/usecase/repo"
	"github.com/arsnazarenko/devops-basketball/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/oapi-codegen/nethttp-middleware"
)

func Run() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	swagger, err := gen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec\n: %s", err)
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
	// add middleware for validation
	r.Use(nethttpmiddleware.OapiRequestValidator(swagger))
	gen.HandlerFromMux(server, r)
	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort(config.Host, config.Port),
	}

	log.Fatal(s.ListenAndServe())
}
