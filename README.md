# CHAT.CONNOR.FU-
# N

File structure:
```
chat-connor-fun
    -config
    -controllers
    -db
    -frontend
        -css
        -js
        index.html
    -models
    main.go
```

The `frontend` folder will later be packaged into the executable

`index.html` will probably just be a blank html document with the needed 
script tags  

`main.go` is the entry-point into the server

## Details
`chat.connor.fun` will have a REST api located at `chat.connor.fun/api/v1/`. 
I'll figure out what it will do later. I have no idea how to write Go...

## Building

### Requirements
You need Go SDK 1.9 and a recent version of PostgreSQL

### Package Path

The repository must be placed at `$GOPATH/src/github.com/aaronaaeng/chat.connor.fun`

### Go Dependencies

Just do 'dep ensure'

### Running 

In order to run the server, Postgres must be running at `postgresql://localhost:5432`
(this should be the default). 

then run `go run` (or to build `go build`)

The Postgres connection string can be overridden by setting
the `DATABASE_URL` environment variable. By default, the server will run in debug mode. This
can be overridden by setting the `CHAT_CONNOR_FUN_PROD`. The HMAC key used to sign JWTs
can be overridden using the `SECRET_KEY` environment variable (this should always be done
in production). If email verification is enabled, `MAILJET_PUBLIC_KEY` and `MAILJET_PRIVATE_KEY` must be set as well.
