# codelabs-preview-go

## development

prepare the following environment 
- `GOOGLE_APPLICATION_CREDENTIALS`
- `GOOGLE_CLIENT_ID`
- `GOOGLE_CLIENT_SECRET`
- `GOOGLE_REDIRECT_URL`

```bash
go run ./cmd/playground/main.go
```

```bash
curl localhost:3000/?file_id=1tkNrHr_ZnWhsPhrVEP3zYcyZJSD7w502atugh300EEA
```

## APIs

### DRAFT

#### request
`POST http://localhost:3000/v/`

__body__

```json
{
    "data": {
        "title": "super cool learning",
        "summary": "summary super cool learning",
        "slug": "super_cool_learning",
        "type": "testing",
        "tags": "super,cool,learning",
        "status": "draft",
        "feedbackLink": "http://example.com",
        "author": "super@cool.learning",
        "authorLDAP": "LDAP://super.cool.learning",
        "analyticsAccount": "GA100200300"
    }
}
```

#### response
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "fileId": "1gGHmoegaMH3anPmvvZxM0eVrxoRVAUe1MjezZMFGDiE"
    }
}
```

### PREVIEW
#### request
`GET http://localhost:3000/v/{{fileid}}/preview`

#### response
`codelabs content in html`

### PUBLISH
#### request
`POST http://localhost:3000/v/{{fileid}}`

#### response
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "revision": 1,
        "meta": {
            "exportedDate": "2020-12-19T20:32:50.840291+07:00",
            "fileId": "1gGHmoegaMH3anPmvvZxM0eVrxoRVAUe1MjezZMFGDiE",
            "meta": {
                "authors": "super@cool.learning",
                "category": [
                    "${category}"
                ],
                "duration": 0,
                "feedback": "http://example.com",
                "ga": "GA100200300",
                "id": "super-cool-learning",
                "source": "",
                "status": [
                    "draft"
                ],
                "summary": "summary super cool learning",
                "tags": [
                    "cool",
                    "learning",
                    "super"
                ],
                "theme": "category-",
                "title": "super cool learning",
                "totalChapters": 1,
                "url": "super-cool-learning"
            },
            "revision": 1
        }
    }
}
```

### VIEW
#### request
`GET http://localhost:3000/v/{{fileid}}/`

`GET http://localhost:3000/v/{{fileid}}/latest/`

`GET http://localhost:3000/v/{{fileid}}/{{revsion}}/`

#### response
`codelabs content in html`

### META
#### request
`GET http://localhost:3000/v/{{fileid}}/meta/latest`

`GET http://localhost:3000/v/{{fileid}}/meta/{{revsion}}`

#### response
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "meta": {
            "exportedDate": "2020-12-19T20:32:50.840291+07:00",
            "fileId": "1gGHmoegaMH3anPmvvZxM0eVrxoRVAUe1MjezZMFGDiE",
            "meta": {
                "authors": "super@cool.learning",
                "category": [
                    "${category}"
                ],
                "duration": 0,
                "feedback": "http://example.com",
                "ga": "GA100200300",
                "id": "super-cool-learning",
                "source": "",
                "status": [
                    "draft"
                ],
                "summary": "summary super cool learning",
                "tags": [
                    "cool",
                    "learning",
                    "super"
                ],
                "theme": "category-",
                "title": "super cool learning",
                "totalChapters": 1,
                "url": "super-cool-learning"
            },
            "revision": 1
        }
    }
}
```
