# Stack API

Disclaimer: This project started purely to satisy requirements specific to my preferences. Though, feel free to submit a PR/feedback!  

Stack API is an API gateway library, for quickly creating JSON API microservices. Stack API handles service registration, routing, and rending JSON.

## Project goals
Create API's and microservices with ease, and attempt to solve common headaches when developing API's and microservices. TThis project is heavily inspired by Go micro.

The reasons I'm not just using Go micro, Go micro didn't have a provider for Kong and I had difficulty creating one. Also, as a microservice newbie, I was a little bewildered by Go micro, so wanted to create something barebones for my own understanding. 

Finally, as someone who has used several different Go frameworks in production, I also wanted to take parts I liked of each, and babake them into this project.

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
    data := api.D{"message": "Hello world!"}
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

## Service discovery 
Stack can create a registry instance, using the stack-registrar library. Registrar, accepts instances of service discovery providers and allows you to pass metadata describing your apis and services to your discovery provider. 

The default provider is Kong... 

```go
// Register your service
registry := registrar.Init(registrar.SetHost("http://localhost:8080"))
```

If you wanted to use Consul instead... (warning, this is still under development. See the following interface for how to create additional providers, please submit a PR! https://github.com/EwanValentine/stack-registrar/blob/master/providers/provider.go)

```go
// Create a new instance of the Consul provider
consul := providers.Consul("http://consul")

// Register your service
registry := registrar.Init(registrar.SetProvider(consul))<Paste>
```

## Self documenting API endpoints
Stack // API lets you generate endpoint documentation as your write your endpoints. 

To enable this feature, create a file `docs.html` in the root of your API codebase. Include a template within that file, such as... 

```html
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

## CORS
Stack handles CORS out-of-the-box. The default rules are... 

```
c := cors.New(cors.Options{
   AllowedHeaders: []string{"*"},
   AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},       
})
```

## Todo
1. Make CORS rules customisable.
2. Create new stack-registrar providers (Consul, etcd, zookeeper etc). 
3. Make docs generation more customisable.
