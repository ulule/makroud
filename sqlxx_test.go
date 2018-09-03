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
	Regions       []*ExoRegion
	Organizations []*ExoOrganization
	Avatars       []*ExoAvatar
	Profiles      []*ExoProfile
	Users         []*ExoUser
	Groups        []*ExoGroup
	Buckets       []*ExoBucket
	Directories   []*ExoDirectory
	Files         []*ExoFile
	Modes         []*ExoChunkMode
	Chunks        []*ExoChunk
	Signatures    []*ExoChunkSignature
}

func GenerateExoCloudFixtures(ctx context.Context, driver sqlxx.Driver, is *require.Assertions) *ExoCloudFixtures {
	fixtures := &ExoCloudFixtures{
		Regions:       []*ExoRegion{},
		Organizations: []*ExoOrganization{},
		Avatars:       []*ExoAvatar{},
		Profiles:      []*ExoProfile{},
		Users:         []*ExoUser{},
		Groups:        []*ExoGroup{},
		Buckets:       []*ExoBucket{},
		Directories:   []*ExoDirectory{},
		Files:         []*ExoFile{},
		Modes:         []*ExoChunkMode{},
		Chunks:        []*ExoChunk{},
		Signatures:    []*ExoChunkSignature{},
	}

	region1 := &ExoRegion{
		Name:     "eu-west-1",
		Hostname: "eu-west-1.exocloud.com",
	}
	fixtures.Regions = append(fixtures.Regions, region1)

	region2 := &ExoRegion{
		Name:     "eu-west-2",
		Hostname: "eu-west-2.exocloud.com",
	}
	fixtures.Regions = append(fixtures.Regions, region2)

	region3 := &ExoRegion{
		Name:     "eu-west-3",
		Hostname: "eu-west-3.exocloud.com",
	}
	fixtures.Regions = append(fixtures.Regions, region3)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Regions {
			err := sqlxx.Save(ctx, dbx, fixtures.Regions[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Regions[i].ID)
		}
		return nil
	})

	organization1 := &ExoOrganization{
		Name:    "Nemo World",
		Website: "https://www.nemoworld.com",
	}
	fixtures.Organizations = append(fixtures.Organizations, organization1)

	organization2 := &ExoOrganization{
		Name:    "Shadow Navigations",
		Website: "https://www.shadow-navigations.co.uk",
	}
	fixtures.Organizations = append(fixtures.Organizations, organization2)

	organization3 := &ExoOrganization{
		Name:    "Ansoft",
		Website: "https://ansoft.io",
	}
	fixtures.Organizations = append(fixtures.Organizations, organization3)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Organizations {
			err := sqlxx.Save(ctx, dbx, fixtures.Organizations[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Organizations[i].ID)
		}
		return nil
	})

	avatar1 := &ExoAvatar{
		URL:      "https://img.exocloud.com/5/b/1c1ca2b0b9581a5ff913ad17ab5884d2.png",
		Path:     "avatar://5/b/1c1ca2b0b9581a5ff913ad17ab5884d2.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar1)

	avatar2 := &ExoAvatar{
		URL:      "https://img.exocloud.com/b/8/acec4a0f22f5e2e635b801dee5d50057.png",
		Path:     "avatar://b/8/acec4a0f22f5e2e635b801dee5d50057.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar2)

	avatar3 := &ExoAvatar{
		URL:      "https://img.exocloud.com/9/2/c5fbb7bfbb5a7ba77064d8f954ec722f.png",
		Path:     "avatar://9/2/c5fbb7bfbb5a7ba77064d8f954ec722f.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar3)

	avatar4 := &ExoAvatar{
		URL:      "https://img.exocloud.com/6/f/b3e4292ae4c02c1e4948a43d4d828df9.png",
		Path:     "avatar://6/f/b3e4292ae4c02c1e4948a43d4d828df9.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar4)

	avatar5 := &ExoAvatar{
		URL:      "https://img.exocloud.com/f/e/93dece774de7f5e69b0b83accd2563c8.png",
		Path:     "avatar://f/e/93dece774de7f5e69b0b83accd2563c8.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar5)

	avatar6 := &ExoAvatar{
		URL:      "https://img.exocloud.com/f/e/26acaf9b1f30d7706c3ac2083069570d.png",
		Path:     "avatar://f/e/26acaf9b1f30d7706c3ac2083069570d.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar6)

	avatar7 := &ExoAvatar{
		URL:      "https://img.exocloud.com/d/3/565b8a99829b8cccd6d48a98920949f5.png",
		Path:     "avatar://d/3/565b8a99829b8cccd6d48a98920949f5.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar7)

	avatar8 := &ExoAvatar{
		URL:      "https://img.exocloud.com/3/2/e0ed03feba7d0a0affb8094d59c69ab6.png",
		Path:     "avatar://3/2/e0ed03feba7d0a0affb8094d59c69ab6.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar8)

	avatar9 := &ExoAvatar{
		URL:      "https://img.exocloud.com/0/0/a4fe73ceaf03c8235f84a2e0df974076.png",
		Path:     "avatar://0/0/a4fe73ceaf03c8235f84a2e0df974076.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar9)

	avatar10 := &ExoAvatar{
		URL:      "https://img.exocloud.com/7/3/a33e81dd43475fc30b0d45e50dc4c112.png",
		Path:     "avatar://7/3/a33e81dd43475fc30b0d45e50dc4c112.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar10)

	avatar11 := &ExoAvatar{
		URL:      "https://img.exocloud.com/7/a/d4ac5c5cf4cd5d93dd70513771eaae48.png",
		Path:     "avatar://7/a/d4ac5c5cf4cd5d93dd70513771eaae48.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar11)

	avatar12 := &ExoAvatar{
		URL:      "https://img.exocloud.com/c/e/ac5fcf3560e3fd85c4b624c51da75725.png",
		Path:     "avatar://c/e/ac5fcf3560e3fd85c4b624c51da75725.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar12)

	avatar13 := &ExoAvatar{
		URL:      "https://img.exocloud.com/1/3/e1e5bad178d670e83836bbf2408608d8.png",
		Path:     "avatar://1/3/e1e5bad178d670e83836bbf2408608d8.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar13)

	avatar14 := &ExoAvatar{
		URL:      "https://img.exocloud.com/a/5/eb66a46ad1b1a6bed540924663c0fd15.png",
		Path:     "avatar://a/5/eb66a46ad1b1a6bed540924663c0fd15.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar14)

	avatar15 := &ExoAvatar{
		URL:      "https://img.exocloud.com/d/8/74fa532311263faf7ad7d2ffe4cf9d50.png",
		Path:     "avatar://d/8/74fa532311263faf7ad7d2ffe4cf9d50.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar15)

	avatar16 := &ExoAvatar{
		URL:      "https://img.exocloud.com/a/6/f71c369e1bba835a7e59b12b30fc3da5.png",
		Path:     "avatar://a/6/f71c369e1bba835a7e59b12b30fc3da5.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar16)

	avatar17 := &ExoAvatar{
		URL:      "https://img.exocloud.com/6/b/9d5b29a1fb9faa66333b3ff79dc667ae.png",
		Path:     "avatar://6/b/9d5b29a1fb9faa66333b3ff79dc667ae.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar17)

	avatar18 := &ExoAvatar{
		URL:      "https://img.exocloud.com/4/9/7404100d189065881516fe2409076448.png",
		Path:     "avatar://4/9/7404100d189065881516fe2409076448.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar18)

	avatar19 := &ExoAvatar{
		URL:      "https://img.exocloud.com/5/8/3fabc1bdcba2cc7975b810f21b9f9a18.png",
		Path:     "avatar://5/8/3fabc1bdcba2cc7975b810f21b9f9a18.png",
		MimeType: "image/png",
	}
	fixtures.Avatars = append(fixtures.Avatars, avatar19)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Avatars {
			err := sqlxx.Save(ctx, dbx, fixtures.Avatars[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Avatars[i].ID)
		}
		return nil
	})

	profile1 := &ExoProfile{
		FirstName: "Natalie",
		LastName:  "Davis",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar1.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile1)

	profile2 := &ExoProfile{
		FirstName: "Michael",
		LastName:  "Thomas",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar2.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile2)

	profile3 := &ExoProfile{
		FirstName: "Elizabeth",
		LastName:  "Davis",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar3.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile3)

	profile4 := &ExoProfile{
		FirstName: "Noah",
		LastName:  "Johnson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar4.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile4)

	profile5 := &ExoProfile{
		FirstName: "Alexander",
		LastName:  "Jackson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar5.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile5)

	profile6 := &ExoProfile{
		FirstName: "Noah",
		LastName:  "Thomas",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar6.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile6)

	profile7 := &ExoProfile{
		FirstName: "Ava",
		LastName:  "Anderson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar7.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile7)

	profile8 := &ExoProfile{
		FirstName: "Elizabeth",
		LastName:  "Harris",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar8.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile8)

	profile9 := &ExoProfile{
		FirstName: "Sophia",
		LastName:  "Martin",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar9.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile9)

	profile10 := &ExoProfile{
		FirstName: "Lily",
		LastName:  "Robinson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar10.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile10)

	profile11 := &ExoProfile{
		FirstName: "Matthew",
		LastName:  "Garcia",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar11.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile11)

	profile12 := &ExoProfile{
		FirstName: "Zoey",
		LastName:  "Taylor",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar12.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile12)

	profile13 := &ExoProfile{
		FirstName: "Olivia",
		LastName:  "Miller",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar13.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile13)

	profile14 := &ExoProfile{
		FirstName: "Jayden",
		LastName:  "Garcia",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar14.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile14)

	profile15 := &ExoProfile{
		FirstName: "Joseph",
		LastName:  "Johnson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar15.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile15)

	profile16 := &ExoProfile{
		FirstName: "Michael",
		LastName:  "Thompson",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar16.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile16)

	profile17 := &ExoProfile{
		FirstName: "Avery",
		LastName:  "Taylor",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar17.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile17)

	profile18 := &ExoProfile{
		FirstName: "Mason",
		LastName:  "Harris",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar18.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile18)

	profile19 := &ExoProfile{
		FirstName: "Jayden",
		LastName:  "Jones",
		AvatarID: sql.NullString{
			Valid:  true,
			String: avatar19.ID,
		},
	}
	fixtures.Profiles = append(fixtures.Profiles, profile19)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Profiles {
			err := sqlxx.Save(ctx, dbx, fixtures.Profiles[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Profiles[i].ID)
		}
		return nil
	})

	user1 := &ExoUser{
		Email:     "liamwilliams013@example.com",
		Password:  "$2a$10$UrdRyep23chaGNdO4eucpeRxtyLanosfDHyUl1sHeClnP.oq55o9K",
		Country:   "HR",
		Locale:    "en",
		ProfileID: profile1.ID,
	}
	fixtures.Users = append(fixtures.Users, user1)

	user2 := &ExoUser{
		Email:     "aidenjohnson827@test.com",
		Password:  "$2a$10$AVZ54nYb9dtzYLleSlsR8OGVX.ZYlARaknoV4wv9209FtsFJ5X08.",
		Country:   "SN",
		Locale:    "en",
		ProfileID: profile2.ID,
	}
	fixtures.Users = append(fixtures.Users, user2)

	user3 := &ExoUser{
		Email:     "abigailanderson703@example.net",
		Password:  "$2a$10$DnO5zvN3SsrLwD7LK3UCLuWTiTO.Ji21ef6ykqHL67MGftrsZtqmm",
		Country:   "MX",
		Locale:    "en",
		ProfileID: profile3.ID,
	}
	fixtures.Users = append(fixtures.Users, user3)

	user4 := &ExoUser{
		Email:     "masonwhite402@example.com",
		Password:  "$2a$10$4qJsXhFjPG5Et4EG9fvJ3uenHWvrocqn4/mfG2d6.XGGQFVBA3de.",
		Country:   "MN",
		Locale:    "en",
		ProfileID: profile4.ID,
	}
	fixtures.Users = append(fixtures.Users, user4)

	user5 := &ExoUser{
		Email:     "ethangarcia238@example.net",
		Password:  "$2a$10$zN.AGkO/rz3pNTCr8Omnbeuxcj0PKIxFBOQC5hKpQdldxfSWmnEZG",
		Country:   "GW",
		Locale:    "en",
		ProfileID: profile5.ID,
	}
	fixtures.Users = append(fixtures.Users, user5)

	user6 := &ExoUser{
		Email:     "madisonthompson331@example.com",
		Password:  "$2a$10$B.64vM2ooNm9HVXmobL7suxZ2GbA2xuNuEreogEK8SEupqHJPS4LW",
		Country:   "KE",
		Locale:    "en",
		ProfileID: profile6.ID,
	}
	fixtures.Users = append(fixtures.Users, user6)

	user7 := &ExoUser{
		Email:     "madisonthomas556@test.net",
		Password:  "$2a$10$logRwyAUMknL5b9Pg12UvOrRsWdQJYLSJiixTBHTB5roA5ZfS16xG",
		Country:   "EE",
		Locale:    "en",
		ProfileID: profile7.ID,
	}
	fixtures.Users = append(fixtures.Users, user7)

	user8 := &ExoUser{
		Email:     "ethanwhite052@example.net",
		Password:  "$2a$10$JGLbKyfSTa2gfzb3y29Wju9kqDusmEcN6JLnD6Bj71/wfDzdGqUYu",
		Country:   "VC",
		Locale:    "en",
		ProfileID: profile8.ID,
	}
	fixtures.Users = append(fixtures.Users, user8)

	user9 := &ExoUser{
		Email:     "oliviajones344@test.org",
		Password:  "$2a$10$HZMOgZ5nrfPoxITMB4eTEOl//9lEsFDI8gemw2uo9noYQeGC7HuZe",
		Country:   "LU",
		Locale:    "en",
		ProfileID: profile9.ID,
	}
	fixtures.Users = append(fixtures.Users, user9)

	user10 := &ExoUser{
		Email:     "michaelgarcia365@test.net",
		Password:  "$2a$10$pQPj5iE4A9ALx8uhBKJAFenEp.TZeMpZKAPEzTQP25nzNmDRbHtJ6",
		Country:   "BW",
		Locale:    "en",
		ProfileID: profile10.ID,
	}
	fixtures.Users = append(fixtures.Users, user10)

	user11 := &ExoUser{
		Email:     "matthewmoore855@test.com",
		Password:  "$2a$10$7voBreEKUoEmfUoamh3E.eLtoRex2v/lByI7.Sl/pZyPohcuD4fn.",
		Country:   "PT",
		Locale:    "en",
		ProfileID: profile11.ID,
	}
	fixtures.Users = append(fixtures.Users, user11)

	user12 := &ExoUser{
		Email:     "joshuawilliams636@test.org",
		Password:  "$2a$10$StC8vDdXpnnnstcUhy54Gu8DerV3y/AGc1X3WQn1Fn6p65/4Ted.K",
		Country:   "SL",
		Locale:    "en",
		ProfileID: profile12.ID,
	}
	fixtures.Users = append(fixtures.Users, user12)

	user13 := &ExoUser{
		Email:     "jacobwilliams832@example.org",
		Password:  "$2a$10$M2lw7qvMe2VmSYdKqsVDD.NwGw.7ghFMiyOu/a/OKPOziA.kBcKY6",
		Country:   "MG",
		Locale:    "en",
		ProfileID: profile13.ID,
	}
	fixtures.Users = append(fixtures.Users, user13)

	user14 := &ExoUser{
		Email:     "abigailjohnson307@example.net",
		Password:  "$2a$10$B8nxAx5JMcMiHemS/N8tjumQNyroQStItJCwh9KXK.SlCJQTPW1uu",
		Country:   "TM",
		Locale:    "en",
		ProfileID: profile14.ID,
	}
	fixtures.Users = append(fixtures.Users, user14)

	user15 := &ExoUser{
		Email:     "alexanderwhite167@test.net",
		Password:  "$2a$10$YvItrKNIzuL4C7s6KrOe0OCzFIvNr80AWj.rr7je1k05XQi97Vfnm",
		Country:   "IQ",
		Locale:    "en",
		ProfileID: profile15.ID,
	}
	fixtures.Users = append(fixtures.Users, user15)

	user16 := &ExoUser{
		Email:     "addisontaylor728@test.net",
		Password:  "$2a$10$5l.2zlX8hwgH.RpyiaNynuSsUM.ocHcgIZ4nVVzr8fSPk/A0l1sxS",
		Country:   "BH",
		Locale:    "en",
		ProfileID: profile16.ID,
	}
	fixtures.Users = append(fixtures.Users, user16)

	user17 := &ExoUser{
		Email:     "joshuawilson362@test.com",
		Password:  "$2a$10$f/Kx.yP5cQVgzA4bXgpmLudCWXkh0kV2UIqS/7063aRKFELb9aW2C",
		Country:   "PS",
		Locale:    "en",
		ProfileID: profile17.ID,
	}
	fixtures.Users = append(fixtures.Users, user17)

	user18 := &ExoUser{
		Email:     "chloetaylor700@test.com",
		Password:  "$2a$10$FIGsTTJ2YVdrdStzziLjZ.WJ/jQAI8/KLPA88.cKg4CLs4lqhsXRS",
		Country:   "TV",
		Locale:    "en",
		ProfileID: profile18.ID,
	}
	fixtures.Users = append(fixtures.Users, user18)

	user19 := &ExoUser{
		Email:     "avajackson701@example.com",
		Password:  "$2a$10$qNYKEoFH.NDSBPOon.uvqOnmWK4N2x7W.QFp/i2uJXmiNQFqIM/om",
		Country:   "PN",
		Locale:    "en",
		ProfileID: profile19.ID,
	}
	fixtures.Users = append(fixtures.Users, user19)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Users {
			err := sqlxx.Save(ctx, dbx, fixtures.Users[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Users[i].ID)
		}
		return nil
	})

	group1 := &ExoGroup{
		Role:           "user",
		UserID:         user16.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group1)

	group2 := &ExoGroup{
		Role:           "user",
		UserID:         user14.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group2)

	group3 := &ExoGroup{
		Role:           "user",
		UserID:         user2.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group3)

	group4 := &ExoGroup{
		Role:           "user",
		UserID:         user8.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group4)

	group5 := &ExoGroup{
		Role:           "user",
		UserID:         user5.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group5)

	group6 := &ExoGroup{
		Role:           "user",
		UserID:         user10.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group6)

	group7 := &ExoGroup{
		Role:           "user",
		UserID:         user3.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group7)

	group8 := &ExoGroup{
		Role:           "user",
		UserID:         user17.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group8)

	group9 := &ExoGroup{
		Role:           "user",
		UserID:         user7.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group9)

	group10 := &ExoGroup{
		Role:           "user",
		UserID:         user16.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group10)

	group11 := &ExoGroup{
		Role:           "user",
		UserID:         user5.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group11)

	group12 := &ExoGroup{
		Role:           "user",
		UserID:         user9.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group12)

	group13 := &ExoGroup{
		Role:           "user",
		UserID:         user4.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group13)

	group14 := &ExoGroup{
		Role:           "user",
		UserID:         user6.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group14)

	group15 := &ExoGroup{
		Role:           "user",
		UserID:         user11.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group15)

	group16 := &ExoGroup{
		Role:           "user",
		UserID:         user9.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group16)

	group17 := &ExoGroup{
		Role:           "user",
		UserID:         user6.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group17)

	group18 := &ExoGroup{
		Role:           "user",
		UserID:         user4.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group18)

	group19 := &ExoGroup{
		Role:           "user",
		UserID:         user19.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group19)

	group20 := &ExoGroup{
		Role:           "user",
		UserID:         user13.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group20)

	group21 := &ExoGroup{
		Role:           "user",
		UserID:         user11.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group21)

	group22 := &ExoGroup{
		Role:           "user",
		UserID:         user8.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group22)

	group23 := &ExoGroup{
		Role:           "user",
		UserID:         user7.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group23)

	group24 := &ExoGroup{
		Role:           "user",
		UserID:         user9.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group24)

	group25 := &ExoGroup{
		Role:           "user",
		UserID:         user16.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group25)

	group26 := &ExoGroup{
		Role:           "user",
		UserID:         user5.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group26)

	group27 := &ExoGroup{
		Role:           "user",
		UserID:         user14.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group27)

	group28 := &ExoGroup{
		Role:           "user",
		UserID:         user6.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group28)

	group29 := &ExoGroup{
		Role:           "user",
		UserID:         user13.ID,
		OrganizationID: organization3.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group29)

	group30 := &ExoGroup{
		Role:           "user",
		UserID:         user18.ID,
		OrganizationID: organization2.ID,
	}
	fixtures.Groups = append(fixtures.Groups, group30)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Groups {
			err := sqlxx.Save(ctx, dbx, fixtures.Groups[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Groups[i].ID)
		}
		return nil
	})

	bucket1 := &ExoBucket{
		Name:           "com.nemoworld.sandbox.media",
		Description:    "Media bucket for sandbox env",
		RegionID:       region1.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Buckets = append(fixtures.Buckets, bucket1)

	bucket2 := &ExoBucket{
		Name:           "com.nemoworld.production.media",
		Description:    "Media bucket for production env",
		RegionID:       region1.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Buckets = append(fixtures.Buckets, bucket2)

	bucket3 := &ExoBucket{
		Name:           "com.nemoworld.sandbox.static",
		Description:    "Assets for sandbox env",
		RegionID:       region3.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Buckets = append(fixtures.Buckets, bucket3)

	bucket4 := &ExoBucket{
		Name:           "com.nemoworld.production.static",
		Description:    "Assets for production env",
		RegionID:       region3.ID,
		OrganizationID: organization1.ID,
	}
	fixtures.Buckets = append(fixtures.Buckets, bucket4)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Buckets {
			err := sqlxx.Save(ctx, dbx, fixtures.Buckets[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Buckets[i].ID)
		}
		return nil
	})

	directory1 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "A",
	}
	fixtures.Directories = append(fixtures.Directories, directory1)

	directory2 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "B",
	}
	fixtures.Directories = append(fixtures.Directories, directory2)

	directory3 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "C",
	}
	fixtures.Directories = append(fixtures.Directories, directory3)

	directory4 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "D",
	}
	fixtures.Directories = append(fixtures.Directories, directory4)

	directory5 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "E",
	}
	fixtures.Directories = append(fixtures.Directories, directory5)

	directory6 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "F",
	}
	fixtures.Directories = append(fixtures.Directories, directory6)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Directories {
			err := sqlxx.Save(ctx, dbx, fixtures.Directories[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Directories[i].ID)
		}
		return nil
	})

	directory7 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AA",
	}
	fixtures.Directories = append(fixtures.Directories, directory7)

	directory8 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AB",
	}
	fixtures.Directories = append(fixtures.Directories, directory8)

	directory9 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AC",
	}
	fixtures.Directories = append(fixtures.Directories, directory9)

	directory10 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AD",
	}
	fixtures.Directories = append(fixtures.Directories, directory10)

	directory11 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AE",
	}
	fixtures.Directories = append(fixtures.Directories, directory11)

	directory12 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory1.ID,
		},
		Path: "AF",
	}
	fixtures.Directories = append(fixtures.Directories, directory12)

	directory13 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory3.ID,
		},
		Path: "CB",
	}
	fixtures.Directories = append(fixtures.Directories, directory13)

	directory14 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory3.ID,
		},
		Path: "CC",
	}
	fixtures.Directories = append(fixtures.Directories, directory14)

	directory15 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		ParentID: sql.NullString{
			Valid:  true,
			String: directory3.ID,
		},
		Path: "CE",
	}
	fixtures.Directories = append(fixtures.Directories, directory15)

	directory16 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket1.ID,
		Path:           "deploy",
	}
	fixtures.Directories = append(fixtures.Directories, directory16)

	directory17 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket3.ID,
		Path:           "fonts",
	}
	fixtures.Directories = append(fixtures.Directories, directory17)

	directory18 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket3.ID,
		Path:           "js",
	}
	fixtures.Directories = append(fixtures.Directories, directory18)

	directory19 := &ExoDirectory{
		OrganizationID: organization1.ID,
		BucketID:       bucket3.ID,
		Path:           "css",
	}
	fixtures.Directories = append(fixtures.Directories, directory19)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Directories {
			if fixtures.Directories[i].ID != "" {
				continue
			}
			err := sqlxx.Save(ctx, dbx, fixtures.Directories[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Directories[i].ID)
		}
		return nil
	})

	file1 := &ExoFile{
		Path:           "aae95479d44ea1049741b15a5d0a91.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory14.ID,
		OrganizationID: organization1.ID,
		UserID:         user14.ID,
	}
	fixtures.Files = append(fixtures.Files, file1)

	file2 := &ExoFile{
		Path:           "ac4786d9f32740acffe4891ced59187.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory14.ID,
		OrganizationID: organization1.ID,
		UserID:         user16.ID,
	}
	fixtures.Files = append(fixtures.Files, file2)

	file3 := &ExoFile{
		Path:           "9f9edba51f17c13710fb22011fad1d2.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory10.ID,
		OrganizationID: organization1.ID,
		UserID:         user10.ID,
	}
	fixtures.Files = append(fixtures.Files, file3)

	file4 := &ExoFile{
		Path:           "fcde3c87c4d9db92348528b0fb09335.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory12.ID,
		OrganizationID: organization1.ID,
		UserID:         user17.ID,
	}
	fixtures.Files = append(fixtures.Files, file4)

	file5 := &ExoFile{
		Path:           "d60bed35b17618884a33b39a946c366.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory12.ID,
		OrganizationID: organization1.ID,
		UserID:         user5.ID,
	}
	fixtures.Files = append(fixtures.Files, file5)

	file6 := &ExoFile{
		Path:           "7e23423599035878e4744496ea234a6.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory12.ID,
		OrganizationID: organization1.ID,
		UserID:         user4.ID,
	}
	fixtures.Files = append(fixtures.Files, file6)

	file7 := &ExoFile{
		Path:           "69bc68e8b05ef0bb04fa7afe7600645.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory12.ID,
		OrganizationID: organization1.ID,
		UserID:         user11.ID,
	}
	fixtures.Files = append(fixtures.Files, file7)

	file8 := &ExoFile{
		Path:           "73bc153f32c881101a436e516422720.db",
		BucketID:       bucket1.ID,
		DirectoryID:    directory8.ID,
		OrganizationID: organization1.ID,
		UserID:         user13.ID,
	}
	fixtures.Files = append(fixtures.Files, file8)

	file9 := &ExoFile{
		Path:           "nemoworld-api-latest",
		BucketID:       bucket1.ID,
		DirectoryID:    directory16.ID,
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
	}
	fixtures.Files = append(fixtures.Files, file9)

	file10 := &ExoFile{
		Path:           "font-awesome.ttf",
		BucketID:       bucket3.ID,
		DirectoryID:    directory17.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file10)

	file11 := &ExoFile{
		Path:           "source-code-pro.ttf",
		BucketID:       bucket3.ID,
		DirectoryID:    directory17.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file11)

	file12 := &ExoFile{
		Path:           "account.js",
		BucketID:       bucket3.ID,
		DirectoryID:    directory18.ID,
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
	}
	fixtures.Files = append(fixtures.Files, file12)

	file13 := &ExoFile{
		Path:           "admin.js",
		BucketID:       bucket3.ID,
		DirectoryID:    directory18.ID,
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
	}
	fixtures.Files = append(fixtures.Files, file13)

	file14 := &ExoFile{
		Path:           "api.js",
		BucketID:       bucket3.ID,
		DirectoryID:    directory18.ID,
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
	}
	fixtures.Files = append(fixtures.Files, file14)

	file15 := &ExoFile{
		Path:           "index.js",
		BucketID:       bucket3.ID,
		DirectoryID:    directory18.ID,
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
	}
	fixtures.Files = append(fixtures.Files, file15)

	file16 := &ExoFile{
		Path:           "account.css",
		BucketID:       bucket3.ID,
		DirectoryID:    directory19.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file16)

	file17 := &ExoFile{
		Path:           "admin.css",
		BucketID:       bucket3.ID,
		DirectoryID:    directory19.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file17)

	file18 := &ExoFile{
		Path:           "common.css",
		BucketID:       bucket3.ID,
		DirectoryID:    directory19.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file18)

	file19 := &ExoFile{
		Path:           "index.css",
		BucketID:       bucket3.ID,
		DirectoryID:    directory19.ID,
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
	}
	fixtures.Files = append(fixtures.Files, file19)

	file20 := &ExoFile{
		Path:           "nemoworld-supervisor-latest",
		BucketID:       bucket1.ID,
		DirectoryID:    directory16.ID,
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
	}
	fixtures.Files = append(fixtures.Files, file20)

	file21 := &ExoFile{
		Path:           "nemoworld-front-latest",
		BucketID:       bucket1.ID,
		DirectoryID:    directory16.ID,
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
	}
	fixtures.Files = append(fixtures.Files, file21)

	file22 := &ExoFile{
		Path:           "nemoworld-worker-latest",
		BucketID:       bucket1.ID,
		DirectoryID:    directory16.ID,
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
	}
	fixtures.Files = append(fixtures.Files, file22)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Files {
			err := sqlxx.Save(ctx, dbx, fixtures.Files[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Files[i].ID)
		}
		return nil
	})

	mode1 := &ExoChunkMode{
		Mode: "rwx",
	}
	fixtures.Modes = append(fixtures.Modes, mode1)

	mode2 := &ExoChunkMode{
		Mode: "r-x",
	}
	fixtures.Modes = append(fixtures.Modes, mode2)

	mode3 := &ExoChunkMode{
		Mode: "r--",
	}
	fixtures.Modes = append(fixtures.Modes, mode3)

	mode4 := &ExoChunkMode{
		Mode: "rw-",
	}
	fixtures.Modes = append(fixtures.Modes, mode4)

	mode5 := &ExoChunkMode{
		Mode: "udp-stream",
	}
	fixtures.Modes = append(fixtures.Modes, mode5)

	mode6 := &ExoChunkMode{
		Mode: "tcp-stream",
	}
	fixtures.Modes = append(fixtures.Modes, mode6)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Modes {
			err := sqlxx.Save(ctx, dbx, fixtures.Modes[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Modes[i].ID)
		}
		return nil
	})

	chunk1 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user14.ID,
		FileID:         file1.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"ca2e9ab8ee532ac3ed61125dbafe50ecbe0db1810b4fca91f31b3733ba71bb5f",
			"460c2b3394a4b039b6a1d536832c1d133c67cf456363d745d0b4b7795cc3459d",
			"1e272e396cc6024713be6835ce29d4b2a541105f8bd66f72df8d6da4770e7ee5",
			"f7bdaf5ffef17f0288d7486e0579ec6892084a9de41b942d8fa9fd730ef2a862",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk1)

	chunk2 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user14.ID,
		FileID:         file1.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"dcff2761d6cf115159a312ad24dd33b6cd6cfa4693645836ce6e9bd9b1be9154",
			"49e569066658e1341fc6bc93421202d637892c655f85d87141a411d23f1de7e9",
			"e3ad2985983fb71d6850a3c53699882d03e0ca198765c91f6db8c97eba690801",
			"6394d4c36ac9ef6a69f09922e3b485f397da9a08127ab7fb63e68229ac2e415c",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk2)

	chunk3 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user14.ID,
		FileID:         file1.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"3fcd7cc7e5a80e537364e9fd01b181c5d5f4d7c9ed49b18a40e65556aa565a09",
			"b7283b165bb0c19ee762e79c6df94592ace276545bfec2230949cfb4ed0ca166",
			"02b6521d4409b3d7b3b8588127f4a2b00e9ad2fa787bbe74d960ea18d4c65a9b",
			"3eb3b7b37035a4e528608728b48b5e5048371340812e15c6225b36fee1a9d481",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk3)

	chunk4 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user16.ID,
		FileID:         file2.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"2f3deec366a1e02d095ec84fbfa670b9b1f0a142f5b7eee9969819d3e7443341",
			"539372e9d1d0440a78c0cc0c965212e96e0368793d4668c4cc76cdf3357122c8",
			"377d1a36565f1904cfb14c5685ead00ec38e01509a0f2cfc49958b1392b8982c",
			"50d3e108d50774aa55f9de2430fb9e5a8f475dceb909eb6d28c7646ac3415b9c",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk4)

	chunk5 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user16.ID,
		FileID:         file2.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"8e267c00e6c7491307bd7d3b7a7a3e312e86bca93a9b80ff93595493053c082a",
			"73450dd00ee722dadb5e0af0dd8a02f6230ce9c42d5e2fe20c23d15c34f9b743",
			"58dd05511372e2d0ff0ca76c9f7245571cec5b4c9d280c46fe797c6dc0db98ac",
			"d7c2cff65bd1bec5d23b964fbed63b5102d0399a7d8a6d8ef5d5082e111ad8ac",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk5)

	chunk6 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user16.ID,
		FileID:         file2.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"b65c6d37f5ee831c8b23ea5f39683711c3ed6e858c3dc40ea5b74d268897ae22",
			"1f8ce95338e25485b21a41094328a79ae49842d6adb24fde192fd71d176367ad",
			"e76c94df54d5fad6c23ed517025c5df913283b42c496521ca6b2b2649c70991d",
			"51e8a441a8098254fee39831060fe7805b6465a581f4845767d8c625357cc544",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk6)

	chunk7 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user10.ID,
		FileID:         file3.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"ff188812f84193e0cd27e3cd12112cef8531e7bc6498ac28b29f5e73b844f061",
			"c95e458ce9d6b4501ea9afedf152df74bdc4b9ce8f9919bc23be9fbc107f0a47",
			"48abfa3b567d1919b54961d607fd31d4c2fedbee48b52691e884b25c725e22a9",
			"2f0079224820884374bc6c2ddee51d3b11b33fcde1755d06202705eb4ba2c093",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk7)

	chunk8 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user10.ID,
		FileID:         file3.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"a473b7f8e7e7b4d7222885ed92ef268512ee2a173330cc4399c106fefc7a3455",
			"bd610e2dff86a3829ed4a4b5fd90b5cbcee08945ac808143c4bf91bb9b6d4706",
			"459b81be4e49ef87d9f107ae48352a7fbae47194ca6d64533a92f740a3eb644c",
			"fd3c8cad884257c3a33e770eb86bd6ed270b70f7e1f023e460772ec9add805cd",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk8)

	chunk9 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user10.ID,
		FileID:         file3.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"24ee0d8af51b89dae519f19bbfdaa10eafc1f13047a9c97c27191a13487da3f3",
			"5b9689acf65377aaabd6d6900c2c1147327b830dcf9ef741fcde359e2de34dcb",
			"66bd37ba56571fa216a216924156ba605ca5d8cf7927d40e542884fb7f3cec4f",
			"2330351206dc1ede9b1a4a2672a8ee5aa33da5c457c9db0037ccafa8902d5aaf",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk9)

	chunk10 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user17.ID,
		FileID:         file4.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"66b0d9467477dc1268b82e568aa1b38e2223d16c27e19ea53770214fe669185e",
			"220a78b7571e509cf5593c0f763959677ddecf47da1c01fe19459ebf4d7817d3",
			"ef955679707e26f2928cd643cc88e546a0aebe956e6fead2f9c6f5c23aa82d43",
			"564c02189c639379eeccf317c2d26c8e379878199efabc05f360cbca256170e5",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk10)

	chunk11 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user17.ID,
		FileID:         file4.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"431a855f852d6a9fb66785cb819a7d964f3241d0c854eef1ea92beae3ba7ce71",
			"55cbb46624681f7999bcfe6ac4dac0c5de98e3d4b281842fd81eb551e2d56052",
			"4c4d0231bb06877aeaeac636537af9c830a55626ef80f9e18611158d6792acb6",
			"104e60238011662129b9880881f0552dbe18ca4bd175425c5ee5b8f8769ad550",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk11)

	chunk12 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user17.ID,
		FileID:         file4.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"d165cde9856ecf92f294027efa84d6548ff98e2bec6dd8ac38fbd8919f6a0a38",
			"b3ff45d15a7f4d81b26040360b77335850a35a93e18644eeec59d9ca04557a16",
			"4c70784dd32c3791209bd229d1be8525ec6621928d3a709f207e2a2d35551d39",
			"5346a9c919a06a3d16a55d3fa3be67ca16eaf8293cdb3b52a43e3ad7d4ef3c64",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk12)

	chunk13 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user5.ID,
		FileID:         file5.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"054e4a8698badd34de1686e087ca48ee0bd5225824ff1a066c547ca42d27d815",
			"366c913495ca074c489dbb841e5abcd04e9e90649a6a6a45e98263f252e2c0ac",
			"964cb6ba01173ef95011bed3c903adc0c81332baeea5a610bde148fbb4cafd60",
			"a5639dc59c6282b7fcb5307304c6061ac6c16cb42b721ffb01316b2e6daf76bf",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk13)

	chunk14 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user5.ID,
		FileID:         file5.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"db26234f24bc91e3da4dd6fb485f051aa3ea7297c860327cf53ff30fdd7d6bf4",
			"92acfbdb80dc283b0832b26dfa3470d7ebb5937be740fb4ab07cfeb972817e5e",
			"8c9a423b0065d86a71211277e59b0ff94d3c3202ca091f8b5cded30a6856f246",
			"30cf1887fab252a8b973216c3f4f5cd04b7c34bc7ffbd4a70f647d6c098ccf07",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk14)

	chunk15 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user5.ID,
		FileID:         file5.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"dbc071f238ee00268dba4bb7e0288c25ecadc8d2aca94a60339427c430f0bb8b",
			"3836042faee6f7234d48f4b029d0d46055518df7674273f6b0af9f9c7374751e",
			"a5c4b948387b46d917df61cae4dfec456bc1a9b339c8861777403e7a292e723b",
			"74d351fdfb7055bcb981f680698ee16bb4afd919ead2b5bed639f03b67dff329",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk15)

	chunk16 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user4.ID,
		FileID:         file6.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"6f5971ffbfe59ad67992639fd4f6101a514c300930e38a3219379eb8b170d897",
			"3081918d3e7c35e9ec7ecdba958b9f07f6b9d89bb55b3db86efbc142f26c57c0",
			"eac92e18bea9c6118bfd44ae3c70bc9fe702c265efed53a261c703ad9dc28545",
			"a5f21dd0d44a2a02b26b1159f2205e1d2471c27de297523c31f3b4407d67259d",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk16)

	chunk17 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user4.ID,
		FileID:         file6.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"35a3c6a1abe731e93ea2e13b18f068b88caf76625340792288addfed9c3221a9",
			"d7961904d5892400055e579dc338d1916b45cc4a6860c3779e2d561867604284",
			"13f3932f52ef068a95ac34d135d3847b09ca16928f8bc2520b2ca7594fa11635",
			"dc06d648ddd52b126ea8e67cb3f9efc91033678a44e04b41f82736eedf4d83c9",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk17)

	chunk18 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user4.ID,
		FileID:         file6.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"529480cccfa8f3ab5b1513b2e9bcaf48d9ed6f62eedd00646d483675f6753275",
			"54805a6215e3557f9f216a5aa1ca6dc0c1d627e27482b56d285ddc328d9b6d90",
			"1895f287a4764c6e21589f747b388db9f3f7d80be284929823d142d29c5fb74a",
			"4c04a85306ec608eafce2a880ab519d01b3b04f5e4b44febb139e9a8efdd8b88",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk18)

	chunk19 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user11.ID,
		FileID:         file7.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"99bc0b6a09d3e578de8db5844765ffd43cc0b95836db8addc570f3d61d38b7a8",
			"c0d70a71412d98bb3ac95ffbecb77b2cc9c64008739d54e41e5cf850c1f59628",
			"1cd7f2ac3038ff88d9407b95c36dd10c02e9a75c510f224f4a4c372125d3be34",
			"1960fe47742db6a877dbee5add40038f9a13e9c3f74f0f888c1e9f83f96b926f",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk19)

	chunk20 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user11.ID,
		FileID:         file7.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"012c544663972b0074221cc08341ed6e0818a4de482d043cda23cd3b5a46c8e6",
			"07cb6114b65604770c5aaae0063c5a34529ac306552af676d7a9d68b955baa81",
			"399ecfc6114ed2d5bb466eca504140e64a0de62c6c43918f2a6bf8ec93d4983b",
			"3954569d24e1886bb992356558b845d8b1a3267012fb952f795a6cb84ad3f337",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk20)

	chunk21 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user11.ID,
		FileID:         file7.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"d996a3f969ce8a443d64dbfcf83d358a86cc49fb5b67ac58498be03fa648159d",
			"f6d85492cd4d0d2f09045c2326b8e4b7c53238a1c5318c7500ea5bb32f6e3149",
			"62cd34e35a544f21fef7adbd2ae489b2f6b2b579d1b2182ec5e8069c2f95103e",
			"7c6b901c1d7ac07cb18f0a894948d743bd65db709dc46e33402ffea4be36106d",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk21)

	chunk22 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user13.ID,
		FileID:         file8.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"85e8dafb0637d400f94e3f8cedd9522aa071060f9411b68e28c5d0d0ef80f6f4",
			"f4779b54b74076a1294f4de0439fab4b7f7c01eba72663a2278518daf747743d",
			"59a32dceae6db785b6c4737632512470ebc9366c448fdb72bbc08a77e317102d",
			"0298393fc20ac6e51430104c1b99bc6e5ac8c172d8f086a33be9fdbae48e1815",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk22)

	chunk23 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user13.ID,
		FileID:         file8.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"3fd92b0e76ed8c818f0890596aa51e998e6ae78aa792e0fe3d57646d2dd6014a",
			"4cd7020b2541db4deeb935118df9205aff4e0cb8e9b50476d502f20b431b801a",
			"8df3e2d6bf731bd385ce07abb33412ca12ec54a066d348f1aa09ef73a290791e",
			"24a0637e78dc8969f5fd675257c6bca6a1049b7fffeb1c3a31e3bbedc6f57875",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk23)

	chunk24 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user13.ID,
		FileID:         file8.ID,
		ModeID:         mode4.ID,
		Bytes: fmt.Sprint(
			"e272e70882901a3e4b7428493e9fd8329a1e5cff4197200ad18b4994f2f395d8",
			"23ab45af11f3da32926b17e3d5ae1c6a8726a485dd28d98b825194007486779a",
			"e2906d9f3a09c35beecc774e663a06b8f3f3dc3b08ff5065ab540f379efe7e5b",
			"32a392818c9f501941bbeb530d15a3687877d6e25973a2d2839cb71b078defaf",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk24)

	chunk25 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file9.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"6b57c20b075abca4eca4b668b8fd6606b1892f5338a4b870737821fd0feae3db",
			"270db604af9b2ab9bd9b37d1df41af88ab01400a93e63d2e1f571aa46913afd5",
			"a28665a55785bc6d58ad23ed5462de85fc1e9f4c4e208741359853fcdd1ef837",
			"bb47b3b03eb72087eb946f0400a23c2cb6facd9d3c61a20222cbc74867216eb6",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk25)

	chunk26 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file9.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"b957c4330f341b0d757d8e85e774964bddd4d6bb132966c28d5fe537e45352b6",
			"b596750c920079192da955a5aa34eab26700cf9a671e1e005af867a98b06c6da",
			"89ceb7194c11bc2c6f7c34c2e5f3532276bda7232bac5f48f255e338e10bc352",
			"d7b2c58d7fbff52d39764bb3ef3995e22592d68717f9baaed385c97485eb97d6",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk26)

	chunk27 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file9.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"1353aa20b8af8398b98bf30f79fd3639a1a6d67b0a70d325b4ca6e5116719fe1",
			"b697abb0c453de2b114771d752b45950ab8b390b6e30f160d358b43e2cc61060",
			"10ba7313484a40f755c498f22a2fde85589f948efc69961439215edb0920585c",
			"67dfde88220b3e5675296f8abcfaa1268537e08503766a0d65ec66a3efad7ce2",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk27)

	chunk28 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file9.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"95096f2269c9d6305f9d94eb26fcd515e1622d56003eef4c0070b272a5a1949b",
			"c9a0bf84e93940cfa37768389acdfd38a3eba593968b799420da3f77407a37ba",
			"667e9d93806c4a93614ccf405de5edac064e17e40a81af034cc10c087c3f99c7",
			"3713a4735e189e525ec0f7f2a0d891a3e6d9de4cc61bb583722aec7270107fe0",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk28)

	chunk29 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file10.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"9d7024e7b8b8f43042b8d1c7ed5fb0dc2503e19fb863aa327399357f3c0fe265",
			"7605c344a33354693b8ef4bd6eb8a2b064e38ad6beff03ce8c898d473c11b462",
			"d50715b55b32b4521a60ff0c04b5e1484ebb73b0d35352fe9f1923b5521d6f1c",
			"8d202f3ba33aa2b704fdc8b6eff6fea17732816da9f5d4d470e30fd1021a92c2",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk29)

	chunk30 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file11.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"7c5c6d9451eecf41e96760f49ef7881adf9bd7352183ad43346de60dd39ae4f9",
			"febac191e7a659f44aac4ad4d8662b793c037d55eb4c20c6c3e00b7f1ea40c06",
			"3d1ddc4120526bf352718b96106e9d6bba54c2449b95e8cee17f5d299c5d49d7",
			"43fb196f8743ebf527829665d37b902635ebfcd34cae442c66418b025dca794e",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk30)

	chunk31 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file12.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"dd74c6c33827fe0211616fc188af4d6f19096bc94e1b1c2cb200dad696bad87d",
			"9f0def42367da23cc1700ed8df33d8888432a5de239c38318bfe22d31a50d1d1",
			"bc261b147fec4df1bb3c30c7e862663756b6bd2c55a73c9fd4c0054471995c7a",
			"34270b907e37ce899e0853169b7c6cb0d56cafb19edb16edccc14c484be4edf3",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk31)

	chunk32 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file12.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"e5bfd5c217e3ebc27e3b3a38e8c08027078476ab67ded461658f759b38bb2a47",
			"6bb1572338024e2e3136616075321e8f6ccac575234d2df63ac42cedd11a2c1e",
			"a0ae57baf72137b8dc1286369bde55a241ccfcf0571ccada80dd1ee8d868438b",
			"bece5c025b6332c099dde127ade55ff1d04f378b12a4209c5bce985943788576",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk32)

	chunk33 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file13.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"52774ab0fd2e806530c3cfb33c9188a2b88dd3ed5c65ba9ccd57a63a67d6d07a",
			"cb8b33220d72549280b8baa5859198c50678b60cc29c955f00eb30cc06f8aaea",
			"f7aab60ff7a1e95a2788d3345f3e7ca806a5114156e713a0132d5bfd2b8be222",
			"5730de124b49e760ceec93caa63e414119a562c3559e0914957719602e13f268",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk33)

	chunk34 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file13.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"06a72e82739c1f156c4a72db998d0ff75ac7ce57f5df4426c20eb747e581c73c",
			"b1227981b4a6d5343b14fd7d2015866b428af7d2f6b3ca3c43d7dabdb24bc081",
			"087e6a9c4e5e7092effbf82a4689fd780f8406a34cd1294ebe9d506b7fc0b427",
			"f365696c7f1baa07eff4ed6154f175a878f393a2343a25f7cfe561f89adbb5ee",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk34)

	chunk35 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file14.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"5b3b31038eae66f292036170a4c489a04861b73b8c5a6154a56253a267182074",
			"3bd77e6175d5ed5adaa8e68273d35db817d8e21a54e00c127b331252a8f8e4fe",
			"b0cea834fa60fa3eb7de809b9c9df8cb1c2a5b2ac89b5b48c0aec4128ae2831f",
			"35419044d2b5b81d439a215a298de51ff5de152d4d04737b5c755fd6000d662f",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk35)

	chunk36 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file14.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"f3e8a120378102f0b5b04f535b697595a8b55c310be688bfdad99b830918d976",
			"eceb0c7433db1e45b967545527ce0991c11f517c5e1161b6f864a6df7b170b3b",
			"13e638b273c4bf1307a211a3060a51750c9812a7f507caffa276b61e96b3058c",
			"f5a59692847ef9d1b8a8d93405307cfaf993a05486354db441fbde371f3914ac",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk36)

	chunk37 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file15.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"52774ab0fd2e806530c3cfb33c9188a2b88dd3ed5c65ba9ccd57a63a67d6d07a",
			"cb8b33220d72549280b8baa5859198c50678b60cc29c955f00eb30cc06f8aaea",
			"f7aab60ff7a1e95a2788d3345f3e7ca806a5114156e713a0132d5bfd2b8be222",
			"5730de124b49e760ceec93caa63e414119a562c3559e0914957719602e13f268",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk37)

	chunk38 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user6.ID,
		FileID:         file15.ID,
		ModeID:         mode2.ID,
		Bytes: fmt.Sprint(
			"041ba28963b17aa4f09728d0d624b645727abfe1a8d05ac44d89e0e7b7c6260d",
			"bf15a4d61900eef7ac1e51d08abf82ce5d5767db5529f081a71f4e9e59f38f86",
			"6492bac76d1df317674eb9a01bddf7e6d9a99fefd7676ddc0807a14dd9f5886d",
			"3353a4dc6ebd009134791049026190f56fec32f039d5f7ec0740bf0b7cdbbcc8",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk38)

	chunk39 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file16.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"d72009be8f86ee4c32997370e27eb9db2b5a4896dde818cc5b0cede719b382eb",
			"50df8ba4e93627610cf6c6945d855c8d1a3d9816c677b4919278993cf8141274",
			"a4a59d52b5c01ab61a62cf0040cf6a77ae65a65a2db9f24b4f91db015a1890d9",
			"e7bdba85a553982209f7bf7876aae260dde689bd13f7c451752c5a306880c233",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk39)

	chunk40 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file17.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"f98875543c1e01f36021a434dd51cf68c198a9d20b7e32d43499f4bc87a093de",
			"b2f4de9a540718dbba9be4889919bf01bdcc41f38866d64a077b5594490d88ba",
			"9d7ae8d5f52c68142ba1c3119d9fdd93f5c552be54553b1e2d447049493b57eb",
			"d3ae2720f89964c55faa29d7f43dd211cb256ba3aa66c0d4ef1af27815daadce",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk40)

	chunk41 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file18.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"5da9618c4a55a46ead6d385179fa978643af21d4960114064efe5dfc221d05f9",
			"8f6a45787e40f023cbd43bb3c96289026d289590d8c68f939eb4f5e1ae8281da",
			"8d19a9041bc673f56544aaf2fc17e908d876536b5e17b494b9f437c449e7b0cf",
			"193a9ee490fd66b60220da8c884edf5a1786a816a6dbf3c984bc9d67f503eb76",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk41)

	chunk42 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user8.ID,
		FileID:         file19.ID,
		ModeID:         mode3.ID,
		Bytes: fmt.Sprint(
			"1308a143b732bda61e2641d66e467f3930c47e1b983843d21c7042ca3ae977e8",
			"9e59026de02ffc4191a7655d0fa06b883b9aff6554edae2c8f4223fd38ffb4dc",
			"a9dc33e43721acb747d74cbf222cfa81eb0799aac59852a28d9b0326ade3d2bb",
			"8ebab91a31623613cfaf68993f4eb86b9d27e16e81a1c59eb2eb02009ac1ffc2",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk42)

	chunk43 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file20.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"c11481d8b7a166d5ff5ca95fa99992de310aa47f76aa8fafc24600a144a63e0e",
			"a8492de65061de31cfd39b45e8c94d01925a4cc0e85528ec97b1705a843fb6b9",
			"f1ab71d3b6b1f2cc5f7a982677d41634b667602e1b6e2c82ddbcba7fd980f9d0",
			"82df649dfd77a48deba7e34aad0f7e40ea0b7bc6dbb011f38d212b72db4982be",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk43)

	chunk44 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file20.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"4a075b111a1945ff9ceee41e358ff0fdf13a4c8345e00abfe4154ab652b70baf",
			"398415bda0600b23e0eb1d0e9113310f5a06f480dcaab308b0e39b635b4ea8ad",
			"584c96efa9f791180fd8f956dab8cfcac08721a94c88384a21599a3da9e87ba7",
			"50b970a54c9c8dbbf9806f01e195976f64ebed30a6c181326473849b0fee133d",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk44)

	chunk45 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file20.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"c8e95d3b518ebfae2a49576b1672f6a7e3cd33e05984b1c56877b3f71f916ca2",
			"cd014b7c394103a8a4c87d21809a1d79d4ee2c72b39264d84e7be30458880429",
			"feac70d51fae7e0203ff5b68fe9d4e00ccf31bc9d601b77c1ce43e51e434a16e",
			"c22be66f7edfc0e9c0843c343a1a7cfba2b2aaa9f7897b41da18e02b92811688",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk45)

	chunk46 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file21.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"3d16394c70da72b10c953d93e9ab6c0823b69c8f3de1f6fdcf88c2285fda282c",
			"8e5be87de17cc60bb1e79b26c60a63b580189f5a10036727b2938a01d98f8c83",
			"fdf1edcb09ef81a646a5c9d8e3c0f8392ed6a7f7ec4ffd8470d6770c04fa81d2",
			"c893f1494f54d00741b1dbf1013ddd4a8738072449712ea61d10157da06bfab6",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk46)

	chunk47 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file21.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"8b7a0ffce702fb710fc07cc4ccebc8f782c6c3845eb9f8d681a7b03c628eebcc",
			"811dc693476c1c1f65846ead336e5693a009a2cddab2f77453fc550c6f860747",
			"1e52cc0c0c4a6b742cb1b3a3001b18adabd46ab789a2a06444caa3ededdc2dcc",
			"fb56950150dc4c08799ed6cf72dbf5f5c53880776150145a3dcead5783d1585d",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk47)

	chunk48 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file21.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"8678244a91779cc39c51ec747810ac099911c31990112082e7d595e24df9b5ae",
			"f9aa8c09a426f93fc0053479c3fe80bb2d213018464377d06ab9d9c1cec84650",
			"d7fff4a9fb6e26a6dfd8b5ac14cd9ed69deff87fbe24789b4f8d5f6e3f96ccb5",
			"ac3cf53544799bbe4e536fc4135ce265cc1d1f72fc7b202a1c067bfb087beb4d",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk48)

	chunk49 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file22.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"f6adbfa9a39343323b34686b5f16fc5f196a89c19225709a1a0432b812354c06",
			"f75322f82147ec11260156f1f66761510a07b87c62b5f581154f35b3a8862c4c",
			"997e3b475ee168238e5158762efbf598e52f4be5ad64f04a177a30c0972d0729",
			"1bf8682af57f3a67e0ef7680b20d639230786c126b41dcce7d32719e86df9242",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk49)

	chunk50 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file22.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"b81a07c3a65f6deaf99199caa983ea1b35e72f8eb81b4791e96c9cd1b7c0eb00",
			"a50eb4573233f30ec5eb42ead07a07c0e1a0e1fc2da1d56b7208731bcd546e9f",
			"f6dd528a7c3f982932a9402e3c8abef482e52da9df7859dcafc6f71f74ae5022",
			"bd7af4ceff04e6494f571a879e0f35282d8c48e52deddfa6d1583d064d930356",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk50)

	chunk51 := &ExoChunk{
		OrganizationID: organization1.ID,
		UserID:         user19.ID,
		FileID:         file22.ID,
		ModeID:         mode1.ID,
		Bytes: fmt.Sprint(
			"e65db167fab1ec513f9bcb6d57db1677cfad2f7709d5824149f4ab6b90f0b577",
			"4754ac3a931d54fda097fdc82b53cb009e4b3df7c82f4215f1dda618baadf9f1",
			"b73ee2223a64a15564776ce6b6133aced5b211b39ac322874412069540a3dc33",
			"2ca3122da843da5255711466dd655cac10065fef137dffdd682f11030b262d23",
		),
	}
	fixtures.Chunks = append(fixtures.Chunks, chunk51)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Chunks {
			err := sqlxx.Save(ctx, dbx, fixtures.Chunks[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Chunks[i].Hash)
		}
		return nil
	})

	signature1 := &ExoChunkSignature{
		ChunkID: chunk25.Hash,
		Bytes:   "f4de76a30a229af6001e5abc3ad11ec9ab1e0fef11902cd601afd8ef175bc60a",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature1)

	signature2 := &ExoChunkSignature{
		ChunkID: chunk26.Hash,
		Bytes:   "cf5f0db262636a7e6aed1ca2294bb5c2a1602f5e3a9787906f10e14f640e8156",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature2)

	signature3 := &ExoChunkSignature{
		ChunkID: chunk27.Hash,
		Bytes:   "b14411f946f177d69779f3dc2be1bbbd0cb8ecc07a529a0f479ed5ba95d52726",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature3)

	signature4 := &ExoChunkSignature{
		ChunkID: chunk28.Hash,
		Bytes:   "d0313cf6130c81cd121a1dd554ea2e86f0ef8c055d52fcd08a6e560b7c85e1a7",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature4)

	signature5 := &ExoChunkSignature{
		ChunkID: chunk43.Hash,
		Bytes:   "12643bb7f36e9425f19403504639c0b0f801e10ee3a9fb9b5c0ae4841b1b0c45",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature5)

	signature6 := &ExoChunkSignature{
		ChunkID: chunk44.Hash,
		Bytes:   "62dc59c52c1567b620b4b42eee2b913d46b19b29b984aa6504932e89457593ab",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature6)

	signature7 := &ExoChunkSignature{
		ChunkID: chunk45.Hash,
		Bytes:   "1ad9a59a19679198102a942b8346efbe157aab8a2a3e78aed3eb98b6199e3bd0",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature7)

	signature8 := &ExoChunkSignature{
		ChunkID: chunk46.Hash,
		Bytes:   "34829c657f138d12fdcfe2312bb5168c8cd832784d69385c1903372a08cf0b30",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature8)

	signature9 := &ExoChunkSignature{
		ChunkID: chunk47.Hash,
		Bytes:   "ca3c7ffd33bd55aa77e7c4652bce95e982f2706ef465537b96625d1eeb566f96",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature9)

	signature10 := &ExoChunkSignature{
		ChunkID: chunk48.Hash,
		Bytes:   "71c207c8f2e58354d018163f2fec62b005481ce9f761cfea4cb61bdb9e1f9d97",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature10)

	signature11 := &ExoChunkSignature{
		ChunkID: chunk49.Hash,
		Bytes:   "779f873f13bf8cf022f08fa07ca84f1cdc69eacf34d16c48e59edb419db59a0b",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature11)

	signature12 := &ExoChunkSignature{
		ChunkID: chunk50.Hash,
		Bytes:   "5a9715c56168b3195aad9422fe30b4a39a73642ed7282527a5b216409957b5bc",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature12)

	signature13 := &ExoChunkSignature{
		ChunkID: chunk51.Hash,
		Bytes:   "f27455efab1b070edce688ff96f4bd2e00f8921df9eb6b961f266ab1835e94fd",
	}
	fixtures.Signatures = append(fixtures.Signatures, signature13)

	_ = sqlxx.Transaction(driver, func(dbx sqlxx.Driver) error {
		for i := range fixtures.Signatures {
			err := sqlxx.Save(ctx, dbx, fixtures.Signatures[i])
			is.NoError(err)
			is.NotEmpty(fixtures.Signatures[i].ID)
		}
		return nil
	})

	return fixtures
}

