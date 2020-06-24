protoc:
	protoc -I/usr/local/include -I${GOPATH}/src:. \
	--go_out=plugins=grpc:. \
	--govalidators_out=. \
	proto/*.proto
