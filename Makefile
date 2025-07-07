build:
	mkdir -p out
	go build -o ./out/ ./...

docker:
	docker build -t nocderechte/bridget .

clean:
	rm -r ./out

.PHONY: build docker clean
