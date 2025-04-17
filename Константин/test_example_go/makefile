.SILENT:

run:
	go run cmd/main.go

docker-build:
	sudo docker build -t manage-sys .

docker-run:
	sudo docker run --rm --network=host --name manageSys manage-sys