build:
	go build -o ./bin/ksedit

install-goreleaser:
	mkdir -p ./bin
	curl -L https://github.com/goreleaser/goreleaser/releases/download/v0.123.3/goreleaser_Darwin_x86_64.tar.gz -o ./bin/goreleaser.tar.gz
	cd ./bin/; tar zxvf ./goreleaser.tar.gz

# before release, cut tag. $ git tag vX.X.X
release:
	./bin/goreleaser --rm-dist
