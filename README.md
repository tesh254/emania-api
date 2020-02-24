# emania-api
Esports platform for news, game trackers and so much more

## Development

### Prerequisites

* Docker installed
* Golang installed
* Docker Compose installed

### Setup

* Pull mongo image `docker pull mongo:latest`
* Run `docker-compose up -d` to run container on the background
* Open container in bash shell `docker exec -it packit bash`
    * To check contents of its database
        * Enter the mongo shell `mongo -u admin -p admin123 --authenticationDatabase packitdatabase`
        * To show databases inside the container `show dbs`
        * To access main db `use packit`
        * To check all collections `show tables`, `show collections` and `db.getCollectionNames()`
* Create 
