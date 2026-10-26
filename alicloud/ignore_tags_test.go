package alicloud

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUnitIgnoreTagsProviderSchema(t *testing.T) {
	p := Provider()
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("provider schema is invalid: %s", err)
	}

	s, ok := p.Schema["ignore_tags"]
	if !ok {
		t.Fatal("provider block is missing ignore_tags")
	}
	if s.Type != schema.TypeList {
		t.Fatalf("ignore_tags type = %v, want TypeList", s.Type)
	}
	if !s.Optional || s.Required || s.Computed {
		t.Fatalf("ignore_tags must be optional only: optional=%v required=%v computed=%v", s.Optional, s.Required, s.Computed)
	}
	if s.MaxItems != 1 {
		t.Fatalf("ignore_tags MaxItems = %d, want 1", s.MaxItems)
	}
	elem, ok := s.Elem.(*schema.Resource)
	if !ok {
		t.Fatalf("ignore_tags element is %T, want *schema.Resource", s.Elem)
	}
	for _, field := range []string{"keys", "key_prefixes"} {
		fieldSchema, ok := elem.Schema[field]
		if !ok {
			t.Fatalf("ignore_tags is missing %s", field)
		}
		if fieldSchema.Type != schema.TypeSet || !fieldSchema.Optional {
			t.Fatalf("%s must be an optional set: type=%v optional=%v", field, fieldSchema.Type, fieldSchema.Optional)
		}
		element, ok := fieldSchema.Elem.(*schema.Schema)
		if !ok {
			t.Fatalf("%s element is %T, want *schema.Schema", field, fieldSchema.Elem)
		}
		if element.Type != schema.TypeString {
			t.Fatalf("%s element type = %v, want TypeString", field, element.Type)
		}
		if element.ValidateFunc == nil {
			t.Fatalf("%s must reject empty elements", field)
		}
		if _, errs := element.ValidateFunc("", field); len(errs) == 0 {
			t.Fatalf("%s accepted an empty string", field)
		}
		if _, errs := element.ValidateFunc("env", field); len(errs) != 0 {
			t.Fatalf("%s rejected a valid value: %v", field, errs)
		}
	}
	if _, ok := p.Schema["default_tags"]; ok {
		t.Fatal("default_tags is out of scope for this feature and must not be added")
	}
}

func TestUnitExpandIgnoreTags(t *testing.T) {
	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"ignore_tags": []interface{}{map[string]interface{}{
			"keys":         []interface{}{"owner", "env"},
			"key_prefixes": []interface{}{"tf-", "cost:"},
		}},
	})
	raw, ok := d.Get("ignore_tags").([]interface{})
	if !ok || len(raw) != 1 {
		t.Fatalf("ignore_tags = %#v, want one block", d.Get("ignore_tags"))
	}
	got := expandIgnoreTags(raw)
	want := features.IgnoreTags{
		Keys:        []string{"env", "owner"},
		KeyPrefixes: []string{"cost:", "tf-"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expandIgnoreTags = %#v, want %#v", got, want)
	}

	for name, input := range map[string][]interface{}{
		"nil":         nil,
		"empty list":  {},
		"nil block":   {nil},
		"empty block": {map[string]interface{}{}},
	} {
		if expanded := expandIgnoreTags(input); !expanded.Empty() {
			t.Fatalf("%s expanded to %#v, want empty", name, expanded)
		}
	}
}

const ignoreTagsProviderBlock = `
provider "alicloud" {
  ignore_tags {
    keys         = ["env"]
    key_prefixes = ["tf-ignored-"]
  }
}
`

func ignoreTagsVpcConfig(name string, thirdOctet int) string {
	return fmt.Sprintf(`%s
resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "10.%d.0.0/16"

  tags = {
    team = "core"
  }
}
`, ignoreTagsProviderBlock, name, thirdOctet)
}

func ignoreTagsVpcConfigWithDescription(name string, thirdOctet int, description string) string {
	return fmt.Sprintf(`%s
resource "alicloud_vpc" "default" {
  vpc_name    = %q
  description = %q
  cidr_block  = "10.%d.0.0/16"

  tags = {
    team = "core"
  }
}
`, ignoreTagsProviderBlock, name, description, thirdOctet)
}

