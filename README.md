# go-Intern-Bookmark
## ディレクトリー構成
- db: データベースのテーブルスキーマ
- config: 環境変数から読み込む設定
- model: モデル層: 型定義を書きます
- repository: データベースにアクセスするリポジトリ層
- service: Bookmarkアプリケーションを定義するサービス層
- web: webサーバーのルーティングやリクエストの解釈、レスポンスを実装するweb層
- ui: フロントエンド
- templates: HTMLテンプレート
- resolver: GraphQLスキーマとクエリの実行
- loader: GraphQL用のデータローダー

## 実行
サーバー起動
```sh
docker-compose up --build
open http://localhost:8000
```

GraphiQL
```
open http://localhost:8000/graphiql
```

テスト実行
```sh
docker-compose build && docker-compose run --rm app make test
```
または、手元のGoでテストする (早い)
```sh
DATABASE_DSN_TEST=root@(localhost:3306)/intern_bookmark_test make test
```

MySQL
```
docker-compose exec db mysql
mysql> use intern_bookmark;
mysql> select * from user;
```

Dockerコンテナを作り直すために一旦削除
```
docker ps -a
docker stop <container id>
docker rm <container id>
# または
docker rm -f <container id>
```

例: データベーススキーマを変更する時はmysqlコンテナを作り直す
```
docker ps -a | grep mysql
docker rm -f 83ece5d09062
```

## フロントエンド
サーバーを起動すると自動的にwatchされる
```
docker-compose up
```

テスト実行
```
docker-compose run --rm node yarn test
```
