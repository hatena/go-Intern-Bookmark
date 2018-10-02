#!/bin/bash
set -xe

mysqladmin -uroot create intern_bookmark
mysqladmin -uroot create intern_bookmark_test

mysql -uroot intern_bookmark < /app/db/schema.sql
mysql -uroot intern_bookmark_test < /app/db/schema.sql
