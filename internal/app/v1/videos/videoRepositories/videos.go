package videoRepositories

import (
	"context"
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
	query := "SELECT id, title, subtitle, video_src, category FROM videos"

	rows, err := r.db.Query(ctx, query)
	var videos []models.Video

	if err != nil {
		return videos, err
	}

	defer rows.Close()

	for rows.Next() {
		var e models.Video

		err := rows.Scan(&e.ID, &e.Title, &e.Subtitle, &e.VideoSrc, &e.Category)
		if err != nil {
			return videos, err
		}

		videos = append(videos, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