func ignoreTagsOssBucketConfig(bucket string, tagsBlock string) string {
	return fmt.Sprintf(`%s
resource "alicloud_oss_bucket" "default" {
  bucket = %q

  %s
}
`, ignoreTagsProviderBlock, bucket, tagsBlock)
}

func ignoreTagsEcsSnapshotPolicyConfig(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_ecs_auto_snapshot_policy" "default" {
  auto_snapshot_policy_name = %q
  repeat_weekdays           = ["1"]
  retention_days            = 1
  time_points               = ["1"]

  tags = {
    team = "core"
  }
}
`, ignoreTagsProviderBlock, name)
}

func ignoreTagsSlbConfig(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = %q
  load_balancer_spec   = "slb.s1.small"
  internet_charge_type = "PayByTraffic"

  tags = {
    team = "core"
  }
}
`, ignoreTagsProviderBlock, name)
}

func ignoreTagsBaselineVpcConfig(name string, thirdOctet int) string {
	return fmt.Sprintf(`
resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "10.%d.0.0/16"

  tags = {
    team = "core"
  }
}
`, name, thirdOctet)
}

func ignoreTagsCaptureID(resourceId string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s is not in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource %s has an empty id", resourceId)
		}
		*target = rs.Primary.ID
		return nil
	}
}

func ignoreTagsFromObject(raw interface{}) map[string]string {
	result := make(map[string]string)
	list := raw
	if container, ok := raw.(map[string]interface{}); ok {
		list = container["Tag"]
	}
	items, ok := list.([]interface{})
	if !ok {
		return result
	}
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		key := ""
		if v, ok := m["TagKey"]; ok && v != nil {
			key = fmt.Sprint(v)
		} else if v, ok := m["Key"]; ok && v != nil {
			key = fmt.Sprint(v)
		}
		if key == "" {
			continue
		}
		value := ""
		if v, ok := m["TagValue"]; ok && v != nil {
			value = fmt.Sprint(v)
		} else if v, ok := m["Value"]; ok && v != nil {
			value = fmt.Sprint(v)
		}
		result[key] = value
	}
	return result
}

func ignoreTagsHasKeys(tags map[string]string, want map[string]string) bool {
	for key, value := range want {
		if tags[key] != value {
			return false
		}
	}
	return true
}

func ignoreTagsWaitForKeys(t *testing.T, what string, fetch func() (map[string]string, error), want map[string]string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Minute)
	var last map[string]string
	var lastErr error
	for {
		last, lastErr = fetch()
		if lastErr == nil && ignoreTagsHasKeys(last, want) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s tags %v not visible after 3m, last=%v err=%v", what, want, last, lastErr)
		}
		time.Sleep(5 * time.Second)
	}
}

func ignoreTagsCloudTagsPresent(what string, fetch func() (map[string]string, error), want map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tags, err := fetch()
		if err != nil {
			return fmt.Errorf("failed to read %s tags: %s", what, err)
		}
		for key, value := range want {
			if tags[key] != value {
				return fmt.Errorf("%s lost tag %q=%q, still present: %v", what, key, value, tags)
			}
		}
		return nil
	}
}

func ignoreTagsTagVpc(t *testing.T, client *connectivity.AliyunClient, id string, tags map[string]string) {
	t.Helper()
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "VPC",
		"ResourceId.1": id,
	}
	i := 1
	for key, value := range tags {
		request[fmt.Sprintf("Tag.%d.Key", i)] = key
		request[fmt.Sprintf("Tag.%d.Value", i)] = value
		i++
	}
	if _, err := client.RpcPost("Vpc", "2016-04-28", "TagResources", map[string]interface{}{}, request, true); err != nil {
		t.Fatalf("failed to tag vpc %s out of band: %s", id, err)
	}
}

func ignoreTagsTagEcsSnapshotPolicy(t *testing.T, client *connectivity.AliyunClient, id string, tags map[string]string) {
	t.Helper()
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "snapshotpolicy",
		"ResourceId.1": id,
	}
	i := 1
	for key, value := range tags {
		request[fmt.Sprintf("Tag.%d.Key", i)] = key
		request[fmt.Sprintf("Tag.%d.Value", i)] = value
		i++
	}
	if _, err := client.RpcPost("Ecs", "2014-05-26", "TagResources", nil, request, false); err != nil {
		t.Fatalf("failed to tag ecs auto snapshot policy %s out of band: %s", id, err)
	}
}

