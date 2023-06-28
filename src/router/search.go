package router

import (
	"encoding/json"
	"fmt"
	"gallery/src/database"
	"gallery/src/utils"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		searchPost(w, r)
	}
	if r.Method == http.MethodGet {
		searchGet(w, r)
	}
}

func searchGet(w http.ResponseWriter, r *http.Request) {
	type Tag struct {
		TagName string `json:"tag_name"`
	}

	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT tag_name FROM Tags")
	utils.CatchErr(err)
	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.TagName)
		utils.CatchErr(err)

		tags = append(tags, tag)
	}
	rows.Close()

	// Encode the array as JSON and send it in the response
	err = json.NewEncoder(w).Encode(tags)
	utils.CatchErr(err)
}

func searchPost(w http.ResponseWriter, r *http.Request) {

	tagsRaw := r.FormValue("tags")

	if tagsRaw == "" {
		fmt.Printf("No tags provided")
		return
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}
	page_length, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		http.Error(w, "Invalid images per page value", http.StatusBadRequest)
		return
	}

	tags := []string{}
	err = json.Unmarshal([]byte(tagsRaw), &tags)
	utils.CatchErr(err, "Cant get tags")
	tags = utils.RemoveDuplicates(tags)

	if len(tags) < 1 {
		http.Error(w, "empty search!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	stmt, err := database.DB.Prepare(`
		SELECT DISTINCT i.*
		FROM Images i
		JOIN Image_Tags it ON i.image_id = it.image_id
		JOIN Tags t ON it.tag_id = t.tag_id
		WHERE t.tag_name IN (?` + strings.Repeat(",?", len(tags)-1) + `)
		LIMIT ? OFFSET ?
	`)
	utils.CatchErr(err)
	defer stmt.Close()

	countStmt, err := database.DB.Prepare(`
		SELECT COUNT(DISTINCT i.image_id)
		FROM Images i
		JOIN Image_Tags it ON i.image_id = it.image_id
		JOIN Tags t ON it.tag_id = t.tag_id
		WHERE t.tag_name IN (?` + strings.Repeat(",?", len(tags)-1) + `)
	`)
	utils.CatchErr(err)
	defer countStmt.Close()

	tagParams := make([]interface{}, len(tags))
	for i, tag := range tags {
		tagParams[i] = tag
	}

	offset := (page - 1) * page_length

	rows, err := stmt.Query(append(tagParams, page_length, offset)...)
	utils.CatchErr(err)
	defer rows.Close()

	var count int
	err = countStmt.QueryRow(tagParams...).Scan(&count)
	utils.CatchErr(err)

	var images []database.Images
	for rows.Next() {
		var img database.Images
		err := rows.Scan(&img.ID, &img.User_ID, &img.Name, &img.Description, &img.Path, &img.Metadata, &img.Type, &img.Size, &img.CreatedOn, &img.UpdatedOn)
		utils.CatchErr(err)
		images = append(images, img)
	}

	err = json.NewEncoder(w).Encode(struct {
		Count  int
		Images []database.Images
	}{
		count,
		images,
	})
	utils.CatchErr(err)
}
