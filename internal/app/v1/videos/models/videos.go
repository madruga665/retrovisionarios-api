package models

type VideoCategory string

const (
	CategoryOriginalSong VideoCategory = "ORIGINAL SONG"
	CategoryCover        VideoCategory = "COVER"
)

type Video struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	Subtitle string        `json:"subtitle"`
	VideoSrc string        `json:"videoSrc"`
	Category VideoCategory `json:"category"`
}
