package resolver

import (
	"io/ioutil"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/hatena/go-Intern-Bookmark/service"
)

//go:generate go-assets-builder --package=resolver --output=./schema-gen.go --strip-prefix="/" --variable=Schema schema.graphql

func loadGraphQLSchema() ([]byte, error) {
	file, err := Schema.Open("schema.graphql")
	if err != nil {
		return nil, err
	}
	schemaBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil
	}
	return schemaBytes, nil
}

func NewHandler(app service.BookmarkApp) http.Handler {
	graphqlSchema, err := loadGraphQLSchema()
	if err != nil {
		panic(err)
	}
	schema := graphql.MustParseSchema(string(graphqlSchema), newResolver(app))
	return &relay.Handler{Schema: schema}
}
