package with

//Closer interface with Close() error function
type Closer interface {
	Close() error
}

//Close is convenient wrapper for closing after performing some action
func Close(cl Closer, fn func()) {
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()

	fn()
}
