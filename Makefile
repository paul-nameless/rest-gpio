build:
	env GOOS=linux GOARCH=arm GOARM=5 go build

deploy:
	rsync -r . pi:home

run:
	docker run --restart unless-stopped -d --device /dev/mem:/dev/mem --device /dev/gpiomem:/dev/gpiomem --network host home
