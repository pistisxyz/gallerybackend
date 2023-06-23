package router

import (
	"encoding/json"
	"gallery/src/database"
	"gallery/src/utils"
	"net/http"
	"strconv"
)

func ImagesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	page_length, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, "Invalid images per page value", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * page_length

	stmt, err := database.DB.Prepare("SELECT user_id, image_name, image_description, image_path, type, size FROM Images LIMIT ? OFFSET ?")
	utils.CatchErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(page_length, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Images struct {
		User_ID     uint64 `db:"user_id"`
		Name        string `db:"image_name"`
		Description string `db:"image_description"`
		Path        string `db:"image_path"`
		Type        string `db:"type"`
		Size        uint   `db:"size"`
	}

	var images []Images
	for rows.Next() {
		var img Images
		err := rows.Scan(&img.User_ID, &img.Name, &img.Description, &img.Path, &img.Type, &img.Size)
		if err != nil {
			utils.CatchErr(err)
		}
		images = append(images, img)
	}
	if err := rows.Err(); err != nil {
		utils.CatchErr(err)
	}

	stmt, err = database.DB.Prepare("SELECT COUNT(*) FROM Images")
	utils.CatchErr(err)
	defer stmt.Close()

	// Execute the query and retrieve the count
	var count int
	err = stmt.QueryRow().Scan(&count)
	utils.CatchErr(err)

	// Encode the array as JSON and send it in the response
	err = json.NewEncoder(w).Encode(struct {
		Count  int
		Images []Images
	}{
		count,
		images,
	})
	utils.CatchErr(err)
}
