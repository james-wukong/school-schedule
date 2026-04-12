package http

import (
	"github.com/gin-gonic/gin"
	"github.com/james-wukong/school-schedule/internal/interface/http/middleware"
)

// RouterRegister is the interface a module must implement
type RouterRegister interface {
	Register(mw *middleware.Manager, router *gin.RouterGroup)
}

type Router struct {
	engine *gin.Engine
	mw     *middleware.Manager
}

func NewRouter(engine *gin.Engine, mw *middleware.Manager) *Router {
	return &Router{engine: engine, mw: mw}
}

// RegisterModules takes a slice of modules and lets them register themselves
func (r *Router) RegisterModules(v1Group *gin.RouterGroup, modules ...RouterRegister) {
	for _, m := range modules {
		m.Register(r.mw, v1Group)
	}
}
