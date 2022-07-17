.PHONY: test
test:
	 @go test ./...

.PHONY: build
build:
	@go build ./...

.PHONY: netlify
netlify:
	mkdir -p functions
	go get ./...
	go build -o functions/jma-lambda ./src/main.go

.PHONY: update-json
update-json:
	@go run gen/update_area_json.go

src/area.json:
	@make update-json
