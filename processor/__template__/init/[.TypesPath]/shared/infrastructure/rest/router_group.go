package rest

import (
	"github.com/mercadolibre/fury_go-core/pkg/web"
	"github.com/mercadolibre/fury_go-platform/pkg/fury"
)

type SharedRouter interface {
	ReaderRoutes(app *fury.Application, prefixGroup *web.RouteGroup)
	WriterRoutes(app *fury.Application, prefixGroup *web.RouteGroup)
	WorkerRoutes(app *fury.Application, prefixGroup *web.RouteGroup)
}

type RouterGroup struct {
	Type       string
	PrefixPath string
	Routers    []SharedRouter
}

func NewRouterGroup(typ string, prefixPath string, routers ...SharedRouter) Router {
	return &RouterGroup{
		Type:       typ,
		PrefixPath: prefixPath,
		Routers:    routers,
	}
}

func NewReaderRouter(prefixPath string, routers ...SharedRouter) Router {
	return NewRouterGroup("reader", prefixPath, routers...)
}

func NewWriterRouter(prefixPath string, routers ...SharedRouter) Router {
	return NewRouterGroup("writer", prefixPath, routers...)
}

func NewWorkerRouter(prefixPath string, routers ...SharedRouter) Router {
	return NewRouterGroup("worker", prefixPath, routers...)
}

func (registrator *RouterGroup) RouteURLs(app *fury.Application) {
	prefixGroup := app.Router.Group(registrator.PrefixPath)

	switch registrator.Type {
	case "reader":
		for _, router := range registrator.Routers {
			router.ReaderRoutes(app, prefixGroup)
		}
	case "writer":
		for _, router := range registrator.Routers {
			router.WriterRoutes(app, prefixGroup)
		}
	default:
		for _, router := range registrator.Routers {
			router.WorkerRoutes(app, prefixGroup)
		}
	}
}

func (registrator *RouterGroup) API() *fury.Application {
	return API(registrator)
}
