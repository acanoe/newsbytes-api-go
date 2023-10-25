clean:
	@rm -f ./sources/progscrape.so
	@rm -f db.sqlite

build-sources:
	@go build -buildmode=plugin -o ./sources/progscrape.so ./sources/progscrape.go

run:
	@go run main.go

run-local: build-sources run