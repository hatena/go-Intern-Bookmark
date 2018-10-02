# resolver
このパッケージでは、GraphQLスキーマと、GraphQLのクエリを受けるHTTP handlerを実装します。

`schema.graphql` にはGraphQLスキーマを書きます。
```graphql
schema {
    query: Query
    mutation: Mutation
}

type Query {
    getBookmark(bookmarkId: ID!): Bookmark!
}

type Mutation {
    createBookmark(url: String!, comment: String!): Bookmark!
    deleteBookmark(bookmarkId: ID!): Boolean!
}

type Bookmark {
    // ...
}
```

`handler.go` では、このスキーマをloadし (実際にはgo-assets-builderによって実行ファイルに埋め込まれる)、 `http.Handler` を返す実装を行います。
```go
func NewHandler(app service.BookmarkApp) http.Handler {
	graphqlSchema, err := loadGraphQLSchema()
	if err != nil {
		panic(err)
	}
	schema := graphql.MustParseSchema(string(graphqlSchema), newResolver(app))
	return &relay.Handler{Schema: schema}
}
```

`Resolver` は、GraphQLクエリの引数を受け取り、データを引いてきて返します。
```go
// getBookmark(bookmarkId: ID!): Bookmark! に対応する

func (r *resolver) GetBookmark(ctx context.Context, args struct{ BookmarkID string }) (*bookmarkResolver, error) {
	// ...
	return &bookmarkResolver{bookmark: bookmark}, nil
}
```

各モデルのプロパティーをたどってデータを取得できるように、スキーマに対応する名前の公開関数を実装します。
```go
// type Bookmark {
//     id: ID!
//     comment: String!
//     user: User!
//     entry: Entry!
// }

func (b *bookmarkResolver) ID(ctx context.Context) graphql.ID {
	return graphql.ID(fmt.Sprint(b.bookmark.ID))
}

func (b *bookmarkResolver) Comment(ctx context.Context) string {
	return b.bookmark.Comment
}

func (b *bookmarkResolver) User(ctx context.Context) (*userResolver, error) {
	// ...
	return &userResolver{user: user}, nil
}

func (b *bookmarkResolver) Entry(ctx context.Context) (*entryResolver, error) {
	// ...
	return &entryResolver{entry: entry}, nil
}
```

`User` と `Entry` を返すには、更にデータベースを引かなければいけません。
`resolver` が `app` を持っているので、 `bookmarkResolver` にも `app` を持たせて直接引いても構いません。
なぜ `Context` を引数にとり、 `loader` パッケージを介してデータを取得しているのでしょうか。
