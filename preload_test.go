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

func TestPreload_ExoRegion_One(t *testing.T) {
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

			err := sqlxx.Preload(ctx, driver, region1, "Buckets")
			is.NoError(err)
			is.NotNil(region1.Buckets)
			is.NotEmpty((*region1.Buckets))
			is.Len((*region1.Buckets), 2)
			is.Contains((*region1.Buckets), *fixtures.Buckets[0])
			is.Contains((*region1.Buckets), *fixtures.Buckets[1])

			region2 := fixtures.Regions[1]

			err = sqlxx.Preload(ctx, driver, region2, "Buckets")
			is.NoError(err)
			is.NotNil(region2.Buckets)
			is.Empty((*region2.Buckets))
			is.Len((*region2.Buckets), 0)

			region3 := fixtures.Regions[2]

			err = sqlxx.Preload(ctx, driver, region3, "Buckets")
			is.NoError(err)
			is.NotNil(region3.Buckets)
			is.NotEmpty((*region3.Buckets))
			is.Len((*region3.Buckets), 2)
			is.Contains((*region3.Buckets), *fixtures.Buckets[2])
			is.Contains((*region3.Buckets), *fixtures.Buckets[3])

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			region1 := fixtures.Regions[0]

			err := sqlxx.Preload(ctx, driver, &region1, "Buckets")
			is.NoError(err)
			is.NotNil(region1.Buckets)
			is.NotEmpty((*region1.Buckets))
			is.Len((*region1.Buckets), 2)
			is.Contains((*region1.Buckets), *fixtures.Buckets[0])
			is.Contains((*region1.Buckets), *fixtures.Buckets[1])

			region2 := fixtures.Regions[1]

			err = sqlxx.Preload(ctx, driver, &region2, "Buckets")
			is.NoError(err)
			is.NotNil(region2.Buckets)
			is.Empty((*region2.Buckets))
			is.Len((*region2.Buckets), 0)

			region3 := fixtures.Regions[2]

			err = sqlxx.Preload(ctx, driver, &region3, "Buckets")
			is.NoError(err)
			is.NotNil(region3.Buckets)
			is.NotEmpty((*region3.Buckets))
			is.Len((*region3.Buckets), 2)
			is.Contains((*region3.Buckets), *fixtures.Buckets[2])
			is.Contains((*region3.Buckets), *fixtures.Buckets[3])

		}
	})
}

func TestPreload_ExoRegion_Many(t *testing.T) {
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

			regions := []ExoRegion{
				*fixtures.Regions[0],
				*fixtures.Regions[1],
				*fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets")
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Contains((*regions[0].Buckets), *fixtures.Buckets[0])
			is.Contains((*regions[0].Buckets), *fixtures.Buckets[1])

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Contains((*regions[2].Buckets), *fixtures.Buckets[2])
			is.Contains((*regions[2].Buckets), *fixtures.Buckets[3])

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := []*ExoRegion{
				fixtures.Regions[0],
				fixtures.Regions[1],
				fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets")
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Contains((*regions[0].Buckets), *fixtures.Buckets[0])
			is.Contains((*regions[0].Buckets), *fixtures.Buckets[1])

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Contains((*regions[2].Buckets), *fixtures.Buckets[2])
			is.Contains((*regions[2].Buckets), *fixtures.Buckets[3])

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := &[]ExoRegion{
				*fixtures.Regions[0],
				*fixtures.Regions[1],
				*fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets")
			is.NoError(err)
			is.Len((*regions), 3)
			is.Equal(fixtures.Regions[0].ID, (*regions)[0].ID)
			is.Equal(fixtures.Regions[1].ID, (*regions)[1].ID)
			is.Equal(fixtures.Regions[2].ID, (*regions)[2].ID)

			is.NotNil((*regions)[0].Buckets)
			is.NotEmpty((*(*regions)[0].Buckets))
			is.Len((*(*regions)[0].Buckets), 2)
			is.Contains((*(*regions)[0].Buckets), *fixtures.Buckets[0])
			is.Contains((*(*regions)[0].Buckets), *fixtures.Buckets[1])

			is.NotNil((*regions)[1].Buckets)
			is.Empty((*(*regions)[1].Buckets))
			is.Len((*(*regions)[1].Buckets), 0)

			is.NotNil((*regions)[2].Buckets)
			is.NotEmpty((*(*regions)[2].Buckets))
			is.Len((*(*regions)[2].Buckets), 2)
			is.Contains((*(*regions)[2].Buckets), *fixtures.Buckets[2])
			is.Contains((*(*regions)[2].Buckets), *fixtures.Buckets[3])

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := &[]*ExoRegion{
				fixtures.Regions[0],
				fixtures.Regions[1],
				fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets")
			is.NoError(err)
			is.Len((*regions), 3)
			is.Equal(fixtures.Regions[0].ID, (*regions)[0].ID)
			is.Equal(fixtures.Regions[1].ID, (*regions)[1].ID)
			is.Equal(fixtures.Regions[2].ID, (*regions)[2].ID)

			is.NotNil((*regions)[0].Buckets)
			is.NotEmpty((*(*regions)[0].Buckets))
			is.Len((*(*regions)[0].Buckets), 2)
			is.Contains((*(*regions)[0].Buckets), *fixtures.Buckets[0])
			is.Contains((*(*regions)[0].Buckets), *fixtures.Buckets[1])

			is.NotNil((*regions)[1].Buckets)
			is.Empty((*(*regions)[1].Buckets))
			is.Len((*(*regions)[1].Buckets), 0)

			is.NotNil((*regions)[2].Buckets)
			is.NotEmpty((*(*regions)[2].Buckets))
			is.Len((*(*regions)[2].Buckets), 2)
			is.Contains((*(*regions)[2].Buckets), *fixtures.Buckets[2])
			is.Contains((*(*regions)[2].Buckets), *fixtures.Buckets[3])

		}
	})
}

func TestPreload_ExoBucket_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Empty(fixtures.Buckets[0].Region)
			is.Empty(fixtures.Buckets[1].Region)
			is.Empty(fixtures.Buckets[2].Region)
			is.Empty(fixtures.Buckets[3].Region)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			bucket1 := fixtures.Buckets[0]

			err := sqlxx.Preload(ctx, driver, bucket1, "Region")
			is.NoError(err)
			is.NotEmpty(bucket1.Region)
			is.Equal(fixtures.Regions[0].ID, bucket1.Region.ID)
			is.Equal(fixtures.Regions[0].Name, bucket1.Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, bucket1.Region.Hostname)

			bucket2 := fixtures.Buckets[1]

			err = sqlxx.Preload(ctx, driver, bucket2, "Region")
			is.NoError(err)
			is.NotEmpty(bucket2.Region)
			is.Equal(fixtures.Regions[0].ID, bucket2.Region.ID)
			is.Equal(fixtures.Regions[0].Name, bucket2.Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, bucket2.Region.Hostname)

			bucket3 := fixtures.Buckets[2]

			err = sqlxx.Preload(ctx, driver, bucket3, "Region")
			is.NoError(err)
			is.NotEmpty(bucket3.Region)
			is.Equal(fixtures.Regions[2].ID, bucket3.Region.ID)
			is.Equal(fixtures.Regions[2].Name, bucket3.Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, bucket3.Region.Hostname)

			bucket4 := fixtures.Buckets[3]

			err = sqlxx.Preload(ctx, driver, bucket4, "Region")
			is.NoError(err)
			is.NotEmpty(bucket4.Region)
			is.Equal(fixtures.Regions[2].ID, bucket4.Region.ID)
			is.Equal(fixtures.Regions[2].Name, bucket4.Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, bucket4.Region.Hostname)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			bucket1 := fixtures.Buckets[0]

			err := sqlxx.Preload(ctx, driver, &bucket1, "Region")
			is.NoError(err)
			is.NotEmpty(bucket1.Region)
			is.Equal(fixtures.Regions[0].ID, bucket1.Region.ID)
			is.Equal(fixtures.Regions[0].Name, bucket1.Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, bucket1.Region.Hostname)

			bucket2 := fixtures.Buckets[1]

			err = sqlxx.Preload(ctx, driver, &bucket2, "Region")
			is.NoError(err)
			is.NotEmpty(bucket2.Region)
			is.Equal(fixtures.Regions[0].ID, bucket2.Region.ID)
			is.Equal(fixtures.Regions[0].Name, bucket2.Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, bucket2.Region.Hostname)

			bucket3 := fixtures.Buckets[2]

			err = sqlxx.Preload(ctx, driver, &bucket3, "Region")
			is.NoError(err)
			is.NotEmpty(bucket3.Region)
			is.Equal(fixtures.Regions[2].ID, bucket3.Region.ID)
			is.Equal(fixtures.Regions[2].Name, bucket3.Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, bucket3.Region.Hostname)

			bucket4 := fixtures.Buckets[3]

			err = sqlxx.Preload(ctx, driver, &bucket4, "Region")
			is.NoError(err)
			is.NotEmpty(bucket4.Region)
			is.Equal(fixtures.Regions[2].ID, bucket4.Region.ID)
			is.Equal(fixtures.Regions[2].Name, bucket4.Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, bucket4.Region.Hostname)

		}
	})
}

func TestPreload_ExoBucket_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Empty(fixtures.Buckets[0].Region)
			is.Empty(fixtures.Buckets[1].Region)
			is.Empty(fixtures.Buckets[2].Region)
			is.Empty(fixtures.Buckets[3].Region)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			buckets := []ExoBucket{
				*fixtures.Buckets[0],
				*fixtures.Buckets[1],
				*fixtures.Buckets[2],
				*fixtures.Buckets[3],
			}

			err := sqlxx.Preload(ctx, driver, &buckets, "Region")
			is.NoError(err)
			is.Len(buckets, 4)
			is.Equal(fixtures.Buckets[0].ID, buckets[0].ID)
			is.Equal(fixtures.Buckets[1].ID, buckets[1].ID)
			is.Equal(fixtures.Buckets[2].ID, buckets[2].ID)
			is.Equal(fixtures.Buckets[3].ID, buckets[3].ID)

			is.NotEmpty(buckets[0].Region)
			is.Equal(fixtures.Regions[0].ID, buckets[0].Region.ID)
			is.Equal(fixtures.Regions[0].Name, buckets[0].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, buckets[0].Region.Hostname)

			is.NotEmpty(buckets[1].Region)
			is.Equal(fixtures.Regions[0].ID, buckets[1].Region.ID)
			is.Equal(fixtures.Regions[0].Name, buckets[1].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, buckets[1].Region.Hostname)

			is.NotEmpty(buckets[2].Region)
			is.Equal(fixtures.Regions[2].ID, buckets[2].Region.ID)
			is.Equal(fixtures.Regions[2].Name, buckets[2].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, buckets[2].Region.Hostname)

			is.NotEmpty(buckets[3].Region)
			is.Equal(fixtures.Regions[2].ID, buckets[3].Region.ID)
			is.Equal(fixtures.Regions[2].Name, buckets[3].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, buckets[3].Region.Hostname)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			buckets := []*ExoBucket{
				fixtures.Buckets[0],
				fixtures.Buckets[1],
				fixtures.Buckets[2],
				fixtures.Buckets[3],
			}

			err := sqlxx.Preload(ctx, driver, &buckets, "Region")
			is.NoError(err)
			is.Len(buckets, 4)
			is.Equal(fixtures.Buckets[0].ID, buckets[0].ID)
			is.Equal(fixtures.Buckets[1].ID, buckets[1].ID)
			is.Equal(fixtures.Buckets[2].ID, buckets[2].ID)
			is.Equal(fixtures.Buckets[3].ID, buckets[3].ID)

			is.NotEmpty(buckets[0].Region)
			is.Equal(fixtures.Regions[0].ID, buckets[0].Region.ID)
			is.Equal(fixtures.Regions[0].Name, buckets[0].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, buckets[0].Region.Hostname)

			is.NotEmpty(buckets[1].Region)
			is.Equal(fixtures.Regions[0].ID, buckets[1].Region.ID)
			is.Equal(fixtures.Regions[0].Name, buckets[1].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, buckets[1].Region.Hostname)

			is.NotEmpty(buckets[2].Region)
			is.Equal(fixtures.Regions[2].ID, buckets[2].Region.ID)
			is.Equal(fixtures.Regions[2].Name, buckets[2].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, buckets[2].Region.Hostname)

			is.NotEmpty(buckets[3].Region)
			is.Equal(fixtures.Regions[2].ID, buckets[3].Region.ID)
			is.Equal(fixtures.Regions[2].Name, buckets[3].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, buckets[3].Region.Hostname)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			buckets := &[]ExoBucket{
				*fixtures.Buckets[0],
				*fixtures.Buckets[1],
				*fixtures.Buckets[2],
				*fixtures.Buckets[3],
			}

			err := sqlxx.Preload(ctx, driver, &buckets, "Region")
			is.NoError(err)
			is.Len((*buckets), 4)
			is.Equal(fixtures.Buckets[0].ID, (*buckets)[0].ID)
			is.Equal(fixtures.Buckets[1].ID, (*buckets)[1].ID)
			is.Equal(fixtures.Buckets[2].ID, (*buckets)[2].ID)
			is.Equal(fixtures.Buckets[3].ID, (*buckets)[3].ID)

			is.NotEmpty((*buckets)[0].Region)
			is.Equal(fixtures.Regions[0].ID, (*buckets)[0].Region.ID)
			is.Equal(fixtures.Regions[0].Name, (*buckets)[0].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, (*buckets)[0].Region.Hostname)

			is.NotEmpty((*buckets)[1].Region)
			is.Equal(fixtures.Regions[0].ID, (*buckets)[1].Region.ID)
			is.Equal(fixtures.Regions[0].Name, (*buckets)[1].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, (*buckets)[1].Region.Hostname)

			is.NotEmpty((*buckets)[2].Region)
			is.Equal(fixtures.Regions[2].ID, (*buckets)[2].Region.ID)
			is.Equal(fixtures.Regions[2].Name, (*buckets)[2].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, (*buckets)[2].Region.Hostname)

			is.NotEmpty((*buckets)[3].Region)
			is.Equal(fixtures.Regions[2].ID, (*buckets)[3].Region.ID)
			is.Equal(fixtures.Regions[2].Name, (*buckets)[3].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, (*buckets)[3].Region.Hostname)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			buckets := &[]*ExoBucket{
				fixtures.Buckets[0],
				fixtures.Buckets[1],
				fixtures.Buckets[2],
				fixtures.Buckets[3],
			}

			err := sqlxx.Preload(ctx, driver, &buckets, "Region")
			is.NoError(err)
			is.Len((*buckets), 4)
			is.Equal(fixtures.Buckets[0].ID, (*buckets)[0].ID)
			is.Equal(fixtures.Buckets[1].ID, (*buckets)[1].ID)
			is.Equal(fixtures.Buckets[2].ID, (*buckets)[2].ID)
			is.Equal(fixtures.Buckets[3].ID, (*buckets)[3].ID)

			is.NotEmpty((*buckets)[0].Region)
			is.Equal(fixtures.Regions[0].ID, (*buckets)[0].Region.ID)
			is.Equal(fixtures.Regions[0].Name, (*buckets)[0].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, (*buckets)[0].Region.Hostname)

			is.NotEmpty((*buckets)[1].Region)
			is.Equal(fixtures.Regions[0].ID, (*buckets)[1].Region.ID)
			is.Equal(fixtures.Regions[0].Name, (*buckets)[1].Region.Name)
			is.Equal(fixtures.Regions[0].Hostname, (*buckets)[1].Region.Hostname)

			is.NotEmpty((*buckets)[2].Region)
			is.Equal(fixtures.Regions[2].ID, (*buckets)[2].Region.ID)
			is.Equal(fixtures.Regions[2].Name, (*buckets)[2].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, (*buckets)[2].Region.Hostname)

			is.NotEmpty((*buckets)[3].Region)
			is.Equal(fixtures.Regions[2].ID, (*buckets)[3].Region.ID)
			is.Equal(fixtures.Regions[2].Name, (*buckets)[3].Region.Name)
			is.Equal(fixtures.Regions[2].Hostname, (*buckets)[3].Region.Hostname)

		}
	})
}

