# wcache

[![License][lic-img]][lic-url]
[![Godoc Reference][godoc-img]][godoc-url]

Implements in-memory [cache with write-back strategy](https://en.wikipedia.org/wiki/Cache_(computing)#Writing_policies)
(a cache with expiration).
An expiration callback(`expireFn`) will be called when record is expired.

### Features
* Thread-safe.
* Per-key or global TTL.
* Can trigger a custom callback on key expiration(could be used as Pub/Sub with aggregation).
* Can trigger a custom callback on key collisions(e.g. for aggregating metrics). 
By default, will overwrite value.
* Graceful shutdown(using context). Will call expiration callback for all records,
ignoring they TTL.

### Install
`go get github.com/vtopc/wcache`

### TODO
* Optional auto-extending expiration on `Get` and/or `Set`.

[godoc-url]: https://godoc.org/github.com/vtopc/wcache
[godoc-img]: https://godoc.org/github.com/vtopc/wcache?status.svg

[lic-url]: https://github.com/vtopc/wcache/blob/master/LICENSE
[lic-img]: http://img.shields.io/badge/license-MIT-red.svg?style=flat
