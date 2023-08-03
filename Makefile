default:
	echo "hi there!"
start:	
	docker compose -f "compose.yml" down
	docker compose -f "compose.yml" up -d --build