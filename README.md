# Twitch API with Go 🐹

Go to my Twitch at MajesticCodingTwitch! This Go site uses templates with the Twitch API.

![Image of Gopher Interviwing](https://github.com/smithlabs/github-assets/blob/main/web/dancing-gopher-hello-world.gif?raw=true)

## Twitch

🎥 **Twitch**: Watch my coding Streams [Twitch channel](https://www.twitch.tv/majesticcodingtwitch).

## Run with Kubernetes 📦
```
kubectl create configmap majesticcodingtwitch-env --from-env-file=.env
kubectl apply -f k8s-go.yaml
kubectl expose deployment majesticcodingtwitch --type=LoadBalancer --name=majesticcodingtwitch --port=8080 --target-port=8080
```

## Run with Docker 🐳
```
docker compose up
```

## Run with Go 🐹
```
go mod tidy
go run main.go
```
