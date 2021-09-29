package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"math/rand"
	"net/http"
	"time"
)

type Person struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Age  float64 `json:"age"`
}

var persons = make([]Person, 0)

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "person",
	Description: "person definitions",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"name": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"age": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Query",
	Description:  "query type",
	Fields:      graphql.Fields{
		"person": &graphql.Field{
			Type: personType,
			Description: "get person by id",
			Args:              graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:         graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for i := range persons {
						if persons[i].ID == id {
							return persons[i], nil
						}
					}
				}
				return nil, nil
			},
		},
		"allPersons": &graphql.Field{
			Description: "get all persons",
			Type: graphql.NewList(personType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return persons, nil
			},
		},
	},
})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mutation",
	Fields:      graphql.Fields{
		"createPerson": &graphql.Field{
			Description:              "create new person",
			Type:              personType,
			Args:              graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"age": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, nameOk := p.Args["name"].(string)
				age, ageOk := p.Args["age"].(float64)
				rand.Seed(time.Now().Unix())
				person := Person{
					ID:   rand.Intn(10000),
				}
				if nameOk {
					person.Name = name
				}
				if ageOk {
					person.Age = age
				}
				persons = append(persons, person)
				return person, nil
			},
		},
		"updatePerson": &graphql.Field{
			Description:       "update person by id",
			Type:              personType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type:         graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"age": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				name, nameOk := p.Args["name"].(string)
				age, ageOk := p.Args["age"].(float64)
				person := Person{}
				for i := range persons {
					if persons[i].ID == id {
						if nameOk {
							persons[i].Name = name
						}
						if ageOk {
							persons[i].Age = age
						}
					}
					person = persons[i]
					break
				}
				return person, nil
			},
		},
		"deletePerson": &graphql.Field{
			Description:       "delete person by id",
			Type:              personType,
			Args:              graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				person := Person{}
				if ok {
					for i := range persons {
						if persons[i].ID == id {
							person = persons[i]
							persons = append(persons[:i], persons[i+1:]...)
						}
						break
					}
				}
				return person, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:         queryType,
	Mutation:     mutationType,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/person", func(res http.ResponseWriter, req *http.Request) {
		result := executeQuery(req.URL.Query().Get("query"), schema)
		json.NewEncoder(res).Encode(result)
	})
	fmt.Println("server is running on port: 8080")
	http.ListenAndServe(":8080", nil)
}

// Create person:: http://localhost:8080/person?query=mutation+_{createPerson(name:"Sagor",age:26){id,name,age}}

// Get single person by id:: http://localhost:8080/person?query={person(id:2388){id,name,age}}
// Get all persons:: http://localhost:8080/person?query={allPersons{id,name,age}}

// Update person info:: http://localhost:8080/person?query=mutation+_{updatePerson(id:2388,name:"sayf azad"){id,name,age}}

// Delete person:: http://localhost:8080/person?query=mutation+_{deletePerson(id:2388){id,name,age}}
