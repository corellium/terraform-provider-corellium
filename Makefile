default: install

build:
	go build .

install: build
	go install .

test:
	go test -count=1 -parallel=4 ./...

testacc:
ifeq ($(NON_ENTERPRISE),1)
	TF_ACC=1 go test -count=1 -parallel=4 -v ./... -run ".*non_enterprise";
else
	TF_ACC=1 go test -count=1 -parallel=4 -v ./... -skip ".*non_enterprise"
endif
