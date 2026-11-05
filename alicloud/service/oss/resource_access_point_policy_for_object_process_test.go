package oss_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The Object FC access point that the policy attaches to is not yet a Terraform
// resource, but a regular OSS access point is. The test therefore creates the
// prerequisite chain (bucket + regular access point) as Terraform resources,
// then creates the Object FC access point itself through the raw OSS API in a
// check, polls it to ready, and finally exercises the policy resource lifecycle
// on top of it. Per the OSS API contract the function configuration of an
// Object FC access point is optional, so the access point is created with the
// minimal body (just the referenced regular access point name); the policy
// under test does not invoke the function and only needs the access point to
// exist.

const (
	accessPointPolicyForObjectProcessAction  = "/?accessPointPolicyForObjectProcess"
	accessPointForObjectProcessAction        = "/?accessPointForObjectProcess"
	accessPointPolicyForObjectProcessProduct = "Oss"
	accessPointPolicyForObjectProcessVersion = "2019-05-17"
	// The Object FC access point creation takes roughly ten minutes per the OSS
	// user guide; bound the polling well beyond that.
	accessPointForObjectProcessStatusTimeout = 20 * time.Minute
	accessPointForObjectProcessDeleteTimeout = 10 * time.Minute
)

func TestAccAliCloudOssAccessPointPolicyForObjectProcess_basic(t *testing.T) {
	resourceId := "alicloud_oss_access_point_policy_for_object_process.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	name := fmt.Sprintf("tf-testacc%sossappfop%d", region, rand)
	// The Object FC access point name is a separate identifier (carried by the
	// x-oss-access-point-for-object-process-name header); it must be lowercase,
	// digits and hyphens only.
	objectFcApName := name + "-fcap"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", region)
			// Idempotent pre-clean of any leftover Object FC access point and
			// policy from a prior aborted run.
			tearDownObjectFcAccessPoint(name, objectFcApName)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointPolicyForObjectProcessDestroy(name, objectFcApName),
		Steps: []resource.TestStep{
			{
				// Prerequisites only: the bucket and the regular access point must
				// exist before the Object FC access point can be created.
				Config: testAccAliCloudOssAccessPointPolicyForObjectProcessPrereq(name, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCreateObjectFcAccessPoint(name, objectFcApName),
				),
			},
			{
				Config: testAccAliCloudOssAccessPointPolicyForObjectProcessConfig(name, region, objectFcApName, testAccAliCloudOssAccessPointPolicyForObjectProcessDocument("Allow")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "access_point_name", objectFcApName),
					resource.TestCheckResourceAttr(resourceId, "id", fmt.Sprintf("%s:%s", name, objectFcApName)),
					testAccCheckOssAccessPointPolicyForObjectProcessEffect(resourceId, "Allow"),
				),
			},
			{
				// Update: replace the policy document with a non-canonical but
				// equivalent encoding and flip the Effect; the API re-encodes on
				// read, so the semantic-equivalence modifier must keep it clean.
				Config: testAccAliCloudOssAccessPointPolicyForObjectProcessFormatted(name, region, objectFcApName, "Deny"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyForObjectProcessEffect(resourceId, "Deny"),
				),
			},
			{
				Config: testAccAliCloudOssAccessPointPolicyForObjectProcessConfig(name, region, objectFcApName, testAccAliCloudOssAccessPointPolicyForObjectProcessDocument("Allow")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssAccessPointPolicyForObjectProcessEffect(resourceId, "Allow"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				// Destroy: tear down the out-of-band Object FC access point first
				// (and its policy) so Terraform can delete the regular access point
				// and bucket, then let Terraform clean up the rest.
				PreConfig: func() {
					tearDownObjectFcAccessPoint(name, objectFcApName)
				},
				Config:  testAccAliCloudOssAccessPointPolicyForObjectProcessPrereq(name, region),
				Destroy: true,
			},
		},
	})
}

