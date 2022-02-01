# Go RBAC Examples

This repository contains a number of Go applications, each demonstrates a different approach to implementing
_Resource-Based Authorization_ (RBAC) for a simple HTTP API.

## Running the Examples

To run an example, `cd` to its directory and run `go run .`.

For example, to run the `casbin` example:

```sh
$ cd casbin
$ go run .
Staring server on 0.0.0.0:8080
```

You can now send requests using `curl`. For example:

```sh
$ curl -X POST -f -u dianet@acmecorp.com:asdfasdf http://localhost:8080/api/asset2
curl: (22) The requested URL returned error: 403

$ curl -X GET -f -u dianet@acmecorp.com:asdfasdf http://localhost:8080/api/asset2
Got permission⏎
```
