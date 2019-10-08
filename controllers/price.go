package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theMomax/GolangPizzaServiceManager/model"
)

const fixed = "fixed"
const normal = "normal"

func init() {
	profiles := os.Getenv("PRICE")

	if gin.Mode() == gin.ReleaseMode && strings.Contains(profiles, fixed) {
		pc = &fixedPriceCalculator{}
	} else if gin.Mode() == gin.ReleaseMode && strings.Contains(profiles, normal) {
		pc = &normalPriceCalculator{}
	} else {
		pc = &mockPriceCalculator{}
	}
}

var pc priceCalculator

// Price -> TODO: api-link
func Price(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id-format: " + c.Param("id")})
		return
	}
	price, err := pc.priceOf(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id " + c.Param("id")})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": price})
}

type priceCalculator interface {
	priceOf(recipeID uint) (float64, error)
}

type mockPriceCalculator struct{}

func (m *mockPriceCalculator) priceOf(recipeID uint) (float64, error) {
	return 0.0, nil
}

type fixedPriceCalculator struct{}

func (f *fixedPriceCalculator) priceOf(recipeID uint) (float64, error) {
	return 7.0, nil
}

type normalPriceCalculator struct{}

func (n *normalPriceCalculator) priceOf(recipeID uint) (float64, error) {
	var recipe = model.Recipe{}
	recipe.ID = recipeID
	err := recipe.Read()
	if err != nil {
		return 0.0, err
	}

	var price float64
	for _, r := range recipe.Resources {
		price += r.Amount
	}

	return price, nil
}