func ignoreTagsTagSlb(t *testing.T, client *connectivity.AliyunClient, id string, tags map[string]string) {
	t.Helper()
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": "instance",
		"ResourceId.1": id,
	}
	i := 1
	for key, value := range tags {
		request[fmt.Sprintf("Tag.%d.Key", i)] = key
		request[fmt.Sprintf("Tag.%d.Value", i)] = value
		i++
	}
	if _, err := client.RpcPost("Slb", "2014-05-15", "TagResources", nil, request, false); err != nil {
		t.Fatalf("failed to tag slb %s out of band: %s", id, err)
	}
}

func ignoreTagsTagOssBucket(t *testing.T, client *connectivity.AliyunClient, bucket string, tags map[string]string) {
	t.Helper()
	var tagging oss.Tagging
	for key, value := range tags {
		tagging.Tags = append(tagging.Tags, oss.Tag{Key: key, Value: value})
	}
	if _, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return nil, ossClient.SetBucketTagging(bucket, tagging)
	}); err != nil {
		t.Fatalf("failed to tag oss bucket %s out of band: %s", bucket, err)
	}
}

func ignoreTagsFetchOssBucketTags(client *connectivity.AliyunClient, bucket string) (map[string]string, error) {
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.GetBucketTagging(bucket)
	})
	if err != nil {
		return nil, err
	}
	tagging, _ := raw.(oss.GetBucketTaggingResult)
	result := make(map[string]string)
	for _, tag := range tagging.Tags {
		result[tag.Key] = tag.Value
	}
	return result, nil
}

func ignoreTagsOssBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket" {
			continue
		}
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.IsBucketExist(rs.Primary.ID)
		})
		if err != nil {
			return WrapError(err)
		}
		if exist, _ := raw.(bool); exist {
			return fmt.Errorf("oss bucket %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func ignoreTagsVpcDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vpc" {
			continue
		}
		if _, err := (&VpcServiceV2{client}).DescribeVpcVpc(rs.Primary.ID); err == nil {
			return fmt.Errorf("vpc %s still exists", rs.Primary.ID)
		} else if !NotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

func ignoreTagsEcsSnapshotPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_auto_snapshot_policy" {
			continue
		}
		if _, err := (&EcsServiceV2{client}).DescribeEcsAutoSnapshotPolicy(rs.Primary.ID); err == nil {
			return fmt.Errorf("ecs auto snapshot policy %s still exists", rs.Primary.ID)
		} else if !NotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

func ignoreTagsSlbDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_load_balancer" {
			continue
		}
		if _, err := (&SlbService{client}).DescribeSlbLoadBalancer(rs.Primary.ID); err == nil {
			return fmt.Errorf("slb %s still exists", rs.Primary.ID)
		} else if !NotFoundError(err) {
			return WrapError(err)
		}
	}
	return nil
}

