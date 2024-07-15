iriswagui
=========

This package allows the user to easily add a [Swagger UI](https://github.com/swagger-api/swagger-ui) to an [Iris](https://www.iris-go.com/) application.
It does so by embedding the Swagger UI assets and adding a party to a parent party that implements all the necessary endpoints.

Basic usage example:
```go
package main

import (
	"flag"
	"log"

	"github.com/alchemy/iriswagui"
	"github.com/kataras/iris/v12"
)

func main() {
	specPath := flag.String("spec", "data/swagger.json", "specification file `path`")
	addr := flag.String("addr", "127.0.0.1:8888", "listening address `host:port`")

	app := iris.New()
	party := app.Party("/documentation")

	config := iriswagui.Config{
		SpecRefs: []iriswagui.SwaggerSpecRef{
			{Name: "MyAPI", URL: *specPath},
		},
	}
    // Swagger UI will be served from /documentation/swagger
	err := iriswagui.HandleSwaggerUI(party, "/swagger", config)
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Listen(*addr)
	if err != nil {
		log.Fatalln(err)
	}
}
```

Swagger UI configuration
------------------------
Swagger UI is highly customizable. This package Config struct offers only minimal configuration options.  
For instance, it accepts only a list of specification files and wether to use deep linking or not.  
If only one specification file is added, Swagger UI will be configured to use `BasicLayout`, otherwise `StandaloneLayout` will be used along with the `SwaggerUIBundle.plugins.Topbar` plugin.  
However, a complete custom configration can be specified by using the `ExternalConfigURL` field of the `Config` struct. When this field is not empty
the configuration object will be fetched from this URL, ignoring all other fields.
