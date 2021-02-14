package bootstrap

type AfterServerStopFunc func() error

// CloseConfig ...
func CloseConfig() AfterServerStopFunc {
	return func() error {
		// ...
		return nil
	}
}
