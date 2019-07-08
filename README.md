# Go-Kit Demo
Demo API using Go-Kit, Zipkin and Prometheus

### Installation & Running

 - `docker-machine create dev --virtualbox-memory 2048 --virtualbox-disk-size 10000`
 - `eval $(docker-machine env dev)`
 - `export COMPOSE_TLS_VERSION=TLSv1_2`
 - `docker-compose up`

To run migrations use:

 - `docker exec -it users migrate -source file://migrations -database postgres://postgres:@db:5432/postgres?sslmode=disable up`

To create a new migration file use:

 - `docker exec -it users migrate create -ext sql -dir migrations create_user`
