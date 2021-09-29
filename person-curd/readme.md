# Go GraphQL CRUD example

Implement create, read, update and delete on Go.

To run the program:

1. go to the directory: `cd person-crud`
2. Run the example: `go run main.go`

## Create

`http://localhost:8080/person?query=mutation+_{createPerson(name:"Sagor",age:26){id,name,age}}`

## Read

* Get person by id: `http://localhost:8080/person?query={person(id:<id>){id,name,age}}`
* Get person list: `http://localhost:8080/person?query={allPersons{id,name,age}}`

## Update

`http://localhost:8080/person?query=mutation+_{updatePerson(id:<id>,name:"sayf azad"){id,name,age}}`

## Delete

<<<<<<< HEAD
` http://localhost:8080/person?query=mutation+_{deletePerson(id:2388){id,name,age}}`
=======
` http://localhost:8080/person?query=mutation+_{deletePerson(id:<id>){id,name,age}}`
>>>>>>> add all
