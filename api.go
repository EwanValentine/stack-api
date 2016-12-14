package api

import (
	"log"
	"net/http"
	"os"

	registrar "github.com/ewanvalentine/stack-registrar"
	"github.com/ewanvalentine/stack-registrar/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type API interface {
	Register(service services.Service) error
	Run(options ...ApiOptions)
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
func (api *APIGateway) Register(service services.Service) error {
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
	route := api.buildRoute("GET", path, handler)
	api.AddRoute(route)
}

func (api *APIGateway) Post(path string, handler Handler) {
	route := api.buildRoute("POST", path, handler)
	api.AddRoute(route)
}

func (api *APIGateway) Patch(path string, handler Handler) {
	route := api.buildRoute("PATCH", path, handler)
	api.AddRoute(route)
}

func (api *APIGateway) Put(path string, handler Handler) {
	route := api.buildRoute("PUT", path, handler)
	api.AddRoute(route)
}

func (api *APIGateway) Delete(path string, handler Handler) {
	route := api.buildRoute("DELETE", path, handler)
	api.AddRoute(route)
}

func (api *APIGateway) buildRoute(method, path string, handler Handler) Route {
	return Route{
		Method:      method,
		Pattern:     path,
		HandlerFunc: handler,
	}
}

// withContext - Wraps a context object around each request
func (api *APIGateway) withContext(next Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		next(ctx)
	})
}

// ApiOption - Api option object
type ApiOption struct {
	port string
}

// ApiOptions - Option type
type ApiOptions func(*ApiOption) error

// SetPort - Sets API port
func SetPort(port string) ApiOptions {
	return func(opt *ApiOption) error {
		opt.port = port
		return nil
	}
}

// Run - Run api
func (api *APIGateway) Run(options ...ApiOptions) {

	opt := &ApiOption{}

	for _, op := range options {
		err := op(opt)
		if err != nil {
			log.Fatalf("Error rending configuration: %v", err)
		}
	}

	port := ":8080"
	envPort := os.Getenv("STACK_SERVICE_PORT")

	if envPort != "" {
		port = envPort
	} else if opt.port != "" {
		port = opt.port
	}

	log.Printf("Running on port: %v", port)

	// @todo - make this customisable
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	})

	handler := c.Handler(api.router)

	// Run
	log.Fatal(http.ListenAndServe(port, handler))
}
