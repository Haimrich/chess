# Chess
An online chess platform

Final Course Project UniCT 2022 

Advanced Programming Languages | Distributed Systems and Big Data

#### Technologies

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![C#](https://img.shields.io/badge/c%23-%23239120.svg?style=flat&logo=c-sharp&logoColor=white)
![C++](https://img.shields.io/badge/c++-%2300599C.svg?style=flat&logo=c%2B%2B&logoColor=white)
![Python](https://img.shields.io/badge/python-3670A0?style=flat&logo=python&logoColor=white)
![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?styleflat&logo=css3&logoColor=white)
![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=flat&logo=html5&logoColor=white)
![.Net](https://img.shields.io/badge/.NET-5C2D91?style=flat&logo=.net&logoColor=white)
![Blazor](https://img.shields.io/badge/Blazor-512BD4?style=flat&logo=blazor&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=flat&logo=mongodb&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=flat&logo=JSON%20web%20tokens)
![Nginx](https://img.shields.io/badge/nginx-%23009639.svg?style=flat&logo=nginx&logoColor=white)
![Kafka](https://img.shields.io/badge/kafka-%23231F20.svg?style=flat&logo=apachekafka&logoColor=white)
![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=flat&logo=grafana&logoColor=white)
![Prometheus](https://img.shields.io/badge/prometheus-%23E6522C.svg?style=flat&logo=prometheus&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=flat&logo=kubernetes&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)

## Description
- [Backend](https://github.com/Haimrich/chess/tree/main/backend) written in Go using [Gin](https://github.com/gin-gonic/gin) web framework and [Gorilla](https://github.com/gorilla/websocket) websockets. This monolithic component has been splitted in the following microservices:
  - [User Service](https://github.com/Haimrich/chess/tree/main/user) handles user signup, login and authentication. It also generates JWT Tokens for authorization in other microservices.
  - [WebSocket Node](https://github.com/Haimrich/chess/tree/main/wsnode) handles WebSocket connections with clients.
  - [Dispatcher](https://github.com/Haimrich/chess/tree/main/dispatcher) routes messages from WS Nodes to Game and Challenge Services.
  - [Challenge Service](https://github.com/Haimrich/chess/tree/main/dispatcher) handles challenge sending and accepting.
  - [Game Service](https://github.com/Haimrich/chess/tree/main/dispatcher) handles games status, timer, legal moves, etc.
- [Frontend](https://github.com/Haimrich/chess/tree/main/frontend) written in C# using [.NET Blazor WebAssembly](https://dotnet.microsoft.com/en-us/apps/aspnet/web-apps/blazor)
- [Engine](https://github.com/Haimrich/chess/tree/main/engine) written in C++ using [Pistache](https://github.com/pistacheio/pistache) HTTP server. [Bitboards](https://www.chessprogramming.org/Bitboards) were adopted for board representation while search logic was inspired by [carnatus](https://github.com/zserge/carnatus) and [sunfish](https://github.com/thomasahle/sunfish).

## Usage

### Docker compose

- Monolitic backend version

  ```
  docker-compose -f docker/docker-compose.yml up
  ```
- Microservices and kafka version

  ```
  docker-compose -f docker/docker-compose.kafka.yml up
  ```

### Minikube
Start minikube and install addons
```shell 
minikube start
eval $(minikube -p minikube docker-env)
minikube addons enable ingress
minikube addons enable metrics-server
```
Build docker images
```shell 
docker build -t chess_frontend ./frontend
docker build -t chess_user ./user
docker build -t chess_wsnode ./wsnode
docker build -t chess_game ./game
docker build -t chess_engine ./engine
docker build -t chess_challenge ./challenge
docker build -t chess_dispatcher ./dispatcher
docker build -t chess_forecaster ./forecaster
docker build -t chess_loadgenerator ./load_generator
```
Add Prometheus, requires [Helm](https://helm.sh/)
```shell 
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack
```
Add [Strimzi](https://strimzi.io/)
```shell 
kubectl apply -f 'https://strimzi.io/install/latest?namespace=default' -n default
```
Deploy
```shell 
kubectl apply -f k8s
# To disable the load generator
# kubectl delete -f k8s/60-deployment-loadgenerator.yaml
```
Append entries in /etc/hosts
```shell 
echo "$(minikube ip) chess.example grafana.chess.example prometheus.chess.example" | sudo tee -a /etc/hosts
```

## Screenshots
<img src="https://user-images.githubusercontent.com/7826610/150849609-abf87a14-e959-440a-b1b9-d5a8e7fadacb.PNG" height="200"> <img src="https://user-images.githubusercontent.com/7826610/150849995-91c326fc-e26b-4bed-8185-f898219bc3e8.png" height="200">
