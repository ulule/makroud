# sqlxx

[![CircleCI][circle-img]][circle-url]
[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*A high level SQL Connector.*


## Introduction

TODO


## Installation

Using [dep](https://github.com/golang/dep)

```console
dep ensure -add github.com/ulule/sqlxx@master
```

or `go get`

```console
go get -u github.com/ulule/sqlxx
```


## Usage

### Create a Driver

A **Driver** is a high level abstraction of a database connection or a transaction. It's almost required everytime alongside a `context.Context` to manipulate rows.

```go
driver, err := sqlxx.New(
	sqlxx.Host(cfg.Host),
	sqlxx.Port(cfg.Port),
	sqlxx.User(cfg.User),
	sqlxx.Password(cfg.Password),
	sqlxx.Database(cfg.Name),
	sqlxx.SSLMode(cfg.SSLMode),
	sqlxx.MaxOpenConnections(cfg.MaxOpenConnections),
	sqlxx.MaxIdleConnections(cfg.MaxIdleConnections),
)
```

Also, you can use directly a struct if you don't need to use
[functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis):

```go
driver, err := sqlxx.NewWithOptions(&sqlxx.ClientOptions{
	Host:               cfg.Host,
	Port:               cfg.Port,
	User:               cfg.User,
	Password:           cfg.Password,
	Database:           cfg.Name,
	SSLMode:            cfg.SSLMode,
	MaxOpenConnections: cfg.MaxOpenConnections,
	MaxIdleConnections: cfg.MaxIdleConnections,
})
```


### Define a Model

With the [**Active Record**](https://en.wikipedia.org/wiki/Active_record_pattern) approach, you have to define a
model that wraps your database table _(or view)_ columns into properties.

Model are struct that contains basic go types, pointers, or `sql.Scanner`, `driver.Valuer` or `Model` interface.
All the fields of this struct will be columns in the database table.

#### An example

```go
type User struct {
	// Columns
	ID        string `sqlxx:"column:id,pk:ulid"`
	Email     string `sqlxx:"column:email"`
	Password  string `sqlxx:"column:password"`
	Country   string `sqlxx:"column:country"`
	Locale    string `sqlxx:"column:locale"`
	ProfileID string `sqlxx:"column:profile_id,fk:profiles"`
	// Relationships
	Group   []Group
	Profile *Profile
}


func (User) TableName() string {
	return "users"
}
```

#### What does it means ?

First of all, you have to define a `TableName` method that returns the database table name _(or view)_.
Without that information, `sqlxx` cannot uses that struct as a `Model`.

Then, you have to define your model columns using struct tags:

 * **column**(`string`): Define column name.
 * **pk**(`bool|string`): Define column as a primary key, it accepts the following argument:
   * **true**: Uses internal db mechanism to define primary key value
   * **db**: Uses internal db mechanism to define primary key value
   * **ulid**: Generate a [ULID](https://github.com/ulid/spec) to define primary key value
 * **default**(`bool`): On insert, if model has a zero value, it will use the db default value.
 * **pk**(`string`): Define column as a foreign key, reference table must be provided.
 * **-**(`bool`): Ignore this field.

> **NOTE:** Tags of type `bool` can be set as `key:true` or just `key` for implicit `true`.

> **NOTE:** Tags must be separated by a comma (`tagA, tagB, tagC`).

Keep in mind that a model **requires one primary key (and just one)**. It's a known limitation that only one primary
key can be specified and it can't be a composite key.

After that, you can define optional relationships _(or associations)_ that can be preloaded later.
The preload mechanism, which enables you to fetch relationships from database, support these types:

 * `Model`
 * `*Model`
 * `[]Model`
 * `[]*Model`
 * `*[]Model`
 * `*[]*Model`

#### Conventions

##### ID as Primary Key

By default, if `pk` tag is undefined, `sqlxx` will uses the field named `ID` as primary key with
this configuration: `pk:db`

```go
type User struct {
	ID   string `sqlxx:"column:id"`   // Field named ID will be used a primary key by default.
	Name string `sqlxx:"column:name"`
}
```

##### Snake Case Column Name

By default, if `column` tag is undefined, `sqlxx` will transform field name to lower snake case as column name.

```go
type User struct {
	ID   string `sqlxx:"pk"` // Column name is `id`
	Name string `sqlxx:""`   // Column name is `name`
}
```

##### CreatedAt tracking

For models having `CreatedAt` field, it will be set to current time when record is first created.

```go
type User struct {
	ID        string    `sqlxx:"column:id,pk"`
	Name      string    `sqlxx:"column:name"`
	CreatedAt time.Time `sqlxx:"column:created_at"`
}
```

You can override the default field name and/or column name by adding this method:

```go
func (User) CreatedKey() string {
	return "created"
}
```

##### UpdatedAt tracking

For models having `UpdatedAt` field, it will be set to current time when record are updated.

```go
type User struct {
	ID        string    `sqlxx:"column:id,pk"`
	Name      string    `sqlxx:"column:name"`
	UpdatedAt time.Time `sqlxx:"column:updated_at"`
}
```

You can override the default field name and/or column name by adding this method:

```go
func (User) UpdatedKey() string {
	return "updated"
}
```

##### DeletedAt tracking

For models having `DeletedAt` field, it will be set to current time when record are archived.

```go
type User struct {
	ID        string      `sqlxx:"column:id,pk"`
	Name      string      `sqlxx:"column:name"`
	DeletedAt pq.NullTime `sqlxx:"column:deleted_at"`
}
```

You can override the default field name and/or column name by adding this method:

```go
func (User) DeletedKey() string {
	return "deleted"
}
```

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

#### Notes

If you have to examine rows generated from unit test, you prevent the test suite to cleanup by using:

```
DB_KEEP=true scripts/test
```

Then, you can access the database with:

```
scripts/database --client
```

### Random

Because sometimes it's hard to think of a good test fixture, using generators can save your productivity.

This website was a great help to write unit test: http://www.fantasynamegenerators.com


## License

This is Free Software, released under the [`MIT License`][license-url].


## Contributing

* Ping us on twitter:
  * [@novln_](https://twitter.com/novln_)
  * [@oibafsellig](https://twitter.com/oibafsellig)
  * [@thoas](https://twitter.com/thoas)
* Fork the [project](https://github.com/ulule/loukoum)
* Fix [bugs](https://github.com/ulule/loukoum/issues)

**Don't hesitate ;)**

[godoc-url]: https://godoc.org/github.com/ulule/sqlxx
[godoc-img]: https://godoc.org/github.com/ulule/sqlxx?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
[sql-url]: https://golang.org/pkg/database/sql/
[sqlx-url]: https://github.com/jmoiron/sqlx
[circle-url]: https://circleci.com/gh/ulule/sqlxx/tree/master
[circle-img]: https://circleci.com/gh/ulule/sqlxx.svg?style=shield&circle-token=e53497efffde023bac7f2710bd12c5d0e71f5af4
