deployment:
	@kubectl apply \
		--filename k8s/ \
		--filename user-service/k8s/

deployment/destroy:
	@kubectl delete \
		--filename k8s/ \
		--filename user-service/k8s/

gomod/tidy:
	@cd user-service && go mod tidy
	@cd token-service && go mod tidy

compose:
	@docker compose --file infra/compose.dev.yaml \
		up --detach

compose/down:
	@docker compose --file infra/compose.dev.yaml down

compose/fresh:
	@docker compose --file infra/compose.dev.yaml \
		up --detach --build

docker/prune:
	@docker system prune --force

docker/clean:
	@docker rm --force $$(docker ps --all --quiet)

.PHONY: deployment deployment/destroy gomod/tidy compose compose/down
.PHONY: compose/fresh docker/prune docker/clean
