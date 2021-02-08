# Book Lending

## How to run?

### Graphql

### Query
```
query User {
  user(id: "6bf35476-b943-44e4-bf2d-fc4ca3cf21d5") {
    name,
    email
  }
}

```

### Mutations

#### CreateUser

```
mutation CreateUser{
  createUser(input: {name: "Júpiter Stein", email: "jupiter.steinn@gmail.com"}){
    id,
    email,
    name,
    createdAt
  }
}

```

### AddBookToMyCollection

```
mutation AddBookToMyCollection{
  addBookToMyCollection(loggedUserId: "a3a226e4-673d-4695-ba52-592bce01ca0a", input: {
    title: "A revolução dos Bichos",
    pages: 184
  }){
    id,
    title,
    createdAt,
    pages
  }
}

```

### LendBook

```
mutation LendBook{
  lendBook(loggedUserId: "a3a226e4-673d-4695-ba52-592bce01ca0a", input: {bookId: "clean-code-2", toUserId: "20c70153-1a0a-4261-a448-9f23f23be825"}){
    lentAt,
    book{
      title,
      pages
      createdAt
    },
    toUser,
    fromUser
  }
}
```

### ReturnBook

```
mutation ReturnBook{
  returnBook(loggedUserId: "20c70153-1a0a-4261-a448-9f23f23be825", bookId: "a-revolucao-dos-bichos-2"){
    book{
      id,
     	title,
      createdAt,
      pages
    },
    toUser,
    fromUser
  }
}
```


