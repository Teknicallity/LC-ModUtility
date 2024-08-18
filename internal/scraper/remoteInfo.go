package scraper

import "time"

type RemoteInfo struct {
	ModVersion               string
	DownloadUrl              string
	LastUpdatedHumanReadable string
	LastUpdatedTime          time.Time
}
