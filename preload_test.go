package sqlxx_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestPreload_CommonFailure(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		{
			value := Human{}
			err := sqlxx.Preload(ctx, nil, &value, "Cat")
			is.Error(err)
			is.Equal(sqlxx.ErrInvalidDriver, errors.Cause(err))
		}
		{
			value := 12
			err := sqlxx.Preload(ctx, driver, &value, "Cat")
			is.Error(err)
			is.Equal(sqlxx.ErrPreloadInvalidSchema, errors.Cause(err))
		}
		{
			value := Human{}
			err := sqlxx.Preload(ctx, driver, value, "Cat")
			is.Error(err)
			is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		}
		{
			value := Human{}
			err := sqlxx.Preload(ctx, driver, &value, "X", "Y", "Z")
			is.Error(err)
			is.Equal(sqlxx.ErrPreloadInvalidPath, errors.Cause(err))
		}
	})
}

func TestPreload_ExoChunk_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures ExoCloudFixtures) {
			is.Nil(fixtures.Chunks[0].Mode)
			is.Nil(fixtures.Chunks[1].Mode)
			is.Nil(fixtures.Chunks[2].Mode)
			is.Nil(fixtures.Chunks[3].Mode)
			is.Nil(fixtures.Chunks[4].Mode)
			is.Nil(fixtures.Chunks[5].Mode)
			is.Nil(fixtures.Chunks[6].Mode)
			is.Nil(fixtures.Chunks[0].Signature)
			is.Nil(fixtures.Chunks[1].Signature)
			is.Nil(fixtures.Chunks[2].Signature)
			is.Nil(fixtures.Chunks[3].Signature)
			is.Nil(fixtures.Chunks[4].Signature)
			is.Nil(fixtures.Chunks[5].Signature)
			is.Nil(fixtures.Chunks[6].Signature)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunk1 := fixtures.Chunks[0]

			err := sqlxx.Preload(ctx, driver, chunk1, "Mode")
			is.NoError(err)
			is.NotNil(chunk1.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk1.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk1.Mode.Mode)
			is.Nil(chunk1.Signature)

			err = sqlxx.Preload(ctx, driver, chunk1, "Signature")
			is.NoError(err)
			is.NotNil(chunk1.Signature)
			is.Equal(fixtures.Signatures[0].ID, chunk1.Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunk1.Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunk1.Signature.Bytes)

			chunk2 := fixtures.Chunks[1]

			err = sqlxx.Preload(ctx, driver, chunk2, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk2.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk2.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk2.Mode.Mode)
			is.NotNil(chunk2.Signature)
			is.Equal(fixtures.Signatures[1].ID, chunk2.Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunk2.Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunk2.Signature.Bytes)

			chunk4 := fixtures.Chunks[3]

			err = sqlxx.Preload(ctx, driver, chunk4, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk4.Mode)
			is.Equal(fixtures.Modes[1].ID, chunk4.Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunk4.Mode.Mode)
			is.Nil(chunk4.Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunk1 := fixtures.Chunks[0]

			err := sqlxx.Preload(ctx, driver, &chunk1, "Mode")
			is.NoError(err)
			is.NotNil(chunk1.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk1.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk1.Mode.Mode)
			is.Nil(chunk1.Signature)

			err = sqlxx.Preload(ctx, driver, &chunk1, "Signature")
			is.NoError(err)
			is.NotNil(chunk1.Signature)
			is.Equal(fixtures.Signatures[0].ID, chunk1.Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunk1.Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunk1.Signature.Bytes)

			chunk2 := fixtures.Chunks[1]

			err = sqlxx.Preload(ctx, driver, &chunk2, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk2.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk2.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk2.Mode.Mode)
			is.NotNil(chunk2.Signature)
			is.Equal(fixtures.Signatures[1].ID, chunk2.Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunk2.Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunk2.Signature.Bytes)

			chunk4 := fixtures.Chunks[3]

			err = sqlxx.Preload(ctx, driver, &chunk4, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk4.Mode)
			is.Equal(fixtures.Modes[1].ID, chunk4.Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunk4.Mode.Mode)
			is.Nil(chunk4.Signature)

		}
	})
}

