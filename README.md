[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/udhos/forward/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/udhos/forward)](https://goreportcard.com/report/github.com/udhos/forward)
[![Go Reference](https://pkg.go.dev/badge/github.com/udhos/forward.svg)](https://pkg.go.dev/github.com/udhos/forward)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/forward)](https://artifacthub.io/packages/search?repo=forward)
[![Docker Pulls forward](https://img.shields.io/docker/pulls/udhos/forward)](https://hub.docker.com/r/udhos/forward)

# forward

[forward](https://github.com/udhos/forward) forwards requests to other services.

* [Build](#build)
* [Usage](#usage)
* [Docker](#docker)
* [Helm chart](#helm-chart)

## Build

```
git clone https://github.com/udhos/forward

cd forward

./build.sh
```

## Usage

Send this request body to `forward`.

```
body: put-actual-body-to-forward-here
method: GET
set_headers:
  key1: value1
  key2: value2
url: http://localhost:2000/v1/hello # put full request URL here
```

### Example for Google

```
# start forward
forward
```

```
# call forward

$ more sample/google.yaml
url: https://www.google.com/search?q=golang
method: GET

$ curl --data-binary @sample/google.yaml localhost:8080/forward
<output omitted>
```

### Example for miniapi

Forward to [miniapi](https://github.com/udhos/miniapi).

```
# start miniapi on port 2000
export ADDR=:2000
miniapi
```

```
# start forward
forward
```

```
# call forward

$ more sample/miniapi.yaml
body: aaaaa
method: PUT
set_headers:
  a: b
url: http://localhost:2000/v1/hello

$ curl --data-binary @sample/miniapi.yaml localhost:8080/forward

{"request":{"headers":{"A":["b"],"Accept-Encoding":["gzip"],"Content-Length":["5"],"User-Agent":["Go-http-client/1.1"],"X-B3-Sampled":["1"],"X-B3-Spanid":["23e1685fce0c9374"],"X-B3-Traceid":["1c8318e574e71c17b61eec3814ede924"]},"method":"PUT","uri":"/v1/hello","host":"localhost:2000","body":"aaaaa","form_query":{},"form_post":{},"parameters":{"param1":"","param2":""}},"message":"ok","status":200,"server_hostname":"ubuntu","server_version":"1.0.5"}
```

## Docker

Docker hub:

https://hub.docker.com/r/udhos/forward

Pull from docker hub:

```
docker pull udhos/forward:0.0.0
```

Build recipe:

```
./docker/build.sh
```

## Helm chart

### Using the repository

See https://udhos.github.io/forward/.

### Create

```
mkdir charts
cd charts
helm create forward
```

Then edit files.

### Lint

```
helm lint ./charts/forward --values charts/forward/values.yaml
```

### Test rendering chart templates locally

```
helm template forward ./charts/forward --values charts/forward/values.yaml
```

### Render templates at server

```
helm install forward ./charts/forward --values charts/forward/values.yaml --dry-run
```

### Generate files for a chart repository

A chart repository is an HTTP server that houses one or more packaged charts.
A chart repository is an HTTP server that houses an index.yaml file and optionally (*) some packaged charts.

(*) Optionally since the package charts could be hosted elsewhere and referenced by the index.yaml file.

    docs
    ├── index.yaml
    └── forward-0.1.0.tgz

See script [update-charts.sh](update-charts.sh):

    # generate chart package from source
    helm package ./charts/forward -d ./docs

    # regenerate the index from existing chart packages
    helm repo index ./docs --url https://udhos.github.io/forward/

### Install

```
helm install forward ./charts/forward --values charts/forward/values.yaml
```

### Upgrade

```
helm upgrade forward ./charts/forward --values charts/forward/values.yaml
```

### Uninstall

```
helm uninstall forward
```
