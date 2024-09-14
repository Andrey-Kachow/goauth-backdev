package db

import (
	"fmt"
	"os"
)

type InMemoryTokenDB struct {
	rows []struct {
		guid  string
		email string
	}
}

func (tdb *InMemoryTokenDB) SaveHashedRefreshToken(userGUID string, refreshTokenHash string) error {
	fmt.Printf("DB: Saved %s token to database for user %s", refreshTokenHash, userGUID)
	return nil
}

func (tdb *InMemoryTokenDB) FetchHashedRefreshTokenFromDB(userGUID string) (string, error) {
	fmt.Printf("DB: A hashed refresh token of user %s has been accessed", userGUID)
	return "sample db data", nil
}

func (tdb *InMemoryTokenDB) GetEmailAddressFromGUID(userGUID string) (string, error) {
	return os.Getenv("GOAUTH_BACKDEV_EMAIL_USERNAME"), nil
}
