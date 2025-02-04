
Installing RabbitMQ and Erlang dependencies:
  * `sudo su`
  * `apt-get update`
  * `apt-get upgrade`
  * `apt-get install erlang`
  * `apt-get install rabbitmq-server`

Setting up RabbitMQ:
  * `systemctl enable rabbitmq-server`
  * `systemctl start rabbitmq-server`
  * `systemctl status rabbitmq-server`
  * `rabbitmq-plugins enable rabbitmq_management`

Creating admin user on (http://localhost:15672/#/):
* `rabbitmqctl add_user admin admin`
* `rabbitmqctl set_user_tags admin administrator`
* `rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"`

Useful RabbitMQ commands:
* `rabbitmq-plugins list`
* `rabbitmqctl status`
* `rabbitmqctl list_queues`
* `rabbitmqctl cluster_status`

Install PostgreSQL:
* `sudo apt-get update`
* `sudo apt-get install postgresql postgresql-contrib`
* `ls /etc/postgresql/<version>/main/`
* `service postgresql status`

Command line tool:
* `man psql`
* `sudo su postgres`
* `psql`
* `\l`  - List databases
* `\du` - List roles

Change default user password and create new user with proper rights:
* `ALTER USER postgres WITH PASSWORD 'admin';`
* `CREATE USER goadmin WITH PASSWORD 'goadmin';`
* `ALTER USER goadmin WITH SUPERUSER;`
* `DROP USER <user>;`

Install pgAdmin 3 from the Ubuntu Software Center and setup with:
* `Host: 127.0.0.1` and `Port: 5432`

Project dependencies:
* Go RabbitMQ Client Library (https://godoc.org/github.com/streadway/amqp): `go get github.com/streadway/amqp`
* Go PostgreSQL Connector (https://godoc.org/github.com/lib/pq): `go get github.com/lib/pq`
* Go web sockets connector (http://www.gorillatoolkit.org/pkg/): `go get github.com/gorilla/websocket`
