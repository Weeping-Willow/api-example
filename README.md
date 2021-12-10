## Description

Creating an example api using golang

Note: .env is included in the repo and contains the default password auth details for development mongo instance incase, so this should only be used for dev or testing 

## Requirments

- port 8080 and 27017 free
- docker and docker compose installed 

## How To Run
- run command `dokcer-compose up -d` will start api on port 8080
## How To Use

1. set a header `Authorization` with value `Bearer complicated-token`

2. by default if you lunch the api using docker it should run on port 8080

3. to check if you've added auth headers correctly call this route `Get /api/`, if it returns 200 it's working ok, if it returns 401 the Authorization header is not set correctly

## Using the api

### Routes
/api routes need Authorization token 

- `GET /`
- `GET /api/`
- `GET /api/score`
Gets leaderboard from highest rank to lowest, has query keys -> `page`, `page[size]`, `name`

Curl example
```curl
curl --location -g --request GET 'localhost:8080/api/score?page=2&page[size]=20&name=NsGR' \
--header 'Authorization: Bearer complicated-token'
```

note you might have to change the name that matches some name in the leaderboard

response Example
```json
{
    "results": [
        {
            "score": 420,
            "name": "YzRy",
            "rank": 9
        },
        {
            "score": 417,
            "name": "RzLN",
            "rank": 10
        }
    ],
    "aroundMe": [
        {
            "score": 481,
            "name": "kjQZ",
            "rank": 1
        },
        {
            "score": 468,
            "name": "NsGR",
            "rank": 2
        },
        {
            "score": 464,
            "name": "GyRA",
            "rank": 3
        }
    ],
    "meta": {
        "current_page": 5,
        "total_pages": 25,
        "total_count": 50,
        "per_page": 2,
        "next_page": 6
    }
}
```
- `POST /api/score`
posts or updates scores in the leaderboard

body example
```json
{
    "name": "test",
    "score": 460
}
```

curl example
```curl
curl --location --request POST 'localhost:8080/api/score' \
--header 'Authorization: Bearer complicated-token' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "test",
    "score": 460
}'
```

response example
```json
{
    "score": 460,
    "name": "test",
    "rank": 4
}
```

## Note
- Atleast to my knowledge there really isn't a good way to get rankings in mongodb, this would be a lot easier to do with sql in something like Postgresql or MySQL, but as this is just a small exercise, sp the wrong choice for db doesn't really matter
- The Get route for scores doesn't have integration tests, same with the services it uses, but you can get the jist on my way of doing tests from looking at the other services or controllers
- Seeding is not that great either, as i didn't really wanna spend a lot of time one making something more proper
- Lots of the errors just return 500, which is not correct, but could be fixed by implementing the interface for service errors, the interface and error handler would have to be refactored a bit, but nothing major
