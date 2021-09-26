package with


type Locker interface {
	Lock()
	Unlock()
}
//Lock simple wrapper for locking -> performing action -> unlocking
func Lock(l Locker, fn func()) {
	l.Lock()
	defer l.Unlock()

	fn()
}
