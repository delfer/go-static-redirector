# Go Static Redirerctor
[![Docker Stars](https://img.shields.io/docker/stars/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/) [![Docker Pulls](https://img.shields.io/docker/pulls/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/) [![Docker Automated build](https://img.shields.io/docker/automated/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/) [![Docker Build Status](https://img.shields.io/docker/build/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/) [![MicroBadger Layers](https://img.shields.io/microbadger/layers/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/) [![MicroBadger Size](https://img.shields.io/microbadger/image-size/delfer/go-static-redirector.svg)](https://hub.docker.com/r/delfer/go-static-redirector/)

Send redirect (301 or 302) HTTP response fast, flexible and (optional) log visits to ClickHouse 

## Features

- Response very fast
- Log every visit to Yandex ClickHouse
- Configurable by environment variables
- Shows buffer size on `/load`

## Configuration

Environment variables:
- `PORT` - HTTP listen port (8080 by default)
- `REDIRECTS` - list of redirects (see section "REDIRECTS")
- `PERMANENTLY` - set `true` to response `301 Moved Permanently` instead of `302 Found` (false=`302 Found` by default)
- `BUFFER` - buffer size (in requests) between HTTP server and DB writer (100,000 by default)
- `DISABLE_CH` - set `true` to disable writing to ClickHouse
- ClickHouse connection
  - `CH_HOST` - host (127.0.0.1 by default)
  - `CH_PORT` - port (9000 by default)
  - `CH_DEBUG` - debug enabled true/false (false by default)
  - `CH_USER` - user (empty=default by default)
  - `CH_PASSWORD` - password (nothing by default)
  - `CH_DB` - database (empty=default by default)

## REDIRECTS
Syntax:
```
<uri1> <redirect1>[|<uri2> <redirect2>[|...]]]
```
`uri_` may contains:
- `/*<word>` - at the end with any word (actually ignored), with fits `/`, `/any`, `/any/` `/any/thing` etc.

`redirect_` may contains:
- `{URI}` - replases with uri from request
- `{HOST}` - replases with host from request

Example:
`/ https://{HOST}/y|/photos https://jsonplaceholder.typicode.com{URI}|/photos/*a https://jsonplaceholder.typicode.com{URI}|/y http://yandex.ru`
> should be in one line, strongly like in syntax
- `/ https://{HOST}/y`
  - `curl -i http://127.0.0.1:8080` -> `https://127.0.0.1:8080/y`
- `/photos https://jsonplaceholder.typicode.com{URI}`
  - `curl -i http://127.0.0.1:8080/photos` -> `https://jsonplaceholder.typicode.com/photos`
  - `curl -i http://127.0.0.1:8080/photos?id=101` -> `https://jsonplaceholder.typicode.com/photos?id=101`
- `photos/*a https://jsonplaceholder.typicode.com{URI}`
  - `curl -i http://127.0.0.1:8080/photos/101` -> `https://jsonplaceholder.typicode.com/photos/101`
- `/y http://yandex.ru`
  - `curl -i http://127.0.0.1:8080/y` -> `http://yandex.ru`

## Attention!
Do not make redirect loops. Folks haven't like this ;)
## Usage

```
docker run -d --restart always \
    -e CH_HOST=10.0.0.1 \
    -e CH_PASSWORD="password" \
    -e REDIRECTS="/ http://google.com|/y http://yandex.ru" \
    -p 8080:8080 delfer/go-static-redirector
```
Open http://10.0.0.1/ in you browser or by `curl` to make new visit,  
open http://10.0.0.1/load to get current buffer usage

## License

MIT
