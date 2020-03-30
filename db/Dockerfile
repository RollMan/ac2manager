FROM mariadb

RUN apt update \
 && apt install -y apache2-utils

EXPOSE 3306

COPY init_sql.sh /docker-entrypoint-initdb.d

CMD ["mysqld"]
