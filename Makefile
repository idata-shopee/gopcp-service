GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

init-dep:
	@dep init

status-dep:
	@dep status

dep:
	@dep ensure

update-dep:
	@dep ensure -update

test:
	go test -cover

test-only:
	go test -run $(CASE) -cover
