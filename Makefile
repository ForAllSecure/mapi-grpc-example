API = api/v1

.PHONY: all
.DEFAULT: all
all: api-stubs docker-image

api-stubs:
	$(MAKE) -C $(API) all

.PHONY: docker-image
docker-image:
	docker build -f fuzz.Dockerfile . -t mapi-grpc-example

.PHONY: swagger
swagger: docker-image
	docker run -t --rm -d --name mapi-grpc-example-tmp mapi-grpc-example
	docker cp mapi-grpc-example-tmp:/opt/grpc-example/api/v1/example-api.swagger.json .
	docker ps | grep mapi-grpc-example-tmp | cut -d" " -f1 | xargs docker rm -f

.PHONY: run
run: docker-image
	docker run -it --rm -d -p 8081:8081 --name mapi-grpc-example mapi-grpc-example

.PHONY: stop
stop:
	docker rm -f mapi-grpc-example

clean:
	$(MAKE) -C $(API) clean