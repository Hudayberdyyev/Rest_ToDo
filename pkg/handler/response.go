package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	ctx.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
