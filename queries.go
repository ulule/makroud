package sqlxx

import (
	sqlDriver "database/sql/driver"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ulule/sqlxx/reflekt"
)

// GetPrimaryKeys returns primary keys for the given interface.
func GetPrimaryKeys(out interface{}, name string) ([]interface{}, error) {
	value := reflekt.ReflectValue(InterfaceToModel(out))
	_, isNull := NullFieldTypes[value.FieldByName(name).Type()]

	pks, err := reflekt.GetFieldValues(out, name)
	if err != nil {
		return nil, err
	}

	var values []interface{}

	for i := range pks {
		if !isNull {
			if reflekt.IsZeroValue(pks[i]) {
				return nil, fmt.Errorf("Cannot perform query on zero value (%s=%v)", name, pks[i])
			}
			values = append(values, pks[i])
		} else {
			valuer := reflekt.Copy(pks[i]).(sqlDriver.Valuer)
			if v, err := valuer.Value(); err == nil && v != nil {
				values = append(values, v)
			}
		}
	}

	return values, nil
}

// SoftDelete soft deletes the model in the database
func SoftDelete(driver Driver, out interface{}, fieldName string) error {
	schema, err := InterfaceToSchema(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryField
	pkValue, err := reflekt.GetFieldValue(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if reflekt.IsZeroValue(pkValue) {
		return fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	field := schema.Fields[fieldName]

	now := time.Now()

	query := fmt.Sprintf("UPDATE %s SET %s = :%s WHERE %s = :%s",
		schema.TableName,
		field.ColumnName,
		field.ColumnName,
		pkField.ColumnName,
		pkField.ColumnName)

	m := map[string]interface{}{
		field.ColumnName:   now,
		pkField.ColumnName: pkValue,
	}

	_, err = driver.NamedExec(query, m)
	if err != nil {
		return err
	}

	return nil

}

// Delete deletes the model in the database
func Delete(driver Driver, out interface{}) error {
	schema, err := InterfaceToSchema(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryField
	pkValue, _ := reflekt.GetFieldValue(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if reflekt.IsZeroValue(pkValue) {
		return fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = :%s",
		schema.TableName,
		pkField.ColumnName,
		pkField.ColumnName)

	_, err = driver.NamedExec(query, out)
	if err != nil {
		return err
	}

	return nil
}

// Save saves the model and populate it to the database
func Save(driver Driver, out interface{}) error {
	schema, err := InterfaceToSchema(out)
	if err != nil {
		return err
	}

	var (
		columns        = []string{}
		ignoredColumns = []string{}
		values         = []string{}
	)

	for _, column := range schema.Fields {
		var (
			isIgnored    bool
			hasDefault   bool
			defaultValue string
		)

		if tag := column.Tags.Get(StructTagName); tag != nil {
			isIgnored = len(tag.Get(StructTagIgnored)) != 0
			defaultValue = tag.Get(StructTagDefault)
			hasDefault = len(defaultValue) != 0
		}

		if isIgnored || hasDefault {
			ignoredColumns = append(ignoredColumns, column.ColumnName)
		}

		if !isIgnored {
			columns = append(columns, column.ColumnName)

			if hasDefault {
				values = append(values, defaultValue)
			} else {
				values = append(values, fmt.Sprintf(":%s", column.ColumnName))
			}
		}
	}

	var query string

	pkField := schema.PrimaryField
	pkValue, _ := reflekt.GetFieldValue(out, pkField.Name)

	if reflekt.IsZeroValue(pkValue) {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			schema.TableName,
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
	} else {
		updates := []string{}

		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s = :%s",
			schema.TableName,
			strings.Join(updates, ", "),
			pkField.ColumnName,
			pkField.ColumnName)
	}

	if len(ignoredColumns) > 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(ignoredColumns, ", "))
	}

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.Get(out, out)
	if err != nil {
		return err
	}

	return nil
}

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, fields ...string) error {
	var (
		err error
	)

	schema, err := InterfaceToSchema(out)
	if err != nil {
		return err
	}

	type ChildRelation struct {
		field         string
		relationField string
		relation      Relation
	}

	var (
		relations      []Relation
		childRelations []ChildRelation
		relationPaths  = schema.RelationPaths()
	)

	for _, field := range fields {
		relation, ok := relationPaths[field]
		if !ok {
			return fmt.Errorf("%s is not a valid relation", field)
		}

		splits := strings.Split(field, ".")
		if len(splits) == 1 {
			relations = append(relations, relation)
		}

		if len(splits) == 2 {
			childRelations = append(childRelations, ChildRelation{
				field:         splits[0],
				relationField: strings.Join(splits[1:], "."),
				relation:      relation,
			})
		}
	}

	if err = preloadRelations(driver, out, relations); err != nil {
		return err
	}

	type Relationship struct {
		item                       interface{}
		itemValue                  reflect.Value
		itemPK                     int64
		itemChild                  interface{}
		itemChildFieldName         string
		itemChildPK                interface{}
		itemChildPKFieldName       string
		itemChildRelationPK        int64
		itemChildRelationFieldName string
	}

	for _, child := range childRelations {
		// Articles
		if reflekt.IsSlice(out) {
			var (
				slice               = reflect.ValueOf(out).Elem()
				relationships       []Relationship
				childrenRelationPKs []int64
			)

			for i := 0; i < slice.Len(); i++ {
				var (
					// Article
					itemValue            = slice.Index(i)
					item                 = itemValue.Interface()
					itemPKFieldName      = schema.PrimaryField.Name
					itemChildFieldName   = child.field
					itemChildPKFieldName = child.relation.RelatedFKField()
				)

				// Retrieve Article.ID
				v, err := reflekt.GetFieldValue(item, itemPKFieldName)
				if err != nil {
					return err
				}

				itemPK, _ := IntToInt64(v)

				// Retrieve Article.User previously fetched
				itemChild, err := reflekt.GetFieldValue(item, itemChildFieldName)
				if err != nil {
					return err
				}

				// Retrieve Article.UserID
				v, err = reflekt.GetFieldValue(itemChild, child.relation.ParentSchema.PrimaryField.Name)
				if err != nil {
					return err
				}

				itemChildPK, _ := IntToInt64(v)

				// Retrieve Article.User.APIKeyID (for SELECT IN)
				v, err = reflekt.GetFieldValue(itemChild, child.relation.RelatedFKField())
				if err != nil {
					return err
				}

				if valuer, ok := v.(sqlDriver.Valuer); ok {
					if valuerValuer, err := valuer.Value(); err == nil {
						v = valuerValuer
					}
				}

				itemChildRelationPK, _ := IntToInt64(v)

				relationships = append(relationships, Relationship{
					item:                       item,
					itemValue:                  itemValue,
					itemPK:                     itemPK,
					itemChild:                  itemChild,
					itemChildFieldName:         itemChildFieldName,
					itemChildPK:                itemChildPK,
					itemChildPKFieldName:       itemChildPKFieldName,
					itemChildRelationPK:        itemChildRelationPK,
					itemChildRelationFieldName: child.relationField,
				})

				var exists bool
				for _, v := range childrenRelationPKs {
					if v == itemChildRelationPK {
						exists = true
						break
					}
				}

				if !exists {
					childrenRelationPKs = append(childrenRelationPKs, itemChildRelationPK)
				}

			}

			// Build a []APIKey slice
			t := reflect.SliceOf(reflekt.ReflectType(child.relation.Model))
			s := reflect.New(t)
			s.Elem().Set(reflect.MakeSlice(t, 0, 0))

			// SELECT * FROM from api_keys WHERE id IN childrenRelationPKs
			if err = FindByParams(driver, s.Interface(), map[string]interface{}{child.relation.Schema.PrimaryField.ColumnName: childrenRelationPKs}); err != nil {
				return err
			}

			instances := s.Elem()

			// Set the relations
			for _, relationship := range relationships {
				for i := 0; i < instances.Len(); i++ {
					instance := instances.Index(i)

					// APIKey.ID
					v, err := reflekt.GetFieldValue(instance, child.relation.Schema.PrimaryField.Name)
					if err != nil {
						return err
					}

					instancePK, _ := IntToInt64(v)

					if relationship.itemChildRelationPK == instancePK {
						itemChildCopy := reflekt.Copy(relationship.itemChild)
						if err = reflekt.SetFieldValue(itemChildCopy, relationship.itemChildRelationFieldName, instance.Interface()); err != nil {
							return err
						}
						if err = reflekt.SetFieldValue(relationship.itemValue, relationship.itemChildFieldName, itemChildCopy); err != nil {
							return err
						}
					}
				}
			}

			continue
		}

		instance, err := reflekt.GetFieldValue(out, child.field)
		if err != nil {
			return err
		}

		cp := reflekt.Copy(instance)

		if err = Preload(driver, cp, child.relationField); err != nil {
			return err
		}

		if err = reflekt.SetFieldValue(out, child.field, cp); err != nil {
			return err
		}
	}

	return nil
}

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, true)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, false)
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(model Model, params map[string]interface{}, fetchOne bool) (string, []interface{}, error) {
	schema, err := GetSchema(model)
	if err != nil {
		return "", nil, err
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		schema.ColumnPaths(),
		model.TableName(),
		schema.WhereColumnPaths(params))

	if fetchOne {
		q = fmt.Sprintf("%s LIMIT 1", q)
	}

	query, args, err := sqlx.Named(q, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

// where executes a where clause.
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) error {
	model := InterfaceToModel(out)

	query, args, err := whereQuery(model, params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(out, driver.Rebind(query), args...)
	}

	return driver.Select(out, driver.Rebind(query), args...)
}
