# Comic Archiver
*Go Edition*

I'm rewriting some of my comic archive scripts in Go because I can.

## Requirements

A recent (1.8+ probably) Go install. No external packages are used, so no need to put it in GOPATH or install anything extra.

## Running

You can either run the scripts directly, for example:

```bash
go run doa.go
```

Or you can build and run them:

```bash
go build doa.go
./doa
```

In either case, the comics will be downloaded to `./comics/<comic name>`.