func TestPreload_ExoDirectory_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Nil(fixtures.Directories[0].Files)
			is.Nil(fixtures.Directories[1].Files)
			is.Nil(fixtures.Directories[2].Files)
			is.Nil(fixtures.Directories[11].Files)
			is.Nil(fixtures.Directories[13].Files)
			is.Nil(fixtures.Directories[15].Files)
			is.Nil(fixtures.Directories[17].Files)
			is.Nil(fixtures.Directories[0].Directories)
			is.Nil(fixtures.Directories[1].Directories)
			is.Nil(fixtures.Directories[2].Directories)
			is.Nil(fixtures.Directories[11].Directories)
			is.Nil(fixtures.Directories[13].Directories)
			is.Nil(fixtures.Directories[15].Directories)
			is.Nil(fixtures.Directories[17].Directories)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directory1 := fixtures.Directories[0]

			err := sqlxx.Preload(ctx, driver, directory1, "Files")
			is.NoError(err)
			is.Empty(directory1.Files)
			is.Empty(directory1.Directories)

			err = sqlxx.Preload(ctx, driver, directory1, "Directories")
			is.NoError(err)
			is.NotEmpty(directory1.Directories)
			is.Len(directory1.Directories, 6)
			is.Equal(fixtures.Directories[6].ID, directory1.Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, directory1.Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, directory1.Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, directory1.Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, directory1.Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, directory1.Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, directory1.Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, directory1.Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, directory1.Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, directory1.Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, directory1.Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, directory1.Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, directory1.Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, directory1.Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, directory1.Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, directory1.Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, directory1.Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, directory1.Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, directory1.Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, directory1.Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, directory1.Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, directory1.Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, directory1.Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, directory1.Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, directory1.Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, directory1.Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, directory1.Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, directory1.Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, directory1.Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, directory1.Directories[5].ParentID)

			directory2 := fixtures.Directories[1]

			err = sqlxx.Preload(ctx, driver, directory2, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory2.Files)
			is.Empty(directory2.Directories)

			directory3 := fixtures.Directories[2]

			err = sqlxx.Preload(ctx, driver, directory3, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory3.Files)
			is.NotEmpty(directory3.Directories)
			is.Len(directory3.Directories, 3)
			is.Equal(fixtures.Directories[12].ID, directory3.Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, directory3.Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, directory3.Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, directory3.Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, directory3.Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, directory3.Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, directory3.Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, directory3.Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, directory3.Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, directory3.Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, directory3.Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, directory3.Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, directory3.Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, directory3.Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, directory3.Directories[2].ParentID)

			directory4 := fixtures.Directories[11]

			err = sqlxx.Preload(ctx, driver, directory4, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory4.Directories)
			is.NotEmpty(directory4.Files)
			is.Len(directory4.Files, 4)
			is.Equal(fixtures.Files[3].ID, directory4.Files[0].ID)
			is.Equal(fixtures.Files[3].Path, directory4.Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, directory4.Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, directory4.Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, directory4.Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, directory4.Files[1].ID)
			is.Equal(fixtures.Files[4].Path, directory4.Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, directory4.Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, directory4.Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, directory4.Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, directory4.Files[2].ID)
			is.Equal(fixtures.Files[5].Path, directory4.Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, directory4.Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, directory4.Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, directory4.Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, directory4.Files[3].ID)
			is.Equal(fixtures.Files[6].Path, directory4.Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, directory4.Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, directory4.Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, directory4.Files[3].UserID)

			directory5 := fixtures.Directories[13]

			err = sqlxx.Preload(ctx, driver, directory5, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory5.Directories)
			is.NotEmpty(directory5.Files)
			is.Len(directory5.Files, 2)
			is.Equal(fixtures.Files[0].ID, directory5.Files[0].ID)
			is.Equal(fixtures.Files[0].Path, directory5.Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, directory5.Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, directory5.Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, directory5.Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, directory5.Files[1].ID)
			is.Equal(fixtures.Files[1].Path, directory5.Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, directory5.Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, directory5.Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, directory5.Files[1].UserID)

			directory6 := fixtures.Directories[15]

			err = sqlxx.Preload(ctx, driver, directory6, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory6.Directories)
			is.NotEmpty(directory6.Files)
			is.Len(directory6.Files, 4)
			is.Equal(fixtures.Files[8].ID, directory6.Files[0].ID)
			is.Equal(fixtures.Files[8].Path, directory6.Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, directory6.Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, directory6.Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, directory6.Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, directory6.Files[1].ID)
			is.Equal(fixtures.Files[19].Path, directory6.Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, directory6.Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, directory6.Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, directory6.Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, directory6.Files[2].ID)
			is.Equal(fixtures.Files[20].Path, directory6.Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, directory6.Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, directory6.Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, directory6.Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, directory6.Files[3].ID)
			is.Equal(fixtures.Files[21].Path, directory6.Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, directory6.Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, directory6.Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, directory6.Files[3].UserID)

			directory7 := fixtures.Directories[17]

			err = sqlxx.Preload(ctx, driver, directory7, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory7.Directories)
			is.NotEmpty(directory7.Files)
			is.Len(directory7.Files, 4)
			is.Equal(fixtures.Files[11].ID, directory7.Files[0].ID)
			is.Equal(fixtures.Files[11].Path, directory7.Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, directory7.Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, directory7.Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, directory7.Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, directory7.Files[1].ID)
			is.Equal(fixtures.Files[12].Path, directory7.Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, directory7.Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, directory7.Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, directory7.Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, directory7.Files[2].ID)
			is.Equal(fixtures.Files[13].Path, directory7.Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, directory7.Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, directory7.Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, directory7.Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, directory7.Files[3].ID)
			is.Equal(fixtures.Files[14].Path, directory7.Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, directory7.Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, directory7.Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, directory7.Files[3].UserID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directory1 := fixtures.Directories[0]

			err := sqlxx.Preload(ctx, driver, &directory1, "Files")
			is.NoError(err)
			is.Empty(directory1.Files)
			is.Empty(directory1.Directories)

			err = sqlxx.Preload(ctx, driver, &directory1, "Directories")
			is.NoError(err)
			is.NotEmpty(directory1.Directories)
			is.Len(directory1.Directories, 6)
			is.Equal(fixtures.Directories[6].ID, directory1.Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, directory1.Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, directory1.Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, directory1.Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, directory1.Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, directory1.Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, directory1.Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, directory1.Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, directory1.Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, directory1.Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, directory1.Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, directory1.Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, directory1.Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, directory1.Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, directory1.Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, directory1.Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, directory1.Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, directory1.Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, directory1.Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, directory1.Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, directory1.Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, directory1.Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, directory1.Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, directory1.Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, directory1.Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, directory1.Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, directory1.Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, directory1.Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, directory1.Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, directory1.Directories[5].ParentID)

			directory2 := fixtures.Directories[1]

			err = sqlxx.Preload(ctx, driver, &directory2, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory2.Files)
			is.Empty(directory2.Directories)

			directory3 := fixtures.Directories[2]

			err = sqlxx.Preload(ctx, driver, &directory3, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory3.Files)
			is.NotEmpty(directory3.Directories)
			is.Len(directory3.Directories, 3)
			is.Equal(fixtures.Directories[12].ID, directory3.Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, directory3.Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, directory3.Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, directory3.Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, directory3.Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, directory3.Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, directory3.Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, directory3.Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, directory3.Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, directory3.Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, directory3.Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, directory3.Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, directory3.Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, directory3.Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, directory3.Directories[2].ParentID)

			directory4 := fixtures.Directories[11]

			err = sqlxx.Preload(ctx, driver, &directory4, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory4.Directories)
			is.NotEmpty(directory4.Files)
			is.Len(directory4.Files, 4)
			is.Equal(fixtures.Files[3].ID, directory4.Files[0].ID)
			is.Equal(fixtures.Files[3].Path, directory4.Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, directory4.Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, directory4.Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, directory4.Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, directory4.Files[1].ID)
			is.Equal(fixtures.Files[4].Path, directory4.Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, directory4.Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, directory4.Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, directory4.Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, directory4.Files[2].ID)
			is.Equal(fixtures.Files[5].Path, directory4.Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, directory4.Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, directory4.Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, directory4.Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, directory4.Files[3].ID)
			is.Equal(fixtures.Files[6].Path, directory4.Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, directory4.Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, directory4.Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, directory4.Files[3].UserID)

			directory5 := fixtures.Directories[13]

			err = sqlxx.Preload(ctx, driver, &directory5, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory5.Directories)
			is.NotEmpty(directory5.Files)
			is.Len(directory5.Files, 2)
			is.Equal(fixtures.Files[0].ID, directory5.Files[0].ID)
			is.Equal(fixtures.Files[0].Path, directory5.Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, directory5.Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, directory5.Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, directory5.Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, directory5.Files[1].ID)
			is.Equal(fixtures.Files[1].Path, directory5.Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, directory5.Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, directory5.Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, directory5.Files[1].UserID)

			directory6 := fixtures.Directories[15]

			err = sqlxx.Preload(ctx, driver, &directory6, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory6.Directories)
			is.NotEmpty(directory6.Files)
			is.Len(directory6.Files, 4)
			is.Equal(fixtures.Files[8].ID, directory6.Files[0].ID)
			is.Equal(fixtures.Files[8].Path, directory6.Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, directory6.Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, directory6.Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, directory6.Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, directory6.Files[1].ID)
			is.Equal(fixtures.Files[19].Path, directory6.Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, directory6.Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, directory6.Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, directory6.Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, directory6.Files[2].ID)
			is.Equal(fixtures.Files[20].Path, directory6.Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, directory6.Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, directory6.Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, directory6.Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, directory6.Files[3].ID)
			is.Equal(fixtures.Files[21].Path, directory6.Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, directory6.Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, directory6.Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, directory6.Files[3].UserID)

			directory7 := fixtures.Directories[17]

			err = sqlxx.Preload(ctx, driver, &directory7, "Files", "Directories")
			is.NoError(err)
			is.Empty(directory7.Directories)
			is.NotEmpty(directory7.Files)
			is.Len(directory7.Files, 4)
			is.Equal(fixtures.Files[11].ID, directory7.Files[0].ID)
			is.Equal(fixtures.Files[11].Path, directory7.Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, directory7.Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, directory7.Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, directory7.Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, directory7.Files[1].ID)
			is.Equal(fixtures.Files[12].Path, directory7.Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, directory7.Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, directory7.Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, directory7.Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, directory7.Files[2].ID)
			is.Equal(fixtures.Files[13].Path, directory7.Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, directory7.Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, directory7.Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, directory7.Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, directory7.Files[3].ID)
			is.Equal(fixtures.Files[14].Path, directory7.Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, directory7.Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, directory7.Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, directory7.Files[3].UserID)

		}
	})
}

