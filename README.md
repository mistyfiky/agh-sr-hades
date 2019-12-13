# hades

## usage
```shell script
docker run --name hades-db -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=hades -d mysql
docker build -t hades -f .docker/Dockerfile .                                             
docker run --name hades --network host -e JWT_KEY=secret -e DB_DSN=root:root@localhost/hades -d hades
curl http://localhost:8081/ping
```
