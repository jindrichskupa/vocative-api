package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/vocative-api/app/handler"
	"github.com/jindrichskupa/vocative-api/app/model"
	"github.com/jindrichskupa/vocative-api/config"
	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize application with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.Dialect,
		config.DB.Username,
		config.DB.Password,
		config.DB.Hostname,
		config.DB.Port,
		config.DB.Name)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	for tries := 0; tries < config.DB.Retries && err != nil; tries++ {
		time.Sleep(5 * time.Second)
		log.Println("Trying to reconnect (", tries, ") db: ", err)
		db, err = gorm.Open(config.DB.Dialect, dbURI)
	}
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}
	db.LogMode(true)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.DB.Prefix + defaultTableName
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
	log.Println("VOCATIVE API application started")
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/healtz", a.GetHealtStatus)

	a.Get("/vocative/firstnames", a.GetAllFirstNames)
	a.Get("/vocative/firstnames/search", a.SearchFirstNames)
	a.Get("/vocative/firstnames/{name}", a.GetFirstName)

	a.Get("/vocative/surnames", a.GetAllSurNames)
	a.Get("/vocative/surnames/search", a.SearchSurNames)
	a.Get("/vocative/surnames/{name}", a.GetSurName)

	a.Get("/vocative", a.GetVocative)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// GetHealtStatus retuns application status info
func (a *App) GetHealtStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetHealtStatus(a.DB, w, r)
}

// GetAllFirstNames handlers to get all Firstnames vocative
func (a *App) GetAllFirstNames(w http.ResponseWriter, r *http.Request) {
	handler.GetAllFirstNames(a.DB, w, r)
}

// GetFirstName handlers to get Firstname vocative
func (a *App) GetFirstName(w http.ResponseWriter, r *http.Request) {
	handler.GetFirstName(a.DB, w, r)
}

// SearchFirstNames handlers to search Firstnames vocative based on filter
func (a *App) SearchFirstNames(w http.ResponseWriter, r *http.Request) {
	handler.SearchFirstNames(a.DB, w, r)
}

// GetAllSurNames handlers to get all Surnames vocative
func (a *App) GetAllSurNames(w http.ResponseWriter, r *http.Request) {
	handler.GetAllSurNames(a.DB, w, r)
}

// GetSurName handlers to get Surname vocative
func (a *App) GetSurName(w http.ResponseWriter, r *http.Request) {
	handler.GetSurName(a.DB, w, r)
}

// SearchSurNames handlers to search Surnames vocative based on filter
func (a *App) SearchSurNames(w http.ResponseWriter, r *http.Request) {
	handler.SearchSurNames(a.DB, w, r)
}

// SearchSurNames handlers to search Surnames vocative based on filter
func (a *App) GetVocative(w http.ResponseWriter, r *http.Request) {
	handler.GetVocative(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
