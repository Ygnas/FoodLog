.PHONY: apply
apply:
	kubectl create secret generic foodlog-credentials --from-file=foodlog-credentials.json | kubectl apply -k .