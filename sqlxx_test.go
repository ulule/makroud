package sqlxx_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

var dbDefaultOptions = map[string]sqlxx.Option{
	"USER":     sqlxx.User("postgres"),
	"PASSWORD": sqlxx.Password(""),
	"HOST":     sqlxx.Host("localhost"),
	"PORT":     sqlxx.Port(5432),
	"NAME":     sqlxx.Database("sqlxx_test"),
}

// ----------------------------------------------------------------------------
// Miscellaneous models
// ----------------------------------------------------------------------------

type Elements struct {
	Air     string `db:"air"`
	Fire    string `sqlxx:"column:fire"`
	Water   string `sqlxx:"-"`
	Earth   string `sqlxx:"column:earth,default"`
	Fifth   string
	enabled bool
}

func (Elements) TableName() string {
	return "rune_elements"
}

// ----------------------------------------------------------------------------
// Object storage application
// ----------------------------------------------------------------------------

type ExoCloudFixtures struct {
	Regions    []*ExoRegion
	Buckets    []*ExoBucket
	Modes      []*ExoChunkMode
	Chunks     []*ExoChunk
	Signatures []*ExoChunkSignature
}

func GenerateExoCloudFixtures(ctx context.Context, driver sqlxx.Driver, is *require.Assertions) ExoCloudFixtures {
	fixtures := ExoCloudFixtures{
		Regions:    []*ExoRegion{},
		Modes:      []*ExoChunkMode{},
		Chunks:     []*ExoChunk{},
		Signatures: []*ExoChunkSignature{},
	}

	region1 := &ExoRegion{
		Name:     "eu-west-1",
		Hostname: "eu-west-1.exocloud.com",
	}
	err := sqlxx.Save(ctx, driver, region1)
	is.NoError(err)
	is.NotEmpty(region1.ID)
	fixtures.Regions = append(fixtures.Regions, region1)

	region2 := &ExoRegion{
		Name:     "eu-west-2",
		Hostname: "eu-west-2.exocloud.com",
	}
	err = sqlxx.Save(ctx, driver, region2)
	is.NoError(err)
	is.NotEmpty(region2.ID)
	fixtures.Regions = append(fixtures.Regions, region2)

	region3 := &ExoRegion{
		Name:     "eu-west-3",
		Hostname: "eu-west-3.exocloud.com",
	}
	err = sqlxx.Save(ctx, driver, region3)
	is.NoError(err)
	is.NotEmpty(region3.ID)
	fixtures.Regions = append(fixtures.Regions, region3)

	bucket1 := &ExoBucket{
		Name:        "com.nemoworld.sandbox.media",
		Description: "Media bucket for sandbox env",
		RegionID:    region1.ID,
	}
	err = sqlxx.Save(ctx, driver, bucket1)
	is.NoError(err)
	is.NotEmpty(bucket1.ID)
	fixtures.Buckets = append(fixtures.Buckets, bucket1)

	bucket2 := &ExoBucket{
		Name:        "com.nemoworld.production.media",
		Description: "Media bucket for production env",
		RegionID:    region1.ID,
	}
	err = sqlxx.Save(ctx, driver, bucket2)
	is.NoError(err)
	is.NotEmpty(bucket2.ID)
	fixtures.Buckets = append(fixtures.Buckets, bucket2)

	bucket3 := &ExoBucket{
		Name:        "com.nemoworld.sandbox.static",
		Description: "Assets for sandbox env",
		RegionID:    region3.ID,
	}
	err = sqlxx.Save(ctx, driver, bucket3)
	is.NoError(err)
	is.NotEmpty(bucket3.ID)
	fixtures.Buckets = append(fixtures.Buckets, bucket3)

	bucket4 := &ExoBucket{
		Name:        "com.nemoworld.production.static",
		Description: "Assets for production env",
		RegionID:    region3.ID,
	}
	err = sqlxx.Save(ctx, driver, bucket4)
	is.NoError(err)
	is.NotEmpty(bucket4.ID)
	fixtures.Buckets = append(fixtures.Buckets, bucket4)

	mode1 := &ExoChunkMode{
		Mode: "rwx",
	}
	err = sqlxx.Save(ctx, driver, mode1)
	is.NoError(err)
	is.NotEmpty(mode1.ID)
	fixtures.Modes = append(fixtures.Modes, mode1)

	mode2 := &ExoChunkMode{
		Mode: "r-x",
	}
	err = sqlxx.Save(ctx, driver, mode2)
	is.NoError(err)
	is.NotEmpty(mode2.ID)
	fixtures.Modes = append(fixtures.Modes, mode2)

	mode3 := &ExoChunkMode{
		Mode: "r-x",
	}
	err = sqlxx.Save(ctx, driver, mode3)
	is.NoError(err)
	is.NotEmpty(mode3.ID)
	fixtures.Modes = append(fixtures.Modes, mode3)

	mode4 := &ExoChunkMode{
		Mode: "rw-",
	}
	err = sqlxx.Save(ctx, driver, mode4)
	is.NoError(err)
	is.NotEmpty(mode4.ID)
	fixtures.Modes = append(fixtures.Modes, mode4)

	mode5 := &ExoChunkMode{
		Mode: "udp-stream",
	}
	err = sqlxx.Save(ctx, driver, mode5)
	is.NoError(err)
	is.NotEmpty(mode5.ID)
	fixtures.Modes = append(fixtures.Modes, mode5)

	mode6 := &ExoChunkMode{
		Mode: "tcp-stream",
	}
	err = sqlxx.Save(ctx, driver, mode6)
	is.NoError(err)
	is.NotEmpty(mode6.ID)
	fixtures.Modes = append(fixtures.Modes, mode6)

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
	fixtures.Chunks = append(fixtures.Chunks, chunk1)

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
	fixtures.Chunks = append(fixtures.Chunks, chunk2)

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
	fixtures.Chunks = append(fixtures.Chunks, chunk3)

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
	fixtures.Chunks = append(fixtures.Chunks, chunk4)

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
	fixtures.Chunks = append(fixtures.Chunks, chunk5)

	chunk6 := &ExoChunk{
		ModeID: mode3.ID,
		Bytes: fmt.Sprint(
			"52a5b98773dc9adc76eb1813c8766a27ac300bee2941c84947592fb75b65de10",
			"48962fd9ed4f3a6bc1e1f869073de7d31748df3ffbbf6bb0526a466db2abc06b",
			"732a3018c21bc1426dde5bd75c57d494fd94c9a8b8de5905673996364e1bc0ac",
		),
	}
	err = sqlxx.Save(ctx, driver, chunk6)
	is.NoError(err)
	is.NotEmpty(chunk6.Hash)
	fixtures.Chunks = append(fixtures.Chunks, chunk6)

	chunk7 := &ExoChunk{
		ModeID: mode4.ID,
		Bytes: fmt.Sprint(
			"3d394c85e961e8ec976162377f46287d47a76968edcc2c1aa08d6667e73cbaf6",
			"678cc918133eb85e78891ae1982bbb9b2306a8aae19f1e8cbbe00909233c5392",
			"183d4adae01c481da776cb27ca5c3a56a2e43ee721d1bcfbd391abcb4ebf10b4",
		),
	}
	err = sqlxx.Save(ctx, driver, chunk7)
	is.NoError(err)
	is.NotEmpty(chunk7.Hash)
	fixtures.Chunks = append(fixtures.Chunks, chunk7)

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
	fixtures.Signatures = append(fixtures.Signatures, signature1)

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
	fixtures.Signatures = append(fixtures.Signatures, signature2)

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
	fixtures.Signatures = append(fixtures.Signatures, signature3)

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
	fixtures.Signatures = append(fixtures.Signatures, signature4)

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
	fixtures.Signatures = append(fixtures.Signatures, signature5)

	return fixtures
}

