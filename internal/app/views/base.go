package views

import (
	"encoding/hex"
	"github.com/ervand7/urlshortener/internal/app/models"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const ctxTime time.Duration = 2

type Server struct {
	Storage models.Storage
}

func (server Server) GetOrCreateUserIDFromCookie(
	w http.ResponseWriter, r *http.Request,
) (userID string) {

	userIDFromCookie, err := r.Cookie("userID")
	if err != nil {
		userID = uuid.New().String()
		encoded := hex.EncodeToString([]byte(userID))
		server.setCookie(encoded, w)
	} else {
		encoded := userIDFromCookie.Value
		decoded, err := hex.DecodeString(encoded)
		if err != nil {
			utils.Logger.Error(err.Error())
			return ""
		}
		userID = string(decoded)
	}
	return userID
}

func (server Server) setCookie(encodedUserID string, w http.ResponseWriter) {
	cookie := &http.Cookie{Name: "userID", Value: encodedUserID, HttpOnly: true}
	http.SetCookie(w, cookie)
}

func (server Server) Write(msg []byte, w http.ResponseWriter) {
	_, err := w.Write(msg)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (server Server) CloseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		utils.Logger.Warn(err.Error())
	}
}
