start:	
	make down
	make up

gen-pb:
	protoc --go_out=. --go-grpc_out=. broker/broker.proto

up:
	docker compose -f "compose.yml" up -d --build --remove-orphans

down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)

prep-deps:
	npm create vite@latest

watch: # this is still in testing
	docker compose alpha watch