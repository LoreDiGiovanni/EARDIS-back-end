
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
