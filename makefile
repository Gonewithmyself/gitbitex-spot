
DIR=`pwd`/pb

.PHONY:

clean:
	rm -rf ${DIR}/*.go

g:
	cd ${DIR} && protoc -I . --go_out=plugins=grpc:. *.proto