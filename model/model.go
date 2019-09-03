package model

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // blank import required by gorm
)

func init() {
	//open a db connection
	var err error
	DB, err = gorm.Open("mysql", "root:root@/pizza_golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	//Migrate the schema
	DB.AutoMigrate(&Store{}, &Recipe{}, &Resource{})

	// load or create Store
	S = &Store{}
	DB.Where("id = ?", 1).First(S)
	DB.Save(S)
}

type (
	// Store is the databse model for a pizza-services ingredient-storage
	Store struct {
		gorm.Model
		Items []*Resource `json:"items"`
	}

	// StoreDescription is a model for the default store-model presented via the REST-api
	StoreDescription struct {
		Items []*ResourceDescription `json:"items"`
	}
)

type (
	// Recipe is the database model for a pizza-recipe
	Recipe struct {
		gorm.Model
		Title     string      `json:"title"`
		Resources []*Resource `json:"resources"`
	}

	// RecipeDescription is a model for the default recipe-model presented via the REST-api
	RecipeDescription struct {
		ID        uint
		Title     string                 `json:"title"`
		Resources []*ResourceDescription `json:"resources"`
	}
)

type (
	// Resource is the database model for a pizza-ingredient
	Resource struct {
		gorm.Model
		Name     string  `json:"name"`
		Amount   float64 `json:"amount"`
		RecipeID uint
		StoreID  uint
	}

	// ResourceDescription is a model for the default resource-model presented via the REST-api
	ResourceDescription struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}
)

// DB the underlying database
var DB *gorm.DB

// S the central ingredient-storage
var S *Store
var sM = &sync.Mutex{}

// Description generates a StoreDescription from the Store
func (s *Store) Description() *StoreDescription {
	items := s.List()
	itemDesc := make([]*ResourceDescription, len(items))
	for i, it := range items {
		itemDesc[i] = it.Description()
	}
	return &StoreDescription{
		Items: itemDesc,
	}
}

// List retrives all ingredients related to this store from the database
func (s *Store) List() (items []*Resource) {
	DB.Model(S).Related(&items)
	return
}

// Add items to the store
func (s *Store) Add(items ...*Resource) {
	for i := range items {
		items[i].ID = 0
	}
	sM.Lock()
	s.Items = append(s.List(), items...)
	DB.Save(s)
	sM.Unlock()
}

// Remove items from the store. If there is not enough of one ingredient
// available, the method returns an error
func (s *Store) Remove(items ...*Resource) error {
	log.Println(items)
	changed := make([]*Resource, 0)
	sM.Lock()
	defer sM.Unlock()
	for _, r := range items {
		var available []*Resource
		DB.Model(S).Where("name = ?", r.Name).Related(&available)
		log.Println("amount of", r.Name, ":", len(available))
		i := 0
		for i < len(available) && r.Amount > 0 {
			min := math.Min(available[i].Amount, r.Amount)
			available[i].Amount -= min
			fmt.Println("amount1", available[i].Amount)
			r.Amount -= min
			changed = append(changed, available[i])
			fmt.Println("amount2", available[i].Amount)
			i++
		}
		log.Println(r.Name, r.Amount)
		if r.Amount > 0 {
			return errors.New("not enough " + r.Name + " available")
		}
	}
	log.Println(len(changed))
	for i := range changed {
		log.Println(changed[i])
		fmt.Println("amount3", changed[i].Amount)
		DB.Save(changed[i])
		if changed[i].Amount <= 0 {
			DB.Delete(changed[i])
		}
	}

	return nil
}

// Description generates a RecipeDescription from the Recipe
func (r *Recipe) Description() *RecipeDescription {
	resourceDesc := make([]*ResourceDescription, len(r.Resources))
	for i, r := range r.Resources {
		resourceDesc[i] = r.Description()
	}

	return &RecipeDescription{
		ID:        r.ID,
		Title:     r.Title,
		Resources: resourceDesc,
	}
}

// Create db operation
func (r *Recipe) Create() {
	DB.Save(r)
}

// Read db operation
func (r *Recipe) Read() error {
	DB.Where("id = ?", r.ID).First(&r)
	if r.ID == 0 {
		return errors.New("recipe with id " + strconv.FormatUint(uint64(r.ID), 10) + " not found")
	}
	DB.Model(&r).Related(&r.Resources)
	return nil
}

// Update db operation
func (r *Recipe) Update() error {
	var duplicate Recipe
	DB.Where("id = ?", r.ID).First(&duplicate)
	if duplicate.Title == "" {
		return errors.New("there is no recipe with id " + strconv.Itoa(int(r.ID)))
	}
	DB.Model(&duplicate).Related(&duplicate.Resources)
	for _, r := range duplicate.Resources {
		r.delete()
	}

	DB.Save(&r)
	return nil
}

// Delete db operation
func (r *Recipe) Delete() error {
	var recipe Recipe
	DB.Where("id = ?", r.ID).First(&recipe)
	if recipe.ID == 0 {
		return errors.New("recipe with id " + strconv.FormatUint(uint64(recipe.ID), 10) + " did not exist")
	}

	DB.Delete(&recipe)
	return nil
}

// Description generates a ResourceDescription from the Resource
func (r *Resource) Description() *ResourceDescription {
	return &ResourceDescription{
		Name:   r.Name,
		Amount: r.Amount,
	}
}

func (r *Resource) delete() error {
	var res Resource
	DB.Where("id = ?", r.ID).First(&res)
	if res.ID == 0 {
		return errors.New("recipe with id " + strconv.FormatUint(uint64(res.ID), 10) + " did not exist")
	}

	DB.Delete(r)
	return nil
}
