package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/vocative-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllSurNames from db
func GetAllSurNames(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	SurNames := []model.SurName{}
	db.Order("count desc").Find(&SurNames)
	respondJSON(w, http.StatusOK, SurNames)
}

// GetSurName by name
func GetSurName(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["name"] == "" {
		respondError(w, http.StatusBadRequest, "")
	}

	SurName := getSurNameOr404(db, vars["name"], w, r)
	if SurName == nil {
		return
	}
	respondJSON(w, http.StatusOK, SurName)
}

// getCityOr404 gets a City instance if exists, or respond the 404 error otherwise
func getSurNameOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.SurName {

	SurNames := []model.SurName{}
	err := db.Where("name=?", name).Find(&SurNames).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &SurNames
}

// SearchSurNames by filter
func SearchSurNames(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	SurName := searchSurNameOr404(db, w, r)
	if SurName == nil {
		return
	}
	respondJSON(w, http.StatusOK, SurName)
}

// searchSurNameOr404 gets a City instance if exists, or respond the 404 error otherwise
func searchSurNameOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) *[]model.SurName {
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

	SurNames := []model.SurName{}
	err := tx.Order("count desc").Find(&SurNames).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &SurNames
}