type ExoRegion struct {
	// Columns
	ID       string `sqlxx:"column:id,pk:ulid"`
	Name     string `sqlxx:"column:name"`
	Hostname string `sqlxx:"column:hostname"`
	// Relationships
	Buckets *[]ExoBucket
}

func (ExoRegion) TableName() string {
	return "exo_region"
}

type ExoBucket struct {
	// Columns
	ID          string `sqlxx:"column:id,pk:ulid"`
	Name        string `sqlxx:"column:name"`
	Description string `sqlxx:"column:description"`
	RegionID    string `sqlxx:"column:region_id,fk:exo_region"`
	// Relationships
	Region ExoRegion
}

func (ExoBucket) TableName() string {
	return "exo_bucket"
}

// type ExoFile struct {
// 	// Columns
// 	ID   string `sqlxx:"column:id,pk:ulid"`
// 	Name string `sqlxx:"column:name"`
// 	Path string `sqlxx:"column:path"`
// 	// Relationships
// 	Chunk []ExoChunk
// }
//
// func (ExoFile) TableName() string {
// 	return "exo_file"
// }
//
// type ExoFileChunk struct {
// 	// Columns
// 	ID      string `sqlxx:"column:id,pk:ulid"`
// 	FileID  string `sqlxx:"column:file_id,fk:exo_file"`
// 	ChunkID string `sqlxx:"column:chunk_id,fk:exo_chunk"`
// }
//
// func (ExoFileChunk) TableName() string {
// 	return "exo_file_chunk"
// }

type ExoChunk struct {
	// Columns
	Hash   string `sqlxx:"column:hash,pk:ulid"`
	Bytes  string `sqlxx:"column:bytes"`
	ModeID string `sqlxx:"column:mode_id,fk:exo_chunk_mode"`
	// Relationships
	Signature *ExoChunkSignature
	Mode      *ExoChunkMode
}

