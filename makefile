
DIR=`pwd`/pb
kafkadir=/Users/mac/go/src/github.com/Gonewithmyself/Road/docker/kafka
.PHONY:

clean:
	cd sql && ./build_sql.sh clean && ./build_sql.sh create
	cd ${kafkadir} && docker-compose stop && rm -rf kafka* zoo1 && docker-compose up -d
	docker exec -t redis redis-cli -n 1 flushdb	
	go clean
	
start:
	#cd ${DIR} && protoc -I . --go_out=plugins=grpc:. *.proto
	go build -o app
	./app