type ExoOrganization struct {
	// Columns
	ID      string `sqlxx:"column:id,pk:ulid"`
	Name    string `sqlxx:"column:name"`
	Website string `sqlxx:"column:website"`
}

func (ExoOrganization) TableName() string {
	return "exo_organization"
}

type ExoUser struct {
	// Columns
	ID        string `sqlxx:"column:id,pk:ulid"`
	Email     string `sqlxx:"column:email"`
	Password  string `sqlxx:"column:password"`
	Country   string `sqlxx:"column:country"`
	Locale    string `sqlxx:"column:locale"`
	ProfileID string `sqlxx:"column:profile_id,fk:exo_profile"`
	// Relationships
	Group   *ExoGroup
	Profile *ExoProfile
}

func (ExoUser) TableName() string {
	return "exo_user"
}

type ExoGroup struct {
	// Columns
	ID             string `sqlxx:"column:id,pk:ulid"`
	Role           string `sqlxx:"column:role"`
	UserID         string `sqlxx:"column:user_id,fk:exo_user"`
	OrganizationID string `sqlxx:"column:organization_id,fk:exo_organization"`
	// Relationships
	User         *ExoUser
	Organization *ExoOrganization
}

func (ExoGroup) TableName() string {
	return "exo_group"
}

type ExoProfile struct {
	// Columns
	ID          string         `sqlxx:"column:id,pk:ulid"`
	FirstName   string         `sqlxx:"column:first_name"`
	LastName    string         `sqlxx:"column:last_name"`
	AvatarID    sql.NullString `sqlxx:"column:avatar_id,fk:exo_avatar"`
	DisplayName sql.NullString `sqlxx:"column:display_name"`
	Description sql.NullString `sqlxx:"column:description"`
	Website     sql.NullString `sqlxx:"column:website"`
	// Relationships
	Avatar *ExoAvatar
}

