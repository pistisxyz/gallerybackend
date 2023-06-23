package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gallery/src/database"
	. "gallery/src/router"
	"gallery/src/utils"

	"github.com/joho/godotenv"
)

func _init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	UPLOAD_DIR = os.Getenv("UPLOAD_DIR")

	utils.CreateDirIfNotExist(UPLOAD_DIR)

	database.ConnectToDb()
	database.ConnectToRdb()
}

func getMetaData(file_location string) {
	cmdStruct := exec.Command("exiftool", file_location)

	_, err := cmdStruct.Output()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	_init()

	http.HandleFunc("/upload", CORS(AUTH_CHECK(UploadHandler)))

	http.HandleFunc("/images", CORS(AUTH_CHECK(ImagesHandler)))

	http.HandleFunc("/search", CORS(AUTH_CHECK(SearchHandler)))

	http.HandleFunc("/profile", CORS(AUTH_CHECK(ProfileHandler)))

	// Register the wrapped file server to handle requests at the "/image/" URL path
	http.Handle("/image/", http.StripPrefix("/image/", AuthMiddleware(http.FileServer(http.Dir("./uploaded")))))

	// Register the wrapped file server to handle requests at the "/compressed/" URL path
	http.Handle("/compressed/", http.StripPrefix("/compressed/", AuthMiddleware(http.FileServer(http.Dir("./compressed")))))

	fmt.Printf("Starting server at port %v\n", os.Getenv("PORT"))

	if err := http.ListenAndServe(os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
