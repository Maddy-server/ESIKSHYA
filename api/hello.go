package api
import (
	"github.com/gin-gonic/gin"
)
 func (server *Server) Hello(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello",
	})
 }
