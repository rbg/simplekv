GO_FILES := $(wildcard *.go) $(wildcard cmd/*.go) $(wildcard api/*.go) $(wildcard store/*.go)  

skv: $(GO_FILES)
	go build -o skv

container: 
	docker build --force-rm=true -t rbg/skv:latest .
