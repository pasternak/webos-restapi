# webos-restapi

LG's WebOS REST API with IFTTT (Google assistant + webhook) integration to control TV using your voice.

## Getting started

### Prerequisites

webos-restapi requires ssl key and crt files to be present in the workdir. You can either create them yourself:

```bash
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

or use any SSL provider. [Let's Encrypt](https://letsencrypt.org/) should be just fine. I recomend using [acme.sh](acme.sh) for that.

### Installing

1. Using source code

    ```bash
    git clone https://github.com/pasternak/webos-restapi.git
    cd webos-restapi
    ./build.sh
    ```
    Above process will build a binary for your current platform. If you prefer to crosscompile for either linux or arm (raspberry pi), run
    ```bash
    ./build.sh (linux|pi)
    ```

1. Using shipped binary in release tab:

    TODO: 
    ```bash
    curl -O ...
    ```

### Deployment

Using the binary (cert and key have to be present in the same directory as `webos`)

```bash
./webos -username=<username> -password=<password>
```

Using docker container

```bash
docker run -d pasternak/webos-restapi -e WEBOS_USERNAME=<username> -e WEBOS_PASSWORD=<password> -v <PATH_TO_SSL_FILES>:/ --name=webos-restapi
```

### IFTTT + Google assistant integration

TODO: add screenshots