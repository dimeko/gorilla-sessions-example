package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type RDBMS struct {
	Db *sql.DB
}

func env() map[string]string {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		panic("Cannot find .env file")
	}
	return map[string]string{
		"username": os.Getenv("POSTGRES_USER"),
		"host":     os.Getenv("POSTGRES_HOST"),
		"password": os.Getenv("POSTGRES_PASSWORD"),
		"db_name":  os.Getenv("POSTGRES_DB"),
		"port":     os.Getenv("POSTGRES_PORT"),
	}
}

func NewRDBMS() *RDBMS {
	envVars := env()
	host := envVars["host"]
	port := envVars["port"]
	user := envVars["username"]
	password := envVars["password"]
	dbname := envVars["db_name"]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	rdbms := &RDBMS{
		Db: db,
	}

	return rdbms
}

func (s *RDBMS) ListProductsRDBMS(limit string, offset string, filter string) ([]*Product, error) {
	rows, err := s.Db.Query("SELECT title, description, price FROM products WHERE title LIKE $1 LIMIT $2 OFFSET $3", filter+"%", limit, offset)
	if err != nil {
		log.Errorf("Error trying to query products. Params: %s, %s", limit, offset)
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Title, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (s *RDBMS) TotalProductsRDBMS() int {
	var total_count int
	err := s.Db.QueryRow("SELECT COUNT(*) FROM products").Scan(&total_count)
	if err != nil {
		panic(err.Error())
	}
	return total_count
}

func (s *RDBMS) UserExists(username string, password string) bool {
	var _username string
	err := s.Db.QueryRow("SELECT username FROM users WHERE username='?' AND password='?'", username, password).Scan(&_username)

	switch {
	case err == sql.ErrNoRows:
		return false
	default:
		return true
	}
}
