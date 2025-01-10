package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler) {
	r.POST("/orders", handler.Create)
	r.DELETE("/orders/:id", handler.Cancel)
}
