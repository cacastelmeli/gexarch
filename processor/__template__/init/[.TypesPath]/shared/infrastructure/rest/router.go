package rest

import (
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
)

type Router interface {
	API() *fury.Application
	RouteURLs(app *fury.Application)
}

// API constructs an http.Handler with all application routes defined.
func API(handler Router) *fury.Application {
	app, err := fury.NewWebApplication()
	if err != nil {
		panic(err.Error())
	}

	handler.RouteURLs(app)
	return app
}
