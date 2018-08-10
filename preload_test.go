package sqlxx_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestPreload_ExoChunk(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		mode := &ExoChunkMode{
			Mode: "rwx",
		}
		err := sqlxx.Save(driver, mode)
		is.NoError(err)
		is.NotEmpty(mode.ID)

		chunk := &ExoChunk{
			ModeID: mode.ID,
			Bytes: fmt.Sprint(
				"4e919ca20b565bb5b03c65130e018ad23d489412352ac8f25f3d0f8dd64905bd",
				"8bf1ee3f3f3a1715656f6c39631a5072e5d2afa23ecebe00c61fb05b54652bdb",
				"ea7548dec5b924a5a7ff2bd94dbe9a109849a3ea322919cc672980d037a325da",
			),
		}
		err = sqlxx.Save(driver, chunk)
		is.NoError(err)
		is.NotEmpty(chunk.Hash)

		signature := &ExoChunkSignature{
			ChunkID: chunk.Hash,
			Bytes: fmt.Sprint(
				"ed4709d761b35df76c1ecf6990f7703bb3e5027a5a3a434b3a4af92afcf9bcb1",
				"67c92c907edf6a68847e3aab6210ff1537e3e1ae079177feded543bb8ee35132",
			),
		}
		err = sqlxx.Save(driver, signature)
		is.NoError(err)
		is.NotEmpty(signature.ID)

		is.Nil(chunk.Mode)
		is.Nil(chunk.Signature)

		err = sqlxx.Preload(driver, chunk, "Mode")
		is.NoError(err)
		is.NotNil(chunk.Mode)
		is.Equal(mode.ID, chunk.ModeID)
		is.Equal(mode.ID, chunk.Mode.ID)
		is.Equal(mode.Mode, chunk.Mode.Mode)

		err = sqlxx.Preload(driver, chunk, "Signature")
		is.NoError(err)
		is.NotNil(chunk.Signature)
		is.Equal(signature.ChunkID, chunk.Hash)
		is.Equal(signature.ID, chunk.Signature.ID)
		is.Equal(signature.ChunkID, chunk.Signature.ChunkID)
		is.Equal(signature.Bytes, chunk.Signature.Bytes)

		chunk = &ExoChunk{
			Hash:   sqlxx.GenerateULID(driver),
			ModeID: 6000,
			Bytes: fmt.Sprint(
				"2eaf31b43c3c215c2aaaa7a5825c68fb97ad4913eedee90f16792e2d4881a7ef",
				"45ac26550ac888b33d52ce69bad114135ce591397d1d23dd2a4021dfb09de3f0",
				"86213a4e0fa96d84b1b997de3d9552b1a1d53c559e61abdc486d36735a06b9de",
			),
		}
		err = sqlxx.Preload(driver, chunk, "Mode")
		is.Error(err)
		is.Equal(sqlxx.ErrPreloadInvalidModel, errors.Cause(err))
		is.Nil(chunk.Mode)

	})
}

func TestPreload_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		group := &Group{
			Name: "spring",
		}
		err := sqlxx.Save(driver, group)
		is.NoError(err)
		is.NotEmpty(group.ID)

		owl1 := &Owl{
			Name:         "Pyro",
			FeatherColor: "Timeless Sanguine",
			FavoriteFood: "Ginger Mooncake",
			GroupID: sql.NullInt64{
				Valid: true,
				Int64: group.ID,
			},
		}
		err = sqlxx.Save(driver, owl1)
		is.NoError(err)
		is.NotEmpty(owl1.ID)

		owl2 := &Owl{
			Name:         "Bungee",
			FeatherColor: "Peaceful Peach",
			FavoriteFood: "Lemon Venison",
			GroupID: sql.NullInt64{
				Valid: false,
			},
		}
		err = sqlxx.Save(driver, owl2)
		is.NoError(err)
		is.NotEmpty(owl2.ID)

		is.Nil(owl1.Group)
		is.Nil(owl2.Group)

		err = sqlxx.Preload(driver, owl1, "Group")
		is.NoError(err)
		is.NotNil(owl1.Group)
		is.True(owl1.GroupID.Valid)
		is.Equal(group.ID, owl1.GroupID.Int64)
		is.Equal(group.ID, owl1.Group.ID)
		is.Equal(group.Name, owl1.Group.Name)

		err = sqlxx.Preload(driver, owl2, "Group")
		is.NoError(err)
		is.Nil(owl2.Group)
		is.False(owl2.GroupID.Valid)

	})
}

func TestPreload_Cat(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		cat1 := &Cat{
			Name: "Dinky",
		}
		err := sqlxx.Save(driver, cat1)
		is.NoError(err)
		is.NotEmpty(cat1.ID)

		cat2 := &Cat{
			Name: "Flick",
		}
		err = sqlxx.Save(driver, cat2)
		is.NoError(err)
		is.NotEmpty(cat2.ID)

		human1 := &Human{
			Name: "Larry Golade",
		}
		err = sqlxx.Save(driver, human1)
		is.NoError(err)
		is.NotEmpty(human1.ID)

		human2 := &Human{
			Name: "Roland Culé",
			CatID: sql.NullString{
				Valid:  true,
				String: cat2.ID,
			},
		}
		err = sqlxx.Save(driver, human2)
		is.NoError(err)
		is.NotEmpty(human2.ID)

		is.Nil(cat1.Owner)
		is.Nil(cat2.Owner)

		err = sqlxx.Preload(driver, cat1, "Owner")
		is.NoError(err)
		is.Nil(cat1.Owner)

		err = sqlxx.Preload(driver, cat2, "Owner")
		is.NoError(err)
		is.NotNil(cat2.Owner)
		is.Equal(human2.ID, cat2.Owner.ID)
		is.Equal(human2.Name, cat2.Owner.Name)

	})
}
