package main

import (
	"encoding/json"
	"golang-crowdfunding/handler"
	"golang-crowdfunding/user"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db, errMain = gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/golang_crowdfund?charset=utf8mb4&parseTime=True&loc=Local"))

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}


func main(){
	if errMain != nil {
		log.Fatal("Gagal koneksi ke database")
	}
	log.Info("Starting server")

	router := mux.NewRouter()
	

	handlerMain := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	userRepository :=  user.NewRepository(db)
	userService := user.NewService(userRepository)

	userHandler := handler.NewUserHandler(userService)

	router.HandleFunc("/users", getUsers)
	router.HandleFunc("/create-users", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/check-email", userHandler.CheckEmailAvailability).Methods("POST")

	http.ListenAndServe(":8000", handlerMain)
}

func getUsers(w http.ResponseWriter, r * http.Request){
	log.Info("fetch all users")
	var users []user.User
	db.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}