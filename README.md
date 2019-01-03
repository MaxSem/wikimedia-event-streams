# wikimediastreams

[![GoDoc](https://godoc.org/github.com/MaxSem/wikimediastreams?status.svg)](https://godoc.org/github.com/MaxSem/wikimediastreams)

Package wikimediastreams provides functionality to receive notifications about changes
on Wikimedia wikis, such as Wikipedia, using Server-Sent Events. See ps://wikitech.wikimedia.org/wiki/EventStreams

Example usage:
    var stream wikimediastreams.RecentChangesStream

    // Optional configuration

    // By default, you receive events from all wikis. Filter by wiki domain:
    stream.FilterByDomain("en.wikipedia.org")
    // or use a wildcard:
    stream.FilterByDomain("*.wikipedia.org")

    // If you had to reconnect but need to not miss any events, pass the
    // last received event's Meta.DateTime to StartSince():
    stream.StartSince("<last timestamp here>")

    // To connect to some other, non-Wikimedia stream, use SetStreamURL():
    stream.SetStreamURL("https://example.com/stream")

    // End optional configuration

    // Connect to the server and wait for events indefinitely.
	err := stream.Run(func(event *wikimediastreams.RecentChangesEvent) {
		fmt.Println(*event)
	}, func(err error) {
		fmt.Fprintln(os.Stderr, err)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}