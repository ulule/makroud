# sqlxx

**[sqlx](https://github.com/jmoiron/sqlx) superset.**

## Installation

```console
$ go get -u github.com/ulule/sqlxx
```

## API

### `GetSchema(model) (*Schema, error)`

Returns `model` schema (your model must be conform to `Model` interface).

```go
schema, err := sqlxx.GetSchema(model)
if err != nil {
    log.Fatal(err)
}

// The model primary key
pk := schema.PrimaryKey

// The model fields
fields := schema.Fields

// The model relations
relations := schema.Relations
```

### `GetByParams(db, out, params) error`

Executes a `WHERE` query with `params`, returning the first matching result into `out` interface.

```go
user := User{}

if err := sqlxx.GetByParams(db, &user, map[string]interface{}{"username": "jdoe"}); err != nil {
    fmt.Println(user.Username)
}
```

### `FindByParams(db, out, params) error`

Executes a `WHERE` query with `params`, returning all matching results into `out` interface.

```go
users := []User{}

if err := sqlxx.FindByParams(db, &users, map[string]interface{}{"is_active": true}); err != nil {
    for _, user := range users {
        fmt.Println(user)
    }
}
```

### `Save(db, out) error`

Executes either an `INSERT` or an `UPDATE` on `out` instance values, depending on primary key existance.

```go
// INSERT

// Here, no primary key yet
user := User{Username: "jdoe"}

if err := sqlxx.Save(db, &user); err != nil {
    // This user has been created
    fmt.Println(user)
}

// UPDATE

// Here, we already have a primary key
fmt.Println(user.ID)

// Let's update the username
user.Username = "johndoe"

if err := sqlxx.Save(db, &user); err != nil {
    // This user has been updated. Username is now "johndoe".
    fmt.Println(user)
}
```

### `Delete(db, out) error`

Executes a `DELETE` on `out` instance primary key.

```go
user := User{Username: "jdoe"}

// Create a user
if err := sqlxx.Save(db, &user); err != nil {
    fmt.Println(user)
}

// Delete it
if err := sqlxx.Delete(db, &user); err != nil {
    fmt.Println(user)
}
```

### `Archive(db, out, field) error`

Executes an `UPDATE` on `field` value from `out` instance.

```go
user := User{Username: "jdoe"}

// Create user
if err := sqlxx.Save(db, &user); err != nil {
    fmt.Println(user)
}

// Archive it by setting deleted_at column
if err := sqlxx.Archive(db, &user, "DeletedAt"); err != nil {
    fmt.Println(user)
}
```

## Struct tags

`sqlxx` tags must be separated by a semicolon (example: `tag1; tag2; tag3;`).

| Key           | Type    | Value                                                   |
|---------------|---------|---------------------------------------------------------|
| `primary_key` | `bool`  | if `true`, field is consired as a primary key           |
| `ignored`     | `bool`  | if `true`, field is ignored                             |
| `default`     | `value` | Default value for the field (example: `default:now()`) |

Tags of type `bool` can be set as `key:true` or just `key` for implicit `true`.

Example:

```go
type User struct {
    // We set as true
    ID int `sqlxx:"primary_key:true"`

    // But it is the same as
    ID int `sqlxx:"primary_key"`
}
```

## Development

### Docker

The test suite is running on PostgreSQL. We use [Docker](https://docs.docker.com/install/) to create a running container using [scripts/database](scripts/database).

### Testing

To run the test suite, simply execute:

```
scripts/test
```

Also, you can execute the linter with:

```
scripts/lint
```

### Random

Because sometimes it's hard to think of a good test fixture, using generators can save your productivity.
This is a small list used to write unit test:

 * http://www.fantasynamegenerators.com/pet-cat-names.php
 * http://www.fantasynamegenerators.com/pet-owl-names.php
 * http://www.fantasynamegenerators.com/food-names.php
 * http://www.fantasynamegenerators.com/color-names.php
