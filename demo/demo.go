package main

import (
	"fmt"
	"github.com/MaxSem/wikimediastreams"
)

func main() {
	var stream wikimediastreams.RecentChangesStream
	stream.Run(func(event *wikimediastreams.RecentChangesEvent) {
		fmt.Println(*event)
	}, func(err error) {
		fmt.Println(err)
	})
}