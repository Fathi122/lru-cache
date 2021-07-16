# simple LRU cache test
[![CITest](https://github.com/Fathi122/lru-cache/workflows/CITest/badge.svg)](https://github.com/Fathi122/lru-cache/actions)
[![codecov](https://codecov.io/gh/Fathi122/lru-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/Fathi122/lru-cache)

## Run test locally
```
$ go test -count=2 -v ./...
$ go test -cpu=1 -bench=. -run='^$' ./...
```
