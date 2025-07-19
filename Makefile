APP=golang-mock
push-production:
	@go mod vendor
	@docker buildx build --platform=linux/amd64 -t docker-regis.ottodigital.id/test/$(APP):latest .
	@docker image push docker-regis.ottodigital.id/test/$(APP):latest
	@rm -rf vendor