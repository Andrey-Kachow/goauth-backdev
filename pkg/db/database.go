package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type UserData struct {
	GUID             string
	Email            string
	RecentIP         string
	RefreshTokenHash string
}

type TokenDB interface {
	SaveUserData(userGUID string, userEmail string, clientIP string, refreshTokenHash string) error
	FetchUserData(userGUID string) (UserData, error)
}

func ProvideApplicationTokenDB() TokenDB {
	if os.Getenv("GOAUTH_BACKDEV_MODE") == "development" {
		fmt.Println("Using the in-memory database")
		return CreateNewInMemoryTokenDB()
	}
	sqlDB, err := InitPostgresDB()
	if err != nil {
		log.Fatal("Failed to Initialize PostgreSQL DB")
	}
	return &PostgreSQLTokenDB{
		DatabaseSQL: sqlDB,
	}
}

type PostgreSQLTokenDB struct {
	DatabaseSQL *sql.DB
}

func InitPostgresDB() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

func (tdb *PostgreSQLTokenDB) SaveUserData(userGUID string, userEmail string, clientIP string, refreshTokenHash string) error {
	_, err := tdb.DatabaseSQL.Exec(`
		INSERT INTO users (guid, email, client_ip, refresh_token)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (guid) 
		DO UPDATE SET 
			refresh_token = EXCLUDED.refresh_token,
			client_ip = EXCLUDED.client_ip;
		`,
		userGUID, userEmail, clientIP, refreshTokenHash)
	if err != nil {
		return fmt.Errorf("failed to save user data: %v", err)
	}
	return nil
}

func (tdb *PostgreSQLTokenDB) FetchUserData(userGUID string) (UserData, error) {
	var userData UserData
	err := tdb.DatabaseSQL.QueryRow(`
		SELECT guid, email, client_ip, refresh_token
		FROM users
		WHERE guid = $1`, userGUID).Scan(&userData.GUID, &userData.Email, &userData.RecentIP, &userData.RefreshTokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserData{}, fmt.Errorf("user not found")
		}
		return UserData{}, fmt.Errorf("failed to fetch user data: %v", err)
	}
	return userData, nil
}
