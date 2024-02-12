# go-url-shortener

## Creating a URL shortener service in Golang using Redis as a database.

### Setup
- Clone the repo

- Install dependencies
```
go mod download
```
- Run the server using docker-compose
```
docker-compose up -d --build
```
### API Documentation
#### Create a shortened URL
Send a POST request to the /api/v1 endpoint with the following body to create a shortened URL

```
curl --location 'http://localhost:3000/api/v1' \
--header 'Content-Type: application/json' \
--data '{
    "url": "https://www.youtube.com/"
}'
```

Response
```
{
    "url": "https://www.youtube.com/",
    "short": "localhost:3000/fcc056",
    "expiry": 24,
    "rate_limit": 9,
    "rate_limit_resets": 30
}
```
Once the limit is reached, the same API will return the following response:

```
{
    "error": "Rate limit exceeded",
    "rate_limit_reset": 28
}
```
