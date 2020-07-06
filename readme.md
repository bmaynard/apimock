# APIMock

APIMock allows you to serve fake REST API requests for your application so you are able to test your API requests. It works with any programming language as it runs a HTTP(S) server and with domain support, allows for zero configuration changes to your application. APIMock can also run a proxy server that will record any API requests your application makes and then creates a mock file for you!

# Table of Contents

- [Features](#Features)
- [Getting Started](#Getting-Started)
  - [Mock Server](#Mock-Server)
  - [Proxy Server](#Proxy-Server)
    - [Kubernetes Sidecar Proxy](#Kubernetes-Sidecar-Proxy)
  - [Configuration File](#Configuration-File)
- [Writing your own mocks](#Writing-your-own-mocks)
  - [Sample](#Sample)
- [Roadmap](#Roadmap)

## Features

- HTTP(S) server that listens for requests and returns the mock response.
- Domain support, separate your mocks by domain so you don't have to change any configuration to your code.
- Can run on Windows/Mac/Linux.
- Proxy server to record your applications REST API requests so you don't have to write any mocks.
- Kubernetes support. Run the proxy as a sidecar to record API requests.

## Getting Started

By default, the application will look for a `./mocks` directory in the current directory you are in. You can also create your mocks directory elsewhere and specify the path through the configuration file. The API mock server has the ability to serve different requests depending on the domain being requested so you have the ability to have different mocks for different services that have the same path. See [Writing your own mocks](#Writing-your-own-mocks) for more information. 

### Mock Server

To start the mock server to serve your API mocks, run:

```bash
$ apimock server
```

You can also specify different configuration options to change the listen address or run the server under TLS. For full list of configuration options, see the `help` command.

```
$ apimock server --help
Run the mock server to respond to API requests

Usage:
  apimock server [flags]

Flags:
  -a, --addr string            The listen addr eg: 127.0.0.1:8000 (default "127.0.0.1:8000")
  -g, --graceful-timeout int   The time for which the server will gracefully wait for existing connections to finish (default 15)
  -h, --help                   help for server
  -k, --keyPath string         The path to the key file
  -p, --pemPath string         The path to the pem file

Global Flags:
      --config string   config file (default is $HOME/.apimock.yaml)
```

### Proxy Server

You can run a proxy server to intercept your API requests which will record and save the responses as mock so you don't have to write any mocks yourself. Currently only HTTP traffic is supported to record the response but the server will proxy your TLS requests.

When the proxy captures a HTTP request, it will create the mock file in the mocks directory under a folder of the domain name for that request. The proxy will create a MD5 hash of the response so while you can have different responses for the same request, it will only create one file per response content. You can also run the proxy under a Kubernetes mode that acts as sidecar proxy which in this scenario, will re-write all requests to the application pod using `localhost`. See [Kubernetes Sidecar Proxy](#Kubernetes-Sidecar-Proxy) for more information.

```bash
$ apimock proxy
```

You can also specify different configuration options to change the listen address or run the proxy under TLS. For full list of configuration options, see the `help` command.

```
$ apimock proxy --help
Run a proxy server to capture requests and save as mock files

Usage:
  apimock proxy [flags]

Flags:
  -a, --addr string      The listen address (default: 127.0.0.1:8888) (default "127.0.0.1:8888")
  -h, --help             help for proxy
  -K, --k8s              Is a proxy for a Kubernetes service (will rewrite all requests to localhost)
  -k, --keyPath string   The path to the key file
  -p, --pemPath string   The path to the pem file

Global Flags:
      --config string   config file (default is $HOME/.apimock.yaml)
```

#### Kubernetes Sidecar Proxy

To run the proxy as a sidecar in your Kubernetes deployment, you can use the [bmaynard/apimock-proxy-kubernetes](https://hub.docker.com/r/bmaynard/apimock-proxy-kubernetes) docker image. You can change `SERVICE_HOST_NAME` to the name of your service so it saves the mocks under the correct domain name for the service, otherwise it will save them in the `localhost` folder. You can then use `kubectl cp pod-name:/app/mocks/ mocks/ -c sidecar-proxy` command to copy the mocks to your local filesystem.

```yaml
...
    spec:
      containers:
      - name: sidecar-proxy
        image:  bmaynard/apimock-proxy-kubernetes
        imagePullPolicy: Always
        env:
        - name: SERVICE_HOST_NAME
          value: "service-host"
        ports:
        - containerPort: 8888
     - name: your-application-container
....
```

### Configuration File

The application will look for a `.apimock.yaml` in your home directory and will load the file if it exists. Currently, you are able to supply the path to your mocks directory. Future enhancements will include being able to store your mock files in S3.

**Sample:**

```yaml
---
filesystem:
  adapter: local
  meta:
    mock_path: "/path/to/your/mocks/directory"
```

## Writing your own mocks

To write your own mocks, you will need to create a base directory that will hold all the different domains your API requests will be served from. You can also create a directory called `_all_` which those mocks will be served under any domain name. A sample directory structure might look like:

```
mocks
├── _all_
│   ├── request-one.json
│   ├── request-two.json
└── articles.apimock.benmaynard.dev
|   ├── articles.json
|   └── article_1.json
└── users.apimock.benmaynard.dev
    ├── users.json
    └── user_1.json
```

A DNS record has been created for `*.apimock.benmaynard.dev` that points to `127.0.0.1` so you can request mocks for different domains on your local machine. **Note:** You will not be able to request the API mock through a browser under the .dev domain unless you use a valid SSL certificate (See: https://get.dev/)

A mock contains two top level sections, `response` and `meta`. The `response` key is the JSON you wish to return and `meta` contains information about the request e.g. the `status_code`, `request_path` and `method`. In the request path, you can use variables in the path so you don't have to create a mock for every possible request. The application is built on top of [gorilla/mux](https://github.com/gorilla/mux), see [https://github.com/gorilla/mux#registered-urls](https://github.com/gorilla/mux#registered-urls) for more information.

### Sample

```json
{
    "response": {
        "articles": [
            {
                "id": 1,
                "title": "this is a test"
            },
            {
                "id": 2,
                "title": "this is a test two"
            }
        ]
    },
    "meta": {
        "status_code": 200,
        "request_path": "/articles",
        "method": "GET"
    }
}
```

## Roadmap

- [ ] S3 Support
- [ ] Generate Host files
- [ ] Live reloading of mock files
- [ ] Request Input Matching
- [ ] Capture TLS proxy requests
- [ ] gRPC Support
- [ ] XML/Other media types Support
- [ ] Generate Kubernetes YAML files
- [ ] Add delay to responses
- [ ] CORS configuration
- [ ] Ability to specify customer headers