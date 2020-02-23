package lrucache

import (
	"strconv"
	"testing"
)

// internal key/val struct wraper
type keyVal struct {
	key int
	val int
}

// init cache with low size capacity
var cache LRUCache = constructor(7)
var lruCacheTests = [][]struct {
	method   string
	intput   interface{}
	expected interface{}
}{
	{
		{"PUT", keyVal{2, 6}, true},
		{"PUT", keyVal{1, 5}, true},
		{"PUT", keyVal{1, 2}, true},
		{"PUT", keyVal{3, 7}, true},
	},
	{
		{"PUT", keyVal{5, 6}, true},
		{"PUT", keyVal{23, 8}, true},
	},
	{
		{"PUT", keyVal{9, 12}, true},
		{"PUT", keyVal{99, 1122}, true},
	},
	{
		{"GET", 2, 6},
		{"GET", 1, 2},
		{"GET", 99, 1122},
	},
	{
		{"GET", 5, 6},
		{"GET", 9, 12},
	},
	{
		{"GET", 8, -1},
		{"GET", 3, 7},
	},
	// eviction
	{
		{"PUT", keyVal{10, 11}, true},
		{"GET", 10, 11},
		{"PUT", keyVal{11, 15}, true},

		{"PUT", keyVal{12, 14}, true},
		{"GET", 12, 14},

		{"PUT", keyVal{12, 34}, true},
		{"GET", 12, 34},
		{"PUT", keyVal{5, 18}, true},
		{"GET", 3, 7},
		{"GET", 5, 18},
		{"PUT", keyVal{13, 17}, true},
		{"GET", 13, 17},
		{"GET", 5, 18},
		{"PUT", keyVal{14, 57}, true},
		{"PUT", keyVal{14, 58}, true},
		{"GET", 14, 58},
	},
}
var benchLruCacheTests = []struct {
	method string
	intput []interface{}
}{
	{"PUT", []interface{}{keyVal{2, 6}, keyVal{1, 5}, keyVal{1, 2}, keyVal{3, 7}, keyVal{5, 6}, keyVal{23, 8}, keyVal{9, 12}, keyVal{99, 1122}}},
	{"GET", []interface{}{2, 1, 99, 9, 5, 3, 23, 8, 99, 35, 53, 5, 3}},
}

// lRUAccess test processor
func lRUAccess(index int, t *testing.T) {
	for _, k := range lruCacheTests[index] {
		switch k.method {
		case "GET":
			if actual := cache.Get(k.intput.(int)); actual != k.expected.(int) {
				t.Errorf("expected %d, actual %d for key %d", k.expected.(int), actual, k.intput.(int))
			}
		case "PUT":
			cache.Put(k.intput.(keyVal).key, k.intput.(keyVal).val)
		default:
			t.Errorf("Unknown Method")
		}
	}
}

func TestLRUParallel(t *testing.T) {
	// parallel tests wrapper
	parallelTest := []func(t *testing.T){
		func(t *testing.T) {
			lRUAccess(0, t)
		},
		func(t *testing.T) {
			lRUAccess(1, t)
		},
		func(t *testing.T) {
			lRUAccess(2, t)
		},
	}
	// test vector
	testGroup := []map[int]func(t *testing.T){
		{1: parallelTest[0]},
		{2: parallelTest[1]},
		{3: parallelTest[2]},
	}

	// runs tests in parallel.
	for k, fc := range testGroup {
		k := k
		fc := fc
		t.Run("Test"+strconv.Itoa(k), func(t *testing.T) {
			t.Parallel()
			// call test function
			fc[k+1](t)
		})
	}
}
func TestReads(t *testing.T) {
	// sequential test wrapper
	seqTest := []func(t *testing.T){
		func(t *testing.T) {
			lRUAccess(3, t)
		},
		func(t *testing.T) {
			lRUAccess(4, t)
		},
		func(t *testing.T) {
			lRUAccess(5, t)
		},
		func(t *testing.T) {
			lRUAccess(6, t)
		},
	}
	// run group sequentially
	t.Run("seqgroup", func(t *testing.T) {
		for k := range seqTest {
			t.Run("TestSeq"+strconv.Itoa(k+1), seqTest[k])
		}
	})
}
func BenchmarkLRU(b *testing.B) {
	for i := range benchLruCacheTests {
		switch benchLruCacheTests[i].method {
		case "PUT":
			b.Run("PutLRU", func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						for _, k := range benchLruCacheTests[i].intput {
							cache.Put(k.(keyVal).key, k.(keyVal).val)
						}
					}
				})
			})
		case "GET":
			b.Run("GetLRU", func(b *testing.B) {
				for kk := 0; kk < b.N; kk++ {
					for _, k := range benchLruCacheTests[i].intput {
						cache.Get(k.(int))
					}
				}
			})
		default:
			b.Errorf("Unknown Method")
		}
	}
}
