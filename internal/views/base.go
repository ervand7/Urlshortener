package views

import (
	"context"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
)

const ctxTime time.Duration = 2

// Server struct for work with views.
type Server struct {
	Storage models.Storage
}

// GetOrCreateUserIDFromCookie returns UserID from cookie if it is there.
// Otherwise, creates new.
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
			logger.Logger.Error(err.Error())
			return ""
		}
		userID = string(decoded)
	}
	return userID
}

// GetStats get statistic of shortened urls.
func (server Server) GetStats(ctx context.Context) (models.Stats, error) {
	var numberOfURLs, numberOfUsers int
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		res, err := server.Storage.GetNumberOfURLs(ctx)
		if err != nil {
			logger.Logger.Error(err.Error())
			return err
		}
		numberOfURLs = res
		return nil
	})

	grp.Go(func() error {
		res, err := server.Storage.GetNumberOfUsers(ctx)
		if err != nil {
			logger.Logger.Error(err.Error())
			return err
		}
		numberOfUsers = res
		return nil
	})

	if err := grp.Wait(); err != nil {
		return models.Stats{}, err
	}
	return models.Stats{
		NumberOfURLs:  numberOfURLs,
		NumberOfUsers: numberOfUsers,
	}, nil
}

// Write writes resp data.
func (server Server) Write(msg []byte, w http.ResponseWriter) {
	_, err := w.Write(msg)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

// CloseBody closes resp body.
func (server Server) CloseBody(r *http.Request) {
	err := r.Body.Close()
	if err != nil {
		logger.Logger.Warn(err.Error())
	}
}

func (server Server) setCookie(encodedUserID string, w http.ResponseWriter) {
	cookie := &http.Cookie{Name: "userID", Value: encodedUserID, HttpOnly: true}
	http.SetCookie(w, cookie)
}
