build:
	env GOOS=linux GOARCH=arm GOARM=5 go build

docker-build:
	docker build -t home:latest .

sync:
	rsync -r . pi:home

run:
	docker run --name home --restart unless-stopped -d --device /dev/mem:/dev/mem --device /dev/gpiomem:/dev/gpiomem --network host home
