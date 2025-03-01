run: 
	@command -v docker >/dev/null 2>&1 || { \
		echo "Docker is not installed. Installing now..."; \
		sudo apt install docker.io; \
		echo "Docker installed !"; \
	}
	sudo docker build -t global-chat-backend:v0 .
	@read -p "Enter server port: " PORT; \
	read -p "Enter MongoDB host(strictly type ip): " MONGO_HOST; \
	echo "SERVER_PORT=$$PORT" > .env; \
	echo "DB_URL=mongodb://$$MONGO_HOST:27017/global-chat" >> .env; \
	sudo docker run -it --init --rm -p $$PORT:5000 --name gc-backend --env-file .env global-chat-backend:v0