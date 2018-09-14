package makroud_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud"
	"github.com/ulule/makroud/reflectx"
)

func TestTags_Analyze(t *testing.T) {
	is := require.New(t)

	elements := &Elements{}

	{
		field, ok := reflectx.GetFieldByName(elements, "Air")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Len(properties, 1)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("air", properties[0].Value())
	}
	{
		field, ok := reflectx.GetFieldByName(elements, "Fire")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Len(properties, 1)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("fire", properties[0].Value())
	}
	{
		field, ok := reflectx.GetFieldByName(elements, "Water")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Len(properties, 1)
		is.Equal(makroud.TagKeyIgnored, properties[0].Key())
		is.Equal("true", properties[0].Value())
	}
	{
		field, ok := reflectx.GetFieldByName(elements, "Earth")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Len(properties, 2)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("earth", properties[0].Value())
		is.Equal(makroud.TagKeyDefault, properties[1].Key())
		is.Equal("true", properties[1].Value())
	}
	{
		field, ok := reflectx.GetFieldByName(elements, "Fifth")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 0)
	}
	{
		field, ok := reflectx.GetFieldByName(elements, "Sixth")
		is.False(ok)
		is.Empty(field)
	}

	chunk := &ExoChunk{}

	{
		field, ok := reflectx.GetFieldByName(chunk, "Hash")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("hash", properties[0].Value())
		is.Equal(makroud.TagKeyPrimaryKey, properties[1].Key())
		is.Equal("ulid", properties[1].Value())
	}

	signature := &ExoChunkSignature{}

	{
		field, ok := reflectx.GetFieldByName(signature, "ChunkID")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("chunk_id", properties[0].Value())
		is.Equal(makroud.TagKeyForeignKey, properties[1].Key())
		is.Equal("exo_chunk", properties[1].Value())
	}

	region := &ExoRegion{}

	{
		field, ok := reflectx.GetFieldByName(region, "ID")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("id", properties[0].Value())
		is.Equal(makroud.TagKeyPrimaryKey, properties[1].Key())
		is.Equal("ulid", properties[1].Value())
	}

	owl := &Owl{}

	{
		field, ok := reflectx.GetFieldByName(owl, "ID")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("id", properties[0].Value())
		is.Equal(makroud.TagKeyPrimaryKey, properties[1].Key())
		is.Equal("true", properties[1].Value())
	}

	pack := &Package{}

	{
		field, ok := reflectx.GetFieldByName(pack, "ID")
		is.True(ok)
		is.NotEmpty(field)

		tags := makroud.GetTags(field)
		is.Len(tags, 1)
		name := tags[0].Name()
		properties := tags[0].Properties()
		is.Equal(makroud.TagName, name)
		is.Equal(makroud.TagKeyColumn, properties[0].Key())
		is.Equal("id", properties[0].Value())
	}

}
