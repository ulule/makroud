package sqlxx_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestPreload_ExoChunk_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		mode := &ExoChunkMode{
			Mode: "rwx",
		}
		err := sqlxx.Save(ctx, driver, mode)
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
		err = sqlxx.Save(ctx, driver, chunk)
		is.NoError(err)
		is.NotEmpty(chunk.Hash)

		signature := &ExoChunkSignature{
			ChunkID: chunk.Hash,
			Bytes: fmt.Sprint(
				"ed4709d761b35df76c1ecf6990f7703bb3e5027a5a3a434b3a4af92afcf9bcb1",
				"67c92c907edf6a68847e3aab6210ff1537e3e1ae079177feded543bb8ee35132",
			),
		}
		err = sqlxx.Save(ctx, driver, signature)
		is.NoError(err)
		is.NotEmpty(signature.ID)

		is.Nil(chunk.Mode)
		is.Nil(chunk.Signature)

		err = sqlxx.Preload(ctx, driver, chunk, "Mode")
		is.NoError(err)
		is.NotNil(chunk.Mode)
		is.Equal(mode.ID, chunk.ModeID)
		is.Equal(mode.ID, chunk.Mode.ID)
		is.Equal(mode.Mode, chunk.Mode.Mode)

		err = sqlxx.Preload(ctx, driver, chunk, "Signature")
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
		err = sqlxx.Preload(ctx, driver, chunk, "Mode")
		is.NoError(err)
		is.Nil(chunk.Mode)

	})
}

