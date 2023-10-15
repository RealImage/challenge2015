package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// New creates a router and registers all the routes for the
// service and returns it.
func New() http.Handler {
	router := gin.Default()
	setPingRoutes(router)
	setMovieRoutes(router)
	return router
}
