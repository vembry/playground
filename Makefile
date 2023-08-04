default:
	echo "hi there!"

start:	
	make tear-down
	docker compose -f "compose.yml" up -d --build

tear-down:
	docker compose -f "compose.yml" down
	(echo "y" | docker volume prune)