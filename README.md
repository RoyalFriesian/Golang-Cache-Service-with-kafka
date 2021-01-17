# Golang-Cache-Service-with-kafka
Do not take this code seriously! It is just a 4 days work. There is huge scope of improvement.  


# Steps to run the setup
  Suppose your work directory is 
    ```
    C:\Go
    ```
### 1. Setup Mongo
  * First install mongo
  * Create Users
  * Start Mongod in one terminal from c:\Go 
    ```
    mongod --port 9999 --auth --dbpath cache
    ```
  * Start Mongo in another terminal create two users
    * Login with admin 
    ```
    mongo --port 9999  --authenticationDatabase "admin" -u "myUserAdmin" -p"{ADMINPASS}"
    ```
    * Create One User (myTester) for cachedb and disconnect mongo
    ```
        use cachedb
        db.createUser(
          {
            user: "myTester",
            pwd:  passwordPrompt(),   // or cleartext password
            roles: [ { role: "readWrite", db: "cachedb" }]
          }
        )
    ```
    * Again Login with admin 
    ```
    mongo --port 9999  --authenticationDatabase "admin" -u "myUserAdmin" -p"{ADMINPASS}"
    ```
    * Create another User (myTester2) for backupdb and disconnect mongo
    ```
        use cachedb
        db.createUser(
          {
            user: "myTester2",
            pwd:  passwordPrompt(),   // or cleartext password
            roles: [ { role: "readWrite", db: "backupdb" }]
          }
        )
    ```
    
    
 ### 2. Setup Kafka and Zookeeper
 * install Docker and run the compose file inside kafka-installation
```
version: '3'

services:
    zookeeper:
      image: wurstmeister/zookeeper
      container_name: zookeeper
      ports:
        - "2181:2181"

    kafka:
      image: wurstmeister/kafka
      container_name: kafka
      ports:
        - "9092:9092"
      environment:
        KAFKA_ADVERTISED_HOST_NAME: localhost
        KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181 
```

### 3. Setup Environment Variables in your machine
```
        KAFKA_ADVERTISED_HOST_NAME: localhost
        KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181 
```


### 4. Update Golang file with mongo and kafka credentials 

```
const (
	hosts      = "localhost:9999"
	database   = "cachedb"
	username   = "myTester"
	password   = "mytesterpass"
	collection = "data"

	hosts_bk      = "localhost:9999"
	database_bk   = "backupdb"
	username_bk   = "myTester2"
	password_bk   = "mytester2pass"
	collection_bk = "data"
	
	kafkaBroker	  = "localhost:9092"
)
```

### 5. Run the Kafka and Zookeeper in docker 
### 6. Run golang program main.go

### 7. Open postman and push some entries in database 
  * post variables are key and data
```
http://localhost:8080/insert
```
### 8. Read data from database with pagination
  * read url with pagination
```
http://localhost:8080/data/page/?page=1
```

### 8. Reload data from backup db to cachedb from kafka data pipline 
  * Open kafka producer 
```
C:\WINDOWS\system32>docker exec -it kafka /bin/sh
C:\WINDOWS\system32>cd/ opt/kafka_2.13-2.7.0/bin
/opt/kafka_2.13-2.7.0/bin # kafka-console-producer.sh  --broker-list kafka:9092 --topic myTopic
```
  * Now send message as producer. This will reload the cachedb from backup db
```
> reload
```

