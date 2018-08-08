package wikimediastreams

import (
	"encoding/json"
	"github.com/r3labs/sse"
)

// RecentChangesStream receives events about everything goin on in a wiki
type RecentChangesStream struct {
	Stream
}

// RecentChangesEvent contains information about recent changes
type RecentChangesEvent struct {
	Event

	Bot       bool          `json:"bot"`
	Comment   string        `json:"comment"`
	Length    newOldNumbers `json:"length"`
	Minor     bool          `json:"minor"`
	Namespace int           `json:"namespace"`
	Title     string        `json:"title"`
	Patrolled bool          `json:"patrolled"`
	Revision  newOldNumbers `json:"revision"`
	Domain    string        `json:"server_name"`
	Timestamp int           `json:"timestamp"`
	Type      string        `json:"type"`
	LogType   string        `json:"log_type"`
	User      string        `json:"user"`
	Wiki      string        `json:"wiki"`
}

// Run connects to the server and starts an infinite loop
func (s *RecentChangesStream) Run(receive func(*RecentChangesEvent), handleError func(error)) error {
	return s.runStream("recentchange", "mediawiki/recentchange/2", func(e *sse.Event) {
		var re RecentChangesEvent

		err := json.Unmarshal(e.Data, &re)
		if err != nil {
			handleError(err)
			return
		}
		valid, err := s.validateMetadata(&re.Meta)
		if err != nil {
			handleError(err)
			return
		}
		if valid {
			receive(&re)
		}
	})
}
