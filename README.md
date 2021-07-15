# simple LRU cache test
[![Actions Status](https://github.com/Fathi122/lru-cache/workflows/build/badge.svg)](https://github.com/Fathi122/lru-cache/actions)
[![codecov](https://codecov.io/gh/Fathi122/lru-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/Fathi122/lru-cache)
[![Known Vulnerabilities](https://snyk.io/test/github/Fathi122/lru-cache/badge.svg)](https://snyk.io/test/github/Fathi122/lru-cache)

## Run test locally
```
$ go test -count=2 -v ./...
$ go test -cpu=1 -bench=. -run='^$' ./...
```