// testAccCreateObjectFcAccessPoint creates the Object FC access point via the
// raw OSS API, referencing the already-applied regular access point, and polls
// it until it reaches the ready state.
func testAccCreateObjectFcAccessPoint(bucket, objectFcApName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// The regular access point name equals the bucket name in this test's
		// config; verify it is present in state before creating the Object FC
		// access point that references it.
		root := s.RootModule()
		bucketRs, ok := root.Resources["alicloud_oss_bucket.default"]
		if !ok {
			return fmt.Errorf("alicloud_oss_bucket.default not found in state")
		}
		apRs, ok := root.Resources["alicloud_oss_access_point.default"]
		if !ok {
			return fmt.Errorf("alicloud_oss_access_point.default not found in state")
		}
		regularAccessPointName := apRs.Primary.Attributes["access_point_name"]
		if regularAccessPointName == "" {
			regularAccessPointName = bucketRs.Primary.ID
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		if err := createObjectFcAccessPoint(client, bucket, objectFcApName, regularAccessPointName); err != nil {
			return fmt.Errorf("creating Object FC access point %s: %s", objectFcApName, err)
		}
		stateConf := helper.BuildStateConf(
			[]string{"creating"},
			[]string{"enable"},
			accessPointForObjectProcessStatusTimeout,
			5*time.Second,
			objectFcAccessPointStateRefreshFunc(client, fmt.Sprintf("%s:%s", bucket, objectFcApName)),
		)
		if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
			return fmt.Errorf("waiting for Object FC access point %s to become ready: %s", objectFcApName, err)
		}
		return nil
	}
}

