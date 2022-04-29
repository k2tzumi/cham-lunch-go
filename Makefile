.DEFAULT_GOAL := up

.PHONY: up
up:
	@docker-compose up -d --build
	@docker container exec -it python3.dev.local zsh

down:
	@docker-compose down
