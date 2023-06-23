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

		ctx := context.Background()
		user_map := database.AuthRdb.HGetAll(ctx, strings.Split(Auth_header, " ")[1]).Val()

		profile_data := struct {
			Username string
		}{
			Username: user_map["username"],
		}

		// Encode the array as JSON and send it in the response
		err := json.NewEncoder(w).Encode(profile_data)
		utils.CatchErr(err)

	}
}
