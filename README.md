# hades

## usage
```shell script
# prerequisites
docker --version
curl --version
jq --version
# build
docker build -t hades .
#prepare
docker network create --driver bridge hades-net
# run
docker run --rm --name hades-db -p 3306:3306 --network hades-net\
 -e MYSQL_ROOT_PASSWORD='root'\
 -e MYSQL_DATABASE='hades'\
 -e MYSQL_USER='user'\
 -e MYSQL_PASSWORD='pass'\
 -d mysql
docker logs -f hades-db
docker run --rm --name hades -p 8081:80 --network hades-net\
 -e JWT_KEY='secret'\
 -e DB_DSN='root:root@tcp(hades-db:3306)/hades'\
 -e APP_PORT='80'\
 -d hades
# test
curl -s http://localhost:8081/ping
token=$(curl -s http://localhost:8081/authenticate --data '{"username":"user","password":"pass"}' | tee $(tty) | jq -r .data.token)
curl -s http://localhost:8081/me --header "Authorization: Bearer ${token}"
# cleanup
docker stop hades hades-db
docker network rm hades-net
```
