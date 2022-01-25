# Chess
An online chess platform

Final Course Project UniCT 2022  
Advanced Programming Languages  
Distributed Systems and Big Data

#### Technologies

![Go](https://img.shields.io/badge/Go-white.svg?style=flat&logo=go&logoColor=00ADD8)
![C#](https://img.shields.io/badge/C%23-white.svg?style=flat&logo=c-sharp&logoColor=239120)
![C++](https://img.shields.io/badge/C++-white.svg?style=flat&logo=c%2B%2B&logoColor=00599C)
![Python](https://img.shields.io/badge/Python-white?style=flat&logo=python&logoColor=3670A0)
![CSS3](https://img.shields.io/badge/CSS3-white.svg?styleflat&logo=css3&logoColor=1572B6)
![HTML5](https://img.shields.io/badge/HTML5-white.svg?style=flat&logo=html5&logoColor=E34F26)  
![.Net](https://img.shields.io/badge/.NET-white?style=flat&logo=.net&logoColor=5C2D91)
![Blazor](https://img.shields.io/badge/Blazor-white?style=flat&logo=blazor&logoColor=512BD4)
![MongoDB](https://img.shields.io/badge/MongoDB-white.svg?style=flat&logo=mongodb&logoColor=4ea94b)
![Nginx](https://img.shields.io/badge/Nginx-white.svg?style=flat&logo=nginx&logoColor=009639)
![Kafka](https://img.shields.io/badge/Kafka-white.svg?style=flat&logo=apachekafka&logoColor=231F20)
![JWT](https://img.shields.io/badge/JWT-white?style=flat&logo=JSON%20web%20tokens&logoColor=black)  
![Grafana](https://img.shields.io/badge/Grafana-white.svg?style=flat&logo=grafana&logoColor=F46800)
![Prometheus](https://img.shields.io/badge/Prometheus-white.svg?style=flat&logo=prometheus&logoColor=E6522C)
![Kubernetes](https://img.shields.io/badge/Kubernetes-white.svg?style=flat&logo=kubernetes&logoColor=326ce5)
![Docker](https://img.shields.io/badge/Docker-white.svg?style=flat&logo=docker&logoColor=0db7ed)

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

- Monolithic backend version

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
helm install prometheus prometheus-community/kube-prometheus-stack --set "grafana.plugins={grafana-simple-json-datasource,grafana-piechart-panel}"
```
Add [Strimzi](https://strimzi.io/)
```shell 
kubectl apply -f 'https://strimzi.io/install/latest?namespace=default' -n default
```
Deploy
```shell 
kubectl apply -f k8s
```
(Optional) Disable the load generator
```shell 
kubectl delete -f k8s/60-deployment-loadgenerator.yaml
```
Append entries in /etc/hosts
```shell 
echo "$(minikube ip) chess.example grafana.chess.example prometheus.chess.example" | sudo tee -a /etc/hosts
```
You can try the application at <http://chess.example/>.  
The Grafana dashboard will be available at <http://grafana.chess.example/d/Iahds4b7k>.

## Screenshots
<img src="https://user-images.githubusercontent.com/7826610/150849609-abf87a14-e959-440a-b1b9-d5a8e7fadacb.PNG" height="200"> <img src="https://user-images.githubusercontent.com/7826610/150849995-91c326fc-e26b-4bed-8185-f898219bc3e8.png" height="200">
