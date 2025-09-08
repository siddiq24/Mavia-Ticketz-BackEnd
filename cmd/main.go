package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/siddiq24/Tickitz-DB/internal/configs"
	"github.com/siddiq24/Tickitz-DB/internal/routers"
)

// @title TICKITZ
// @version 1.0
// @description RESTful API created using TICKITZ
// @host localhost:8085
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db, err := configs.InitDB()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	defer db.Close()

	r := routers.InitRouter(db)

	port := ":8085"
	log.Println("Server running on:", port)
	DB := os.Getenv("DBNAME")
	log.Printf("\n\nCONNECT TO DATABASE : %s <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n\n", DB)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
