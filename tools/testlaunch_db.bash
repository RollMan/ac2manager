#!/bin/bash
docker build . -t mysql_ac2 && docker run -it -e MYSQL_ROOT_PASSWORD=mysqlrootpw -e DEFAULT_USER_ID=admin -e DEFAULT_USER_PW=pw -p3306:3306  mysql_ac2
