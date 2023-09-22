default:
	echo "hi there!"

start:	
	make down
	make up

up:
	docker compose -f "compose.deps.yml" -f "compose.apps.yml" -f "compose.tools.yml" up -d --build

down:
	docker compose -f "compose.deps.yml" -f "compose.apps.yml" -f "compose.tools.yml" down
	(echo "y" | docker volume prune)

watch:
	docker compose alpha watch