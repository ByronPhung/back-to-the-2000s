# back-to-the-2000s
A short view back to the past when gamers were all about the forums... at least that's what the theme of this project reminds me of lol.

In all seriousness, this introductory project explores implementing a standalone API in Go, which you can ultimately test and play around with following the instructions below.

# Prerequisites
* **An installation of Go 1.17 or later.** For installation instructions, see [Installing Go](https://go.dev/doc/install).
* **A command terminal.** Your installation of Go should already have setup for any Linux and Mac terminal as well as PowerShell and cmd on Windows. If you're missing this for whatever reason, then you can double-check [Go's Installation Page](https://go.dev/doc/install) and attempt to reinstall one more time.
* **The curl tool.** So we can actually test and use our program's server :). This should already be installed on Linux and Mac as well as Windows 10 Insider build 17063 and later. If absolutely needed, you can [download curl directly from the main website](https://curl.se/download.html).

# Execution

## Starting the Server
1. Clone this repository to any folder/workspace of your liking.
2. Open up a command terminal and change directory to the clone. For example:
```
cd D:\Workspace\back-to-the-2000s
```
3. If you haven't already or are running into any issues with unresolved dependencies, then you can run the following command to resolve and download them:
```
go get .
```
4. Run the following command in your terminal to start the standalone Go HTTP server:
```
go run main.go
```

## Using the API

This entire section relies on the curl tool as mentioned above in the Prequisities sesction. You can optionally use a UI tool, such as [Postman](https://www.postman.com), but for the sake of the most common use case and simplicity, the following instructions will be using the command terminal + the curl tool.

1. Open up a new command terminal. The current directory does not matter (assuming your curl tool is already symlinked and available in your environment variables).
2. The structure of the available API is as follows that you can run and test:
```
http://localhost:8080/v1/user-posts/:userId
```
For example, if you wanted to attempt to query a user with userId=4, then you can use the following command:
```
curl http://localhost:8080/v1/user-posts/4
```
which would return the following response based on the current mock data:
```
$ curl http://localhost:8080/v1/user-posts/4
{
    "id": 4,
    "userInfo": {
        "name": "Patricia Lebsack",
        "username": "Karianne",
        "email": "Julianne.OConner@kory.org"
    },
    "posts": [
        {
            "id": 31,
            "title": "ullam ut quidem id aut vel consequuntur",
            "body": "debitis eius sed quibusdam non quis consectetur vitae\nimpedit ut qui consequatur sed aut in\nquidem sit nostrum et maiores adipisci atque\nquaerat voluptatem adipisci repudiandae"
        },
        {
            "id": 32,
            "title": "doloremque illum aliquid sunt",
            "body": "deserunt eos nobis asperiores et hic\nest debitis repellat molestiae optio\nnihil ratione ut eos beatae quibusdam distinctio maiores\nearum voluptates et aut adipisci ea maiores voluptas maxime"       
        },
        {
            "id": 33,
            "title": "qui explicabo molestiae dolorem",
            "body": "rerum ut et numquam laborum odit est sit\nid qui sint in\nquasi tenetur tempore aperiam et quaerat qui in\nrerum officiis sequi cumque quod"
        },
        {
            "id": 34,
            "title": "magnam ut rerum iure",
            "body": "ea velit perferendis earum ut voluptatem voluptate itaque iusto\ntotam pariatur in\nnemo voluptatem voluptatem autem magni tempora minima in\nest distinctio qui assumenda accusamus dignissimos officia nesciunt nobis"
        },
        {
            "id": 35,
            "title": "id nihil consequatur molestias animi provident",
            "body": "nisi error delectus possimus ut eligendi vitae\nplaceat eos harum cupiditate facilis reprehenderit voluptatem beatae\nmodi ducimus quo illum voluptas eligendi\net nobis quia fugit"
        },
        {
            "id": 36,
            "title": "fuga nam accusamus voluptas reiciendis itaque",
            "body": "ad mollitia et omnis minus architecto odit\nvoluptas doloremque maxime aut non ipsa qui alias veniam\nblanditiis culpa aut quia nihil cumque facere et occaecati\nqui aspernatur quia eaque ut aperiam inventore"
        },
        {
            "id": 37,
            "title": "provident vel ut sit ratione est",
            "body": "debitis et eaque non officia sed nesciunt pariatur vel\nvoluptatem iste vero et ea\nnumquam aut expedita ipsum nulla in\nvoluptates omnis consequatur aut enim officiis in quam qui"
        },
        {
            "id": 38,
            "title": "explicabo et eos deleniti nostrum ab id repellendus",
            "body": "animi esse sit aut sit nesciunt assumenda eum voluptas\nquia voluptatibus provident quia necessitatibus ea\nrerum repudiandae quia voluptatem delectus fugit aut id quia\nratione optio eos iusto veniam iure"
        },
        {
            "id": 39,
            "title": "eos dolorem iste accusantium est eaque quam",
            "body": "corporis rerum ducimus vel eum accusantium\nmaxime aspernatur a porro possimus iste omnis\nest in deleniti asperiores fuga aut\nvoluptas sapiente vel dolore minus voluptatem incidunt ex"
        },
        {
            "id": 40,
            "title": "enim quo cumque",
            "body": "ut voluptatum aliquid illo tenetur nemo sequi quo facilis\nipsum rem optio mollitia quas\nvoluptatem eum voluptas qui\nunde omnis voluptatem iure quasi maxime voluptas nam"
        }
    ]
}
```
If you want to test additional details, such as HTTP response codes, then you can add the `-v` flag to the curl command to get a verbose log, which will include details, such as headers, status code, etc.
```
$ curl -v -XGET 'http://localhost:8080/v1/user-posts/4'
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying ::1:8080...
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /v1/user-posts/4 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.71.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Tue, 11 Jan 2022 05:55:33 GMT
< Transfer-Encoding: chunked
<
{
    "id": 4,
    "userInfo": {
        "name": "Patricia Lebsack",
        "username": "Karianne",
        "email": "Julianne.OConner@kory.org"
    },
    "posts": [
        {
            "id": 31,
            "title": "ullam ut quidem id aut vel consequuntur",
            "body": "debitis eius sed quibusdam non quis consectetur vitae\nimpedit ut qui consequatur sed aut in\nquidem sit nostrum et maiores adipisci atque\nquaerat voluptatem adipisci repudiandae"
        },
        {
            "id": 32,
            "title": "doloremque illum aliquid sunt",
            "body": "deserunt eos nobis asperiores et hic\nest debitis repellat molestiae optio\nnihil ratione ut eos beatae quibusdam distinctio maiores\nearum voluptates et aut adipisci ea maiores voluptas maxime"
        },
        {
            "id": 33,
            "title": "qui explicabo molestiae dolorem",
            "body": "rerum ut et numquam laborum odit est sit\nid qui sint in\nquasi tenetur tempore aperiam et quaerat qui in\nrerum officiis sequi cumque quod"
        },
        {
            "id": 34,
            "title": "magnam ut rerum iure",
            "body": "ea velit perferendis earum ut voluptatem voluptate itaque iusto\ntotam pariatur in\nnemo voluptatem voluptatem autem magni tempora minima in\nest distinctio qui assumenda accusamus dignissimos officia nesciunt nobis"
        },
        {
            "id": 35,
            "title": "id nihil consequatur molestias animi provident",
            "body": "nisi error delectus possimus ut eligendi vitae\nplaceat eos harum cupiditate facilis reprehenderit voluptatem beatae\nmodi ducimus quo illum voluptas eligendi\net nobis quia fugit"
        },
        {
            "id": 36,
            "title": "fuga nam accusamus voluptas reiciendis itaque",
            "body": "ad mollitia et omnis minus architecto odit\nvoluptas doloremque maxime aut non ipsa qui alias veniam\nblanditiis culpa aut quia nihil cumque facere et occaecati\nqui aspernatur quia eaque ut aperiam inventore"
        },
        {
            "id": 37,
            "title": "provident vel ut sit ratione est",
            "body": "debitis et eaque non officia sed nesciunt pariatur vel\nvoluptatem iste vero et ea\nnumquam aut expedita ipsum nulla in\nvoluptates omnis consequatur aut enim officiis in quam qui"
        },
        {
            "id": 38,
            "title": "explicabo et eos deleniti nostrum ab id repellendus",
            "body": "animi esse sit aut sit nesciunt assumenda eum voluptas\nquia voluptatibus provident quia necessitatibus ea\nrerum repudiandae quia voluptatem delectus fugit aut id quia\nratione optio eos iusto veniam iure"
        },
        {
            "id": 39,
            "title": "eos dolorem iste accusantium est eaque quam",
            "body": "corporis rerum ducimus vel eum accusantium\nmaxime aspernatur a porro possimus iste omnis\nest in deleniti asperiores fuga aut\nvoluptas sapiente vel dolore minus voluptatem incidunt ex"
        },
        {
            "id": 40,
            "title": "enim quo cumque",
            "body": "ut voluptatum aliquid illo tenetur nemo sequi quo facilis\nipsum rem optio mollitia quas\nvoluptatem eum voluptas qui\nunde omnis voluptatem iure quasi maxime voluptas nam"
        }
    ]
}
```

## Test Data and Scenarios

### Valid User IDs

The current mock data has valid data for 10 user IDs starting from 1 and ending at 10, which should return a 200 Ok with a JSON response in the following format:
```
{
    "id": 4,
    "userInfo": {
        "name": "Patricia Lebsack",
        "username": "Karianne",
        "email": "Julianne.OConner@kory.org"
    },
    "posts": [
        {
            "id": 31,
            "title": "ullam ut quidem id aut vel consequuntur",
            "body": "debitis eius sed quibusdam non quis consectetur vitae\nimpedit ut qui consequatur sed aut in\nquidem sit nostrum et maiores adipisci atque\nquaerat voluptatem adipisci repudiandae"
        },
        ... // Can return any number of posts associated with this user.
    ]
}
```

### User ID Not Found

Any request with a non-existent (but valid integer) user ID should expect a 404 Not Found:
```
$ curl -v -XGET 'http://localhost:8080/v1/user-posts/123456'
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying ::1:8080...
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /v1/user-posts/123456 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.71.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 404 Not Found
< Content-Type: application/json; charset=utf-8
< Date: Tue, 11 Jan 2022 05:56:49 GMT
< Content-Length: 49
<
{
    "message": "Could not find userId=123456"
}
```

### Bad Request

This API validates that the userId input is in the expected integer format as that is the source data's model schema. If the API detects a non-integer input, then it should return a 400 Bad Request:
```
$ curl -v -XGET 'http://localhost:8080/v1/user-posts/test-123'
Note: Unnecessary use of -X or --request, GET is already inferred.
*   Trying ::1:8080...
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /v1/user-posts/test-123 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.71.1
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
< Date: Tue, 11 Jan 2022 05:57:38 GMT
< Content-Length: 78
<
{
    "message": "Expected ID in integer format, but got 'test-123' instead"
}
```

### All Other Errors

If for whatever reason the mock server is down or something has changed internally that drastically breaks the expected JSON models, then you should expect a 500 error code with the same "message" JSON blob and a detailed error message:
```
{
    "message": "Unexpected error creating client request for Cool Vendor's Get User API: error=blablablablab"
}
```

Some other explicitly handled errors are:
```
{
    "message": "Unexpected communication or client policy error occurred trying to fetch userId=123456 from Cool Vendor: error=blablablablab"
}
```
```
{
    "message": "Unable to parse response body as 'user' JSON for Cool Vendor's Get User By ID API: error=blablablablab"
}
```
```
{
    "message": "Unexpected error trying to read response body for server error trying to fetch userId=123456 from Cool Vendor: error=blablablablab"
}
```
```
{
    "message": "Unexpected server error occurred trying to fetch userId=123456 from Cool Vendor: error=blablablablab"
}
```
```
{
    "message": "Unexpected error creating client request for Cool Vendor's Get Posts API: error=blablablablab"
}
```
```
{
    "message": "Unexpected communication or client policy error occurred trying to fetch posts for userId=123456 from Cool Vendor: error=blablablablab"
}
```
```
{
    "message": "Unable to parse response body as '[]postSummary' JSON for Cool Vendor's Get Posts API: error=blablablablab"
}
```
```
{
    "message": "Unexpected error trying to read response body for server error trying to fetch posts for userId=123456 from Cool Vendor: error=blablablablab"
}
```
```
{
    "message": "Unexpected server error occurred trying to fetch posts for userId=123456 from Cool Vendor: error=blablablablab"
}
```

## Running Unit Tests

You can run unit tests by calling `go test` in the clone's main directory. For example:
```
cd D:\Workspace\back-to-the-2000s
go test
```
