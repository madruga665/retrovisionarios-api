package models

type Video struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	VideoSrc string `json:"videoSrc"`
	Category string `json:"category"`
}
