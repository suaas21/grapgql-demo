package graphql_objects

import (
	"fmt"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/suaas21/grapgql-demo/books-authors-query/model"
	"github.com/suaas21/grapgql-demo/books-authors-query/storage"
	"github.com/suaas21/grapgql-demo/books-authors-query/utils"
	"strings"
)

func init() {
	BookType.AddFieldConfig("authors", &graphql.Field{
		Type: graphql.NewList(AuthorType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var (
				book, bookOk = p.Source.(model.Book)
				loaders      = p.Context.Value("loaders").(map[string]*dataloader.Loader)
				handleErrors = func(errors []error) error {
					if len(errors) == 0 {
						return nil
					}
					var errs []string
					for _, e := range errors {
						errs = append(errs, e.Error())
					}
					return fmt.Errorf(strings.Join(errs, "\n"))
				}
			)
			if !bookOk {
				return nil, nil
			}
			var keys dataloader.Keys
			for i := range book.AuthorIDs {
				keys = append(keys, utils.NewResolverKey(book.AuthorIDs[i]))
			}

			thunk := loaders[utils.BookAuthorIds].LoadMany(p.Context, keys)
			return func() (interface{}, error) {
				res, errors := thunk()
				return res, handleErrors(errors)
			}, nil
		},
	})
	AuthorType.AddFieldConfig("books", &graphql.Field{
		Type: graphql.NewList(BookType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var (
				author, authorOk = p.Source.(model.Author)
				loaders          = p.Context.Value("loaders").(map[string]*dataloader.Loader)
				handleErrors     = func(errors []error) error {
					if len(errors) == 0 {
						return nil
					}
					var errs []string
					for _, e := range errors {
						errs = append(errs, e.Error())
					}
					return fmt.Errorf(strings.Join(errs, "\n"))
				}
			)
			if !authorOk {
				return nil, nil
			}
			var keys dataloader.Keys
			for i := range author.BookIDs {
				keys = append(keys, utils.NewResolverKey(author.BookIDs[i]))
			}

			thunk := loaders[utils.AuthorBookIds].LoadMany(p.Context, keys)
			return func() (interface{}, error) {
				res, errors := thunk()
				return res, handleErrors(errors)
			}, nil
		},
	})
}

var BookType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Book",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		// "authors" type & resolver define init() func
	},
})

var AuthorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		// "books" type & resolver define init() func
	},
})

var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"book": &graphql.Field{
			Type:        BookType,
			Description: "get book by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, book := range storage.ListBook {
						if book.ID == uint(id) {
							return book, nil
						}
					}
				}
				return nil, nil
			},
		},
		"author": &graphql.Field{
			Type:        AuthorType,
			Description: "get author by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, author := range storage.ListAuthor {
						if author.ID == uint(id) {
							return author, nil
						}
					}
				}
				return nil, nil
			},
		},
	},
})

var MutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"book": &graphql.Field{
			Description: "create new book",
			Type:        BookType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"author_ids": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(graphql.Int)),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, idOk := p.Args["id"].(int)
				name, nameOk := p.Args["name"].(string)
				description, desOk := p.Args["description"].(string)
				authorIds, authorIdsOk := p.Args["author_ids"].([]interface{})
				book := model.Book{}
				if idOk {
					book.ID = uint(id)
				}
				if nameOk {
					book.Name = name
				}
				if desOk {
					book.Description = description
				}
				if authorIdsOk {
					var auids = make([]uint, 0)
					for _, aid := range authorIds {
						auids = append(auids, uint(aid.(int)))
					}
					book.AuthorIDs = auids
				}
				storage.ListBook = append(storage.ListBook, book)
				fmt.Println(book)
				return book, nil
			},
		},
		"author": &graphql.Field{
			Description: "create new author",
			Type:        BookType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"book_ids": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, idOk := p.Args["id"].(int)
				name, nameOk := p.Args["name"].(string)
				bookIds, bookIdsOk := p.Args["book_ids"].([]interface{})
				author := model.Author{}
				if idOk {
					author.ID = uint(id)
				}
				if nameOk {
					author.Name = name
				}
				if bookIdsOk {
					var auids = make([]uint, 0)
					for _, aid := range bookIds {
						auids = append(auids, uint(aid.(int)))
					}
					author.BookIDs = auids
				}
				storage.ListAuthor = append(storage.ListAuthor, author)
				return author, nil
			},
		},
	},
})
