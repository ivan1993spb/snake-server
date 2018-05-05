IMAGE=ivan1993spb/snake-server:latest

docker/build:
	@docker build -t $(IMAGE) .

docker/push:
	@docker push $(IMAGE)
