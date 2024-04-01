.PHONY: apply
apply:
	kubectl create secret generic foodlog-credentials --from-file=foodlog-credentials.json | kubectl apply -k .

.PHONY: delete
delete:
	kubectl delete secret foodlog-credentials | kubectl delete -k .

.PHONY: update
update:
	kubectl rollout restart deployment foodlog-deployment