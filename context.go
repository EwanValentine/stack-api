package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Context - Http Context object.
type Context struct {

	// Database connection
	Datastore *Datastore

	// Http data
	request  *http.Request
	response http.ResponseWriter

	// Mux params
	Params map[string]string
}

// NewContext - Return a new context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		request:  r,
		response: w,
	}
	ctx.Init()
	return ctx
}

// Init - Initialise http context.
func (c *Context) Init() {
	c.Params = mux.Vars(c.request)
}

// JSON - Set header to JSON
func (c *Context) JSON(data interface{}, code int) {
	c.SetHttpCode(code)
	c.response.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)

	if err != nil {
		http.Error(c.response, err.Error(), http.StatusInternalServerError)
		return
	}

	c.response.Write(json)
}

// Set - Sets http response header code.
func (c *Context) SetHttpCode(code int) {
	c.response.WriteHeader(code)
}
