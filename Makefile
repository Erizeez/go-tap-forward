.PHONY: build
build:
	go build -o bin/test

.PHONY: run
run: build
	sudo ./bin/test

.PHONY: adjust
adjust:
	sudo ip neigh change 192.168.52.2 lladdr 00:11:22:33:44:55 dev tap-test