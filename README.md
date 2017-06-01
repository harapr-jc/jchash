# jchash - Crypto Hashing Service

## Summary

Provides an executable (jchash) to run a server that services requests for generating
the cryptographical hash for a given password.


## Dependencies

Runs on Linux.
Depends on Go standard library and OS commands uuidgen and grep.
Depends on package `harapr-jc/hashgen` for server and utilties for reuse.

### Starting the `jchash` server

```
./jchash --host <hostname> --port <port number>
```

For a help message, use `./jchash -h`.

### Endpoints

* POST /hash to submit a job

Example:
```
$ curl --data "password=angryMonkey" http://localhost:8080/hash
48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
```
Submits a job.
Returns a unique job id that can be used to fetch the crypto hash for the given password.

* GET /hash/{id}

Returns result of given job if available. Result is base64 encoded crypto hash.
Example:
```
$ curl http://localhost:8080/hash/48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
ZEHhWB65gUlzdVwtDQArEyx-KVLzp_aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A-gf7Q==
```

* GET /stats

Returns statistics for all GET and POST hash requests for current server instance.

Example:
```
$ curl http://localhost:8080/stats
{"total":2001,"average":2543630,"units":"ns"}
```

### Adding Salt

To reduce the risk of brute-force attacks on crypto hashes, each request can ask for the addition
of a random salt during hashing. Use the option 'salt=yes' to add a random salt during hashing.
The random salt is stored persistently for the request (user account).

```
$ curl --data "password=angryMonkey;salt=yes" http://localhost:8080/hash'
```

When salt is used, two distinct requests using the same password value will yield different crypto hash bytes.

### Stopping the `jchash` server

The server supports a graceful shutdown. Simply issue a SIGINT to the process.
All running jobs will complete before shutdown.
No additional job requests can start while shutdown is pending.
Example:
```
2017/06/01 12:58:39 Starting crypto hash server at "localhost:8080"
2017/06/01 12:58:47 getting record with uuid: 48765c95-9ea8-43bd-bc0b-5ecb6897f2b9
^C2017/06/01 12:58:58 Shutdown requested...
2017/06/01 12:58:58 Done waiting
2017/06/01 12:58:58 Shutting down server...
2017/06/01 12:58:58 Server shutdown complete. Bye!
```

# Development

## Testing

The `hashgen` package has tests.
```
$ cd hashgen
$ go test -v
```

### Code Coverage

Sample coverage check:

```
$ cd hashgen
$ go test -coverprofile=cover.out
...
$ go tool cover -func=cover.out
github.com/harapr-jc/hashgen/dao.go:24:		New			100.0%
github.com/harapr-jc/hashgen/dao.go:29:		Append			71.4%
github.com/harapr-jc/hashgen/dao.go:65:		Get			90.9%
github.com/harapr-jc/hashgen/hasher.go:15:	getSalt			80.0%
github.com/harapr-jc/hashgen/hasher.go:26:	getCryptoHash		100.0%
github.com/harapr-jc/hashgen/lru.go:39:		NewCache		100.0%
github.com/harapr-jc/hashgen/lru.go:53:		Add			77.8%
github.com/harapr-jc/hashgen/lru.go:72:		Get			87.5%
github.com/harapr-jc/hashgen/server.go:51:	HandleGetHashRequest	100.0%
github.com/harapr-jc/hashgen/server.go:73:	HandleHashRequest	93.1%
github.com/harapr-jc/hashgen/server.go:133:	HandleStatsRequest	100.0%
github.com/harapr-jc/hashgen/server.go:138:	StartServer		90.3%
github.com/harapr-jc/hashgen/stats.go:38:	Accumulate		72.7%
github.com/harapr-jc/hashgen/stats.go:59:	GetJson			85.7%
github.com/harapr-jc/hashgen/uuid.go:11:	getUuid			75.0%
total:						(statements)		86.8%
```

### Benchmarking

Some of the tests support benchmarks.

```
$ cd hashgen
$ go test -bench=".*"
```

## Implementation Choices

### Job Id
The use of uuid instead of monotonically increasing id for a job identifier allows server
restart without having to persist the next id. Furthermore, no synchronization is required
on next id.

A drawback is that, without a database back end, ordering of uuid is less efficient for
persistent storage. A consequence is that records that are evicted from the cache will
take longer to find.

### LRU cache
The purpose for using the cache is to put an upper bound on memory usage for the server.

The design intent is that n most recently used jobs are held in cache, all are passed thru
to persistent storage. On a cache miss, job data could be fetched from persistent storage.
The default backing file is called 'backup.json'. One would prefer to use a performant data
store instead. The Go standard library has sql API, but no database drivers are included.

## Outstanding Development Items

* Raise code coverage
* The LRU cache is mock, it doesn't really evict yet
* Add command line option for backing file
* Add guard against running the executable twice on same machine
