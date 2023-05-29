package routes

import "github.com/gin-gonic/gin"

type Routes struct {
	router *gin.Engine
}

func RegisterRoutes() Routes {
	r := Routes{
		router: gin.Default(),
	}

	v1 := r.router.Group("/v1")
	r.videoRoutes(v1)

	return r
}

func (r Routes) Run(addr ...string) error {
	return r.router.Run(addr...)
}
