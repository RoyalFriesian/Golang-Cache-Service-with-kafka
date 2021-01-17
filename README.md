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
    
    
 ### 2. Setup Kafka 
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

