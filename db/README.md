# データベースの設定
テーブルスキーマを `schema.sql` に記述します。

`PRIMARY KEY` や `UNIQUE KEY` は、アプリケーションの特性に応じて考えましょう。
```sql
CREATE TABLE user (
    `id` BIGINT UNSIGNED NOT NULL,

    `name` VARBINARY(32) NOT NULL,
    `password_hash` VARBINARY(254) NOT NULL,

    `created_at` DATETIME(6) NOT NULL,
    `updated_at` DATETIME(6) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY (name),

    KEY (created_at),
    KEY (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

`docker/init.sh` はdockerイメージの初期化時に自動で呼ばれます。
