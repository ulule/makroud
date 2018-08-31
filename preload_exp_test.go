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

func TestPreload_ExoRegion_MultiLevel(t *testing.T) {
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

			slice := reflectx.NewReflectSlice(reflect.PtrTo(reflectx.GetIndirectType(ExoBucket{})))

			buckets, err := reflectx.GetFieldValue(region1, "Buckets")
			is.NoError(err)
			is.True(reflectx.IsSlice(buckets))

			val1 := reflectx.GetIndirectValue(buckets)

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

			{
				val1, err := reflectx.GetFieldValue(region1, "Buckets")
				is.NoError(err)
				is.NotNil(val1)

				val3 := reflectx.GetIndirectValue(val1)
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

					val2, err := reflectx.GetFieldValue(val4, "Directories")
					is.NoError(err)
					if val2 != nil {
						fmt.Printf("::10 %+v\n", val2)

						if reflectx.IsSlice(val2) {

							val5 := reflectx.GetIndirectValue(val2)
							for i := 0; i < val5.Len(); i++ {
								val6 := val5.Index(i)

								if val6.Kind() == reflect.Interface {
									val6 = reflect.ValueOf(val6.Interface())
									if val6.IsNil() {
										continue
									}
								}

								if val6.Kind() != reflect.Ptr && val6.CanAddr() {
									val6 = val6.Addr()
								}

								err := sqlxx.Preload(ctx, driver, val6.Interface(), "Files")
								is.NoError(err)

								val7, err := reflectx.GetFieldValue(val6, "Files")
								is.NoError(err)
								fmt.Printf("::1 %+v\n", val7)

							}

						} else {

							val7, err := reflectx.GetFieldValue(val2, "Files")
							is.NoError(err)
							fmt.Printf("::1 %+v\n", val7)

						}

					}
				}

				fmt.Printf("::region %+v\n", region1)
				for _, bucket := range *region1.Buckets {
					fmt.Printf("::bucket %+v\n", bucket)
					for _, directory := range bucket.Directories {
						fmt.Printf("::directory %+v\n", directory)
						for _, file := range directory.Files {
							fmt.Printf("::file %+v\n", file)
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
						}
					}
				}

			}
		}
	})
}
