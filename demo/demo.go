package main

import (
	"fmt"
	"os"

	"github.com/MaxSem/wikimediastreams"
)

func main() {
	var stream wikimediastreams.RecentChangesStream

	err := stream.Run(func(event *wikimediastreams.RecentChangesEvent) {
		fmt.Println(*event)
	}, func(err error) {
		fmt.Fprintln(os.Stderr, err)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
