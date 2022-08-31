# Contact Manager
A simple contact manager used as a sample for the 2022 $COMPANY Summer Campus. 
Every subfolder of this project contains the project at a different 
step of the development process.

## Requirements
You need to install the golang package from the website 
https://go.dev/doc/install

## Run the project
To run the project in development mode please type
```bash
go run main.go
```

## Build the project
The build is an executable file under windows.
```bash
go build
```

## Development process
Those are the steps that the we have to perform in order to make our application 
working. 

### Setup the environment
We can easy achieve that by building a simple hello world and make sure that it 
will work.


### Setup the dependencies. 
We have decided to use those dependencies.

1. gin-gonic as our web framework. A web framework hides the complexity of 
having to manage the sockets, the path's addressing etc.. 

2. gorm with sqlite as a database. Gorm allows us to abstract from SQL 
queries and to think more about objects when writing/reading database. Row
are just the way that we are using to store the informations.

Follow the guide at https://gorm.io/docs/ in order to install the package.

We have also to create a PostgresSQL database in order to persist our data 
https://www.postgresql.org/download/. I suggest to download the zip archive 
instead of the installer. Doing this you can run the database using the command 
line in a folder of your choice. I've prepared a data folder that can do the 
job for us.

```sh
initdb.exe -D data --encoding UTF8

postgres -D data
```

And then run the queries with the `psql` command.

```sh
psql -d postgres -c "CREATE USER gorm WITH PASSWORD 'gorm'"
psql -d postgres -c "CREATE DATABASE contacts WITH OWNER 'gorm'"
```
In order to connect from golang we have to install the driver. Gorm supports
different types of database, we may choose the correct one.

```sh
go get gorm.io/driver/postgres
```
### Create the base objects application. 
In this example is only one (`Contact`). We also try to store and read the value
from the database. We choose to transmit our data using JSON so we have to test 
also that the object get correctly marshalled. 
Operations that I have to perform on the object are:
- _create_ creating an object in the database and generate an `id`.
- _list_ returning all the object that are into the table.
- _read_ reading the values of an object stored into the database using his 
  own `id.`
- _update_ updating the values of an object that are stored into the 
database using object own `id`.
- _delete_ remove a row from the database using object own `id`.

### Bring the web framework in. 
Create the REST compliant endpoints for reading an writing the resource. The url
endpoints are described using an openapi specification. _.. set here url of the
swagger .._

The application can be easily deployed on a server. I did it for you at the 
endpoint address _somewhere on the net_.

## Appendix
### Swagger generation
Run this command into the `05-release` folder.
```
swag init --output ./docs/swagger
```