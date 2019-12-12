# hades

## usage
```shell script
docker build -t hades -f .docker/Dockerfile .
docker run -p 8081:80 -d --name hades hades
curl http://localhost:8081/ping
```
