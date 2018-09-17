// Package makroud provides a high level SQL Connector.
// At the moment, only PostgreSQL is supported.
//
// It's an advanced mapper and/or a lightweight ORM that relies on reflection.
//
// For further informations, you can read this documentation:
// https://github.com/ulule/makroud/blob/master/README.md
//
// Or you can discover makroud with these examples.
// First, you have to create a driver:
//
//   driver, err := makroud.New(
//       makroud.Host(cfg.Host),
//       makroud.Port(cfg.Port),
//       makroud.User(cfg.User),
//       makroud.Password(cfg.Password),
//       makroud.Database(cfg.Name),
//       makroud.SSLMode(cfg.SSLMode),
//       makroud.MaxOpenConnections(cfg.MaxOpenConnections),
//       makroud.MaxIdleConnections(cfg.MaxIdleConnections),
//   )
//
// Then, define a model:
//
//   type User struct {
//       ID        string `makroud:"column:id,pk:ulid"`
//       Email     string `makroud:"column:email"`
//       Password  string `makroud:"column:password"`
//       Country   string `makroud:"column:country"`
//       Locale    string `makroud:"column:locale"`
//   }
//
// Execute an insert:
//
//   user := &User{
//       Email:    "gilles@ulule.com",
//       Password: "019a7bdf56b9f48e18096d62b21f",
//       Country:  "FR",
//       Locale:   "fr",
//   }
//
//   err := makroud.Save(ctx, driver, user)
//
// Or an update:
//
//   user.Email = "gilles.fabio@ulule.com"
//
//   err := makroud.Save(ctx, driver, user)
//
// Or execute a simple query without model:
//
//   import "github.com/ulule/loukoum"
//
//   list := []string{}
//
//   stmt := loukoum.Update("users").
//       Set(
//           loukoum.Pair("updated_at", loukoum.Raw("NOW()")),
//           loukoum.Pair("status", status),
//       ).
//       Where(loukoum.Condition("group_id").Equal(gid)).
//       Returning("id")
//
//   err := makroud.Exec(ctx, driver, stmt, &list)
//
package makroud