func TestPreload_ExoDirectory_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Nil(fixtures.Directories[0].Files)
			is.Nil(fixtures.Directories[1].Files)
			is.Nil(fixtures.Directories[2].Files)
			is.Nil(fixtures.Directories[11].Files)
			is.Nil(fixtures.Directories[13].Files)
			is.Nil(fixtures.Directories[15].Files)
			is.Nil(fixtures.Directories[17].Files)
			is.Nil(fixtures.Directories[0].Directories)
			is.Nil(fixtures.Directories[1].Directories)
			is.Nil(fixtures.Directories[2].Directories)
			is.Nil(fixtures.Directories[11].Directories)
			is.Nil(fixtures.Directories[13].Directories)
			is.Nil(fixtures.Directories[15].Directories)
			is.Nil(fixtures.Directories[17].Directories)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directories := []ExoDirectory{
				*fixtures.Directories[0],
				*fixtures.Directories[1],
				*fixtures.Directories[2],
				*fixtures.Directories[11],
				*fixtures.Directories[13],
				*fixtures.Directories[15],
				*fixtures.Directories[17],
			}

			err := sqlxx.Preload(ctx, driver, &directories, "Files", "Directories")
			is.NoError(err)
			is.Len(directories, 7)
			is.Equal(fixtures.Directories[0].ID, directories[0].ID)
			is.Equal(fixtures.Directories[1].ID, directories[1].ID)
			is.Equal(fixtures.Directories[2].ID, directories[2].ID)
			is.Equal(fixtures.Directories[11].ID, directories[3].ID)
			is.Equal(fixtures.Directories[13].ID, directories[4].ID)
			is.Equal(fixtures.Directories[15].ID, directories[5].ID)
			is.Equal(fixtures.Directories[17].ID, directories[6].ID)

			is.Empty(directories[0].Files)
			is.NotEmpty(directories[0].Directories)
			is.Len(directories[0].Directories, 6)
			is.Equal(fixtures.Directories[6].ID, directories[0].Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, directories[0].Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, directories[0].Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, directories[0].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, directories[0].Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, directories[0].Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, directories[0].Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, directories[0].Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, directories[0].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, directories[0].Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, directories[0].Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, directories[0].Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, directories[0].Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, directories[0].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, directories[0].Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, directories[0].Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, directories[0].Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, directories[0].Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, directories[0].Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, directories[0].Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, directories[0].Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, directories[0].Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, directories[0].Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, directories[0].Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, directories[0].Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, directories[0].Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, directories[0].Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, directories[0].Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, directories[0].Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, directories[0].Directories[5].ParentID)

			is.Empty(directories[1].Files)
			is.Empty(directories[1].Directories)

			is.Empty(directories[2].Files)
			is.NotEmpty(directories[2].Directories)
			is.Len(directories[2].Directories, 3)
			is.Equal(fixtures.Directories[12].ID, directories[2].Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, directories[2].Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, directories[2].Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, directories[2].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, directories[2].Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, directories[2].Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, directories[2].Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, directories[2].Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, directories[2].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, directories[2].Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, directories[2].Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, directories[2].Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, directories[2].Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, directories[2].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, directories[2].Directories[2].ParentID)

			is.Empty(directories[3].Directories)
			is.NotEmpty(directories[3].Files)
			is.Len(directories[3].Files, 4)
			is.Equal(fixtures.Files[3].ID, directories[3].Files[0].ID)
			is.Equal(fixtures.Files[3].Path, directories[3].Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, directories[3].Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, directories[3].Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, directories[3].Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, directories[3].Files[1].ID)
			is.Equal(fixtures.Files[4].Path, directories[3].Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, directories[3].Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, directories[3].Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, directories[3].Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, directories[3].Files[2].ID)
			is.Equal(fixtures.Files[5].Path, directories[3].Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, directories[3].Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, directories[3].Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, directories[3].Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, directories[3].Files[3].ID)
			is.Equal(fixtures.Files[6].Path, directories[3].Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, directories[3].Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, directories[3].Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, directories[3].Files[3].UserID)

			is.Empty(directories[4].Directories)
			is.NotEmpty(directories[4].Files)
			is.Len(directories[4].Files, 2)
			is.Equal(fixtures.Files[0].ID, directories[4].Files[0].ID)
			is.Equal(fixtures.Files[0].Path, directories[4].Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, directories[4].Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, directories[4].Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, directories[4].Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, directories[4].Files[1].ID)
			is.Equal(fixtures.Files[1].Path, directories[4].Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, directories[4].Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, directories[4].Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, directories[4].Files[1].UserID)

			is.Empty(directories[5].Directories)
			is.NotEmpty(directories[5].Files)
			is.Len(directories[5].Files, 4)
			is.Equal(fixtures.Files[8].ID, directories[5].Files[0].ID)
			is.Equal(fixtures.Files[8].Path, directories[5].Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, directories[5].Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, directories[5].Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, directories[5].Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, directories[5].Files[1].ID)
			is.Equal(fixtures.Files[19].Path, directories[5].Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, directories[5].Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, directories[5].Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, directories[5].Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, directories[5].Files[2].ID)
			is.Equal(fixtures.Files[20].Path, directories[5].Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, directories[5].Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, directories[5].Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, directories[5].Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, directories[5].Files[3].ID)
			is.Equal(fixtures.Files[21].Path, directories[5].Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, directories[5].Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, directories[5].Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, directories[5].Files[3].UserID)

			is.Empty(directories[6].Directories)
			is.NotEmpty(directories[6].Files)
			is.Len(directories[6].Files, 4)
			is.Equal(fixtures.Files[11].ID, directories[6].Files[0].ID)
			is.Equal(fixtures.Files[11].Path, directories[6].Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, directories[6].Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, directories[6].Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, directories[6].Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, directories[6].Files[1].ID)
			is.Equal(fixtures.Files[12].Path, directories[6].Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, directories[6].Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, directories[6].Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, directories[6].Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, directories[6].Files[2].ID)
			is.Equal(fixtures.Files[13].Path, directories[6].Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, directories[6].Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, directories[6].Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, directories[6].Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, directories[6].Files[3].ID)
			is.Equal(fixtures.Files[14].Path, directories[6].Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, directories[6].Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, directories[6].Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, directories[6].Files[3].UserID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directories := []*ExoDirectory{
				fixtures.Directories[0],
				fixtures.Directories[1],
				fixtures.Directories[2],
				fixtures.Directories[11],
				fixtures.Directories[13],
				fixtures.Directories[15],
				fixtures.Directories[17],
			}

			err := sqlxx.Preload(ctx, driver, &directories, "Files", "Directories")
			is.NoError(err)
			is.Len(directories, 7)
			is.Equal(fixtures.Directories[0].ID, directories[0].ID)
			is.Equal(fixtures.Directories[1].ID, directories[1].ID)
			is.Equal(fixtures.Directories[2].ID, directories[2].ID)
			is.Equal(fixtures.Directories[11].ID, directories[3].ID)
			is.Equal(fixtures.Directories[13].ID, directories[4].ID)
			is.Equal(fixtures.Directories[15].ID, directories[5].ID)
			is.Equal(fixtures.Directories[17].ID, directories[6].ID)

			is.Empty(directories[0].Files)
			is.NotEmpty(directories[0].Directories)
			is.Len(directories[0].Directories, 6)
			is.Equal(fixtures.Directories[6].ID, directories[0].Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, directories[0].Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, directories[0].Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, directories[0].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, directories[0].Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, directories[0].Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, directories[0].Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, directories[0].Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, directories[0].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, directories[0].Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, directories[0].Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, directories[0].Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, directories[0].Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, directories[0].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, directories[0].Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, directories[0].Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, directories[0].Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, directories[0].Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, directories[0].Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, directories[0].Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, directories[0].Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, directories[0].Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, directories[0].Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, directories[0].Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, directories[0].Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, directories[0].Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, directories[0].Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, directories[0].Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, directories[0].Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, directories[0].Directories[5].ParentID)

			is.Empty(directories[1].Files)
			is.Empty(directories[1].Directories)

			is.Empty(directories[2].Files)
			is.NotEmpty(directories[2].Directories)
			is.Len(directories[2].Directories, 3)
			is.Equal(fixtures.Directories[12].ID, directories[2].Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, directories[2].Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, directories[2].Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, directories[2].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, directories[2].Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, directories[2].Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, directories[2].Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, directories[2].Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, directories[2].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, directories[2].Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, directories[2].Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, directories[2].Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, directories[2].Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, directories[2].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, directories[2].Directories[2].ParentID)

			is.Empty(directories[3].Directories)
			is.NotEmpty(directories[3].Files)
			is.Len(directories[3].Files, 4)
			is.Equal(fixtures.Files[3].ID, directories[3].Files[0].ID)
			is.Equal(fixtures.Files[3].Path, directories[3].Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, directories[3].Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, directories[3].Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, directories[3].Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, directories[3].Files[1].ID)
			is.Equal(fixtures.Files[4].Path, directories[3].Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, directories[3].Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, directories[3].Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, directories[3].Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, directories[3].Files[2].ID)
			is.Equal(fixtures.Files[5].Path, directories[3].Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, directories[3].Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, directories[3].Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, directories[3].Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, directories[3].Files[3].ID)
			is.Equal(fixtures.Files[6].Path, directories[3].Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, directories[3].Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, directories[3].Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, directories[3].Files[3].UserID)

			is.Empty(directories[4].Directories)
			is.NotEmpty(directories[4].Files)
			is.Len(directories[4].Files, 2)
			is.Equal(fixtures.Files[0].ID, directories[4].Files[0].ID)
			is.Equal(fixtures.Files[0].Path, directories[4].Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, directories[4].Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, directories[4].Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, directories[4].Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, directories[4].Files[1].ID)
			is.Equal(fixtures.Files[1].Path, directories[4].Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, directories[4].Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, directories[4].Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, directories[4].Files[1].UserID)

			is.Empty(directories[5].Directories)
			is.NotEmpty(directories[5].Files)
			is.Len(directories[5].Files, 4)
			is.Equal(fixtures.Files[8].ID, directories[5].Files[0].ID)
			is.Equal(fixtures.Files[8].Path, directories[5].Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, directories[5].Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, directories[5].Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, directories[5].Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, directories[5].Files[1].ID)
			is.Equal(fixtures.Files[19].Path, directories[5].Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, directories[5].Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, directories[5].Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, directories[5].Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, directories[5].Files[2].ID)
			is.Equal(fixtures.Files[20].Path, directories[5].Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, directories[5].Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, directories[5].Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, directories[5].Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, directories[5].Files[3].ID)
			is.Equal(fixtures.Files[21].Path, directories[5].Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, directories[5].Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, directories[5].Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, directories[5].Files[3].UserID)

			is.Empty(directories[6].Directories)
			is.NotEmpty(directories[6].Files)
			is.Len(directories[6].Files, 4)
			is.Equal(fixtures.Files[11].ID, directories[6].Files[0].ID)
			is.Equal(fixtures.Files[11].Path, directories[6].Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, directories[6].Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, directories[6].Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, directories[6].Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, directories[6].Files[1].ID)
			is.Equal(fixtures.Files[12].Path, directories[6].Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, directories[6].Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, directories[6].Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, directories[6].Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, directories[6].Files[2].ID)
			is.Equal(fixtures.Files[13].Path, directories[6].Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, directories[6].Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, directories[6].Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, directories[6].Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, directories[6].Files[3].ID)
			is.Equal(fixtures.Files[14].Path, directories[6].Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, directories[6].Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, directories[6].Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, directories[6].Files[3].UserID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directories := &[]ExoDirectory{
				*fixtures.Directories[0],
				*fixtures.Directories[1],
				*fixtures.Directories[2],
				*fixtures.Directories[11],
				*fixtures.Directories[13],
				*fixtures.Directories[15],
				*fixtures.Directories[17],
			}

			err := sqlxx.Preload(ctx, driver, &directories, "Files", "Directories")
			is.NoError(err)
			is.Len((*directories), 7)
			is.Equal(fixtures.Directories[0].ID, (*directories)[0].ID)
			is.Equal(fixtures.Directories[1].ID, (*directories)[1].ID)
			is.Equal(fixtures.Directories[2].ID, (*directories)[2].ID)
			is.Equal(fixtures.Directories[11].ID, (*directories)[3].ID)
			is.Equal(fixtures.Directories[13].ID, (*directories)[4].ID)
			is.Equal(fixtures.Directories[15].ID, (*directories)[5].ID)
			is.Equal(fixtures.Directories[17].ID, (*directories)[6].ID)

			is.Empty((*directories)[0].Files)
			is.NotEmpty((*directories)[0].Directories)
			is.Len((*directories)[0].Directories, 6)
			is.Equal(fixtures.Directories[6].ID, (*directories)[0].Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, (*directories)[0].Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, (*directories)[0].Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, (*directories)[0].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, (*directories)[0].Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, (*directories)[0].Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, (*directories)[0].Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, (*directories)[0].Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, (*directories)[0].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, (*directories)[0].Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, (*directories)[0].Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, (*directories)[0].Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, (*directories)[0].Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, (*directories)[0].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, (*directories)[0].Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, (*directories)[0].Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, (*directories)[0].Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, (*directories)[0].Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, (*directories)[0].Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, (*directories)[0].Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, (*directories)[0].Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, (*directories)[0].Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, (*directories)[0].Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, (*directories)[0].Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, (*directories)[0].Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, (*directories)[0].Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, (*directories)[0].Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, (*directories)[0].Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, (*directories)[0].Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, (*directories)[0].Directories[5].ParentID)

			is.Empty((*directories)[1].Files)
			is.Empty((*directories)[1].Directories)

			is.Empty((*directories)[2].Files)
			is.NotEmpty((*directories)[2].Directories)
			is.Len((*directories)[2].Directories, 3)
			is.Equal(fixtures.Directories[12].ID, (*directories)[2].Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, (*directories)[2].Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, (*directories)[2].Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, (*directories)[2].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, (*directories)[2].Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, (*directories)[2].Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, (*directories)[2].Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, (*directories)[2].Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, (*directories)[2].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, (*directories)[2].Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, (*directories)[2].Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, (*directories)[2].Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, (*directories)[2].Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, (*directories)[2].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, (*directories)[2].Directories[2].ParentID)

			is.Empty((*directories)[3].Directories)
			is.NotEmpty((*directories)[3].Files)
			is.Len((*directories)[3].Files, 4)
			is.Equal(fixtures.Files[3].ID, (*directories)[3].Files[0].ID)
			is.Equal(fixtures.Files[3].Path, (*directories)[3].Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, (*directories)[3].Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, (*directories)[3].Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, (*directories)[3].Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, (*directories)[3].Files[1].ID)
			is.Equal(fixtures.Files[4].Path, (*directories)[3].Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, (*directories)[3].Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, (*directories)[3].Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, (*directories)[3].Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, (*directories)[3].Files[2].ID)
			is.Equal(fixtures.Files[5].Path, (*directories)[3].Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, (*directories)[3].Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, (*directories)[3].Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, (*directories)[3].Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, (*directories)[3].Files[3].ID)
			is.Equal(fixtures.Files[6].Path, (*directories)[3].Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, (*directories)[3].Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, (*directories)[3].Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, (*directories)[3].Files[3].UserID)

			is.Empty((*directories)[4].Directories)
			is.NotEmpty((*directories)[4].Files)
			is.Len((*directories)[4].Files, 2)
			is.Equal(fixtures.Files[0].ID, (*directories)[4].Files[0].ID)
			is.Equal(fixtures.Files[0].Path, (*directories)[4].Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, (*directories)[4].Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, (*directories)[4].Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, (*directories)[4].Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, (*directories)[4].Files[1].ID)
			is.Equal(fixtures.Files[1].Path, (*directories)[4].Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, (*directories)[4].Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, (*directories)[4].Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, (*directories)[4].Files[1].UserID)

			is.Empty((*directories)[5].Directories)
			is.NotEmpty((*directories)[5].Files)
			is.Len((*directories)[5].Files, 4)
			is.Equal(fixtures.Files[8].ID, (*directories)[5].Files[0].ID)
			is.Equal(fixtures.Files[8].Path, (*directories)[5].Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, (*directories)[5].Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, (*directories)[5].Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, (*directories)[5].Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, (*directories)[5].Files[1].ID)
			is.Equal(fixtures.Files[19].Path, (*directories)[5].Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, (*directories)[5].Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, (*directories)[5].Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, (*directories)[5].Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, (*directories)[5].Files[2].ID)
			is.Equal(fixtures.Files[20].Path, (*directories)[5].Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, (*directories)[5].Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, (*directories)[5].Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, (*directories)[5].Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, (*directories)[5].Files[3].ID)
			is.Equal(fixtures.Files[21].Path, (*directories)[5].Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, (*directories)[5].Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, (*directories)[5].Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, (*directories)[5].Files[3].UserID)

			is.Empty((*directories)[6].Directories)
			is.NotEmpty((*directories)[6].Files)
			is.Len((*directories)[6].Files, 4)
			is.Equal(fixtures.Files[11].ID, (*directories)[6].Files[0].ID)
			is.Equal(fixtures.Files[11].Path, (*directories)[6].Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, (*directories)[6].Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, (*directories)[6].Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, (*directories)[6].Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, (*directories)[6].Files[1].ID)
			is.Equal(fixtures.Files[12].Path, (*directories)[6].Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, (*directories)[6].Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, (*directories)[6].Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, (*directories)[6].Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, (*directories)[6].Files[2].ID)
			is.Equal(fixtures.Files[13].Path, (*directories)[6].Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, (*directories)[6].Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, (*directories)[6].Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, (*directories)[6].Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, (*directories)[6].Files[3].ID)
			is.Equal(fixtures.Files[14].Path, (*directories)[6].Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, (*directories)[6].Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, (*directories)[6].Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, (*directories)[6].Files[3].UserID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			directories := &[]*ExoDirectory{
				fixtures.Directories[0],
				fixtures.Directories[1],
				fixtures.Directories[2],
				fixtures.Directories[11],
				fixtures.Directories[13],
				fixtures.Directories[15],
				fixtures.Directories[17],
			}

			err := sqlxx.Preload(ctx, driver, &directories, "Files", "Directories")
			is.NoError(err)
			is.Len((*directories), 7)
			is.Equal(fixtures.Directories[0].ID, (*directories)[0].ID)
			is.Equal(fixtures.Directories[1].ID, (*directories)[1].ID)
			is.Equal(fixtures.Directories[2].ID, (*directories)[2].ID)
			is.Equal(fixtures.Directories[11].ID, (*directories)[3].ID)
			is.Equal(fixtures.Directories[13].ID, (*directories)[4].ID)
			is.Equal(fixtures.Directories[15].ID, (*directories)[5].ID)
			is.Equal(fixtures.Directories[17].ID, (*directories)[6].ID)

			is.Empty((*directories)[0].Files)
			is.NotEmpty((*directories)[0].Directories)
			is.Len((*directories)[0].Directories, 6)
			is.Equal(fixtures.Directories[6].ID, (*directories)[0].Directories[0].ID)
			is.Equal(fixtures.Directories[6].Path, (*directories)[0].Directories[0].Path)
			is.Equal(fixtures.Directories[6].BucketID, (*directories)[0].Directories[0].BucketID)
			is.Equal(fixtures.Directories[6].OrganizationID, (*directories)[0].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[6].ParentID, (*directories)[0].Directories[0].ParentID)
			is.Equal(fixtures.Directories[7].ID, (*directories)[0].Directories[1].ID)
			is.Equal(fixtures.Directories[7].Path, (*directories)[0].Directories[1].Path)
			is.Equal(fixtures.Directories[7].BucketID, (*directories)[0].Directories[1].BucketID)
			is.Equal(fixtures.Directories[7].OrganizationID, (*directories)[0].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[7].ParentID, (*directories)[0].Directories[1].ParentID)
			is.Equal(fixtures.Directories[8].ID, (*directories)[0].Directories[2].ID)
			is.Equal(fixtures.Directories[8].Path, (*directories)[0].Directories[2].Path)
			is.Equal(fixtures.Directories[8].BucketID, (*directories)[0].Directories[2].BucketID)
			is.Equal(fixtures.Directories[8].OrganizationID, (*directories)[0].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[8].ParentID, (*directories)[0].Directories[2].ParentID)
			is.Equal(fixtures.Directories[9].ID, (*directories)[0].Directories[3].ID)
			is.Equal(fixtures.Directories[9].Path, (*directories)[0].Directories[3].Path)
			is.Equal(fixtures.Directories[9].BucketID, (*directories)[0].Directories[3].BucketID)
			is.Equal(fixtures.Directories[9].OrganizationID, (*directories)[0].Directories[3].OrganizationID)
			is.Equal(fixtures.Directories[9].ParentID, (*directories)[0].Directories[3].ParentID)
			is.Equal(fixtures.Directories[10].ID, (*directories)[0].Directories[4].ID)
			is.Equal(fixtures.Directories[10].Path, (*directories)[0].Directories[4].Path)
			is.Equal(fixtures.Directories[10].BucketID, (*directories)[0].Directories[4].BucketID)
			is.Equal(fixtures.Directories[10].OrganizationID, (*directories)[0].Directories[4].OrganizationID)
			is.Equal(fixtures.Directories[10].ParentID, (*directories)[0].Directories[4].ParentID)
			is.Equal(fixtures.Directories[11].ID, (*directories)[0].Directories[5].ID)
			is.Equal(fixtures.Directories[11].Path, (*directories)[0].Directories[5].Path)
			is.Equal(fixtures.Directories[11].BucketID, (*directories)[0].Directories[5].BucketID)
			is.Equal(fixtures.Directories[11].OrganizationID, (*directories)[0].Directories[5].OrganizationID)
			is.Equal(fixtures.Directories[11].ParentID, (*directories)[0].Directories[5].ParentID)

			is.Empty((*directories)[1].Files)
			is.Empty((*directories)[1].Directories)

			is.Empty((*directories)[2].Files)
			is.NotEmpty((*directories)[2].Directories)
			is.Len((*directories)[2].Directories, 3)
			is.Equal(fixtures.Directories[12].ID, (*directories)[2].Directories[0].ID)
			is.Equal(fixtures.Directories[12].Path, (*directories)[2].Directories[0].Path)
			is.Equal(fixtures.Directories[12].BucketID, (*directories)[2].Directories[0].BucketID)
			is.Equal(fixtures.Directories[12].OrganizationID, (*directories)[2].Directories[0].OrganizationID)
			is.Equal(fixtures.Directories[12].ParentID, (*directories)[2].Directories[0].ParentID)
			is.Equal(fixtures.Directories[13].ID, (*directories)[2].Directories[1].ID)
			is.Equal(fixtures.Directories[13].Path, (*directories)[2].Directories[1].Path)
			is.Equal(fixtures.Directories[13].BucketID, (*directories)[2].Directories[1].BucketID)
			is.Equal(fixtures.Directories[13].OrganizationID, (*directories)[2].Directories[1].OrganizationID)
			is.Equal(fixtures.Directories[13].ParentID, (*directories)[2].Directories[1].ParentID)
			is.Equal(fixtures.Directories[14].ID, (*directories)[2].Directories[2].ID)
			is.Equal(fixtures.Directories[14].Path, (*directories)[2].Directories[2].Path)
			is.Equal(fixtures.Directories[14].BucketID, (*directories)[2].Directories[2].BucketID)
			is.Equal(fixtures.Directories[14].OrganizationID, (*directories)[2].Directories[2].OrganizationID)
			is.Equal(fixtures.Directories[14].ParentID, (*directories)[2].Directories[2].ParentID)

			is.Empty((*directories)[3].Directories)
			is.NotEmpty((*directories)[3].Files)
			is.Len((*directories)[3].Files, 4)
			is.Equal(fixtures.Files[3].ID, (*directories)[3].Files[0].ID)
			is.Equal(fixtures.Files[3].Path, (*directories)[3].Files[0].Path)
			is.Equal(fixtures.Files[3].BucketID, (*directories)[3].Files[0].BucketID)
			is.Equal(fixtures.Files[3].OrganizationID, (*directories)[3].Files[0].OrganizationID)
			is.Equal(fixtures.Files[3].UserID, (*directories)[3].Files[0].UserID)
			is.Equal(fixtures.Files[4].ID, (*directories)[3].Files[1].ID)
			is.Equal(fixtures.Files[4].Path, (*directories)[3].Files[1].Path)
			is.Equal(fixtures.Files[4].BucketID, (*directories)[3].Files[1].BucketID)
			is.Equal(fixtures.Files[4].OrganizationID, (*directories)[3].Files[1].OrganizationID)
			is.Equal(fixtures.Files[4].UserID, (*directories)[3].Files[1].UserID)
			is.Equal(fixtures.Files[5].ID, (*directories)[3].Files[2].ID)
			is.Equal(fixtures.Files[5].Path, (*directories)[3].Files[2].Path)
			is.Equal(fixtures.Files[5].BucketID, (*directories)[3].Files[2].BucketID)
			is.Equal(fixtures.Files[5].OrganizationID, (*directories)[3].Files[2].OrganizationID)
			is.Equal(fixtures.Files[5].UserID, (*directories)[3].Files[2].UserID)
			is.Equal(fixtures.Files[6].ID, (*directories)[3].Files[3].ID)
			is.Equal(fixtures.Files[6].Path, (*directories)[3].Files[3].Path)
			is.Equal(fixtures.Files[6].BucketID, (*directories)[3].Files[3].BucketID)
			is.Equal(fixtures.Files[6].OrganizationID, (*directories)[3].Files[3].OrganizationID)
			is.Equal(fixtures.Files[6].UserID, (*directories)[3].Files[3].UserID)

			is.Empty((*directories)[4].Directories)
			is.NotEmpty((*directories)[4].Files)
			is.Len((*directories)[4].Files, 2)
			is.Equal(fixtures.Files[0].ID, (*directories)[4].Files[0].ID)
			is.Equal(fixtures.Files[0].Path, (*directories)[4].Files[0].Path)
			is.Equal(fixtures.Files[0].BucketID, (*directories)[4].Files[0].BucketID)
			is.Equal(fixtures.Files[0].OrganizationID, (*directories)[4].Files[0].OrganizationID)
			is.Equal(fixtures.Files[0].UserID, (*directories)[4].Files[0].UserID)
			is.Equal(fixtures.Files[1].ID, (*directories)[4].Files[1].ID)
			is.Equal(fixtures.Files[1].Path, (*directories)[4].Files[1].Path)
			is.Equal(fixtures.Files[1].BucketID, (*directories)[4].Files[1].BucketID)
			is.Equal(fixtures.Files[1].OrganizationID, (*directories)[4].Files[1].OrganizationID)
			is.Equal(fixtures.Files[1].UserID, (*directories)[4].Files[1].UserID)

			is.Empty((*directories)[5].Directories)
			is.NotEmpty((*directories)[5].Files)
			is.Len((*directories)[5].Files, 4)
			is.Equal(fixtures.Files[8].ID, (*directories)[5].Files[0].ID)
			is.Equal(fixtures.Files[8].Path, (*directories)[5].Files[0].Path)
			is.Equal(fixtures.Files[8].BucketID, (*directories)[5].Files[0].BucketID)
			is.Equal(fixtures.Files[8].OrganizationID, (*directories)[5].Files[0].OrganizationID)
			is.Equal(fixtures.Files[8].UserID, (*directories)[5].Files[0].UserID)
			is.Equal(fixtures.Files[19].ID, (*directories)[5].Files[1].ID)
			is.Equal(fixtures.Files[19].Path, (*directories)[5].Files[1].Path)
			is.Equal(fixtures.Files[19].BucketID, (*directories)[5].Files[1].BucketID)
			is.Equal(fixtures.Files[19].OrganizationID, (*directories)[5].Files[1].OrganizationID)
			is.Equal(fixtures.Files[19].UserID, (*directories)[5].Files[1].UserID)
			is.Equal(fixtures.Files[20].ID, (*directories)[5].Files[2].ID)
			is.Equal(fixtures.Files[20].Path, (*directories)[5].Files[2].Path)
			is.Equal(fixtures.Files[20].BucketID, (*directories)[5].Files[2].BucketID)
			is.Equal(fixtures.Files[20].OrganizationID, (*directories)[5].Files[2].OrganizationID)
			is.Equal(fixtures.Files[20].UserID, (*directories)[5].Files[2].UserID)
			is.Equal(fixtures.Files[21].ID, (*directories)[5].Files[3].ID)
			is.Equal(fixtures.Files[21].Path, (*directories)[5].Files[3].Path)
			is.Equal(fixtures.Files[21].BucketID, (*directories)[5].Files[3].BucketID)
			is.Equal(fixtures.Files[21].OrganizationID, (*directories)[5].Files[3].OrganizationID)
			is.Equal(fixtures.Files[21].UserID, (*directories)[5].Files[3].UserID)

			is.Empty((*directories)[6].Directories)
			is.NotEmpty((*directories)[6].Files)
			is.Len((*directories)[6].Files, 4)
			is.Equal(fixtures.Files[11].ID, (*directories)[6].Files[0].ID)
			is.Equal(fixtures.Files[11].Path, (*directories)[6].Files[0].Path)
			is.Equal(fixtures.Files[11].BucketID, (*directories)[6].Files[0].BucketID)
			is.Equal(fixtures.Files[11].OrganizationID, (*directories)[6].Files[0].OrganizationID)
			is.Equal(fixtures.Files[11].UserID, (*directories)[6].Files[0].UserID)
			is.Equal(fixtures.Files[12].ID, (*directories)[6].Files[1].ID)
			is.Equal(fixtures.Files[12].Path, (*directories)[6].Files[1].Path)
			is.Equal(fixtures.Files[12].BucketID, (*directories)[6].Files[1].BucketID)
			is.Equal(fixtures.Files[12].OrganizationID, (*directories)[6].Files[1].OrganizationID)
			is.Equal(fixtures.Files[12].UserID, (*directories)[6].Files[1].UserID)
			is.Equal(fixtures.Files[13].ID, (*directories)[6].Files[2].ID)
			is.Equal(fixtures.Files[13].Path, (*directories)[6].Files[2].Path)
			is.Equal(fixtures.Files[13].BucketID, (*directories)[6].Files[2].BucketID)
			is.Equal(fixtures.Files[13].OrganizationID, (*directories)[6].Files[2].OrganizationID)
			is.Equal(fixtures.Files[13].UserID, (*directories)[6].Files[2].UserID)
			is.Equal(fixtures.Files[14].ID, (*directories)[6].Files[3].ID)
			is.Equal(fixtures.Files[14].Path, (*directories)[6].Files[3].Path)
			is.Equal(fixtures.Files[14].BucketID, (*directories)[6].Files[3].BucketID)
			is.Equal(fixtures.Files[14].OrganizationID, (*directories)[6].Files[3].OrganizationID)
			is.Equal(fixtures.Files[14].UserID, (*directories)[6].Files[3].UserID)

		}
	})
}

func TestPreload_ExoChunk_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Nil(fixtures.Chunks[23].Mode)
			is.Nil(fixtures.Chunks[24].Mode)
			is.Nil(fixtures.Chunks[25].Mode)
			is.Nil(fixtures.Chunks[26].Mode)
			is.Nil(fixtures.Chunks[27].Mode)
			is.Nil(fixtures.Chunks[28].Mode)
			is.Nil(fixtures.Chunks[23].Signature)
			is.Nil(fixtures.Chunks[24].Signature)
			is.Nil(fixtures.Chunks[25].Signature)
			is.Nil(fixtures.Chunks[26].Signature)
			is.Nil(fixtures.Chunks[27].Signature)
			is.Nil(fixtures.Chunks[28].Signature)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunk1 := fixtures.Chunks[24]

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

			chunk2 := fixtures.Chunks[25]

			err = sqlxx.Preload(ctx, driver, chunk2, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk2.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk2.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk2.Mode.Mode)
			is.NotNil(chunk2.Signature)
			is.Equal(fixtures.Signatures[1].ID, chunk2.Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunk2.Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunk2.Signature.Bytes)

			chunk3 := fixtures.Chunks[26]

			err = sqlxx.Preload(ctx, driver, chunk3, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk3.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk3.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk3.Mode.Mode)
			is.NotNil(chunk3.Signature)
			is.Equal(fixtures.Signatures[2].ID, chunk3.Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunk3.Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunk3.Signature.Bytes)

			chunk4 := fixtures.Chunks[27]

			err = sqlxx.Preload(ctx, driver, chunk4, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk4.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk4.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk4.Mode.Mode)
			is.NotNil(chunk4.Signature)
			is.Equal(fixtures.Signatures[3].ID, chunk4.Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunk4.Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunk4.Signature.Bytes)

			chunk5 := fixtures.Chunks[28]

			err = sqlxx.Preload(ctx, driver, chunk5, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk5.Mode)
			is.Equal(fixtures.Modes[2].ID, chunk5.Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunk5.Mode.Mode)
			is.Nil(chunk5.Signature)

			chunk6 := fixtures.Chunks[23]

			err = sqlxx.Preload(ctx, driver, chunk6, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk6.Mode)
			is.Equal(fixtures.Modes[3].ID, chunk6.Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunk6.Mode.Mode)
			is.Nil(chunk6.Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunk1 := fixtures.Chunks[24]

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

			chunk2 := fixtures.Chunks[25]

			err = sqlxx.Preload(ctx, driver, &chunk2, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk2.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk2.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk2.Mode.Mode)
			is.NotNil(chunk2.Signature)
			is.Equal(fixtures.Signatures[1].ID, chunk2.Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunk2.Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunk2.Signature.Bytes)

			chunk3 := fixtures.Chunks[26]

			err = sqlxx.Preload(ctx, driver, &chunk3, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk3.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk3.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk3.Mode.Mode)
			is.NotNil(chunk3.Signature)
			is.Equal(fixtures.Signatures[2].ID, chunk3.Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunk3.Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunk3.Signature.Bytes)

			chunk4 := fixtures.Chunks[27]

			err = sqlxx.Preload(ctx, driver, &chunk4, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk4.Mode)
			is.Equal(fixtures.Modes[0].ID, chunk4.Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunk4.Mode.Mode)
			is.NotNil(chunk4.Signature)
			is.Equal(fixtures.Signatures[3].ID, chunk4.Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunk4.Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunk4.Signature.Bytes)

			chunk5 := fixtures.Chunks[28]

			err = sqlxx.Preload(ctx, driver, &chunk5, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk5.Mode)
			is.Equal(fixtures.Modes[2].ID, chunk5.Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunk5.Mode.Mode)
			is.Nil(chunk5.Signature)

			chunk6 := fixtures.Chunks[23]

			err = sqlxx.Preload(ctx, driver, &chunk6, "Mode", "Signature")
			is.NoError(err)
			is.NotNil(chunk6.Mode)
			is.Equal(fixtures.Modes[3].ID, chunk6.Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunk6.Mode.Mode)
			is.Nil(chunk6.Signature)

		}
	})
}

