package api

import (
	"encoding/json"
	"net/http"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/auth"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/db"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/msg"
)

var tokenDatabase = db.ProvideApplicationTokenDB()
var emailNotificationService = msg.ProvideNotificationService()

type loginRequestBody struct {
	GUID  string `json:"guid"`
	Email string `json:"email"`
}

type refreshRequestBody struct {
	RefreshToken string `json:"refresh_token"`
}

func sendErrorText(writer http.ResponseWriter, status int) {
	http.Error(writer, http.StatusText(status), status)
}

func exitWithError(err error, status int, writer http.ResponseWriter) bool {
	if err != nil {
		sendErrorText(writer, status)
		return true
	}
	return false
}

func AccessHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		sendErrorText(writer, http.StatusMethodNotAllowed)
		return
	}

	var requestBody loginRequestBody
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if exitWithError(err, http.StatusBadRequest, writer) {
		return
	}

	userGUID := requestBody.GUID
	userEmail := requestBody.Email
	clientIP := request.RemoteAddr

	accessToken, refreshToken, err := auth.GeneratePair(userGUID, clientIP, userEmail, tokenDatabase)
	if exitWithError(err, http.StatusBadRequest, writer) {
		return
	}

	_, _, err = auth.ValidateAccessTokenClaims(accessToken, clientIP, userEmail, emailNotificationService)
	if exitWithError(err, http.StatusUnauthorized, writer) {
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
	err := json.NewDecoder(request.Body).Decode(&requestBody)
	if exitWithError(err, http.StatusBadRequest, writer) {
		return
	}
	refreshToken := requestBody.RefreshToken

	userGUID, userEmail, err := auth.ValidateRefreshTokenAndPassword(refreshToken, tokenDatabase)
	if exitWithError(err, http.StatusUnauthorized, writer) {
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(userGUID, request.RemoteAddr, userEmail)
	if exitWithError(err, http.StatusInternalServerError, writer) {
		return
	}

	json.NewEncoder(writer).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
