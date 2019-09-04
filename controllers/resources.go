package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/theMomax/GolangPizzaServiceManager/model"
)

// Order -> TODO: api-link
func Order(c *gin.Context) {
	var recipe model.Recipe
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id-format: " + c.Param("id")})
		return
	}
	recipe.ID = uint(id)
	err = recipe.Read()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = model.S.Remove(recipe.Resources...)
	if err != nil {
		c.JSON(http.StatusIMUsed, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ordered " + recipe.Title})
}

// FetchAvailable -> TODO: api-link
func FetchAvailable(c *gin.Context) {
	items := model.S.Description().Items
	if len(items) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "no items found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": strconv.Itoa(len(items)) + " resources found", "data": items})
}

// AddResource -> TODO: api-link
func AddResource(c *gin.Context) {
	var items []*model.Resource
	if err := c.BindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "illegal format for resources: " + err.Error()})
		return
	}

	model.S.Add(items...)
	c.JSON(http.StatusOK, gin.H{"message": strconv.Itoa(len(items)) + " resources added to store"})
}
