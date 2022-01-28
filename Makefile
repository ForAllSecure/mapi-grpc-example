API = api/v1
IMAGE_NAME = mapi-grpc-example
CONTAINER_NAME= mapi-grpc-example

.PHONY: all
.DEFAULT: all
all: api-stubs docker-image

api-stubs:
	$(MAKE) -C $(API) all

.PHONY: docker-image
docker-image:
	docker build -f fuzz.Dockerfile . -t $(IMAGE_NAME)

.PHONY: swagger
swagger: docker-image
	docker run -t --rm -d --name $(CONTAINER_NAME)-tmp $(IMAGE_NAME)
	docker cp $(CONTAINER_NAME)-tmp:/opt/grpc-example/api/v1/example-api.swagger.json .
	docker rm -f $(CONTAINER_NAME)-tmp

.PHONY: run
run: swagger
	docker rm -f $(CONTAINER_NAME) || true
	docker run -it --rm -d -p 8081:8081 --name $(CONTAINER_NAME) $(IMAGE_NAME)

.PHONY: stop
stop:
	docker rm -f $(CONTAINER_NAME)

clean:
	$(MAKE) -C $(API) clean