package handler

import (
	"encoding/json"
	"fmt"
	"golang-crowdfunding/auth"
	"golang-crowdfunding/helper"
	"golang-crowdfunding/user"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator"
)


type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input user.RegisterUserInput
	log.Info("User registration called")
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Invalid request body")
		return
	}
	json.Unmarshal(reqBody, &input)

	errValidate := validatorCustom(input)
	if errValidate != nil {
		response := helper.ApiResponse("Invalid input", http.StatusUnprocessableEntity, "failed", errValidate)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Can't react database, please try again later")
		return
	}

	token, err := h.authService.GenerateToken(newUser.Id)
	if err != nil {
		response := helper.ApiResponse("Invalid input", http.StatusUnprocessableEntity, "failed", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	formatter := user.FormatUser(newUser, token)
	
	response := helper.ApiResponse("User created successfully", http.StatusOK, "success", formatter)
	
	
	json.NewEncoder(w).Encode(response)
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input user.LoginInput
	log.Info("Login endpoint called")
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Fatal("Can't react database, please try again later")
		return
	}
	json.Unmarshal(reqBody, &input)
	// w.Header().Set("KONTOL", "GEDE")

	loggedUser, err := h.userService.Login(input)
	if err != nil {
		response := helper.ApiResponse("Invalid email or password", http.StatusForbidden, "failed", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := h.authService.GenerateToken(loggedUser.Id)

	if err != nil {
		response := helper.ApiResponse("Failed to generate token", http.StatusUnprocessableEntity, "failed", err)
		json.NewEncoder(w).Encode(response)
		return
	}

	formatter := user.FormatUser(loggedUser, token)
	response := helper.ApiResponse("Berhasil login", http.StatusOK, "success", formatter)
	
	json.NewEncoder(w).Encode(response)

}

func (h *userHandler) CheckEmailAvailability(w http.ResponseWriter, r *http.Request){
	var input user.CheckEmailInput
	var response helper.Response
	log.Info("Check email avaibility called")
	reqBody, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Info("Can't react database, please try again later")
		return
	}
	json.Unmarshal(reqBody, &input)
	available, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		log.Info("Can't react database, please try again later")
		return
	}

	var isAvailable helper.EmailAvaibility
	if !available {
		isAvailable.IsAvailable = false
		response = helper.ApiResponse("Email is already taken", http.StatusForbidden, "failed", isAvailable)
	} else {
		isAvailable.IsAvailable = true
		response = helper.ApiResponse("Email is available", http.StatusOK, "success", isAvailable)
	}
	
	json.NewEncoder(w).Encode(response)
	
}

func (h *userHandler) UploadAvatar(w http.ResponseWriter, r *http.Request){
	var response helper.Response
	var isUploaded helper.UploadStatus

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("avatar")
	fileExtension := getFileExtension(handler)

	if err != nil {
		isUploaded.IsUploaded = false
		response = helper.ApiResponse("File uploading error", http.StatusForbidden, "failed", isUploaded)
		log.Info("Can't process file")
	}
	defer file.Close()

	path := "images/"
	tempFile, err := ioutil.TempFile(path, fmt.Sprintf("avatar-*.%s", fileExtension))
	fileName := tempFile.Name()

	if err != nil {
        log.Info(err)
    }
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
        log.Info(err)
    }

	tempFile.Write(fileBytes)
	userId := w.Header().Get("currentUser")
	userIdConv, _ := strconv.Atoi(userId)
	_, err = h.userService.SaveAvatar(userIdConv, fileName)

	if err != nil {
        log.Info(err)
    } else {
		isUploaded.IsUploaded = true
		response = helper.ApiResponse("File successfully uploaded", http.StatusOK, "success", isUploaded)
	}

	json.NewEncoder(w).Encode(response)
}

func validatorCustom(input interface{}) error {
	var validate *validator.Validate = validator.New()
	errValidate := validate.Struct(input)
	return errValidate
}