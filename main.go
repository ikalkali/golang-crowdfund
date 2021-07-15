package main

import (
	"fmt"
	"golang-crowdfunding/user"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}


func main(){
	dsn := "root:root@tcp(127.0.0.1:3306)/golang_crowdfund?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Warn("Gagal koneksi ke database")
	}
	log.Info("Starting server")

	var users []user.User

	db.First(&users)

	nama := users[0].Name
	fmt.Println(nama)
}