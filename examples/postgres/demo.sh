#!/bin/bash -x

kubectl apply -f postgres-sc.yaml
kubectl apply -f postgres-vol.yaml
kubectl apply -f postgres-app.yaml