// createObjectFcAccessPoint issues CreateAccessPointForObjectProcess with the
// minimal body: the referenced regular access point name. The function
// configuration is optional per the API contract and is not required for the
// policy lifecycle under test.
func createObjectFcAccessPoint(client *connectivity.AliyunClient, bucket, objectFcApName, regularAccessPointName string) error {
	query := map[string]*string{
		"x-oss-access-point-for-object-process-name": helper.StringPointer(objectFcApName),
	}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	body := map[string]interface{}{
		"CreateAccessPointForObjectProcessConfiguration": map[string]interface{}{
			"AccessPointName": regularAccessPointName,
		},
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	return retry.Retry(accessPointForObjectProcessStatusTimeout, func() *retry.RetryError {
		_, err := client.Do(accessPointPolicyForObjectProcessProduct, helper.XmlParam("PUT", accessPointPolicyForObjectProcessVersion, "CreateAccessPointForObjectProcess", accessPointForObjectProcessAction), query, body, nil, hostMap, false)
		if err != nil {
			// A leftover access point from a prior run is as good as created for
			// the test's purposes; the pre-clean aims to remove it but tolerates
			// the race.
			if errs.IsExpectedErrors(err, []string{"AccessPointForObjectProcessAlreadyExists"}) {
				return nil
			}
			// The referenced regular access point or its bucket may still be
			// settling; retry those transient conflicts.
			if errs.IsExpectedErrors(err, []string{"NoSuchAccessPoint", "NoSuchAccessPointForObjectProcess", "AccessPointCreatingConflict", "AccessPointForObjectProcessCreatingConflict"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

// tearDownObjectFcAccessPoint removes the policy and the Object FC access point
// out-of-band, polling the access point to gone. OSS may require the policy to
// be deleted before the access point, so the policy is removed first.
func tearDownObjectFcAccessPoint(bucket, objectFcApName string) {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	if client == nil {
		return
	}
	deleteAccessPointPolicyForObjectProcess(client, bucket, objectFcApName)
	deleteObjectFcAccessPoint(client, bucket, objectFcApName)
}

func deleteAccessPointPolicyForObjectProcess(client *connectivity.AliyunClient, bucket, objectFcApName string) {
	query := map[string]*string{
		"x-oss-access-point-for-object-process-name": helper.StringPointer(objectFcApName),
	}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	_ = retry.Retry(2*time.Minute, func() *retry.RetryError {
		_, err := client.Do(accessPointPolicyForObjectProcessProduct, helper.JsonXmlParam("DELETE", accessPointPolicyForObjectProcessVersion, "DeleteAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessAction), query, nil, nil, hostMap, false)
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicyForObjectProcess", "NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
}

func deleteObjectFcAccessPoint(client *connectivity.AliyunClient, bucket, objectFcApName string) {
	query := map[string]*string{
		"x-oss-access-point-for-object-process-name": helper.StringPointer(objectFcApName),
	}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	_ = retry.Retry(accessPointForObjectProcessDeleteTimeout, func() *retry.RetryError {
		_, err := client.Do(accessPointPolicyForObjectProcessProduct, helper.XmlParam("DELETE", accessPointPolicyForObjectProcessVersion, "DeleteAccessPointForObjectProcess", accessPointForObjectProcessAction), query, nil, nil, hostMap, false)
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) {
				return nil
			}
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		// The delete is asynchronous; keep polling until the access point is gone.
		if objectFcAccessPointExists(client, bucket, objectFcApName) {
			wait()
			return retry.RetryableError(fmt.Errorf("Object FC access point %s still deleting", objectFcApName))
		}
		return nil
	})
}

func objectFcAccessPointExists(client *connectivity.AliyunClient, bucket, objectFcApName string) bool {
	query := map[string]*string{
		"x-oss-access-point-for-object-process-name": helper.StringPointer(objectFcApName),
	}
	hostMap := map[string]*string{
		"bucket": helper.StringPointer(bucket),
	}
	_, err := client.Do(accessPointPolicyForObjectProcessProduct, helper.XmlParam("GET", accessPointPolicyForObjectProcessVersion, "GetAccessPointForObjectProcess", accessPointForObjectProcessAction), query, nil, nil, hostMap, true)
	if err != nil {
		return !errs.NotFoundError(err) && !errs.IsExpectedErrors(err, []string{"NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"})
	}
	return true
}

func objectFcAccessPointStateRefreshFunc(client *connectivity.AliyunClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		parts := strings.Split(id, ":")
		if len(parts) != 2 {
			return nil, "", fmt.Errorf("invalid Resource Id %s. Expected bucket:access_point_name", id)
		}
		query := map[string]*string{
			"x-oss-access-point-for-object-process-name": helper.StringPointer(parts[1]),
		}
		hostMap := map[string]*string{
			"bucket": helper.StringPointer(parts[0]),
		}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do(accessPointPolicyForObjectProcessProduct, helper.XmlParam("GET", accessPointPolicyForObjectProcessVersion, "GetAccessPointForObjectProcess", accessPointForObjectProcessAction), query, nil, nil, hostMap, true)
			if requestErr != nil {
				if helper.NeedRetry(requestErr) {
					wait()
					return retry.RetryableError(requestErr)
				}
				return retry.NonRetryableError(requestErr)
			}
			return nil
		})
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) {
				return nil, "", nil
			}
			return nil, "", err
		}
		if response == nil {
			return nil, "", nil
		}
		result, _ := jsonpath.Get("$.GetAccessPointForObjectProcessResult", response)
		if result == nil {
			return nil, "", nil
		}
		obj, ok := result.(map[string]interface{})
		if !ok {
			return nil, "", nil
		}
		status := fmt.Sprint(obj["Status"])
		return obj, status, nil
	}
}

func testAccCheckOssAccessPointPolicyForObjectProcessDestroy(bucket, objectFcApName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_oss_access_point_policy_for_object_process" {
				continue
			}
			parts := strings.SplitN(rs.Primary.ID, ":", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return fmt.Errorf("invalid resource id %s: expected bucket:access_point_name", rs.Primary.ID)
			}
			query := map[string]*string{
				"x-oss-access-point-for-object-process-name": helper.StringPointer(parts[1]),
			}
			hostMap := map[string]*string{
				"bucket": helper.StringPointer(parts[0]),
			}
			var response map[string]interface{}
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err := retry.Retry(1*time.Minute, func() *retry.RetryError {
				var requestErr error
				response, requestErr = client.Do(accessPointPolicyForObjectProcessProduct, helper.XmlJsonParam("GET", accessPointPolicyForObjectProcessVersion, "GetAccessPointPolicyForObjectProcess", accessPointPolicyForObjectProcessAction), query, nil, nil, hostMap, true)
				if requestErr != nil {
					if helper.NeedRetry(requestErr) {
						wait()
						return retry.RetryableError(requestErr)
					}
					return retry.NonRetryableError(requestErr)
				}
				return nil
			})
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"NoSuchAccessPointPolicyForObjectProcess", "NoSuchAccessPointForObjectProcess", "NoSuchAccessPoint"}) || errs.NotFoundError(err) {
					continue
				}
				return fmt.Errorf("describing access point policy for object process %s failed: %s", rs.Primary.ID, err)
			}
			if response != nil {
				return fmt.Errorf("access point policy for object process %s still exists", rs.Primary.ID)
			}
		}
		// Belt-and-suspenders: ensure the out-of-band Object FC access point is
		// gone even if a step failed before the Destroy step could clean it up.
		tearDownObjectFcAccessPoint(bucket, objectFcApName)
		return nil
	}
}

