SHELL=/bin/bash


DATE=$(shell date +%Y%m_%d_%H%M)

.PHONY: noop
noop:
	echo "use arguments"

all: rotate deploy

deploy: deploy-nginx deploy-mysql

deploy-nginx:
	sudo cp ./etc/nginx/sites-enabled/isuports.conf /etc/nginx/sites-enabled/isuports.conf
	sudo cp ./etc/nginx/nginx.conf /etc/nginx/nginx.conf
	sudo systemctl restart nginx

deploy-mysql:
	sudo cp ./etc/mysql/conf.d/my.cnf /etc/mysql/conf.d/my.cnf
	sudo touch /var/log/mysql/mysql-slow.log
	sudo chmod 777 /var/log/mysql/mysql-slow.log
	sudo systemctl restart mysql

rotate: nginx-rotate mysql-rotate

analyze: analyze-alp analyze-sql

analyze-alp:
	sudo cat /var/log/nginx/access.log | alp ltsv -m '/api/estate/req_doc/.\d+,/api/estate/.\d+,/api/chair/.\d+,/api/recommended_estate/.\d+,/api/chair/buy/.\d+' --sort=sum -r | tee logs/nginx/alp.txt

analyze-sql:
	sudo pt-query-digest /var/log/mysql/mysql-slow.log  | tee logs/mysql/digest.log

nginx-rotate:
	mkdir -p logs/nginx/backup
	mv /etc/nginx/access.log /etc/nginx/backup/access.log.$(DATE) | :
	mv logs/nginx/alp.log logs/nginx/backup/alp.log.$(DATE) | :
	sudo touch /var/log/nginx/access.log
	sudo chmod 777 /var/log/nginx/access.log

mysql-rotate:
	mkdir -p logs/mysql/backup
	mv logs/mysql/mysql-slow.log logs/mysql/backup/mysql-slow.log.$(DATE) | :
	mv logs/nginx/digest.log logs/nginx/backup/digest.log.$(DATE) | :
	sudo touch /var/log/mysql/mysql-slow.log
	sudo chmod 777 /var/log//mysql/mysql-slow.log
