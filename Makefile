include .env
export $(shell sed 's/=.*//' .env)

restart: 
ifdef svc
	make rm svc="$(svc)" && make up svc="$(svc)"
else
	docker-compose restart 
endif
.PHONY: restart

rebuild:
	 docker-compose up --build --force-recreate --no-deps $1
.PHONY: rebuild

ps:
	docker-compose ps
.PHONY: ps

up:
ifdef svc
	docker-compose up -d --build $(svc)
else
	docker-compose up -d --build
endif
.PHONY: up

rm:
ifdef svc
	docker stop $(svc) && docker rm $(svc)
else
	docker-compose down
endif
.PHONY: rm

nuke:
	docker rmi $(docker image ls -q)
.PHONY: nuke

