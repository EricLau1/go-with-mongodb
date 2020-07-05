# Go with MongoDB + Docker

### Dependencies

```bash
    go get -u github.com/joho/godotenv
    
    go get -u gopkg.in/mgo.v2

    go get -u github.com/emicklei/go-restful

```

> Install Docker!

---

### Config Files


__init-mongo.js__

```js
db.createUser({
    user: 'root',
    pwd: 'root',
    roles: [
        {
            'role': 'readWrite',
            'db': 'mydb'
        }
    ]
});
```


__docker-compose.yml__

```yml
version: '3'
services:
    database:
        image: 'mongo'
        container_name: 'mongodb_container'
        environment: 
            - MONGO_INITDB_DATABASE=mydb
            - MONGO_INITDB_USERNAME=root
            - MONGO_INITDB_PASSWORD=root
        volumes:
            - ./init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro
            - ./data:/data/db
        ports:
            - '27017-27019:27017-27019'
```

### Run MongoDb with Docker

```bash
    cd docker

    docker-compose up -d
```








