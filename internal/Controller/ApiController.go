package Controller

import (
	"GoCacher/internal"
	"GoCacher/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() {

}

func Insert(c *gin.Context) {

	coordinator := internal.GetCoordinator()

	var data model.InsertRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance := coordinator.GetInstance(data.Identifier)
	instance.SetValue(data.Key, data.Value)

}

func Fetch(c *gin.Context) {
	coordinator := internal.GetCoordinator()

	var data model.FetchRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instance := coordinator.GetInstance(data.Identifier)
	value, _ := instance.GetValue(data.Key)
	response := model.FetchResponse{Key: data.Key, Value: value}

	c.JSON(200, response)
}
