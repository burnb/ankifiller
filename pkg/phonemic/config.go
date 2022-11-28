package phonemic

type Config interface {
	GetLocale() string
	GetSystem() string
}
