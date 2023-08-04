default:
	echo "hi there!"

start:	
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)
	docker compose -f "compose.yml" up -d --build