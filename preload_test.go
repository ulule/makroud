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

		center1 := &Center{
			Name: "Soul",
			Area: "Lancaster",
		}
		err = sqlxx.Save(driver, center1)
		is.NoError(err)
		is.NotEmpty(center1.ID)

		center2 := &Center{
			Name: "Cloud",
			Area: "Nancledra",
		}
		err = sqlxx.Save(driver, center2)
		is.NoError(err)
		is.NotEmpty(center2.ID)

		center3 := &Center{
			Name: "Gold",
			Area: "Woodhurst",
		}
		err = sqlxx.Save(driver, center3)
		is.NoError(err)
		is.NotEmpty(center3.ID)

		center4 := &Center{
			Name: "Moonstone",
			Area: "Armskirk",
		}
		err = sqlxx.Save(driver, center4)
		is.NoError(err)
		is.NotEmpty(center4.ID)

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

		owl3 := &Owl{
			Name:         "Wacky",
			FeatherColor: "Harsh Cyan",
			FavoriteFood: "Pecan Trifle",
			GroupID: sql.NullInt64{
				Valid: true,
				Int64: group.ID,
			},
		}
		err = sqlxx.Save(driver, owl3)
		is.NoError(err)
		is.NotEmpty(owl3.ID)

		pack1 := &Package{
			SenderID:   center2.ID,
			ReceiverID: center1.ID,
			Status:     "processing",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl1.ID,
			},
		}
		err = sqlxx.Save(driver, pack1)
		is.NoError(err)
		is.NotEmpty(pack1.ID)

		pack2 := &Package{
			SenderID:   center2.ID,
			ReceiverID: center4.ID,
			Status:     "delivered",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl1.ID,
			},
		}
		err = sqlxx.Save(driver, pack2)
		is.NoError(err)
		is.NotEmpty(pack2.ID)

		pack3 := &Package{
			SenderID:   center2.ID,
			ReceiverID: center4.ID,
			Status:     "waiting",
			TransporterID: sql.NullInt64{
				Valid: false,
			},
		}
		err = sqlxx.Save(driver, pack3)
		is.NoError(err)
		is.NotEmpty(pack3.ID)

		pack4 := &Package{
			SenderID:   center1.ID,
			ReceiverID: center4.ID,
			Status:     "processing",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl2.ID,
			},
		}
		err = sqlxx.Save(driver, pack4)
		is.NoError(err)
		is.NotEmpty(pack4.ID)

		pack5 := &Package{
			SenderID:   center3.ID,
			ReceiverID: center4.ID,
			Status:     "delivered",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl3.ID,
			},
		}
		err = sqlxx.Save(driver, pack5)
		is.NoError(err)
		is.NotEmpty(pack5.ID)

		pack6 := &Package{
			SenderID:   center4.ID,
			ReceiverID: center3.ID,
			Status:     "delivered",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl3.ID,
			},
		}
		err = sqlxx.Save(driver, pack6)
		is.NoError(err)
		is.NotEmpty(pack6.ID)

		pack7 := &Package{
			SenderID:   center3.ID,
			ReceiverID: center2.ID,
			Status:     "delivered",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl3.ID,
			},
		}
		err = sqlxx.Save(driver, pack7)
		is.NoError(err)
		is.NotEmpty(pack7.ID)

		pack8 := &Package{
			SenderID:   center2.ID,
			ReceiverID: center3.ID,
			Status:     "delivered",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl3.ID,
			},
		}
		err = sqlxx.Save(driver, pack8)
		is.NoError(err)
		is.NotEmpty(pack8.ID)

		pack9 := &Package{
			SenderID:   center3.ID,
			ReceiverID: center1.ID,
			Status:     "processing",
			TransporterID: sql.NullInt64{
				Valid: true,
				Int64: owl3.ID,
			},
		}
		err = sqlxx.Save(driver, pack9)
		is.NoError(err)
		is.NotEmpty(pack9.ID)

		is.Nil(owl1.Group)
		is.Nil(owl2.Group)
		is.Nil(owl3.Group)
		is.Empty(owl1.Packages)
		is.Empty(owl2.Packages)
		is.Empty(owl3.Packages)

		err = sqlxx.Preload(driver, owl1, "Group")
		is.NoError(err)
		is.NotNil(owl1.Group)
		is.True(owl1.GroupID.Valid)
		is.Equal(group.ID, owl1.GroupID.Int64)
		is.Equal(group.ID, owl1.Group.ID)
		is.Equal(group.Name, owl1.Group.Name)
		is.Empty(owl1.Packages)

		err = sqlxx.Preload(driver, &owl1, "Packages")
		is.NoError(err)
		is.NotEmpty(owl1.Packages)
		is.Len(owl1.Packages, 2)
		is.Contains(owl1.Packages, *pack1)
		is.Contains(owl1.Packages, *pack2)

		err = sqlxx.Preload(driver, owl2, "Group", "Packages")
		is.NoError(err)
		is.Nil(owl2.Group)
		is.False(owl2.GroupID.Valid)
		is.NotEmpty(owl2.Packages)
		is.Len(owl2.Packages, 1)
		is.Contains(owl2.Packages, *pack4)

		err = sqlxx.Preload(driver, *owl3, "Group", "Packages")
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		is.Nil(owl3.Group)
		is.Empty(owl3.Packages)

		err = sqlxx.Preload(driver, owl3, "Group", "Packages")
		is.NoError(err)
		is.NotNil(owl3.Group)
		is.True(owl3.GroupID.Valid)
		is.Equal(group.ID, owl3.GroupID.Int64)
		is.Equal(group.ID, owl3.Group.ID)
		is.Equal(group.Name, owl3.Group.Name)
		is.NotEmpty(owl3.Packages)
		is.Len(owl3.Packages, 5)

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

		cat3 := &Cat{
			Name: "Icarus",
		}
		err = sqlxx.Save(driver, cat3)
		is.NoError(err)
		is.NotEmpty(cat3.ID)

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

		meow1 := &Meow{
			Body:  "Meow !",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(driver, meow1)
		is.NoError(err)
		is.NotEmpty(meow1.Hash)

		meow2 := &Meow{
			Body:  "Meow meow...",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(driver, meow2)
		is.NoError(err)
		is.NotEmpty(meow2.Hash)

		meow3 := &Meow{
			Body:  "Meow meow ? meeeeeoooow ?!",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(driver, meow3)
		is.NoError(err)
		is.NotEmpty(meow3.Hash)

		meow4 := &Meow{
			Body:  "Meow, meow meow.",
			CatID: cat3.ID,
		}
		err = sqlxx.Save(driver, meow4)
		is.NoError(err)
		is.NotEmpty(meow4.Hash)

		is.Nil(cat1.Owner)
		is.Empty(cat1.Meows)
		is.Nil(cat2.Owner)
		is.Empty(cat2.Meows)
		is.Nil(cat3.Owner)
		is.Empty(cat3.Meows)

		err = sqlxx.Preload(driver, cat1, "Owner")
		is.NoError(err)
		is.Nil(cat1.Owner)

		err = sqlxx.Preload(driver, cat2, "Owner")
		is.NoError(err)
		is.NotNil(cat2.Owner)
		is.Equal(human2.ID, cat2.Owner.ID)
		is.Equal(human2.Name, cat2.Owner.Name)

		err = sqlxx.Preload(driver, &cat1, "Meows")
		is.NoError(err)
		is.NotEmpty(cat1.Meows)
		is.Len(cat1.Meows, 3)
		is.Contains(cat1.Meows, meow1)
		is.Contains(cat1.Meows, meow2)
		is.Contains(cat1.Meows, meow3)

		err = sqlxx.Preload(driver, &cat2, "Meows")
		is.NoError(err)
		is.Empty(cat2.Meows)

		err = sqlxx.Preload(driver, cat3, "Owner", "Meows")
		is.NoError(err)
		is.Nil(cat3.Owner)
		is.NotEmpty(cat3.Meows)
		is.Len(cat3.Meows, 1)
		is.Contains(cat3.Meows, meow4)

	})
}
