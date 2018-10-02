FROM node:8 as node

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install

COPY ui ui

RUN yarn build


FROM golang:1.10

ENV WORKDIR /go/src/github.com/hatena/go-Intern-Bookmark
WORKDIR $WORKDIR

COPY Makefile Gopkg.toml Gopkg.lock ./
RUN make setup

COPY . $WORKDIR

COPY --from=node /app/static $WORKDIR/static

CMD ["./script/localup"]
