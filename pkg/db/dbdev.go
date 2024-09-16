package db

import "fmt"

type userInfo struct {
	GUID         string
	Email        string
	RefreshToken string
}

type InMemoryTokenDB struct {
	rows map[string]userInfo
}

func CreateNewInMemoryTokenDB() TokenDB {
	return &InMemoryTokenDB{
		rows: make(map[string]userInfo),
	}
}

func (tdb *InMemoryTokenDB) SaveUserData(userGUID string, userEmail string, refreshTokenHash string) error {
	tdb.rows[userGUID] =
		userInfo{
			GUID:         userGUID,
			Email:        userEmail,
			RefreshToken: refreshTokenHash,
		}
	return nil
}

func (tdb *InMemoryTokenDB) FetchHashedRefreshTokenFromDB(userGUID string) (string, error) {
	if userInfo, exists := tdb.rows[userGUID]; exists {
		return userInfo.RefreshToken, nil
	}
	return "", fmt.Errorf("user with GUID %s not found", userGUID)
}
