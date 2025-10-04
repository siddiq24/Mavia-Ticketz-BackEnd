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
	// init database
	db, err := configs.InitDB()
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}
	defer db.Close()

	// init redis
	rc, err := configs.InitRedis()
	if err != nil {
		log.Fatal("failed to connect redis: ", err)
	}
	defer rc.Close()

	// init router dengan dependency db & redis
	r := routers.InitRouter(db, rc)

	// server port
	port := ":8085"
	log.Println("Server running on:", port)

	DB := os.Getenv("DB_NAME")
	log.Printf("\n\nCONNECT TO DATABASE : %s <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n\n", DB)

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
