package reflectx_test

import (
	"reflect"
	"testing"
	"time"

	"database/sql"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

// nolint: structcheck,megacheck
type Elements struct {
	Spirit     bool
	Air        uint8
	Water      uint16
	Earth      uint32
	Fire       uint64
	Combustion *uint8
	Snow       *uint16
	Plant      *uint32
	Lava       *uint64
	Blast      int8
	Sound      int16
	Ice        int32
	Metal      int64
	Cristal    *int8
	Gas        *int16
	Acid       *int32
	Petrolium  *int64
	Sand       string
	Lightning  *string
	Vaccum     float32
	Blood      *float32
	Plastic    float64
	Glass      *float64
	Oxygen     pq.NullTime
	Poison     *pq.NullTime
	Vapor      sql.NullInt64
	Dust       *sql.NullInt64
	Laser      sql.NullString
	Deflect    *sql.NullString
	Corrosion  sql.NullBool
	Rubber     *sql.NullBool
	Fiber      sql.NullFloat64
	Dioxide    *sql.NullFloat64
	Gravity    ElementInterface
	Pressure   ElementStruct
	Push       []byte
	Absortion  []rune
	Fragture   []int64
	Friction   []string
	String     time.Time
	Tension    *time.Time
	xA         []byte
	xB         bool
	xC         string
	xD         int64
}

type ElementInterface interface {
}

type ElementStruct struct {
}

func TestReflectx_GetFields(t *testing.T) {
	is := require.New(t)

	expected := []string{
		"Spirit",
		"Air",
		"Water",
		"Earth",
		"Fire",
		"Combustion",
		"Snow",
		"Plant",
		"Lava",
		"Blast",
		"Sound",
		"Ice",
		"Metal",
		"Cristal",
		"Gas",
		"Acid",
		"Petrolium",
		"Sand",
		"Lightning",
		"Vaccum",
		"Blood",
		"Plastic",
		"Glass",
		"Oxygen",
		"Poison",
		"Vapor",
		"Dust",
		"Laser",
		"Deflect",
		"Corrosion",
		"Rubber",
		"Fiber",
		"Dioxide",
		"Gravity",
		"Pressure",
		"Push",
		"Absortion",
		"Fragture",
		"Friction",
		"String",
		"Tension",
	}

	fields, err := reflectx.GetFields(&Elements{})
	is.NoError(err)
	is.NotEmpty(fields)
	is.Equal(expected, fields)

	fields, err = reflectx.GetFields(Elements{})
	is.NoError(err)
	is.NotEmpty(fields)
	is.Equal(expected, fields)

	fields, err = reflectx.GetFields(func() bool {
		return false
	})
	is.Error(err)
	is.Empty(fields)

	fields, err = reflectx.GetFields("hello world!")
	is.Error(err)
	is.Empty(fields)

	fields, err = reflectx.GetFields(Elements{}.Sand)
	is.Error(err)
	is.Empty(fields)
}

func TestReflectx_GetFieldByName(t *testing.T) {
	is := require.New(t)

	field, ok := reflectx.GetFieldByName(&Elements{}, "Snow")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Snow", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Ptr, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Snow")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Snow", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Ptr, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "Gravity")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Gravity", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Gravity")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Gravity", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "Metal")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Metal", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Metal")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Metal", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "xA")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xA", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Slice, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xB")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xB", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Bool, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xC")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xC", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.String, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xD")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xD", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xE")
	is.False(ok)
	is.Empty(field)
}

func TestReflectx_GetFieldValue(t *testing.T) {
	is := require.New(t)

	a := false
	b := "89db"
	c := time.Now()
	d := uint8(3)
	e := sql.NullFloat64{
		Valid:   true,
		Float64: 3.333,
	}
	f := ElementStruct{}

	elements := &Elements{
		Spirit:   a,
		Air:      d,
		Tension:  &c,
		Sand:     b,
		Fiber:    e,
		Pressure: f,
		xA:       []byte("hello"),
	}

	value, err := reflectx.GetFieldValue(elements, "Spirit")
	is.NoError(err)
	is.Equal(false, value)

	value, err = reflectx.GetFieldValue(elements, "Air")
	is.NoError(err)
	is.Equal(uint8(3), value)

	value, err = reflectx.GetFieldValue(elements, "Tension")
	is.NoError(err)
	is.Equal(&c, value)

	value, err = reflectx.GetFieldValue(elements, "Sand")
	is.NoError(err)
	is.Equal("89db", value)

	value, err = reflectx.GetFieldValue(elements, "Pressure")
	is.NoError(err)
	is.Equal(f, value)

	value, err = reflectx.GetFieldValue(elements, "xA")
	is.Error(err)
	is.Nil(value)
}

