## Installation
### Prerequisites
- Permission of AWS EC2 StartInstances, StopInstances, DescribeInstances and DescribeInstanceStatuses
- EC2 AMI which contains `/home/ubuntu/server/{accServer.exe, cfg/}`
- A computer for the web server (e.g. AWS Lightsail) with a user of name `admin`

### Installation
#### EC2 Instance Image
- Download source code
- Place the contents in `game_server/`
- Set `IP_ADDRESS` with Lightsail address on `launch_acc_server.sh`
- Create ssh key and upload it to `authorized_keys` in Lightsail
- `systemctl daemon-reload`
- `systemctl enable ac2manager.service`

#### Lightsail
##### WEB
TODO

##### EC2CTL
- Fill a variable of `id` in `jobmng/jobmng.go`
- Modify default configurations in `confjson/*.json`

#### Tokens
Fill the below:

```
AC2_DB_USERNAME=
MYSQL_ROOT_PASSWORD=
MYSQL_PASSWORD=
AC2_APP_ADMINUSERNAME=
AC2_APP_ADMINPASSWORD_HASHED=
JWT_SIGNING_KEY=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
```

and write as a file named `.env` at the root path of this project. `AC2_APP_ADMINPASSWORD_HASHED` is a password hashed by bcrypt like `\$2a\$nn\$xxxxxxxxxxxxxxxxxxxxxxxxxx`.


### Launch
`docker-compose up -d`
