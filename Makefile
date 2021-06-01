build:
	go build ./...

test:
	go test -race ./...

patch:
	./script_patch.sh