FROM mysql:5.7

VOLUME ["/var/lib/mysql"]
COPY database.sql /docker-entrypoint-initdb.d/database.sql


EXPOSE 3306