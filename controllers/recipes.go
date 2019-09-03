package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/theMomax/GolangPizzaServiceManager/model"
)

// CreateRecipe -> TODO: api-link
func CreateRecipe(c *gin.Context) {
	var recipe model.Recipe
	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal format for recipe: " + err.Error()})
		return
	}

	if recipe.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"status": "cannot update/create recipe with specific id"})
		return
	}

	if recipe.Title == "" || recipe.Resources == nil || len(recipe.Resources) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content for recipe"})
		return
	}

	var duplicate model.Recipe
	model.DB.Where("title = ?", recipe.Title).First(&duplicate)
	if duplicate.Title != "" {
		fmt.Println(duplicate)
		c.JSON(http.StatusConflict, gin.H{"status": "recipe with title " + duplicate.Title + " already exists"})
		return
	}

	recipe.Create()
	c.JSON(http.StatusCreated, gin.H{"resourceID": recipe.ID})
}

// UpdateRecipe -> TODO: api-link
func UpdateRecipe(c *gin.Context) {
	var recipe model.Recipe
	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal format for recipe: " + err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "invalid id-format: " + c.Param("id")})
		return
	}
	recipe.ID = uint(id)

	if recipe.ID == 0 {
		c.JSON(http.StatusConflict, gin.H{"status": "cannot update recipe without id"})
		return
	}

	if recipe.Title == "" || recipe.Resources == nil || len(recipe.Resources) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "illegal content for recipe"})
		return
	}

	err = recipe.Update()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resourceID": recipe.ID})
}

// Fetch -> TODO: api-link
func Fetch(c *gin.Context) {
	var recipe model.Recipe
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "invalid id-format: " + c.Param("id")})
		return
	}
	recipe.ID = uint(id)
	err = recipe.Read()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "found", "data": recipe.Description()})
}

// FetchAll -> TODO: api-link
func FetchAll(c *gin.Context) {
	var recipies []model.Recipe
	model.DB.Find(&recipies)

	if len(recipies) <= 0 {
		c.JSON(http.StatusNoContent, gin.H{"status": "no recipes found"})
		return
	}

	var descriptions []*model.RecipeDescription
	for _, r := range recipies {
		model.DB.Model(&r).Related(&r.Resources)
		descriptions = append(descriptions, r.Description())
	}

	c.JSON(http.StatusOK, gin.H{"status": strconv.Itoa(len(descriptions)) + " recipes found", "data": descriptions})
}

// DeleteRecipe -> TODO: api-link
func DeleteRecipe(c *gin.Context) {
	var recipe model.Recipe
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "invalid id-format: " + c.Param("id")})
		return
	}
	recipe.ID = uint(id)
	err = recipe.Delete()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "recipe was deleted succesfully"})
}
