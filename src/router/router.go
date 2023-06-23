package router

import (
	"fmt"
	"gallery/src/database"
	"net/http"
	"strings"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		UploadFile(w, r)
	}
}

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ImagesGet(w, r)
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		SearchPost(w, r)
	}
	if r.Method == http.MethodGet {
		SearchGet(w, r)
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

func AUTH_CHECK(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Auth_header := r.Header.Get("Authorization")
		if Auth_header == "" {
			fmt.Fprint(w, "not logged in")
			return
		}

		token_arr := strings.Split(Auth_header, " ")

		if len(token_arr) < 1 {
			fmt.Fprint(w, "error getting auth token... try to log out and log back in!")
			return
		}

		token := database.RedisGetAuth(token_arr[1])

		if token != "" {
			next(w, r)
		} else {
			fmt.Fprint(w, "Expired Session!")
		}
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		var Auth_header string
		// Loop through the cookies and print their names and values
		for _, cookie := range cookies {
			if cookie.Name == "token" {
				Auth_header = cookie.Value
			}
		}

		if Auth_header == "" {
			fmt.Fprint(w, "not logged in")
			return
		}

		token := database.RedisGetAuth(Auth_header)

		if token == "" {
			fmt.Fprint(w, "Expired Session!")
			return
		}
		next.ServeHTTP(w, r)
	})
}
