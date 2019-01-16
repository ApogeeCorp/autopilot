#!/bin/bash -x

kubectl delete -f postgres-app.yaml
kubectl delete -f postgres-vol.yaml
kubectl delete -f postgres-sc.yaml
