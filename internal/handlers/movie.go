package handlers

import (
	"fmt"
	"net/http"

	movieService "myproject/challenge2015/internal/services/movie"

	"github.com/gin-gonic/gin"
)

func setMovieRoutes(router *gin.Engine) {
	router.GET("/smallest-degree-of-separation/:person1/:person2", GetMinimumDegreeOfSeperation)
}

func GetMinimumDegreeOfSeperation(c *gin.Context) {

	person1 := c.Param("person1")
	person2 := c.Param("person2")

	//validatation
	if !(len(person1) > 0 || len(person2) > 0) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Please Provide Both Person Name"})
	}
	service := movieService.New()
	seperation, err := service.GetMinimumDegreeOfSeperation(person1, person2)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("error in service layer %s", err.Error())})
	}

	message := fmt.Sprintf("Smallest Degree Of Separation between %s and %s is %d", person1, person2, seperation)

	c.IndentedJSON(http.StatusOK, gin.H{"message": message})
}