func TestPreload_ExoChunk_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		mode1 := &ExoChunkMode{
			Mode: "rwx",
		}
		err := sqlxx.Save(ctx, driver, mode1)
		is.NoError(err)
		is.NotEmpty(mode1.ID)

		mode2 := &ExoChunkMode{
			Mode: "r-x",
		}
		err = sqlxx.Save(ctx, driver, mode2)
		is.NoError(err)
		is.NotEmpty(mode2.ID)

		mode3 := &ExoChunkMode{
			Mode: "r-x",
		}
		err = sqlxx.Save(ctx, driver, mode3)
		is.NoError(err)
		is.NotEmpty(mode3.ID)

		mode4 := &ExoChunkMode{
			Mode: "rw-",
		}
		err = sqlxx.Save(ctx, driver, mode4)
		is.NoError(err)
		is.NotEmpty(mode4.ID)

		chunk1 := &ExoChunk{
			ModeID: mode1.ID,
			Bytes: fmt.Sprint(
				"4e919ca20b565bb5b03c65130e018ad23d489412352ac8f25f3d0f8dd64905bd",
				"8bf1ee3f3f3a1715656f6c39631a5072e5d2afa23ecebe00c61fb05b54652bdb",
				"ea7548dec5b924a5a7ff2bd94dbe9a109849a3ea322919cc672980d037a325da",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk1)
		is.NoError(err)
		is.NotEmpty(chunk1.Hash)

		chunk2 := &ExoChunk{
			ModeID: mode1.ID,
			Bytes: fmt.Sprint(
				"455a9ccb2316fbbeefc621809bdb020986337fcd99a2e497f39aa76b67840a21",
				"fd6ad7d1877b18ef29026d497a70a573a0bb0329739fa1e51ad0f864b43474c2",
				"007ffae918ef2b7f9918db0da37b75c32d70e496f0a630ccec0fb0d1c12bdfbc",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk2)
		is.NoError(err)
		is.NotEmpty(chunk2.Hash)

		chunk3 := &ExoChunk{
			ModeID: mode1.ID,
			Bytes: fmt.Sprint(
				"5c4ff25049b5f36621f4c9e1e3723b43fc21d9008b1fb7bb90ce7e7d2dec11c9",
				"9eceb0c2e250ab4ebc895db91f27d861762d6167d316af43a2bdb6c777778bac",
				"e9a51d6f44d013a38220246bba5f5bf9c46e78941ebe6669fdf141477dddffc6",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk3)
		is.NoError(err)
		is.NotEmpty(chunk3.Hash)

		chunk4 := &ExoChunk{
			ModeID: mode2.ID,
			Bytes: fmt.Sprint(
				"d9c234d4934ba063bf9c80fda227529f344580ef237a53419b6710d3184dfb3f",
				"a59e8a8b81057cd9e4cf54c8af62a4e465f154d2102b4ddd0806ff51de08bd6e",
				"083e2565506d8a0f59437a58685f468eb3177f26190552cfac7d93631542918d",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk4)
		is.NoError(err)
		is.NotEmpty(chunk4.Hash)

		chunk5 := &ExoChunk{
			ModeID: mode2.ID,
			Bytes: fmt.Sprint(
				"f1676c02af03cfe8485f82825d71afa09184ba49d0304a8988dd2c75fb593858",
				"e59108f9554f7143ab0bb851cd6a301c2bce487f3398976a5b4f18f576e61b1c",
				"b532922f3c436f92b350871cbb39eabce016d4ab0e465eb22db4a221be13c985",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk5)
		is.NoError(err)
		is.NotEmpty(chunk5.Hash)

		chunk6 := &ExoChunk{
			ModeID: mode3.ID,
			Bytes: fmt.Sprint(
				"f1676c02af03cfe8485f82825d71afa09184ba49d0304a8988dd2c75fb593858",
				"e59108f9554f7143ab0bb851cd6a301c2bce487f3398976a5b4f18f576e61b1c",
				"b532922f3c436f92b350871cbb39eabce016d4ab0e465eb22db4a221be13c985",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk6)
		is.NoError(err)
		is.NotEmpty(chunk6.Hash)

		chunk7 := &ExoChunk{
			ModeID: mode4.ID,
			Bytes: fmt.Sprint(
				"f1676c02af03cfe8485f82825d71afa09184ba49d0304a8988dd2c75fb593858",
				"e59108f9554f7143ab0bb851cd6a301c2bce487f3398976a5b4f18f576e61b1c",
				"b532922f3c436f92b350871cbb39eabce016d4ab0e465eb22db4a221be13c985",
			),
		}
		err = sqlxx.Save(ctx, driver, chunk7)
		is.NoError(err)
		is.NotEmpty(chunk7.Hash)

		signature1 := &ExoChunkSignature{
			ChunkID: chunk1.Hash,
			Bytes: fmt.Sprint(
				"e5544175f3b80b488234a708be21ca3586b8dc5c9b42bfc2fcce883759bc7f4c",
				"b360e640a0264737e0e96cad0df2de60daa2920a58e4080672e1596c6b613c5e",
			),
		}
		err = sqlxx.Save(ctx, driver, signature1)
		is.NoError(err)
		is.NotEmpty(signature1.ID)

		signature2 := &ExoChunkSignature{
			ChunkID: chunk2.Hash,
			Bytes: fmt.Sprint(
				"e8cfe78d34e59d612e39c1a5c6e6b560b765b699aa4777da26d23a5042e7c104",
				"090888070233ec70f7035ead153fd217c8b0204687fc1b5ebd3e5701692467d2",
			),
		}
		err = sqlxx.Save(ctx, driver, signature2)
		is.NoError(err)
		is.NotEmpty(signature2.ID)

		signature3 := &ExoChunkSignature{
			ChunkID: chunk3.Hash,
			Bytes: fmt.Sprint(
				"53c69ef75f8ce6f6d876a095cc6adda9e32e417334b14d35b6de783235332729",
				"cba39a4675a475d7b33bd5eec55f1391eaa7e7484feffdd303f62c37a1665009",
			),
		}
		err = sqlxx.Save(ctx, driver, signature3)
		is.NoError(err)
		is.NotEmpty(signature3.ID)

		signature4 := &ExoChunkSignature{
			ChunkID: chunk6.Hash,
			Bytes: fmt.Sprint(
				"e867e17445730f555950d16ee4bbdef7c37cb4b687cc7fccdd99e0e92000102b",
				"c0d4e571cf9b0201f6c68e3725baa3876108ab6881969e38c86bb6308dd729a6",
			),
		}
		err = sqlxx.Save(ctx, driver, signature4)
		is.NoError(err)
		is.NotEmpty(signature4.ID)

		signature5 := &ExoChunkSignature{
			ChunkID: chunk7.Hash,
			Bytes: fmt.Sprint(
				"e70f8cdcf10c65d955b638c45ea0a10c34eb8c9c7a82de4bea897512070d5de9",
				"1c258d4369f3446ad9836c287f8c8b92844b57ce540fffd3a4dca6b053a1a1ea",
			),
		}
		err = sqlxx.Save(ctx, driver, signature5)
		is.NoError(err)
		is.NotEmpty(signature5.ID)

		is.Nil(chunk1.Mode)
		is.Nil(chunk1.Signature)
		is.Nil(chunk2.Mode)
		is.Nil(chunk2.Signature)
		is.Nil(chunk3.Mode)
		is.Nil(chunk3.Signature)
		is.Nil(chunk4.Mode)
		is.Nil(chunk4.Signature)
		is.Nil(chunk5.Mode)
		is.Nil(chunk5.Signature)
		is.Nil(chunk6.Mode)
		is.Nil(chunk6.Signature)
		is.Nil(chunk7.Mode)
		is.Nil(chunk7.Signature)

		{
			chunks := []ExoChunk{*chunk1, *chunk2, *chunk3, *chunk4, *chunk5, *chunk6, *chunk7}
			err = sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 7)
			is.Equal(chunk1.Hash, chunks[0].Hash)
			is.Equal(chunk2.Hash, chunks[1].Hash)
			is.Equal(chunk3.Hash, chunks[2].Hash)
			is.Equal(chunk4.Hash, chunks[3].Hash)
			is.Equal(chunk5.Hash, chunks[4].Hash)
			is.Equal(chunk6.Hash, chunks[5].Hash)
			is.Equal(chunk7.Hash, chunks[6].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(mode1.ID, chunks[0].Mode.ID)
			is.Equal(mode1.Mode, chunks[0].Mode.Mode)
			is.NotNil(chunks[0].Signature)
			is.Equal(signature1.ID, chunks[0].Signature.ID)
			is.Equal(signature1.ChunkID, chunks[0].Signature.ChunkID)
			is.Equal(signature1.Bytes, chunks[0].Signature.Bytes)

			is.NotNil(chunks[1].Mode)
			is.Equal(mode1.ID, chunks[1].Mode.ID)
			is.Equal(mode1.Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(signature2.ID, chunks[1].Signature.ID)
			is.Equal(signature2.ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(signature2.Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(mode1.ID, chunks[2].Mode.ID)
			is.Equal(mode1.Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(signature3.ID, chunks[2].Signature.ID)
			is.Equal(signature3.ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(signature3.Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(mode2.ID, chunks[3].Mode.ID)
			is.Equal(mode2.Mode, chunks[3].Mode.Mode)
			is.Nil(chunks[3].Signature)

			is.NotNil(chunks[4].Mode)
			is.Equal(mode2.ID, chunks[4].Mode.ID)
			is.Equal(mode2.Mode, chunks[4].Mode.Mode)
			is.Nil(chunks[4].Signature)

			is.NotNil(chunks[5].Mode)
			is.Equal(mode3.ID, chunks[5].Mode.ID)
			is.Equal(mode3.Mode, chunks[5].Mode.Mode)
			is.NotNil(chunks[5].Signature)
			is.Equal(signature4.ID, chunks[5].Signature.ID)
			is.Equal(signature4.ChunkID, chunks[5].Signature.ChunkID)
			is.Equal(signature4.Bytes, chunks[5].Signature.Bytes)

			is.NotNil(chunks[6].Mode)
			is.Equal(mode4.ID, chunks[6].Mode.ID)
			is.Equal(mode4.Mode, chunks[6].Mode.Mode)
			is.NotNil(chunks[6].Signature)
			is.Equal(signature5.ID, chunks[6].Signature.ID)
			is.Equal(signature5.ChunkID, chunks[6].Signature.ChunkID)
			is.Equal(signature5.Bytes, chunks[6].Signature.Bytes)

		}

		{
			chunks := []*ExoChunk{chunk1, chunk2, chunk3, chunk4, chunk5, chunk6, chunk7}
			err = sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 7)
			is.Equal(chunk1.Hash, chunks[0].Hash)
			is.Equal(chunk2.Hash, chunks[1].Hash)
			is.Equal(chunk3.Hash, chunks[2].Hash)
			is.Equal(chunk4.Hash, chunks[3].Hash)
			is.Equal(chunk5.Hash, chunks[4].Hash)
			is.Equal(chunk6.Hash, chunks[5].Hash)
			is.Equal(chunk7.Hash, chunks[6].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(mode1.ID, chunks[0].Mode.ID)
			is.Equal(mode1.Mode, chunks[0].Mode.Mode)
			is.NotNil(chunks[0].Signature)
			is.Equal(signature1.ID, chunks[0].Signature.ID)
			is.Equal(signature1.ChunkID, chunks[0].Signature.ChunkID)
			is.Equal(signature1.Bytes, chunks[0].Signature.Bytes)

			is.NotNil(chunks[1].Mode)
			is.Equal(mode1.ID, chunks[1].Mode.ID)
			is.Equal(mode1.Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(signature2.ID, chunks[1].Signature.ID)
			is.Equal(signature2.ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(signature2.Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(mode1.ID, chunks[2].Mode.ID)
			is.Equal(mode1.Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(signature3.ID, chunks[2].Signature.ID)
			is.Equal(signature3.ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(signature3.Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(mode2.ID, chunks[3].Mode.ID)
			is.Equal(mode2.Mode, chunks[3].Mode.Mode)
			is.Nil(chunks[3].Signature)

			is.NotNil(chunks[4].Mode)
			is.Equal(mode2.ID, chunks[4].Mode.ID)
			is.Equal(mode2.Mode, chunks[4].Mode.Mode)
			is.Nil(chunks[4].Signature)

			is.NotNil(chunks[5].Mode)
			is.Equal(mode3.ID, chunks[5].Mode.ID)
			is.Equal(mode3.Mode, chunks[5].Mode.Mode)
			is.NotNil(chunks[5].Signature)
			is.Equal(signature4.ID, chunks[5].Signature.ID)
			is.Equal(signature4.ChunkID, chunks[5].Signature.ChunkID)
			is.Equal(signature4.Bytes, chunks[5].Signature.Bytes)

			is.NotNil(chunks[6].Mode)
			is.Equal(mode4.ID, chunks[6].Mode.ID)
			is.Equal(mode4.Mode, chunks[6].Mode.Mode)
			is.NotNil(chunks[6].Signature)
			is.Equal(signature5.ID, chunks[6].Signature.ID)
			is.Equal(signature5.ChunkID, chunks[6].Signature.ChunkID)
			is.Equal(signature5.Bytes, chunks[6].Signature.Bytes)

		}
	})
}

func TestPreload_Owl_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		group := &Group{
			Name: "spring",
		}
		err := sqlxx.Save(ctx, driver, group)
		is.NoError(err)
		is.NotEmpty(group.ID)

		center1 := &Center{
			Name: "Soul",
			Area: "Lancaster",
		}
		err = sqlxx.Save(ctx, driver, center1)
		is.NoError(err)
		is.NotEmpty(center1.ID)

		center2 := &Center{
			Name: "Cloud",
			Area: "Nancledra",
		}
		err = sqlxx.Save(ctx, driver, center2)
		is.NoError(err)
		is.NotEmpty(center2.ID)

		center3 := &Center{
			Name: "Gold",
			Area: "Woodhurst",
		}
		err = sqlxx.Save(ctx, driver, center3)
		is.NoError(err)
		is.NotEmpty(center3.ID)

		center4 := &Center{
			Name: "Moonstone",
			Area: "Armskirk",
		}
		err = sqlxx.Save(ctx, driver, center4)
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
		err = sqlxx.Save(ctx, driver, owl1)
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
		err = sqlxx.Save(ctx, driver, owl2)
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
		err = sqlxx.Save(ctx, driver, owl3)
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
		err = sqlxx.Save(ctx, driver, pack1)
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
		err = sqlxx.Save(ctx, driver, pack2)
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
		err = sqlxx.Save(ctx, driver, pack3)
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
		err = sqlxx.Save(ctx, driver, pack4)
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
		err = sqlxx.Save(ctx, driver, pack5)
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
		err = sqlxx.Save(ctx, driver, pack6)
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
		err = sqlxx.Save(ctx, driver, pack7)
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
		err = sqlxx.Save(ctx, driver, pack8)
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
		err = sqlxx.Save(ctx, driver, pack9)
		is.NoError(err)
		is.NotEmpty(pack9.ID)

		is.Nil(owl1.Group)
		is.Nil(owl2.Group)
		is.Nil(owl3.Group)
		is.Empty(owl1.Packages)
		is.Empty(owl2.Packages)
		is.Empty(owl3.Packages)

		err = sqlxx.Preload(ctx, driver, owl1, "Group")
		is.NoError(err)
		is.NotNil(owl1.Group)
		is.True(owl1.GroupID.Valid)
		is.Equal(group.ID, owl1.GroupID.Int64)
		is.Equal(group.ID, owl1.Group.ID)
		is.Equal(group.Name, owl1.Group.Name)
		is.Empty(owl1.Packages)

		err = sqlxx.Preload(ctx, driver, &owl1, "Packages")
		is.NoError(err)
		is.NotEmpty(owl1.Packages)
		is.Len(owl1.Packages, 2)
		is.Contains(owl1.Packages, *pack1)
		is.Contains(owl1.Packages, *pack2)

		err = sqlxx.Preload(ctx, driver, owl2, "Group", "Packages")
		is.NoError(err)
		is.Nil(owl2.Group)
		is.False(owl2.GroupID.Valid)
		is.NotEmpty(owl2.Packages)
		is.Len(owl2.Packages, 1)
		is.Contains(owl2.Packages, *pack4)

		err = sqlxx.Preload(ctx, driver, *owl3, "Group", "Packages")
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		is.Nil(owl3.Group)
		is.Empty(owl3.Packages)

		err = sqlxx.Preload(ctx, driver, owl3, "Group", "Packages")
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

func TestPreload_Cat_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{
			Name: "Dinky",
		}
		err := sqlxx.Save(ctx, driver, cat1)
		is.NoError(err)
		is.NotEmpty(cat1.ID)

		cat2 := &Cat{
			Name: "Flick",
		}
		err = sqlxx.Save(ctx, driver, cat2)
		is.NoError(err)
		is.NotEmpty(cat2.ID)

		cat3 := &Cat{
			Name: "Icarus",
		}
		err = sqlxx.Save(ctx, driver, cat3)
		is.NoError(err)
		is.NotEmpty(cat3.ID)

		human1 := &Human{
			Name: "Larry Golade",
		}
		err = sqlxx.Save(ctx, driver, human1)
		is.NoError(err)
		is.NotEmpty(human1.ID)

		human2 := &Human{
			Name: "Roland Culé",
			CatID: sql.NullString{
				Valid:  true,
				String: cat2.ID,
			},
		}
		err = sqlxx.Save(ctx, driver, human2)
		is.NoError(err)
		is.NotEmpty(human2.ID)

		meow1 := &Meow{
			Body:  "Meow !",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow1)
		is.NoError(err)
		is.NotEmpty(meow1.Hash)

		meow2 := &Meow{
			Body:  "Meow meow...",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow2)
		is.NoError(err)
		is.NotEmpty(meow2.Hash)

		meow3 := &Meow{
			Body:  "Meow meow ? meeeeeoooow ?!",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow3)
		is.NoError(err)
		is.NotEmpty(meow3.Hash)

		meow4 := &Meow{
			Body:  "Meow, meow meow.",
			CatID: cat3.ID,
		}
		err = sqlxx.Save(ctx, driver, meow4)
		is.NoError(err)
		is.NotEmpty(meow4.Hash)

		is.Nil(cat1.Owner)
		is.Empty(cat1.Meows)
		is.Nil(cat2.Owner)
		is.Empty(cat2.Meows)
		is.Nil(cat3.Owner)
		is.Empty(cat3.Meows)

		err = sqlxx.Preload(ctx, driver, cat1, "Owner")
		is.NoError(err)
		is.Nil(cat1.Owner)

		err = sqlxx.Preload(ctx, driver, cat2, "Owner")
		is.NoError(err)
		is.NotNil(cat2.Owner)
		is.Equal(human2.ID, cat2.Owner.ID)
		is.Equal(human2.Name, cat2.Owner.Name)

		err = sqlxx.Preload(ctx, driver, &cat1, "Meows")
		is.NoError(err)
		is.NotEmpty(cat1.Meows)
		is.Len(cat1.Meows, 3)
		is.Contains(cat1.Meows, meow1)
		is.Contains(cat1.Meows, meow2)
		is.Contains(cat1.Meows, meow3)

		err = sqlxx.Preload(ctx, driver, &cat2, "Meows")
		is.NoError(err)
		is.Empty(cat2.Meows)

		err = sqlxx.Preload(ctx, driver, cat3, "Owner", "Meows")
		is.NoError(err)
		is.Nil(cat3.Owner)
		is.NotEmpty(cat3.Meows)
		is.Len(cat3.Meows, 1)
		is.Contains(cat3.Meows, meow4)

	})
}

func TestPreload_Cat_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{
			Name: "Eagle",
		}
		err := sqlxx.Save(ctx, driver, cat1)
		is.NoError(err)
		is.NotEmpty(cat1.ID)

		cat2 := &Cat{
			Name: "Zigzag",
		}
		err = sqlxx.Save(ctx, driver, cat2)
		is.NoError(err)
		is.NotEmpty(cat2.ID)

		cat3 := &Cat{
			Name: "Scully",
		}
		err = sqlxx.Save(ctx, driver, cat3)
		is.NoError(err)
		is.NotEmpty(cat3.ID)

		human1 := &Human{
			Name: "André Naline",
			CatID: sql.NullString{
				Valid:  true,
				String: cat1.ID,
			},
		}
		err = sqlxx.Save(ctx, driver, human1)
		is.NoError(err)
		is.NotEmpty(human1.ID)

		human2 := &Human{
			Name: "Garcin Lazare",
			CatID: sql.NullString{
				Valid:  true,
				String: cat2.ID,
			},
		}
		err = sqlxx.Save(ctx, driver, human2)
		is.NoError(err)
		is.NotEmpty(human2.ID)

		meow1 := &Meow{
			Body:  "Meow !",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow1)
		is.NoError(err)
		is.NotEmpty(meow1.Hash)

		meow2 := &Meow{
			Body:  "Meow meow...",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow2)
		is.NoError(err)
		is.NotEmpty(meow2.Hash)

		meow3 := &Meow{
			Body:  "Meow meow ? meeeeeoooow ?!",
			CatID: cat1.ID,
		}
		err = sqlxx.Save(ctx, driver, meow3)
		is.NoError(err)
		is.NotEmpty(meow3.Hash)

		meow4 := &Meow{
			Body:  "Meow, meow meow.",
			CatID: cat3.ID,
		}
		err = sqlxx.Save(ctx, driver, meow4)
		is.NoError(err)
		is.NotEmpty(meow4.Hash)

		is.Nil(cat1.Owner)
		is.Empty(cat1.Meows)
		is.Nil(cat2.Owner)
		is.Empty(cat2.Meows)
		is.Nil(cat3.Owner)
		is.Empty(cat3.Meows)

		{
			cats := []Cat{*cat1, *cat2, *cat3}
			err = sqlxx.Preload(ctx, driver, &cats, "Owner", "Meows")
			is.NoError(err)
			is.Len(cats, 3)
			is.Equal(cat1.ID, cats[0].ID)
			is.Equal(cat2.ID, cats[1].ID)
			is.Equal(cat3.ID, cats[2].ID)

			is.NotNil(cats[0].Owner)
			is.Equal(human1.ID, cats[0].Owner.ID)
			is.Equal(human1.Name, cats[0].Owner.Name)
			is.NotEmpty(cats[0].Meows)
			is.Len(cats[0].Meows, 3)
			is.Contains(cats[0].Meows, meow1)
			is.Contains(cats[0].Meows, meow2)
			is.Contains(cats[0].Meows, meow3)

			is.NotNil(cats[1].Owner)
			is.Equal(human2.ID, cats[1].Owner.ID)
			is.Equal(human2.Name, cats[1].Owner.Name)
			is.Empty(cats[1].Meows)

			is.Nil(cats[2].Owner)
			is.NotEmpty(cats[2].Meows)
			is.Len(cats[2].Meows, 1)
			is.Contains(cats[2].Meows, meow4)
		}
		{
			cats := []*Cat{cat1, cat2, cat3}
			err = sqlxx.Preload(ctx, driver, &cats, "Owner", "Meows")
			is.NoError(err)
			is.Len(cats, 3)
			is.Equal(cat1.ID, cats[0].ID)
			is.Equal(cat2.ID, cats[1].ID)
			is.Equal(cat3.ID, cats[2].ID)

			is.NotNil(cats[0].Owner)
			is.Equal(human1.ID, cats[0].Owner.ID)
			is.Equal(human1.Name, cats[0].Owner.Name)
			is.NotEmpty(cats[0].Meows)
			is.Len(cats[0].Meows, 3)
			is.Contains(cats[0].Meows, meow1)
			is.Contains(cats[0].Meows, meow2)
			is.Contains(cats[0].Meows, meow3)

			is.NotNil(cats[1].Owner)
			is.Equal(human2.ID, cats[1].Owner.ID)
			is.Equal(human2.Name, cats[1].Owner.Name)
			is.Empty(cats[1].Meows)

			is.Nil(cats[2].Owner)
			is.NotEmpty(cats[2].Meows)
			is.Len(cats[2].Meows, 1)
			is.Contains(cats[2].Meows, meow4)
		}

	})
}
