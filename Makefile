build: docker-native-inspect

docker-build:
	docker build -t docker-native-inspect .
	docker run --rm docker-native-inspect tar czf - docker-native-inspect | tar vxzf -

clean:
	rm -f docker-native-inspect*

docker-native-inspect: *.go
	go build -o docker-native-inspect -ldflags "-s"

.PHONY: build docker-build clean
