package cache

import (
	"testing"
)

func TestSetAndGet(t *testing.T) {
	cache := New[int, int](2)
	cache.Set(1, 10)
	value := cache.Get(1)
	if nil == value {
		t.Errorf("cache.Get(1) = %d; want 10", value)
	} else if 10 != *value {
		t.Errorf("cache.Get(1) = %d; want 10", *value)
	}
}

func TestDel(t *testing.T) {
	cache := New[int, int](2)
	cache.Set(1, 10)
	cache.Del(1)
	value := cache.Get(1)
	if nil != value {
		t.Errorf("cache.Get(1) = %d; want 10", value)
	}

	size := cache.Len()
	if 0 != size {
		t.Errorf("cache.Len() = %d; want 0", size)
	}
}

func TestCapacity(t *testing.T) {
	var size int
	cache := New[int, int](2)

	cache.Set(1, 10)
	size = cache.Len()
	if 1 != size {
		t.Errorf("cache.Len() = %d; want 1", size)
	}

	cache.Set(2, 10)
	size = cache.Len()
	if 2 != size {
		t.Errorf("cache.Len() = %d; want 2", size)
	}

	cache.Set(3, 10)
	size = cache.Len()
	if 2 != size {
		t.Errorf("cache.Len() = %d; want 2", size)
	}
}

func BenchmarkSetWithSmallCapacity(b *testing.B) {
	cache := New[int, int](2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i, 10)
	}
}

func BenchmarkSetWithLargeCapacity(b *testing.B) {
	cache := New[int, int](5000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(i, 10)
	}
}