func TestAccAliCloudIgnoreTagsVpc(t *testing.T) {
	resourceId := "alicloud_vpc.default"
	rand := acctest.RandIntRange(100, 250)
	name := fmt.Sprintf("tf-testAccIgnoreTagsVpc%d", rand)
	config := ignoreTagsVpcConfig(name, rand)
	updatedConfig := ignoreTagsVpcConfigWithDescription(name, rand, "ignore-tags update")

	var vpcId string
	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		object, err := (&VpcServiceV2{client}).DescribeVpcVpc(vpcId)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &vpcId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsTagVpc(t, client, vpcId, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
					ignoreTagsWaitForKeys(t, "vpc", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.tf-ignored-created-by"),
					ignoreTagsCloudTagsPresent("vpc", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"}),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "ignore-tags update"),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.tf-ignored-created-by"),
					ignoreTagsCloudTagsPresent("vpc", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"}),
				),
			},
			{
				PreConfig: func() {
					ignoreTagsTagVpc(t, client, vpcId, map[string]string{"cost-center": "cc"})
					ignoreTagsWaitForKeys(t, "vpc", fetchTags, map[string]string{"cost-center": "cc"})
				},
				Config:             updatedConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAliCloudIgnoreTagsEcsAutoSnapshotPolicy(t *testing.T) {
	resourceId := "alicloud_ecs_auto_snapshot_policy.default"
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf-testAccIgnoreTagsSnapshot%d", rand)
	config := ignoreTagsEcsSnapshotPolicyConfig(name)

	var policyId string
	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		object, err := (&EcsServiceV2{client}).DescribeEcsAutoSnapshotPolicy(policyId)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsEcsSnapshotPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &policyId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsTagEcsSnapshotPolicy(t, client, policyId, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
					ignoreTagsWaitForKeys(t, "ecs auto snapshot policy", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.tf-ignored-created-by"),
					ignoreTagsCloudTagsPresent("ecs auto snapshot policy", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"}),
				),
			},
			{
				PreConfig: func() {
					ignoreTagsTagEcsSnapshotPolicy(t, client, policyId, map[string]string{"cost-center": "cc"})
					ignoreTagsWaitForKeys(t, "ecs auto snapshot policy", fetchTags, map[string]string{"cost-center": "cc"})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAliCloudIgnoreTagsSlbLoadBalancer(t *testing.T) {
	resourceId := "alicloud_slb_load_balancer.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccIgnoreTagsSlb%d", rand)
	config := ignoreTagsSlbConfig(name)

	var slbId string
	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		raw, err := (&SlbService{client}).ListTagResources(slbId, "instance")
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(raw), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsSlbDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &slbId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsTagSlb(t, client, slbId, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
					ignoreTagsWaitForKeys(t, "slb", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"})
				},
				Config:   config,
				PlanOnly: true,
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.tf-ignored-created-by"),
					ignoreTagsCloudTagsPresent("slb", fetchTags, map[string]string{"env": "prod", "tf-ignored-created-by": "acc"}),
				),
			},
			{
				PreConfig: func() {
					ignoreTagsTagSlb(t, client, slbId, map[string]string{"cost-center": "cc"})
					ignoreTagsWaitForKeys(t, "slb", fetchTags, map[string]string{"cost-center": "cc"})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAliCloudIgnoreTagsBaselineDriftVpc(t *testing.T) {
	resourceId := "alicloud_vpc.default"
	rand := acctest.RandIntRange(100, 250)
	name := fmt.Sprintf("tf-testAccIgnoreTagsBaseline%d", rand)
	config := ignoreTagsBaselineVpcConfig(name, rand)

	var vpcId string
	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		object, err := (&VpcServiceV2{client}).DescribeVpcVpc(vpcId)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &vpcId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsTagVpc(t, client, vpcId, map[string]string{"env": "prod"})
					ignoreTagsWaitForKeys(t, "vpc", fetchTags, map[string]string{"env": "prod"})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Removing the last managed tag makes OSS delete taggings wholesale unless the ignored tags already
// on the bucket are written back, so this case pins the ignored tag to the cloud across that update.
func TestAccAliCloudIgnoreTagsOssBucketRemoveLastManagedTag(t *testing.T) {
	resourceId := "alicloud_oss_bucket.default"
	bucket := fmt.Sprintf("tf-testacc-ignore-tags-%d", acctest.RandIntRange(100000, 999999))
	configWithTags := ignoreTagsOssBucketConfig(bucket, `tags = {
    team = "core"
  }`)
	configWithoutTags := ignoreTagsOssBucketConfig(bucket, `tags = {}`)

	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		return ignoreTagsFetchOssBucketTags(client, bucket)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsTagOssBucket(t, client, bucket, map[string]string{"team": "core", "env": "prod"})
					ignoreTagsWaitForKeys(t, "oss bucket", fetchTags, map[string]string{"env": "prod"})
				},
				Config:   configWithTags,
				PlanOnly: true,
			},
			{
				Config: configWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("oss bucket", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			{
				Config: configWithoutTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "0"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("oss bucket", fetchTags, map[string]string{"env": "prod"}),
				),
			},
		},
	})
}

func ignoreTagsOssBucketConfigWithoutIgnoreTags(bucket string, tagsBlock string) string {
	return fmt.Sprintf(`resource "alicloud_oss_bucket" "default" {
  bucket = %q

  %s
}
`, bucket, tagsBlock)
}

// F3 regression: with no ignore_tags configuration, an OSS whole-set tag write must keep its
// baseline semantics — the configured tags are sent as-is, so changing a legal key the
// system-prefix filter recognizes (sae.do.not.delete) from old to new must actually write new,
// not write the current value back.
func TestAccAliCloudIgnoreTagsOssBucketUpdateWithoutConfiguration(t *testing.T) {
	resourceId := "alicloud_oss_bucket.default"
	bucket := fmt.Sprintf("tf-testacc-ignore-tags-%d", acctest.RandIntRange(100000, 999999))
	configWithOld := ignoreTagsOssBucketConfigWithoutIgnoreTags(bucket, `tags = {
    "sae.do.not.delete" = "old"
  }`)
	configWithNew := ignoreTagsOssBucketConfigWithoutIgnoreTags(bucket, `tags = {
    "sae.do.not.delete" = "new"
  }`)

	fetchTags := func() (map[string]string, error) {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		return ignoreTagsFetchOssBucketTags(client, bucket)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithOld,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.sae.do.not.delete", "old"),
					ignoreTagsCloudTagsPresent("oss bucket", fetchTags, map[string]string{"sae.do.not.delete": "old"}),
				),
			},
			{
				Config: configWithNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.sae.do.not.delete", "new"),
					ignoreTagsCloudTagsPresent("oss bucket", fetchTags, map[string]string{"sae.do.not.delete": "new"}),
				),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// G2 regressions: whole-set tag write candidates whose update API replaces the
// whole tag set (OOS UpdateTemplate, CMS PutResourceMetricRule labels, ESS
// ModifyScalingConfiguration, ESS ModifyEciScalingConfiguration, CAS
// CreateCertificateName, ACK ModifyClusterNodePool). With ignore_tags configured
// a managed tag change must send the ignored tags back instead of deleting them.
// ---------------------------------------------------------------------------

func ignoreTagsFromValueMap(raw interface{}) map[string]string {
	result := make(map[string]string)
	if m, ok := raw.(map[string]interface{}); ok {
		for key, value := range m {
			result[key] = fmt.Sprint(value)
		}
	}
	return result
}

func ignoreTagsRpcPost(t *testing.T, client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}) {
	t.Helper()
	if _, err := client.RpcPost(product, version, action, nil, request, false); err != nil {
		t.Fatalf("failed to call %s %s out of band: %s", product, action, err)
	}
}

func ignoreTagsJsonTags(tags map[string]string) string {
	raw := make(map[string]interface{}, len(tags))
	for key, value := range tags {
		raw[key] = value
	}
	encoded, err := convertMaptoJsonString(raw)
	if err != nil {
		panic(err)
	}
	return encoded
}

// OOS template: UpdateTemplate replaces the whole tag set when Tags is present.
func TestAccAliCloudIgnoreTagsOosTemplateManagedTagChange(t *testing.T) {
	resourceId := "alicloud_oos_template.default"
	name := fmt.Sprintf("tf-testacc-ignore-tags-oos-%d", acctest.RandIntRange(100000, 999999))
	content := `{"FormatVersion":"OOS-2019-06-01","Description":"ignore tags acc","Tasks":[{"Name":"noop","Action":"ACS::Sleep","Properties":{"Duration":"PT1S"}}]}`
	config := func(tagsBlock string) string {
		return fmt.Sprintf(`%s
resource "alicloud_oos_template" "default" {
  template_name = %q
  content       = %q

  %s
}
`, ignoreTagsProviderBlock, name, content, tagsBlock)
	}
	configCore := config(`tags = {
    team = "core"
  }`)
	configGold := config(`tags = {
    team = "gold"
  }`)

	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		object, err := (&OosService{client}).DescribeOosTemplate(name)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromValueMap(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsOosTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsRpcPost(t, client, "oos", "2019-06-01", "UpdateTemplate", map[string]interface{}{
						"TemplateName": name,
						"Content":      content,
						"RegionId":     client.RegionId,
						"Tags":         ignoreTagsJsonTags(map[string]string{"team": "core", "env": "prod"}),
					})
					ignoreTagsWaitForKeys(t, "oos template", fetchTags, map[string]string{"env": "prod"})
				},
				Config:   configCore,
				PlanOnly: true,
			},
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("oos template", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			{
				// GetTemplate lags behind UpdateTemplate intermittently, so the trailing Read
				// can still see the previous tag set; assert only the cloud state here.
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCloudTagsPresent("oos template", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			// After the cloud converges, a no-op apply refreshes the state to the written
			// tags and the plan becomes empty again.
			{
				PreConfig: func() {
					ignoreTagsWaitForKeys(t, "oos template", fetchTags, map[string]string{"env": "prod", "team": "gold"})
				},
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "gold"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("oos template", fetchTags, map[string]string{"env": "prod", "team": "gold"}),
				),
			},
		},
	})
}

func ignoreTagsOosTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oos_template" {
			continue
		}
		if _, err := (&OosService{client}).DescribeOosTemplate(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("oos template %s still exists", rs.Primary.ID)
	}
	return nil
}