func (ExoProfile) TableName() string {
	return "exo_profile"
}

type ExoAvatar struct {
	// Columns
	ID       string `sqlxx:"column:id,pk:ulid"`
	URL      string `sqlxx:"column:url"`
	Path     string `sqlxx:"column:path"`
	MimeType string `sqlxx:"column:mime_type"`
}

func (ExoAvatar) TableName() string {
	return "exo_avatar"
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
	ID             string `sqlxx:"column:id,pk:ulid"`
	Name           string `sqlxx:"column:name"`
	Description    string `sqlxx:"column:description"`
	RegionID       string `sqlxx:"column:region_id,fk:exo_region"`
	OrganizationID string `sqlxx:"column:organization_id,fk:exo_organization"`
	// Relationships
	Region      ExoRegion
	Directories []ExoDirectory
	Files       []ExoFile
}

func (ExoBucket) TableName() string {
	return "exo_bucket"
}

type ExoDirectory struct {
	// Columns
	ID             string         `sqlxx:"column:id,pk:ulid"`
	Path           string         `sqlxx:"column:path"`
	OrganizationID string         `sqlxx:"column:organization_id,fk:exo_organization"`
	BucketID       string         `sqlxx:"column:bucket_id,fk:exo_bucket"`
	ParentID       sql.NullString `sqlxx:"column:parent_id,fk:exo_directory"`
	// Relationships
	Directories []*ExoDirectory
	Files       []*ExoFile
}

