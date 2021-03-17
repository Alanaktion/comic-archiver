# Comic Archiver
*Go Edition*

This is a project to build simplistic archive tools for a large number of webcomics. I'm doing this mostly because I've had too many comics I love disappear with no reasonable way of finding them again.

Completed comics in particular are important to archive, because no matter how good they are and how much people love them, they disappear a lot more often than active ones.

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

## Disclaimer

This tool is written for archival purposes only, and as long as the comics are online you should make every effort to read them on their official websites! Support the artists wherever you can, including on Patreon, via merch, and using their sites without ad blockers. They depend on every bit of that income to keep creating their incredible comics, and deserve all the attention and support we can give them.

If the worst should happen though, running these scripts in advance before a comic disappears can keep the creator's work alive when they're not able to keep their original site online anymore. I recommend not publicly listing any of the archived strips as long as the comic is still online in some official form, but once one does go offline, it seems right to me to make your archives available to fans looking for the comic. If you're not sure, you could always try contacting the artist to make sure it's okay with them!
