package views

import (
	"context"
	"encoding/hex"
	"github.com/ervand7/urlshortener/internal/app/models"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

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

func (server Server) SaveURL(userID, shorten, origin string, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := server.Storage.Set(ctx, userID, shorten, origin); err != nil {
		return err
	}
	return nil
}

func (server Server) SaveURLs(dbEntries []utils.DBEntry, r *http.Request) error {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := server.Storage.SetMany(ctx, dbEntries); err != nil {
		return err
	}
	return nil
}

func (server Server) GetOriginByShort(
	short string, r *http.Request,
) (origin string, err error) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	origin, err = server.Storage.Get(ctx, short)
	if err != nil {
		return "", err
	}
	return origin, nil
}

func (server Server) GetUserURLs(
	userID string, r *http.Request,
) (userURLs []map[string]string, err error) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	userURLs, err = server.Storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userURLs, nil
}

func (server Server) DeleteURLs(urlsFromRequest []string) {
	server.Storage.DeleteURLs(urlsFromRequest)
}
