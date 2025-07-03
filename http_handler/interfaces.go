package http_handler

type CustomCacheInterface interface {
	Put(key string, value string)
	Get(key string) (string, error)
}
