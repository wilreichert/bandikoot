.PHONY: bdkt python mariadb memcached rabbitmq keystone

all: bdkt python mariadb memcached python rabbitmq keystone

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
	docker rmi python:2.7-alpine3.5
	rm bkdt/bdkt
	rm keystone/Dockerfile
	rm glance/Dockerfile

build-%:
	docker build --rm -t $* --pull --build-arg TEST=asd -f $*/Dockerfile .
