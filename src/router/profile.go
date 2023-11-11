package router

import (
	"context"
	"encoding/json"
	"gallery/src/database"
	"gallery/src/utils"
	"net/http"
	"strings"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		Auth_header := r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")

		ctx := context.Background()
		user_map := database.AuthRdb.HGetAll(ctx, strings.Split(Auth_header, " ")[1]).Val()

		profile_data := struct {
			Username     string
			TotalUploads string
		}{
			Username:     user_map["username"],
			TotalUploads: user_map["total_uploads"],
		}

		if profile_data.TotalUploads == "" {
			err := database.DB.QueryRow("SELECT COUNT(*) FROM Images WHERE user_id = ?", user_map["user:id"]).Scan(&profile_data.TotalUploads)
			utils.CatchErr(err)
		}

		// Encode the array as JSON and send it in the response
		err := json.NewEncoder(w).Encode(profile_data)
		utils.CatchErr(err)

	}
}
