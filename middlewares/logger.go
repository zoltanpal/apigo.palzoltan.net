package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorLogger prints out any errors that handlers push to the Gin Context.
func ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute handler

		for _, ginErr := range c.Errors {
			// You can inspect ginErr.Type, .Meta, etc.
			log.Printf("[GIN ERROR] %v\n", ginErr.Err)
		}
	}
}
