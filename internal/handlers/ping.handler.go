package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

type PingResponse struct {
	Message     string `example:"pong"`
	RequestId   string `example:"123"`
	ContentType string `example:"application/json"`
}

// GetPing godoc
// @Summary Get Ping
// @Description Returns pong response
// @Tags ping
// @Produce json
// @Param X-Request-ID header int true "Header for requestID"
// @Param Content-Type header string true "Header for body type"
// @Success 200 {object} PingResponse
// @Router /ping [get]
func (p *PingHandler) GetPing(ctx *gin.Context) {
	requestId := ctx.GetHeader("X-Request-ID")
	contentType := ctx.GetHeader("Content-Type")
	ctx.JSON(http.StatusOK, PingResponse{
		Message:     "oke",
		RequestId:   requestId,
		ContentType: contentType,
	})
}
