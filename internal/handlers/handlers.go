package handlers

import (
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

	g.Static("/js", "./web/static/js")
	g.Static("/css", "./web/static/css")
	g.StaticFile("/favicon.ico", "./web/static/images/favicon.ico")
	g.LoadHTMLFiles("./web/templates/index.gohtml")

	g.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.gohtml", nil)
	})
	g.GET("/msg/:id", h.fetchMessageById)

	return g
}
