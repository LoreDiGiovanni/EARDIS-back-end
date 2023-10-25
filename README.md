
## Mongodb docker 
mongodb docker sample
 ``` yaml
services:
  mongo:
    image: mongo
    container_name: mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: example
    ports:
      - 27017:27017
 ```

## Export

 ``` 
export MONGODB_URI="mongodb://root:example@127.0.0.1:27017/"
export JWT_SECRET="testpass"
 ```
