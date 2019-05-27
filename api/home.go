package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Homepage",
	})
}
