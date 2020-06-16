package makroud

// Observer is a collector that gathers various runtime error.
type Observer interface {
	// OnClose
	OnClose(err error, flags map[string]string)
	// OnRollback
	OnRollback(err error, flags map[string]string)
}
