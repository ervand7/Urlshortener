package views

import (
	"encoding/hex"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/models/url"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/google/uuid"
	"net/http"
)

type Server struct {
	MemoryStorage *url.MemoryStorage
	FileStorage   *url.FileStorage
	DBStorage     *url.DBStorage
}

func (server Server) GetOrCreateUserIDFromCookie(w http.ResponseWriter, r *http.Request) (userID string) {
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

func (server Server) SaveURL(userID, shorten, origin string, w http.ResponseWriter) {
	if config.GetConfig().DatabaseDSN != "" {
		if err := server.DBStorage.Set(userID, shorten, origin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	if config.GetConfig().FileStoragePath != "" {
		if err := server.FileStorage.Set(shorten, origin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	server.MemoryStorage.Set(userID, shorten, origin)
}

func (server Server) GetOriginByShort(short string) (origin string, err error) {
	if config.GetConfig().DatabaseDSN != "" {
		origin, err = server.DBStorage.Get(short)
		if err != nil {
			return "", err
		}
		return origin, nil
	}
	if config.GetConfig().FileStoragePath != "" {
		origin, err = server.FileStorage.Get(short)
		if err != nil {
			return "", nil
		}
		return origin, nil
	}
	origin = server.MemoryStorage.Get(short)
	return origin, nil
}

func (server Server) GetUserURLs(userID string) (userURLs []map[string]string, err error) {
	if config.GetConfig().DatabaseDSN != "" {
		userURLs, err = server.DBStorage.GetUserURLs(userID)
		if err != nil {
			return nil, err
		}
		return userURLs, nil
	}

	userURLs, err = server.MemoryStorage.GetUserURLs(userID)
	if err != nil {
		return nil, err
	}
	return userURLs, nil
}
