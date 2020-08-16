package url

import (
	"time"
)

type UrlModel struct {
	Id       int            `json:"id"`
	Url      string         `json:"url"`
	Interval float32        `json:"interval"`
	History  []HistoryEntry `json:"-"`
}

type HistoryEntry struct {
	Response  string    `json:"response"`
	Duration  float32   `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}
