package sqlxx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestPreloadX_ExoUser_MultiLevel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		fixtures := GenerateExoCloudFixtures(ctx, driver, is)

		user1 := fixtures.Users[1]

		err := sqlxx.Preload(ctx, driver, user1,
			"Profile",
			"Group",
			"Profile.Avatar",
			"Group.User",
			"Group.Organization",
		)
		is.NoError(err)

		fmt.Println("::1-1", user1.ProfileID)
		fmt.Println("::1-2", user1.Profile.ID)
		fmt.Println("::1-1", user1.Profile.AvatarID)
		fmt.Println("::1-2", user1.Profile.Avatar)

	})
}

func TestPreloadX_ExoRegion_MultiLevel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Empty(fixtures.Regions[0].Buckets)
			is.Empty(fixtures.Regions[1].Buckets)
			is.Empty(fixtures.Regions[2].Buckets)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			region1 := fixtures.Regions[0]
			region2 := fixtures.Regions[1]
			region3 := fixtures.Regions[2]
			regions := []*ExoRegion{
				region1, region2, region3,
			}

			err := sqlxx.Preload(ctx, driver, &regions,
				"Buckets",
				"Buckets.Region",
				"Buckets.Directories",
				"Buckets.Files",
				"Buckets.Files.Chunks",
				"Buckets.Directories.Directories",
				"Buckets.Directories.Directories.Files",
				"Buckets.Directories.Directories.Files.Chunks",
				"Buckets.Directories.Files",
				"Buckets.Directories.Files.Chunks",
			)
			is.NoError(err)

		}
	})
}
