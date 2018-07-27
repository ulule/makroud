package sqlxx

import (
	"fmt"

	"github.com/pkg/errors"
)

// ForeignKey is a composite object that define a foreign key for a model.
// This foreign key will be used later for Preload...
//
// For example: If we have an User, we could have this primary key defined in User's schema.
//
//     ForeignKey {
//         ModelName: User,
//         TableName: users,
//         FieldName: AvatarID,
//         ColumnName: avatar_id,
//         ColumnPath: users.avatar_id,
//         Reference: avatars,
//     }
//
type ForeignKey struct {
	Field
	fkTableName string
}

// NewForeignKey creates a foreign key from a field instance.
func NewForeignKey(field *Field) (*ForeignKey, error) {
	pk := &ForeignKey{
		Field:       *field,
		fkTableName: field.GetForeignKey(),
	}
	return pk, nil
}

// Reference returns the foreign key's table name.
func (key ForeignKey) Reference() string {
	return key.fkTableName
}

type Reference struct {
	Field
	fkModelName  string
	fkTableName  string
	fkFieldName  string
	fkColumnPath string
	fkColumnName string
}

// NewReference creates a reference from a field instance.
func NewReference(driver Driver, local *Schema, field *Field) (*Reference, error) {
	reference := ToModel(field.rtype)
	if reference == nil {
		return nil, errors.Errorf("invalid model: %s", field.rtype.String())
	}

	// TODO (novln): Allow circular reference...
	remote, err := GetSchema(driver, reference)
	if err != nil {
		return nil, err
	}

	// Article.Author(User)
	if field.associationType == AssociationTypeOne {
		fmt.Printf("::7 has_once\n")

		for i := range local.references {
			fmt.Println("::8", local.references[i])
			if local.references[i].TableName() == remote.TableName() {
				fmt.Println("::10", local.references[i])
			}
		}

		for i := range remote.references {
			fmt.Println("::9", remote.references[i])
			if remote.references[i].TableName() == remote.TableName() {
				fmt.Println("::11", remote.references[i])
			}
		}

		// local.
		//
		// fk := &Reference{
		// 	Field:        *field,
		// 	fkModelName:  schema.modelName,
		// 	fkTableName:  schema.tableName,
		// 	fkFieldName:  schema.pk.fieldName,
		// 	fkColumnPath: schema.pk.columnPath,
		// 	fkColumnName: schema.pk.columnName,
		// }
		// fmt.Printf("::8 %s %s %s %s %s\n", fk.fkModelName, fk.fkTableName, fk.fkFieldName, fk.fkColumnName, fk.fkColumnPath)
		// return fk, nil
	}

	// // User.Avatars(Avatar) -- Avatar.UserID
	// if field.AssociationType == AssociationTypeMany {
	// 	fmt.Printf("::7 has_many\n")
	// 	fk := &ForeignKey{
	// 		Model:                referenceModel,                                                                                   // Avatar model
	// 		ModelName:            referenceModelName,                                                                               // Avatar
	// 		TableName:            referenceTableName,                                                                               // avatars
	// 		FieldName:            fmt.Sprintf("%s%s", field.ModelName, PrimaryKeyFieldName),                                        // UserID
	// 		ColumnName:           fmt.Sprintf("%s_%s", snaker.CamelToSnake(field.ModelName), strings.ToLower(PrimaryKeyFieldName)), // user_id
	// 		AssociationFieldName: field.ModelName,                                                                                  // User
	//
	// 		Reference: &ForeignKey{
	// 			Model:                field.Model,                          // User model
	// 			ModelName:            field.ModelName,                      // User
	// 			TableName:            field.TableName,                      // users
	// 			FieldName:            PrimaryKeyFieldName,                  // ID
	// 			ColumnName:           strings.ToLower(PrimaryKeyFieldName), // id
	// 			AssociationFieldName: field.FieldName,                      // Avatars
	// 		},
	// 	}
	// 	spew.Dump(fk)
	// 	return fk, nil
	// }
	// fmt.Printf("::7 has_none ?\n")
	panic("Fuck")
	return nil, nil
}
