#https://github.com/darshanadinushal/rabbitmq-cluster/tree/master/deployment/rabbitmq-cluster
kubectl create ns rabbits
kubectl apply -n rabbits -f configmap.yaml
kubectl apply -n rabbits -f rbac.yaml
kubectl apply -n rabbits -f cookie_secret.yaml
kubectl apply -n rabbits -f service.yaml
kubectl apply -n rabbits -f statefulset.yaml
