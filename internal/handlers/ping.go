package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setPingRoutes(router *gin.Engine) {
	router.GET("/ping", Ping)
	router.GET("/albums", getAlbums)
}

func Ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "ping")
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
