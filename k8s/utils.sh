
# helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
# prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack
kubectl port-forward deployment/prometheus-grafana 3000
kubectl port-forward pod/prometheus-prometheus-kube-prometheus-prometheus-0 9090
# serve per l'horizontal pod autoscaler
#helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
#helm upgrade --install metrics-server metrics-server/metrics-server


# minikube start --memory=4096 --vm-driver=kvm2
# eval $(minikube -p minikube docker-env)

# minikube addons enable ingress
# minikube addons enable metrics-server
# minikube addons enable ingress-dns

# kubectl apply -f 'https://strimzi.io/install/latest?namespace=default' -n default

# echo $(minikube ip)


