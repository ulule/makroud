package makroud

import (
	"bytes"
	"strconv"
)

// ----------------------------------------------------------------------------
// Debug high-level components
// ----------------------------------------------------------------------------

// DebugTagProperty returns a human readable version of given instance.
func DebugTagProperty(prop TagProperty) string {
	buffer := &bytes.Buffer{}
	debugTagProperty(prop).write(buffer)
	return buffer.String()
}

// DebugTag returns a human readable version of given instance.
func DebugTag(tag Tag) string {
	buffer := &bytes.Buffer{}
	debugTag(tag).write(buffer)
	return buffer.String()
}

// DebugTags returns a human readable version of given instance.
func DebugTags(tags Tags) string {
	buffer := &bytes.Buffer{}
	debugTags(tags).write(buffer)
	return buffer.String()
}

// DebugField returns a human readable version of given instance.
func DebugField(field Field) string {
	buffer := &bytes.Buffer{}
	debugField(field).write(buffer)
	return buffer.String()
}

// ----------------------------------------------------------------------------
// Debug low-level components
// ----------------------------------------------------------------------------

type debugWriter interface {
	write(buffer *bytes.Buffer)
}

type debugValue struct {
	k string
	v string
}

func (val debugValue) write(buffer *bytes.Buffer) {
	buffer.WriteString(`"`)
	buffer.WriteString(val.k)
	buffer.WriteString(`":"`)
	buffer.WriteString(val.v)
	buffer.WriteString(`"`)
}

type debugValues struct {
	k string
	v []debugWriter
}

func (val debugValues) write(buffer *bytes.Buffer) {
	buffer.WriteString(`"`)
	buffer.WriteString(val.k)
	buffer.WriteString(`":{`)
	for i := range val.v {
		if i != 0 {
			buffer.WriteString(",")
		}
		val.v[i].write(buffer)
	}
	buffer.WriteString(`}`)
}

type debugWrap struct {
	k string
	v debugWriter
}

func (val debugWrap) write(buffer *bytes.Buffer) {
	buffer.WriteString(`"`)
	buffer.WriteString(val.k)
	buffer.WriteString(`": `)
	val.v.write(buffer)
}

type debugObj []debugWriter

func (val debugObj) write(buffer *bytes.Buffer) {
	buffer.WriteString("{")
	for i := range val {
		if i != 0 {
			buffer.WriteString(",")
		}
		val[i].write(buffer)
	}
	buffer.WriteString("}")
}

type debugArr []debugWriter

func (val debugArr) write(buffer *bytes.Buffer) {
	buffer.WriteString("[")
	for i := range val {
		if i != 0 {
			buffer.WriteString(", ")
		}
		val[i].write(buffer)
	}
	buffer.WriteString(`]`)
}

func debugTagProperty(prop TagProperty) debugWriter {
	return debugObj{
		debugValue{
			k: "key",
			v: prop.Key(),
		},
		debugValue{
			k: "value",
			v: prop.Value(),
		},
	}
}

func debugTag(tag Tag) debugWriter {
	props := make([]debugWriter, 0, len(tag.properties))
	for i := range tag.properties {
		props[i] = debugTagProperty(tag.properties[i])
	}
	return debugObj{
		debugValue{
			k: "name",
			v: tag.Name(),
		},
		debugValues{
			k: "properties",
			v: props,
		},
	}
}

func debugTags(tags Tags) debugWriter {
	props := make([]debugWriter, 0, len(tags))
	for i := range tags {
		props[i] = debugTag(tags[i])
	}
	return debugArr(props)
}

func debugField(field Field) debugWriter {
	return debugObj{
		debugValue{
			k: "model_name",
			v: field.ModelName(),
		},
		debugValue{
			k: "table_name",
			v: field.TableName(),
		},
		debugValue{
			k: "field_name",
			v: field.FieldName(),
		},
		debugValue{
			k: "column_path",
			v: field.ColumnPath(),
		},
		debugValue{
			k: "column_name",
			v: field.ColumnName(),
		},
		debugValue{
			k: "is_primary_key",
			v: strconv.FormatBool(field.IsPrimaryKey()),
		},
		debugValue{
			k: "is_foreign_key",
			v: strconv.FormatBool(field.IsForeignKey()),
		},
		debugValue{
			k: "foreign_key",
			v: field.ForeignKey(),
		},
		debugValue{
			k: "is_association",
			v: strconv.FormatBool(field.IsAssociation()),
		},
		debugValue{
			k: "is_excluded",
			v: strconv.FormatBool(field.IsExcluded()),
		},
		debugValue{
			k: "has_default",
			v: strconv.FormatBool(field.HasDefault()),
		},
		debugValue{
			k: "has_ulid",
			v: strconv.FormatBool(field.HasULID()),
		},
		debugValue{
			k: "is_created_key",
			v: strconv.FormatBool(field.IsCreatedKey()),
		},
		debugValue{
			k: "is_updated_key",
			v: strconv.FormatBool(field.IsUpdatedKey()),
		},
		debugValue{
			k: "is_deleted_key",
			v: strconv.FormatBool(field.IsDeletedKey()),
		},
		debugValue{
			k: "reflect_type",
			v: field.rtype.String(),
		},
		debugValue{
			k: "association_type",
			v: field.associationType.String(),
		},
	}
}

func debugReference(reference Reference) debugWriter {
	return debugObj{
		debugValue{
			k: "model_name",
			v: reference.ModelName(),
		},
		debugValue{
			k: "table_name",
			v: reference.TableName(),
		},
		debugValue{
			k: "field_name",
			v: reference.FieldName(),
		},
		debugValue{
			k: "is_local",
			v: strconv.FormatBool(reference.IsLocal()),
		},
		debugWrap{
			k: "local",
			v: debugReferenceObject(reference.Local()),
		},
		debugWrap{
			k: "remote",
			v: debugReferenceObject(reference.Remote()),
		},
	}
}

func debugReferenceObject(reference ReferenceObject) debugWriter {
	return debugObj{
		debugValue{
			k: "schema",
			v: reference.Schema().ModelName(),
		},
		debugValue{
			k: "model_name",
			v: reference.ModelName(),
		},
		debugValue{
			k: "table_name",
			v: reference.TableName(),
		},
		debugValue{
			k: "field_name",
			v: reference.FieldName(),
		},
		debugValue{
			k: "column_path",
			v: reference.ColumnPath(),
		},
		debugValue{
			k: "column_name",
			v: reference.ColumnName(),
		},
		debugValue{
			k: "is_primary_key",
			v: strconv.FormatBool(reference.isPrimaryKey),
		},
		debugValue{
			k: "primary_key_type",
			v: reference.pkType.String(),
		},
		debugValue{
			k: "is_foreign_key",
			v: strconv.FormatBool(reference.isForeignKey),
		},
		debugValue{
			k: "foreign_key_type",
			v: reference.fkType.String(),
		},
	}
}
