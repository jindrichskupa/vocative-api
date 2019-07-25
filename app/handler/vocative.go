package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jindrichskupa/vocative-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetVocative by filter
func GetVocative(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Vocatives := getVocativeOr404(db, w, r)
	if Vocatives == nil {
		return
	}
	respondJSON(w, http.StatusOK, Vocatives)
}

// searchSurNameOr404 gets a City instance if exists, or respond the 404 error otherwise
func getVocativeOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) *[]model.VocativeJSON {
	queryValues := r.URL.Query()

	log.Println(queryValues)
	filters := map[string]string{
		"firstname": "",
		"surname":   "",
		"gender":    "",
	}
	column := map[string]string{
		"firstname": "name_search like ?",
		"surname":   "name_search like ?",
		"gender":    "gender = ?",
	}

	tx := db.Where("1 = 1")

	for key, filter := range filters {

		if queryValues[key] != nil {
			value := normalizeNameSearch(queryValues[key][0])
			if key == "firstname" || key == "surname" {
				value = "%" + value + "%"
			}
			log.Println("Key [", key, filter, "]: ", value)
			filters[key] = value
		}
	}

	if queryValues["limit"] != nil {
		limit, err := strconv.Atoi(queryValues["limit"][0])
		if err != nil {
			limit = 10
		}
		tx = tx.Limit(limit)
	} else {
		tx = tx.Limit(10)
	}
	tx = tx.Order("count desc")
	tx = tx.Where(column["gender"], filters["gender"])

	FirstNames := []model.FirstName{}
	err := tx.Where(column["firstname"], filters["firstname"]).Find(&FirstNames).Error
	SurNames := []model.SurName{}
	err = tx.Where(column["surname"], filters["surname"]).Find(&SurNames).Error

	Vocatives := []model.VocativeJSON{}

	for _, firstname := range FirstNames {
		for _, surname := range SurNames {
			Vocatives = append(Vocatives, model.VocativeJSON{
				Name:     firstname.Name + " " + surname.Name,
				Vocative: firstname.Vocative + " " + surname.Vocative,
				Count:    firstname.Count + surname.Count,
				Gender:   firstname.Gender,
			})
		}
	}

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &Vocatives
}
