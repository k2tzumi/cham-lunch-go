.PHONY: test
test:
	 @go test ./...

.PHONY: build
build:
	@go build ./...

.PHONY: update-json
update-json:
	@go run gen/update_area_json.go

src/area.json:
	@make update-json
