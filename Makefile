default:
	echo "hi there!"

start:	
	make down
	make up

up:
	docker compose -f "compose.yml" up -d --build --remove-orphans

down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)

watch:
	docker compose alpha watch