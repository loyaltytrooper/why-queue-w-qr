package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
	"why-queue-w-qr/attendance/attendanceHandlers"
	"why-queue-w-qr/login/loginHandlers"
	"why-queue-w-qr/utils"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Loading .env data: %s", err.Error())
	}
}

func main() {
	postgreConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"))
	var err error

	loginHandlers.StudentsDB, err = sql.Open("postgres", postgreConn)
	if err != nil {
		log.Fatalf("Error Connecting to the students DB: %s", err.Error())
	}

	defer func(StudentsDB *sql.DB) {
		err := StudentsDB.Close()
		if err != nil {
			log.Fatalf("Error closing your Database! %s", err)
		}
	}(loginHandlers.StudentsDB)

	attendanceHandlers.MongoCTX, _ = context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(attendanceHandlers.MongoCTX, options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(attendanceHandlers.MongoCTX)

	app := fiber.New()
	v1 := app.Group("/v1")

	v1.Post(utils.Login, loginHandlers.LoginFunc)
	v1.Post(utils.MarkAttendance, attendanceHandlers.MarkAttendance)
	v1.Post(utils.AddExcusedAttendance, attendanceHandlers.AddExcusedAttendance) // TODO change to POST
	v1.Get(utils.GetAll, loginHandlers.GetAll)
	v1.Get(utils.GetStudentAttendance, attendanceHandlers.GetStudentAttendance)
	v1.Get(utils.GetClassAttendance, attendanceHandlers.GetClassAttendance)

	log.Fatal(app.Listen("127.0.0.1:9010"))
}