func testAccCheckOssAccessPointPolicyForObjectProcessEffect(resourceId, effect string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(resourceId, "policy", func(value string) error {
		var document map[string]interface{}
		if err := json.Unmarshal([]byte(value), &document); err != nil {
			return fmt.Errorf("policy is not valid JSON: %s", err)
		}
		statements, ok := document["Statement"].([]interface{})
		if !ok {
			return fmt.Errorf("policy has no Statement array: %s", value)
		}
		for _, statement := range statements {
			entry, ok := statement.(map[string]interface{})
			if !ok {
				continue
			}
			if entry["Effect"] == effect {
				return nil
			}
		}
		return fmt.Errorf("policy has no %q statement: %s", effect, value)
	})
}

func testAccAliCloudOssAccessPointPolicyForObjectProcessDocument(effect string) string {
	return fmt.Sprintf(`jsonencode({
    Version = "1"
    Statement = [
      {
        Action    = ["oss:GetObject"]
        Effect    = "%s"
        Principal = [data.alicloud_caller_identity.current.account_id]
        Resource  = ["acs:oss:${var.region}:${data.alicloud_caller_identity.current.account_id}:accesspointforobjectprocess/${var.access_point_name}/object/*"]
      }
    ]
  })`, effect)
}

func testAccAliCloudOssAccessPointPolicyForObjectProcessFormatted(name, region, objectFcApName, effect string) string {
	policy := fmt.Sprintf("replace(%s, \",\", \",\\n\")", testAccAliCloudOssAccessPointPolicyForObjectProcessDocument(effect))
	return testAccAliCloudOssAccessPointPolicyForObjectProcessConfig(name, region, objectFcApName, policy)
}

func testAccAliCloudOssAccessPointPolicyForObjectProcessConfig(name, region, objectFcApName, policy string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

variable "access_point_name" {
  default = "%s"
}

provider "alicloud" {
  region = var.region
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_oss_access_point" "default" {
  access_point_name = var.name
  bucket            = alicloud_oss_bucket.default.bucket
  network_origin    = "vpc"
  vpc_configuration {
    vpc_id = "vpc-abctest"
  }
  public_access_block_configuration {
    block_public_access = true
  }
}

data "alicloud_caller_identity" "current" {}

resource "alicloud_oss_access_point_policy_for_object_process" "default" {
  access_point_name = var.access_point_name
  bucket            = alicloud_oss_bucket.default.bucket
  policy            = %s
}
`, name, region, objectFcApName, policy)
}

func testAccAliCloudOssAccessPointPolicyForObjectProcessPrereq(name, region string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

provider "alicloud" {
  region = var.region
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
}

resource "alicloud_oss_access_point" "default" {
  access_point_name = var.name
  bucket            = alicloud_oss_bucket.default.bucket
  network_origin    = "vpc"
  vpc_configuration {
    vpc_id = "vpc-abctest"
  }
  public_access_block_configuration {
    block_public_access = true
  }
}
`, name, region)
}
