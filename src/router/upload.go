package router

import (
	"encoding/json"
	"fmt"
	"gallery/src/database"
	"gallery/src/utils"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var UPLOAD_DIR string

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	// r.ParseMultipartForm(100 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	// fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	tagsRaw := r.FormValue("tags")
	tags := []string{}
	err = json.Unmarshal([]byte(tagsRaw), &tags)
	utils.CatchErr(err, "Cant get tags")
	tags = utils.RemoveDuplicates(tags)

	extension := filepath.Ext(handler.Filename)
	newFilename := utils.GenerateUniqueFilename() + extension

	// Create file
	dst, err := os.Create(UPLOAD_DIR + "/" + newFilename)
	if dst != nil {
		defer dst.Close()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the INSERT statement
	stmt, err := database.DB.Prepare("INSERT INTO Images (image_name, user_id, image_description, image_path, type) VALUES (?, ?, ?, ?, ?)")
	utils.CatchErr(err)
	defer stmt.Close()

	Auth_header := r.Header.Get("Authorization")

	user := database.RedisGetAuth(strings.Split(Auth_header, " ")[1])

	value, err := strconv.ParseUint(user, 0, 64)
	if err != nil {
		panic(err)
	}

	file_path := UPLOAD_DIR + "/" + newFilename

	file_ext := path.Ext(file_path)

	var file_type = "unknown"
	if utils.ContainsString(utils.ImageFormats, file_ext) {
		file_type = "image"
		go utils.MakeMediaThumbnail("./"+file_path, file_type)
	} else if utils.ContainsString(utils.VideoExtensions, file_ext) {
		file_type = "video"
		go utils.MakeMediaThumbnail("./"+file_path, file_type)
	}

	// Execute the INSERT statement with the image values
	result, err := stmt.Exec(handler.Filename, value, "Example image", file_path, file_type)
	utils.CatchErr(err)

	imageId, err := result.LastInsertId()
	utils.CatchErr(err)

	go func() {
		rows, err := database.DB.Query("SELECT * FROM Tags")
		utils.CatchErr(err)
		var tagsDb []database.TagDb
		for rows.Next() {
			var tag database.TagDb
			err := rows.Scan(&tag.TagId, &tag.TagName)
			utils.CatchErr(err)

			tagsDb = append(tagsDb, tag)
		}
		rows.Close()

		for _, tag := range tags {
			var exists, tagId = database.TagsContainsString(tagsDb, tag)

			if !exists {
				query, err := database.DB.Prepare("INSERT INTO Tags (tag_name) VALUES (?)")
				utils.CatchErr(err)
				result, err := query.Exec(tag)
				utils.CatchErr(err)
				lastId, err := result.LastInsertId()
				utils.CatchErr(err)

				tagId = uint64(lastId)
			}

			query, err := database.DB.Prepare("INSERT INTO Image_Tags (image_id, tag_id) VALUES (?, ?)")
			utils.CatchErr(err)
			_, err = query.Exec(imageId, tagId)
			utils.CatchErr(err)
		}
	}()

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
