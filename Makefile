SHELL=/bin/bash


DATE=$(shell date +%Y%m_%d_%H%M)


all: deploy-nginx

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

nginx-rotate:
	mkdir -p logs/nginx/backup
	mv logs/nginx/access.log logs/nginx/backup/access.log.$(DATE) | :
	mv logs/nginx/alp.log logs/nginx/backup/alp.log.$(DATE) | :
	sudo touch logs/nginx/access.log
	sudo chmod 777 logs/nginx/access.log

mysql-rotate:
	mkdir -p logs/mysql/backup
	mv logs/mysql/mysql-slow.log logs/mysql/backup/mysql-slow.log.$(DATE) | :
	mv logs/nginx/digest.log logs/nginx/backup/digest.log.$(DATE) | :
	sudo touch logs/mysql/mysql-slow.log
	sudo chmod 777 logs/mysql/mysql-slow.log
