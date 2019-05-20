package piaas

type Event struct {
	Path string
}

type Monitor interface {
	// Start watching for file system changes. Specify the debounceTime so the events are collected within the debounce time.
	Start(debounceTime uint64)
	// Register a observer to receive the fs events
	Subscribe(observer chan<- []string)
	// Stop file system watching.
	Stop()
}
