package cache

type ICache[K comparable, V any] interface {
	Get(key K) *V
	Set(key K, value V)
	Del(key K)
	Has(key K) bool
	Len() int
}
