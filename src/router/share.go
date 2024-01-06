package router

import (
	"context"
	"encoding/json"
	"gallery/src/database"
	"gallery/src/utils"
	"io"
	"net/http"
	"strings"
)

const shareRdbPrefix = "gallery_share_"

func ShareHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		SharePost(w, r)
	}
}

func SharePost(w http.ResponseWriter, r *http.Request) { //TODO: add pages

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	result := database.DataRdb.HGetAll(ctx, shareRdbPrefix+string(body)).Val()

	tags := []string{}
	err = json.Unmarshal([]byte(result["tags"]), &tags)
	if utils.CatchErr(err, "Cant get tags") {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}
	tags = utils.RemoveDuplicates(tags)

	w.Header().Set("Content-Type", "application/json")

	stmt, err := database.DB.Prepare(`
		SELECT DISTINCT i.*
		FROM Images i
		JOIN Image_Tags it ON i.image_id = it.image_id
		JOIN Tags t ON it.tag_id = t.tag_id
		WHERE t.tag_name IN (?` + strings.Repeat(",?", len(tags)-1) + `)
	`)
	utils.CatchErr(err)
	defer stmt.Close()

	tagParams := make([]interface{}, len(tags))
	for i, tag := range tags {
		tagParams[i] = tag
	}

	rows, err := stmt.Query(tagParams...)
	utils.CatchErr(err)
	defer rows.Close()

	var images []database.Images
	for rows.Next() {
		var img database.Images
		err := rows.Scan(&img.ID, &img.User_ID, &img.Name, &img.Description, &img.Path, &img.Metadata, &img.Type, &img.Size, &img.CreatedOn, &img.UpdatedOn)
		utils.CatchErr(err)
		images = append(images, img)
	}

	err = json.NewEncoder(w).Encode(struct {
		Images []database.Images
	}{
		images,
	})
	utils.CatchErr(err)

}

//Establishing contact with transdimensional wizards since 1996.
