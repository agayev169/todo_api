# API for TODO list app

Basic API for TODO list appliactions. It receives your requests, handles them and places the DB (see [Usage](#usage)) in a file to restore it after its restart. All the requests and responses are sent in JSON format.

## Install
```go get github.com/agayev169/todo_api```
will download and install the api to your local machine into your `$GOPATH`.
Or simply download the executable from [Releases](https://github.com/agayev169/todo_api/releases/tag/v1.0).

## Usage
run the api with the following command:
```./todo_api[.exe] [-port PORT] [-db DB]```

The program has the following flags:
```
  -db string
        The name of the DB file (default "db.json")
  -port int
        The port number to run the API on (default 1337)
```

In order to get usage of the program just use
```./todo_api[.exe] -help```

## Types
`TodoItem`
| Field    | Type   | Description                        |
| -------- | ------ | ---------------------------------- |
| id       | int    | Unique identifier of the TODO item |
| name     | string | Name of the task                   |
| priority | int    | Priority of task from 1 to 10      |

## Endpoints
The API has the following endpoints:
- `/add` adds a new TODO item to DB. It expects `TodoItem` to add. If there is an item with the same ID, you will receive a response with `ok=false` and `description="Repeated ID"` (see [Response](#response)). In case of succesful addition of the item, you will received the added item (of `TodoItem` type) in the `result` field of the response.
- `/remove` removes an item from DB. It expects the ID of a `TodoItem` to remove. If there is no item with the ID received, you will receive a response with `ok=false` and `description="ID not in DB"` (see [Response](#response)). In case of succesful deletion of the item, you will receive the removed item (of `TodoItem` type) in the `result` field of the response.
- `/get` returns a `TodoItem` from DB. It expects the ID of a `TodoItem` to return. If there is no item with the ID received, you will receive a response with `ok=false` and `description="Not in DB"` (see [Response](#response)). Otherwise, the requested `TodoItem` will be placed in the `result` field of the response.
- `/getAll` returns a list of all `TodoItem`s from DB.

## Response
Response has the following fields:
- `ok (bool)` showing whether the request was OK.
- `description (string)` contains the message for the request. Can be empty if `ok=true`. If `ok=false` contains the error message.
- `result` contains the result of the request itself. Is `TodoItem` for `/add`, `remove` and `get`, and `[TodoItem]` for `/getAll` 

## Examples
Request: (Adding new item)
```/add
{
    "id": 1,
    "name": "Finish API",
    "priority": 5
}
```
Response:
```
{
    "ok": true,
    "decsription": "Added to the DB",
    "result": {
        "id": 1,
        "name": "Finish API",
        "priority": 5
    }
}
```
Request: (Adding another item)
/add 
```
{
    "id": 2,
    "name": "Learn GO",
    "priority": 7
}
```
Response: 
```
{
    "ok": true,
    "decsription": "Added to the DB",
    "result": {
        "id": 2,
        "name": "Learn GO",
        "priority": 7
    }
}
```
Request: (Adding an item with existing ID)
/add 
```
{
    "id": 2,
    "name": "Sleep",
    "priority": 1
}
```
Response:
```
{
    "ok": false,
    "decsription": "Repeated ID",
    "result": {
        "id": 2,
        "name": "Learn GO",
        "priority": 7
    }
}
```
Request: (Adding an item with invalid priority value)
/add
```
{
    "id": 1, 
    "name": "Learn golang", 
    "priority": 10
}
```
Response:
```
{
    "ok": false,
    "decsription": "Invalid request",
    "result": null
}
```
Request: (Getting an item)
```
/get 
{
    "id": 2
}
```
Response:
```
{
    "ok": true,
    "decsription": "Retrieved succesfully",
    "result": {
        "id": 2,
        "name": "Learn GO",
        "priority": 7
    }
}
```
Request: (Getting a non-existent item)
```
/get
{
    "id": 1337
}
```
Response:
```
{
    "ok": false,
    "decsription": "Not in DB",
    "result": null
}
```
Request: (Getting all items)
```
/getAll
```
Response:
```
{
    "ok": true,
    "decsription": "Retrieved succesfully",
    "result": [
        {
            "id": 1,
            "name": "Finish API",
            "priority": 5
        },
        {
            "id": 2,
            "name": "Learn GO",
            "priority": 7
        }
    ]
}
```
Request: (Removing item)
```
{
    "id": 1
}
```
Response:
```
{
    "ok": true,
    "decsription": "Removed from DB",
    "result": {
        "id": 1,
        "name": "Finish API",
        "priority": 5
    }
}
```
Request: (Removing item)
```
{
    "id": 2
}
```
Response:
```
{
    "ok": true,
    "decsription": "Removed from DB",
    "result": {
        "id": 2,
        "name": "Learn GO",
        "priority": 7
    }
}
```
Request: (Removing a non-existent item)
```
{
    "id": 123
}
```
Response:
```
{
    "ok": false,
    "decsription": "ID not in DB",
    "result": null
}
```
