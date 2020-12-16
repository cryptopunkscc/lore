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

    $ lore-cli create            # reads and save data from stdin and returns the id
    $ lore-cli list              # list all available files
    $ lore-cli read <id>         # streams the file to stdout
    $ lore-cli play <id>         # streams the file to ffplay
    $ lore-cli delete <id>       # deletes the file

Uses ffplay to play files. Make sure it's in your $PATH.