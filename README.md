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
# Write your query or mutation here
mutation CreateUser{
  createUser(input: {name: "Júpiter Stein", email: "jupiter.steinn@gmail.com"}){
    id,
    email,
    name,
    createdAt
  }
}

```
