# urlshortner

URLshortner is webserver providing functionality to create short url for given url. Application is written in Domain-Driven-Design architecture. It uses inmemory repository to create new records. 

## Setup

### Run on local machine
```
go mod tidy
go test ./...
go run cmd/apiserver/main.go

```

### Building docker image 
```
docker build -t app .
docker tag app urlshortner:latest
```

## Run Server using docker
```
docker run -d -p 8080:8080 urlshortner:latest
```

## API Documentation 

Server exposes following API: 

- GET /urlshortnerservice/v1/healthcheck

    API for backend server

- POST /urlshortnerservice/v1/url

    API for creating new short url 

    payload :  {"url": < string > }

- GET /< shorturl >

    API for redirecting given short url to source url location. 
    If short-url found then it redirect to desire location with HTTP code 302.

- GET /urlshortnerservice/v1/metrics?headcount=3

    API implements metrics for shortner web server. In current implementation, it gives top url request received on web server. 

### Sample API Demo
1.  Healthcheck API

```
curl --location 'http://localhost:8080/urlshortnerservice/v1/metrics?headcount=3'
```
response :
```
{
    "result": "ok",
    "success": true
}
```

2. Create New short URL
```
curl --location 'http://localhost:8080/urlshortnerservice/v1/url' \
--header 'Content-Type: application/json' \
--data '{
    "url": "https://timesofindia.indiatimes.com/india/chandrayaan-3-mission-nasas-lro-captures-vikram-landing-site/articleshow/103425084.cms?from=mdr"
}'
```
Response:
```
{
    "result": {
        "shortUrl": "30TI2j"
    },
    "success": true
}
```

3. Get source url for given short url

```
curl --location 'http://localhost:8080/30TI2j'
```
Note : use above url in browser , It will redirect to original source url page

4. Get Metrics 
```
http://localhost:8080/urlshortnerservice/v1/metrics?headcount=3
```
Response: 
```
{
    "result": {
        "data": [
            {
                "domain": "timesofindia.indiatimes.com",
                "count": 1
            }
        ]
    },
    "success": true
}
```

