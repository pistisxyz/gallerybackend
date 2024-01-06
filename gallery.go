package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func main() {
	_init()

	http.HandleFunc("/upload", CORS(AUTH_CHECK(UploadHandler)))

	http.HandleFunc("/images", CORS(AUTH_CHECK(ImagesHandler)))

	http.HandleFunc("/search", CORS(AUTH_CHECK(SearchHandler)))

	http.HandleFunc("/notes", CORS(AUTH_CHECK(NotesHandler)))

	http.HandleFunc("/profile", CORS(AUTH_CHECK(ProfileHandler)))

	http.HandleFunc("/tags", CORS(AUTH_CHECK(TagsHandler)))

	// Register the wrapped file server to handle requests at the "/image/" URL path
	http.Handle("/file/", http.StripPrefix("/file/", AuthMiddleware(http.FileServer(http.Dir("./uploaded")))))

	// Register the wrapped file server to handle requests at the "/compressed/" URL path
	http.Handle("/compressed/", http.StripPrefix("/compressed/", AuthMiddleware(http.FileServer(http.Dir("./compressed")))))

	// Register the wrapped file server to handle requests at the "/share/" URL path
	http.Handle("/share/file/", http.StripPrefix("/share/file/", ShareAuthMiddleware(http.FileServer(http.Dir("./compressed")))))

	fmt.Printf("Starting server at port %v\n", os.Getenv("PORT"))

	if err := http.ListenAndServe(os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
