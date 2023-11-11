package router

import (
	"encoding/json"
	"gallery/src/database"
	"gallery/src/utils"
	"net/http"
)

func TagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tagsGet(w, r)
	}
}

type TagImagesResponse struct {
	TagID            int     `json:"tag_id"`
	TagName          string  `json:"tag_name"`
	ImageID          int     `json:"image_id"`
	ImageName        string  `json:"image_name"`
	ImageDescription string  `json:"image_description"`
	ImagePath        string  `json:"image_path"`
	Metadata         string  `json:"metadata"`
	Type             string  `json:"type"`
	Size             int     `json:"size"`
	CreatedOn        []uint8 `json:"created_on"`
	UpdatedOn        []uint8 `json:"updated_on"`
}

const SQL_STRING = `WITH TagImages AS (
    SELECT
        t.tag_id,
        t.tag_name,
        i.image_id,
        i.image_name,
        i.image_description,
        i.image_path,
        i.metadata,
        i.type,
        i.size,
        i.createdOn,
        i.updatedOn,
        ROW_NUMBER() OVER (PARTITION BY t.tag_id ORDER BY i.image_id) AS row_num
    FROM
        Tags t
    LEFT JOIN
        Image_Tags it ON t.tag_id = it.tag_id
    LEFT JOIN
        Images i ON it.image_id = i.image_id
)
SELECT
    tag_id,
    tag_name,
    image_id,
    image_name,
    image_description,
    image_path,
    metadata,
    type,
    size,
    createdOn,
    updatedOn
FROM
    TagImages
WHERE
    row_num = 1;`

func tagsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := database.DB.Query(SQL_STRING)
	utils.CatchErr(err)
	var tags []TagImagesResponse
	for rows.Next() {
		var tag TagImagesResponse
		err := rows.Scan(
			&tag.TagID,
			&tag.TagName,
			&tag.ImageID,
			&tag.ImageName,
			&tag.ImageDescription,
			&tag.ImagePath,
			&tag.Metadata,
			&tag.Type,
			&tag.Size,
			&tag.CreatedOn,
			&tag.UpdatedOn,
		)
		utils.CatchErr(err)

		tags = append(tags, tag)
	}
	rows.Close()

	err = json.NewEncoder(w).Encode(tags)
	utils.CatchErr(err)
}
