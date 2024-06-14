# ramchi
`ramchi` is an extension to `chi` for rapid &amp; modular development of sites.
It allows the user to serve both backend and frontend data inside of the same go project,
without the need for coding any javascript (`ramchi` does all this work for you).
`ramchi` is based upon developer experience and usage, while still making your website fast and responsive.

## Install

`go get -u github.com/etwodev/ramchi`

## Config

When you create [your first server](), `ramchi` will generate a `ramchi.config.json` file,
which allows you to configure aspects of the server.

```json
{
 "port": "8080",
 "address": "localhost",
 "experimental": false
}
```

