package runtimeutil

func Go(fn func()) {
	defer func() {
		accessRecoverMu.RLock()
		defer accessRecoverMu.RUnlock()
		Recover()
	}()
	fn()
}