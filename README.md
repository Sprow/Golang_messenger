change variable CHOOSE_DATABASE=(MONGO_DB or POSTGRE_SQL) to switch for another DB



docker PostgreSQL

docker run --name postgreDB -p 5432:5432 -e POSTGRES_USER=sprow -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=messenger --rm postgres


docker MongoDB

docker run  -p 27017:27017  --rm --name mymongo mongo:latest