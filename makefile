build:
	go build -o ./bin/ksedit

install-goreleaser:
	curl -L https://github.com/goreleaser/goreleaser/releases/download/v0.123.3/goreleaser_Darwin_x86_64.tar.gz -o ./bin/goreleaser.tar.gz
	tar zxvf ./bin/goreleaser.tar.gz

# before release, cut tag
release:
	./bin/goreleaser --rm-dist
