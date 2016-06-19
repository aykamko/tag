SHELL:=/bin/bash

GOBUILD=go build -v github.com/aykamko/tag >/dev/null

generate_release_binaries:
	mkdir -p release; \
	cd release; \
	export GOOS=darwin; export GOARCH=386; \
	${GOBUILD} && zip tag_$${GOOS}_$${GOARCH} tag; \
	export GOOS=darwin; export GOARCH=amd64; \
	${GOBUILD} && zip tag_$${GOOS}_$${GOARCH} tag; \
	export GOOS=linux; export GOARCH=386; \
	${GOBUILD} && tar -cvzf tag_$${GOOS}_$${GOARCH}.tar.gz tag; \
	export GOOS=linux; export GOARCH=amd64; \
	${GOBUILD} && tar -cvzf tag_$${GOOS}_$${GOARCH}.tar.gz tag; \
	export GOOS=linux; export GOARCH=arm; \
	${GOBUILD} && tar -cvzf tag_$${GOOS}_$${GOARCH}.tar.gz tag; \
	export GOOS=windows; export GOARCH=386; \
	${GOBUILD} && tar -cvzf tag_$${GOOS}_$${GOARCH}.tar.gz tag; \
	export GOOS=windows; export GOARCH=amd64; \
	${GOBUILD} && tar -cvzf tag_$${GOOS}_$${GOARCH}.tar.gz tag;
