package main

import (
	"encoding/json"
	"fmt"
	"golang-crowdfunding/auth"
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/handler"
	"golang-crowdfunding/helper"
	"golang-crowdfunding/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
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
	campaignRepository := campaign.NewRepository(db)
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router.HandleFunc("/users", getUsers)
	router.HandleFunc("/create-user", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/check-email", userHandler.CheckEmailAvailability).Methods("POST")

	router.HandleFunc("/campaign", campaignHandler.GetCampaigns)
	router.HandleFunc("/campaign", campaignHandler.GetCampaigns).Queries("userId", "{userId}").Name("campaignSingular")
	router.HandleFunc("/campaign/{campaign_id}", campaignHandler.GetCampaignById)

	router.Handle("/upload-avatar", authMiddleware(http.HandlerFunc(userHandler.UploadAvatar), userService, authService)).Methods("POST")

	router.Handle("/create-campaign", authMiddleware(http.HandlerFunc(campaignHandler.CreateCampaign),userService, authService)).Methods("POST")
	router.Handle("/update-campaign/{campaign_id}", authMiddleware(http.HandlerFunc(campaignHandler.UpdateCampaign),userService, authService)).Methods("POST")

	http.ListenAndServe("127.0.0.1:8000", handlerMain)
}

func getUsers(w http.ResponseWriter, r * http.Request){
	log.Info("fetch all users")
	var users []user.User
	db.Find(&users)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func authMiddleware(h http.Handler, userService user.Service, authService auth.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")

	var response helper.Response
	
	if !strings.Contains(authHeader, "Bearer"){
		response = helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenString := ""
	arrayToken := strings.Split(authHeader, " ")
	fmt.Println(arrayToken)
	if len(arrayToken) == 2 {
		tokenString = arrayToken[1]
	}
	token, err := authService.ValidateToken(tokenString)
	if err != nil {
		response = helper.ApiResponse("Invalid token", http.StatusUnauthorized, "error", nil)
		fmt.Println(response)
		json.NewEncoder(w).Encode(response)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		response = helper.ApiResponse("Invalid token", http.StatusUnauthorized, "error", nil)
		json.NewEncoder(w).Encode(response)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userId := int(claims["user_id"].(float64))

	user, err := userService.GetUserById(userId)
	if err != nil {
		response = helper.ApiResponse("Invalid token", http.StatusUnauthorized, "error", nil)
		json.NewEncoder(w).Encode(response)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userIdConv := strconv.Itoa(user.Id)

	w.Header().Set("currentUser", userIdConv)
	w.Header().Set("KONTOL", "GEDE")
	

	h.ServeHTTP(w, r)
	})
	
}
