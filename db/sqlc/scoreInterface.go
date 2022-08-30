package db

import (
	"context"
	"database/sql"
)

type ScoreInterface interface {
	CreateScorePoint(ctx context.Context, arg CreateScorePointParams) error
	ScoreDetailsList(ctx context.Context, userID int32) ([]Score, error)
	ScoreDetailsListByCountry(ctx context.Context, country sql.NullString) ([]ScoreDetailsListByCountryRow, error)
	ScoreDetailsListByState(ctx context.Context, state sql.NullString) ([]ScoreDetailsListByStateRow, error)
	ScoreDetailsSats(ctx context.Context, arg ScoreDetailsSatsParams) ([]Score, error)
}
