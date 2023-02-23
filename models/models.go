package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type PublicAccount struct {
	ID      int
	BizID   string
	Name    string
	Account string
}

type Article struct {
	ID         int
	URL        string
	AccountID  int
	CreateTime time.Time
}

var db *sql.DB

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	db, err = sql.Open("mysql", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %s", err.Error())
	}
}

func GetPublicAccountByName(name string) (PublicAccount, error) {
	account := PublicAccount{}

	err := db.QueryRow("SELECT id, biz_id, name, account FROM public_accounts WHERE name = ?", name).Scan(&account.ID, &account.BizID, &account.Name, &account.Account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, ErrNotFound
		}
		return account, err
	}

	return account, nil
}

func GetPublicAccountByBizID(bizID string) (PublicAccount, error) {
	account := PublicAccount{}

	err := db.QueryRow("SELECT id, biz_id, name, account FROM public_accounts WHERE biz_id = ?", bizID).Scan(&account.ID, &account.BizID, &account.Name, &account.Account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return account, ErrNotFound
		}
		return account, err
	}

	return account, nil
}

func InsertArticle(article Article) error {
	stmt, err := db.Prepare("INSERT INTO articles(url, account_id, create_time) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.URL, article.AccountID, article.CreateTime)
	if err != nil {
		return err
	}

	return nil
}

func IsArticleExist(url string) (bool, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM articles WHERE url = ?", url).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
