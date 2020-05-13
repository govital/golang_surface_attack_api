package crudProvider

type Interface interface {
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	Set(key, value string) error
	Increment(key string) error
}