// CMS alarm: PutResourceMetricRule replaces the label set when Labels is present
// and clears it when omitted, and the provider re-sends labels on every update.
func TestAccAliCloudIgnoreTagsCmsAlarmManagedTagChange(t *testing.T) {
	resourceId := "alicloud_cms_alarm.default"
	name := fmt.Sprintf("tf-testacc-ignore-tags-cms-%d", acctest.RandIntRange(100000, 999999))
	config := func(tagsBlock string) string {
		return fmt.Sprintf(`%s
resource "alicloud_cms_monitor_group" "default" {
  monitor_group_name = %q
}

resource "alicloud_cms_alarm" "default" {
  name           = %q
  project        = "acs_ecs_dashboard"
  metric         = "cpu_total"
  period         = 15
  contact_groups = [alicloud_cms_monitor_group.default.id]

  metric_dimensions = "[{\"instanceId\":\"i-ignore-tags-none\"}]"

  escalations_critical {
    statistics = "Average"
    threshold  = "90"
    times      = 1
  }

  %s
}
`, ignoreTagsProviderBlock, name, name, tagsBlock)
	}
	configCore := config(`tags = {
    team = "core"
  }`)
	configGold := config(`tags = {
    team = "gold"
  }`)

	var client *connectivity.AliyunClient
	var alarmId string
	fetchLabels := func() (map[string]string, error) {
		object, err := (&CmsService{client}).DescribeAlarm(alarmId)
		if err != nil {
			return nil, err
		}
		raw, ok := object["Labels"].(map[string]interface{})
		if !ok {
			return map[string]string{}, nil
		}
		return ignoreTagsFromObject(raw["Labels"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsCmsAlarmDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &alarmId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsLabelCmsAlarm(t, client, alarmId, map[string]string{"env": "prod"})
					ignoreTagsWaitForKeys(t, "cms alarm labels", fetchLabels, map[string]string{"env": "prod"})
				},
				Config:   configCore,
				PlanOnly: true,
			},
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("cms alarm labels", fetchLabels, map[string]string{"env": "prod"}),
				),
			},
			{
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "gold"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("cms alarm labels", fetchLabels, map[string]string{"env": "prod"}),
				),
			},
		},
	})
}

