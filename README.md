lore
====

An experimental decentralized immutable data graph.

## Status

Highly experimental code. Nobody should use it for anything.

## Install

### Dependencies & requirements

Go is required to build it, ffplay is required to play media.

### Download and install

    $ go get github.com/cryptopunkscc/lore
    $ cd $GOPATH/src/github.com/cryptopunkscc/lore
    $ make install

This will use `go install` to install two binaries - `lored` and `lore-cli`.

## Usage

Start the daemon first:

    $ lored

Then you can use the cli tool to control the daemon:

    $ lore-cli add <path>        # add a local file to shared files
    $ lore-cli list              # list IDs of locally shared files
    $ lore-cli addsource <url>   # add an external lored as a data source
    $ lore-cli listsources       # list data sources
    $ lore-cli play <id>         # find a file with provided ID and play it locally

Uses ffplay to play files. Make sure it's in your $PATH.