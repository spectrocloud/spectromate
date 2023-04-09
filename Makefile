start:
	@echo "Starting the application and ngrok..."
	@docker compose up --detach
	@ngrok http 3000 --log-level=debug > /dev/null & sleep 2
	@echo "Ngrok URL: $$(curl -s localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url')"


stop:
	@echo "Stopping the application and ngrok"
	docker compose down
	pkill ngrok
