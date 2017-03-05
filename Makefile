.PHONY: bdkt python mariadb memcached rabbitmq keystone

all: bootstrap build

bootstrap: bdkt	python

bdkt:
	cd bdkt && go build bdkt.go

python:
	docker build --rm -t python:2.7-alpine3.5 --pull -f python/Dockerfile .

build: keystone

cinder: build-cinder

glance: build-glance

heat: build-heat

keystone: build-keystone

horizon: build-horizon

mariadb: build-mariadb

memcached: build-memcached

neutron: build-neutron

nova: build-nova

rabbitmq: build-rabbitmq

clean:
	docker rmi python:2.7-alpine3.5 | true
	rm -f bdkt/bdkt
	rm -f {keystone,glance}/Dockerfile

build-%:
	bdkt/bdkt -config config.yaml -service $*
	docker build --rm -t $* --build-arg TEST=asd -f $*/Dockerfile .
