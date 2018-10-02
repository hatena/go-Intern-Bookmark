# loader
このパッケージでは、GraphQLのエンドポイントためのデータローダーを実装します。
[graph-gophers/dataloader](https://godoc.org/github.com/graph-gophers/dataloader) パッケージを使います。

GraphQLのリクエストからデータを引くまでの流れは次のようになります。

- リクエストの `Context` に、各モデルの `dataloader.Loader` を保存する (`Attach`)
- クエリに対してモデルに対する `Resolver` を返す
- `Context` から引きたいモデルに対応する `dataloader.Loader` を探す (`getLoader`)
- `dataloader.Loader` 経由でデータを引くと、 `dataloader.BatchFunc` が呼ばれてクエリをまとめることができる
