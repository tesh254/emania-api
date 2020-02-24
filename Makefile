run:
	@echo ":::: App is startin up ::::"
	@echo "CONFIG::  😁 Exporting environemnt variables"
	# Parrot os source alternative
	/bin/sh .env
	@echo "SUCCESS:  ✔ Environment variables exported"
	@echo "INIT::::  ⚡ Running server"
	go run main.go
docker:
	@echo ":::: Starting Container in Bash ::::"
	@echo "🐋 Loading...."
	docker exec -it emania-api bash
docs:
	@echo ":::: Swagger Editor Docs ::::"
	@echo "📃 Loading..."
	docker run -d -p 3500:4000 swagger-editor
run-2:
	@echo ":::: App is starting up ::::"
	@echo "Config:: 😁 Exporting environemnt variables"
	/bin/sh .env
	@echo "SUCCESS: ✔ Environment variables exported"
	@echo "INIT:::: ⚡ Running server"
	./main
deploy:
	@echo ":::: Deploying application ::::"
	export GO111MODULE=on
	gcloud app deploy app.yaml