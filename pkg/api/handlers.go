package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/Andrey-Kachow/goauth-backdev/pkg/app"
	"github.com/Andrey-Kachow/goauth-backdev/pkg/auth"
)

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

func ipAddrFromRequest(request *http.Request) string {
	if os.Getenv("GOAUTH_BACKDEV_MODE") == "development" {
		//
		// Using string of type "ip:port" instead of Ip allows simplier testing IP change from localhost.
		// Sending request from other device will be sent from different port of the Wi-fi router.
		//
		return request.RemoteAddr
	}
	return strings.Split(request.RemoteAddr, ":")[0]
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
	clientIP := ipAddrFromRequest(request)

	accessToken, refreshToken, err := auth.GeneratePair(userGUID, clientIP, userEmail, app.Context())
	if exitWithError(err, http.StatusBadRequest, writer) {
		return
	}

	_, _, err = auth.ValidateAccessTokenClaims(accessToken, clientIP, userEmail)
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

	userGUID, userEmail, err := auth.ValidateRefreshTokenAndPassword(refreshToken, app.Context())
	if exitWithError(err, http.StatusUnauthorized, writer) {
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(userGUID, ipAddrFromRequest(request), userEmail)
	if exitWithError(err, http.StatusInternalServerError, writer) {
		return
	}

	json.NewEncoder(writer).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}
