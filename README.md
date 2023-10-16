mkdir blog_backend
cd blog_backend
go mod init ivy.cave.local/sockleblu/blog_backend
go get github.com/99designs/gqlgen

go run github.com/99designs/gqlgen init

# To regenerate code after changes
go run github.com/99designs/gqlgen generate

# Create User
```
mutation CreateUser ($input: UserInput!) {
  createUser(input: $input) {
    id
    username
    email
    token
  }
}
```

Query Variables:
```
{
  "input": {
    "username": "test_username",
    "password": "test_password",
    "email": "test@email.com"
  }
}
```

# Create Article
```
mutation CreateArticleMutation($input: ArticleInput!) {
        createArticle(input: $input) {
            id
    				slug
            title
            content
            user {
                id
                username
                email
            }
        }
    }
```

Query Variables:
```
{
  "input": {
    "title": "hello there!",
    "tags": [
      {
       "tag": "new_tag"
      }
    ],
    "content": "Umm hi..."
  }
}
```

Header Variables:
```
{
 "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDk0ODc1MTksInVzZXJuYW1lIjoic29ja2xlYmx1MSJ9.uYQiGtS3geXkX8tP1Vmk4KO5LaRMJYxMNPVF0iKTVdo"
}
```

# Delete Article
```
mutation deleteArticleByID ($articleId: Int!) {
  deleteArticleByID(articleId: $articleId)
}
```

Query Variables:
```
{
  "articleId": 1
}
```

# Login
```
mutation Login ($input: LoginInput!) {
  login(input: $input) {
    id
    username
    email
    token
  }
}
```

Query Variables:
```
{
  "input": {
    "username": "test_username",
    "password": "test_password",
  }
}
```

# Update Article
```
mutation updateArticle($articleId: Int!, $input: ArticleInput!) {
  updateArticle(articleId: $articleId, input: $input) {
    id
    content
    slug
  }
}
```

Query Variables:
```
{
  "articleId": 22,
  "input": {
    "content": "hey there!",
    "title": "new title",
    "tags": []
  }
}
```
