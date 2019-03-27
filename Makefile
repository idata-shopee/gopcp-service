GOPATH := $(shell cd ../../../.. && pwd)
export GOPATH

init-dep:
	@dep init

status-dep:
	@dep status

update-dep:
	@dep ensure -update

test:
	go test -cover

restore:
	godep restore -v

test-only:
	go test -run $(CASE) -cover

save:
	godep save

update:
	@godep update github.com/...
