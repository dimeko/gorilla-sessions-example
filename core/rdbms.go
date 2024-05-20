package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type RDBMS struct {
	Db *sql.DB
}

type Product struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
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
	// rows, err := s.Db.Query("SELECT name, title, description, price FROM products WHERE name LIKE '" + filter + "%' OR title LIKE '" + filter + "%' LIMIT " + limit + " OFFSET " + offset)
	rows, err := s.Db.Query("SELECT name, title, description, price FROM products WHERE name LIKE $1 OR title LIKE $1 LIMIT $2 OFFSET $3", filter+"%", limit, offset)
	if err != nil {
		logger.Errorf("Error trying to query products. Params: %s, %s", limit, offset)
		logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var products []*Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name, &p.Title, &p.Description, &p.Price); err != nil {
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

func _compare_passwords(hashed_pwd string, plain_pwd []byte) bool {
	byte_hash := []byte(hashed_pwd)
	err := bcrypt.CompareHashAndPassword(byte_hash, plain_pwd)
	if err != nil {
		logger.Println(err)
		return false
	}

	return true
}

func (s *RDBMS) Login(username string, password string) bool {
	var pwd string
	err := s.Db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&pwd)
	if err != nil {
		logger.Println(err)
		return false
	}
	return _compare_passwords(pwd, []byte(password))
}
