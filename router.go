package foundation

import "github.com/gin-gonic/gin"

var (
	router *gin.Engine
)

func NewGinEngine() *gin.Engine {
	router = gin.New()

	return router
}
