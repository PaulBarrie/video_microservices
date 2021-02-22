include .env
export $(shell sed 's/=.*//' .env)

restart: 
ifdef svc
	make rm svc="$(svc)" && make up svc="$(svc)"
else
	docker-compose restart 
endif

rebuild:
	 docker-compose up --build --force-recreate --no-deps $1

ps:
	docker-compose ps

up:
ifdef svc
	docker-compose up -d --build $(svc)
else
	docker-compose up -d --build
endif

all: apis search 

app:
	docker-compose up -d --build app

apis:
	docker-compose up -d --build minio smtp api video_encoder

db:
	docker-compose up -d msql

search: db
	docker-compose up -d --build logstash elasticsearch 

kafka_connect:
	cd Docker/Debezium/mysql && ./reg_mysql_con.sh && echo "\n[+]Kafka connected to mysql\n" && cd ../es && ./reg_es_con.sh && echo "\n[+]Kafka connected to elasticsearch\n"


dev: 
	docker-compose up -d adminer kibana

rm:
ifdef svc
	docker stop $(svc) && docker rm $(svc)
else
	docker-compose down
endif

nuke:
	docker rmi $(docker image ls -q)

nuke_docker:
	@echo "Rebuilding docker services from scratch..."



doc:
	docker exec $(API_GO) bash -c "swag init -g main.go"

docu: doc restart
	@echo "Restart api"
