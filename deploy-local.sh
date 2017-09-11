deployTimeCapsule () {
	cd /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/tc-server
	echo >&2 "building time capsule executable..."
	GOOS=linux go build

	echo >&2 "building time capsule docker container..."
	docker build -t brendankellogg/iuga-timecapsule .

	# docker pull mysql

	if [ "$(docker ps -aq --filter name=iuga-timecapsule)" ]; then
		docker rm -f iuga-timecapsule
	fi

    startMySQL

	echo >&2 "starting iuga-timecapsule..."
	docker run -d \
	--name iuga-timecapsule \
    --link mysql:mysql \
	-p 4000:4000 \
	-v /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/tc-server/tls:/tls:ro \
    -v /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/tc-server/secret:/secret:ro \
	-e TLSKEY=$TLSKEY \
	-e TLSCERT=$TLSCERT \
    -e PORT=4000 \
    -e MYSQLADDR="mysql:3306" \
    -e MYSQLPASS=$MYSQLPASS \
    -e MYSQLDB="CapsuleDB" \
    -e GOOGLE_REFRESH_TOKEN=$GOOGLE_REFRESH_TOKEN \
	brendankellogg/iuga-timecapsule:latest

}

deployTimeCapsuleNoDocker() {
	cd /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/tc-server
	echo >&2 "building time capsule executable..."
	go install
	tc-server
}

startMySQL() {
    if [ "$(docker ps -aq --filter name=mysql)" ]; then
        docker stop mysql
        docker rm mysql
    fi

    docker run -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=$MYSQLPASS \
    -e MYSQL_DATABASE=CapsuleDB \
    -v /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/mysql:/var/lib/mysql \
    --name mysql \
    mysql
}

if [[ "$1" == "tc" ]]; then
    deployTimeCapsule
elif [[ "$1" == "tc-nd" ]]; then
	deployTimeCapsuleNoDocker
fi