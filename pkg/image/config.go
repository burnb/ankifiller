package image

type Config interface {
	GetAPIKey() string
	GetCx() string
	GetGl() *string
}
