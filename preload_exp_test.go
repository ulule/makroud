package sqlxx_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
	"github.com/ulule/sqlxx/reflectx"
)

func TestPreloadX_ExoUser_MultiLevel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		fixtures := GenerateExoCloudFixtures(ctx, driver, is)

		user1 := fixtures.Users[1]

		err := sqlxx.Preload(ctx, driver, user1,
			"Profile",
			"Profile.Avatar",
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
			err := sqlxx.Preload(ctx, driver, region1, "Buckets")
			is.NoError(err)
			is.NotNil(region1.Buckets)

			region2 := fixtures.Regions[1]
			err = sqlxx.Preload(ctx, driver, region2, "Buckets")
			is.NoError(err)
			is.NotNil(region2.Buckets)

			region3 := fixtures.Regions[2]
			err = sqlxx.Preload(ctx, driver, region3, "Buckets")
			is.NoError(err)
			is.NotNil(region3.Buckets)

			regions := []*ExoRegion{
				region1, region2, region3,
			}

			buckets, err := reflectx.GetFieldValue(region1, "Buckets")
			is.NoError(err)
			is.True(reflectx.IsSlice(buckets))

			val1 := reflectx.GetIndirectValue(buckets)

			slice := reflectx.NewReflectSlice(reflect.PtrTo(reflectx.GetIndirectType(reflectx.GetSliceType(val1))))

			for i := 0; i < val1.Len(); i++ {
				val2 := val1.Index(i)

				if val2.Kind() == reflect.Interface {
					val2 = reflect.ValueOf(val2.Interface())
					if val2.IsNil() {
						continue
					}
				}

				if val2.Kind() != reflect.Ptr && val2.CanAddr() {
					val2 = val2.Addr()
				}

				reflectx.AppendReflectSlice(slice, val2.Interface())
			}

			buckets, err = reflectx.GetFieldValue(region2, "Buckets")
			is.NoError(err)
			is.True(reflectx.IsSlice(buckets))

			val2 := reflectx.GetIndirectValue(buckets)
			for i := 0; i < val2.Len(); i++ {
				val3 := val2.Index(i)

				if val3.Kind() == reflect.Interface {
					val3 = reflect.ValueOf(val3.Interface())
					if val3.IsNil() {
						continue
					}
				}

				if val3.Kind() != reflect.Ptr && val3.CanAddr() {
					val3 = val3.Addr()
				}

				reflectx.AppendReflectSlice(slice, val3.Interface())
			}

			buckets, err = reflectx.GetFieldValue(region3, "Buckets")
			is.NoError(err)
			is.True(reflectx.IsSlice(buckets))

			val3 := reflectx.GetIndirectValue(buckets)
			for i := 0; i < val3.Len(); i++ {
				val4 := val3.Index(i)

				if val4.Kind() == reflect.Interface {
					val4 = reflect.ValueOf(val4.Interface())
					if val4.IsNil() {
						continue
					}
				}

				if val4.Kind() != reflect.Ptr && val4.CanAddr() {
					val4 = val4.Addr()
				}

				reflectx.AppendReflectSlice(slice, val4.Interface())
			}

			err = sqlxx.Preload(ctx, driver, slice.Interface(), "Directories")
			is.NoError(err)

			walker := reflectx.NewWalker(regions)
			err = walker.Find("Buckets.Directories", func(values interface{}) error {
				return sqlxx.Preload(ctx, driver, values, "Files")
			})
			is.NoError(err)

			for i := 0; i < 4; i++ {
				fmt.Println()
			}

			walker = reflectx.NewWalker(regions)
			err = walker.Find("Buckets.Directories.Files", func(values interface{}) error {
				return sqlxx.Preload(ctx, driver, values, "Chunks")
			})
			is.NoError(err)

			for i := 0; i < 4; i++ {
				fmt.Println()
			}

			fmt.Printf("::region %+v\n", region1)
			for _, bucket := range *region1.Buckets {
				fmt.Printf("::bucket %+v\n", bucket)
				for _, directory := range bucket.Directories {
					fmt.Printf("::directory %+v\n", directory)
					for _, file := range directory.Files {
						fmt.Printf("::file %+v\n", file)
						for _, chunk := range file.Chunks {
							fmt.Printf("::chunk %+v\n", chunk)
						}
					}
				}
			}

			fmt.Printf("::region %+v\n", region2)
			for _, bucket := range *region2.Buckets {
				fmt.Printf("::bucket %+v\n", bucket)
				for _, directory := range bucket.Directories {
					fmt.Printf("::directory %+v\n", directory)
					for _, file := range directory.Files {
						fmt.Printf("::file %+v\n", file)
						for _, chunk := range file.Chunks {
							fmt.Printf("::chunk %+v\n", chunk)
						}
					}
				}
			}

			fmt.Printf("::region %+v\n", region3)
			for _, bucket := range *region3.Buckets {
				fmt.Printf("::bucket %+v\n", bucket)
				for _, directory := range bucket.Directories {
					fmt.Printf("::directory %+v\n", directory)
					for _, file := range directory.Files {
						fmt.Printf("::file %+v\n", file)
						for _, chunk := range file.Chunks {
							fmt.Printf("::chunk %+v\n", chunk)
						}
					}
				}
			}
		}
	})
}
