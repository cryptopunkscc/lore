default: build

build:
		go build ./cmd/lored
		go build ./cmd/lore-cli
		go build ./cmd/loreid


install:
		go install ./cmd/lored
		go install ./cmd/lore-cli
		go install ./cmd/loreid

clean:
		rm lored lore-cli