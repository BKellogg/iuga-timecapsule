startMySQL() {
    if [ "$(docker ps -aq --filter name=mysql)" ]; then
		docker rm -f mysql
	fi

    docker run -d \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=$MYSQLPASS \
    -e MYSQL_DATABASE=CapsuleDB \
    -v /Users/Brendan/Documents/UW/INFO344/go/src/github.com/BKellogg/iuga-timecapsule/mysql:/var/lib/mysql \
    --name mysql \
    mysql
}

startMySQL