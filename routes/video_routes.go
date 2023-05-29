package routes

import (
	"encoder/controllers"
	"github.com/gin-gonic/gin"
)

func (r Routes) videoRoutes(rg *gin.RouterGroup) {
	vg := rg.Group("/videos")
	vg.POST("/encode", controllers.EncodeVideo)
}
