package bootstrap

type BeforeServerStartFunc func() error

func InitConfig() BeforeServerStartFunc {
	return func() error {
		// ...
		return nil
	}
}
