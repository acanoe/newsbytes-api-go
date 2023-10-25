.PHONY: clean
clean:
	@rm -f ./sources/progscrape.so
	@rm -f db.sqlite

build-sources:
	go build -buildmode=plugin -o ./sources/progscrape.so ./sources/progscrape.go

.PHONY: run
run:
	go run main.go

run-local: clean build-sources run