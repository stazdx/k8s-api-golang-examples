#/bin/bash

docker build --no-cache -t $1 .
docker push $1

sed "s#DOCKER_IMAGE#$1#g" k8s-manifests.yaml | kubectl apply -f -

kubectl get po -w