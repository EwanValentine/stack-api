# Stack API

Stack API is an API gateway library, for quickly creating JSON API microservices. Stack API handles service registration, routing, and rending JSON.

## Example

```go
package main

import (
    "fmt"
    "log"

    api "github.com/ewanvalentine/stack-api"
    registrar "github.com/ewanvalentine/stack-registrar"
    "github.com/ewanvalentine/stack-registrar/services"
)

// Index handler
func Index(c *api.Context) {
    data := api.D{"test": "testing"}
    c.JSON(data, 200)
}

// Stuff
func Stuff(c *api.Context) {
    var data map[string]string
    c.Bind(&data)
    fmt.Println(api.D{"_message": "Got data"}, 200)
}

func main() {

    // Register your service
    registry := registrar.Init(registrar.SetHost("http://localhost:8080"))

    // Create new app
    app := api.Init(registry)
    err := app.Register(services.Service{
        Name:     "test",
        Host:     "test.com",
        Upstream: "http://localhost:8080",
    })

    // Some basic routes
    app.Get("/", Index)
    app.Post("/test", Stuff)

    if err != nil {
        log.Printf("Error creating service: %v", err)
    }

    // Run app
    app.Run(api.SetPort(":9090"))
}
```

## Self documenting API endpoints
Stack // API lets you generate endpoint documentation as your write your endpoints. 

To enable this feature, create a file `docs.html` in the root of your API codebase. Include a template within that file, such as... 

```
<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title></title>
  </head>
  <body>
    <h1>Docs</h1>

    <section>
      {{range .}}
        Name: {{.Name}}<br />
        Route: {{.Pattern}}<br />
        Method: {{.Method}}<br />
        <br />
      {{end}}
    </section>
  </body>
</html>
```

This allows you to style your docs according to your companies branding etc. 
