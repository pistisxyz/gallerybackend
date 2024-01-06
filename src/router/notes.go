package router

import (
	"encoding/binary"
	"encoding/json"
	"gallery/src/database"
	"gallery/src/utils"
	"net/http"
	"strings"
)

func NotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		notesGet(w, r)
	} else if r.Method == http.MethodPost {
		notesPost(w, r)
	}
}

func notesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// stmt, err := database.DB.Prepare("SELECT user_id, image_name, image_description, image_path, metadata, type, size FROM Images LIMIT ? OFFSET ?")
	stmt, err := database.DB.Prepare("SELECT notes_id, user_id, username, note_name, note, createdOn, updatedOn FROM notes ORDER BY RAND()")
	utils.CatchErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.NotesId, &note.UserId, &note.Username, &note.NoteName, &note.Note, &note.CreatedOn, &note.UpdatedOn)
		if err != nil {
			utils.CatchErr(err)
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		utils.CatchErr(err)
	}
	err = json.NewEncoder(w).Encode(notes)
	if utils.CatchErr(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func notesPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	note := r.FormValue("note")
	noteName := r.FormValue("note_name")
	if note == "" || noteName == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	query, err := database.DB.Prepare("INSERT INTO notes (user_id, username, note_name, note) VALUES (?, ?, ?, ?)")
	if utils.CatchErr(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer query.Close()

	Auth_header := r.Header.Get("Authorization")
	Auth_token := strings.Split(Auth_header, " ")[1]

	body := Note{}
	binary.Read(r.Body, binary.BigEndian, &body)
	_, err = query.Exec(database.RedisGetAuth(Auth_token), database.RedisGetAuthName(Auth_token), noteName, note)
	if utils.CatchErr(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Note struct {
	NotesId   uint64 `db:"notes_id"`
	UserId    uint64 `db:"user_id"`
	Username  string `db:"string"`
	NoteName  string `db:"note_name"`
	Note      string `db:"note"`
	CreatedOn string `db:"createdOn"`
	UpdatedOn string `db:"updatedOn"`
}
