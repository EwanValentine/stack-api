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

// Bind - Bind's json data to empty struct.
func (c *Context) Bind(payload interface{}) error {

	// Decode body string into a string
	decoder := json.NewDecoder(c.request.Body)
	err := decoder.Decode(&payload)

	if err != nil {
		return err
	}

	defer c.request.Body.Close()

	return nil
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

// Param - Get param by name
func (c *Context) Param(name string) string {
	return c.Params[name]
}

// Header - Get a header value
func (c *Context) Header(name string) string {
	return c.request.Header.Get(name)
}
