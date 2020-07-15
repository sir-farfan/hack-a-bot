APP := bot

verify:
	golint
	go test ./...

docker: verify
	$(eval branch = $(shell git branch --show-current) )
	@echo current branch: $(branch)
	docker build . -t $(APP):$(branch)
	#docker image tag $(APP) $(APP):$(branch)
