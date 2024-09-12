package api

import (
	"encoding/json"
	"net/http"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/auth"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type loginRequestBody struct {
	GUID string `json:"guid"`
}

type refreshRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

func sendErrorText(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func exitWith(err error, status int, writer http.ResponseWriter) bool {
	if err != nil {
		sendErrorText(writer, status)
		return true
	}
	return false
}

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		sendErrorText(writer, http.StatusMethodNotAllowed)
		return
	}

	var requestBody loginRequestBody
	json.NewDecoder(request.Body).Decode(&requestBody)
	userGUID := requestBody.GUID

	accessToken, err := auth.GenerateAccessToken(userGUID) // 15 minute expiration
	if exitWith(err, http.StatusBadRequest, writer) {
		return
	}

	refreshToken, hashedRefreshToken, err := auth.GenerateRefreshToken(userGUID) // 7 day expiration
	if exitWith(err, http.StatusBadRequest, writer) {
		return
	}

	err = db.SaveHashedRefreshToken(userGUID, hashedRefreshToken)
	if exitWith(err, http.StatusInternalServerError, writer) {
		return
	}

	json.NewEncoder(writer).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		sendErrorText(writer, http.StatusMethodNotAllowed)
		return
	}

	var requestBody refreshRequestBody
	json.NewDecoder(request.Body).Decode(&requestBody)
	refreshToken := requestBody.RefreshToken

	claims, err := auth.ValidateToken(refreshToken)
	if exitWith(err, http.StatusUnauthorized, writer) {
		return
	}

	userGUID := claims["guid"].(string)
	hashedTokenFromDB, err := db.FetchHashedRefreshTokenFromDB(userGUID)

	if exitWith(err, http.StatusInternalServerError, writer) {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedTokenFromDB), []byte(refreshToken))
	if exitWith(err, http.StatusUnauthorized, writer) {
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(userGUID) // 15 minute expiration
	if exitWith(err, http.StatusInternalServerError, writer) {
		return
	}

	json.NewEncoder(writer).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
