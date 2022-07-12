package handlers

import (
	"github.com/Vallghall/wb-test-l0/internal/model/message"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) fetchMessageById(c *gin.Context) {
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
}
