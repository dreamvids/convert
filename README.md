# convert

DreamVids video conversion server

## setup

Database parameters can be changed via the following environnement variables:

DV_DB_HOST, DV_DB_PORT, DV_DB_USER, DV_DB_PASSWORD, DV_DB_NAME

## manual

Make sure that Go is installed and that your GOPATH is set.

```
go get github.com/dreamvids/convert
$GOPATH/bin/convert
```

## docker

```
git clone https://github.com/dreamvids/convert.git
cd convert
docker build -t dreamvids/convert .
docker run --rm --link your-db-container:db -p 8001:8001 dreamvids/convert
```
