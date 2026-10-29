package cloudphone_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudphone"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudPhone KeyPair. >>> Resource test cases, automatically generated.
// Case chuyuan_createKeyPair_import_prod 9935
func TestAccAliCloudCloudPhoneKeyPair_basic9935(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_phone_key_pair.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudPhoneKeyPairMap9935)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudphone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKeyPair")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudphone%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudPhoneKeyPairBasicDependence9935)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name":   name,
					"public_key_body": "QAAAAH0o+PMrbz9ZlxaNMYlk1rJkN4JXqwSUVYW5YzMW3fWJ7At1XO40GYDEFL43fLob52pmRxRDuRoGAELmS1AyzqUle2v9yGKFziqS/vK/4vM4MW/ppnTmvh9zPXir0fB/uwXS4iS6xt0gmvprgyRNs7hgXtBK9ASiGuPCv47aRJqh9mYzq2pe2rgb+K0OU5/nQXwWKSxYsv+w3KWPshpwx8iF/JWvjixILJ5gygndd+1HyE8jrLVmvm/OitNaMgkolY1bvmRVVKLmzde7FtXw0s4TVfYUvF385gwlrOulKcL7UuMHV87MV/tcvEA0Gg88JrKgI5LmvQ8BDkrfoSi+bchk1KTAqJ8YMvL2pOogXbBoONeJS176zLYpLHmONtIDQFz/gEqAjGQVW+j4J+1w8oWrn8EjtcDe2kY34s3PDLioK3BN9CIBBQur+SH25R0RnEqD0YPFT7/ym0LomtPOS0t72n5JejBTfWaXiqb/I4f2Ypy1PA6fV5UUFIHODpNtuS4g2HKKqDS/sgYRBA2gpN2MmqeqgsQmSy+EljHdUe4KDVqAZ/qxLqnbp47BGHw2xjuZ60nXAoRecWCg2GDbx13ga4dKUQY+ER8Jruz7ILK4MRB7E4SjSUVmgcdh534c51BYIdI2HkQwQU2dgyJyQme9sDQxxGHpYKFQlFSyXUeOSjXLtQEAAQA= qiaozhou_15694163938@h2sqyfpc71g1t2w",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name":   name,
						"public_key_body": "QAAAAH0o+PMrbz9ZlxaNMYlk1rJkN4JXqwSUVYW5YzMW3fWJ7At1XO40GYDEFL43fLob52pmRxRDuRoGAELmS1AyzqUle2v9yGKFziqS/vK/4vM4MW/ppnTmvh9zPXir0fB/uwXS4iS6xt0gmvprgyRNs7hgXtBK9ASiGuPCv47aRJqh9mYzq2pe2rgb+K0OU5/nQXwWKSxYsv+w3KWPshpwx8iF/JWvjixILJ5gygndd+1HyE8jrLVmvm/OitNaMgkolY1bvmRVVKLmzde7FtXw0s4TVfYUvF385gwlrOulKcL7UuMHV87MV/tcvEA0Gg88JrKgI5LmvQ8BDkrfoSi+bchk1KTAqJ8YMvL2pOogXbBoONeJS176zLYpLHmONtIDQFz/gEqAjGQVW+j4J+1w8oWrn8EjtcDe2kY34s3PDLioK3BN9CIBBQur+SH25R0RnEqD0YPFT7/ym0LomtPOS0t72n5JejBTfWaXiqb/I4f2Ypy1PA6fV5UUFIHODpNtuS4g2HKKqDS/sgYRBA2gpN2MmqeqgsQmSy+EljHdUe4KDVqAZ/qxLqnbp47BGHw2xjuZ60nXAoRecWCg2GDbx13ga4dKUQY+ER8Jruz7ILK4MRB7E4SjSUVmgcdh534c51BYIdI2HkQwQU2dgyJyQme9sDQxxGHpYKFQlFSyXUeOSjXLtQEAAQA= qiaozhou_15694163938@h2sqyfpc71g1t2w",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_key_body"},
			},
		},
	})
}

var AlicloudCloudPhoneKeyPairMap9935 = map[string]string{}

func AlicloudCloudPhoneKeyPairBasicDependence9935(name string) string {
	// key_pair is a standalone account-level resource (CreateKeyPair/ImportKeyPair);
	// the generated policy/instance_group/instance fixture is unnecessary and caused
	// policyUndeleteInstance destroy failures (residual instances bound to the policy).
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case chuyuan_createKeyPair_create_prod 10198
func TestAccAliCloudCloudPhoneKeyPair_basic10198(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_phone_key_pair.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudPhoneKeyPairMap10198)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudphone.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeKeyPair")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudphone%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudPhoneKeyPairBasicDependence10198)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_pair_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_pair_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_key_body"},
			},
		},
	})
}

var AlicloudCloudPhoneKeyPairMap10198 = map[string]string{}

func AlicloudCloudPhoneKeyPairBasicDependence10198(name string) string {
	// key_pair is a standalone account-level resource (CreateKeyPair/ImportKeyPair);
	// the generated policy/instance_group/instance fixture is unnecessary and caused
	// policyUndeleteInstance destroy failures (residual instances bound to the policy).
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudPhone KeyPair. <<< Resource test cases, automatically generated.
