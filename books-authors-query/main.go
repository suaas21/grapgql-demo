package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	loader_fn "github.com/suaas21/grapgql-demo/books-authors-query/dataloader"
	graphql_objects "github.com/suaas21/grapgql-demo/books-authors-query/graphql-objects"
	"github.com/suaas21/grapgql-demo/books-authors-query/utils"
	"net/http"
)

func executeQuery(query string, schema graphql.Schema, ctx context.Context) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       ctx,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {

	http.HandleFunc("/graphql", func(res http.ResponseWriter, req *http.Request) {
		var loaders = make(map[string]*dataloader.Loader, 1)
		loaders[utils.BookAuthorIds] = dataloader.NewBatchedLoader(loader_fn.GetAuthorsBatchFn)
		loaders[utils.AuthorBookIds] = dataloader.NewBatchedLoader(loader_fn.GetBooksBatchFn)

		schema, err := graphql.NewSchema(graphql.SchemaConfig{
			Query:    graphql_objects.QueryType,
			Mutation: graphql_objects.MutationType,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx := context.WithValue(context.Background(), "loaders", loaders)
		result := executeQuery(req.URL.Query().Get("query"), schema, ctx)
		_ = json.NewEncoder(res).Encode(result)
	})
	fmt.Println("server is running on port: 8080")
	_ = http.ListenAndServe(":8080", nil)

}

// http://localhost:8080/graphql?query=mutation+_{book(id:1,name:"Sagor",description:"childhood",author_ids:[1]){id,name,description}}
// http://localhost:8080/graphql?query=mutation+_{author(id:1,name:"Sagors childhood",book_ids:[1]){id,name}}

// http://localhost:8080/graphql?query={book(id:1){id,name,description,authors{id,name}}}
// http://localhost:8080/graphql?query={author(id:1){id,name,books{id,name}}}