func TestPreload_ExoChunk_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures ExoCloudFixtures) {
			is.Nil(fixtures.Chunks[0].Mode)
			is.Nil(fixtures.Chunks[1].Mode)
			is.Nil(fixtures.Chunks[2].Mode)
			is.Nil(fixtures.Chunks[3].Mode)
			is.Nil(fixtures.Chunks[4].Mode)
			is.Nil(fixtures.Chunks[5].Mode)
			is.Nil(fixtures.Chunks[6].Mode)
			is.Nil(fixtures.Chunks[0].Signature)
			is.Nil(fixtures.Chunks[1].Signature)
			is.Nil(fixtures.Chunks[2].Signature)
			is.Nil(fixtures.Chunks[3].Signature)
			is.Nil(fixtures.Chunks[4].Signature)
			is.Nil(fixtures.Chunks[5].Signature)
			is.Nil(fixtures.Chunks[6].Signature)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := []ExoChunk{
				*fixtures.Chunks[0],
				*fixtures.Chunks[1],
				*fixtures.Chunks[2],
				*fixtures.Chunks[3],
				*fixtures.Chunks[4],
				*fixtures.Chunks[5],
				*fixtures.Chunks[6],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 7)
			is.Equal(fixtures.Chunks[0].Hash, chunks[0].Hash)
			is.Equal(fixtures.Chunks[1].Hash, chunks[1].Hash)
			is.Equal(fixtures.Chunks[2].Hash, chunks[2].Hash)
			is.Equal(fixtures.Chunks[3].Hash, chunks[3].Hash)
			is.Equal(fixtures.Chunks[4].Hash, chunks[4].Hash)
			is.Equal(fixtures.Chunks[5].Hash, chunks[5].Hash)
			is.Equal(fixtures.Chunks[6].Hash, chunks[6].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[0].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[0].Mode.Mode)
			is.NotNil(chunks[0].Signature)
			is.Equal(fixtures.Signatures[0].ID, chunks[0].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunks[0].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunks[0].Signature.Bytes)

			is.NotNil(chunks[1].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(fixtures.Signatures[1].ID, chunks[1].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(fixtures.Signatures[2].ID, chunks[2].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(fixtures.Modes[1].ID, chunks[3].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunks[3].Mode.Mode)
			is.Nil(chunks[3].Signature)

			is.NotNil(chunks[4].Mode)
			is.Equal(fixtures.Modes[1].ID, chunks[4].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunks[4].Mode.Mode)
			is.Nil(chunks[4].Signature)

			is.NotNil(chunks[5].Mode)
			is.Equal(fixtures.Modes[2].ID, chunks[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunks[5].Mode.Mode)
			is.NotNil(chunks[5].Signature)
			is.Equal(fixtures.Signatures[3].ID, chunks[5].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunks[5].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunks[5].Signature.Bytes)

			is.NotNil(chunks[6].Mode)
			is.Equal(fixtures.Modes[3].ID, chunks[6].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunks[6].Mode.Mode)
			is.NotNil(chunks[6].Signature)
			is.Equal(fixtures.Signatures[4].ID, chunks[6].Signature.ID)
			is.Equal(fixtures.Signatures[4].ChunkID, chunks[6].Signature.ChunkID)
			is.Equal(fixtures.Signatures[4].Bytes, chunks[6].Signature.Bytes)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := []*ExoChunk{
				fixtures.Chunks[0],
				fixtures.Chunks[1],
				fixtures.Chunks[2],
				fixtures.Chunks[3],
				fixtures.Chunks[4],
				fixtures.Chunks[5],
				fixtures.Chunks[6],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 7)
			is.Equal(fixtures.Chunks[0].Hash, chunks[0].Hash)
			is.Equal(fixtures.Chunks[1].Hash, chunks[1].Hash)
			is.Equal(fixtures.Chunks[2].Hash, chunks[2].Hash)
			is.Equal(fixtures.Chunks[3].Hash, chunks[3].Hash)
			is.Equal(fixtures.Chunks[4].Hash, chunks[4].Hash)
			is.Equal(fixtures.Chunks[5].Hash, chunks[5].Hash)
			is.Equal(fixtures.Chunks[6].Hash, chunks[6].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[0].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[0].Mode.Mode)
			is.NotNil(chunks[0].Signature)
			is.Equal(fixtures.Signatures[0].ID, chunks[0].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunks[0].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunks[0].Signature.Bytes)

			is.NotNil(chunks[1].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(fixtures.Signatures[1].ID, chunks[1].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(fixtures.Signatures[2].ID, chunks[2].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(fixtures.Modes[1].ID, chunks[3].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunks[3].Mode.Mode)
			is.Nil(chunks[3].Signature)

			is.NotNil(chunks[4].Mode)
			is.Equal(fixtures.Modes[1].ID, chunks[4].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, chunks[4].Mode.Mode)
			is.Nil(chunks[4].Signature)

			is.NotNil(chunks[5].Mode)
			is.Equal(fixtures.Modes[2].ID, chunks[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunks[5].Mode.Mode)
			is.NotNil(chunks[5].Signature)
			is.Equal(fixtures.Signatures[3].ID, chunks[5].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunks[5].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunks[5].Signature.Bytes)

			is.NotNil(chunks[6].Mode)
			is.Equal(fixtures.Modes[3].ID, chunks[6].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunks[6].Mode.Mode)
			is.NotNil(chunks[6].Signature)
			is.Equal(fixtures.Signatures[4].ID, chunks[6].Signature.ID)
			is.Equal(fixtures.Signatures[4].ChunkID, chunks[6].Signature.ChunkID)
			is.Equal(fixtures.Signatures[4].Bytes, chunks[6].Signature.Bytes)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := &[]ExoChunk{
				*fixtures.Chunks[0],
				*fixtures.Chunks[1],
				*fixtures.Chunks[2],
				*fixtures.Chunks[3],
				*fixtures.Chunks[4],
				*fixtures.Chunks[5],
				*fixtures.Chunks[6],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len((*chunks), 7)
			is.Equal(fixtures.Chunks[0].Hash, (*chunks)[0].Hash)
			is.Equal(fixtures.Chunks[1].Hash, (*chunks)[1].Hash)
			is.Equal(fixtures.Chunks[2].Hash, (*chunks)[2].Hash)
			is.Equal(fixtures.Chunks[3].Hash, (*chunks)[3].Hash)
			is.Equal(fixtures.Chunks[4].Hash, (*chunks)[4].Hash)
			is.Equal(fixtures.Chunks[5].Hash, (*chunks)[5].Hash)
			is.Equal(fixtures.Chunks[6].Hash, (*chunks)[6].Hash)

			is.NotNil((*chunks)[0].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[0].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[0].Mode.Mode)
			is.NotNil((*chunks)[0].Signature)
			is.Equal(fixtures.Signatures[0].ID, (*chunks)[0].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, (*chunks)[0].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, (*chunks)[0].Signature.Bytes)

			is.NotNil((*chunks)[1].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[1].Mode.Mode)
			is.NotNil((*chunks)[1].Signature)
			is.Equal(fixtures.Signatures[1].ID, (*chunks)[1].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, (*chunks)[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, (*chunks)[1].Signature.Bytes)

			is.NotNil((*chunks)[2].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[2].Mode.Mode)
			is.NotNil((*chunks)[2].Signature)
			is.Equal(fixtures.Signatures[2].ID, (*chunks)[2].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, (*chunks)[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, (*chunks)[2].Signature.Bytes)

			is.NotNil((*chunks)[3].Mode)
			is.Equal(fixtures.Modes[1].ID, (*chunks)[3].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, (*chunks)[3].Mode.Mode)
			is.Nil((*chunks)[3].Signature)

			is.NotNil((*chunks)[4].Mode)
			is.Equal(fixtures.Modes[1].ID, (*chunks)[4].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, (*chunks)[4].Mode.Mode)
			is.Nil((*chunks)[4].Signature)

			is.NotNil((*chunks)[5].Mode)
			is.Equal(fixtures.Modes[2].ID, (*chunks)[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, (*chunks)[5].Mode.Mode)
			is.NotNil((*chunks)[5].Signature)
			is.Equal(fixtures.Signatures[3].ID, (*chunks)[5].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, (*chunks)[5].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, (*chunks)[5].Signature.Bytes)

			is.NotNil((*chunks)[6].Mode)
			is.Equal(fixtures.Modes[3].ID, (*chunks)[6].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, (*chunks)[6].Mode.Mode)
			is.NotNil((*chunks)[6].Signature)
			is.Equal(fixtures.Signatures[4].ID, (*chunks)[6].Signature.ID)
			is.Equal(fixtures.Signatures[4].ChunkID, (*chunks)[6].Signature.ChunkID)
			is.Equal(fixtures.Signatures[4].Bytes, (*chunks)[6].Signature.Bytes)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := &[]*ExoChunk{
				fixtures.Chunks[0],
				fixtures.Chunks[1],
				fixtures.Chunks[2],
				fixtures.Chunks[3],
				fixtures.Chunks[4],
				fixtures.Chunks[5],
				fixtures.Chunks[6],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len((*chunks), 7)
			is.Equal(fixtures.Chunks[0].Hash, (*chunks)[0].Hash)
			is.Equal(fixtures.Chunks[1].Hash, (*chunks)[1].Hash)
			is.Equal(fixtures.Chunks[2].Hash, (*chunks)[2].Hash)
			is.Equal(fixtures.Chunks[3].Hash, (*chunks)[3].Hash)
			is.Equal(fixtures.Chunks[4].Hash, (*chunks)[4].Hash)
			is.Equal(fixtures.Chunks[5].Hash, (*chunks)[5].Hash)
			is.Equal(fixtures.Chunks[6].Hash, (*chunks)[6].Hash)

			is.NotNil((*chunks)[0].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[0].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[0].Mode.Mode)
			is.NotNil((*chunks)[0].Signature)
			is.Equal(fixtures.Signatures[0].ID, (*chunks)[0].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, (*chunks)[0].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, (*chunks)[0].Signature.Bytes)

			is.NotNil((*chunks)[1].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[1].Mode.Mode)
			is.NotNil((*chunks)[1].Signature)
			is.Equal(fixtures.Signatures[1].ID, (*chunks)[1].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, (*chunks)[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, (*chunks)[1].Signature.Bytes)

			is.NotNil((*chunks)[2].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[2].Mode.Mode)
			is.NotNil((*chunks)[2].Signature)
			is.Equal(fixtures.Signatures[2].ID, (*chunks)[2].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, (*chunks)[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, (*chunks)[2].Signature.Bytes)

			is.NotNil((*chunks)[3].Mode)
			is.Equal(fixtures.Modes[1].ID, (*chunks)[3].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, (*chunks)[3].Mode.Mode)
			is.Nil((*chunks)[3].Signature)

			is.NotNil((*chunks)[4].Mode)
			is.Equal(fixtures.Modes[1].ID, (*chunks)[4].Mode.ID)
			is.Equal(fixtures.Modes[1].Mode, (*chunks)[4].Mode.Mode)
			is.Nil((*chunks)[4].Signature)

			is.NotNil((*chunks)[5].Mode)
			is.Equal(fixtures.Modes[2].ID, (*chunks)[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, (*chunks)[5].Mode.Mode)
			is.NotNil((*chunks)[5].Signature)
			is.Equal(fixtures.Signatures[3].ID, (*chunks)[5].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, (*chunks)[5].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, (*chunks)[5].Signature.Bytes)

			is.NotNil((*chunks)[6].Mode)
			is.Equal(fixtures.Modes[3].ID, (*chunks)[6].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, (*chunks)[6].Mode.Mode)
			is.NotNil((*chunks)[6].Signature)
			is.Equal(fixtures.Signatures[4].ID, (*chunks)[6].Signature.ID)
			is.Equal(fixtures.Signatures[4].ChunkID, (*chunks)[6].Signature.ChunkID)
			is.Equal(fixtures.Signatures[4].Bytes, (*chunks)[6].Signature.Bytes)

		}
		{

			chunks := []ExoChunk{}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 0)

		}
	})
}

func TestPreload_Owl_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckOwlFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Owls[0].Group)
			is.Nil(fixtures.Owls[1].Group)
			is.Nil(fixtures.Owls[2].Group)
			is.Nil(fixtures.Owls[3].Group)
			is.Nil(fixtures.Owls[4].Group)
			is.Nil(fixtures.Owls[5].Group)
			is.Empty(fixtures.Owls[0].Group)
			is.Empty(fixtures.Owls[1].Packages)
			is.Empty(fixtures.Owls[2].Packages)
			is.Empty(fixtures.Owls[3].Packages)
			is.Empty(fixtures.Owls[4].Packages)
			is.Empty(fixtures.Owls[5].Packages)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owl1 := fixtures.Owls[0]

			err := sqlxx.Preload(ctx, driver, owl1, "Group")
			is.NoError(err)
			is.NotNil(owl1.Group)
			is.Equal(fixtures.Groups[0].ID, owl1.Group.ID)
			is.Equal(fixtures.Groups[0].Name, owl1.Group.Name)
			is.Empty(owl1.Packages)

			err = sqlxx.Preload(ctx, driver, owl1, "Packages")
			is.NoError(err)
			is.NotEmpty(owl1.Packages)
			is.Len(owl1.Packages, 2)
			is.Contains(owl1.Packages, *fixtures.Packages[0])
			is.Contains(owl1.Packages, *fixtures.Packages[1])

			owl2 := fixtures.Owls[1]

			err = sqlxx.Preload(ctx, driver, owl2, "Group", "Packages")
			is.NoError(err)
			is.Nil(owl2.Group)
			is.NotEmpty(owl2.Packages)
			is.Len(owl2.Packages, 1)
			is.Contains(owl2.Packages, *fixtures.Packages[3])

			owl3 := fixtures.Owls[2]

			err = sqlxx.Preload(ctx, driver, owl3, "Group", "Packages")
			is.NoError(err)
			is.NotNil(owl3.Group)
			is.Equal(fixtures.Groups[0].ID, owl3.Group.ID)
			is.Equal(fixtures.Groups[0].Name, owl3.Group.Name)
			is.NotEmpty(owl3.Packages)
			is.Len(owl3.Packages, 5)
			is.Contains(owl3.Packages, *fixtures.Packages[4])
			is.Contains(owl3.Packages, *fixtures.Packages[5])
			is.Contains(owl3.Packages, *fixtures.Packages[6])
			is.Contains(owl3.Packages, *fixtures.Packages[7])
			is.Contains(owl3.Packages, *fixtures.Packages[8])

			owl5 := fixtures.Owls[4]

			err = sqlxx.Preload(ctx, driver, owl5, "Group", "Packages")
			is.NoError(err)
			is.NotNil(owl5.Group)
			is.Equal(fixtures.Groups[2].ID, owl5.Group.ID)
			is.Equal(fixtures.Groups[2].Name, owl5.Group.Name)
			is.Empty(owl5.Packages)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owl1 := fixtures.Owls[0]

			err := sqlxx.Preload(ctx, driver, &owl1, "Group", "Packages")
			is.NoError(err)
			is.NotNil(owl1.Group)
			is.Equal(fixtures.Groups[0].ID, owl1.Group.ID)
			is.Equal(fixtures.Groups[0].Name, owl1.Group.Name)
			is.NotEmpty(owl1.Packages)
			is.Len(owl1.Packages, 2)
			is.Contains(owl1.Packages, *fixtures.Packages[0])
			is.Contains(owl1.Packages, *fixtures.Packages[1])

			owl2 := fixtures.Owls[1]

			err = sqlxx.Preload(ctx, driver, &owl2, "Packages")
			is.NoError(err)
			is.Nil(owl2.Group)
			is.NotEmpty(owl2.Packages)
			is.Len(owl2.Packages, 1)
			is.Contains(owl2.Packages, *fixtures.Packages[3])

			owl5 := fixtures.Owls[4]

			err = sqlxx.Preload(ctx, driver, &owl5, "Group", "Packages")
			is.NoError(err)
			is.NotNil(owl5.Group)
			is.Equal(fixtures.Groups[2].ID, owl5.Group.ID)
			is.Equal(fixtures.Groups[2].Name, owl5.Group.Name)
			is.Empty(owl5.Packages)

		}
	})
}

func TestPreload_Owl_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckOwlFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Owls[0].Group)
			is.Nil(fixtures.Owls[1].Group)
			is.Nil(fixtures.Owls[2].Group)
			is.Nil(fixtures.Owls[3].Group)
			is.Nil(fixtures.Owls[4].Group)
			is.Nil(fixtures.Owls[5].Group)
			is.Empty(fixtures.Owls[0].Group)
			is.Empty(fixtures.Owls[1].Packages)
			is.Empty(fixtures.Owls[2].Packages)
			is.Empty(fixtures.Owls[3].Packages)
			is.Empty(fixtures.Owls[4].Packages)
			is.Empty(fixtures.Owls[5].Packages)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owls := []Owl{
				*fixtures.Owls[0],
				*fixtures.Owls[1],
				*fixtures.Owls[2],
				*fixtures.Owls[3],
				*fixtures.Owls[4],
				*fixtures.Owls[5],
			}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len(owls, 6)
			is.Equal(fixtures.Owls[0].ID, owls[0].ID)
			is.Equal(fixtures.Owls[1].ID, owls[1].ID)
			is.Equal(fixtures.Owls[2].ID, owls[2].ID)
			is.Equal(fixtures.Owls[3].ID, owls[3].ID)
			is.Equal(fixtures.Owls[4].ID, owls[4].ID)
			is.Equal(fixtures.Owls[5].ID, owls[5].ID)

			is.NotNil(owls[0].Group)
			is.Equal(fixtures.Groups[0].ID, owls[0].Group.ID)
			is.Equal(fixtures.Groups[0].Name, owls[0].Group.Name)
			is.NotEmpty(owls[0].Packages)
			is.Len(owls[0].Packages, 2)
			is.Contains(owls[0].Packages, *fixtures.Packages[0])
			is.Contains(owls[0].Packages, *fixtures.Packages[1])

			is.Nil(owls[1].Group)
			is.NotEmpty(owls[1].Packages)
			is.Len(owls[1].Packages, 1)
			is.Contains(owls[1].Packages, *fixtures.Packages[3])

			is.NotNil(owls[2].Group)
			is.Equal(fixtures.Groups[0].ID, owls[2].Group.ID)
			is.Equal(fixtures.Groups[0].Name, owls[2].Group.Name)
			is.NotEmpty(owls[2].Packages)
			is.Len(owls[2].Packages, 5)
			is.Contains(owls[2].Packages, *fixtures.Packages[4])
			is.Contains(owls[2].Packages, *fixtures.Packages[5])
			is.Contains(owls[2].Packages, *fixtures.Packages[6])
			is.Contains(owls[2].Packages, *fixtures.Packages[7])
			is.Contains(owls[2].Packages, *fixtures.Packages[8])

			is.NotNil(owls[3].Group)
			is.Equal(fixtures.Groups[1].ID, owls[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, owls[3].Group.Name)
			is.NotEmpty(owls[3].Packages)
			is.Len(owls[3].Packages, 4)
			is.Contains(owls[3].Packages, *fixtures.Packages[9])
			is.Contains(owls[3].Packages, *fixtures.Packages[10])
			is.Contains(owls[3].Packages, *fixtures.Packages[11])
			is.Contains(owls[3].Packages, *fixtures.Packages[12])

			is.NotNil(owls[4].Group)
			is.Equal(fixtures.Groups[2].ID, owls[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, owls[4].Group.Name)
			is.Empty(owls[4].Packages)

			is.NotNil(owls[5].Group)
			is.Equal(fixtures.Groups[3].ID, owls[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, owls[5].Group.Name)
			is.NotEmpty(owls[5].Packages)
			is.Len(owls[5].Packages, 3)
			is.Contains(owls[5].Packages, *fixtures.Packages[13])
			is.Contains(owls[5].Packages, *fixtures.Packages[14])
			is.Contains(owls[5].Packages, *fixtures.Packages[15])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owls := []*Owl{
				fixtures.Owls[0],
				fixtures.Owls[1],
				fixtures.Owls[2],
				fixtures.Owls[3],
				fixtures.Owls[4],
				fixtures.Owls[5],
			}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len(owls, 6)
			is.Equal(fixtures.Owls[0].ID, owls[0].ID)
			is.Equal(fixtures.Owls[1].ID, owls[1].ID)
			is.Equal(fixtures.Owls[2].ID, owls[2].ID)
			is.Equal(fixtures.Owls[3].ID, owls[3].ID)
			is.Equal(fixtures.Owls[4].ID, owls[4].ID)
			is.Equal(fixtures.Owls[5].ID, owls[5].ID)

			is.NotNil(owls[0].Group)
			is.Equal(fixtures.Groups[0].ID, owls[0].Group.ID)
			is.Equal(fixtures.Groups[0].Name, owls[0].Group.Name)
			is.NotEmpty(owls[0].Packages)
			is.Len(owls[0].Packages, 2)
			is.Contains(owls[0].Packages, *fixtures.Packages[0])
			is.Contains(owls[0].Packages, *fixtures.Packages[1])

			is.Nil(owls[1].Group)
			is.NotEmpty(owls[1].Packages)
			is.Len(owls[1].Packages, 1)
			is.Contains(owls[1].Packages, *fixtures.Packages[3])

			is.NotNil(owls[2].Group)
			is.Equal(fixtures.Groups[0].ID, owls[2].Group.ID)
			is.Equal(fixtures.Groups[0].Name, owls[2].Group.Name)
			is.NotEmpty(owls[2].Packages)
			is.Len(owls[2].Packages, 5)
			is.Contains(owls[2].Packages, *fixtures.Packages[4])
			is.Contains(owls[2].Packages, *fixtures.Packages[5])
			is.Contains(owls[2].Packages, *fixtures.Packages[6])
			is.Contains(owls[2].Packages, *fixtures.Packages[7])
			is.Contains(owls[2].Packages, *fixtures.Packages[8])

			is.NotNil(owls[3].Group)
			is.Equal(fixtures.Groups[1].ID, owls[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, owls[3].Group.Name)
			is.NotEmpty(owls[3].Packages)
			is.Len(owls[3].Packages, 4)
			is.Contains(owls[3].Packages, *fixtures.Packages[9])
			is.Contains(owls[3].Packages, *fixtures.Packages[10])
			is.Contains(owls[3].Packages, *fixtures.Packages[11])
			is.Contains(owls[3].Packages, *fixtures.Packages[12])

			is.NotNil(owls[4].Group)
			is.Equal(fixtures.Groups[2].ID, owls[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, owls[4].Group.Name)
			is.Empty(owls[4].Packages)

			is.NotNil(owls[5].Group)
			is.Equal(fixtures.Groups[3].ID, owls[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, owls[5].Group.Name)
			is.NotEmpty(owls[5].Packages)
			is.Len(owls[5].Packages, 3)
			is.Contains(owls[5].Packages, *fixtures.Packages[13])
			is.Contains(owls[5].Packages, *fixtures.Packages[14])
			is.Contains(owls[5].Packages, *fixtures.Packages[15])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owls := &[]Owl{
				*fixtures.Owls[0],
				*fixtures.Owls[1],
				*fixtures.Owls[2],
				*fixtures.Owls[3],
				*fixtures.Owls[4],
				*fixtures.Owls[5],
			}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len((*owls), 6)
			is.Equal(fixtures.Owls[0].ID, (*owls)[0].ID)
			is.Equal(fixtures.Owls[1].ID, (*owls)[1].ID)
			is.Equal(fixtures.Owls[2].ID, (*owls)[2].ID)
			is.Equal(fixtures.Owls[3].ID, (*owls)[3].ID)
			is.Equal(fixtures.Owls[4].ID, (*owls)[4].ID)
			is.Equal(fixtures.Owls[5].ID, (*owls)[5].ID)

			is.NotNil((*owls)[0].Group)
			is.Equal(fixtures.Groups[0].ID, (*owls)[0].Group.ID)
			is.Equal(fixtures.Groups[0].Name, (*owls)[0].Group.Name)
			is.NotEmpty((*owls)[0].Packages)
			is.Len((*owls)[0].Packages, 2)
			is.Contains((*owls)[0].Packages, *fixtures.Packages[0])
			is.Contains((*owls)[0].Packages, *fixtures.Packages[1])

			is.Nil((*owls)[1].Group)
			is.NotEmpty((*owls)[1].Packages)
			is.Len((*owls)[1].Packages, 1)
			is.Contains((*owls)[1].Packages, *fixtures.Packages[3])

			is.NotNil((*owls)[2].Group)
			is.Equal(fixtures.Groups[0].ID, (*owls)[2].Group.ID)
			is.Equal(fixtures.Groups[0].Name, (*owls)[2].Group.Name)
			is.NotEmpty((*owls)[2].Packages)
			is.Len((*owls)[2].Packages, 5)
			is.Contains((*owls)[2].Packages, *fixtures.Packages[4])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[5])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[6])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[7])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[8])

			is.NotNil((*owls)[3].Group)
			is.Equal(fixtures.Groups[1].ID, (*owls)[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, (*owls)[3].Group.Name)
			is.NotEmpty((*owls)[3].Packages)
			is.Len((*owls)[3].Packages, 4)
			is.Contains((*owls)[3].Packages, *fixtures.Packages[9])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[10])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[11])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[12])

			is.NotNil((*owls)[4].Group)
			is.Equal(fixtures.Groups[2].ID, (*owls)[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, (*owls)[4].Group.Name)
			is.Empty((*owls)[4].Packages)

			is.NotNil((*owls)[5].Group)
			is.Equal(fixtures.Groups[3].ID, (*owls)[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, (*owls)[5].Group.Name)
			is.NotEmpty((*owls)[5].Packages)
			is.Len((*owls)[5].Packages, 3)
			is.Contains((*owls)[5].Packages, *fixtures.Packages[13])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[14])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[15])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owls := &[]*Owl{
				fixtures.Owls[0],
				fixtures.Owls[1],
				fixtures.Owls[2],
				fixtures.Owls[3],
				fixtures.Owls[4],
				fixtures.Owls[5],
			}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len((*owls), 6)
			is.Equal(fixtures.Owls[0].ID, (*owls)[0].ID)
			is.Equal(fixtures.Owls[1].ID, (*owls)[1].ID)
			is.Equal(fixtures.Owls[2].ID, (*owls)[2].ID)
			is.Equal(fixtures.Owls[3].ID, (*owls)[3].ID)
			is.Equal(fixtures.Owls[4].ID, (*owls)[4].ID)
			is.Equal(fixtures.Owls[5].ID, (*owls)[5].ID)

			is.NotNil((*owls)[0].Group)
			is.Equal(fixtures.Groups[0].ID, (*owls)[0].Group.ID)
			is.Equal(fixtures.Groups[0].Name, (*owls)[0].Group.Name)
			is.NotEmpty((*owls)[0].Packages)
			is.Len((*owls)[0].Packages, 2)
			is.Contains((*owls)[0].Packages, *fixtures.Packages[0])
			is.Contains((*owls)[0].Packages, *fixtures.Packages[1])

			is.Nil((*owls)[1].Group)
			is.NotEmpty((*owls)[1].Packages)
			is.Len((*owls)[1].Packages, 1)
			is.Contains((*owls)[1].Packages, *fixtures.Packages[3])

			is.NotNil((*owls)[2].Group)
			is.Equal(fixtures.Groups[0].ID, (*owls)[2].Group.ID)
			is.Equal(fixtures.Groups[0].Name, (*owls)[2].Group.Name)
			is.NotEmpty((*owls)[2].Packages)
			is.Len((*owls)[2].Packages, 5)
			is.Contains((*owls)[2].Packages, *fixtures.Packages[4])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[5])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[6])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[7])
			is.Contains((*owls)[2].Packages, *fixtures.Packages[8])

			is.NotNil((*owls)[3].Group)
			is.Equal(fixtures.Groups[1].ID, (*owls)[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, (*owls)[3].Group.Name)
			is.NotEmpty((*owls)[3].Packages)
			is.Len((*owls)[3].Packages, 4)
			is.Contains((*owls)[3].Packages, *fixtures.Packages[9])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[10])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[11])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[12])

			is.NotNil((*owls)[4].Group)
			is.Equal(fixtures.Groups[2].ID, (*owls)[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, (*owls)[4].Group.Name)
			is.Empty((*owls)[4].Packages)

			is.NotNil((*owls)[5].Group)
			is.Equal(fixtures.Groups[3].ID, (*owls)[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, (*owls)[5].Group.Name)
			is.NotEmpty((*owls)[5].Packages)
			is.Len((*owls)[5].Packages, 3)
			is.Contains((*owls)[5].Packages, *fixtures.Packages[13])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[14])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[15])

		}
		{

			owls := []Owl{}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len(owls, 0)

		}
	})
}

