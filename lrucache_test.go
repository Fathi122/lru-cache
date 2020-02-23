package lrucache

import (
	"strconv"
	"testing"
)

type keyVal struct {
	key int
	val int
}

var cache3 LRUCache = constructor(6)
var lrucacheTests = [][]struct {
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
		{"GET", 2, -1}, // early eviction
		{"GET", 1, 2},
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
		{"GET", 2, -1},
		{"PUT", keyVal{12, 14}, true},
		{"GET", 12, 14},
		{"GET", 1, -1},
		{"PUT", keyVal{12, 34}, true},
		{"GET", 12, 34},
		{"PUT", keyVal{5, 18}, true},
		{"GET", 5, 18},
		{"PUT", keyVal{13, 17}, true},
		{"GET", 13, 17},
		{"PUT", keyVal{14, 57}, true},
		{"GET", 14, 57},
	},
}

// lRUAccess test processor
func lRUAccess(index int, t *testing.T) {
	for _, k := range lrucacheTests[index] {
		switch k.method {
		case "GET":
			if actual := cache3.Get(k.intput.(int)); actual != k.expected.(int) {
				t.Errorf("expected %d, actual %d for key %d", k.expected.(int), actual, k.intput.(int))
			}
		case "PUT":
			cache3.Put(k.intput.(keyVal).key, k.intput.(keyVal).val)
		default:
			t.Errorf("Wrong Method Name")
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
