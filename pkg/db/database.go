package db

import "fmt"

func SaveHashedRefreshToken(userGUID string, token string) error {
	fmt.Printf("DB: Saved %s token to database for user %s", token, userGUID)
	return nil
}

func FetchHashedRefreshTokenFromDB(userGUID string) (string, error) {
	fmt.Printf("DB: A hashed refresh token of user %s has been accessed", userGUID)
	return "sample db data", nil
}
