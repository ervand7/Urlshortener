package views

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
)

// GetStats get statistic of shortened urls.
func GetStats(ctx context.Context, storage models.Storage) (models.Stats, error) {
	var numberOfURLs, numberOfUsers int
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		res, err := storage.GetNumberOfURLs(ctx)
		if err != nil {
			logger.Logger.Error(err.Error())
			return err
		}
		numberOfURLs = res
		return nil
	})

	grp.Go(func() error {
		res, err := storage.GetNumberOfUsers(ctx)
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