func ignoreTagsLabelCmsAlarm(t *testing.T, client *connectivity.AliyunClient, id string, labels map[string]string) {
	t.Helper()
	object, err := (&CmsService{client}).DescribeAlarm(id)
	if err != nil {
		t.Fatalf("failed to describe cms alarm %s: %s", id, err)
	}
	request := map[string]interface{}{
		"RuleId":        id,
		"RuleName":      object["RuleName"],
		"Namespace":     object["Namespace"],
		"MetricName":    object["MetricName"],
		"ContactGroups": object["ContactGroups"],
		"Resources":     object["Resources"],
	}
	if escalations, ok := object["Escalations"].(map[string]interface{}); ok {
		if critical, ok := escalations["Critical"].(map[string]interface{}); ok {
			request["Escalations.Critical.Statistics"] = critical["Statistics"]
			request["Escalations.Critical.ComparisonOperator"] = critical["ComparisonOperator"]
			request["Escalations.Critical.Threshold"] = critical["Threshold"]
			request["Escalations.Critical.Times"] = critical["Times"]
		}
	}
	merged := make(map[string]string)
	if raw, ok := object["Labels"].(map[string]interface{}); ok {
		for key, value := range ignoreTagsFromObject(raw["Labels"]) {
			merged[key] = value
		}
	}
	for key, value := range labels {
		merged[key] = value
	}
	i := 1
	for key, value := range merged {
		request[fmt.Sprintf("Labels.%d.Key", i)] = key
		request[fmt.Sprintf("Labels.%d.Value", i)] = value
		i++
	}
	ignoreTagsRpcPost(t, client, "Cms", "2019-01-01", "PutResourceMetricRule", request)
}

func ignoreTagsCmsAlarmDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_alarm" {
			continue
		}
		if _, err := (&CmsService{client}).DescribeAlarm(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("cms alarm %s still exists", rs.Primary.ID)
	}
	return nil
}

// SSL certificate: CreateCertificateName replaces the whole tag set.
func TestAccAliCloudIgnoreTagsSslCertificateManagedTagChange(t *testing.T) {
	resourceId := "alicloud_ssl_certificates_service_certificate.default"
	name := fmt.Sprintf("tf-testacc-ignore-tags-cas-%d", acctest.RandIntRange(100000, 999999))
	config := func(tagsBlock string) string {
		return fmt.Sprintf(`%s
%s

resource "alicloud_ssl_certificates_service_certificate" "default" {
  certificate_name = %q
  cert             = var.cert
  key              = var.key

  %s
}
`, ignoreTagsProviderBlock, AliCloudSslCertificatesServiceCertificateBasicDependence0(name), name, tagsBlock)
	}
	configCore := config(`tags = {
    team = "core"
  }`)
	configGold := config(`tags = {
    team = "gold"
  }`)

	var certId string
	var client *connectivity.AliyunClient
	fetchTags := func() (map[string]string, error) {
		object, err := (&SslCertificatesServiceServiceV2{client}).DescribeSslCertificatesServiceCertificate(certId)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsSslCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &certId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					request := map[string]interface{}{
						"CertificateId": certId,
						"RegionId":      client.RegionId,
					}
					i := 1
					for key, value := range map[string]string{"team": "core", "env": "prod"} {
						request[fmt.Sprintf("Tags.%d.Key", i)] = key
						request[fmt.Sprintf("Tags.%d.Value", i)] = value
						i++
					}
					ignoreTagsRpcPost(t, client, "cas", "2018-08-13", "CreateCertificateName", request)
					ignoreTagsWaitForKeys(t, "ssl certificate", fetchTags, map[string]string{"env": "prod"})
				},
				Config:   configCore,
				PlanOnly: true,
			},
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ssl certificate", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			{
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "gold"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ssl certificate", fetchTags, map[string]string{"env": "prod"}),
				),
			},
		},
	})
}

func ignoreTagsSslCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ssl_certificates_service_certificate" {
			continue
		}
		if _, err := (&SslCertificatesServiceServiceV2{client}).DescribeSslCertificatesServiceCertificate(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("ssl certificate %s still exists", rs.Primary.ID)
	}
	return nil
}

// ESS scaling configuration: ModifyScalingConfiguration replaces the whole tag
// set when Tags is present.
func TestAccAliCloudIgnoreTagsEssScalingConfigurationManagedTagChange(t *testing.T) {
	resourceId := "alicloud_ess_scaling_configuration.default"
	name := fmt.Sprintf("tf-testacc-ignore-tags-ess-%d", acctest.RandIntRange(100000, 999999))
	config := func(tagsBlock string) string {
		return fmt.Sprintf(`%s
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.c5"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_22_04_x64"
  most_recent = true
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = %q
}
resource "alicloud_security_group" "default" {
  name   = %q
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 1
  scaling_group_name = %q
  vswitch_ids        = [alicloud_vswitch.default.id]
}
resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true

  %s
}
`, ignoreTagsProviderBlock, name, name, name, name, tagsBlock)
	}
	configCore := config(`tags = {
    team = "core"
  }`)
	configGold := config(`tags = {
    team = "gold"
  }`)

	var client *connectivity.AliyunClient
	var configId string
	fetchTags := func() (map[string]string, error) {
		object, err := (&EssService{client}).DescribeEssScalingConfigurationByCommonApi(configId)
		if err != nil {
			return nil, err
		}
		raw, ok := object["Tags"].(map[string]interface{})
		if !ok {
			return map[string]string{}, nil
		}
		return ignoreTagsFromObject(raw["Tag"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsEssScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &configId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					ignoreTagsRpcPost(t, client, "Ess", "2014-08-28", "ModifyScalingConfiguration", map[string]interface{}{
						"ScalingConfigurationId": configId,
						"Tags":                   ignoreTagsJsonTags(map[string]string{"team": "core", "env": "prod"}),
					})
					ignoreTagsWaitForKeys(t, "ess scaling configuration", fetchTags, map[string]string{"env": "prod"})
				},
				Config:   configCore,
				PlanOnly: true,
			},
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ess scaling configuration", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			{
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "gold"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ess scaling configuration", fetchTags, map[string]string{"env": "prod"}),
				),
			},
		},
	})
}

func ignoreTagsEssScalingConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_scaling_configuration" {
			continue
		}
		if _, err := (&EssService{client}).DescribeEssScalingConfigurationByCommonApi(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("ess scaling configuration %s still exists", rs.Primary.ID)
	}
	return nil
}

// ESS ECI scaling configuration: ModifyEciScalingConfiguration replaces the
// whole tag set when Tag.N is present.
func TestAccAliCloudIgnoreTagsEssEciScalingConfigurationManagedTagChange(t *testing.T) {
	resourceId := "alicloud_ess_eci_scaling_configuration.default"
	name := fmt.Sprintf("tf-testacc-ignore-tags-esseci-%d", acctest.RandIntRange(100000, 999999))
	config := func(tagsBlock string) string {
		return fmt.Sprintf(`%s
data "alicloud_zones" "default" {
}
resource "alicloud_vpc" "default" {
  vpc_name   = %q
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = %q
}
resource "alicloud_security_group" "default" {
  name   = %q
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 1
  scaling_group_name = %q
  vswitch_ids        = [alicloud_vswitch.default.id]
  group_type         = "ECI"
}
resource "alicloud_ess_eci_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  security_group_id = alicloud_security_group.default.id

  container_group_name = %q
  restart_policy       = "Never"

  containers {
    name   = "g2probe"
    image  = "registry-vpc.aliyuncs.com/eci_open/alpine:3.5"
    cpu    = "1"
    memory = "2"
  }

  %s
}
`, ignoreTagsProviderBlock, name, name, name, name, name, tagsBlock)
	}
	configCore := config(`tags = {
    team = "core"
  }`)
	configGold := config(`tags = {
    team = "gold"
  }`)

	var client *connectivity.AliyunClient
	var configId string
	fetchTags := func() (map[string]string, error) {
		object, err := (&EssService{client}).DescribeEssEciScalingConfiguration(configId)
		if err != nil {
			return nil, err
		}
		return ignoreTagsFromObject(object["Tags"]), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      ignoreTagsEssEciScalingConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					ignoreTagsCaptureID(resourceId, &configId),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
				),
			},
			{
				PreConfig: func() {
					client = testAccProvider.Meta().(*connectivity.AliyunClient)
					object, err := (&EssService{client}).DescribeEssEciScalingConfiguration(configId)
					if err != nil {
						t.Fatalf("failed to describe ess eci scaling configuration: %s", err)
					}
					request := map[string]interface{}{
						"ScalingConfigurationId":   configId,
						"ScalingConfigurationName": object["ScalingConfigurationName"],
						"ContainerGroupName":       object["ContainerGroupName"],
						"SecurityGroupId":          object["SecurityGroupId"],
					}
					i := 1
					for key, value := range map[string]string{"team": "core", "env": "prod"} {
						request[fmt.Sprintf("Tag.%d.Key", i)] = key
						request[fmt.Sprintf("Tag.%d.Value", i)] = value
						i++
					}
					ignoreTagsRpcPost(t, client, "Ess", "2014-08-28", "ModifyEciScalingConfiguration", request)
					ignoreTagsWaitForKeys(t, "ess eci scaling configuration", fetchTags, map[string]string{"env": "prod"})
				},
				Config:   configCore,
				PlanOnly: true,
			},
			{
				Config: configCore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "core"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ess eci scaling configuration", fetchTags, map[string]string{"env": "prod"}),
				),
			},
			{
				Config: configGold,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.team", "gold"),
					resource.TestCheckNoResourceAttr(resourceId, "tags.env"),
					ignoreTagsCloudTagsPresent("ess eci scaling configuration", fetchTags, map[string]string{"env": "prod"}),
				),
			},
		},
	})
}

func ignoreTagsEssEciScalingConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ess_eci_scaling_configuration" {
			continue
		}
		if _, err := (&EssService{client}).DescribeEssEciScalingConfiguration(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return err
		}
		return fmt.Errorf("ess eci scaling configuration %s still exists", rs.Primary.ID)
	}
	return nil
}
