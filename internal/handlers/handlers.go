package handlers

import (
	"github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/Vallghall/wb-test-l0/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *services.Service
}

func New(s *services.Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Routes() *gin.Engine {
	g := gin.Default()

	g.GET("/favicon.ico")

	g.GET("/msg/:id", func(c *gin.Context) {
		msg, err := h.service.LoadCachedMsgById(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, map[string]message.Message{
			"response": *msg,
		})
	})

	return g
}