func TestPreload_Cat_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckCatFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Cats[0].Feeder)
			is.Nil(fixtures.Cats[1].Feeder)
			is.Nil(fixtures.Cats[2].Feeder)
			is.Nil(fixtures.Cats[3].Feeder)
			is.Nil(fixtures.Cats[4].Feeder)
			is.Nil(fixtures.Cats[5].Feeder)
			is.Nil(fixtures.Cats[6].Feeder)
			is.Nil(fixtures.Cats[7].Feeder)
			is.Empty(fixtures.Cats[0].Meows)
			is.Empty(fixtures.Cats[1].Meows)
			is.Empty(fixtures.Cats[2].Meows)
			is.Empty(fixtures.Cats[3].Meows)
			is.Empty(fixtures.Cats[4].Meows)
			is.Empty(fixtures.Cats[5].Meows)
			is.Empty(fixtures.Cats[6].Meows)
			is.Empty(fixtures.Cats[7].Meows)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cat1 := fixtures.Cats[0]

			err := sqlxx.Preload(ctx, driver, cat1, "Feeder")
			is.NoError(err)
			is.NotNil(cat1.Feeder)
			is.Equal(fixtures.Humans[0].ID, cat1.Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, cat1.Feeder.Name)
			is.Empty(cat1.Meows)

			err = sqlxx.Preload(ctx, driver, cat1, "Meows")
			is.NoError(err)
			is.NotEmpty(cat1.Meows)
			is.Len(cat1.Meows, 3)
			is.Contains(cat1.Meows, fixtures.Meows[0])
			is.Contains(cat1.Meows, fixtures.Meows[1])
			is.Contains(cat1.Meows, fixtures.Meows[2])

			cat2 := fixtures.Cats[1]

			err = sqlxx.Preload(ctx, driver, cat2, "Feeder", "Meows")
			is.NoError(err)
			is.NotNil(cat2.Feeder)
			is.Equal(fixtures.Humans[1].ID, cat2.Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, cat2.Feeder.Name)
			is.Empty(cat2.Meows)

			cat3 := fixtures.Cats[2]

			err = sqlxx.Preload(ctx, driver, cat3, "Feeder", "Meows")
			is.NoError(err)
			is.Nil(cat3.Feeder)
			is.NotEmpty(cat3.Meows)
			is.Len(cat3.Meows, 1)
			is.Contains(cat3.Meows, fixtures.Meows[3])

			cat6 := fixtures.Cats[5]

			err = sqlxx.Preload(ctx, driver, cat6, "Feeder", "Meows")
			is.NoError(err)
			is.Nil(cat6.Feeder)
			is.Empty(cat6.Meows)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cat1 := fixtures.Cats[0]

			err := sqlxx.Preload(ctx, driver, &cat1, "Feeder")
			is.NoError(err)
			is.NotNil(cat1.Feeder)
			is.Equal(fixtures.Humans[0].ID, cat1.Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, cat1.Feeder.Name)
			is.Empty(cat1.Meows)

			err = sqlxx.Preload(ctx, driver, &cat1, "Meows")
			is.NoError(err)
			is.NotEmpty(cat1.Meows)
			is.Len(cat1.Meows, 3)
			is.Contains(cat1.Meows, fixtures.Meows[0])
			is.Contains(cat1.Meows, fixtures.Meows[1])
			is.Contains(cat1.Meows, fixtures.Meows[2])

			cat2 := fixtures.Cats[1]

			err = sqlxx.Preload(ctx, driver, &cat2, "Feeder", "Meows")
			is.NoError(err)
			is.NotNil(cat2.Feeder)
			is.Equal(fixtures.Humans[1].ID, cat2.Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, cat2.Feeder.Name)
			is.Empty(cat2.Meows)

			cat3 := fixtures.Cats[2]

			err = sqlxx.Preload(ctx, driver, &cat3, "Feeder", "Meows")
			is.NoError(err)
			is.Nil(cat3.Feeder)
			is.NotEmpty(cat3.Meows)
			is.Len(cat3.Meows, 1)
			is.Contains(cat3.Meows, fixtures.Meows[3])

			cat6 := fixtures.Cats[5]

			err = sqlxx.Preload(ctx, driver, &cat6, "Feeder", "Meows")
			is.NoError(err)
			is.Nil(cat6.Feeder)
			is.Empty(cat6.Meows)

		}
	})
}

