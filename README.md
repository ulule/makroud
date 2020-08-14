# Makroud

[![CircleCI][circle-img]][circle-url]
[![Documentation][godoc-img]][godoc-url]
![License][license-img]

_A high level SQL Connector._

## Introduction

Makroud is a high level SQL Connector that only support **PostgreSQL** at the moment.

It's an advanced mapper and/or a lightweight ORM that relies on reflection to generate queries.
**Using reflection has its flaws**, type safety is not guaranteed and a panic is still possible,
even if you are very careful and vigilant. **However,** development is super easy and straightforward
since it doesn't rely on code generation.

Makroud isn't a migration tool and it doesn't inspect the database to define the application data model
_(since there is no code generation)_. It's **really** important to have Active Record that are synchronized with your
data model in your database.

It also support simple associations _(one, many)_ with preloading.

Under the hood, it relies on [Loukoum](https://github.com/ulule/loukoum) for query generation.
In addition, it's heavily inspired by [Sqlx](https://github.com/jmoiron/sqlx) for its extended mapper and
[Sqalx](https://github.com/heetch/sqalx) to support nested transaction.

## Installation

Using [Go Modules](https://github.com/golang/go/wiki/Modules)

```console
go get github.com/ulule/makroud@v0.7.1
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

Also, you can use a struct directly if you don't need to use
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

### Advanced mapper

If a lightweight ORM doesn't fit your requirements and an advanced mapper is enough for your usecase:
makroud does that very well.

You could either use primitive and compound types for your queries:

```go
import "github.com/ulule/makroud"

stmt := `SELECT id FROM users WHERE email = 'john.doe@example.com'`
id := int64(0)

err := makroud.RawExec(ctx, driver, stmt, &id)
if err != nil {
	return err
}

stmt = `SELECT created_at FROM users WHERE id = 42`
created := time.Time{}

err := makroud.RawExec(ctx, driver, stmt, &created)
if err != nil {
	return err
}

stmt = `SELECT email FROM users WHERE id IN (1, 2, 3, 4)`
list := []string{}

err := makroud.RawExec(ctx, driver, stmt, &list)
if err != nil {
	return err
}
```

Or, define a struct that contains your database table _(or view)_ columns:

```go
import "github.com/ulule/makroud"
import "github.com/ulule/loukoum/v3"

type User struct {
	ID       int64  `mk:"id"`
	Email    string `mk:"email"`
	Password string `mk:"password"`
	Country  string `mk:"country"`
	Locale   string `mk:"locale"`
}

users := []User{}
stmt := loukoum.Select("*").
	From("users").
	Where(loukoum.Condition("id").In(1, 2, 3, 4))

err := makroud.Exec(ctx, driver, stmt, &users)
if err != nil {
	return err
}
```

### Lightweight ORM

#### Define a Model

With the [**Active Record**](https://en.wikipedia.org/wiki/Active_record_pattern) approach, you have to define a
model that wraps your database table _(or view)_ columns into properties.

Models are structs that contain basic go types, pointers, `sql.Scanner`, `driver.Valuer` or `Model` interface.
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
	Roles    []Role    `makroud:"relation:roles.user_id"`
	Profile  *Profile  `makroud:"relation:profile_id"`
}

func (User) TableName() string {
	return "users"
}
```

#### What does it means ?

First of all, you have to define a `TableName` method that returns the database table name _(or view)_.
Without that information, `makroud` cannot uses that struct as a `Model`.

Then, you have to define your model columns using struct tags:

- **column**(`string`): Define column name.
- **pk**(`bool|string`): Define column as a primary key, it accepts the following argument:
  - **true**: Uses internal db mechanism to define primary key value
  - **ulid**: Generate a [ULID](https://github.com/ulid/spec) to define primary key value
  - **uuid-v1**: Generate a [UUID V1](<https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_1_(date-time_and_MAC_address)>)
    to define primary key value
  - **uuid-v4**: Generate a [UUID V4](<https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)>)
    to define primary key value
- **default**(`bool`): On insert, if model has a zero value, it will use the db default value.
- **fk**(`string`): Define column as a foreign key, reference table must be provided.
- **relation**(`string`): Define which column to use for preload. The column must be prefixed by the table name
  if it's not the model table name _(However, the prefix is optional if the table name is the same as the model)_.
  See [Preload](https://github.com/ulule/makroud#preload) section for further information.
- **-**(`bool`): Ignore this field.

> **NOTE:** Tags of type `bool` can be set as `key:true` or just `key` for implicit `true`.

> **NOTE:** Tags must be separated by a comma (`tagA, tagB, tagC`).

Keep in mind that a model **requires one primary key (and just one)**. It's a known limitation that only one primary
key can be specified and it can't be a composite key.

After that, you can define optional relationships _(or associations)_ that can be preloaded later.
The preload mechanism, which enables you to fetch relationships from database, support these types:

- `Model`
- `*Model`
- `[]Model`
- `[]*Model`
- `*[]Model`
- `*[]*Model`

> **NOTE:** You could either use `makroud` or `mk` as tag identifier.

#### Conventions

##### ID as Primary Key

By default, if the `pk` tag is undefined, `makroud` will use the field named `ID` as primary key with
this configuration: `pk:db`

```go
type User struct {
	ID   string `makroud:"column:id"`   // Field named ID will be used as a primary key by default.
	Name string `makroud:"column:name"`
}
```

##### Snake Case Column Name

By default, if the `column` tag is undefined, `makroud` will transform field name to lower snake case as column name.

```go
type User struct {
	ID   string `makroud:"pk"` // Column name is `id`
	Name string `makroud:""`   // Column name is `name`
}
```

##### Preload relationships

By default, if the `relation` tag is undefined, `makroud` will infer the column name to use for the preload mechanism.

**Local foreign key:**

Let's define a user with a profile:

```go
type User struct {
	ID       string  `makroud:"column:id,pk"`
	Email    string  `makroud:"column:email"`
	PID      string  `makroud:"column:profile_id,fk:profiles"`
	Profile  *Profile
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	ID         string  `makroud:"column:id,pk:ulid"`
	FirstName  string  `makroud:"column:first_name"`
	LastName   string  `makroud:"column:last_name"`
	Enabled    bool    `makroud:"column:enabled"`
}

func (Profile) TableName() string {
	return "profiles"
}
```

Since the field `Profile` in the `User` has no `relation` tag, `makroud` will try to find, in the first pass,
the field with the name `ProfileID` _(FieldName + ID)_ in `User` model.
It's mandatory that this field is also a foreign key to the `profiles` table.

Unfortunately for us, `User` model has no such field. So, `makroud` will try to find, in the second and final pass,
the first field that is a foreign key to the `profiles` table. In our example, it will use the field `PID`.

**Remote foreign key:**

Let's define a user with a profile:

```go
type User struct {
	ID       string  `makroud:"column:id,pk"`
	Email    string  `makroud:"column:email"`
	Profile  *Profile
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	ID         string  `makroud:"column:id,pk:ulid"`
	FirstName  string  `makroud:"column:first_name"`
	LastName   string  `makroud:"column:last_name"`
	Enabled    bool    `makroud:"column:enabled"`
	UID        string  `makroud:"column:user_id,fk:users"`
}

func (Profile) TableName() string {
	return "profiles"
}
```

Since the field `Profile` in the `User` has no `relation` tag, `makroud` will try to find, in the first pass,
the field with the name `UserID` _(ModelName + ID)_ in `Profile` model.
It's mandatory that this field is also a foreign key to the `users` table.

Unfortunately for us, `Profile` model has no such field. So, `makroud` will try to find, in the second and final pass,
the first field that is a foreign key to the `users` table. In our example, it will use the field `UID`.

##### CreatedAt tracking

For models having a `CreatedAt` field, it will be set to current time when the record is first created.

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

For models having a `UpdatedAt` field, it will be set to current time when records are updated.

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

For models having a `DeletedAt` field, it will be set to current time when records are archived.

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

Or for more complex statements, use a [Loukoum](https://github.com/ulule/loukoum) `InsertBuilder` alongside the model.

```go
import "github.com/ulule/loukoum/v3"

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

For a simple update, asumming your model have a primary key defined, you can save it by executing:

```go
func UpdateUser(ctx context.Context, driver makroud.Driver, user *User, name string) error {
	user.Name = name
	return makroud.Save(ctx, driver, user)
}
```

Or for more complex statements, use a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum/v3"

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
you can delete it using:

```go
func DeleteUser(ctx context.Context, driver makroud.Driver, user *User) error {
	return makroud.Delete(ctx, driver, user)
}
```

Or for more complex statements, use a [Loukoum](https://github.com/ulule/loukoum) `DeleteBuilder` alongside the model.

```go
import "github.com/ulule/loukoum/v3"

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

> **NOTE**: If the model has no `DeletedAt` field, an error is returned.

Or for more complex statements, use a [Loukoum](https://github.com/ulule/loukoum) `UpdateBuilder` alongside the model.

```go
import "github.com/ulule/loukoum/v3"

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
import "github.com/ulule/loukoum/v3"

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

Also, it supports query without `Model`.

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

On models having associations, you can execute a preload to fetch these relationships from the database.

Let's define a user with a profile:

```go
type User struct {
	ID       string   `makroud:"column:id,pk"`
	Email    string   `makroud:"column:email"`
	Profile  *Profile `makroud:"relation:profiles.user_id"`
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	ID         string  `makroud:"column:id,pk:ulid"`
	FirstName  string  `makroud:"column:first_name"`
	LastName   string  `makroud:"column:last_name"`
	UserID     string  `makroud:"column:user_id,fk:users"`
	Enabled    bool    `makroud:"column:enabled"`
}

func (Profile) TableName() string {
	return "profiles"
}
```

Once you obtain a user record, you can preload its profile by executing:

```go
err := makroud.Preload(ctx, driver, &user, makroud.WithPreloadField("Profile"))
```

**Or,** if preloading requires specific conditions, you can use a callback like this:

```go
import "github.com/ulule/loukoum/v3/builder"

err := makroud.Preload(ctx, driver, &user,
	makroud.WithPreloadCallback("Profile", func(query builder.Select) builder.Select {
		return query.Where(loukoum.Condition("enabled").Equal(true))
	}),
)
```

If there is no error and if the user record has a profile, then you should have the `Profile` value loaded.

<!---

## Benchmarks

A [benchmark repository](https://github.com/ulule/makroud-benchmarks) containing result with
other ORM or Mapper using reflection is available.

> **NOTE:** A benchmark is always an observation, not a measurement of performance.

![SelectAll NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_all_nsop.png)
![SelectAll Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_all_bop.png)
![SelectAll Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_all_aop.png)
![SelectSubset NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_subset_nsop.png)
![SelectSubset Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_subset_bop.png)
![SelectSubset Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_subset_aop.png)
![SelectComplex NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_complex_nsop.png)
![SelectComplex Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_complex_bop.png)
![SelectComplex Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/select_complex_aop.png)
![Insert NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/insert_nsop.png)
![Insert Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/insert_bop.png)
![Insert Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/insert_aop.png)
![Update NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/update_nsop.png)
![Update Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/update_bop.png)
![Update Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/update_aop.png)
![Delete NsOp](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/delete_nsop.png)
![Delete Bop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/delete_bop.png)
![Delete Aop](https://raw.githubusercontent.com/ulule/makroud-benchmarks/master/graph/images/delete_aop.png)

-->

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

If you have to examine rows generated from unit test, you can prevent the test suite to cleanup by using:

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

- Ping us on twitter:
  - [@novln\_](https://twitter.com/novln_)
  - [@oibafsellig](https://twitter.com/oibafsellig)
  - [@thoas](https://twitter.com/thoas)
- Fork the [project](https://github.com/ulule/loukoum)
- Fix [bugs](https://github.com/ulule/loukoum/issues)

**Don't hesitate ;)**

[godoc-url]: https://godoc.org/github.com/ulule/makroud
[godoc-img]: https://godoc.org/github.com/ulule/makroud?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
[circle-url]: https://circleci.com/gh/ulule/makroud/tree/master
[circle-img]: https://circleci.com/gh/ulule/makroud.svg?style=shield&circle-token=e53497efffde023bac7f2710bd12c5d0e71f5af4
