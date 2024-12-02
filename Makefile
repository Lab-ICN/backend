deployment:
	@kubectl apply \
		--filename k8s/ \
		--filename user-service/k8s/

deployment/destroy:
	@kubectl delete \
		--filename k8s/ \
		--filename user-service/k8s/

.PHONY: deployment deployment/destroy
