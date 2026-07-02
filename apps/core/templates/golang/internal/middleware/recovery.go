package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {

	return gin.CustomRecovery(func(c *gin.Context, err any) {

		log.Error(
			"panic",

			slog.Any("error", err),
		)

		c.JSON(500, gin.H{
			"message": "internal server error",
		})

	})

}
