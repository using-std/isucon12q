SHELL=/bin/bash

BRANCH:=master

DATE=$(shell date +%Y%m_%d_%H%M)

.PHONY: noop
noop:
	echo "use arguments"

all: rotate deploy

# make checkout BRANCH=feature/checkout
# make deploy BRANCH=feature/checkout
# など
.PHONY: checkout
checkout:
	git fetch && \
	git reset --hard origin/$(BRANCH)

deploy: checkout deploy-nginx deploy-mysql deploy-go

deploy-nginx:
	sudo cp ./etc/nginx/sites-enabled/isuports.conf /etc/nginx/sites-enabled/isuports.conf
	sudo cp ./etc/nginx/nginx.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

deploy-mysql:
	sudo cp ./etc/mysql/conf.d/my.cnf /etc/mysql/conf.d/my.cnf
	sudo touch /var/log/mysql/mysql-slow.log
	sudo chmod 777 /var/log/mysql/mysql-slow.log
	sudo systemctl restart mysql

deploy-go:
	sudo systemctl restart isuports.service 

rotate: nginx-rotate mysql-rotate sqlite-rotate

analyze: analyze-alp analyze-sql

analyze-alp:
	sudo cat /var/log/nginx/access.log | alp ltsv -m "/api/player/player/*,/api/player/competition/.*/ranking,/api/organizer/competition/.*/score,/api/organizer/competition/.*/finish,/api/organizer/player/.*/disqualified" --sort sum -r | tee logs/nginx/alp.log

analyze-sql:
	sudo pt-query-digest /var/log/mysql/mysql-slow.log  | tee logs/mysql/digest.log

nginx-rotate:
	mkdir -p logs/nginx/backup
	sudo mv /var/log/nginx/access.log logs/nginx/backup/access.log.$(DATE) | :
	mv logs/nginx/alp.log logs/nginx/backup/alp.log.$(DATE) | :
	sudo touch /var/log/nginx/access.log
	sudo chmod 777 /var/log/nginx/access.log
	sudo systemctl restart nginx

mysql-rotate:
	mkdir -p logs/mysql/backup
	sudo mv /var/log/mysql/mysql-slow.log logs/mysql/backup/mysql-slow.log.$(DATE) | :
	mv logs/mysql/digest.log logs/mysql/backup/digest.log.$(DATE) | :
	sudo touch /var/log/mysql/mysql-slow.log
	sudo chmod 777 /var/log//mysql/mysql-slow.log
	sudo systemctl restart mysql

sqlite-rotate:
	mkdir -p logs/sqlite/backup
	mv logs/sqlite/trace.log logs/sqlite/backup/trace.log.$(DATE) | :
	touch logs/sqlite/trace.log