func TestReflectx_PushFieldValue_Users(t *testing.T) {
	is := require.New(t)

	type User struct {
		ID        string
		Enabled   bool
		Type      uint8
		Role      *uint8
		Counter   int32
		CreatedAt time.Time
		ProfileID sql.NullInt64
		AvatarID  *sql.NullString
		LastLogin *pq.NullTime
		Links     []string
	}

	u1 := User{}
	u2 := &User{}
	u3 := &User{}
	u4 := &User{}
	u5 := &User{}

	v1 := "01CNBT7JEJ1W7BTZGABTRZ0DE4"
	v2 := false
	v3 := true
	v4 := uint8(10)
	v5 := uint8(20)
	v6 := int32(200)
	v7 := int64(2000)
	v8 := sql.NullInt64{
		Valid: true,
		Int64: int64(v4),
	}
	v9 := sql.NullInt64{
		Valid: true,
		Int64: v7,
	}
	v10 := "01CNBYJ67KQHS263ERMS1PNAYT"
	v11 := "01CNBYJ7VFEBBD9ZQ7WA9P641K"
	v12 := sql.NullString{
		Valid:  true,
		String: v10,
	}
	v13 := sql.NullString{
		Valid:  true,
		String: v11,
	}
	v14 := time.Now()
	v15 := pq.NullTime{
		Valid: false,
	}
	v16 := pq.NullTime{
		Valid: true,
		Time:  v14,
	}
	v17 := "link1"
	v18 := "link2"
	v19 := "link3"
	v20 := []string{v17, v18, v19}

	err := reflectx.PushFieldValue(u1, "Hash", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(u1, "Name", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(u1, "ID", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(u2, "ID", v1, false)
	is.NoError(err)
	is.Equal(v1, u2.ID)

	err = reflectx.PushFieldValue(u3, "ID", &v1, false)
	is.NoError(err)
	is.Equal(v1, u3.ID)

	err = reflectx.PushFieldValue(u4, "ID", v1, true)
	is.NoError(err)
	is.Equal(v1, u4.ID)

	err = reflectx.PushFieldValue(u5, "ID", &v1, true)
	is.NoError(err)
	is.Equal(v1, u5.ID)

	err = reflectx.PushFieldValue(u2, "Enabled", v2, false)
	is.NoError(err)
	is.Equal(v2, u2.Enabled)

	err = reflectx.PushFieldValue(u3, "Enabled", &v2, false)
	is.NoError(err)
	is.Equal(v2, u3.Enabled)

	err = reflectx.PushFieldValue(u4, "Enabled", v2, true)
	is.NoError(err)
	is.Equal(v2, u4.Enabled)

	err = reflectx.PushFieldValue(u5, "Enabled", &v2, true)
	is.NoError(err)
	is.Equal(v2, u5.Enabled)

	err = reflectx.PushFieldValue(u2, "Enabled", v3, false)
	is.NoError(err)
	is.Equal(v3, u2.Enabled)

	err = reflectx.PushFieldValue(u3, "Enabled", &v3, false)
	is.NoError(err)
	is.Equal(v3, u3.Enabled)

	err = reflectx.PushFieldValue(u4, "Enabled", v3, true)
	is.NoError(err)
	is.Equal(v3, u4.Enabled)

	err = reflectx.PushFieldValue(u5, "Enabled", &v3, true)
	is.NoError(err)
	is.Equal(v3, u5.Enabled)

	err = reflectx.PushFieldValue(u2, "Type", v4, true)
	is.NoError(err)
	is.Equal(v4, u2.Type)

	err = reflectx.PushFieldValue(u3, "Type", &v4, true)
	is.NoError(err)
	is.Equal(v4, u3.Type)

	err = reflectx.PushFieldValue(u2, "Type", v5, true)
	is.NoError(err)
	is.Equal(v5, u2.Type)

	err = reflectx.PushFieldValue(u3, "Type", &v5, true)
	is.NoError(err)
	is.Equal(v5, u3.Type)

	err = reflectx.PushFieldValue(u4, "Role", v4, true)
	is.NoError(err)
	is.Equal(&v4, u4.Role)

	err = reflectx.PushFieldValue(u5, "Role", &v4, true)
	is.NoError(err)
	is.Equal(&v4, u5.Role)

	err = reflectx.PushFieldValue(u4, "Role", v5, true)
	is.NoError(err)
	is.Equal(&v5, u4.Role)

	err = reflectx.PushFieldValue(u5, "Role", &v5, true)
	is.NoError(err)
	is.Equal(&v5, u5.Role)

	err = reflectx.PushFieldValue(u2, "Counter", v6, true)
	is.NoError(err)
	is.Equal(v6, u2.Counter)

	err = reflectx.PushFieldValue(u2, "Counter", v5, true)
	is.Error(err)
	is.Equal(v6, u2.Counter)

	err = reflectx.PushFieldValue(u2, "Counter", v4, true)
	is.Error(err)
	is.Equal(v6, u2.Counter)

	err = reflectx.PushFieldValue(u3, "CreatedAt", v14, true)
	is.NoError(err)
	is.Equal(v14, u3.CreatedAt)

	err = reflectx.PushFieldValue(u5, "CreatedAt", v14, false)
	is.NoError(err)
	is.Equal(v14, u5.CreatedAt)

	err = reflectx.PushFieldValue(u2, "LastLogin", v14, true)
	is.NoError(err)
	is.NotNil(u2.LastLogin)
	is.True(u2.LastLogin.Valid)
	is.Equal(v14, u2.LastLogin.Time)

	err = reflectx.PushFieldValue(u4, "LastLogin", v15, true)
	is.NoError(err)
	is.NotNil(u4.LastLogin)
	is.Equal(&v15, u4.LastLogin)

	err = reflectx.PushFieldValue(u5, "LastLogin", v16, true)
	is.NoError(err)
	is.NotNil(u5.LastLogin)
	is.Equal(&v16, u5.LastLogin)

	err = reflectx.PushFieldValue(u2, "ProfileID", v7, true)
	is.NoError(err)
	is.True(u2.ProfileID.Valid)
	is.Equal(v7, u2.ProfileID.Int64)

	err = reflectx.PushFieldValue(u2, "ProfileID", v4, true)
	is.NoError(err)
	is.True(u2.ProfileID.Valid)
	is.Equal(int64(v4), u2.ProfileID.Int64)

	err = reflectx.PushFieldValue(u2, "ProfileID", v5, true)
	is.NoError(err)
	is.True(u2.ProfileID.Valid)
	is.Equal(int64(v5), u2.ProfileID.Int64)

	err = reflectx.PushFieldValue(u3, "ProfileID", &v7, true)
	is.NoError(err)
	is.True(u3.ProfileID.Valid)
	is.Equal(v7, u3.ProfileID.Int64)

	err = reflectx.PushFieldValue(u3, "ProfileID", &v4, true)
	is.NoError(err)
	is.True(u3.ProfileID.Valid)
	is.Equal(int64(v4), u3.ProfileID.Int64)

	err = reflectx.PushFieldValue(u3, "ProfileID", &v5, true)
	is.NoError(err)
	is.True(u3.ProfileID.Valid)
	is.Equal(int64(v5), u3.ProfileID.Int64)

	err = reflectx.PushFieldValue(u4, "ProfileID", v8, true)
	is.NoError(err)
	is.Equal(v8, u4.ProfileID)

	err = reflectx.PushFieldValue(u4, "ProfileID", v9, true)
	is.NoError(err)
	is.Equal(v9, u4.ProfileID)

	err = reflectx.PushFieldValue(u5, "ProfileID", &v8, true)
	is.NoError(err)
	is.Equal(v8, u5.ProfileID)

	err = reflectx.PushFieldValue(u5, "ProfileID", &v9, true)
	is.NoError(err)
	is.Equal(v9, u5.ProfileID)

	err = reflectx.PushFieldValue(u2, "AvatarID", v10, true)
	is.NoError(err)
	is.NotNil(u2.AvatarID)
	is.True(u2.AvatarID.Valid)
	is.Equal(v10, u2.AvatarID.String)

	err = reflectx.PushFieldValue(u2, "AvatarID", v11, true)
	is.NoError(err)
	is.NotNil(u2.AvatarID)
	is.True(u2.AvatarID.Valid)
	is.Equal(v11, u2.AvatarID.String)

	err = reflectx.PushFieldValue(u3, "AvatarID", &v10, true)
	is.NoError(err)
	is.NotNil(u3.AvatarID)
	is.True(u3.AvatarID.Valid)
	is.Equal(v10, u3.AvatarID.String)

	err = reflectx.PushFieldValue(u3, "AvatarID", &v11, true)
	is.NoError(err)
	is.NotNil(u3.AvatarID)
	is.True(u3.AvatarID.Valid)
	is.Equal(v11, u3.AvatarID.String)

	err = reflectx.PushFieldValue(u4, "AvatarID", v12, true)
	is.NoError(err)
	is.NotNil(u4.AvatarID)
	is.Equal(&v12, u4.AvatarID)

	err = reflectx.PushFieldValue(u4, "AvatarID", v13, true)
	is.NoError(err)
	is.NotNil(u4.AvatarID)
	is.Equal(&v13, u4.AvatarID)

	err = reflectx.PushFieldValue(u5, "AvatarID", &v12, true)
	is.NoError(err)
	is.NotNil(u5.AvatarID)
	is.Equal(&v12, u5.AvatarID)

	err = reflectx.PushFieldValue(u5, "AvatarID", &v13, true)
	is.NoError(err)
	is.NotNil(u5.AvatarID)
	is.Equal(&v13, u5.AvatarID)

	err = reflectx.PushFieldValue(u2, "Links", v17, false)
	is.NoError(err)
	is.NotEmpty(u2.Links)
	is.Len(u2.Links, 1)
	is.Equal(v17, u2.Links[0])

	err = reflectx.PushFieldValue(u2, "Links", v18, false)
	is.NoError(err)
	is.NotEmpty(u2.Links)
	is.Len(u2.Links, 2)
	is.Equal(v17, u2.Links[0])
	is.Equal(v18, u2.Links[1])

	err = reflectx.PushFieldValue(u2, "Links", v19, false)
	is.NoError(err)
	is.NotEmpty(u2.Links)
	is.Len(u2.Links, 3)
	is.Equal(v17, u2.Links[0])
	is.Equal(v18, u2.Links[1])
	is.Equal(v19, u2.Links[2])

	err = reflectx.PushFieldValue(u3, "Links", &v17, false)
	is.NoError(err)
	is.NotEmpty(u3.Links)
	is.Len(u3.Links, 1)
	is.Equal(v17, u3.Links[0])

	err = reflectx.PushFieldValue(u3, "Links", &v18, false)
	is.NoError(err)
	is.NotEmpty(u3.Links)
	is.Len(u3.Links, 2)
	is.Equal(v17, u3.Links[0])
	is.Equal(v18, u3.Links[1])

	err = reflectx.PushFieldValue(u3, "Links", &v19, false)
	is.NoError(err)
	is.NotEmpty(u3.Links)
	is.Len(u3.Links, 3)
	is.Equal(v17, u3.Links[0])
	is.Equal(v18, u3.Links[1])
	is.Equal(v19, u3.Links[2])

	err = reflectx.PushFieldValue(u4, "Links", v20, false)
	is.NoError(err)
	is.NotEmpty(u4.Links)
	is.Len(u4.Links, 3)
	is.Equal(v17, u4.Links[0])
	is.Equal(v18, u4.Links[1])
	is.Equal(v19, u4.Links[2])

	err = reflectx.PushFieldValue(u5, "Links", &v20, false)
	is.NoError(err)
	is.NotEmpty(u5.Links)
	is.Len(u5.Links, 3)
	is.Equal(v17, u5.Links[0])
	is.Equal(v18, u5.Links[1])
	is.Equal(v19, u5.Links[2])

}

func TestReflectx_PushFieldValue_Containers(t *testing.T) {
	is := require.New(t)

	type Element struct {
		ID   string
		Name string
	}

	type Container struct {
		ElemA Element
		ElemB *Element
		ElemC []Element
		ElemD []*Element
		ElemE *[]Element
		ElemF *[]*Element
	}

	v1 := Element{
		ID:   "01CNBDHSMKWACCNYTS4PG3D55A",
		Name: "Nora",
	}
	v2 := Element{
		ID:   "01CNBDJM17F5TE2S00YGKKS3PK",
		Name: "Costin",
	}
	v3 := Element{
		ID:   "01CNBDKACVQ1WTRHRZ98Y22NZZ",
		Name: "Najwa",
	}
	v4 := Element{
		ID:   "01CNBDM9H8WTJN7T62Q5893PHM",
		Name: "Reynard",
	}
	v5 := Element{
		ID:   "01CNBDMJJJ85DFRGV4S5J8N3VG",
		Name: "Chen",
	}
	v6 := Element{
		ID:   "01CNBDMVX8N1FMWMER2YB8WSTG",
		Name: "Ljupcho",
	}
	v7 := []Element{
		v1, v2, v3,
	}
	v8 := []*Element{
		&v4, &v5, &v6,
	}

	c1 := Container{}
	c2 := &Container{}
	c3 := &Container{}
	c4 := &Container{}
	c5 := &Container{}
	c6 := &Container{}
	c7 := &Container{}
	c8 := &Container{}

	err := reflectx.PushFieldValue(c1, "Hash", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(c1, "Name", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(c1, "ElemA", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(c1, "ElemB", v1, false)
	is.Error(err)

	err = reflectx.PushFieldValue(c2, "ElemA", v1, false)
	is.NoError(err)
	is.Equal(v1, c2.ElemA)

	err = reflectx.PushFieldValue(c2, "ElemB", v1, false)
	is.NoError(err)
	is.Equal(&v1, c2.ElemB)

	err = reflectx.PushFieldValue(c2, "ElemA", v2, false)
	is.NoError(err)
	is.Equal(v2, c2.ElemA)

	err = reflectx.PushFieldValue(c2, "ElemB", v2, false)
	is.NoError(err)
	is.Equal(&v2, c2.ElemB)

	err = reflectx.PushFieldValue(c2, "ElemA", &v1, false)
	is.NoError(err)
	is.Equal(v1, c2.ElemA)

	err = reflectx.PushFieldValue(c2, "ElemB", &v1, false)
	is.NoError(err)
	is.Equal(&v1, c2.ElemB)

	err = reflectx.PushFieldValue(c3, "ElemC", v1, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 1)
	is.Equal(v1, c3.ElemC[0])

	err = reflectx.PushFieldValue(c3, "ElemC", v2, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 2)
	is.Equal(v1, c3.ElemC[0])
	is.Equal(v2, c3.ElemC[1])

	err = reflectx.PushFieldValue(c3, "ElemC", &v3, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 3)
	is.Equal(v1, c3.ElemC[0])
	is.Equal(v2, c3.ElemC[1])
	is.Equal(v3, c3.ElemC[2])

	err = reflectx.PushFieldValue(c3, "ElemC", &v4, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 4)
	is.Equal(v1, c3.ElemC[0])
	is.Equal(v2, c3.ElemC[1])
	is.Equal(v3, c3.ElemC[2])
	is.Equal(v4, c3.ElemC[3])

	err = reflectx.PushFieldValue(c3, "ElemC", v5, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 5)
	is.Equal(v1, c3.ElemC[0])
	is.Equal(v2, c3.ElemC[1])
	is.Equal(v3, c3.ElemC[2])
	is.Equal(v4, c3.ElemC[3])
	is.Equal(v5, c3.ElemC[4])

	err = reflectx.PushFieldValue(c3, "ElemC", &v6, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemC)
	is.Len(c3.ElemC, 6)
	is.Equal(v1, c3.ElemC[0])
	is.Equal(v2, c3.ElemC[1])
	is.Equal(v3, c3.ElemC[2])
	is.Equal(v4, c3.ElemC[3])
	is.Equal(v5, c3.ElemC[4])
	is.Equal(v6, c3.ElemC[5])

	err = reflectx.PushFieldValue(c3, "ElemD", v2, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 1)
	is.Equal(&v2, c3.ElemD[0])

	err = reflectx.PushFieldValue(c3, "ElemD", v1, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 2)
	is.Equal(&v2, c3.ElemD[0])
	is.Equal(&v1, c3.ElemD[1])

	err = reflectx.PushFieldValue(c3, "ElemD", &v3, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 3)
	is.Equal(&v2, c3.ElemD[0])
	is.Equal(&v1, c3.ElemD[1])
	is.Equal(&v3, c3.ElemD[2])

	err = reflectx.PushFieldValue(c3, "ElemD", &v6, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 4)
	is.Equal(&v2, c3.ElemD[0])
	is.Equal(&v1, c3.ElemD[1])
	is.Equal(&v3, c3.ElemD[2])
	is.Equal(&v6, c3.ElemD[3])

	err = reflectx.PushFieldValue(c3, "ElemD", &v5, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 5)
	is.Equal(&v2, c3.ElemD[0])
	is.Equal(&v1, c3.ElemD[1])
	is.Equal(&v3, c3.ElemD[2])
	is.Equal(&v6, c3.ElemD[3])
	is.Equal(&v5, c3.ElemD[4])

	err = reflectx.PushFieldValue(c3, "ElemD", v4, false)
	is.NoError(err)
	is.NotEmpty(c3.ElemD)
	is.Len(c3.ElemD, 6)
	is.Equal(&v2, c3.ElemD[0])
	is.Equal(&v1, c3.ElemD[1])
	is.Equal(&v3, c3.ElemD[2])
	is.Equal(&v6, c3.ElemD[3])
	is.Equal(&v5, c3.ElemD[4])
	is.Equal(&v4, c3.ElemD[5])

	err = reflectx.PushFieldValue(c4, "ElemE", v1, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 1)
	is.Equal(v1, (*c4.ElemE)[0])

	err = reflectx.PushFieldValue(c4, "ElemE", v4, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 2)
	is.Equal(v1, (*c4.ElemE)[0])
	is.Equal(v4, (*c4.ElemE)[1])

	err = reflectx.PushFieldValue(c4, "ElemE", &v3, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 3)
	is.Equal(v1, (*c4.ElemE)[0])
	is.Equal(v4, (*c4.ElemE)[1])
	is.Equal(v3, (*c4.ElemE)[2])

	err = reflectx.PushFieldValue(c4, "ElemE", v2, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 4)
	is.Equal(v1, (*c4.ElemE)[0])
	is.Equal(v4, (*c4.ElemE)[1])
	is.Equal(v3, (*c4.ElemE)[2])
	is.Equal(v2, (*c4.ElemE)[3])

	err = reflectx.PushFieldValue(c4, "ElemE", &v5, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 5)
	is.Equal(v1, (*c4.ElemE)[0])
	is.Equal(v4, (*c4.ElemE)[1])
	is.Equal(v3, (*c4.ElemE)[2])
	is.Equal(v2, (*c4.ElemE)[3])
	is.Equal(v5, (*c4.ElemE)[4])

	err = reflectx.PushFieldValue(c4, "ElemE", &v6, false)
	is.NoError(err)
	is.NotNil(c4.ElemE)
	is.Len((*c4.ElemE), 6)
	is.Equal(v1, (*c4.ElemE)[0])
	is.Equal(v4, (*c4.ElemE)[1])
	is.Equal(v3, (*c4.ElemE)[2])
	is.Equal(v2, (*c4.ElemE)[3])
	is.Equal(v5, (*c4.ElemE)[4])
	is.Equal(v6, (*c4.ElemE)[5])

	err = reflectx.PushFieldValue(c4, "ElemF", v6, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 1)
	is.Equal(&v6, (*c4.ElemF)[0])

	err = reflectx.PushFieldValue(c4, "ElemF", v1, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 2)
	is.Equal(&v6, (*c4.ElemF)[0])
	is.Equal(&v1, (*c4.ElemF)[1])

	err = reflectx.PushFieldValue(c4, "ElemF", &v5, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 3)
	is.Equal(&v6, (*c4.ElemF)[0])
	is.Equal(&v1, (*c4.ElemF)[1])
	is.Equal(&v5, (*c4.ElemF)[2])

	err = reflectx.PushFieldValue(c4, "ElemF", &v2, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 4)
	is.Equal(&v6, (*c4.ElemF)[0])
	is.Equal(&v1, (*c4.ElemF)[1])
	is.Equal(&v5, (*c4.ElemF)[2])
	is.Equal(&v2, (*c4.ElemF)[3])

	err = reflectx.PushFieldValue(c4, "ElemF", v4, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 5)
	is.Equal(&v6, (*c4.ElemF)[0])
	is.Equal(&v1, (*c4.ElemF)[1])
	is.Equal(&v5, (*c4.ElemF)[2])
	is.Equal(&v2, (*c4.ElemF)[3])
	is.Equal(&v4, (*c4.ElemF)[4])

	err = reflectx.PushFieldValue(c4, "ElemF", &v3, false)
	is.NoError(err)
	is.NotNil(c4.ElemF)
	is.Len((*c4.ElemF), 6)
	is.Equal(&v6, (*c4.ElemF)[0])
	is.Equal(&v1, (*c4.ElemF)[1])
	is.Equal(&v5, (*c4.ElemF)[2])
	is.Equal(&v2, (*c4.ElemF)[3])
	is.Equal(&v4, (*c4.ElemF)[4])
	is.Equal(&v3, (*c4.ElemF)[5])

	err = reflectx.PushFieldValue(c5, "ElemC", &v7, true)
	is.NoError(err)
	is.NotEmpty(c5.ElemC)
	is.Len(c5.ElemC, 3)
	is.Equal(v1, c5.ElemC[0])
	is.Equal(v2, c5.ElemC[1])
	is.Equal(v3, c5.ElemC[2])

	err = reflectx.PushFieldValue(c6, "ElemC", v7, true)
	is.NoError(err)
	is.NotEmpty(c6.ElemC)
	is.Len(c6.ElemC, 3)
	is.Equal(v1, c6.ElemC[0])
	is.Equal(v2, c6.ElemC[1])
	is.Equal(v3, c6.ElemC[2])

	err = reflectx.PushFieldValue(c5, "ElemD", &v8, true)
	is.NoError(err)
	is.NotEmpty(c5.ElemD)
	is.Len(c5.ElemD, 3)
	is.Equal(&v4, c5.ElemD[0])
	is.Equal(&v5, c5.ElemD[1])
	is.Equal(&v6, c5.ElemD[2])

	err = reflectx.PushFieldValue(c6, "ElemD", v8, true)
	is.NoError(err)
	is.NotEmpty(c6.ElemD)
	is.Len(c6.ElemD, 3)
	is.Equal(&v4, c6.ElemD[0])
	is.Equal(&v5, c6.ElemD[1])
	is.Equal(&v6, c6.ElemD[2])

	err = reflectx.PushFieldValue(c7, "ElemE", &v7, true)
	is.NoError(err)
	is.NotNil(c7.ElemE)
	is.Len((*c7.ElemE), 3)
	is.Equal(v1, (*c7.ElemE)[0])
	is.Equal(v2, (*c7.ElemE)[1])
	is.Equal(v3, (*c7.ElemE)[2])

	err = reflectx.PushFieldValue(c8, "ElemE", v7, true)
	is.NoError(err)
	is.NotNil(c8.ElemE)
	is.Len((*c8.ElemE), 3)
	is.Equal(v1, (*c8.ElemE)[0])
	is.Equal(v2, (*c8.ElemE)[1])
	is.Equal(v3, (*c8.ElemE)[2])

	err = reflectx.PushFieldValue(c7, "ElemF", &v8, true)
	is.NoError(err)
	is.NotNil(c7.ElemF)
	is.Len((*c7.ElemF), 3)
	is.Equal(&v4, (*c7.ElemF)[0])
	is.Equal(&v5, (*c7.ElemF)[1])
	is.Equal(&v6, (*c7.ElemF)[2])

	err = reflectx.PushFieldValue(c8, "ElemF", v8, true)
	is.NoError(err)
	is.NotNil(c8.ElemF)
	is.Len((*c8.ElemF), 3)
	is.Equal(&v4, (*c8.ElemF)[0])
	is.Equal(&v5, (*c8.ElemF)[1])
	is.Equal(&v6, (*c8.ElemF)[2])

}
