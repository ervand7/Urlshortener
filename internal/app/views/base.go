package views

import (
	"fmt"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/controllers/usertoken"
	"github.com/ervand7/urlshortener/internal/app/models/url"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/google/uuid"
	"net/http"
)

type Server struct {
	MemoryStorage *url.MemoryStorage
	FileStorage   *url.FileStorage
	UserToken     *usertoken.UserToken
}

func (server Server) CloseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		utils.Logger.Warn(err.Error())
	}
}

func (server Server) Write(msg []byte, w http.ResponseWriter) {
	_, err := w.Write(msg)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (server Server) SaveURL(userID, shorten, origin string, w http.ResponseWriter) {
	switch config.GetConfig().FileStoragePath {
	case "":
		server.MemoryStorage.Set(userID, shorten, origin)
	default:
		if err := server.FileStorage.Set(shorten, origin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (server Server) GetUserURLs(userID string) (userURLs []map[string]string) {
	userURLs = make([]map[string]string, 0)
	for _, data := range server.MemoryStorage.HashTable {
		if data.UserID == userID {
			pair := map[string]string{
				"short_url":    data.Shorten,
				"original_url": data.Origin,
			}
			userURLs = append(userURLs, pair)
		}
	}
	return userURLs
}

func (server Server) GetUserIDFromCookie(w http.ResponseWriter, r *http.Request) (userID string) {
	userIDFromCookie, err := r.Cookie("userID")
	if err != nil {
		userID = uuid.New().String()
		server.setCookie(userID, w)
	} else {
		userID = userIDFromCookie.Value
	}
	return userID
}

func (server Server) setCookie(userID string, w http.ResponseWriter) {
	encodedUserID, _ := server.UserToken.Encode(userID)
	fmt.Println("+++++++++++Закодированный UserID, который пойдет в куку: encodedUserID")
	fmt.Println("+++++++++++Key на момент кодирования: ", server.UserToken.Key)
	fmt.Println("+++++++++++Nonce на момент кодирования: ", server.UserToken.Nonce)
	fmt.Println("+++++++++++AesGCM на момент кодирования: ", server.UserToken.AesGCM)
	cookie := &http.Cookie{Name: "userID", Value: encodedUserID, HttpOnly: true}
	http.SetCookie(w, cookie)
}
