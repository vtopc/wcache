# wcache

[![Godoc Reference][godoc-img]][godoc]

Write cache with delayed sync.

Implements a cache with TTL(expiration).
`expireFn` will be called when record is expired.

Context could be used for flush during graceful shutdown.

[godoc]: https://godoc.org/github.com/vtopc/wcache
[godoc-img]: https://godoc.org/github.com/vtopc/wcache?status.svg
