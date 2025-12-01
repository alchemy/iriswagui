package iriswagui

import (
	"embed"
	"net/http"
	"net/url"
	"os"
	"path"
	"text/template"

	"github.com/kataras/iris/v12"
)

//go:embed swagger-ui
var swaggerUI embed.FS

var indexTmpl *template.Template
var swaggerInitTmpl *template.Template

func init() {
	var err error
	indexTmpl, err = template.ParseFS(swaggerUI, "swagger-ui/index.html.tmpl")
	if err != nil {
		panic(err)
	}

	swaggerInitTmpl, err = template.ParseFS(swaggerUI, "swagger-ui/swagger-initializer.js.tmpl")
	if err != nil {
		panic(err)
	}
}

// SwaggerSpecRef represents a single swagger specification file.
// The URL field can be and absolute URL or a local path. In the
// letter case a route to serve that file will be added to the
// Swagger UI party.
type SwaggerSpecRef struct {
	Name string
	URL  string
}

// Config represents a (minimal) swagger UI configuration.
type Config struct {
	// ExternalConfigURL is the URL of an external configuration JavaScript object.
	// If it is non empty the UI configuration will be based on this object only, all
	// other parameters will be ingored.
	ExternalConfigURL string
	// SpecRefs is a list of API specifications [SwaggerSpecRef].
	// If only one SpecRef is present Swagger UI will be configured with the BasicLayout,
	// otherwise with the StandaloneLayout and the SwaggerUIBundle.plugins.Topbar plugin.
	SpecRefs []SwaggerSpecRef
	// DeepLinking will determine the value of the corresponfing parameter
	// of the Swagger UI configuration.
	DeepLinking bool
}

type swaggerUIConfig struct {
	Config
	BaseURL string
}

// HandleSwaggerUI adds a new sub-party to party with all the necessary endpoints
// to serve the Swagger UI based on config.
func HandleSwaggerUI(party iris.Party, relativePath string, config Config) error {
	swaggerParty := party.Party(relativePath)
	swaggerParty.HandleDir("/", swaggerUI, iris.DirOptions{IndexName: "/index.html", Compress: true})
	uiConfig := &swaggerUIConfig{
		Config:  config,
		BaseURL: swaggerParty.GetRelPath(),
	}

	specsParty := swaggerParty.Party("/specs")
	for i, specRef := range uiConfig.SpecRefs {
		parsedUrl, err := url.Parse(specRef.URL)
		if err != nil {
			return err
		}
		if !parsedUrl.IsAbs() {
			base := path.Base(parsedUrl.Path)
			relPath := "/" + base
			specsParty.Get(relPath, func(ctx iris.Context) {
				specContent, err := os.ReadFile(parsedUrl.Path)
				if err != nil {
					ctx.StopWithError(http.StatusInternalServerError, err)
				}
				ctx.Write(specContent)
			})
			specsPartyRelPath := specsParty.GetRelPath()
			uiConfig.SpecRefs[i].URL, err = url.JoinPath(specsPartyRelPath, relPath)
			if err != nil {
				return err
			}
		}
	}

	swaggerParty.Get("/", func(ctx iris.Context) {
		err := indexTmpl.ExecuteTemplate(ctx.ResponseWriter(), "index.html.tmpl", uiConfig)
		if err != nil {
			ctx.StopWithError(http.StatusInternalServerError, err)
		}
	})

	swaggerParty.Get("/swagger-initializer.js", func(ctx iris.Context) {
		err := swaggerInitTmpl.ExecuteTemplate(ctx.ResponseWriter(), "swagger-initializer.js.tmpl", uiConfig)
		if err != nil {
			ctx.StopWithError(http.StatusInternalServerError, err)
		}
	})

	return nil
}
