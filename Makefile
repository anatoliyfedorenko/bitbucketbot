TARGET=bitbucketbot

all: clean fmt build

clean:
	rm -rf $(TARGET)

build:
	go build -o $(TARGET) main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o $(TARGET) main.go

build_docker:
	docker build -t bitbucketbot .

docker: build_linux build_docker
