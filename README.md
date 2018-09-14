# Makroud

[![CircleCI][circle-img]][circle-url]
[![Documentation][godoc-img]][godoc-url]
![License][license-img]

*A high level SQL Connector.*

## Introduction

Makroud is a high level SQL Connector that only support **PostgreSQL** at the moment.

It's an advanced mapper and/or a lightweight ORM that relies on reflection to generate queries.
**Using reflection has it's flaws**, type safety is not guaranteed and a panic is still possible,
even if you are every careful and vigilant. **However,** development is super easy and straightforward
since it doesn't relies on code generation.

Makroud isn't a migration tools, and doesn't inspects the database to define the application data model
since there is no code generation. It's important to have Active Record that are synchronized with your
data model in your database.

It also support simple associations with preloading.

Under the hood, it relies on three components:

 * [Loukoum](https://github.com/ulule/loukoum) for query generation
 * [Sqlx](https://github.com/ulule/sqlx) to have an extended mapper
 * [Sqalx](https://github.com/ulule/sqalx) to support nested transaction

## Installation

Using [dep](https://github.com/golang/dep)

```console
dep ensure -add github.com/ulule/makroud@master
```

or `go get`

```console
go get -u github.com/ulule/makroud
```

## Usage

### Create a Driver

A **Driver** is a high level abstraction of a database connection or a transaction.
It's almost required everytime alongside a `context.Context` to manipulate rows.

```go
driver, err := makroud.New(
	makroud.Host(cfg.Host),
	makroud.Port(cfg.Port),
	makroud.User(cfg.User),
	makroud.Password(cfg.Password),
	makroud.Database(cfg.Name),
	makroud.SSLMode(cfg.SSLMode),
	makroud.MaxOpenConnections(cfg.MaxOpenConnections),
	makroud.MaxIdleConnections(cfg.MaxIdleConnections),
)
```

Also, you can use directly a struct if you don't need to use
[functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis):

```go
driver, err := makroud.NewWithOptions(&makroud.ClientOptions{
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
	ID        string `makroud:"column:id,pk:ulid"`
	Email     string `makroud:"column:email"`
	Password  string `makroud:"column:password"`
	Country   string `makroud:"column:country"`
	Locale    string `makroud:"column:locale"`
	ProfileID string `makroud:"column:profile_id,fk:profiles"`
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
Without that information, `makroud` cannot uses that struct as a `Model`.

Then, you have to define your model columns using struct tags:

 * **column**(`string`): Define column name.
 * **pk**(`bool|string`): Define column as a primary key, it accepts the following argument:
   * **true**: Uses internal db mechanism to define primary key value
   * **ulid**: Generate a [ULID](https://github.com/ulid/spec) to define primary key value
   * **uuid-v1**: Generate a [UUID V1](https://en.wikipedia.org/wiki/Universally_unique_identifier)
     to define primary key value
   * **uuid-v4**: Generate a [UUID V4](https://en.wikipedia.org/wiki/Universally_unique_identifier)
     to define primary key value
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

By default, if `pk` tag is undefined, `makroud` will uses the field named `ID` as primary key with
this configuration: `pk:db`

```go
type User struct {
	ID   string `makroud:"column:id"`   // Field named ID will be used a primary key by default.
	Name string `makroud:"column:name"`
}
```

##### Snake Case Column Name

By default, if `column` tag is undefined, `makroud` will transform field name to lower snake case as column name.

```go
type User struct {
	ID   string `makroud:"pk"` // Column name is `id`
	Name string `makroud:""`   // Column name is `name`
}
```

##### CreatedAt tracking

For models having `CreatedAt` field, it will be set to current time when record is first created.

```go
type User struct {
	ID        string    `makroud:"column:id,pk"`
	Name      string    `makroud:"column:name"`
	CreatedAt time.Time `makroud:"column:created_at"`
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
	ID        string    `makroud:"column:id,pk"`
	Name      string    `makroud:"column:name"`
	UpdatedAt time.Time `makroud:"column:updated_at"`
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
	ID        string      `makroud:"column:id,pk"`
	Name      string      `makroud:"column:name"`
	DeletedAt pq.NullTime `makroud:"column:deleted_at"`
}
```

You can override the default field name and/or column name by adding this method:

```go
func (User) DeletedKey() string {
	return "deleted"
}
```

### Operations

For the following sections, we assume that you have a `context.Context` and a `makroud.Driver` instance.

#### Insert

For a simple insert, you can use save a model like this:

```go
func CreateUser(ctx context.Context, driver makroud.Driver, name string) (*User, error) {
	user := &User{
		Name: name,
	}

	err := makroud.Save(ctx, driver, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `InsertBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func CreateUser(ctx context.Context, driver makroud.Driver, name string) (*User, error) {
	user := &User{
		Name: name,
	}

	stmt := loukoum.Insert("users").
		Set(loukoum.Pair("name", user.Name)).
		Returning("id, created_at, updated_at")

	err := makroud.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Update

For a simple update, asumming your model have a primary key defined, you can use save it by executing:

```go
func UpdateUser(ctx context.Context, driver makroud.Driver, user *User, name string) error {
	user.Name = name
	return makroud.Save(ctx, driver, user)
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func UpdateUser(ctx context.Context, driver makroud.Driver, user *User, name string) error {
	user.Name = name

	stmt := loukoum.Update("users").
		Set(
			loukoum.Pair("updated_at", loukoum.Raw("NOW()")),
			loukoum.Pair("name", user.Name),
		).
		Where(loukoum.Condition("id").Equal(user.ID)).
		Returning("updated_at")

	err := makroud.Exec(ctx, driver, stmt, user)
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
func DeleteUser(ctx context.Context, driver makroud.Driver, user *User) error {
	return makroud.Delete(ctx, driver, user)
}
```

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `DeleteBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func DeleteUser(ctx context.Context, driver makroud.Driver, user *User) error {
	stmt := loukoum.Delete("users").Where(loukoum.Condition("id").Equal(user.ID))

	err := makroud.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Archive

Archive executes an `UPDATE` on `DeletedAt` field on given value.

```go
func ArchiveUser(ctx context.Context, driver makroud.Driver, user *User) error {
	return makroud.Archive(ctx, driver, user)
}
```

> **NOTE**: If model has no `DeletedAt` field, an error is returned.

Or for more complex statements, uses a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum"

func ArchiveUser(ctx context.Context, driver makroud.Driver, user *User) error {
	user.Name = name

	stmt := loukoum.Update("users").
		Set(
			loukoum.Pair("deleted_at", loukoum.Raw("NOW()")),
			loukoum.Pair("name", ""),
		).
		Where(loukoum.Condition("id").Equal(user.ID)).
		Returning("deleted_at")

	err := makroud.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

#### Query

By using a [Loukoum](https://github.com/ulule/loukoum) `SelectBuilder`.

```go
import "github.com/ulule/loukoum"

func GetUserByID(ctx context.Context, driver makroud.Driver, id string) (*User, error) {
	user := &User{}

	columns, err := makroud.GetColumns(driver, user)
	if err != nil {
		return nil, err
	}

	stmt := loukoum.Select(columns...).
		From(user.TableName()).
		Where(loukoum.Condition("id").Equal(id))

	err := makroud.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByName(ctx context.Context, driver makroud.Driver, name string) (*User, error) {
	user := &User{}

	columns, err := makroud.GetColumns(driver, user)
	if err != nil {
		return nil, err
	}

	stmt := loukoum.Select(columns...).
		From(user.TableName()).
		Where(loukoum.Condition("name").Equal(name))

	err := makroud.Exec(ctx, driver, stmt, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
```

Also, it's support query without `Model`.

```go
func FindUserIDWithStaffRole(ctx context.Context, driver makroud.Driver) ([]string, error) {
	list := []string{}

	stmt := loukoum.Select("id").
		From("users").
		Where(loukoum.Condition("role").Equal("staff"))

	err := makroud.Exec(ctx, driver, stmt, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
```

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

[godoc-url]: https://godoc.org/github.com/ulule/makroud
[godoc-img]: https://godoc.org/github.com/ulule/makroud?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
[circle-url]: https://circleci.com/gh/ulule/makroud/tree/master
[circle-img]: https://circleci.com/gh/ulule/makroud.svg?style=shield&circle-token=e53497efffde023bac7f2710bd12c5d0e71f5af4
