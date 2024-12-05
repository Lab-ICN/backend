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

.PHONY: deployment deployment/destroy gomod/tidy
