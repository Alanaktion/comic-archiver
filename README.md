# Comic Archiver
*Go Edition*

I'm rewriting some of my comic archive scripts in Go because I can.

## Requirements

A recent Go install. You should properly configure GOPATH to build it.

## Building

```bash
go build .
```

## Running

Build the app, then run the resulting binary:

```bash
./comic-archiver <action> [args ...]
```

For example to archive _xkcd_ and _Whomp!_:

```bash
./comic-archiver archive xkcd whomp
```

Comics will be downloaded to `./comics/<comic name>`.

Running the `help` action will list all available actions and archivers:

```bash
./comic-archiver help
```
