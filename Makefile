.PHONY: apply
apply:
	kubectl create configmap foodlog-credentials --from-file=foodlog-credentials.json | kubectl apply -k .
