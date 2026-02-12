# taskusama

## About

**taskusama** is a project tracking software with web access

## Pre-reqs

If you want to build and run **taskusama**, then you want to have:

1. Docker
2. Go
3. Make

## How to build via the Makefile

Build the Docker image: `make docker`

## How to run in development mode

1. For development purposes, generate snakeoil certificates: 
```
sudo make-ssl-cert generate-default-snakeoil --force-overwrite
```

2. Copy the certificates to the directory containing the `docker-compose.yml`
```
cp /etc/ssl/certs/ssl-cert-snakeoil.pem .
sudo cp /etc/ssl/private/ssl-cert-snakeoil.key .
```

3. Run the Docker containers: `docker compose up -d`

## How to use

Navigate in a web browser to https://localhost/issues
