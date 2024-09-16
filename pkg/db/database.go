package db

import (
	"fmt"
	"os"
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
	return &PostgreSQLTokenDB{}
}

type PostgreSQLTokenDB struct{}

func (tdb *PostgreSQLTokenDB) SaveUserData(userGUID string, userEmail string, clientIP string, refreshTokenHash string) error {
	fmt.Printf("DB: Saved %s token to database for user %s", refreshTokenHash, userGUID)
	return nil
}

func (tdb *PostgreSQLTokenDB) FetchUserData(userGUID string) (UserData, error) {
	fmt.Printf("DB: A hashed refresh token of user %s has been accessed", userGUID)
	return UserData{}, nil
}
