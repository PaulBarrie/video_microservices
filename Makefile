include .env
export $(shell sed 's/=.*//' .env)

restart_dock: 
	[ !-z "${1}"] && docker-compose restart || docker restart $1

rebuild:
	 docker-compose up --build --force-recreate --no-deps $1
	
nice_message:
	echo "Bravo Polo, tu gères ! Ton api est edispo sur le port 3000"

up:
	docker-compose up $1

down: 
	docker-compose down

nuke:
	docker rmi $(docker image ls -q)

start_search:
	docker-compose up msql debezium kafka elasticsearch kibana
	
nuke_docker:
	@echo "Rebuilding docker services from scratch..."

doc:
	docker exec $(API_GO) bash -c "swag init -g main.go"

docu: doc restart
	@echo "Restart api"
