package db

import "fmt"

type InMemoryTokenDB struct {
	rows map[string]UserData
}

func CreateNewInMemoryTokenDB() TokenDB {
	return &InMemoryTokenDB{
		rows: make(map[string]UserData),
	}
}

func (tdb *InMemoryTokenDB) SaveUserData(userGUID string, userEmail string, clientIP string, refreshTokenHash string) error {
	tdb.rows[userGUID] =
		UserData{
			GUID:             userGUID,
			Email:            userEmail,
			RecentIP:         clientIP,
			RefreshTokenHash: refreshTokenHash,
		}
	return nil
}

func (tdb *InMemoryTokenDB) FetchUserData(userGUID string) (UserData, error) {
	userData, exists := tdb.rows[userGUID]
	if exists {
		return userData, nil
	}
	return UserData{}, fmt.Errorf("user with GUID %s not found", userGUID)
}