func (ExoChunk) TableName() string {
	return "exo_chunk"
}

type ExoChunkSignature struct {
	ID      string `sqlxx:"column:id,pk:ulid"`
	ChunkID string `sqlxx:"column:chunk_id,fk:exo_chunk"`
	Bytes   string `sqlxx:"column:bytes"`
}

func (ExoChunkSignature) TableName() string {
	return "exo_chunk_signature"
}

type ExoChunkMode struct {
	ID   string `sqlxx:"column:id,pk:ulid"`
	Mode string `sqlxx:"column:mode"`
}

func (ExoChunkMode) TableName() string {
	return "exo_chunk_mode"
}

// ----------------------------------------------------------------------------
// Zootopia
// ----------------------------------------------------------------------------

type ZootopiaFixtures struct {
	Groups   []*Group
	Centers  []*Center
	Owls     []*Owl
	Bags     []*Bag
	Packages []*Package
	Cats     []*Cat
	Humans   []*Human
	Meows    []*Meow
}

func GenerateZootopiaFixtures(ctx context.Context, driver sqlxx.Driver, is *require.Assertions) ZootopiaFixtures {
	fixtures := ZootopiaFixtures{
		Groups:   []*Group{},
		Centers:  []*Center{},
		Owls:     []*Owl{},
		Bags:     []*Bag{},
		Packages: []*Package{},
		Cats:     []*Cat{},
		Humans:   []*Human{},
		Meows:    []*Meow{},
	}

	group1 := &Group{
		Name: "Spring",
	}
	err := sqlxx.Save(ctx, driver, group1)
	is.NoError(err)
	is.NotEmpty(group1.ID)
	fixtures.Groups = append(fixtures.Groups, group1)

	group2 := &Group{
		Name: "Summer",
	}
	err = sqlxx.Save(ctx, driver, group2)
	is.NoError(err)
	is.NotEmpty(group2.ID)
	fixtures.Groups = append(fixtures.Groups, group2)

	group3 := &Group{
		Name: "Winter",
	}
	err = sqlxx.Save(ctx, driver, group3)
	is.NoError(err)
	is.NotEmpty(group3.ID)
	fixtures.Groups = append(fixtures.Groups, group3)

	group4 := &Group{
		Name: "Fall",
	}
	err = sqlxx.Save(ctx, driver, group4)
	is.NoError(err)
	is.NotEmpty(group4.ID)
	fixtures.Groups = append(fixtures.Groups, group4)

	center1 := &Center{
		Name: "Soul",
		Area: "Lancaster",
	}
	err = sqlxx.Save(ctx, driver, center1)
	is.NoError(err)
	is.NotEmpty(center1.ID)
	fixtures.Centers = append(fixtures.Centers, center1)

	center2 := &Center{
		Name: "Cloud",
		Area: "Nancledra",
	}
	err = sqlxx.Save(ctx, driver, center2)
	is.NoError(err)
	is.NotEmpty(center2.ID)
	fixtures.Centers = append(fixtures.Centers, center2)

	center3 := &Center{
		Name: "Gold",
		Area: "Woodhurst",
	}
	err = sqlxx.Save(ctx, driver, center3)
	is.NoError(err)
	is.NotEmpty(center3.ID)
	fixtures.Centers = append(fixtures.Centers, center3)

	center4 := &Center{
		Name: "Moonstone",
		Area: "Armskirk",
	}
	err = sqlxx.Save(ctx, driver, center4)
	is.NoError(err)
	is.NotEmpty(center4.ID)
	fixtures.Centers = append(fixtures.Centers, center4)

	center5 := &Center{
		Name: "Celestial",
		Area: "Bayside",
	}
	err = sqlxx.Save(ctx, driver, center5)
	is.NoError(err)
	is.NotEmpty(center5.ID)
	fixtures.Centers = append(fixtures.Centers, center5)

	center6 := &Center{
		Name: "Solitude",
		Area: "Black Castle",
	}
	err = sqlxx.Save(ctx, driver, center6)
	is.NoError(err)
	is.NotEmpty(center6.ID)
	fixtures.Centers = append(fixtures.Centers, center6)

	owl1 := &Owl{
		Name:         "Pyro",
		FeatherColor: "Timeless Sanguine",
		FavoriteFood: "Ginger Mooncake",
		GroupID: sql.NullInt64{
			Valid: true,
			Int64: group1.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, owl1)
	is.NoError(err)
	is.NotEmpty(owl1.ID)
	fixtures.Owls = append(fixtures.Owls, owl1)

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
	fixtures.Owls = append(fixtures.Owls, owl2)

	owl3 := &Owl{
		Name:         "Wacky",
		FeatherColor: "Harsh Cyan",
		FavoriteFood: "Pecan Trifle",
		GroupID: sql.NullInt64{
			Valid: true,
			Int64: group1.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, owl3)
	is.NoError(err)
	is.NotEmpty(owl3.ID)
	fixtures.Owls = append(fixtures.Owls, owl3)

	owl4 := &Owl{
		Name:         "Puffins",
		FeatherColor: "Botanic Ruby",
		FavoriteFood: "Avocado Salmon",
		GroupID: sql.NullInt64{
			Valid: true,
			Int64: group2.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, owl4)
	is.NoError(err)
	is.NotEmpty(owl4.ID)
	fixtures.Owls = append(fixtures.Owls, owl4)

	owl5 := &Owl{
		Name:         "Pistache",
		FeatherColor: "Distorted Cherry",
		FavoriteFood: "Blueberry Milk",
		GroupID: sql.NullInt64{
			Valid: true,
			Int64: group3.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, owl5)
	is.NoError(err)
	is.NotEmpty(owl5.ID)
	fixtures.Owls = append(fixtures.Owls, owl5)

	owl6 := &Owl{
		Name:         "Baloo",
		FeatherColor: "Supreme Mauve",
		FavoriteFood: "Tomato Turkey",
		GroupID: sql.NullInt64{
			Valid: true,
			Int64: group4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, owl6)
	is.NoError(err)
	is.NotEmpty(owl6.ID)
	fixtures.Owls = append(fixtures.Owls, owl6)

	bag1 := &Bag{
		OwlID: owl1.ID,
		Color: "Frosty Cyan",
	}
	err = sqlxx.Save(ctx, driver, bag1)
	is.NoError(err)
	is.NotEmpty(bag1.ID)
	fixtures.Bags = append(fixtures.Bags, bag1)

	bag2 := &Bag{
		OwlID: owl2.ID,
		Color: "Burned Blue",
	}
	err = sqlxx.Save(ctx, driver, bag2)
	is.NoError(err)
	is.NotEmpty(bag2.ID)
	fixtures.Bags = append(fixtures.Bags, bag2)

	bag3 := &Bag{
		OwlID: owl4.ID,
		Color: "Ordinary Maroon",
	}
	err = sqlxx.Save(ctx, driver, bag3)
	is.NoError(err)
	is.NotEmpty(bag3.ID)
	fixtures.Bags = append(fixtures.Bags, bag3)

	bag4 := &Bag{
		OwlID: owl5.ID,
		Color: "Misty Lemon",
	}
	err = sqlxx.Save(ctx, driver, bag4)
	is.NoError(err)
	is.NotEmpty(bag4.ID)
	fixtures.Bags = append(fixtures.Bags, bag4)

	bag5 := &Bag{
		OwlID: owl6.ID,
		Color: "Lustrous Onyx",
	}
	err = sqlxx.Save(ctx, driver, bag5)
	is.NoError(err)
	is.NotEmpty(bag5.ID)
	fixtures.Bags = append(fixtures.Bags, bag5)

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
	fixtures.Packages = append(fixtures.Packages, pack1)

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
	fixtures.Packages = append(fixtures.Packages, pack2)

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
	fixtures.Packages = append(fixtures.Packages, pack3)

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
	fixtures.Packages = append(fixtures.Packages, pack4)
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
	fixtures.Packages = append(fixtures.Packages, pack5)

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
	fixtures.Packages = append(fixtures.Packages, pack6)

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
	fixtures.Packages = append(fixtures.Packages, pack7)

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
	fixtures.Packages = append(fixtures.Packages, pack8)

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
	fixtures.Packages = append(fixtures.Packages, pack9)

	pack10 := &Package{
		SenderID:   center4.ID,
		ReceiverID: center6.ID,
		Status:     "delivered",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack10)
	is.NoError(err)
	is.NotEmpty(pack10.ID)
	fixtures.Packages = append(fixtures.Packages, pack10)

	pack11 := &Package{
		SenderID:   center4.ID,
		ReceiverID: center6.ID,
		Status:     "delivered",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack11)
	is.NoError(err)
	is.NotEmpty(pack11.ID)
	fixtures.Packages = append(fixtures.Packages, pack11)

	pack12 := &Package{
		SenderID:   center4.ID,
		ReceiverID: center6.ID,
		Status:     "delivered",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack12)
	is.NoError(err)
	is.NotEmpty(pack12.ID)
	fixtures.Packages = append(fixtures.Packages, pack12)

	pack13 := &Package{
		SenderID:   center4.ID,
		ReceiverID: center6.ID,
		Status:     "processing",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack13)
	is.NoError(err)
	is.NotEmpty(pack13.ID)
	fixtures.Packages = append(fixtures.Packages, pack13)

	pack14 := &Package{
		SenderID:   center6.ID,
		ReceiverID: center5.ID,
		Status:     "delivered",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl6.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack14)
	is.NoError(err)
	is.NotEmpty(pack14.ID)
	fixtures.Packages = append(fixtures.Packages, pack14)

	pack15 := &Package{
		SenderID:   center6.ID,
		ReceiverID: center5.ID,
		Status:     "delivered",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl6.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack15)
	is.NoError(err)
	is.NotEmpty(pack15.ID)
	fixtures.Packages = append(fixtures.Packages, pack15)

	pack16 := &Package{
		SenderID:   center6.ID,
		ReceiverID: center5.ID,
		Status:     "processing",
		TransporterID: sql.NullInt64{
			Valid: true,
			Int64: owl6.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, pack16)
	is.NoError(err)
	is.NotEmpty(pack16.ID)
	fixtures.Packages = append(fixtures.Packages, pack16)

	cat1 := &Cat{
		Name: "Eagle",
	}
	err = sqlxx.Save(ctx, driver, cat1)
	is.NoError(err)
	is.NotEmpty(cat1.ID)
	fixtures.Cats = append(fixtures.Cats, cat1)

	cat2 := &Cat{
		Name: "Zigzag",
	}
	err = sqlxx.Save(ctx, driver, cat2)
	is.NoError(err)
	is.NotEmpty(cat2.ID)
	fixtures.Cats = append(fixtures.Cats, cat2)

	cat3 := &Cat{
		Name: "Scully",
	}
	err = sqlxx.Save(ctx, driver, cat3)
	is.NoError(err)
	is.NotEmpty(cat3.ID)
	fixtures.Cats = append(fixtures.Cats, cat3)

	cat4 := &Cat{
		Name: "Hooker",
	}
	err = sqlxx.Save(ctx, driver, cat4)
	is.NoError(err)
	is.NotEmpty(cat4.ID)
	fixtures.Cats = append(fixtures.Cats, cat4)

	cat5 := &Cat{
		Name: "Ditty",
	}
	err = sqlxx.Save(ctx, driver, cat5)
	is.NoError(err)
	is.NotEmpty(cat5.ID)
	fixtures.Cats = append(fixtures.Cats, cat5)

	cat6 := &Cat{
		Name: "Dinky",
	}
	err = sqlxx.Save(ctx, driver, cat6)
	is.NoError(err)
	is.NotEmpty(cat6.ID)
	fixtures.Cats = append(fixtures.Cats, cat6)

	cat7 := &Cat{
		Name: "Flick",
	}
	err = sqlxx.Save(ctx, driver, cat7)
	is.NoError(err)
	is.NotEmpty(cat7.ID)
	fixtures.Cats = append(fixtures.Cats, cat7)

	cat8 := &Cat{
		Name: "Icarus",
	}
	err = sqlxx.Save(ctx, driver, cat8)
	is.NoError(err)
	is.NotEmpty(cat8.ID)
	fixtures.Cats = append(fixtures.Cats, cat8)

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
	fixtures.Humans = append(fixtures.Humans, human1)

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
	fixtures.Humans = append(fixtures.Humans, human2)

	human3 := &Human{
		Name: "Larry Golade",
	}
	err = sqlxx.Save(ctx, driver, human3)
	is.NoError(err)
	is.NotEmpty(human3.ID)
	fixtures.Humans = append(fixtures.Humans, human3)

	human4 := &Human{
		Name: "Roland Culé",
		CatID: sql.NullString{
			Valid:  true,
			String: cat4.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, human4)
	is.NoError(err)
	is.NotEmpty(human4.ID)
	fixtures.Humans = append(fixtures.Humans, human4)

	human5 := &Human{
		Name: "Alex Ception",
	}
	err = sqlxx.Save(ctx, driver, human5)
	is.NoError(err)
	is.NotEmpty(human5.ID)
	fixtures.Humans = append(fixtures.Humans, human5)

	human6 := &Human{
		Name: "Djamal Alatête",
	}
	err = sqlxx.Save(ctx, driver, human6)
	is.NoError(err)
	is.NotEmpty(human6.ID)
	fixtures.Humans = append(fixtures.Humans, human6)

	human7 := &Human{
		Name: " Guy Tar",
		CatID: sql.NullString{
			Valid:  true,
			String: cat7.ID,
		},
	}
	err = sqlxx.Save(ctx, driver, human7)
	is.NoError(err)
	is.NotEmpty(human7.ID)
	fixtures.Humans = append(fixtures.Humans, human7)

	meow1 := &Meow{
		Body:  "meow_00000000001",
		CatID: cat1.ID,
	}
	err = sqlxx.Save(ctx, driver, meow1)
	is.NoError(err)
	is.NotEmpty(meow1.Hash)
	fixtures.Meows = append(fixtures.Meows, meow1)

	meow2 := &Meow{
		Body:  "meow_00000000002",
		CatID: cat1.ID,
	}
	err = sqlxx.Save(ctx, driver, meow2)
	is.NoError(err)
	is.NotEmpty(meow2.Hash)
	fixtures.Meows = append(fixtures.Meows, meow2)

	meow3 := &Meow{
		Body:  "meow_00000000003",
		CatID: cat1.ID,
	}
	err = sqlxx.Save(ctx, driver, meow3)
	is.NoError(err)
	is.NotEmpty(meow3.Hash)
	fixtures.Meows = append(fixtures.Meows, meow3)

	meow4 := &Meow{
		Body:  "meow_00000000004",
		CatID: cat3.ID,
	}
	err = sqlxx.Save(ctx, driver, meow4)
	is.NoError(err)
	is.NotEmpty(meow4.Hash)
	fixtures.Meows = append(fixtures.Meows, meow4)

	meow5 := &Meow{
		Body:  "meow_00000000005",
		CatID: cat4.ID,
	}
	err = sqlxx.Save(ctx, driver, meow5)
	is.NoError(err)
	is.NotEmpty(meow5.Hash)
	fixtures.Meows = append(fixtures.Meows, meow5)

	meow6 := &Meow{
		Body:  "meow_00000000006",
		CatID: cat4.ID,
	}
	err = sqlxx.Save(ctx, driver, meow6)
	is.NoError(err)
	is.NotEmpty(meow6.Hash)
	fixtures.Meows = append(fixtures.Meows, meow6)

	meow7 := &Meow{
		Body:  "meow_00000000007",
		CatID: cat4.ID,
	}
	err = sqlxx.Save(ctx, driver, meow7)
	is.NoError(err)
	is.NotEmpty(meow7.Hash)
	fixtures.Meows = append(fixtures.Meows, meow7)

	meow8 := &Meow{
		Body:  "meow_00000000008",
		CatID: cat5.ID,
	}
	err = sqlxx.Save(ctx, driver, meow8)
	is.NoError(err)
	is.NotEmpty(meow8.Hash)
	fixtures.Meows = append(fixtures.Meows, meow8)

	meow9 := &Meow{
		Body:  "meow_00000000009",
		CatID: cat7.ID,
	}
	err = sqlxx.Save(ctx, driver, meow9)
	is.NoError(err)
	is.NotEmpty(meow9.Hash)
	fixtures.Meows = append(fixtures.Meows, meow9)

	meow10 := &Meow{
		Body:  "meow_00000000010",
		CatID: cat7.ID,
	}
	err = sqlxx.Save(ctx, driver, meow10)
	is.NoError(err)
	is.NotEmpty(meow10.Hash)
	fixtures.Meows = append(fixtures.Meows, meow10)

	meow11 := &Meow{
		Body:  "meow_00000000011",
		CatID: cat8.ID,
	}
	err = sqlxx.Save(ctx, driver, meow11)
	is.NoError(err)
	is.NotEmpty(meow11.Hash)
	fixtures.Meows = append(fixtures.Meows, meow11)

	meow12 := &Meow{
		Body:  "meow_00000000012",
		CatID: cat8.ID,
	}
	err = sqlxx.Save(ctx, driver, meow12)
	is.NoError(err)
	is.NotEmpty(meow12.Hash)
	fixtures.Meows = append(fixtures.Meows, meow12)

	meow13 := &Meow{
		Body:  "meow_00000000013",
		CatID: cat8.ID,
	}
	err = sqlxx.Save(ctx, driver, meow13)
	is.NoError(err)
	is.NotEmpty(meow13.Hash)
	fixtures.Meows = append(fixtures.Meows, meow13)

	meow14 := &Meow{
		Body:  "meow_00000000014",
		CatID: cat8.ID,
	}
	err = sqlxx.Save(ctx, driver, meow14)
	is.NoError(err)
	is.NotEmpty(meow14.Hash)
	fixtures.Meows = append(fixtures.Meows, meow14)

	meow15 := &Meow{
		Body:  "meow_00000000015",
		CatID: cat8.ID,
	}
	err = sqlxx.Save(ctx, driver, meow15)
	is.NoError(err)
	is.NotEmpty(meow15.Hash)
	fixtures.Meows = append(fixtures.Meows, meow15)

	return fixtures
}

type Group struct {
	// Columns
	ID   int64  `sqlxx:"column:id,pk"`
	Name string `sqlxx:"column:name"`
}

func (Group) TableName() string {
	return "ztp_group"
}

type Center struct {
	// Columns
	ID   string `sqlxx:"column:id"`
	Name string `sqlxx:"column:name"`
	Area string `sqlxx:"column:area"`
}

func (Center) TableName() string {
	return "ztp_center"
}

type Owl struct {
	// Columns
	ID           int64         `sqlxx:"column:id,pk"`
	Name         string        `sqlxx:"column:name"`
	FeatherColor string        `sqlxx:"column:feather_color"`
	FavoriteFood string        `sqlxx:"column:favorite_food"`
	GroupID      sql.NullInt64 `sqlxx:"column:group_id,fk:ztp_group"`
	// Relationships
	Group    *Group
	Packages []Package
	Bag      *Bag
}

func (Owl) TableName() string {
	return "ztp_owl"
}

type Bag struct {
	// Columns
	ID    int64  `sqlxx:"column:id,pk"`
	Color string `sqlxx:"column:color"`
	OwlID int64  `sqlxx:"column:owl_id,fk:ztp_owl"`
	// Relationships
	Owl Owl
}

func (Bag) TableName() string {
	return "ztp_bag"
}

type Package struct {
	// Columns
	ID            string        `sqlxx:"column:id"`
	Status        string        `sqlxx:"column:status"`
	SenderID      string        `sqlxx:"column:sender_id,fk:ztp_center"`
	ReceiverID    string        `sqlxx:"column:receiver_id,fk:ztp_center"`
	TransporterID sql.NullInt64 `sqlxx:"column:transporter_id,fk:ztp_owl"`
	// Relationships
	Sender   *Center
	Receiver *Center
}

func (Package) TableName() string {
	return "ztp_package"
}

type Cat struct {
	// Columns
	ID        string      `sqlxx:"column:id,pk:ulid"`
	Name      string      `sqlxx:"column:name"`
	CreatedAt time.Time   `sqlxx:"column:created_at,default"`
	UpdatedAt time.Time   `sqlxx:"column:updated_at,default"`
	DeletedAt pq.NullTime `sqlxx:"column:deleted_at"`
	// Relationships
	Feeder *Human
	Meows  []*Meow
}

func (Cat) TableName() string {
	return "ztp_cat"
}

type Meow struct {
	// Columns
	Hash      string      `sqlxx:"column:hash,pk:ulid"`
	Body      string      `sqlxx:"column:body"`
	CatID     string      `sqlxx:"column:cat_id,fk:ztp_cat"`
	CreatedAt time.Time   `sqlxx:"column:created"`
	UpdatedAt time.Time   `sqlxx:"column:updated"`
	DeletedAt pq.NullTime `sqlxx:"column:deleted"`
}

func (Meow) TableName() string {
	return "ztp_meow"
}

func (Meow) CreatedKey() string {
	return "created"
}

func (Meow) UpdatedKey() string {
	return "updated"
}

func (Meow) DeletedKey() string {
	return "deleted"
}

type Human struct {
	// Columns
	ID        string         `sqlxx:"column:id,pk:ulid"`
	Name      string         `sqlxx:"column:name"`
	CreatedAt time.Time      `sqlxx:"column:created_at,default"`
	UpdatedAt time.Time      `sqlxx:"column:updated_at,default"`
	DeletedAt pq.NullTime    `sqlxx:"column:deleted_at"`
	CatID     sql.NullString `sqlxx:"column:cat_id,fk:ztp_cat"`
	// Relationships
	Cat *Cat
}

func (Human) TableName() string {
	return "ztp_human"
}

// ----------------------------------------------------------------------------
// Loader
// ----------------------------------------------------------------------------

type environment struct {
	driver *sqlxx.Client
	is     *require.Assertions
}

func (e *environment) startup(ctx context.Context) {
	DropTables(ctx, e.driver)
	CreateTables(ctx, e.driver)
}

func (e *environment) shutdown(ctx context.Context) {
	value := os.Getenv("DB_KEEP")
	if len(value) == 0 {
		DropTables(ctx, e.driver)
	}
	e.is.NoError(e.driver.Close())
}

func dbParamString(option func(string) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	if len(v) != 0 {
		return option(v)
	}
	for i := range env {
		v = os.Getenv(env[i])
		if len(v) != 0 {
			return option(v)
		}
	}
	return dbDefaultOptions[param]
}

func dbParamInt(option func(int) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	n, err := strconv.Atoi(v)
	if err == nil {
		return option(n)
	}
	for i := range env {
		v = os.Getenv(env[i])
		n, err = strconv.Atoi(v)
		if err == nil {
			return option(n)
		}
	}
	return dbDefaultOptions[param]
}

type SetupCallback func(handler SetupHandler)

type SetupHandler func(driver sqlxx.Driver)

func Setup(t require.TestingT, options ...sqlxx.Option) SetupCallback {
	is := require.New(t)
	ctx := context.Background()
	opts := []sqlxx.Option{
		dbParamString(sqlxx.Host, "host", "PGHOST"),
		dbParamInt(sqlxx.Port, "port", "PGPORT"),
		dbParamString(sqlxx.User, "user", "PGUSER"),
		dbParamString(sqlxx.Password, "password", "PGPASSWORD"),
		dbParamString(sqlxx.Database, "name", "PGDATABASE"),
		sqlxx.Cache(true),
	}
	opts = append(opts, options...)

	db, err := sqlxx.New(opts...)
	is.NoError(err)
	is.NotNil(db)

	env := &environment{
		is:     is,
		driver: db,
	}

	return func(handler SetupHandler) {
		env.startup(ctx)
		handler(db)
		env.shutdown(ctx)
	}
}

func DropTables(ctx context.Context, db *sqlxx.Client) {
	db.MustExec(ctx, `
		-- Simple schema
		DROP TABLE IF EXISTS ztp_human CASCADE;
		DROP TABLE IF EXISTS ztp_package CASCADE;
		DROP TABLE IF EXISTS ztp_bag CASCADE;
		DROP TABLE IF EXISTS ztp_owl CASCADE;
		DROP TABLE IF EXISTS ztp_cat CASCADE;
		DROP TABLE IF EXISTS ztp_meow CASCADE;
		DROP TABLE IF EXISTS ztp_group CASCADE;
		DROP TABLE IF EXISTS ztp_center CASCADE;

		-- Object storage application
		DROP TABLE IF EXISTS exo_chunk_signature CASCADE;
		DROP TABLE IF EXISTS exo_chunk CASCADE;
		DROP TABLE IF EXISTS exo_chunk_mode CASCADE;
		DROP TABLE IF EXISTS exo_bucket CASCADE;
		DROP TABLE IF EXISTS exo_region CASCADE;

	`)
}

func CreateTables(ctx context.Context, db *sqlxx.Client) {
	db.MustExec(ctx, `

		--
		-- Zootopia schema
		--

		CREATE TABLE ztp_group (
			id                SERIAL PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL
		);
		CREATE TABLE ztp_center (
			id                VARCHAR(32) PRIMARY KEY NOT NULL DEFAULT md5(random()::text),
			name              VARCHAR(255) NOT NULL,
			area              VARCHAR(255) NOT NULL
		);
		CREATE TABLE ztp_owl (
			id                SERIAL PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			feather_color     VARCHAR(255) NOT NULL,
			favorite_food     VARCHAR(255) NOT NULL,
			group_id          INTEGER REFERENCES ztp_group(id)
		);
		CREATE TABLE ztp_bag (
			id                SERIAL PRIMARY KEY NOT NULL,
			color             VARCHAR(255) NOT NULL,
			owl_id            INTEGER NOT NULL REFERENCES ztp_owl(id)
		);
		CREATE TABLE ztp_cat (
			id                VARCHAR(26) PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at        TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_meow (
			hash              VARCHAR(26) PRIMARY KEY NOT NULL,
			body              VARCHAR(2048) NOT NULL,
			cat_id            VARCHAR(26) NOT NULL REFERENCES ztp_cat(id),
			created           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted           TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_human (
			id                VARCHAR(26) PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			cat_id            VARCHAR(26) REFERENCES ztp_cat(id),
			created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at        TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_package (
			id                VARCHAR(32) PRIMARY KEY NOT NULL DEFAULT md5(random()::text),
			status            VARCHAR(255) NOT NULL,
			sender_id         VARCHAR(32) NOT NULL REFERENCES ztp_center(id),
			receiver_id       VARCHAR(32) NOT NULL REFERENCES ztp_center(id),
			transporter_id    INTEGER REFERENCES ztp_owl(id)
		);

		--
		-- Object storage application
		--

		CREATE TABLE exo_region (
			id              VARCHAR(26) PRIMARY KEY NOT NULL,
			name            VARCHAR(255) NOT NULL,
			hostname        VARCHAR(2048) NOT NULL
		);
		CREATE TABLE exo_bucket (
			id              VARCHAR(26) PRIMARY KEY NOT NULL,
			name            VARCHAR(512) NOT NULL,
			description     VARCHAR(2048) NOT NULL,
			region_id       VARCHAR(26) NOT NULL REFERENCES exo_region(id)
		);
		CREATE TABLE exo_chunk_mode (
			id              VARCHAR(26) PRIMARY KEY NOT NULL,
			mode            VARCHAR(255) NOT NULL
		);
		CREATE TABLE exo_chunk (
			hash            VARCHAR(26) PRIMARY KEY NOT NULL,
			bytes           VARCHAR(2048) NOT NULL,
			mode_id         VARCHAR(26) NOT NULL REFERENCES exo_chunk_mode(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_chunk_signature (
			id              VARCHAR(26) PRIMARY KEY NOT NULL,
			chunk_id        VARCHAR(26) NOT NULL REFERENCES exo_chunk(hash),
			bytes           VARCHAR(2048) NOT NULL
		);

		--
		-- Application schema
		--

		-- TODO

	`)
}
