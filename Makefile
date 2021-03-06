TAG ?= "latest"

build:
	GOOS=linux go build -o say-hi-linux
	docker build --rm --pull -t say-hi .

deploy: build
	docker tag say-hi bshelton229/say-hi:$(TAG)
	docker tag say-hi quay.io/bshelton229/say-hi:$(TAG)

	docker push bshelton229/say-hi:$(TAG)
	docker push quay.io/bshelton229/say-hi:$(TAG)

run:
	docker run --rm -it -P say-hi:latest

shell:
	docker run --rm -it -P say-hi:latest bash
