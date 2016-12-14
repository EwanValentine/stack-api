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
	Request  *http.Request
	Response http.ResponseWriter

	// Mux params
	Params map[string]string
}

// NewContext - Return a new context.
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Request:  r,
		Response: w,
	}
	ctx.Init()
	return ctx
}

// Init - Initialise http context.
func (c *Context) Init() {
	c.Params = mux.Vars(c.Request)
}

// Bind - Bind's json data to empty struct.
func (c *Context) Bind(payload interface{}) error {

	// Decode body string into a string
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		return err
	}

	defer c.Request.Body.Close()

	return nil
}

// JSON - Set header to JSON
func (c *Context) JSON(data interface{}, code int) {
	c.SetHttpCode(code)
	c.Response.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)

	if err != nil {
		http.Error(c.Response, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Response.Write(json)
}

// Set - Sets http Response header code.
func (c *Context) SetHttpCode(code int) {
	c.Response.WriteHeader(code)
}

// Param - Get param by name
func (c *Context) Param(name string) string {
	return c.Params[name]
}

// Header - Get a header value
func (c *Context) Header(name string) string {
	return c.Request.Header.Get(name)
}

// Param - Get param by name
func (c *Context) Param(name string) string {
	return c.Params[name]
}

// Header - Get a header value
func (c *Context) Header(name string) string {
	return c.request.Header.Get(name)
}
