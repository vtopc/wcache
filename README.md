# wcache - an in-memory cache with expiration

[![Godoc Reference][godoc-img]][godoc]

Implements a cache with TTL(expiration).
`expireFn` will be called when record is expired.

### Features
* Thread-safe.
* Individual expiring time or global expiring time, you can choose.
* Can trigger a custom callback on key expiration(could be used as Pub/Sub with aggregation).
* Can trigger a custom callback on key collisions(e.g. for aggregating metrics). 
By default, will overwrite value.
* Graceful shutdown(using context). Will call expiration callback for all records,
ignoring they TTL.

### TODO
* Auto-Extending expiration on `Get` and/or `Set`.
* Benchmarks.

[godoc]: https://godoc.org/github.com/vtopc/wcache
[godoc-img]: https://godoc.org/github.com/vtopc/wcache?status.svg
