package api

import (
	"log"
	"net/http"
	"os"

	registrar "github.com/ewanvalentine/stack-registrar"
	"github.com/gorilla/mux"
)

type API interface {
	Init()
	RegisterRoutes(routes Routes)
	Run()
}

type APIGateway struct {
	router   *mux.Router
	registry registrar.Registry
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func Init(registry registrar.Registry) *APIGateway {
	return &APIGateway{nil, registry}
}

// Register - Register your API
func (api *APIGateway) Register(service registrar.Service) error {
	return api.registry.Register(service)
}

func (api *APIGateway) RegisterRoutes(routes Routes) *mux.Router {
	if api.router == nil {
		api.router = mux.NewRouter().StrictSlash(true)
	}

	for _, route := range routes {
		api.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return api.router
}

func (api *APIGateway) Run() {
	var port string
	port = os.Getenv("STACK_API_PORT")
	if port == "" {
		port = ":" + "8080"
	}

	// Run
	log.Fatal(http.ListenAndServe(os.Getenv("STACK_API_PORT"), api.router))
}
