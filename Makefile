default: build

build:
		go build ./cmd/lored
		go build ./cmd/lore-cli

install:
		go install ./cmd/lored
		go install ./cmd/lore-cli

clean:
		rm lored lore-cli