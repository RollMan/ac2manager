#!/bin/bash
docker build . -t mysql_ac2 && docker run -it --env-file ../.env -p3306:3306  mysql_ac2
