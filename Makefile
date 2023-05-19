default: install

build:
	go build .

install: build
	go install .

testacc:
ifeq ($(NON_ENTERPRISE),1)
	TF_ACC=1 go test -count=1 -parallel=4 -v ./corellium/... -run ".*non_enterprise";
else
	TF_ACC=1 go test -count=1 -parallel=4 -v ./corellium/... -skip ".*non_enterprise"
endif

testexamples:
	go test -failfast -count=1 -v . -run "TestExamples.*"
