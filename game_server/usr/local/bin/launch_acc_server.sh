#!/bin/bash

### Fetch and place configuration JSONs
APP_ADDRESS=
rsync -r -e "ssh -i ~/.ssh/id_rsa" admin@${APP_ADDRESS}:/opt/${APP_ADDRESS}/

### Launch
wine /home/admin/server/accServer.exe
