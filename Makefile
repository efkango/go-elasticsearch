runES:
	docker run -d --name es01 -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -it docker.elastic.co/elasticsearch/elasticsearch:7.14.0
startES:
	docker start es01
stopES:
	docker stop es01
build:
	go build main.go employee.go client.go
run: build
	./main
