package repository

import (
	"context"
	"fmt"
	"retrovisionarios-api/internal/app/v1/videos/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoRepository struct {
	db *pgxpool.Pool
}

func NewVideoRepository(db *pgxpool.Pool) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) GetAll(ctx context.Context) ([]models.Video, error) {
	query := "SELECT id, title, subtitle, video_src, category FROM videos ORDER BY id ASC"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	videos := make([]models.Video, 0, 10)

	for rows.Next() {
		var v models.Video

		err := rows.Scan(&v.ID, &v.Title, &v.Subtitle, &v.VideoSrc, &v.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan video: %w", err)
		}

		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return videos, nil
}
