ENV=dev

pre:
ifdef e
ENV=${e}
endif

set-env := export ENV=$(ENV) ;\
	export COMPOSE_PATH_SEPARATOR=: ;\
	export COMPOSE_FILE=docker-compose.yml:docker-compose.$(ENV).yml

up: pre
	$(set-env) ;\
	docker-compose up -d

down: pre
	$(set-env) ;\
	docker-compose down

rebuild: pre
	$(set-env) ;\
	docker-compose build --no-cache

clean: pre
	docker system prune -a
