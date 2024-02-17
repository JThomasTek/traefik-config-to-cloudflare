# Cloudflare Traefik Control

Cloudflare Traefik Control (CTC) will read your Traefik config.yaml routes and generate DNS records in Cloudflare. It also acts as a dynamic DNS for the DNS entries it creates.

This application can be ran using Docker (recommended) or ran locally using CLI. To function, the application expects environment variables to be set for configuration.

## Running the App

It is recommended that you run the application using Docker to make it easier to upgrade if new features are added or security patches are implemented. Below you will find a examples of running the application with a simple Docker CLI command and a `docker-compose.yaml` file example for use with Docker Compose.

### Docker CLI

``` bash
docker run -d --restart always -v path/to/local/ctc/dir:/etc/ctc -v /etc/traefik:/etc/traefik ghcr.io/jthomastek/cloudflare-traefik-control:latest
```

### Docker Compose

``` yaml
version: "3"

volumes:
    ctc-dir:

services:
    cloudflare-traefik-control:
        image: ghcr.io/jthomastek/cloudflare-traefik-control:latest
        restart: always
        volumes:
            - ctc-dir:/etc/ctc
            - /etc/traefik:/etc/traefik
```

## Environment Variables

|     Variable Name    |       Default Value      | Required |
| -------------------- | ------------------------ | -------- |
| CLOUDFLARE_API_TOKEN |                          |    Yes   |
|  CLOUDFLARE_ZONE_ID  |                          |    Yes   |
|      LOG_LEVEL       |                          |    No    |
|  TRAEFIK_CONFIG_FILE | /etc/traefik/config.yaml |    Yes   |

At this time, the application only supports providing a Cloudflare API token. If there is large interest in there being API key support then it will be implemented.
