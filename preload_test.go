package sqlxx_test

import (
	"fmt"
	"testing"

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
		is.NotEmpty("mode:", mode.ID)
		fmt.Println("mode:", mode.ID)

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
		fmt.Println("chunck:", chunk.Hash)

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
		fmt.Println("signature:", signature.ID)

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

	})
}
