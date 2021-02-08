# Book Lending

## Como executar?

Inicialize o db usando o comando `make env-up`. Após isso, pode inicializar a aplicação tanto em um container docker quanto chamando a função main. A chamada do DB está isolada para separar a execução da dependência do db da execução da aplicação, assim, subir a aplicação se torna mais rápido.

#### Pelo container Docker:
Gera um container com a aplicação e usa as variáveis de ambiente do docker-compose. Use o comando `make run-docker`.

#### Pela função main:
Executa a função, usando as variáveis de ambiente do arquivo env.local. Use o comando `make run`


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
  addBookToMyCollection(loggedUserId: "user-id", input: {
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
  lendBook(loggedUserId: "user-id-owner-of-a-book", input: {bookId: "clean-code-2", toUserId: "user-id-to-lent-a-book"}){
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
  returnBook(loggedUserId: "user-id-who-lent-a-book", bookId: "a-revolucao-dos-bichos-2"){
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


