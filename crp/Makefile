.PHONY: build

build:
	go build -o main .

.PHONY: docker

docker:
	@echo "Building Docker Image:"
	docker image build -f Dockerfile -t forum-image .
	@echo

	@echo "List of images:"
	docker images
	@echo

	@echo "Initiating Container:"
	docker container run -t -p 8080:8080 --detach --name forum-container forum-image
	@echo

	@echo "Running command:"
	docker exec -it forum-container ls -la
	@echo

	@echo "Running server:"
	docker exec -it forum-container ./main
	@echo

.PHONY: clean

clean:
	@echo "Stopping container:"
	docker stop forum-container
	@echo

	@echo "Removing container:"
	docker rm forum-container
	@echo

	@echo "Deleting images:"
	docker rmi -f forum-image
	@echo

	@echo "List of images and containers now:"
	docker ps -a
	@echo
	docker images
	@echo

	rm -rf main

.DEFAULT_GOAL := build