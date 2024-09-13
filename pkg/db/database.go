package db

import (
	"fmt"
	"os"
)

type TokenDB interface {
	SaveHashedRefreshToken(userGUID string, refreshTokenHash string) error
	FetchHashedRefreshTokenFromDB(userGUID string) (string, error)
	GetEmailAddressFromGUID(userGUID string) (string, error)
}

type PostgreSQLTokenDB struct {
}

func (tdb *PostgreSQLTokenDB) SaveHashedRefreshToken(userGUID string, refreshTokenHash string) error {
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
