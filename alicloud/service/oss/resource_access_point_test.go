package oss_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudOssAccessPoint_basic6642(t *testing.T) {
	resourceId := "alicloud_oss_access_point.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket and access point names only accept lowercase letters, digits
	// and hyphens.
	name := fmt.Sprintf("tf-testacc%sossaccesspoint%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssAccessPoint(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssAccessPointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOSSAccessPointBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "access_point_name", name),
					resource.TestCheckResourceAttrSet(resourceId, "bucket"),
					resource.TestCheckResourceAttr(resourceId, "network_origin", "vpc"),
					resource.TestCheckResourceAttr(resourceId, "vpc_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "vpc_configuration.0.vpc_id", "vpc-abctest"),
					resource.TestCheckResourceAttr(resourceId, "public_access_block_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "public_access_block_configuration.0.block_public_access", "true"),
				),
			},
			{
				Config: testAccAliCloudOSSAccessPointUpdateBlockPublicAccess(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "public_access_block_configuration.0.block_public_access", "false"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccPreCheckOssAccessPoint(t *testing.T) {
	acctest.PreCheck(t)
	// Access points are served from cn-hangzhou; the flat-suite gate this case
	// was migrated from converged to cn-hangzhou for every configured region, so
	// pin it here instead of skipping the case.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssAccessPointDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_access_point" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := oss.NewService(client).DescribeAccessPoint(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing access point %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("access point %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudOSSAccessPointBasic(name string) string {
	return testAccAliCloudOSSAccessPointConfig(name, `
  public_access_block_configuration {
    block_public_access = true
  }`)
}

func testAccAliCloudOSSAccessPointUpdateBlockPublicAccess(name string) string {
	return testAccAliCloudOSSAccessPointConfig(name, `
  public_access_block_configuration {
    block_public_access = false
  }`)
}

func testAccAliCloudOSSAccessPointConfig(name, publicAccessBlock string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
}

resource "alicloud_oss_access_point" "default" {
  access_point_name = var.name
  bucket            = alicloud_oss_bucket.default.bucket
  network_origin    = "vpc"
  vpc_configuration {
    vpc_id = "vpc-abctest"
  }
%s
}
`, name, publicAccessBlock)
}

// An access point that is not VPC-origin reports an empty VpcId. A state entry
// written by an older provider version as vpc_configuration = [{vpc_id = ""}]
// must be cleared by refresh, otherwise the ForceNew attribute keeps planning a
// replacement for a configuration that omits the block.
func TestUnitAliCloudOssAccessPointRefreshVpcConfiguration(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		previous             []interface{}
		networkOrigin        string
		vpcId                string
		omitVpcConfiguration bool
		config               map[string]interface{}
		wantCount            string
		wantChanged          bool
	}{
		{
			name:          "internet origin clears the stale empty block",
			previous:      []interface{}{map[string]interface{}{"vpc_id": ""}},
			networkOrigin: "internet",
			vpcId:         "",
			config: map[string]interface{}{
				"access_point_name": "ap-test",
				"bucket":            "test-bucket",
				"network_origin":    "internet",
			},
			wantCount:   "0",
			wantChanged: false,
		},
		{
			// The response may omit the element instead of returning an empty
			// VpcId; the stale block must still be cleared.
			name:                 "internet origin without a VpcConfiguration element clears the stale block",
			previous:             []interface{}{map[string]interface{}{"vpc_id": ""}},
			networkOrigin:        "internet",
			omitVpcConfiguration: true,
			config: map[string]interface{}{
				"access_point_name": "ap-test",
				"bucket":            "test-bucket",
				"network_origin":    "internet",
			},
			wantCount:   "0",
			wantChanged: false,
		},
		{
			name:          "vpc origin keeps the block",
			previous:      []interface{}{map[string]interface{}{"vpc_id": "vpc-test"}},
			networkOrigin: "vpc",
			vpcId:         "vpc-test",
			config: map[string]interface{}{
				"access_point_name": "ap-test",
				"bucket":            "test-bucket",
				"network_origin":    "vpc",
				"vpc_configuration": []interface{}{map[string]interface{}{"vpc_id": "vpc-test"}},
			},
			wantCount:   "1",
			wantChanged: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := oss.ResourceAccessPoint()
			raw := map[string]interface{}{
				"access_point_name": "ap-test",
				"bucket":            "test-bucket",
				"network_origin":    tc.networkOrigin,
				"vpc_configuration": tc.previous,
			}
			d := schema.TestResourceDataRaw(t, r.Schema, raw)
			d.SetId(":ap-test")
			// An empty bucket keeps the OSS gateway from prefixing the endpoint,
			// so the request lands on the local mock server.
			client := ossAccessPointTestClient(t, tc.networkOrigin, tc.vpcId, tc.omitVpcConfiguration)
			refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), d.State(), client)
			if diags.HasError() {
				t.Fatalf("refresh failed: %v", diags)
			}
			if got := refreshed.Attributes["vpc_configuration.#"]; got != tc.wantCount {
				t.Fatalf("vpc_configuration.# after refresh = %q, want %q, attributes: %#v", got, tc.wantCount, refreshed.Attributes)
			}
			diff, err := r.Diff(context.Background(), refreshed, sdkterraform.NewResourceConfigRaw(tc.config), nil)
			if err != nil {
				t.Fatalf("diff failed: %v", err)
			}
			if changed := diff != nil && !diff.Empty(); changed != tc.wantChanged {
				t.Fatalf("plan changed = %v, want %v, diff: %#v", changed, tc.wantChanged, diff)
			}
		})
	}
}

func ossAccessPointTestClient(t *testing.T, networkOrigin, vpcId string, omitVpcConfiguration bool) *connectivity.AliyunClient {
	t.Helper()
	vpcConfiguration := ""
	if !omitVpcConfiguration {
		vpcConfiguration = fmt.Sprintf("<VpcConfiguration><VpcId>%s</VpcId></VpcConfiguration>", vpcId)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.RawQuery, "accessPoint") {
			t.Errorf("unexpected OSS request: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/xml")
		if _, err := fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<GetAccessPointResult>
  <AccessPointName>ap-test</AccessPointName>
  <Bucket>test-bucket</Bucket>
  <NetworkOrigin>%s</NetworkOrigin>
  <Status>enable</Status>
  %s
  <PublicAccessBlockConfiguration>
    <BlockPublicAccess>false</BlockPublicAccess>
  </PublicAccessBlockConfiguration>
</GetAccessPointResult>`, networkOrigin, vpcConfiguration); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("oss", endpoint)
	return client
}
