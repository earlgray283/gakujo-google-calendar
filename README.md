# gakujo-google-calendar

## For Developer

### Setup

```console
$ echo 'GAKUJO_USERNAME={学情のユーザーネーム}' >> .env
$ echo 'GAKUJO_PASSWORD={学情のパスワード}' >> .env
```

### Requirements

- Go(https://go.dev/doc/install)

### Rules

- branch name: `{username}-#{issue}` (e.g. `earlgray-#1`)
- coding style: obey [Effective Go](https://go.dev/doc/effective_go)
    - naming: the names of functions, valiables, etc. should be camelCase or PascalCase.
    - format: use `$ gofmt .`(You can format files automaticaly by using VScode extension.)
- You should pass CI when you send pull request.
- You should commit frequently.
- The commit messages should be easy to understand.
