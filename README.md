# simple-webapp


## Build and Deploy

After creating an account of heroku, execute commands below.

```shell
$ heroku container:login
$ heroku create [AppName]
$ heroku addons:create heroku-postgresql:hobby-dev -a [AppName]
$ heroku container:push web -a [AppName]
$ heroku open -a [AppName]
```

## Behavior

Behavior table according to http methods.

This characters, ":id", are replaced by integer.

|      | GET | POST | PUT | DELETE |
|:-----|:----:|:----:|:----:|:----:|
|/     | send "Hello World!!" | | | |
|/users| send user list | add user | | |
|/users/:id| show selected user | | modify selected user | delete selected user |

## Data format

This app receives data and returns response in json format.

### Send data format
```json
{
  "name":"test",
  "email":"hoge@example.com"
}
```

### Receive data format

Timezone is "Asia/Japan".
```json
{
  "id": 1,
  "name": "test",
  "email": "hoge@example.com",
  "created_at": "2017-09-08T22:26:32.205803424+09:00",
  "updated_at": "2017-09-08T22:26:32.205803424+09:00"
}
```

## Example

```shell
$ curl -XPOST -H 'Content-Type:application/json' https://(AppName).herokuapp.com/users -d '{"name": "test", "email": "hoge@example.com" }'
```

```json
{
  "id": 1,
  "name": "test",
  "email": "hoge@example.com",
  "created_at": "2017-09-08T22:26:32.205803424+09:00",
  "updated_at": "2017-09-08T22:26:32.205803424+09:00"
}
```

