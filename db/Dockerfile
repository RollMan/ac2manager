FROM mariadb

EXPOSE 3306

COPY init.sql /docker-entrypoint-initdb.d

CMD ["mysqld"]
