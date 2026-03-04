package models

type VideoCategory string

const (
	CategoryOriginalSong VideoCategory = "ORIGINAL SONG"
	CategoryCover        VideoCategory = "COVER"
)

type Video struct {
	ID       int           `json:"id"`
	Title    string        `json:"title" binding:"required"`
	Subtitle string        `json:"subtitle"`
	VideoSrc string        `json:"videoSrc" binding:"required,url"`
	Category VideoCategory `json:"category" binding:"required"`
}
