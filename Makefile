registry=lwabish
time:=$(shell date +"%Y%m%d-%H%M%S")

.PHONY all: clean build image install clean

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin .

image:
	docker build -t $(registry)/scheduler:$(time) ./bin
	docker push $(registry)/scheduler:$(time)
	docker rmi $(registry)/scheduler:$(time)

install:
	kubectl config use-context home
	sed -i 's/appVersion: \S*/appVersion: $(time)/' ./chart/Chart.yaml
	helm upgrade -i -n lwabish-scheduler --create-namespace lwabish-scheduler ./chart
	git restore .

clean:
	rm -f bin/k8s-scheduler
