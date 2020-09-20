#!/bin/bash
set -eux

### Address of destination
IP_ADDRESS=

### Fetch and place configuration JSONs
INSTANCE_ID=$(wget -q -O - http://169.254.169.254/latest/meta-data/instance-id)
rsync -r -e "ssh -i ~/.ssh/id_rsa" "admin@${IP_ADDRESS}:/opt/ac2manager/${INSTANCE_ID}/*" ~/server/cfg

### Launch
exec wine /home/ubuntu/server/accServer.exe
