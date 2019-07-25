package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/vocative-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllFirstNames from db
func GetAllFirstNames(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	FirstNames := []model.FirstName{}
	db.Find(&FirstNames)
	respondJSON(w, http.StatusOK, FirstNames)
}

// GetFirstName by name
func GetFirstName(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["name"] == "" {
		respondError(w, http.StatusBadRequest, "")
	}

	FirstName := getFirstNameOr404(db, vars["name"], w, r)
	if FirstName == nil {
		return
	}
	respondJSON(w, http.StatusOK, FirstName)
}

// getCityOr404 gets a City instance if exists, or respond the 404 error otherwise
func getFirstNameOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.FirstName {

	FirstNames := []model.FirstName{}
	err := db.Where("name=?", name).Find(&FirstNames).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &FirstNames
}

// SearchFirstNames by filter
func SearchFirstNames(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	FirstName := searchFirstNameOr404(db, w, r)
	if FirstName == nil {
		return
	}
	respondJSON(w, http.StatusOK, FirstName)
}

// searchFirstNameOr404 gets a City instance if exists, or respond the 404 error otherwise
func searchFirstNameOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) *[]model.FirstName {
	queryValues := r.URL.Query()

	log.Println(queryValues)
	filters := []string{"name", "gender"}
	column := map[string]string{
		"name":   "name_search like ?",
		"gender": "gender = ?",
	}

	tx := db.Where("1 = 1")

	for _, filter := range filters {
		if queryValues[filter] != nil {
			value := normalizeNameSearch(queryValues[filter][0])
			if filter == "name" {
				value = "%" + value + "%"
			}
			log.Println("Key [", filter, "]: ", value)
			tx = tx.Where(column[filter], value)
		}
	}

	if queryValues["limit"] != nil {
		limit, err := strconv.Atoi(queryValues["limit"][0])
		if err != nil {
			limit = 200
		}
		tx = tx.Limit(limit)
	} else {
		tx = tx.Limit(200)
	}

	FirstNames := []model.FirstName{}
	err := tx.Order("count desc").Find(&FirstNames).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &FirstNames
}
