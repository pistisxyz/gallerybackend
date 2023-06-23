package database

import (
	"context"
	"database/sql"
	"fmt"
	"gallery/src/utils"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

var DB *sql.DB
var AuthRdb *redis.Client

func RedisGetAuth(auth_token string) string { // TODO: add expired token error and logout (login prompt)
	ctx := context.Background()
	user_map := AuthRdb.HGetAll(ctx, auth_token).Val()

	return user_map["user:id"]

}

func ConnectToDb() {
	// Capture connection properties.
	config := mysql.NewConfig()
	config.User = os.Getenv("DB_USER")
	config.Passwd = os.Getenv("DB_PASSWORD")
	config.Net = "tcp"
	config.Addr = os.Getenv("DB_ADDR")
	config.DBName = os.Getenv("DB_NAME")

	// Get a database handle.
	var err error
	db, err := sql.Open("mysql", config.FormatDSN())
	utils.CatchErr(err)

	pingErr := db.Ping()
	if pingErr != nil {
		utils.CatchErr(pingErr)
	}
	fmt.Println("Connected!")
}

func ConnectToRdb() {
	AuthRdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_ADDR"),
		Password: os.Getenv("RDB_PASSWORD"),
		DB:       1,
	})
}

type TagDb struct {
	TagId   uint64 `db:"tag_id"`
	TagName string `db:"tag_name"`
}

type Images struct {
	ID          uint64 `db:"image_id"`
	User_ID     uint64 `db:"user_id"`
	Name        string `db:"image_name"`
	Description string `db:"image_description"`
	Path        string `db:"image_path"`
	Type        string `db:"type"`
	Size        uint   `db:"size"`
	CreatedOn   string `db:"created_on"`
	UpdatedOn   string `db:"updated_on"`
}

func TagsContainsString(arr []TagDb, target string) (bool, uint64) {
	for _, s := range arr {
		if strings.EqualFold(s.TagName, target) {
			return true, uint64(s.TagId)
		}
	}
	return false, 0
}
