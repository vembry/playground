default:
	echo "hi there!"
start:	
	docker build --pull --rm -f "api/Dockerfile" -t credits-api:local "api"
	docker build --pull --rm -f "load-test/Dockerfile" -t credits-load-test:local "load-test"
	docker compose -f "compose.yml" down
	docker compose -f "compose.yml" up -d --build