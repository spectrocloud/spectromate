.PHONY: license

start:
	@echo "Starting the application and ngrok..."
	@docker compose up --detach
	@ngrok http 3000 --log-level=debug > /dev/null & sleep 2
	@echo "Ngrok URL: $$(curl -s localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url')"


stop:
	@echo "Stopping the application and ngrok"
	docker compose down
	pkill ngrok

server:
	@echo "Starting the server..."
	go run main.go

init: 
	@echo "Installing dependencies..."
	go mod download
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/google/go-licenses@latest
	~/go/bin/mockgen -source=internal/cache.go -destination=mock/mock_cache.go -package=mock

opensource:
	@echo "Checking license headers..."
	~/go/bin/go-licenses report spectrocloud.com/spectromate --template=docs/open-source.tpl > docs/open-source.md 


license:
	@echo "Applying license headers..."
	 copywrite headers	