package videoServices

import (
	"context"
	"retrovisionarios-api/internal/app/v1/videos/models"
	"strings"
)

type VideoRepository interface {
	GetAll(ctx context.Context) ([]models.Video, error)
}

type VideoService struct {
	repo VideoRepository
}

func NewVideoService(repository VideoRepository) *VideoService {
	return &VideoService{repo: repository}
}

func (s *VideoService) GetAll(ctx context.Context) (map[string][]models.Video, error) {
	videoList, err := s.repo.GetAll(ctx)
	groupedVideos := make(map[string][]models.Video)

	if err != nil {
		return map[string][]models.Video{}, err
	}

	for _, v := range videoList {
		key := strings.ToLower(strings.ReplaceAll(v.Category, " ", "_"))
		groupedVideos[key] = append(groupedVideos[key], v)
	}
	return groupedVideos, nil
}
