default: run

run:
	go run cmd/ltrain.go

protoc:
	protoc --go_out=. transit_realtime/gtfs-realtime.proto
	cd nyct_subway && protoc --go_out=. nyct-subway.proto && cd ..
	sed -i'' -e 's/import transit_realtime "."/import "github.com\/aoeu\/mta\/transit_realtime"/' nyct_subway/nyct-subway.pb.go
