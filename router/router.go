package router

import (
	"net/http"
	"time"

	"github.com/Weeping-Willow/api-example/service"
	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type router struct {
	engine  *gin.Engine
	Service service.Service
}

func StartServer(s service.Service) error {
	r := newRouter(s).initCors().initRoutes()

	server := &http.Server{
		Addr:         ":" + s.GetConfig().Port,
		Handler:      r.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return server.ListenAndServe()
}

func newRouter(s service.Service) *router {
	e := gin.New()
	e.Use(ginzerolog.Logger("gin"), gin.Recovery())

	return &router{
		engine:  e,
		Service: s,
	}
}

func (r *router) initRoutes() *router {
	r.engine.GET("/", r.greet)

	api := r.engine.Group("/api")
	// routes under /api will be protected by token auth
	api.Use(withError(middlewareTokenAuth(r.Service.TokenService())))

	api.GET("/", r.greet)
	api.POST("/score", withError(controllerPostScore(r.Service.ScoreService())))
	return r
}

func (r *router) initCors() *router {
	conf := cors.DefaultConfig()
	conf.AllowAllOrigins = true
	conf.AllowCredentials = true
	conf.AddAllowHeaders("Authorization")
	r.engine.Use(cors.New(conf))

	return r
}

func (r *router) greet(c *gin.Context) {
	c.JSON(200, gin.H{"greet": "Hello World"})
}