func (ExoDirectory) TableName() string {
	return "exo_directory"
}

type ExoFile struct {
	// Columns
	ID             string `sqlxx:"column:id,pk:ulid"`
	Path           string `sqlxx:"column:path"`
	OrganizationID string `sqlxx:"column:organization_id,fk:exo_organization"`
	UserID         string `sqlxx:"column:user_id,fk:exo_user"`
	BucketID       string `sqlxx:"column:bucket_id,fk:exo_bucket"`
	DirectoryID    string `sqlxx:"column:directory_id,fk:exo_directory"`
	// Relationships
	Chunks []ExoChunk
}

func (ExoFile) TableName() string {
	return "exo_file"
}

type ExoChunk struct {
	// Columns
	Hash           string `sqlxx:"column:hash,pk:ulid"`
	Bytes          string `sqlxx:"column:bytes"`
	OrganizationID string `sqlxx:"column:organization_id,fk:exo_organization"`
	UserID         string `sqlxx:"column:user_id,fk:exo_user"`
	ModeID         string `sqlxx:"column:mode_id,fk:exo_chunk_mode"`
	FileID         string `sqlxx:"column:file_id,fk:exo_file"`
	// Relationships
	File      *ExoFile
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

		--
		-- Zootopia schema
		--

		DROP TABLE IF EXISTS ztp_human CASCADE;
		DROP TABLE IF EXISTS ztp_package CASCADE;
		DROP TABLE IF EXISTS ztp_bag CASCADE;
		DROP TABLE IF EXISTS ztp_owl CASCADE;
		DROP TABLE IF EXISTS ztp_cat CASCADE;
		DROP TABLE IF EXISTS ztp_meow CASCADE;
		DROP TABLE IF EXISTS ztp_group CASCADE;
		DROP TABLE IF EXISTS ztp_center CASCADE;

		--
		-- Object storage application
		--

		DROP TABLE IF EXISTS exo_chunk_signature CASCADE;
		DROP TABLE IF EXISTS exo_chunk CASCADE;
		DROP TABLE IF EXISTS exo_chunk_mode CASCADE;
		DROP TABLE IF EXISTS exo_file CASCADE;
		DROP TABLE IF EXISTS exo_directory CASCADE;
		DROP TABLE IF EXISTS exo_bucket CASCADE;
		DROP TABLE IF EXISTS exo_region CASCADE;
		DROP TABLE IF EXISTS exo_group CASCADE;
		DROP TABLE IF EXISTS exo_user CASCADE;
		DROP TABLE IF EXISTS exo_profile CASCADE;
		DROP TABLE IF EXISTS exo_avatar CASCADE;
		DROP TABLE IF EXISTS exo_organization CASCADE;

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

		CREATE TABLE exo_organization (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			name               VARCHAR(255) NOT NULL,
			website            VARCHAR(2048) NOT NULL
		);
		CREATE TABLE exo_avatar (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			url                VARCHAR(2048) NOT NULL,
			path               VARCHAR(255) NOT NULL,
			mime_type          VARCHAR(64) NOT NULL
		);
		CREATE TABLE exo_profile (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			first_name         VARCHAR(255) NOT NULL,
			last_name          VARCHAR(255) NOT NULL,
			avatar_id          VARCHAR(26) REFERENCES exo_avatar(id) ON DELETE RESTRICT,
			display_name       VARCHAR(255),
			description        VARCHAR(2048),
			website            VARCHAR(2048)
		);
		CREATE TABLE exo_user (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			email              VARCHAR(255) NOT NULL,
			password           VARCHAR(128) NOT NULL,
			country            VARCHAR(2) NOT NULL,
			locale             VARCHAR(5) NOT NULL,
			profile_id         VARCHAR(26) NOT NULL REFERENCES exo_profile(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_group (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			role               VARCHAR(5) NOT NULL,
			user_id            VARCHAR(26) NOT NULL REFERENCES exo_user(id) ON DELETE RESTRICT,
			organization_id    VARCHAR(26) NOT NULL REFERENCES exo_organization(id) ON DELETE RESTRICT
		);
		CREATE UNIQUE INDEX exo_group_unique ON exo_group (user_id, organization_id);
		CREATE TABLE exo_region (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			name               VARCHAR(255) NOT NULL,
			hostname           VARCHAR(2048) NOT NULL
		);
		CREATE TABLE exo_bucket (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			name               VARCHAR(512) NOT NULL,
			description        VARCHAR(2048) NOT NULL,
			region_id          VARCHAR(26) NOT NULL REFERENCES exo_region(id) ON DELETE RESTRICT,
			organization_id    VARCHAR(26) NOT NULL REFERENCES exo_organization(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_directory (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			path               VARCHAR(512) NOT NULL,
			organization_id    VARCHAR(26) NOT NULL REFERENCES exo_organization(id) ON DELETE RESTRICT,
			bucket_id          VARCHAR(26) NOT NULL REFERENCES exo_bucket(id) ON DELETE RESTRICT,
			parent_id          VARCHAR(26) REFERENCES exo_directory(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_file (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			path               VARCHAR(512) NOT NULL,
			organization_id    VARCHAR(26) NOT NULL REFERENCES exo_organization(id) ON DELETE RESTRICT,
			user_id            VARCHAR(26) NOT NULL REFERENCES exo_user(id) ON DELETE RESTRICT,
			bucket_id          VARCHAR(26) NOT NULL REFERENCES exo_bucket(id) ON DELETE RESTRICT,
			directory_id       VARCHAR(26) NOT NULL REFERENCES exo_directory(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_chunk_mode (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			mode               VARCHAR(255) NOT NULL
		);
		CREATE TABLE exo_chunk (
			hash               VARCHAR(26) PRIMARY KEY NOT NULL,
			bytes              VARCHAR(2048) NOT NULL,
			organization_id    VARCHAR(26) NOT NULL REFERENCES exo_organization(id) ON DELETE RESTRICT,
			user_id            VARCHAR(26) NOT NULL REFERENCES exo_user(id) ON DELETE RESTRICT,
			mode_id            VARCHAR(26) NOT NULL REFERENCES exo_chunk_mode(id) ON DELETE RESTRICT,
			file_id            VARCHAR(26) NOT NULL REFERENCES exo_file(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_chunk_signature (
			id                 VARCHAR(26) PRIMARY KEY NOT NULL,
			chunk_id           VARCHAR(26) NOT NULL REFERENCES exo_chunk(hash) ON DELETE RESTRICT,
			bytes              VARCHAR(2048) NOT NULL
		);

	`)
}