func TestPreload_ExoChunk_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckExoCloudFixtures := func(fixtures *ExoCloudFixtures) {
			is.Nil(fixtures.Chunks[23].Mode)
			is.Nil(fixtures.Chunks[24].Mode)
			is.Nil(fixtures.Chunks[25].Mode)
			is.Nil(fixtures.Chunks[26].Mode)
			is.Nil(fixtures.Chunks[27].Mode)
			is.Nil(fixtures.Chunks[28].Mode)
			is.Nil(fixtures.Chunks[23].Signature)
			is.Nil(fixtures.Chunks[24].Signature)
			is.Nil(fixtures.Chunks[25].Signature)
			is.Nil(fixtures.Chunks[26].Signature)
			is.Nil(fixtures.Chunks[27].Signature)
			is.Nil(fixtures.Chunks[28].Signature)
		}

		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := []ExoChunk{
				*fixtures.Chunks[23],
				*fixtures.Chunks[24],
				*fixtures.Chunks[25],
				*fixtures.Chunks[26],
				*fixtures.Chunks[27],
				*fixtures.Chunks[28],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 6)
			is.Equal(fixtures.Chunks[23].Hash, chunks[0].Hash)
			is.Equal(fixtures.Chunks[24].Hash, chunks[1].Hash)
			is.Equal(fixtures.Chunks[25].Hash, chunks[2].Hash)
			is.Equal(fixtures.Chunks[26].Hash, chunks[3].Hash)
			is.Equal(fixtures.Chunks[27].Hash, chunks[4].Hash)
			is.Equal(fixtures.Chunks[28].Hash, chunks[5].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(fixtures.Modes[3].ID, chunks[0].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunks[0].Mode.Mode)
			is.Nil(chunks[0].Signature)

			is.NotNil(chunks[1].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(fixtures.Signatures[0].ID, chunks[1].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(fixtures.Signatures[1].ID, chunks[2].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[3].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[3].Mode.Mode)
			is.NotNil(chunks[3].Signature)
			is.Equal(fixtures.Signatures[2].ID, chunks[3].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunks[3].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunks[3].Signature.Bytes)

			is.NotNil(chunks[4].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[4].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[4].Mode.Mode)
			is.NotNil(chunks[4].Signature)
			is.Equal(fixtures.Signatures[3].ID, chunks[4].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunks[4].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunks[4].Signature.Bytes)

			is.NotNil(chunks[5].Mode)
			is.Equal(fixtures.Modes[2].ID, chunks[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunks[5].Mode.Mode)
			is.Nil(chunks[5].Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := []*ExoChunk{
				fixtures.Chunks[23],
				fixtures.Chunks[24],
				fixtures.Chunks[25],
				fixtures.Chunks[26],
				fixtures.Chunks[27],
				fixtures.Chunks[28],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len(chunks, 6)
			is.Equal(fixtures.Chunks[23].Hash, chunks[0].Hash)
			is.Equal(fixtures.Chunks[24].Hash, chunks[1].Hash)
			is.Equal(fixtures.Chunks[25].Hash, chunks[2].Hash)
			is.Equal(fixtures.Chunks[26].Hash, chunks[3].Hash)
			is.Equal(fixtures.Chunks[27].Hash, chunks[4].Hash)
			is.Equal(fixtures.Chunks[28].Hash, chunks[5].Hash)

			is.NotNil(chunks[0].Mode)
			is.Equal(fixtures.Modes[3].ID, chunks[0].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, chunks[0].Mode.Mode)
			is.Nil(chunks[0].Signature)

			is.NotNil(chunks[1].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[1].Mode.Mode)
			is.NotNil(chunks[1].Signature)
			is.Equal(fixtures.Signatures[0].ID, chunks[1].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, chunks[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, chunks[1].Signature.Bytes)

			is.NotNil(chunks[2].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[2].Mode.Mode)
			is.NotNil(chunks[2].Signature)
			is.Equal(fixtures.Signatures[1].ID, chunks[2].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, chunks[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, chunks[2].Signature.Bytes)

			is.NotNil(chunks[3].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[3].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[3].Mode.Mode)
			is.NotNil(chunks[3].Signature)
			is.Equal(fixtures.Signatures[2].ID, chunks[3].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, chunks[3].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, chunks[3].Signature.Bytes)

			is.NotNil(chunks[4].Mode)
			is.Equal(fixtures.Modes[0].ID, chunks[4].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, chunks[4].Mode.Mode)
			is.NotNil(chunks[4].Signature)
			is.Equal(fixtures.Signatures[3].ID, chunks[4].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, chunks[4].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, chunks[4].Signature.Bytes)

			is.NotNil(chunks[5].Mode)
			is.Equal(fixtures.Modes[2].ID, chunks[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, chunks[5].Mode.Mode)
			is.Nil(chunks[5].Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := &[]ExoChunk{
				*fixtures.Chunks[23],
				*fixtures.Chunks[24],
				*fixtures.Chunks[25],
				*fixtures.Chunks[26],
				*fixtures.Chunks[27],
				*fixtures.Chunks[28],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len((*chunks), 6)
			is.Equal(fixtures.Chunks[23].Hash, (*chunks)[0].Hash)
			is.Equal(fixtures.Chunks[24].Hash, (*chunks)[1].Hash)
			is.Equal(fixtures.Chunks[25].Hash, (*chunks)[2].Hash)
			is.Equal(fixtures.Chunks[26].Hash, (*chunks)[3].Hash)
			is.Equal(fixtures.Chunks[27].Hash, (*chunks)[4].Hash)
			is.Equal(fixtures.Chunks[28].Hash, (*chunks)[5].Hash)

			is.NotNil((*chunks)[0].Mode)
			is.Equal(fixtures.Modes[3].ID, (*chunks)[0].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, (*chunks)[0].Mode.Mode)
			is.Nil((*chunks)[0].Signature)

			is.NotNil((*chunks)[1].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[1].Mode.Mode)
			is.NotNil((*chunks)[1].Signature)
			is.Equal(fixtures.Signatures[0].ID, (*chunks)[1].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, (*chunks)[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, (*chunks)[1].Signature.Bytes)

			is.NotNil((*chunks)[2].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[2].Mode.Mode)
			is.NotNil((*chunks)[2].Signature)
			is.Equal(fixtures.Signatures[1].ID, (*chunks)[2].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, (*chunks)[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, (*chunks)[2].Signature.Bytes)

			is.NotNil((*chunks)[3].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[3].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[3].Mode.Mode)
			is.NotNil((*chunks)[3].Signature)
			is.Equal(fixtures.Signatures[2].ID, (*chunks)[3].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, (*chunks)[3].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, (*chunks)[3].Signature.Bytes)

			is.NotNil((*chunks)[4].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[4].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[4].Mode.Mode)
			is.NotNil((*chunks)[4].Signature)
			is.Equal(fixtures.Signatures[3].ID, (*chunks)[4].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, (*chunks)[4].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, (*chunks)[4].Signature.Bytes)

			is.NotNil((*chunks)[5].Mode)
			is.Equal(fixtures.Modes[2].ID, (*chunks)[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, (*chunks)[5].Mode.Mode)
			is.Nil((*chunks)[5].Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			chunks := &[]*ExoChunk{
				fixtures.Chunks[23],
				fixtures.Chunks[24],
				fixtures.Chunks[25],
				fixtures.Chunks[26],
				fixtures.Chunks[27],
				fixtures.Chunks[28],
			}

			err := sqlxx.Preload(ctx, driver, &chunks, "Mode", "Signature")
			is.NoError(err)
			is.Len((*chunks), 6)
			is.Equal(fixtures.Chunks[23].Hash, (*chunks)[0].Hash)
			is.Equal(fixtures.Chunks[24].Hash, (*chunks)[1].Hash)
			is.Equal(fixtures.Chunks[25].Hash, (*chunks)[2].Hash)
			is.Equal(fixtures.Chunks[26].Hash, (*chunks)[3].Hash)
			is.Equal(fixtures.Chunks[27].Hash, (*chunks)[4].Hash)
			is.Equal(fixtures.Chunks[28].Hash, (*chunks)[5].Hash)

			is.NotNil((*chunks)[0].Mode)
			is.Equal(fixtures.Modes[3].ID, (*chunks)[0].Mode.ID)
			is.Equal(fixtures.Modes[3].Mode, (*chunks)[0].Mode.Mode)
			is.Nil((*chunks)[0].Signature)

			is.NotNil((*chunks)[1].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[1].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[1].Mode.Mode)
			is.NotNil((*chunks)[1].Signature)
			is.Equal(fixtures.Signatures[0].ID, (*chunks)[1].Signature.ID)
			is.Equal(fixtures.Signatures[0].ChunkID, (*chunks)[1].Signature.ChunkID)
			is.Equal(fixtures.Signatures[0].Bytes, (*chunks)[1].Signature.Bytes)

			is.NotNil((*chunks)[2].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[2].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[2].Mode.Mode)
			is.NotNil((*chunks)[2].Signature)
			is.Equal(fixtures.Signatures[1].ID, (*chunks)[2].Signature.ID)
			is.Equal(fixtures.Signatures[1].ChunkID, (*chunks)[2].Signature.ChunkID)
			is.Equal(fixtures.Signatures[1].Bytes, (*chunks)[2].Signature.Bytes)

			is.NotNil((*chunks)[3].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[3].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[3].Mode.Mode)
			is.NotNil((*chunks)[3].Signature)
			is.Equal(fixtures.Signatures[2].ID, (*chunks)[3].Signature.ID)
			is.Equal(fixtures.Signatures[2].ChunkID, (*chunks)[3].Signature.ChunkID)
			is.Equal(fixtures.Signatures[2].Bytes, (*chunks)[3].Signature.Bytes)

			is.NotNil((*chunks)[4].Mode)
			is.Equal(fixtures.Modes[0].ID, (*chunks)[4].Mode.ID)
			is.Equal(fixtures.Modes[0].Mode, (*chunks)[4].Mode.Mode)
			is.NotNil((*chunks)[4].Signature)
			is.Equal(fixtures.Signatures[3].ID, (*chunks)[4].Signature.ID)
			is.Equal(fixtures.Signatures[3].ChunkID, (*chunks)[4].Signature.ChunkID)
			is.Equal(fixtures.Signatures[3].Bytes, (*chunks)[4].Signature.Bytes)

			is.NotNil((*chunks)[5].Mode)
			is.Equal(fixtures.Modes[2].ID, (*chunks)[5].Mode.ID)
			is.Equal(fixtures.Modes[2].Mode, (*chunks)[5].Mode.Mode)
			is.Nil((*chunks)[5].Signature)

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
			is.Empty(fixtures.Owls[0].Packages)
			is.Empty(fixtures.Owls[1].Packages)
			is.Empty(fixtures.Owls[2].Packages)
			is.Empty(fixtures.Owls[3].Packages)
			is.Empty(fixtures.Owls[4].Packages)
			is.Empty(fixtures.Owls[5].Packages)
			is.Empty(fixtures.Owls[0].Bag)
			is.Empty(fixtures.Owls[1].Bag)
			is.Empty(fixtures.Owls[2].Bag)
			is.Empty(fixtures.Owls[3].Bag)
			is.Empty(fixtures.Owls[4].Bag)
			is.Empty(fixtures.Owls[5].Bag)
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
			is.Empty(owl1.Bag)

			err = sqlxx.Preload(ctx, driver, owl1, "Packages")
			is.NoError(err)
			is.NotEmpty(owl1.Packages)
			is.Len(owl1.Packages, 2)
			is.Contains(owl1.Packages, *fixtures.Packages[0])
			is.Contains(owl1.Packages, *fixtures.Packages[1])
			is.Empty(owl1.Bag)

			err = sqlxx.Preload(ctx, driver, owl1, "Bag")
			is.NoError(err)
			is.NotEmpty(owl1.Bag)
			is.Equal(fixtures.Bags[0].ID, owl1.Bag.ID)
			is.Equal(fixtures.Bags[0].Color, owl1.Bag.Color)

			owl2 := fixtures.Owls[1]

			err = sqlxx.Preload(ctx, driver, owl2, "Group", "Bag", "Packages")
			is.NoError(err)
			is.Nil(owl2.Group)
			is.NotEmpty(owl2.Packages)
			is.Len(owl2.Packages, 1)
			is.Contains(owl2.Packages, *fixtures.Packages[3])
			is.NotEmpty(owl2.Bag)
			is.Equal(fixtures.Bags[1].ID, owl2.Bag.ID)
			is.Equal(fixtures.Bags[1].Color, owl2.Bag.Color)

			owl3 := fixtures.Owls[2]

			err = sqlxx.Preload(ctx, driver, owl3, "Group", "Bag", "Packages")
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
			is.Empty(owl3.Bag)

			owl5 := fixtures.Owls[4]

			err = sqlxx.Preload(ctx, driver, owl5, "Group", "Bag", "Packages")
			is.NoError(err)
			is.NotNil(owl5.Group)
			is.Equal(fixtures.Groups[2].ID, owl5.Group.ID)
			is.Equal(fixtures.Groups[2].Name, owl5.Group.Name)
			is.Empty(owl5.Packages)
			is.NotEmpty(owl5.Bag)
			is.Equal(fixtures.Bags[3].ID, owl5.Bag.ID)
			is.Equal(fixtures.Bags[3].Color, owl5.Bag.Color)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckOwlFixtures(fixtures)

			owl1 := fixtures.Owls[0]

			err := sqlxx.Preload(ctx, driver, &owl1, "Group", "Bag", "Packages")
			is.NoError(err)
			is.NotNil(owl1.Group)
			is.Equal(fixtures.Groups[0].ID, owl1.Group.ID)
			is.Equal(fixtures.Groups[0].Name, owl1.Group.Name)
			is.NotEmpty(owl1.Packages)
			is.Len(owl1.Packages, 2)
			is.Contains(owl1.Packages, *fixtures.Packages[0])
			is.Contains(owl1.Packages, *fixtures.Packages[1])
			is.NotEmpty(owl1.Bag)
			is.Equal(fixtures.Bags[0].ID, owl1.Bag.ID)
			is.Equal(fixtures.Bags[0].Color, owl1.Bag.Color)

			owl2 := fixtures.Owls[1]

			err = sqlxx.Preload(ctx, driver, &owl2, "Group", "Bag", "Packages")
			is.NoError(err)
			is.Nil(owl2.Group)
			is.NotEmpty(owl2.Packages)
			is.Len(owl2.Packages, 1)
			is.Contains(owl2.Packages, *fixtures.Packages[3])
			is.NotEmpty(owl2.Bag)
			is.Equal(fixtures.Bags[1].ID, owl2.Bag.ID)
			is.Equal(fixtures.Bags[1].Color, owl2.Bag.Color)

			owl3 := fixtures.Owls[2]

			err = sqlxx.Preload(ctx, driver, &owl3, "Group", "Bag", "Packages")
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
			is.Empty(owl3.Bag)

			owl5 := fixtures.Owls[4]

			err = sqlxx.Preload(ctx, driver, &owl5, "Group", "Bag", "Packages")
			is.NoError(err)
			is.NotNil(owl5.Group)
			is.Equal(fixtures.Groups[2].ID, owl5.Group.ID)
			is.Equal(fixtures.Groups[2].Name, owl5.Group.Name)
			is.Empty(owl5.Packages)
			is.NotEmpty(owl5.Bag)
			is.Equal(fixtures.Bags[3].ID, owl5.Bag.ID)
			is.Equal(fixtures.Bags[3].Color, owl5.Bag.Color)

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
			is.Empty(fixtures.Owls[0].Packages)
			is.Empty(fixtures.Owls[1].Packages)
			is.Empty(fixtures.Owls[2].Packages)
			is.Empty(fixtures.Owls[3].Packages)
			is.Empty(fixtures.Owls[4].Packages)
			is.Empty(fixtures.Owls[5].Packages)
			is.Empty(fixtures.Owls[0].Bag)
			is.Empty(fixtures.Owls[1].Bag)
			is.Empty(fixtures.Owls[2].Bag)
			is.Empty(fixtures.Owls[3].Bag)
			is.Empty(fixtures.Owls[4].Bag)
			is.Empty(fixtures.Owls[5].Bag)
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

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Bag", "Packages")
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
			is.NotEmpty(owls[0].Bag)
			is.Equal(fixtures.Bags[0].ID, owls[0].Bag.ID)
			is.Equal(fixtures.Bags[0].Color, owls[0].Bag.Color)

			is.Nil(owls[1].Group)
			is.NotEmpty(owls[1].Packages)
			is.Len(owls[1].Packages, 1)
			is.Contains(owls[1].Packages, *fixtures.Packages[3])
			is.NotEmpty(owls[1].Bag)
			is.Equal(fixtures.Bags[1].ID, owls[1].Bag.ID)
			is.Equal(fixtures.Bags[1].Color, owls[1].Bag.Color)

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
			is.Empty(owls[2].Bag)

			is.NotNil(owls[3].Group)
			is.Equal(fixtures.Groups[1].ID, owls[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, owls[3].Group.Name)
			is.NotEmpty(owls[3].Packages)
			is.Len(owls[3].Packages, 4)
			is.Contains(owls[3].Packages, *fixtures.Packages[9])
			is.Contains(owls[3].Packages, *fixtures.Packages[10])
			is.Contains(owls[3].Packages, *fixtures.Packages[11])
			is.Contains(owls[3].Packages, *fixtures.Packages[12])
			is.NotEmpty(owls[3].Bag)
			is.Equal(fixtures.Bags[2].ID, owls[3].Bag.ID)
			is.Equal(fixtures.Bags[2].Color, owls[3].Bag.Color)

			is.NotNil(owls[4].Group)
			is.Equal(fixtures.Groups[2].ID, owls[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, owls[4].Group.Name)
			is.Empty(owls[4].Packages)
			is.NotEmpty(owls[4].Bag)
			is.Equal(fixtures.Bags[3].ID, owls[4].Bag.ID)
			is.Equal(fixtures.Bags[3].Color, owls[4].Bag.Color)

			is.NotNil(owls[5].Group)
			is.Equal(fixtures.Groups[3].ID, owls[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, owls[5].Group.Name)
			is.NotEmpty(owls[5].Packages)
			is.Len(owls[5].Packages, 3)
			is.Contains(owls[5].Packages, *fixtures.Packages[13])
			is.Contains(owls[5].Packages, *fixtures.Packages[14])
			is.Contains(owls[5].Packages, *fixtures.Packages[15])
			is.NotEmpty(owls[5].Bag)
			is.Equal(fixtures.Bags[4].ID, owls[5].Bag.ID)
			is.Equal(fixtures.Bags[4].Color, owls[5].Bag.Color)

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

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Bag", "Packages")
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
			is.NotEmpty(owls[0].Bag)
			is.Equal(fixtures.Bags[0].ID, owls[0].Bag.ID)
			is.Equal(fixtures.Bags[0].Color, owls[0].Bag.Color)

			is.Nil(owls[1].Group)
			is.NotEmpty(owls[1].Packages)
			is.Len(owls[1].Packages, 1)
			is.Contains(owls[1].Packages, *fixtures.Packages[3])
			is.NotEmpty(owls[1].Bag)
			is.Equal(fixtures.Bags[1].ID, owls[1].Bag.ID)
			is.Equal(fixtures.Bags[1].Color, owls[1].Bag.Color)

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
			is.Empty(owls[2].Bag)

			is.NotNil(owls[3].Group)
			is.Equal(fixtures.Groups[1].ID, owls[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, owls[3].Group.Name)
			is.NotEmpty(owls[3].Packages)
			is.Len(owls[3].Packages, 4)
			is.Contains(owls[3].Packages, *fixtures.Packages[9])
			is.Contains(owls[3].Packages, *fixtures.Packages[10])
			is.Contains(owls[3].Packages, *fixtures.Packages[11])
			is.Contains(owls[3].Packages, *fixtures.Packages[12])
			is.NotEmpty(owls[3].Bag)
			is.Equal(fixtures.Bags[2].ID, owls[3].Bag.ID)
			is.Equal(fixtures.Bags[2].Color, owls[3].Bag.Color)

			is.NotNil(owls[4].Group)
			is.Equal(fixtures.Groups[2].ID, owls[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, owls[4].Group.Name)
			is.Empty(owls[4].Packages)
			is.NotEmpty(owls[4].Bag)
			is.Equal(fixtures.Bags[3].ID, owls[4].Bag.ID)
			is.Equal(fixtures.Bags[3].Color, owls[4].Bag.Color)

			is.NotNil(owls[5].Group)
			is.Equal(fixtures.Groups[3].ID, owls[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, owls[5].Group.Name)
			is.NotEmpty(owls[5].Packages)
			is.Len(owls[5].Packages, 3)
			is.Contains(owls[5].Packages, *fixtures.Packages[13])
			is.Contains(owls[5].Packages, *fixtures.Packages[14])
			is.Contains(owls[5].Packages, *fixtures.Packages[15])
			is.NotEmpty(owls[5].Bag)
			is.Equal(fixtures.Bags[4].ID, owls[5].Bag.ID)
			is.Equal(fixtures.Bags[4].Color, owls[5].Bag.Color)

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

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Bag", "Packages")
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
			is.NotEmpty((*owls)[0].Bag)
			is.Equal(fixtures.Bags[0].ID, (*owls)[0].Bag.ID)
			is.Equal(fixtures.Bags[0].Color, (*owls)[0].Bag.Color)

			is.Nil((*owls)[1].Group)
			is.NotEmpty((*owls)[1].Packages)
			is.Len((*owls)[1].Packages, 1)
			is.Contains((*owls)[1].Packages, *fixtures.Packages[3])
			is.NotEmpty((*owls)[1].Bag)
			is.Equal(fixtures.Bags[1].ID, (*owls)[1].Bag.ID)
			is.Equal(fixtures.Bags[1].Color, (*owls)[1].Bag.Color)

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
			is.Empty((*owls)[2].Bag)

			is.NotNil((*owls)[3].Group)
			is.Equal(fixtures.Groups[1].ID, (*owls)[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, (*owls)[3].Group.Name)
			is.NotEmpty((*owls)[3].Packages)
			is.Len((*owls)[3].Packages, 4)
			is.Contains((*owls)[3].Packages, *fixtures.Packages[9])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[10])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[11])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[12])
			is.NotEmpty((*owls)[3].Bag)
			is.Equal(fixtures.Bags[2].ID, (*owls)[3].Bag.ID)
			is.Equal(fixtures.Bags[2].Color, (*owls)[3].Bag.Color)

			is.NotNil((*owls)[4].Group)
			is.Equal(fixtures.Groups[2].ID, (*owls)[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, (*owls)[4].Group.Name)
			is.Empty((*owls)[4].Packages)
			is.NotEmpty((*owls)[4].Bag)
			is.Equal(fixtures.Bags[3].ID, (*owls)[4].Bag.ID)
			is.Equal(fixtures.Bags[3].Color, (*owls)[4].Bag.Color)

			is.NotNil((*owls)[5].Group)
			is.Equal(fixtures.Groups[3].ID, (*owls)[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, (*owls)[5].Group.Name)
			is.NotEmpty((*owls)[5].Packages)
			is.Len((*owls)[5].Packages, 3)
			is.Contains((*owls)[5].Packages, *fixtures.Packages[13])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[14])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[15])
			is.NotEmpty((*owls)[5].Bag)
			is.Equal(fixtures.Bags[4].ID, (*owls)[5].Bag.ID)
			is.Equal(fixtures.Bags[4].Color, (*owls)[5].Bag.Color)

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

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Bag", "Packages")
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
			is.NotEmpty((*owls)[0].Bag)
			is.Equal(fixtures.Bags[0].ID, (*owls)[0].Bag.ID)
			is.Equal(fixtures.Bags[0].Color, (*owls)[0].Bag.Color)

			is.Nil((*owls)[1].Group)
			is.NotEmpty((*owls)[1].Packages)
			is.Len((*owls)[1].Packages, 1)
			is.Contains((*owls)[1].Packages, *fixtures.Packages[3])
			is.NotEmpty((*owls)[1].Bag)
			is.Equal(fixtures.Bags[1].ID, (*owls)[1].Bag.ID)
			is.Equal(fixtures.Bags[1].Color, (*owls)[1].Bag.Color)

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
			is.Empty((*owls)[2].Bag)

			is.NotNil((*owls)[3].Group)
			is.Equal(fixtures.Groups[1].ID, (*owls)[3].Group.ID)
			is.Equal(fixtures.Groups[1].Name, (*owls)[3].Group.Name)
			is.NotEmpty((*owls)[3].Packages)
			is.Len((*owls)[3].Packages, 4)
			is.Contains((*owls)[3].Packages, *fixtures.Packages[9])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[10])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[11])
			is.Contains((*owls)[3].Packages, *fixtures.Packages[12])
			is.NotEmpty((*owls)[3].Bag)
			is.Equal(fixtures.Bags[2].ID, (*owls)[3].Bag.ID)
			is.Equal(fixtures.Bags[2].Color, (*owls)[3].Bag.Color)

			is.NotNil((*owls)[4].Group)
			is.Equal(fixtures.Groups[2].ID, (*owls)[4].Group.ID)
			is.Equal(fixtures.Groups[2].Name, (*owls)[4].Group.Name)
			is.Empty((*owls)[4].Packages)
			is.NotEmpty((*owls)[4].Bag)
			is.Equal(fixtures.Bags[3].ID, (*owls)[4].Bag.ID)
			is.Equal(fixtures.Bags[3].Color, (*owls)[4].Bag.Color)

			is.NotNil((*owls)[5].Group)
			is.Equal(fixtures.Groups[3].ID, (*owls)[5].Group.ID)
			is.Equal(fixtures.Groups[3].Name, (*owls)[5].Group.Name)
			is.NotEmpty((*owls)[5].Packages)
			is.Len((*owls)[5].Packages, 3)
			is.Contains((*owls)[5].Packages, *fixtures.Packages[13])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[14])
			is.Contains((*owls)[5].Packages, *fixtures.Packages[15])
			is.NotEmpty((*owls)[5].Bag)
			is.Equal(fixtures.Bags[4].ID, (*owls)[5].Bag.ID)
			is.Equal(fixtures.Bags[4].Color, (*owls)[5].Bag.Color)

		}
		{

			owls := []Owl{}

			err := sqlxx.Preload(ctx, driver, &owls, "Group", "Packages")
			is.NoError(err)
			is.Len(owls, 0)

		}
	})
}

func TestPreload_Bag_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckBagFixtures := func(fixtures ZootopiaFixtures) {
			is.Empty(fixtures.Bags[0].Owl)
			is.Empty(fixtures.Bags[1].Owl)
			is.Empty(fixtures.Bags[2].Owl)
			is.Empty(fixtures.Bags[3].Owl)
			is.Empty(fixtures.Bags[4].Owl)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bag1 := fixtures.Bags[0]

			err := sqlxx.Preload(ctx, driver, bag1, "Owl")
			is.NoError(err)
			is.NotEmpty(bag1.Owl)
			is.Equal(fixtures.Owls[0].ID, bag1.Owl.ID)
			is.Equal(fixtures.Owls[0].Name, bag1.Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, bag1.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, bag1.Owl.FeatherColor)

			bag2 := fixtures.Bags[1]

			err = sqlxx.Preload(ctx, driver, bag2, "Owl")
			is.NoError(err)
			is.NotEmpty(bag2.Owl)
			is.Equal(fixtures.Owls[1].ID, bag2.Owl.ID)
			is.Equal(fixtures.Owls[1].Name, bag2.Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, bag2.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, bag2.Owl.FeatherColor)

			bag3 := fixtures.Bags[2]

			err = sqlxx.Preload(ctx, driver, bag3, "Owl")
			is.NoError(err)
			is.NotEmpty(bag3.Owl)
			is.Equal(fixtures.Owls[3].ID, bag3.Owl.ID)
			is.Equal(fixtures.Owls[3].Name, bag3.Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, bag3.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, bag3.Owl.FeatherColor)

			bag4 := fixtures.Bags[3]

			err = sqlxx.Preload(ctx, driver, bag4, "Owl")
			is.NoError(err)
			is.NotEmpty(bag4.Owl)
			is.Equal(fixtures.Owls[4].ID, bag4.Owl.ID)
			is.Equal(fixtures.Owls[4].Name, bag4.Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, bag4.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, bag4.Owl.FeatherColor)

			bag5 := fixtures.Bags[4]

			err = sqlxx.Preload(ctx, driver, bag5, "Owl")
			is.NoError(err)
			is.NotEmpty(bag5.Owl)
			is.Equal(fixtures.Owls[5].ID, bag5.Owl.ID)
			is.Equal(fixtures.Owls[5].Name, bag5.Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, bag5.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, bag5.Owl.FeatherColor)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bag1 := fixtures.Bags[0]

			err := sqlxx.Preload(ctx, driver, &bag1, "Owl")
			is.NoError(err)
			is.NotEmpty(bag1.Owl)
			is.Equal(fixtures.Owls[0].ID, bag1.Owl.ID)
			is.Equal(fixtures.Owls[0].Name, bag1.Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, bag1.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, bag1.Owl.FeatherColor)

			bag2 := fixtures.Bags[1]

			err = sqlxx.Preload(ctx, driver, &bag2, "Owl")
			is.NoError(err)
			is.NotEmpty(bag2.Owl)
			is.Equal(fixtures.Owls[1].ID, bag2.Owl.ID)
			is.Equal(fixtures.Owls[1].Name, bag2.Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, bag2.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, bag2.Owl.FeatherColor)

			bag3 := fixtures.Bags[2]

			err = sqlxx.Preload(ctx, driver, &bag3, "Owl")
			is.NoError(err)
			is.NotEmpty(bag3.Owl)
			is.Equal(fixtures.Owls[3].ID, bag3.Owl.ID)
			is.Equal(fixtures.Owls[3].Name, bag3.Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, bag3.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, bag3.Owl.FeatherColor)

			bag4 := fixtures.Bags[3]

			err = sqlxx.Preload(ctx, driver, &bag4, "Owl")
			is.NoError(err)
			is.NotEmpty(bag4.Owl)
			is.Equal(fixtures.Owls[4].ID, bag4.Owl.ID)
			is.Equal(fixtures.Owls[4].Name, bag4.Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, bag4.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, bag4.Owl.FeatherColor)

			bag5 := fixtures.Bags[4]

			err = sqlxx.Preload(ctx, driver, &bag5, "Owl")
			is.NoError(err)
			is.NotEmpty(bag5.Owl)
			is.Equal(fixtures.Owls[5].ID, bag5.Owl.ID)
			is.Equal(fixtures.Owls[5].Name, bag5.Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, bag5.Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, bag5.Owl.FeatherColor)

		}
	})
}

func TestPreload_Bag_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckBagFixtures := func(fixtures ZootopiaFixtures) {
			is.Empty(fixtures.Bags[0].Owl)
			is.Empty(fixtures.Bags[1].Owl)
			is.Empty(fixtures.Bags[2].Owl)
			is.Empty(fixtures.Bags[3].Owl)
			is.Empty(fixtures.Bags[4].Owl)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bags := []Bag{
				*fixtures.Bags[0],
				*fixtures.Bags[1],
				*fixtures.Bags[2],
				*fixtures.Bags[3],
				*fixtures.Bags[4],
			}

			err := sqlxx.Preload(ctx, driver, &bags, "Owl")
			is.NoError(err)
			is.Len(bags, 5)
			is.Equal(fixtures.Bags[0].ID, bags[0].ID)
			is.Equal(fixtures.Bags[1].ID, bags[1].ID)
			is.Equal(fixtures.Bags[2].ID, bags[2].ID)
			is.Equal(fixtures.Bags[3].ID, bags[3].ID)
			is.Equal(fixtures.Bags[4].ID, bags[4].ID)

			is.NotEmpty(bags[0].Owl)
			is.Equal(fixtures.Owls[0].ID, bags[0].Owl.ID)
			is.Equal(fixtures.Owls[0].Name, bags[0].Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, bags[0].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, bags[0].Owl.FeatherColor)

			is.NotEmpty(bags[1].Owl)
			is.Equal(fixtures.Owls[1].ID, bags[1].Owl.ID)
			is.Equal(fixtures.Owls[1].Name, bags[1].Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, bags[1].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, bags[1].Owl.FeatherColor)

			is.NotEmpty(bags[2].Owl)
			is.Equal(fixtures.Owls[3].ID, bags[2].Owl.ID)
			is.Equal(fixtures.Owls[3].Name, bags[2].Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, bags[2].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, bags[2].Owl.FeatherColor)

			is.NotEmpty(bags[3].Owl)
			is.Equal(fixtures.Owls[4].ID, bags[3].Owl.ID)
			is.Equal(fixtures.Owls[4].Name, bags[3].Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, bags[3].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, bags[3].Owl.FeatherColor)

			is.NotEmpty(bags[4].Owl)
			is.Equal(fixtures.Owls[5].ID, bags[4].Owl.ID)
			is.Equal(fixtures.Owls[5].Name, bags[4].Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, bags[4].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, bags[4].Owl.FeatherColor)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bags := []*Bag{
				fixtures.Bags[0],
				fixtures.Bags[1],
				fixtures.Bags[2],
				fixtures.Bags[3],
				fixtures.Bags[4],
			}

			err := sqlxx.Preload(ctx, driver, &bags, "Owl")
			is.NoError(err)
			is.Len(bags, 5)
			is.Equal(fixtures.Bags[0].ID, bags[0].ID)
			is.Equal(fixtures.Bags[1].ID, bags[1].ID)
			is.Equal(fixtures.Bags[2].ID, bags[2].ID)
			is.Equal(fixtures.Bags[3].ID, bags[3].ID)
			is.Equal(fixtures.Bags[4].ID, bags[4].ID)

			is.NotEmpty(bags[0].Owl)
			is.Equal(fixtures.Owls[0].ID, bags[0].Owl.ID)
			is.Equal(fixtures.Owls[0].Name, bags[0].Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, bags[0].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, bags[0].Owl.FeatherColor)

			is.NotEmpty(bags[1].Owl)
			is.Equal(fixtures.Owls[1].ID, bags[1].Owl.ID)
			is.Equal(fixtures.Owls[1].Name, bags[1].Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, bags[1].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, bags[1].Owl.FeatherColor)

			is.NotEmpty(bags[2].Owl)
			is.Equal(fixtures.Owls[3].ID, bags[2].Owl.ID)
			is.Equal(fixtures.Owls[3].Name, bags[2].Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, bags[2].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, bags[2].Owl.FeatherColor)

			is.NotEmpty(bags[3].Owl)
			is.Equal(fixtures.Owls[4].ID, bags[3].Owl.ID)
			is.Equal(fixtures.Owls[4].Name, bags[3].Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, bags[3].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, bags[3].Owl.FeatherColor)

			is.NotEmpty(bags[4].Owl)
			is.Equal(fixtures.Owls[5].ID, bags[4].Owl.ID)
			is.Equal(fixtures.Owls[5].Name, bags[4].Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, bags[4].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, bags[4].Owl.FeatherColor)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bags := &[]Bag{
				*fixtures.Bags[0],
				*fixtures.Bags[1],
				*fixtures.Bags[2],
				*fixtures.Bags[3],
				*fixtures.Bags[4],
			}

			err := sqlxx.Preload(ctx, driver, &bags, "Owl")
			is.NoError(err)
			is.Len((*bags), 5)
			is.Equal(fixtures.Bags[0].ID, (*bags)[0].ID)
			is.Equal(fixtures.Bags[1].ID, (*bags)[1].ID)
			is.Equal(fixtures.Bags[2].ID, (*bags)[2].ID)
			is.Equal(fixtures.Bags[3].ID, (*bags)[3].ID)
			is.Equal(fixtures.Bags[4].ID, (*bags)[4].ID)

			is.NotEmpty((*bags)[0].Owl)
			is.Equal(fixtures.Owls[0].ID, (*bags)[0].Owl.ID)
			is.Equal(fixtures.Owls[0].Name, (*bags)[0].Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, (*bags)[0].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, (*bags)[0].Owl.FeatherColor)

			is.NotEmpty((*bags)[1].Owl)
			is.Equal(fixtures.Owls[1].ID, (*bags)[1].Owl.ID)
			is.Equal(fixtures.Owls[1].Name, (*bags)[1].Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, (*bags)[1].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, (*bags)[1].Owl.FeatherColor)

			is.NotEmpty((*bags)[2].Owl)
			is.Equal(fixtures.Owls[3].ID, (*bags)[2].Owl.ID)
			is.Equal(fixtures.Owls[3].Name, (*bags)[2].Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, (*bags)[2].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, (*bags)[2].Owl.FeatherColor)

			is.NotEmpty((*bags)[3].Owl)
			is.Equal(fixtures.Owls[4].ID, (*bags)[3].Owl.ID)
			is.Equal(fixtures.Owls[4].Name, (*bags)[3].Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, (*bags)[3].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, (*bags)[3].Owl.FeatherColor)

			is.NotEmpty((*bags)[4].Owl)
			is.Equal(fixtures.Owls[5].ID, (*bags)[4].Owl.ID)
			is.Equal(fixtures.Owls[5].Name, (*bags)[4].Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, (*bags)[4].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, (*bags)[4].Owl.FeatherColor)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckBagFixtures(fixtures)

			bags := &[]*Bag{
				fixtures.Bags[0],
				fixtures.Bags[1],
				fixtures.Bags[2],
				fixtures.Bags[3],
				fixtures.Bags[4],
			}

			err := sqlxx.Preload(ctx, driver, &bags, "Owl")
			is.NoError(err)
			is.Len((*bags), 5)
			is.Equal(fixtures.Bags[0].ID, (*bags)[0].ID)
			is.Equal(fixtures.Bags[1].ID, (*bags)[1].ID)
			is.Equal(fixtures.Bags[2].ID, (*bags)[2].ID)
			is.Equal(fixtures.Bags[3].ID, (*bags)[3].ID)
			is.Equal(fixtures.Bags[4].ID, (*bags)[4].ID)

			is.NotEmpty((*bags)[0].Owl)
			is.Equal(fixtures.Owls[0].ID, (*bags)[0].Owl.ID)
			is.Equal(fixtures.Owls[0].Name, (*bags)[0].Owl.Name)
			is.Equal(fixtures.Owls[0].FavoriteFood, (*bags)[0].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[0].FeatherColor, (*bags)[0].Owl.FeatherColor)

			is.NotEmpty((*bags)[1].Owl)
			is.Equal(fixtures.Owls[1].ID, (*bags)[1].Owl.ID)
			is.Equal(fixtures.Owls[1].Name, (*bags)[1].Owl.Name)
			is.Equal(fixtures.Owls[1].FavoriteFood, (*bags)[1].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[1].FeatherColor, (*bags)[1].Owl.FeatherColor)

			is.NotEmpty((*bags)[2].Owl)
			is.Equal(fixtures.Owls[3].ID, (*bags)[2].Owl.ID)
			is.Equal(fixtures.Owls[3].Name, (*bags)[2].Owl.Name)
			is.Equal(fixtures.Owls[3].FavoriteFood, (*bags)[2].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[3].FeatherColor, (*bags)[2].Owl.FeatherColor)

			is.NotEmpty((*bags)[3].Owl)
			is.Equal(fixtures.Owls[4].ID, (*bags)[3].Owl.ID)
			is.Equal(fixtures.Owls[4].Name, (*bags)[3].Owl.Name)
			is.Equal(fixtures.Owls[4].FavoriteFood, (*bags)[3].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[4].FeatherColor, (*bags)[3].Owl.FeatherColor)

			is.NotEmpty((*bags)[4].Owl)
			is.Equal(fixtures.Owls[5].ID, (*bags)[4].Owl.ID)
			is.Equal(fixtures.Owls[5].Name, (*bags)[4].Owl.Name)
			is.Equal(fixtures.Owls[5].FavoriteFood, (*bags)[4].Owl.FavoriteFood)
			is.Equal(fixtures.Owls[5].FeatherColor, (*bags)[4].Owl.FeatherColor)

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

func TestPreload_Human_One(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckHumanFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Humans[0].Cat)
			is.Nil(fixtures.Humans[1].Cat)
			is.Nil(fixtures.Humans[2].Cat)
			is.Nil(fixtures.Humans[3].Cat)
			is.Nil(fixtures.Humans[4].Cat)
			is.Nil(fixtures.Humans[5].Cat)
			is.Nil(fixtures.Humans[6].Cat)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			human1 := fixtures.Humans[0]

			err := sqlxx.Preload(ctx, driver, human1, "Cat")
			is.NoError(err)
			is.NotNil(human1.Cat)
			is.Equal(fixtures.Cats[0].ID, human1.Cat.ID)
			is.Equal(fixtures.Cats[0].Name, human1.Cat.Name)

			human2 := fixtures.Humans[1]

			err = sqlxx.Preload(ctx, driver, human2, "Cat")
			is.NoError(err)
			is.NotNil(human2.Cat)
			is.Equal(fixtures.Cats[1].ID, human2.Cat.ID)
			is.Equal(fixtures.Cats[1].Name, human2.Cat.Name)

			human3 := fixtures.Humans[2]

			err = sqlxx.Preload(ctx, driver, human3, "Cat")
			is.NoError(err)
			is.Nil(human3.Cat)

			human7 := fixtures.Humans[6]

			err = sqlxx.Preload(ctx, driver, human7, "Cat")
			is.NoError(err)
			is.NotNil(human7.Cat)
			is.Equal(fixtures.Cats[6].ID, human7.Cat.ID)
			is.Equal(fixtures.Cats[6].Name, human7.Cat.Name)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			human1 := fixtures.Humans[0]

			err := sqlxx.Preload(ctx, driver, &human1, "Cat")
			is.NoError(err)
			is.NotNil(human1.Cat)
			is.Equal(fixtures.Cats[0].ID, human1.Cat.ID)
			is.Equal(fixtures.Cats[0].Name, human1.Cat.Name)

			human2 := fixtures.Humans[1]

			err = sqlxx.Preload(ctx, driver, &human2, "Cat")
			is.NoError(err)
			is.NotNil(human2.Cat)
			is.Equal(fixtures.Cats[1].ID, human2.Cat.ID)
			is.Equal(fixtures.Cats[1].Name, human2.Cat.Name)

			human3 := fixtures.Humans[2]

			err = sqlxx.Preload(ctx, driver, &human3, "Cat")
			is.NoError(err)
			is.Nil(human3.Cat)

			human7 := fixtures.Humans[6]

			err = sqlxx.Preload(ctx, driver, &human7, "Cat")
			is.NoError(err)
			is.NotNil(human7.Cat)
			is.Equal(fixtures.Cats[6].ID, human7.Cat.ID)
			is.Equal(fixtures.Cats[6].Name, human7.Cat.Name)

		}
	})
}

func TestPreload_Human_Many(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		CheckHumanFixtures := func(fixtures ZootopiaFixtures) {
			is.Nil(fixtures.Humans[0].Cat)
			is.Nil(fixtures.Humans[1].Cat)
			is.Nil(fixtures.Humans[2].Cat)
			is.Nil(fixtures.Humans[3].Cat)
			is.Nil(fixtures.Humans[4].Cat)
			is.Nil(fixtures.Humans[5].Cat)
			is.Nil(fixtures.Humans[6].Cat)
		}

		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			humans := []Human{
				*fixtures.Humans[0],
				*fixtures.Humans[1],
				*fixtures.Humans[2],
				*fixtures.Humans[3],
				*fixtures.Humans[4],
				*fixtures.Humans[5],
				*fixtures.Humans[6],
			}

			err := sqlxx.Preload(ctx, driver, &humans, "Cat")
			is.NoError(err)
			is.Len(humans, 7)
			is.Equal(fixtures.Humans[0].ID, humans[0].ID)
			is.Equal(fixtures.Humans[1].ID, humans[1].ID)
			is.Equal(fixtures.Humans[2].ID, humans[2].ID)
			is.Equal(fixtures.Humans[3].ID, humans[3].ID)
			is.Equal(fixtures.Humans[4].ID, humans[4].ID)
			is.Equal(fixtures.Humans[5].ID, humans[5].ID)
			is.Equal(fixtures.Humans[6].ID, humans[6].ID)

			is.NotNil(humans[0].Cat)
			is.Equal(fixtures.Cats[0].ID, humans[0].Cat.ID)
			is.Equal(fixtures.Cats[0].Name, humans[0].Cat.Name)

			is.NotNil(humans[1].Cat)
			is.Equal(fixtures.Cats[1].ID, humans[1].Cat.ID)
			is.Equal(fixtures.Cats[1].Name, humans[1].Cat.Name)

			is.Nil(humans[2].Cat)

			is.NotNil(humans[3].Cat)
			is.Equal(fixtures.Cats[3].ID, humans[3].Cat.ID)
			is.Equal(fixtures.Cats[3].Name, humans[3].Cat.Name)

			is.Nil(humans[4].Cat)

			is.Nil(humans[5].Cat)

			is.NotNil(humans[6].Cat)
			is.Equal(fixtures.Cats[6].ID, humans[6].Cat.ID)
			is.Equal(fixtures.Cats[6].Name, humans[6].Cat.Name)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			humans := []*Human{
				fixtures.Humans[0],
				fixtures.Humans[1],
				fixtures.Humans[2],
				fixtures.Humans[3],
				fixtures.Humans[4],
				fixtures.Humans[5],
				fixtures.Humans[6],
			}

			err := sqlxx.Preload(ctx, driver, &humans, "Cat")
			is.NoError(err)
			is.Len(humans, 7)
			is.Equal(fixtures.Humans[0].ID, humans[0].ID)
			is.Equal(fixtures.Humans[1].ID, humans[1].ID)
			is.Equal(fixtures.Humans[2].ID, humans[2].ID)
			is.Equal(fixtures.Humans[3].ID, humans[3].ID)
			is.Equal(fixtures.Humans[4].ID, humans[4].ID)
			is.Equal(fixtures.Humans[5].ID, humans[5].ID)
			is.Equal(fixtures.Humans[6].ID, humans[6].ID)

			is.NotNil(humans[0].Cat)
			is.Equal(fixtures.Cats[0].ID, humans[0].Cat.ID)
			is.Equal(fixtures.Cats[0].Name, humans[0].Cat.Name)

			is.NotNil(humans[1].Cat)
			is.Equal(fixtures.Cats[1].ID, humans[1].Cat.ID)
			is.Equal(fixtures.Cats[1].Name, humans[1].Cat.Name)

			is.Nil(humans[2].Cat)

			is.NotNil(humans[3].Cat)
			is.Equal(fixtures.Cats[3].ID, humans[3].Cat.ID)
			is.Equal(fixtures.Cats[3].Name, humans[3].Cat.Name)

			is.Nil(humans[4].Cat)

			is.Nil(humans[5].Cat)

			is.NotNil(humans[6].Cat)
			is.Equal(fixtures.Cats[6].ID, humans[6].Cat.ID)
			is.Equal(fixtures.Cats[6].Name, humans[6].Cat.Name)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			humans := &[]Human{
				*fixtures.Humans[0],
				*fixtures.Humans[1],
				*fixtures.Humans[2],
				*fixtures.Humans[3],
				*fixtures.Humans[4],
				*fixtures.Humans[5],
				*fixtures.Humans[6],
			}

			err := sqlxx.Preload(ctx, driver, &humans, "Cat")
			is.NoError(err)
			is.Len((*humans), 7)
			is.Equal(fixtures.Humans[0].ID, (*humans)[0].ID)
			is.Equal(fixtures.Humans[1].ID, (*humans)[1].ID)
			is.Equal(fixtures.Humans[2].ID, (*humans)[2].ID)
			is.Equal(fixtures.Humans[3].ID, (*humans)[3].ID)
			is.Equal(fixtures.Humans[4].ID, (*humans)[4].ID)
			is.Equal(fixtures.Humans[5].ID, (*humans)[5].ID)
			is.Equal(fixtures.Humans[6].ID, (*humans)[6].ID)

			is.NotNil((*humans)[0].Cat)
			is.Equal(fixtures.Cats[0].ID, (*humans)[0].Cat.ID)
			is.Equal(fixtures.Cats[0].Name, (*humans)[0].Cat.Name)

			is.NotNil((*humans)[1].Cat)
			is.Equal(fixtures.Cats[1].ID, (*humans)[1].Cat.ID)
			is.Equal(fixtures.Cats[1].Name, (*humans)[1].Cat.Name)

			is.Nil((*humans)[2].Cat)

			is.NotNil((*humans)[3].Cat)
			is.Equal(fixtures.Cats[3].ID, (*humans)[3].Cat.ID)
			is.Equal(fixtures.Cats[3].Name, (*humans)[3].Cat.Name)

			is.Nil((*humans)[4].Cat)

			is.Nil((*humans)[5].Cat)

			is.NotNil((*humans)[6].Cat)
			is.Equal(fixtures.Cats[6].ID, (*humans)[6].Cat.ID)
			is.Equal(fixtures.Cats[6].Name, (*humans)[6].Cat.Name)

		}
		{

			fixtures := GenerateZootopiaFixtures(ctx, driver, is)
			CheckHumanFixtures(fixtures)

			humans := &[]*Human{
				fixtures.Humans[0],
				fixtures.Humans[1],
				fixtures.Humans[2],
				fixtures.Humans[3],
				fixtures.Humans[4],
				fixtures.Humans[5],
				fixtures.Humans[6],
			}

			err := sqlxx.Preload(ctx, driver, &humans, "Cat")
			is.NoError(err)
			is.Len((*humans), 7)
			is.Equal(fixtures.Humans[0].ID, (*humans)[0].ID)
			is.Equal(fixtures.Humans[1].ID, (*humans)[1].ID)
			is.Equal(fixtures.Humans[2].ID, (*humans)[2].ID)
			is.Equal(fixtures.Humans[3].ID, (*humans)[3].ID)
			is.Equal(fixtures.Humans[4].ID, (*humans)[4].ID)
			is.Equal(fixtures.Humans[5].ID, (*humans)[5].ID)
			is.Equal(fixtures.Humans[6].ID, (*humans)[6].ID)

			is.NotNil((*humans)[0].Cat)
			is.Equal(fixtures.Cats[0].ID, (*humans)[0].Cat.ID)
			is.Equal(fixtures.Cats[0].Name, (*humans)[0].Cat.Name)

			is.NotNil((*humans)[1].Cat)
			is.Equal(fixtures.Cats[1].ID, (*humans)[1].Cat.ID)
			is.Equal(fixtures.Cats[1].Name, (*humans)[1].Cat.Name)

			is.Nil((*humans)[2].Cat)

			is.NotNil((*humans)[3].Cat)
			is.Equal(fixtures.Cats[3].ID, (*humans)[3].Cat.ID)
			is.Equal(fixtures.Cats[3].Name, (*humans)[3].Cat.Name)

			is.Nil((*humans)[4].Cat)

			is.Nil((*humans)[5].Cat)

			is.NotNil((*humans)[6].Cat)
			is.Equal(fixtures.Cats[6].ID, (*humans)[6].Cat.ID)
			is.Equal(fixtures.Cats[6].Name, (*humans)[6].Cat.Name)

		}
		{

			humans := []Human{}

			err := sqlxx.Preload(ctx, driver, &humans, "Cat")
			is.NoError(err)
			is.Len(humans, 0)

		}
	})
}

func TestPreload_ExoCloud_MultiLevel(t *testing.T) {
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

			err := sqlxx.Preload(ctx, driver, region1, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region1.Buckets)
			is.NotEmpty((*region1.Buckets))
			is.Len((*region1.Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*region1.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*region1.Buckets)[1].ID)
			is.Len((*region1.Buckets)[0].Files, 12)
			is.Len((*region1.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*region1.Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*region1.Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[2]).ID, (*region1.Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[3]).ID, (*region1.Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[4]).ID, (*region1.Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[5]).ID, (*region1.Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[6]).ID, (*region1.Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[7]).ID, (*region1.Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[8]).ID, (*region1.Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[19]).ID, (*region1.Buckets)[0].Files[9].ID)
			is.Equal((*fixtures.Files[20]).ID, (*region1.Buckets)[0].Files[10].ID)
			is.Equal((*fixtures.Files[21]).ID, (*region1.Buckets)[0].Files[11].ID)
			is.Len((*region1.Buckets)[0].Directories, 16)
			is.Len((*region1.Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[0]).ID, (*region1.Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[1]).ID, (*region1.Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[2]).ID, (*region1.Buckets)[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[3]).ID, (*region1.Buckets)[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[4]).ID, (*region1.Buckets)[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[5]).ID, (*region1.Buckets)[0].Directories[5].ID)
			is.Equal((*fixtures.Directories[6]).ID, (*region1.Buckets)[0].Directories[6].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*region1.Buckets)[0].Directories[7].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*region1.Buckets)[0].Directories[8].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*region1.Buckets)[0].Directories[9].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*region1.Buckets)[0].Directories[11].ID)
			is.Equal((*fixtures.Directories[12]).ID, (*region1.Buckets)[0].Directories[12].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*region1.Buckets)[0].Directories[13].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*region1.Buckets)[0].Directories[14].ID)
			is.Equal((*fixtures.Directories[15]).ID, (*region1.Buckets)[0].Directories[15].ID)
			is.NotEmpty((*region1.Buckets)[0].Region)
			is.NotEmpty((*region1.Buckets)[1].Region)
			is.Equal(region1.ID, (*region1.Buckets)[0].Region.ID)
			is.Equal(region1.ID, (*region1.Buckets)[1].Region.ID)

			region2 := fixtures.Regions[1]

			err = sqlxx.Preload(ctx, driver, region2, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region2.Buckets)
			is.Empty((*region2.Buckets))
			is.Len((*region2.Buckets), 0)

			region3 := fixtures.Regions[2]

			err = sqlxx.Preload(ctx, driver, region3, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region3.Buckets)
			is.NotEmpty((*region3.Buckets))
			is.Len((*region3.Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*region3.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*region3.Buckets)[1].ID)
			is.Len((*region3.Buckets)[0].Files, 10)
			is.Len((*region3.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*region3.Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*region3.Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[11]).ID, (*region3.Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[12]).ID, (*region3.Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[13]).ID, (*region3.Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[14]).ID, (*region3.Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[15]).ID, (*region3.Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[16]).ID, (*region3.Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[17]).ID, (*region3.Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[18]).ID, (*region3.Buckets)[0].Files[9].ID)
			is.Len((*region3.Buckets)[0].Directories, 3)
			is.Len((*region3.Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[16]).ID, (*region3.Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[17]).ID, (*region3.Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[18]).ID, (*region3.Buckets)[0].Directories[2].ID)
			is.NotEmpty((*region3.Buckets)[0].Region)
			is.NotEmpty((*region3.Buckets)[1].Region)
			is.Equal(region3.ID, (*region3.Buckets)[0].Region.ID)
			is.Equal(region3.ID, (*region3.Buckets)[1].Region.ID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			region1 := fixtures.Regions[0]

			err := sqlxx.Preload(ctx, driver, &region1, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region1.Buckets)
			is.NotEmpty((*region1.Buckets))
			is.Len((*region1.Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*region1.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*region1.Buckets)[1].ID)
			is.Len((*region1.Buckets)[0].Files, 12)
			is.Len((*region1.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*region1.Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*region1.Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[2]).ID, (*region1.Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[3]).ID, (*region1.Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[4]).ID, (*region1.Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[5]).ID, (*region1.Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[6]).ID, (*region1.Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[7]).ID, (*region1.Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[8]).ID, (*region1.Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[19]).ID, (*region1.Buckets)[0].Files[9].ID)
			is.Equal((*fixtures.Files[20]).ID, (*region1.Buckets)[0].Files[10].ID)
			is.Equal((*fixtures.Files[21]).ID, (*region1.Buckets)[0].Files[11].ID)
			is.Len((*region1.Buckets)[0].Directories, 16)
			is.Len((*region1.Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[0]).ID, (*region1.Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[1]).ID, (*region1.Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[2]).ID, (*region1.Buckets)[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[3]).ID, (*region1.Buckets)[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[4]).ID, (*region1.Buckets)[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[5]).ID, (*region1.Buckets)[0].Directories[5].ID)
			is.Equal((*fixtures.Directories[6]).ID, (*region1.Buckets)[0].Directories[6].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*region1.Buckets)[0].Directories[7].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*region1.Buckets)[0].Directories[8].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*region1.Buckets)[0].Directories[9].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*region1.Buckets)[0].Directories[11].ID)
			is.Equal((*fixtures.Directories[12]).ID, (*region1.Buckets)[0].Directories[12].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*region1.Buckets)[0].Directories[13].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*region1.Buckets)[0].Directories[14].ID)
			is.Equal((*fixtures.Directories[15]).ID, (*region1.Buckets)[0].Directories[15].ID)
			is.NotEmpty((*region1.Buckets)[0].Region)
			is.NotEmpty((*region1.Buckets)[1].Region)
			is.Equal(region1.ID, (*region1.Buckets)[0].Region.ID)
			is.Equal(region1.ID, (*region1.Buckets)[1].Region.ID)

			region2 := fixtures.Regions[1]

			err = sqlxx.Preload(ctx, driver, &region2, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region2.Buckets)
			is.Empty((*region2.Buckets))
			is.Len((*region2.Buckets), 0)

			region3 := fixtures.Regions[2]

			err = sqlxx.Preload(ctx, driver, &region3, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.NotNil(region3.Buckets)
			is.NotEmpty((*region3.Buckets))
			is.Len((*region3.Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*region3.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*region3.Buckets)[1].ID)
			is.Len((*region3.Buckets)[0].Files, 10)
			is.Len((*region3.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*region3.Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*region3.Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[11]).ID, (*region3.Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[12]).ID, (*region3.Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[13]).ID, (*region3.Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[14]).ID, (*region3.Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[15]).ID, (*region3.Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[16]).ID, (*region3.Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[17]).ID, (*region3.Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[18]).ID, (*region3.Buckets)[0].Files[9].ID)
			is.Len((*region3.Buckets)[0].Directories, 3)
			is.Len((*region3.Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[16]).ID, (*region3.Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[17]).ID, (*region3.Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[18]).ID, (*region3.Buckets)[0].Directories[2].ID)
			is.NotEmpty((*region3.Buckets)[0].Region)
			is.NotEmpty((*region3.Buckets)[1].Region)
			is.Equal(region3.ID, (*region3.Buckets)[0].Region.ID)
			is.Equal(region3.ID, (*region3.Buckets)[1].Region.ID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := []ExoRegion{
				*fixtures.Regions[0],
				*fixtures.Regions[1],
				*fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*regions[0].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*regions[0].Buckets)[1].ID)
			is.Len((*regions[0].Buckets)[0].Files, 12)
			is.Len((*regions[0].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*regions[0].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[2]).ID, (*regions[0].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[3]).ID, (*regions[0].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[4]).ID, (*regions[0].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[5]).ID, (*regions[0].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[6]).ID, (*regions[0].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[7]).ID, (*regions[0].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[8]).ID, (*regions[0].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[19]).ID, (*regions[0].Buckets)[0].Files[9].ID)
			is.Equal((*fixtures.Files[20]).ID, (*regions[0].Buckets)[0].Files[10].ID)
			is.Equal((*fixtures.Files[21]).ID, (*regions[0].Buckets)[0].Files[11].ID)
			is.Len((*regions[0].Buckets)[0].Directories, 16)
			is.Len((*regions[0].Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[0]).ID, (*regions[0].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[1]).ID, (*regions[0].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[2]).ID, (*regions[0].Buckets)[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[3]).ID, (*regions[0].Buckets)[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[4]).ID, (*regions[0].Buckets)[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[5]).ID, (*regions[0].Buckets)[0].Directories[5].ID)
			is.Equal((*fixtures.Directories[6]).ID, (*regions[0].Buckets)[0].Directories[6].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*regions[0].Buckets)[0].Directories[7].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*regions[0].Buckets)[0].Directories[8].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*regions[0].Buckets)[0].Directories[9].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*regions[0].Buckets)[0].Directories[11].ID)
			is.Equal((*fixtures.Directories[12]).ID, (*regions[0].Buckets)[0].Directories[12].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*regions[0].Buckets)[0].Directories[13].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*regions[0].Buckets)[0].Directories[14].ID)
			is.Equal((*fixtures.Directories[15]).ID, (*regions[0].Buckets)[0].Directories[15].ID)
			is.NotEmpty((*regions[0].Buckets)[0].Region)
			is.NotEmpty((*regions[0].Buckets)[1].Region)
			is.Equal(regions[0].ID, (*regions[0].Buckets)[0].Region.ID)
			is.Equal(regions[0].ID, (*regions[0].Buckets)[1].Region.ID)

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*regions[2].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*regions[2].Buckets)[1].ID)
			is.Len((*regions[2].Buckets)[0].Files, 10)
			is.Len((*regions[2].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*regions[2].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*regions[2].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[11]).ID, (*regions[2].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[12]).ID, (*regions[2].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[13]).ID, (*regions[2].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[14]).ID, (*regions[2].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[15]).ID, (*regions[2].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[16]).ID, (*regions[2].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[17]).ID, (*regions[2].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[18]).ID, (*regions[2].Buckets)[0].Files[9].ID)
			is.Len((*regions[2].Buckets)[0].Directories, 3)
			is.Len((*regions[2].Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[16]).ID, (*regions[2].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[17]).ID, (*regions[2].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[18]).ID, (*regions[2].Buckets)[0].Directories[2].ID)
			is.NotEmpty((*regions[2].Buckets)[0].Region)
			is.NotEmpty((*regions[2].Buckets)[1].Region)
			is.Equal(regions[2].ID, (*regions[2].Buckets)[0].Region.ID)
			is.Equal(regions[2].ID, (*regions[2].Buckets)[1].Region.ID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := []*ExoRegion{
				fixtures.Regions[0],
				fixtures.Regions[1],
				fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions, "Buckets.Files", "Buckets.Directories", "Buckets.Region")
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*regions[0].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*regions[0].Buckets)[1].ID)
			is.Len((*regions[0].Buckets)[0].Files, 12)
			is.Len((*regions[0].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*regions[0].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[2]).ID, (*regions[0].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[3]).ID, (*regions[0].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[4]).ID, (*regions[0].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[5]).ID, (*regions[0].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[6]).ID, (*regions[0].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[7]).ID, (*regions[0].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[8]).ID, (*regions[0].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[19]).ID, (*regions[0].Buckets)[0].Files[9].ID)
			is.Equal((*fixtures.Files[20]).ID, (*regions[0].Buckets)[0].Files[10].ID)
			is.Equal((*fixtures.Files[21]).ID, (*regions[0].Buckets)[0].Files[11].ID)
			is.Len((*regions[0].Buckets)[0].Directories, 16)
			is.Len((*regions[0].Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[0]).ID, (*regions[0].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[1]).ID, (*regions[0].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[2]).ID, (*regions[0].Buckets)[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[3]).ID, (*regions[0].Buckets)[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[4]).ID, (*regions[0].Buckets)[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[5]).ID, (*regions[0].Buckets)[0].Directories[5].ID)
			is.Equal((*fixtures.Directories[6]).ID, (*regions[0].Buckets)[0].Directories[6].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*regions[0].Buckets)[0].Directories[7].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*regions[0].Buckets)[0].Directories[8].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*regions[0].Buckets)[0].Directories[9].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*regions[0].Buckets)[0].Directories[11].ID)
			is.Equal((*fixtures.Directories[12]).ID, (*regions[0].Buckets)[0].Directories[12].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*regions[0].Buckets)[0].Directories[13].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*regions[0].Buckets)[0].Directories[14].ID)
			is.Equal((*fixtures.Directories[15]).ID, (*regions[0].Buckets)[0].Directories[15].ID)
			is.NotEmpty((*regions[0].Buckets)[0].Region)
			is.NotEmpty((*regions[0].Buckets)[1].Region)
			is.Equal(regions[0].ID, (*regions[0].Buckets)[0].Region.ID)
			is.Equal(regions[0].ID, (*regions[0].Buckets)[1].Region.ID)

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*regions[2].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*regions[2].Buckets)[1].ID)
			is.Len((*regions[2].Buckets)[0].Files, 10)
			is.Len((*regions[2].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*regions[2].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*regions[2].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[11]).ID, (*regions[2].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[12]).ID, (*regions[2].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[13]).ID, (*regions[2].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[14]).ID, (*regions[2].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[15]).ID, (*regions[2].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[16]).ID, (*regions[2].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[17]).ID, (*regions[2].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[18]).ID, (*regions[2].Buckets)[0].Files[9].ID)
			is.Len((*regions[2].Buckets)[0].Directories, 3)
			is.Len((*regions[2].Buckets)[1].Directories, 0)
			is.Equal((*fixtures.Directories[16]).ID, (*regions[2].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[17]).ID, (*regions[2].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[18]).ID, (*regions[2].Buckets)[0].Directories[2].ID)
			is.NotEmpty((*regions[2].Buckets)[0].Region)
			is.NotEmpty((*regions[2].Buckets)[1].Region)
			is.Equal(regions[2].ID, (*regions[2].Buckets)[0].Region.ID)
			is.Equal(regions[2].ID, (*regions[2].Buckets)[1].Region.ID)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			region1 := fixtures.Regions[0]
			err := sqlxx.Preload(ctx, driver, region1,
				"Buckets.Files",
				"Buckets.Files.Owner",
				"Buckets.Files.Owner.Profile",
				"Buckets.Files.Owner.Profile.Avatar",
				"Buckets.Files.Chunks",
				"Buckets.Files.Chunks.Mode",
				"Buckets.Files.Chunks.Signature",
			)
			is.NoError(err)
			is.NotNil(region1.Buckets)
			is.NotEmpty((*region1.Buckets))
			is.Len((*region1.Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*region1.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*region1.Buckets)[1].ID)
			is.Len((*region1.Buckets)[0].Files, 12)
			is.Len((*region1.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*region1.Buckets)[0].Files[0].ID)
			is.NotNil((*region1.Buckets)[0].Files[0].Owner)
			is.NotEmpty(*((*region1.Buckets)[0].Files[0].Owner))
			is.Equal((*fixtures.Users[13]).ID, (*region1.Buckets)[0].Files[0].Owner.ID)
			is.NotNil((*region1.Buckets)[0].Files[0].Owner.Profile)
			is.NotEmpty(*((*region1.Buckets)[0].Files[0].Owner.Profile))
			is.Equal((*fixtures.Profiles[13]).ID, (*region1.Buckets)[0].Files[0].Owner.Profile.ID)
			is.NotNil((*region1.Buckets)[0].Files[0].Owner.Profile.Avatar)
			is.NotEmpty(*((*region1.Buckets)[0].Files[0].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[13]).ID, (*region1.Buckets)[0].Files[0].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region1.Buckets)[0].Files[0].Chunks)
			is.Len((*region1.Buckets)[0].Files[0].Chunks, 3)
			is.Equal((*fixtures.Chunks[0]).Hash, (*region1.Buckets)[0].Files[0].Chunks[0].Hash)
			is.Equal((*fixtures.Chunks[1]).Hash, (*region1.Buckets)[0].Files[0].Chunks[1].Hash)
			is.Equal((*fixtures.Chunks[2]).Hash, (*region1.Buckets)[0].Files[0].Chunks[2].Hash)
			is.NotNil((*region1.Buckets)[0].Files[0].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*region1.Buckets)[0].Files[0].Chunks[0].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[0].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*region1.Buckets)[0].Files[0].Chunks[1].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[0].Chunks[2].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*region1.Buckets)[0].Files[0].Chunks[2].Mode.ID)
			is.Nil((*region1.Buckets)[0].Files[0].Chunks[0].Signature)
			is.Nil((*region1.Buckets)[0].Files[0].Chunks[1].Signature)
			is.Nil((*region1.Buckets)[0].Files[0].Chunks[2].Signature)
			is.Equal((*fixtures.Files[8]).ID, (*region1.Buckets)[0].Files[8].ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Owner)
			is.NotEmpty(*((*region1.Buckets)[0].Files[8].Owner))
			is.Equal((*fixtures.Users[18]).ID, (*region1.Buckets)[0].Files[8].Owner.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Owner.Profile)
			is.NotEmpty(*((*region1.Buckets)[0].Files[8].Owner.Profile))
			is.Equal((*fixtures.Profiles[18]).ID, (*region1.Buckets)[0].Files[8].Owner.Profile.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Owner.Profile.Avatar)
			is.NotEmpty(*((*region1.Buckets)[0].Files[8].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[18]).ID, (*region1.Buckets)[0].Files[8].Owner.Profile.Avatar.ID)
			is.Len((*region1.Buckets)[0].Files[8].Chunks, 4)
			is.Equal((*fixtures.Chunks[24]).Hash, (*region1.Buckets)[0].Files[8].Chunks[0].Hash)
			is.Equal((*fixtures.Chunks[25]).Hash, (*region1.Buckets)[0].Files[8].Chunks[1].Hash)
			is.Equal((*fixtures.Chunks[26]).Hash, (*region1.Buckets)[0].Files[8].Chunks[2].Hash)
			is.Equal((*fixtures.Chunks[27]).Hash, (*region1.Buckets)[0].Files[8].Chunks[3].Hash)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*region1.Buckets)[0].Files[8].Chunks[0].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*region1.Buckets)[0].Files[8].Chunks[1].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[2].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*region1.Buckets)[0].Files[8].Chunks[2].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[3].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*region1.Buckets)[0].Files[8].Chunks[3].Mode.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[0].Signature)
			is.Equal((*fixtures.Signatures[0]).ID, (*region1.Buckets)[0].Files[8].Chunks[0].Signature.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[1].Signature)
			is.Equal((*fixtures.Signatures[1]).ID, (*region1.Buckets)[0].Files[8].Chunks[1].Signature.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[2].Signature)
			is.Equal((*fixtures.Signatures[2]).ID, (*region1.Buckets)[0].Files[8].Chunks[2].Signature.ID)
			is.NotNil((*region1.Buckets)[0].Files[8].Chunks[3].Signature)
			is.Equal((*fixtures.Signatures[3]).ID, (*region1.Buckets)[0].Files[8].Chunks[3].Signature.ID)

			region2 := fixtures.Regions[1]
			err = sqlxx.Preload(ctx, driver, region2,
				"Buckets.Files",
				"Buckets.Files.Owner",
				"Buckets.Files.Owner.Profile",
				"Buckets.Files.Owner.Profile.Avatar",
				"Buckets.Files.Chunks",
				"Buckets.Files.Chunks.Mode",
				"Buckets.Files.Chunks.Signature",
			)
			is.NoError(err)
			is.NotNil(region2.Buckets)
			is.Empty((*region2.Buckets))
			is.Len((*region2.Buckets), 0)

			region3 := fixtures.Regions[2]
			err = sqlxx.Preload(ctx, driver, region3,
				"Buckets.Files",
				"Buckets.Files.Owner",
				"Buckets.Files.Owner.Profile",
				"Buckets.Files.Owner.Profile.Avatar",
				"Buckets.Files.Chunks",
				"Buckets.Files.Chunks.Mode",
				"Buckets.Files.Chunks.Signature",
			)
			is.NoError(err)
			is.NotNil(region3.Buckets)
			is.NotEmpty((*region3.Buckets))
			is.Len((*region3.Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*region3.Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*region3.Buckets)[1].ID)
			is.Len((*region3.Buckets)[0].Files, 10)
			is.Len((*region3.Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*region3.Buckets)[0].Files[0].ID)
			is.NotNil((*region3.Buckets)[0].Files[0].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[0].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[0].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[0].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[0].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[0].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[0].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[0].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[0].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[0].Chunks)
			is.Len((*region3.Buckets)[0].Files[0].Chunks, 1)
			is.Equal((*fixtures.Chunks[28]).Hash, (*region3.Buckets)[0].Files[0].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[0].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[0].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[0].Chunks[0].Signature)
			is.Equal((*fixtures.Files[10]).ID, (*region3.Buckets)[0].Files[1].ID)
			is.NotNil((*region3.Buckets)[0].Files[1].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[1].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[1].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[1].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[1].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[1].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[1].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[1].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[1].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[1].Chunks)
			is.Len((*region3.Buckets)[0].Files[1].Chunks, 1)
			is.Equal((*fixtures.Chunks[29]).Hash, (*region3.Buckets)[0].Files[1].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[1].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[1].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[1].Chunks[0].Signature)
			is.Equal((*fixtures.Files[11]).ID, (*region3.Buckets)[0].Files[2].ID)
			is.NotNil((*region3.Buckets)[0].Files[2].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[2].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*region3.Buckets)[0].Files[2].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[2].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[2].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*region3.Buckets)[0].Files[2].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[2].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[2].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*region3.Buckets)[0].Files[2].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[2].Chunks)
			is.Len((*region3.Buckets)[0].Files[2].Chunks, 2)
			is.Equal((*fixtures.Chunks[30]).Hash, (*region3.Buckets)[0].Files[2].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[2].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[2].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[31]).Hash, (*region3.Buckets)[0].Files[2].Chunks[1].Hash)
			is.NotNil((*region3.Buckets)[0].Files[2].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[2].Chunks[1].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[2].Chunks[0].Signature)
			is.Nil((*region3.Buckets)[0].Files[2].Chunks[1].Signature)
			is.Equal((*fixtures.Files[12]).ID, (*region3.Buckets)[0].Files[3].ID)
			is.NotNil((*region3.Buckets)[0].Files[3].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[3].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*region3.Buckets)[0].Files[3].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[3].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[3].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*region3.Buckets)[0].Files[3].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[3].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[3].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*region3.Buckets)[0].Files[3].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[3].Chunks)
			is.Len((*region3.Buckets)[0].Files[3].Chunks, 2)
			is.Equal((*fixtures.Chunks[32]).Hash, (*region3.Buckets)[0].Files[3].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[3].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[3].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[33]).Hash, (*region3.Buckets)[0].Files[3].Chunks[1].Hash)
			is.NotNil((*region3.Buckets)[0].Files[3].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[3].Chunks[1].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[3].Chunks[0].Signature)
			is.Nil((*region3.Buckets)[0].Files[3].Chunks[1].Signature)
			is.Equal((*fixtures.Files[13]).ID, (*region3.Buckets)[0].Files[4].ID)
			is.NotNil((*region3.Buckets)[0].Files[4].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[4].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*region3.Buckets)[0].Files[4].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[4].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[4].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*region3.Buckets)[0].Files[4].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[4].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[4].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*region3.Buckets)[0].Files[4].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[4].Chunks)
			is.Len((*region3.Buckets)[0].Files[4].Chunks, 2)
			is.Equal((*fixtures.Chunks[34]).Hash, (*region3.Buckets)[0].Files[4].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[4].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[4].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[35]).Hash, (*region3.Buckets)[0].Files[4].Chunks[1].Hash)
			is.NotNil((*region3.Buckets)[0].Files[4].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[4].Chunks[1].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[4].Chunks[0].Signature)
			is.Nil((*region3.Buckets)[0].Files[4].Chunks[1].Signature)
			is.Equal((*fixtures.Files[14]).ID, (*region3.Buckets)[0].Files[5].ID)
			is.NotNil((*region3.Buckets)[0].Files[5].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[5].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*region3.Buckets)[0].Files[5].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[5].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[5].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*region3.Buckets)[0].Files[5].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[5].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[5].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*region3.Buckets)[0].Files[5].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[5].Chunks)
			is.Len((*region3.Buckets)[0].Files[5].Chunks, 2)
			is.Equal((*fixtures.Chunks[36]).Hash, (*region3.Buckets)[0].Files[5].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[5].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[5].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[37]).Hash, (*region3.Buckets)[0].Files[5].Chunks[1].Hash)
			is.NotNil((*region3.Buckets)[0].Files[5].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*region3.Buckets)[0].Files[5].Chunks[1].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[5].Chunks[0].Signature)
			is.Nil((*region3.Buckets)[0].Files[5].Chunks[1].Signature)
			is.Equal((*fixtures.Files[15]).ID, (*region3.Buckets)[0].Files[6].ID)
			is.NotNil((*region3.Buckets)[0].Files[6].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[6].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[6].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[6].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[6].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[6].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[6].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[6].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[6].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[6].Chunks)
			is.Len((*region3.Buckets)[0].Files[6].Chunks, 1)
			is.Equal((*fixtures.Chunks[38]).Hash, (*region3.Buckets)[0].Files[6].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[6].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[6].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[6].Chunks[0].Signature)
			is.Equal((*fixtures.Files[16]).ID, (*region3.Buckets)[0].Files[7].ID)
			is.NotNil((*region3.Buckets)[0].Files[7].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[7].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[7].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[7].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[7].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[7].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[7].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[7].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[7].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[7].Chunks)
			is.Len((*region3.Buckets)[0].Files[7].Chunks, 1)
			is.Equal((*fixtures.Chunks[39]).Hash, (*region3.Buckets)[0].Files[7].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[7].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[7].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[7].Chunks[0].Signature)
			is.Equal((*fixtures.Files[17]).ID, (*region3.Buckets)[0].Files[8].ID)
			is.NotNil((*region3.Buckets)[0].Files[8].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[8].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[8].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[8].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[8].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[8].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[8].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[8].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[8].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[8].Chunks)
			is.Len((*region3.Buckets)[0].Files[8].Chunks, 1)
			is.Equal((*fixtures.Chunks[40]).Hash, (*region3.Buckets)[0].Files[8].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[8].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[8].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[8].Chunks[0].Signature)
			is.Equal((*fixtures.Files[18]).ID, (*region3.Buckets)[0].Files[9].ID)
			is.NotNil((*region3.Buckets)[0].Files[9].Owner)
			is.NotEmpty(*((*region3.Buckets)[0].Files[9].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*region3.Buckets)[0].Files[9].Owner.ID)
			is.NotNil((*region3.Buckets)[0].Files[9].Owner.Profile)
			is.NotEmpty(*((*region3.Buckets)[0].Files[9].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*region3.Buckets)[0].Files[9].Owner.Profile.ID)
			is.NotNil((*region3.Buckets)[0].Files[9].Owner.Profile.Avatar)
			is.NotEmpty(*((*region3.Buckets)[0].Files[9].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*region3.Buckets)[0].Files[9].Owner.Profile.Avatar.ID)
			is.NotEmpty((*region3.Buckets)[0].Files[9].Chunks)
			is.Len((*region3.Buckets)[0].Files[9].Chunks, 1)
			is.Equal((*fixtures.Chunks[41]).Hash, (*region3.Buckets)[0].Files[9].Chunks[0].Hash)
			is.NotNil((*region3.Buckets)[0].Files[9].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*region3.Buckets)[0].Files[9].Chunks[0].Mode.ID)
			is.Nil((*region3.Buckets)[0].Files[9].Chunks[0].Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := []ExoRegion{
				*fixtures.Regions[0],
				*fixtures.Regions[1],
				*fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions,
				"Buckets.Files",
				"Buckets.Files.Owner",
				"Buckets.Files.Owner.Profile",
				"Buckets.Files.Owner.Profile.Avatar",
				"Buckets.Files.Chunks",
				"Buckets.Files.Chunks.Mode",
				"Buckets.Files.Chunks.Signature",
			)
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*regions[0].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*regions[0].Buckets)[1].ID)
			is.Len((*regions[0].Buckets)[0].Files, 12)
			is.Len((*regions[0].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Files[0].ID)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Owner)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[0].Owner))
			is.Equal((*fixtures.Users[13]).ID, (*regions[0].Buckets)[0].Files[0].Owner.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Owner.Profile)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[0].Owner.Profile))
			is.Equal((*fixtures.Profiles[13]).ID, (*regions[0].Buckets)[0].Files[0].Owner.Profile.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[0].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[13]).ID, (*regions[0].Buckets)[0].Files[0].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[0].Buckets)[0].Files[0].Chunks)
			is.Len((*regions[0].Buckets)[0].Files[0].Chunks, 3)
			is.Equal((*fixtures.Chunks[0]).Hash, (*regions[0].Buckets)[0].Files[0].Chunks[0].Hash)
			is.Equal((*fixtures.Chunks[1]).Hash, (*regions[0].Buckets)[0].Files[0].Chunks[1].Hash)
			is.Equal((*fixtures.Chunks[2]).Hash, (*regions[0].Buckets)[0].Files[0].Chunks[2].Hash)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*regions[0].Buckets)[0].Files[0].Chunks[0].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*regions[0].Buckets)[0].Files[0].Chunks[1].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[0].Chunks[2].Mode)
			is.Equal((*fixtures.Modes[3]).ID, (*regions[0].Buckets)[0].Files[0].Chunks[2].Mode.ID)
			is.Nil((*regions[0].Buckets)[0].Files[0].Chunks[0].Signature)
			is.Nil((*regions[0].Buckets)[0].Files[0].Chunks[1].Signature)
			is.Nil((*regions[0].Buckets)[0].Files[0].Chunks[2].Signature)
			is.Equal((*fixtures.Files[8]).ID, (*regions[0].Buckets)[0].Files[8].ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Owner)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[8].Owner))
			is.Equal((*fixtures.Users[18]).ID, (*regions[0].Buckets)[0].Files[8].Owner.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Owner.Profile)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[8].Owner.Profile))
			is.Equal((*fixtures.Profiles[18]).ID, (*regions[0].Buckets)[0].Files[8].Owner.Profile.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[0].Buckets)[0].Files[8].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[18]).ID, (*regions[0].Buckets)[0].Files[8].Owner.Profile.Avatar.ID)
			is.Len((*regions[0].Buckets)[0].Files[8].Chunks, 4)
			is.Equal((*fixtures.Chunks[24]).Hash, (*regions[0].Buckets)[0].Files[8].Chunks[0].Hash)
			is.Equal((*fixtures.Chunks[25]).Hash, (*regions[0].Buckets)[0].Files[8].Chunks[1].Hash)
			is.Equal((*fixtures.Chunks[26]).Hash, (*regions[0].Buckets)[0].Files[8].Chunks[2].Hash)
			is.Equal((*fixtures.Chunks[27]).Hash, (*regions[0].Buckets)[0].Files[8].Chunks[3].Hash)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[0].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[1].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[2].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[2].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[3].Mode)
			is.Equal((*fixtures.Modes[0]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[3].Mode.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[0].Signature)
			is.Equal((*fixtures.Signatures[0]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[0].Signature.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[1].Signature)
			is.Equal((*fixtures.Signatures[1]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[1].Signature.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[2].Signature)
			is.Equal((*fixtures.Signatures[2]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[2].Signature.ID)
			is.NotNil((*regions[0].Buckets)[0].Files[8].Chunks[3].Signature)
			is.Equal((*fixtures.Signatures[3]).ID, (*regions[0].Buckets)[0].Files[8].Chunks[3].Signature.ID)

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*regions[2].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*regions[2].Buckets)[1].ID)
			is.Len((*regions[2].Buckets)[0].Files, 10)
			is.Len((*regions[2].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*regions[2].Buckets)[0].Files[0].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[0].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[0].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[0].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[0].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[0].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[0].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[0].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[0].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[0].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[0].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[0].Chunks, 1)
			is.Equal((*fixtures.Chunks[28]).Hash, (*regions[2].Buckets)[0].Files[0].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[0].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[0].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[0].Chunks[0].Signature)
			is.Equal((*fixtures.Files[10]).ID, (*regions[2].Buckets)[0].Files[1].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[1].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[1].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[1].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[1].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[1].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[1].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[1].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[1].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[1].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[1].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[1].Chunks, 1)
			is.Equal((*fixtures.Chunks[29]).Hash, (*regions[2].Buckets)[0].Files[1].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[1].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[1].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[1].Chunks[0].Signature)
			is.Equal((*fixtures.Files[11]).ID, (*regions[2].Buckets)[0].Files[2].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[2].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[2].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*regions[2].Buckets)[0].Files[2].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[2].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[2].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*regions[2].Buckets)[0].Files[2].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[2].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[2].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*regions[2].Buckets)[0].Files[2].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[2].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[2].Chunks, 2)
			is.Equal((*fixtures.Chunks[30]).Hash, (*regions[2].Buckets)[0].Files[2].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[2].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[2].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[31]).Hash, (*regions[2].Buckets)[0].Files[2].Chunks[1].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[2].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[2].Chunks[1].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[2].Chunks[0].Signature)
			is.Nil((*regions[2].Buckets)[0].Files[2].Chunks[1].Signature)
			is.Equal((*fixtures.Files[12]).ID, (*regions[2].Buckets)[0].Files[3].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[3].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[3].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*regions[2].Buckets)[0].Files[3].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[3].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[3].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*regions[2].Buckets)[0].Files[3].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[3].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[3].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*regions[2].Buckets)[0].Files[3].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[3].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[3].Chunks, 2)
			is.Equal((*fixtures.Chunks[32]).Hash, (*regions[2].Buckets)[0].Files[3].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[3].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[3].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[33]).Hash, (*regions[2].Buckets)[0].Files[3].Chunks[1].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[3].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[3].Chunks[1].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[3].Chunks[0].Signature)
			is.Nil((*regions[2].Buckets)[0].Files[3].Chunks[1].Signature)
			is.Equal((*fixtures.Files[13]).ID, (*regions[2].Buckets)[0].Files[4].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[4].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[4].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*regions[2].Buckets)[0].Files[4].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[4].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[4].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*regions[2].Buckets)[0].Files[4].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[4].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[4].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*regions[2].Buckets)[0].Files[4].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[4].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[4].Chunks, 2)
			is.Equal((*fixtures.Chunks[34]).Hash, (*regions[2].Buckets)[0].Files[4].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[4].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[4].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[35]).Hash, (*regions[2].Buckets)[0].Files[4].Chunks[1].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[4].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[4].Chunks[1].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[4].Chunks[0].Signature)
			is.Nil((*regions[2].Buckets)[0].Files[4].Chunks[1].Signature)
			is.Equal((*fixtures.Files[14]).ID, (*regions[2].Buckets)[0].Files[5].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[5].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[5].Owner))
			is.Equal((*fixtures.Users[5]).ID, (*regions[2].Buckets)[0].Files[5].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[5].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[5].Owner.Profile))
			is.Equal((*fixtures.Profiles[5]).ID, (*regions[2].Buckets)[0].Files[5].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[5].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[5].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[5]).ID, (*regions[2].Buckets)[0].Files[5].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[5].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[5].Chunks, 2)
			is.Equal((*fixtures.Chunks[36]).Hash, (*regions[2].Buckets)[0].Files[5].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[5].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[5].Chunks[0].Mode.ID)
			is.Equal((*fixtures.Chunks[37]).Hash, (*regions[2].Buckets)[0].Files[5].Chunks[1].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[5].Chunks[1].Mode)
			is.Equal((*fixtures.Modes[1]).ID, (*regions[2].Buckets)[0].Files[5].Chunks[1].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[5].Chunks[0].Signature)
			is.Nil((*regions[2].Buckets)[0].Files[5].Chunks[1].Signature)
			is.Equal((*fixtures.Files[15]).ID, (*regions[2].Buckets)[0].Files[6].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[6].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[6].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[6].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[6].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[6].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[6].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[6].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[6].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[6].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[6].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[6].Chunks, 1)
			is.Equal((*fixtures.Chunks[38]).Hash, (*regions[2].Buckets)[0].Files[6].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[6].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[6].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[6].Chunks[0].Signature)
			is.Equal((*fixtures.Files[16]).ID, (*regions[2].Buckets)[0].Files[7].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[7].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[7].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[7].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[7].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[7].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[7].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[7].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[7].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[7].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[7].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[7].Chunks, 1)
			is.Equal((*fixtures.Chunks[39]).Hash, (*regions[2].Buckets)[0].Files[7].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[7].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[7].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[7].Chunks[0].Signature)
			is.Equal((*fixtures.Files[17]).ID, (*regions[2].Buckets)[0].Files[8].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[8].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[8].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[8].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[8].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[8].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[8].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[8].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[8].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[8].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[8].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[8].Chunks, 1)
			is.Equal((*fixtures.Chunks[40]).Hash, (*regions[2].Buckets)[0].Files[8].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[8].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[8].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[8].Chunks[0].Signature)
			is.Equal((*fixtures.Files[18]).ID, (*regions[2].Buckets)[0].Files[9].ID)
			is.NotNil((*regions[2].Buckets)[0].Files[9].Owner)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[9].Owner))
			is.Equal((*fixtures.Users[7]).ID, (*regions[2].Buckets)[0].Files[9].Owner.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[9].Owner.Profile)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[9].Owner.Profile))
			is.Equal((*fixtures.Profiles[7]).ID, (*regions[2].Buckets)[0].Files[9].Owner.Profile.ID)
			is.NotNil((*regions[2].Buckets)[0].Files[9].Owner.Profile.Avatar)
			is.NotEmpty(*((*regions[2].Buckets)[0].Files[9].Owner.Profile.Avatar))
			is.Equal((*fixtures.Avatars[7]).ID, (*regions[2].Buckets)[0].Files[9].Owner.Profile.Avatar.ID)
			is.NotEmpty((*regions[2].Buckets)[0].Files[9].Chunks)
			is.Len((*regions[2].Buckets)[0].Files[9].Chunks, 1)
			is.Equal((*fixtures.Chunks[41]).Hash, (*regions[2].Buckets)[0].Files[9].Chunks[0].Hash)
			is.NotNil((*regions[2].Buckets)[0].Files[9].Chunks[0].Mode)
			is.Equal((*fixtures.Modes[2]).ID, (*regions[2].Buckets)[0].Files[9].Chunks[0].Mode.ID)
			is.Nil((*regions[2].Buckets)[0].Files[9].Chunks[0].Signature)

		}
		{

			fixtures := GenerateExoCloudFixtures(ctx, driver, is)
			CheckExoCloudFixtures(fixtures)

			regions := []ExoRegion{
				*fixtures.Regions[0],
				*fixtures.Regions[1],
				*fixtures.Regions[2],
			}

			err := sqlxx.Preload(ctx, driver, &regions,
				"Buckets.Directories",
				"Buckets.Directories.Directories",
				"Buckets.Directories.Directories.Directories",
				"Buckets.Files",
				"Buckets.Directories.Files",
				"Buckets.Directories.Directories.Files",
				"Buckets.Directories.Directories.Directories.Files",
			)
			is.NoError(err)
			is.Len(regions, 3)
			is.Equal(fixtures.Regions[0].ID, regions[0].ID)
			is.Equal(fixtures.Regions[1].ID, regions[1].ID)
			is.Equal(fixtures.Regions[2].ID, regions[2].ID)

			is.NotNil(regions[0].Buckets)
			is.NotEmpty((*regions[0].Buckets))
			is.Len((*regions[0].Buckets), 2)
			is.Equal((*fixtures.Buckets[0]).ID, (*regions[0].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[1]).ID, (*regions[0].Buckets)[1].ID)
			is.Len((*regions[0].Buckets)[0].Directories, 16)
			is.Len((*regions[0].Buckets)[1].Directories, 0)
			SortExoCloudDirectories1(fixtures.Directories, &(*regions[0].Buckets)[0].Directories)
			is.Equal((*fixtures.Directories[0]).ID, (*regions[0].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[1]).ID, (*regions[0].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[2]).ID, (*regions[0].Buckets)[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[3]).ID, (*regions[0].Buckets)[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[4]).ID, (*regions[0].Buckets)[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[5]).ID, (*regions[0].Buckets)[0].Directories[5].ID)
			is.Equal((*fixtures.Directories[6]).ID, (*regions[0].Buckets)[0].Directories[6].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*regions[0].Buckets)[0].Directories[7].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*regions[0].Buckets)[0].Directories[8].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*regions[0].Buckets)[0].Directories[9].ID)
			is.Equal((*fixtures.Directories[10]).ID, (*regions[0].Buckets)[0].Directories[10].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*regions[0].Buckets)[0].Directories[11].ID)
			is.Equal((*fixtures.Directories[12]).ID, (*regions[0].Buckets)[0].Directories[12].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*regions[0].Buckets)[0].Directories[13].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*regions[0].Buckets)[0].Directories[14].ID)
			is.Equal((*fixtures.Directories[15]).ID, (*regions[0].Buckets)[0].Directories[15].ID)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories, 6)
			SortExoCloudDirectories2(fixtures.Directories, &(*regions[0].Buckets)[0].Directories[0].Directories)
			is.Equal((*fixtures.Directories[6]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[7]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[8]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[2].ID)
			is.Equal((*fixtures.Directories[9]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[3].ID)
			is.Equal((*fixtures.Directories[10]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[4].ID)
			is.Equal((*fixtures.Directories[11]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[5].ID)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[0].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[1].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[2].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[3].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[4].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[5].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[1].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[2].Directories, 3)
			SortExoCloudDirectories2(fixtures.Directories, &(*regions[0].Buckets)[0].Directories[2].Directories)
			is.Equal((*fixtures.Directories[12]).ID, (*regions[0].Buckets)[0].Directories[2].Directories[0].ID)
			is.Equal((*fixtures.Directories[13]).ID, (*regions[0].Buckets)[0].Directories[2].Directories[1].ID)
			is.Equal((*fixtures.Directories[14]).ID, (*regions[0].Buckets)[0].Directories[2].Directories[2].ID)
			is.Len((*regions[0].Buckets)[0].Directories[2].Directories[0].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[2].Directories[1].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[2].Directories[2].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[3].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[4].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[5].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[6].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[7].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[8].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[9].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[10].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[11].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[12].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[13].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[14].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Directories[15].Directories, 0)
			is.Len((*regions[0].Buckets)[0].Files, 12)
			is.Len((*regions[0].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*regions[0].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[2]).ID, (*regions[0].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[3]).ID, (*regions[0].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[4]).ID, (*regions[0].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[5]).ID, (*regions[0].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[6]).ID, (*regions[0].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[7]).ID, (*regions[0].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[8]).ID, (*regions[0].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[19]).ID, (*regions[0].Buckets)[0].Files[9].ID)
			is.Equal((*fixtures.Files[20]).ID, (*regions[0].Buckets)[0].Files[10].ID)
			is.Equal((*fixtures.Files[21]).ID, (*regions[0].Buckets)[0].Files[11].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[0].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[1].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[2].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[3].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[4].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[5].Files)
			is.Empty((*regions[0].Buckets)[0].Directories[6].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[7].Files)
			is.Len((*regions[0].Buckets)[0].Directories[7].Files, 1)
			is.Equal((*fixtures.Files[7]).ID, (*regions[0].Buckets)[0].Directories[7].Files[0].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[8].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[9].Files)
			is.Len((*regions[0].Buckets)[0].Directories[9].Files, 1)
			is.Equal((*fixtures.Files[2]).ID, (*regions[0].Buckets)[0].Directories[9].Files[0].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[10].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[11].Files)
			is.Len((*regions[0].Buckets)[0].Directories[11].Files, 4)
			is.Equal((*fixtures.Files[3]).ID, (*regions[0].Buckets)[0].Directories[11].Files[0].ID)
			is.Equal((*fixtures.Files[4]).ID, (*regions[0].Buckets)[0].Directories[11].Files[1].ID)
			is.Equal((*fixtures.Files[5]).ID, (*regions[0].Buckets)[0].Directories[11].Files[2].ID)
			is.Equal((*fixtures.Files[6]).ID, (*regions[0].Buckets)[0].Directories[11].Files[3].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[12].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[13].Files)
			is.Len((*regions[0].Buckets)[0].Directories[13].Files, 2)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Directories[13].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*regions[0].Buckets)[0].Directories[13].Files[1].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[14].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[15].Files)
			is.Len((*regions[0].Buckets)[0].Directories[15].Files, 4)
			is.Equal((*fixtures.Files[8]).ID, (*regions[0].Buckets)[0].Directories[15].Files[0].ID)
			is.Equal((*fixtures.Files[19]).ID, (*regions[0].Buckets)[0].Directories[15].Files[1].ID)
			is.Equal((*fixtures.Files[20]).ID, (*regions[0].Buckets)[0].Directories[15].Files[2].ID)
			is.Equal((*fixtures.Files[21]).ID, (*regions[0].Buckets)[0].Directories[15].Files[3].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[0].Directories[0].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[0].Directories[1].Files)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[1].Files, 1)
			is.Equal((*fixtures.Files[7]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[1].Files[0].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[0].Directories[2].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[0].Directories[3].Files)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[3].Files, 1)
			is.Equal((*fixtures.Files[2]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[3].Files[0].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[0].Directories[4].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[0].Directories[5].Files)
			is.Len((*regions[0].Buckets)[0].Directories[0].Directories[5].Files, 4)
			is.Equal((*fixtures.Files[3]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[5].Files[0].ID)
			is.Equal((*fixtures.Files[4]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[5].Files[1].ID)
			is.Equal((*fixtures.Files[5]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[5].Files[2].ID)
			is.Equal((*fixtures.Files[6]).ID, (*regions[0].Buckets)[0].Directories[0].Directories[5].Files[3].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[2].Directories[0].Files)
			is.NotEmpty((*regions[0].Buckets)[0].Directories[2].Directories[1].Files)
			is.Len((*regions[0].Buckets)[0].Directories[2].Directories[1].Files, 2)
			is.Equal((*fixtures.Files[0]).ID, (*regions[0].Buckets)[0].Directories[2].Directories[1].Files[0].ID)
			is.Equal((*fixtures.Files[1]).ID, (*regions[0].Buckets)[0].Directories[2].Directories[1].Files[1].ID)
			is.Empty((*regions[0].Buckets)[0].Directories[2].Directories[2].Files)

			is.NotNil(regions[1].Buckets)
			is.Empty((*regions[1].Buckets))
			is.Len((*regions[1].Buckets), 0)

			is.NotNil(regions[2].Buckets)
			is.NotEmpty((*regions[2].Buckets))
			is.Len((*regions[2].Buckets), 2)
			is.Equal((*fixtures.Buckets[2]).ID, (*regions[2].Buckets)[0].ID)
			is.Equal((*fixtures.Buckets[3]).ID, (*regions[2].Buckets)[1].ID)
			is.Len((*regions[2].Buckets)[0].Directories, 3)
			is.Len((*regions[2].Buckets)[1].Directories, 0)
			SortExoCloudDirectories1(fixtures.Directories, &(*regions[2].Buckets)[0].Directories)
			is.Equal((*fixtures.Directories[16]).ID, (*regions[2].Buckets)[0].Directories[0].ID)
			is.Equal((*fixtures.Directories[17]).ID, (*regions[2].Buckets)[0].Directories[1].ID)
			is.Equal((*fixtures.Directories[18]).ID, (*regions[2].Buckets)[0].Directories[2].ID)
			is.Len((*regions[2].Buckets)[0].Directories[0].Directories, 0)
			is.Len((*regions[2].Buckets)[0].Directories[1].Directories, 0)
			is.Len((*regions[2].Buckets)[0].Directories[2].Directories, 0)
			is.Len((*regions[2].Buckets)[0].Files, 10)
			is.Len((*regions[2].Buckets)[1].Files, 0)
			is.Equal((*fixtures.Files[9]).ID, (*regions[2].Buckets)[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*regions[2].Buckets)[0].Files[1].ID)
			is.Equal((*fixtures.Files[11]).ID, (*regions[2].Buckets)[0].Files[2].ID)
			is.Equal((*fixtures.Files[12]).ID, (*regions[2].Buckets)[0].Files[3].ID)
			is.Equal((*fixtures.Files[13]).ID, (*regions[2].Buckets)[0].Files[4].ID)
			is.Equal((*fixtures.Files[14]).ID, (*regions[2].Buckets)[0].Files[5].ID)
			is.Equal((*fixtures.Files[15]).ID, (*regions[2].Buckets)[0].Files[6].ID)
			is.Equal((*fixtures.Files[16]).ID, (*regions[2].Buckets)[0].Files[7].ID)
			is.Equal((*fixtures.Files[17]).ID, (*regions[2].Buckets)[0].Files[8].ID)
			is.Equal((*fixtures.Files[18]).ID, (*regions[2].Buckets)[0].Files[9].ID)
			is.NotEmpty((*regions[2].Buckets)[0].Directories[0].Files)
			is.Len((*regions[2].Buckets)[0].Directories[0].Files, 2)
			is.Equal((*fixtures.Files[9]).ID, (*regions[2].Buckets)[0].Directories[0].Files[0].ID)
			is.Equal((*fixtures.Files[10]).ID, (*regions[2].Buckets)[0].Directories[0].Files[1].ID)
			is.NotEmpty((*regions[2].Buckets)[0].Directories[1].Files)
			is.Len((*regions[2].Buckets)[0].Directories[1].Files, 4)
			is.Equal((*fixtures.Files[11]).ID, (*regions[2].Buckets)[0].Directories[1].Files[0].ID)
			is.Equal((*fixtures.Files[12]).ID, (*regions[2].Buckets)[0].Directories[1].Files[1].ID)
			is.Equal((*fixtures.Files[13]).ID, (*regions[2].Buckets)[0].Directories[1].Files[2].ID)
			is.Equal((*fixtures.Files[14]).ID, (*regions[2].Buckets)[0].Directories[1].Files[3].ID)
			is.NotEmpty((*regions[2].Buckets)[0].Directories[2].Files)
			is.Len((*regions[2].Buckets)[0].Directories[2].Files, 4)
			is.Equal((*fixtures.Files[15]).ID, (*regions[2].Buckets)[0].Directories[2].Files[0].ID)
			is.Equal((*fixtures.Files[16]).ID, (*regions[2].Buckets)[0].Directories[2].Files[1].ID)
			is.Equal((*fixtures.Files[17]).ID, (*regions[2].Buckets)[0].Directories[2].Files[2].ID)
			is.Equal((*fixtures.Files[18]).ID, (*regions[2].Buckets)[0].Directories[2].Files[3].ID)

		}
	})
}
