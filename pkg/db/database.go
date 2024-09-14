package db

import (
	"fmt"
	"os"
)

type TokenDB interface {
	SaveUserData(userGUID string, userEmail string, refreshTokenHash string) error
	FetchHashedRefreshTokenFromDB(userGUID string) (string, error)
	GetEmailAddressFromGUID(userGUID string) (string, error)
}

func ProvideApplicationTokenDB() TokenDB {
	return &PostgreSQLTokenDB{}
}

type PostgreSQLTokenDB struct {
}

func (tdb *PostgreSQLTokenDB) SaveUserData(userGUID string, userEmail string, refreshTokenHash string) error {
	fmt.Printf("DB: Saved %s token to database for user %s", refreshTokenHash, userGUID)
	return nil
}

func (tdb *PostgreSQLTokenDB) FetchHashedRefreshTokenFromDB(userGUID string) (string, error) {
	fmt.Printf("DB: A hashed refresh token of user %s has been accessed", userGUID)
	return "sample db data", nil
}

func (tdb *PostgreSQLTokenDB) GetEmailAddressFromGUID(userGUID string) (string, error) {
	return os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME"), nil
}
