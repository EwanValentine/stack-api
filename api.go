package api

import (
	"log"
	"net/http"
	"os"

	registrar "github.com/ewanvalentine/stack-registrar"
	"github.com/gorilla/mux"
)

type API interface {
	Init(registry registrar.Registry) *APIGateway
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
	HandlerFunc Handler
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
		api.AddRoute(route)
	}

	return api.router
}

// AddRoute - Add a route
func (api *APIGateway) AddRoute(route Route) {

	if api.router == nil {
		api.router = mux.NewRouter().StrictSlash(true)
	}

	handler := api.withContext(route.HandlerFunc)

	api.router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

func (api *APIGateway) Get(path string, handler Handler) {
	route := Route{
		Method:      "GET",
		Pattern:     path,
		HandlerFunc: handler,
	}
	api.AddRoute(route)
}

func (api *APIGateway) withContext(next Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		next(ctx)
	})
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
