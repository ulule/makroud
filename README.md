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

Model are struct that contains basic go types, pointers, `sql.Scanner`, `driver.Valuer` or `Model` interface.
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

### Operations

For the following sections, we assume that you have a `context.Context` and a `sqlxx.Driver` instance.

#### Insert

For a simple insert, you can use save a model like this:

```go
func CreateUser(ctx context.Context, driver sqlxx.Driver, name string) (*User, error) {
	user := &User{
		Name: name,
	}

	err := sqlxx.Save(ctx, driver, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `InsertBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func CreateUser(ctx context.Context, driver sqlxx.Driver, name string) (*User, error) {
	user := &User{
		Name: name,
	}

	stmt := loukoum.Insert("users").
		Set(loukoum.Pair("name", user.Name)).
		Returning("id, created_at, updated_at")

	err := sqlxx.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Update

For a simple update, asumming your model have a primary key defined, you can use save it by executing:

```go
func UpdateUser(ctx context.Context, driver sqlxx.Driver, user *User, name string) error {
	user.Name = name
	return sqlxx.Save(ctx, driver, user)
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func UpdateUser(ctx context.Context, driver sqlxx.Driver, user *User, name string) error {
	user.Name = name

	stmt := loukoum.Update("users").
		Set(
			loukoum.Pair("updated_at", loukoum.Raw("NOW()")),
			loukoum.Pair("name", user.Name),
		).
		Where(loukoum.Condition("id").Equal(user.ID)).
		Returning("updated_at")

	err := sqlxx.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Delete

For a simple delete _(using a `DELETE` statement)_, asumming your model have a primary key defined,
you can use delete it using:

```go
func DeleteUser(ctx context.Context, driver sqlxx.Driver, user *User) error {
	return sqlxx.Delete(ctx, driver, user)
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `DeleteBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func DeleteUser(ctx context.Context, driver sqlxx.Driver, user *User) error {
	stmt := loukoum.Delete("users").Where(loukoum.Condition("id").Equal(user.ID))

	err := sqlxx.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Archive

Archive executes an `UPDATE` on `DeletedAt` field on given value.

```go
func ArchiveUser(ctx context.Context, driver sqlxx.Driver, user *User) error {
	return sqlxx.Archive(ctx, driver, user)
}
```

> **NOTE**: If model has no `DeletedAt` field, an error is returned.

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func ArchiveUser(ctx context.Context, driver sqlxx.Driver, user *User) error {
	user.Name = name

	stmt := loukoum.Update("users").
		Set(
			loukoum.Pair("deleted_at", loukoum.Raw("NOW()")),
			loukoum.Pair("name", ""),
		).
		Where(loukoum.Condition("id").Equal(user.ID)).
		Returning("deleted_at")

	err := sqlxx.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Query

TODO

### Preload

TODO

## Development

### Docker

The test suite is running on PostgreSQL. We use [Docker](https://docs.docker.com/install/) to create a running
container using [scripts/database](scripts/database).

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
