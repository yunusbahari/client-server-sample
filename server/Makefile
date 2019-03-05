
.PHONY: compile
VERSION = v1

compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w"

docker:
	docker build -t "asia.gcr.io/cloudeng-poc-staging-cicd/yunus/server:$(VERSION)" .

push:
	docker push "asia.gcr.io/cloudeng-poc-staging-cicd/yunus/server:$(VERSION)"

deploy:
	kubectl apply -f deployment.yml
