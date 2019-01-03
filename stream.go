// Package wikimediastreams provides functionality to receive notifications about changes
// on Wikimedia wikis, such as Wikipedia, using Server-Sent Events. See ps://wikitech.wikimedia.org/wiki/EventStreams
//
// Example usage:
// 	var stream wikimediastreams.RecentChangesStream
//
// 	// Optional configuration
//
// 	// By default, you receive events from all wikis. Filter by wiki domain:
// 	stream.FilterByDomain("en.wikipedia.org")
// 	// or use a wildcard:
// 	stream.FilterByDomain("*.wikipedia.org")
//
// 	// If you had to reconnect but need to not miss any events, pass the
// 	// last received event's Meta.DateTime to StartSince():
// 	stream.StartSince("<last timestamp here>")
//
// 	// To connect to some other, non-Wikimedia stream, use SetStreamURL():
// 	stream.SetStreamURL("https://example.com/stream")
//
// 	// End optional configuration
//
// 	// Connect to the server and wait for events indefinitely.
// 	err := stream.Run(func(event *wikimediastreams.RecentChangesEvent) {
// 		fmt.Println(*event)
// 	}, func(err error) {
// 		fmt.Fprintln(os.Stderr, err)
// 	})
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 	}
package wikimediastreams

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/r3labs/sse"
)

type newOldNumbers struct {
	New int `json:"new"`
	Old int `json:"old"`
}

// Metadata represents metadata present in every stream type
type Metadata struct {
	Domain    string `json:"domain"`
	DateTime  string `json:"dt"`
	ID        string `json:"id"`
	RequestID string `json:"request_id"`
	SchemaURI string `json:"schema_uri"`
	Topic     string `json:"topic"`
	URI       string `json:"uri"`
	Partition uint64 `json:"partition"`
	Offset    uint64 `json:"offset"`
}

// Event received
type Event struct {
	Meta Metadata `json:"meta"`
}

// Stream is a base for type-specific streams
type Stream struct {
	streamURL    string
	client       *sse.Client
	domainRegexp *regexp.Regexp
	schema       string
	since        string
}

// UnexpectedSchemaError is returned when a message with an unexpected schema_uri is received
type UnexpectedSchemaError struct {
	schema   string
	expected string
}

func (e *UnexpectedSchemaError) Error() string {
	return fmt.Sprintf("Received event with schema_uri='%s', '%s' expected", e.schema, e.expected)
}

// SetStreamURL allows to customize the URL to receive data from.
// Does nothing after Run() has been called on the stream.
func (s *Stream) SetStreamURL(url string) *Stream {
	s.streamURL = url
	return s
}

// FilterByDomain allows to filter by domain.
// It can match both literal ("en.wikipedia.org") and
// masked ("*.wikibooks.org") domains.
// Does nothing after Run() has been called on the stream.
func (s *Stream) FilterByDomain(filter string) error {
	re, retval := regexp.Compile("^" + strings.Replace(strings.Replace(filter, ".", "\\.", -1), "*", ".*", -1) + "$")
	s.domainRegexp = re
	return retval
}

// StartSince configures the stream to start reading events from some time in the past,
// represented by an ISO 8601 timestamp. Use it to avoid losing data on reconnects.
func (s *Stream) StartSince(time string) {
	s.since = "?since=" + url.QueryEscape(time)
}

func (s *Stream) validateMetadata(meta *Metadata) (bool, error) {
	if meta.SchemaURI != s.schema {
		return false, &UnexpectedSchemaError{meta.SchemaURI, s.schema}
	}
	if s.domainRegexp != nil && !s.domainRegexp.MatchString(meta.Domain) {
		return false, nil
	}
	return true, nil
}

// runStream connects to the server and begins receiving events,
// calling the callback for each of them. Blocks eternally,
// but can be called from a goroutine.
func (s *Stream) runStream(streamType string, expectedSchema string, callback func(*sse.Event)) error {
	if s.streamURL == "" {
		s.streamURL = "https://stream.wikimedia.org/v2/stream/"
	}
	s.schema = expectedSchema
	url := s.streamURL + streamType + s.since
	s.client = sse.NewClient(url)

	return s.client.Subscribe("messages", func(msg *sse.Event) {
		// Filter for weird empty events
		if len(msg.Data) == 0 {
			return
		}
		callback(msg)
	})
}