func TestPreload_Cat_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckCatFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Cats[0].Feeder)
			is.Nil(fixtures.Cats[1].Feeder)
			is.Nil(fixtures.Cats[2].Feeder)
			is.Nil(fixtures.Cats[3].Feeder)
			is.Nil(fixtures.Cats[4].Feeder)
			is.Nil(fixtures.Cats[5].Feeder)
			is.Nil(fixtures.Cats[6].Feeder)
			is.Nil(fixtures.Cats[7].Feeder)
			is.Empty(fixtures.Cats[0].Meows)
			is.Empty(fixtures.Cats[1].Meows)
			is.Empty(fixtures.Cats[2].Meows)
			is.Empty(fixtures.Cats[3].Meows)
			is.Empty(fixtures.Cats[4].Meows)
			is.Empty(fixtures.Cats[5].Meows)
			is.Empty(fixtures.Cats[6].Meows)
			is.Empty(fixtures.Cats[7].Meows)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cats := []Cat{
				*fixtures.Cats[0],
				*fixtures.Cats[1],
				*fixtures.Cats[2],
				*fixtures.Cats[3],
				*fixtures.Cats[4],
				*fixtures.Cats[5],
				*fixtures.Cats[6],
				*fixtures.Cats[7],
			}

			err := sqlxx.Preload(ctx, driver, &cats, "Feeder", "Meows")
			is.NoError(err)
			is.Len(cats, 8)
			is.Equal(fixtures.Cats[0].ID, cats[0].ID)
			is.Equal(fixtures.Cats[1].ID, cats[1].ID)
			is.Equal(fixtures.Cats[2].ID, cats[2].ID)
			is.Equal(fixtures.Cats[3].ID, cats[3].ID)
			is.Equal(fixtures.Cats[4].ID, cats[4].ID)
			is.Equal(fixtures.Cats[5].ID, cats[5].ID)
			is.Equal(fixtures.Cats[6].ID, cats[6].ID)
			is.Equal(fixtures.Cats[7].ID, cats[7].ID)

			is.NotNil(cats[0].Feeder)
			is.Equal(fixtures.Humans[0].ID, cats[0].Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, cats[0].Feeder.Name)
			is.NotEmpty(cats[0].Meows)
			is.Len(cats[0].Meows, 3)
			is.Contains(cats[0].Meows, fixtures.Meows[0])
			is.Contains(cats[0].Meows, fixtures.Meows[1])
			is.Contains(cats[0].Meows, fixtures.Meows[2])

			is.NotNil(cats[1].Feeder)
			is.Equal(fixtures.Humans[1].ID, cats[1].Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, cats[1].Feeder.Name)
			is.Empty(cats[1].Meows)

			is.Nil(cats[2].Feeder)
			is.NotEmpty(cats[2].Meows)
			is.Len(cats[2].Meows, 1)
			is.Contains(cats[2].Meows, fixtures.Meows[3])

			is.NotNil(cats[3].Feeder)
			is.Equal(fixtures.Humans[3].ID, cats[3].Feeder.ID)
			is.Equal(fixtures.Humans[3].Name, cats[3].Feeder.Name)
			is.NotEmpty(cats[3].Meows)
			is.Len(cats[3].Meows, 3)
			is.Contains(cats[3].Meows, fixtures.Meows[4])
			is.Contains(cats[3].Meows, fixtures.Meows[5])
			is.Contains(cats[3].Meows, fixtures.Meows[6])

			is.Nil(cats[4].Feeder)
			is.NotEmpty(cats[4].Meows)
			is.Len(cats[4].Meows, 1)
			is.Contains(cats[4].Meows, fixtures.Meows[7])

			is.Nil(cats[5].Feeder)
			is.Empty(cats[5].Meows)

			is.NotNil(cats[6].Feeder)
			is.Equal(fixtures.Humans[6].ID, cats[6].Feeder.ID)
			is.Equal(fixtures.Humans[6].Name, cats[6].Feeder.Name)
			is.NotEmpty(cats[6].Meows)
			is.Len(cats[6].Meows, 2)
			is.Contains(cats[6].Meows, fixtures.Meows[8])
			is.Contains(cats[6].Meows, fixtures.Meows[9])

			is.Nil(cats[7].Feeder)
			is.NotEmpty(cats[7].Meows)
			is.Len(cats[7].Meows, 5)
			is.Contains(cats[7].Meows, fixtures.Meows[10])
			is.Contains(cats[7].Meows, fixtures.Meows[11])
			is.Contains(cats[7].Meows, fixtures.Meows[12])
			is.Contains(cats[7].Meows, fixtures.Meows[13])
			is.Contains(cats[7].Meows, fixtures.Meows[14])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cats := []*Cat{
				fixtures.Cats[0],
				fixtures.Cats[1],
				fixtures.Cats[2],
				fixtures.Cats[3],
				fixtures.Cats[4],
				fixtures.Cats[5],
				fixtures.Cats[6],
				fixtures.Cats[7],
			}

			err := sqlxx.Preload(ctx, driver, &cats, "Feeder", "Meows")
			is.NoError(err)
			is.Len(cats, 8)
			is.Equal(fixtures.Cats[0].ID, cats[0].ID)
			is.Equal(fixtures.Cats[1].ID, cats[1].ID)
			is.Equal(fixtures.Cats[2].ID, cats[2].ID)
			is.Equal(fixtures.Cats[3].ID, cats[3].ID)
			is.Equal(fixtures.Cats[4].ID, cats[4].ID)
			is.Equal(fixtures.Cats[5].ID, cats[5].ID)
			is.Equal(fixtures.Cats[6].ID, cats[6].ID)
			is.Equal(fixtures.Cats[7].ID, cats[7].ID)

			is.NotNil(cats[0].Feeder)
			is.Equal(fixtures.Humans[0].ID, cats[0].Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, cats[0].Feeder.Name)
			is.NotEmpty(cats[0].Meows)
			is.Len(cats[0].Meows, 3)
			is.Contains(cats[0].Meows, fixtures.Meows[0])
			is.Contains(cats[0].Meows, fixtures.Meows[1])
			is.Contains(cats[0].Meows, fixtures.Meows[2])

			is.NotNil(cats[1].Feeder)
			is.Equal(fixtures.Humans[1].ID, cats[1].Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, cats[1].Feeder.Name)
			is.Empty(cats[1].Meows)

			is.Nil(cats[2].Feeder)
			is.NotEmpty(cats[2].Meows)
			is.Len(cats[2].Meows, 1)
			is.Contains(cats[2].Meows, fixtures.Meows[3])

			is.NotNil(cats[3].Feeder)
			is.Equal(fixtures.Humans[3].ID, cats[3].Feeder.ID)
			is.Equal(fixtures.Humans[3].Name, cats[3].Feeder.Name)
			is.NotEmpty(cats[3].Meows)
			is.Len(cats[3].Meows, 3)
			is.Contains(cats[3].Meows, fixtures.Meows[4])
			is.Contains(cats[3].Meows, fixtures.Meows[5])
			is.Contains(cats[3].Meows, fixtures.Meows[6])

			is.Nil(cats[4].Feeder)
			is.NotEmpty(cats[4].Meows)
			is.Len(cats[4].Meows, 1)
			is.Contains(cats[4].Meows, fixtures.Meows[7])

			is.Nil(cats[5].Feeder)
			is.Empty(cats[5].Meows)

			is.NotNil(cats[6].Feeder)
			is.Equal(fixtures.Humans[6].ID, cats[6].Feeder.ID)
			is.Equal(fixtures.Humans[6].Name, cats[6].Feeder.Name)
			is.NotEmpty(cats[6].Meows)
			is.Len(cats[6].Meows, 2)
			is.Contains(cats[6].Meows, fixtures.Meows[8])
			is.Contains(cats[6].Meows, fixtures.Meows[9])

			is.Nil(cats[7].Feeder)
			is.NotEmpty(cats[7].Meows)
			is.Len(cats[7].Meows, 5)
			is.Contains(cats[7].Meows, fixtures.Meows[10])
			is.Contains(cats[7].Meows, fixtures.Meows[11])
			is.Contains(cats[7].Meows, fixtures.Meows[12])
			is.Contains(cats[7].Meows, fixtures.Meows[13])
			is.Contains(cats[7].Meows, fixtures.Meows[14])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cats := &[]Cat{
				*fixtures.Cats[0],
				*fixtures.Cats[1],
				*fixtures.Cats[2],
				*fixtures.Cats[3],
				*fixtures.Cats[4],
				*fixtures.Cats[5],
				*fixtures.Cats[6],
				*fixtures.Cats[7],
			}

			err := sqlxx.Preload(ctx, driver, &cats, "Feeder", "Meows")
			is.NoError(err)
			is.Len((*cats), 8)
			is.Equal(fixtures.Cats[0].ID, (*cats)[0].ID)
			is.Equal(fixtures.Cats[1].ID, (*cats)[1].ID)
			is.Equal(fixtures.Cats[2].ID, (*cats)[2].ID)
			is.Equal(fixtures.Cats[3].ID, (*cats)[3].ID)
			is.Equal(fixtures.Cats[4].ID, (*cats)[4].ID)
			is.Equal(fixtures.Cats[5].ID, (*cats)[5].ID)
			is.Equal(fixtures.Cats[6].ID, (*cats)[6].ID)
			is.Equal(fixtures.Cats[7].ID, (*cats)[7].ID)

			is.NotNil((*cats)[0].Feeder)
			is.Equal(fixtures.Humans[0].ID, (*cats)[0].Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, (*cats)[0].Feeder.Name)
			is.NotEmpty((*cats)[0].Meows)
			is.Len((*cats)[0].Meows, 3)
			is.Contains((*cats)[0].Meows, fixtures.Meows[0])
			is.Contains((*cats)[0].Meows, fixtures.Meows[1])
			is.Contains((*cats)[0].Meows, fixtures.Meows[2])

			is.NotNil((*cats)[1].Feeder)
			is.Equal(fixtures.Humans[1].ID, (*cats)[1].Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, (*cats)[1].Feeder.Name)
			is.Empty((*cats)[1].Meows)

			is.Nil((*cats)[2].Feeder)
			is.NotEmpty((*cats)[2].Meows)
			is.Len((*cats)[2].Meows, 1)
			is.Contains((*cats)[2].Meows, fixtures.Meows[3])

			is.NotNil((*cats)[3].Feeder)
			is.Equal(fixtures.Humans[3].ID, (*cats)[3].Feeder.ID)
			is.Equal(fixtures.Humans[3].Name, (*cats)[3].Feeder.Name)
			is.NotEmpty((*cats)[3].Meows)
			is.Len((*cats)[3].Meows, 3)
			is.Contains((*cats)[3].Meows, fixtures.Meows[4])
			is.Contains((*cats)[3].Meows, fixtures.Meows[5])
			is.Contains((*cats)[3].Meows, fixtures.Meows[6])

			is.Nil((*cats)[4].Feeder)
			is.NotEmpty((*cats)[4].Meows)
			is.Len((*cats)[4].Meows, 1)
			is.Contains((*cats)[4].Meows, fixtures.Meows[7])

			is.Nil((*cats)[5].Feeder)
			is.Empty((*cats)[5].Meows)

			is.NotNil((*cats)[6].Feeder)
			is.Equal(fixtures.Humans[6].ID, (*cats)[6].Feeder.ID)
			is.Equal(fixtures.Humans[6].Name, (*cats)[6].Feeder.Name)
			is.NotEmpty((*cats)[6].Meows)
			is.Len((*cats)[6].Meows, 2)
			is.Contains((*cats)[6].Meows, fixtures.Meows[8])
			is.Contains((*cats)[6].Meows, fixtures.Meows[9])

			is.Nil((*cats)[7].Feeder)
			is.NotEmpty((*cats)[7].Meows)
			is.Len((*cats)[7].Meows, 5)
			is.Contains((*cats)[7].Meows, fixtures.Meows[10])
			is.Contains((*cats)[7].Meows, fixtures.Meows[11])
			is.Contains((*cats)[7].Meows, fixtures.Meows[12])
			is.Contains((*cats)[7].Meows, fixtures.Meows[13])
			is.Contains((*cats)[7].Meows, fixtures.Meows[14])

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckCatFixtures(fixtures)

			cats := &[]*Cat{
				fixtures.Cats[0],
				fixtures.Cats[1],
				fixtures.Cats[2],
				fixtures.Cats[3],
				fixtures.Cats[4],
				fixtures.Cats[5],
				fixtures.Cats[6],
				fixtures.Cats[7],
			}

			err := sqlxx.Preload(ctx, driver, &cats, "Feeder", "Meows")
			is.NoError(err)
			is.Len((*cats), 8)
			is.Equal(fixtures.Cats[0].ID, (*cats)[0].ID)
			is.Equal(fixtures.Cats[1].ID, (*cats)[1].ID)
			is.Equal(fixtures.Cats[2].ID, (*cats)[2].ID)
			is.Equal(fixtures.Cats[3].ID, (*cats)[3].ID)
			is.Equal(fixtures.Cats[4].ID, (*cats)[4].ID)
			is.Equal(fixtures.Cats[5].ID, (*cats)[5].ID)
			is.Equal(fixtures.Cats[6].ID, (*cats)[6].ID)
			is.Equal(fixtures.Cats[7].ID, (*cats)[7].ID)

			is.NotNil((*cats)[0].Feeder)
			is.Equal(fixtures.Humans[0].ID, (*cats)[0].Feeder.ID)
			is.Equal(fixtures.Humans[0].Name, (*cats)[0].Feeder.Name)
			is.NotEmpty((*cats)[0].Meows)
			is.Len((*cats)[0].Meows, 3)
			is.Contains((*cats)[0].Meows, fixtures.Meows[0])
			is.Contains((*cats)[0].Meows, fixtures.Meows[1])
			is.Contains((*cats)[0].Meows, fixtures.Meows[2])

			is.NotNil((*cats)[1].Feeder)
			is.Equal(fixtures.Humans[1].ID, (*cats)[1].Feeder.ID)
			is.Equal(fixtures.Humans[1].Name, (*cats)[1].Feeder.Name)
			is.Empty((*cats)[1].Meows)

			is.Nil((*cats)[2].Feeder)
			is.NotEmpty((*cats)[2].Meows)
			is.Len((*cats)[2].Meows, 1)
			is.Contains((*cats)[2].Meows, fixtures.Meows[3])

			is.NotNil((*cats)[3].Feeder)
			is.Equal(fixtures.Humans[3].ID, (*cats)[3].Feeder.ID)
			is.Equal(fixtures.Humans[3].Name, (*cats)[3].Feeder.Name)
			is.NotEmpty((*cats)[3].Meows)
			is.Len((*cats)[3].Meows, 3)
			is.Contains((*cats)[3].Meows, fixtures.Meows[4])
			is.Contains((*cats)[3].Meows, fixtures.Meows[5])
			is.Contains((*cats)[3].Meows, fixtures.Meows[6])

			is.Nil((*cats)[4].Feeder)
			is.NotEmpty((*cats)[4].Meows)
			is.Len((*cats)[4].Meows, 1)
			is.Contains((*cats)[4].Meows, fixtures.Meows[7])

			is.Nil((*cats)[5].Feeder)
			is.Empty((*cats)[5].Meows)

			is.NotNil((*cats)[6].Feeder)
			is.Equal(fixtures.Humans[6].ID, (*cats)[6].Feeder.ID)
			is.Equal(fixtures.Humans[6].Name, (*cats)[6].Feeder.Name)
			is.NotEmpty((*cats)[6].Meows)
			is.Len((*cats)[6].Meows, 2)
			is.Contains((*cats)[6].Meows, fixtures.Meows[8])
			is.Contains((*cats)[6].Meows, fixtures.Meows[9])

			is.Nil((*cats)[7].Feeder)
			is.NotEmpty((*cats)[7].Meows)
			is.Len((*cats)[7].Meows, 5)
			is.Contains((*cats)[7].Meows, fixtures.Meows[10])
			is.Contains((*cats)[7].Meows, fixtures.Meows[11])
			is.Contains((*cats)[7].Meows, fixtures.Meows[12])
			is.Contains((*cats)[7].Meows, fixtures.Meows[13])
			is.Contains((*cats)[7].Meows, fixtures.Meows[14])

		}
		{

			cats := []Cat{}

			err := sqlxx.Preload(ctx, driver, &cats, "Feeder", "Meows")
			is.NoError(err)
			is.Len(cats, 0)

		}
	})
}
