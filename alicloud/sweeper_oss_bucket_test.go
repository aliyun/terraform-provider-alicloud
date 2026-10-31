package alicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_oss_bucket", &resource.Sweeper{
		Name: "alicloud_oss_bucket",
		F:    testSweepOSSBuckets,
	})
}

func testSweepOSSBuckets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf-test-",
		"test-bucket-",
		"tf-oss-test-",
		"tf-object-test-",
		"test-acc-alicloud-",
		"tftest",
		"cri",
		"tf-example",
		"tf_example",
		"terraform-example",
	}

	var options []oss.Option
	options = append(options, oss.MaxKeys(1000))

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.ListBuckets(options...)
	})
	if err != nil {
		return fmt.Errorf("Error retrieving OSS buckets: %s", err)
	}
	resp, _ := raw.(oss.ListBucketsResult)
	for _, v := range resp.Buckets {
		name := v.Name
		if !strings.HasSuffix(v.Location, client.RegionId) {
			continue
		}
		skip := true
		if !sweepAll() {
			if strings.HasPrefix(strings.ToLower(name), "terraform-alicloud-provider") {
				continue
			}
			if strings.HasPrefix(strings.ToLower(name), "terraform-remote-backend") {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping OSS bucket: %s", name)
				continue
			}
		}
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.Bucket(name)
		})
		if err != nil {
			return fmt.Errorf("Error getting bucket (%s): %#v", name, err)
		}
		bucket, _ := raw.(*oss.Bucket)
		var initialOptions []oss.Option
		nextMarker := ""
		for {
			var options []oss.Option
			options = append(options, initialOptions...)
			if nextMarker != "" {
				options = append(options, oss.Marker(nextMarker))
			}
			raw, err := client.WithOssBucketByName(bucket.BucketName, func(bucket *oss.Bucket) (interface{}, error) {
				return bucket.ListObjects(options...)
			})
			if err != nil {
				continue
			}
			response, _ := raw.(oss.ListObjectsResult)

			if response.Objects == nil || len(response.Objects) < 1 {
				break
			}

			if len(response.Objects) < 1 {
				break
			}

			if len(response.Objects) > 0 {
				for _, o := range response.Objects {
					if err := bucket.DeleteObject(o.Key); err != nil {
						log.Printf("[ERROR] Failed to delete object (%s): %s.", o.Key, err)
					}
				}
			}

			nextMarker = response.NextMarker
			if nextMarker == "" {
				break
			}
		}

		log.Printf("[INFO] Deleting OSS bucket: %s", name)

		_, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.DeleteBucket(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OSS bucket (%s): %s", name, err)
		}
	}
	return nil
}
