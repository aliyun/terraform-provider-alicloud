package vpc_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	vpcprovider "github.com/aliyun/terraform-provider-alicloud/alicloud/service/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("alicloud_network_acl_attachment", &resource.Sweeper{
		Name: "alicloud_network_acl_attachment",
		F:    testSweepNetworkAclAttachment,
	})
}

func testSweepNetworkAclAttachment(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var networkAcls []vpc.NetworkAcl
	request := vpc.CreateDescribeNetworkAclsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeNetworkAcls(request)
		})
		if err != nil {
			log.Printf("[ERROR] %s get an error: %#v", request.GetActionName(), err)
		}
		response, _ := raw.(*vpc.DescribeNetworkAclsResponse)
		if len(response.NetworkAcls.NetworkAcl) < 1 {
			break
		}
		networkAcls = append(networkAcls, response.NetworkAcls.NetworkAcl...)

		if len(response.NetworkAcls.NetworkAcl) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	for _, nacl := range networkAcls {
		name := nacl.NetworkAclName
		id := nacl.NetworkAclId
		resources := nacl.Resources.Resource
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Network Acl: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Unassociating Network Acl: %s (%s)", name, id)
		request := vpc.CreateUnassociateNetworkAclRequest()
		request.NetworkAclId = id

		unassociateNetworkAclResource := []vpc.UnassociateNetworkAclResource{}
		for i := 0; i < len(resources); i++ {
			vpcSource := vpc.UnassociateNetworkAclResource{
				ResourceId:   resources[i].ResourceId,
				ResourceType: resources[i].ResourceType,
			}
			unassociateNetworkAclResource = append(unassociateNetworkAclResource, vpcSource)
		}
		request.Resource = &unassociateNetworkAclResource

		_, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateNetworkAcl(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to unassociate Network Acl (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}

func TestAccAliCloudVPCNetworkAclAttachment_basic(t *testing.T) {
	resourceId := "alicloud_network_acl_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"network_acl_id": acctest.CHECKSET,
	})
	rand := tfacctest.RandInt()
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckNetworkAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAttachment_create(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_associate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "2",
					}),
				),
			},
			{
				Config: testAccNetworkAclAttachment_unassociate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudVPCNetworkAclAttachment_multiInstance(t *testing.T) {
	resourceId := "alicloud_network_acl_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, map[string]string{
		"network_acl_id": acctest.CHECKSET,
		"resources.#":    "1",
	})
	resourceId2 := "alicloud_network_acl_attachment.default2"
	ra2 := acctest.ResourceAttrInit(resourceId2, map[string]string{
		"network_acl_id": acctest.CHECKSET,
		"resources.#":    "1",
	})
	rand := tfacctest.RandInt()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "tf-testAcc_network_acl", func(string) string {
		return testAccNetworkAclAttachment_multiInstanceDependence(rand)
	})
	testAccCheck := ra.ResourceAttrMapUpdateSet()
	testAccCheck2 := ra2.ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      testAccCheckNetworkAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// Two attachment instances share one network ACL; each instance owns a
				// distinct vswitch and the subset read must keep both instances bound.
				Config: testAccNetworkAclAttachment_multiInstance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId, "resources.*.resource_id", "alicloud_vswitch.default", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "resources.*", map[string]string{
						"resource_type": "VSwitch",
					}),
					testAccCheckNetworkAclAttachmentExists(resourceId2),
					testAccCheck2(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId2, "resources.*.resource_id", "alicloud_vswitch.default2", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId2, "resources.*", map[string]string{
						"resource_type": "VSwitch",
					}),
				),
			},
			{
				// Move one attachment to a third vswitch while the sibling attachment stays
				// bound to the same network ACL; this is the update path that used to
				// converge on the wrong binding set.
				Config: testAccConfig(map[string]interface{}{
					"network_acl_id": "${alicloud_network_acl.default.id}",
					"resources": []map[string]interface{}{{
						"resource_id":   "${alicloud_vswitch.default3.id}",
						"resource_type": "VSwitch",
					}},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId),
					testAccCheck(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId, "resources.*.resource_id", "alicloud_vswitch.default3", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId, "resources.*", map[string]string{
						"resource_type": "VSwitch",
					}),
					testAccCheckNetworkAclAttachmentExists(resourceId2),
					testAccCheck2(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId2, "resources.*.resource_id", "alicloud_vswitch.default2", "id"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceId2, "resources.*", map[string]string{
						"resource_type": "VSwitch",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_acl_id": "${alicloud_network_acl.default.id}",
					"resources": []map[string]interface{}{{
						"resource_id":   "${alicloud_vswitch.default3.id}",
						"resource_type": "VSwitch",
					}},
				}),
				PlanOnly: true,
			},
			{
				// Destroy one attachment while the sibling keeps its binding on the shared network
				// acl: the deleted instance may only release the resources recorded in its own state.
				Config: testAccNetworkAclAttachment_multiInstanceDropSibling(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkAclAttachmentExists(resourceId2),
					testAccCheck2(map[string]string{
						"network_acl_id": acctest.CHECKSET,
						"resources.#":    "1",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceId2, "resources.*.resource_id", "alicloud_vswitch.default2", "id"),
					testAccCheckNetworkAclBindingAbsent("alicloud_network_acl.default", "alicloud_vswitch.default3"),
				),
			},
			{
				Config:   testAccNetworkAclAttachment_multiInstanceDropSibling(rand),
				PlanOnly: true,
			},
		},
	})
}

func testAccNetworkAclAttachmentResourcesFromState(rs *terraform.ResourceState) []vpc.Resource {
	resources := make([]vpc.Resource, 0)
	for k, v := range rs.Primary.Attributes {
		if strings.HasPrefix(k, "resources.") && strings.HasSuffix(k, ".resource_id") {
			resources = append(resources, vpc.Resource{
				ResourceId:   v,
				ResourceType: rs.Primary.Attributes[strings.TrimSuffix(k, ".resource_id")+".resource_type"],
			})
		}
	}
	return resources
}

func testAccCheckNetworkAclAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(Error("Not found: %s", n))
		}
		if rs.Primary.ID == "" {
			return WrapError(Error("No Network Acl Attachment ID is set"))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := NewVpcService(client)

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		// Build the request from the resources recorded for this attachment only. Using the whole
		// network acl would not prove that this attachment is still bound.
		vpcResources := testAccNetworkAclAttachmentResourcesFromState(rs)
		if len(vpcResources) == 0 {
			return WrapError(Error("No resources recorded for Network Acl Attachment %s", n))
		}
		if _, err := vpcService.DescribeNetworkAclAttachment(context.Background(), networkAclId, vpcResources); err != nil {
			return WrapError(err)
		}
		return nil
	}
}

// testAccCheckNetworkAclBindingAbsent proves that the network acl holds no binding record for the
// given vswitch, so a sibling attachment released only its own resources when it was destroyed.
func testAccCheckNetworkAclBindingAbsent(aclResource, vswitchResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		aclState, ok := s.RootModule().Resources[aclResource]
		if !ok {
			return WrapError(Error("Not found: %s", aclResource))
		}
		vswitchState, ok := s.RootModule().Resources[vswitchResource]
		if !ok {
			return WrapError(Error("Not found: %s", vswitchResource))
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		vpcService := NewVpcService(client)
		object, err := vpcService.DescribeNetworkAcl(aclState.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return WrapError(err)
		}
		expected := []vpc.Resource{{ResourceId: vswitchState.Primary.ID, ResourceType: "VSwitch"}}
		absent, err := networkAclAttachmentHasNoBinding(object, expected)
		if err != nil {
			return WrapError(err)
		}
		if !absent {
			return WrapError(Error("Network Acl %s still has a binding record for %s", aclState.Primary.ID, vswitchState.Primary.ID))
		}
		return nil
	}
}

func testAccCheckNetworkAclAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := NewVpcService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_network_acl_attachment" {
			continue
		}

		parts, err := ParseResourceId(rs.Primary.ID, 2)
		if err != nil {
			return WrapError(err)
		}
		networkAclId := parts[0]

		// CheckDestroy must prove that the bindings recorded for this attachment are gone, not that
		// some later snapshot of the acl happens to have no bindings.
		vpcResource := testAccNetworkAclAttachmentResourcesFromState(rs)
		if len(vpcResource) == 0 {
			continue
		}
		object, err := vpcService.DescribeNetworkAcl(networkAclId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
		absent, err := networkAclAttachmentHasNoBinding(object, vpcResource)
		if err != nil {
			return WrapError(err)
		}
		if !absent {
			return WrapError(Error("Network Acl Attachment %s still has binding records after destroy", rs.Primary.ID))
		}
	}
	return nil
}

func TestUnitAlicloudNetworkAclAttachmentDeletedWait(t *testing.T) {
	aclWith := func(records ...interface{}) map[string]interface{} {
		return map[string]interface{}{
			"Resources": map[string]interface{}{"Resource": records},
		}
	}
	record := func(id, status string) map[string]interface{} {
		return map[string]interface{}{"ResourceId": id, "ResourceType": "VSwitch", "Status": status}
	}
	both := []vpc.Resource{
		{ResourceId: "vsw-1", ResourceType: "VSwitch"},
		{ResourceId: "vsw-2", ResourceType: "VSwitch"},
	}
	cases := []struct {
		name     string
		object   map[string]interface{}
		expected []vpc.Resource
		want     bool
		wantErr  bool
	}{
		{name: "single record still unbinding", object: aclWith(record("vsw-1", "UNBINDING")), expected: both[:1], want: false},
		{name: "single record still binded", object: aclWith(record("vsw-1", "BINDED")), expected: both[:1], want: false},
		{name: "single record still binding", object: aclWith(record("vsw-1", "BINDING")), expected: both[:1], want: false},
		{name: "partial disappearance keeps waiting", object: aclWith(record("vsw-2", "UNBINDING")), expected: both, want: false},
		{name: "partial disappearance with sibling binded keeps waiting", object: aclWith(record("vsw-2", "BINDED")), expected: both, want: false},
		{name: "all records gone", object: aclWith(record("vsw-3", "BINDED")), expected: both, want: true},
		{name: "acl has no binding records", object: aclWith(), expected: both, want: true},
		{name: "no expected resources", object: aclWith(record("vsw-1", "BINDING")), expected: nil, want: true},
		{name: "missing resources is empty", object: map[string]interface{}{}, expected: both, want: true},
		{name: "nil resources is empty", object: map[string]interface{}{"Resources": nil}, expected: both, want: true},
		{name: "unexpected resources type fails closed", object: map[string]interface{}{"Resources": "unexpected"}, expected: both, wantErr: true},
		{name: "missing resource list is empty", object: map[string]interface{}{"Resources": map[string]interface{}{}}, expected: both, want: true},
		{name: "nil resource list is empty", object: map[string]interface{}{"Resources": map[string]interface{}{"Resource": nil}}, expected: both, want: true},
		{name: "unexpected resource list type fails closed", object: map[string]interface{}{"Resources": map[string]interface{}{"Resource": "unexpected"}}, expected: both, wantErr: true},
		{name: "malformed record fails closed", object: aclWith("unexpected", record("vsw-1", "BINDED")), expected: both, wantErr: true},
		{name: "missing record status fails closed", object: aclWith(map[string]interface{}{"ResourceId": "vsw-1", "ResourceType": "VSwitch"}), expected: both, wantErr: true},
	}
	for _, c := range cases {
		got, err := networkAclAttachmentHasNoBinding(c.object, c.expected)
		if (err != nil) != c.wantErr {
			t.Fatalf("%s: error=%v, wantErr=%v", c.name, err, c.wantErr)
		}
		if !c.wantErr && got != c.want {
			t.Fatalf("%s: networkAclAttachmentHasNoBinding=%v, want %v", c.name, got, c.want)
		}
	}
}

func TestUnitAlicloudNetworkAclAttachmentPresentResources(t *testing.T) {
	record := func(id, status string) map[string]interface{} {
		return map[string]interface{}{"ResourceId": id, "ResourceType": "VSwitch", "Status": status}
	}
	aclWith := func(records ...interface{}) map[string]interface{} {
		return map[string]interface{}{"Resources": map[string]interface{}{"Resource": records}}
	}
	expected := []vpc.Resource{
		{ResourceId: "vsw-1", ResourceType: "VSwitch"},
		{ResourceId: "vsw-2", ResourceType: "VSwitch"},
	}
	for _, status := range []string{"BINDING", "BINDED", "UNBINDING"} {
		present, err := networkAclAttachmentPresentResources(aclWith(record("vsw-1", status), record("vsw-sibling", "BINDING")), expected)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", status, err)
		}
		if len(present) != 1 || present[0].ResourceId != "vsw-1" {
			t.Fatalf("%s: present=%v, want only vsw-1", status, present)
		}
	}
	present, err := networkAclAttachmentPresentResources(aclWith(record("vsw-sibling", "BINDED")), expected)
	if err != nil {
		t.Fatalf("all desired records absent: unexpected error: %v", err)
	}
	if len(present) != 0 {
		t.Fatalf("all desired records absent: present=%v, want empty", present)
	}
	if _, err := networkAclAttachmentPresentResources(aclWith("unexpected"), expected); err == nil {
		t.Fatal("malformed record: expected an error")
	}
}

func TestUnitAlicloudNetworkAclAttachmentAvailableWait(t *testing.T) {
	record := func(id, status string) map[string]interface{} {
		return map[string]interface{}{"ResourceId": id, "ResourceType": "VSwitch", "Status": status}
	}
	aclWith := func(status string, records ...interface{}) map[string]interface{} {
		return map[string]interface{}{
			"Status":    status,
			"Resources": map[string]interface{}{"Resource": records},
		}
	}
	expected := []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}}
	cases := []struct {
		name    string
		object  map[string]interface{}
		want    bool
		wantErr bool
	}{
		{name: "desired binded and sibling binding", object: aclWith("Available", record("vsw-1", "BINDED"), record("vsw-sibling", "BINDING")), want: true},
		{name: "desired binded and sibling unbinding", object: aclWith("Available", record("vsw-1", "BINDED"), record("vsw-sibling", "UNBINDING")), want: true},
		{name: "desired binding keeps waiting", object: aclWith("Available", record("vsw-1", "BINDING")), want: false},
		{name: "desired unbinding keeps waiting", object: aclWith("Available", record("vsw-1", "UNBINDING")), want: false},
		{name: "desired missing keeps waiting", object: aclWith("Available", record("vsw-sibling", "BINDED")), want: false},
		{name: "acl status keeps waiting", object: aclWith("Updating", record("vsw-1", "BINDED")), want: false},
		{name: "malformed sibling record fails closed", object: aclWith("Available", record("vsw-1", "BINDED"), "unexpected"), wantErr: true},
	}
	for _, c := range cases {
		got, err := networkAclAttachmentAllBound(c.object, expected, Available)
		if (err != nil) != c.wantErr {
			t.Fatalf("%s: error=%v, wantErr=%v", c.name, err, c.wantErr)
		}
		if !c.wantErr && got != c.want {
			t.Fatalf("%s: networkAclAttachmentAllBound=%v, want %v", c.name, got, c.want)
		}
	}
}

func TestUnitAlicloudNetworkAclAttachmentResourcesToUnbind(t *testing.T) {
	aclWith := func(records ...interface{}) map[string]interface{} {
		return map[string]interface{}{
			"Resources": map[string]interface{}{"Resource": records},
		}
	}
	record := func(id, resourceType string) map[string]interface{} {
		return map[string]interface{}{"ResourceId": id, "ResourceType": resourceType, "Status": "BINDED"}
	}
	ids := func(resources []vpc.UnassociateNetworkAclResource) string {
		values := make([]string, 0, len(resources))
		for _, res := range resources {
			values = append(values, res.ResourceType+":"+res.ResourceId)
		}
		return strings.Join(values, ",")
	}
	cases := []struct {
		name     string
		object   map[string]interface{}
		expected []vpc.Resource
		want     string
		wantErr  bool
	}{
		{
			name:     "only the resources of the deleted attachment are unbound",
			object:   aclWith(record("vsw-1", "VSwitch"), record("vsw-2", "VSwitch")),
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "VSwitch:vsw-1",
		},
		{
			name:     "every expected resource is unbound and sibling bindings are untouched",
			object:   aclWith(record("vsw-1", "VSwitch"), record("vsw-2", "VSwitch"), record("vsw-3", "VSwitch")),
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}, {ResourceId: "vsw-3", ResourceType: "VSwitch"}},
			want:     "VSwitch:vsw-1,VSwitch:vsw-3",
		},
		{
			name:     "already unbound resources produce no request",
			object:   aclWith(record("vsw-2", "VSwitch")),
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "",
		},
		{
			name:     "records with another resource type do not match",
			object:   aclWith(record("vsw-1", "Other")),
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "",
		},
		{
			name:     "missing acl object produces no request",
			object:   nil,
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "",
		},
		{
			name:     "acl without a resources block produces no request",
			object:   map[string]interface{}{},
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "",
		},
		{
			name:     "acl without a resource list produces no request",
			object:   map[string]interface{}{"Resources": map[string]interface{}{}},
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			want:     "",
		},
		{
			name:     "unexpected resources shape fails closed",
			object:   map[string]interface{}{"Resources": "unexpected"},
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			wantErr:  true,
		},
		{
			name:     "malformed records fail closed",
			object:   aclWith("unexpected", record("vsw-1", "VSwitch")),
			expected: []vpc.Resource{{ResourceId: "vsw-1", ResourceType: "VSwitch"}},
			wantErr:  true,
		},
		{
			name:     "no expected resources produces no request",
			object:   aclWith(record("vsw-1", "VSwitch")),
			expected: nil,
			want:     "",
		},
	}
	for _, c := range cases {
		resources, err := networkAclAttachmentResourcesToUnbind(c.object, c.expected)
		if (err != nil) != c.wantErr {
			t.Fatalf("%s: error=%v, wantErr=%v", c.name, err, c.wantErr)
		}
		if got := ids(resources); !c.wantErr && got != c.want {
			t.Fatalf("%s: networkAclAttachmentResourcesToUnbind=%q, want %q", c.name, got, c.want)
		}
	}
}

func testNetworkAclAttachmentClient(t *testing.T, handler http.Handler) *connectivity.AliyunClient {
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := connectivity.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("vpc", endpoint)
	endpoints.Store("Vpc", endpoint)
	return client
}

func testNetworkAclAttachmentDescribeResponse(resourceIds ...string) string {
	return testNetworkAclAttachmentDescribeResponseWithStatus("BINDED", resourceIds...)
}

func testNetworkAclAttachmentDescribeResponseWithStatus(status string, resourceIds ...string) string {
	records := make([]string, 0, len(resourceIds))
	for _, resourceId := range resourceIds {
		records = append(records, fmt.Sprintf(`{"ResourceId":%q,"ResourceType":"VSwitch","Status":%q}`, resourceId, status))
	}
	return fmt.Sprintf(`{"NetworkAclAttribute":{"Status":"Available","Resources":{"Resource":[%s]}}}`, strings.Join(records, ","))
}

func testNetworkAclAttachmentWriteError(w http.ResponseWriter, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = fmt.Fprintf(w, `{"Code":%q,"Message":"fixture error"}`, code)
}

func testNetworkAclAttachmentUpdateData(t *testing.T, oldResourceId, newResourceId string) (*schema.Resource, *schema.ResourceData) {
	res := vpcprovider.ResourceAliyunNetworkAclAttachment()
	attributes := map[string]interface{}{
		"network_acl_id": "nacl-test",
		"resources": []interface{}{map[string]interface{}{
			"resource_id": oldResourceId, "resource_type": "VSwitch",
		}},
	}
	dInit := schema.TestResourceDataRaw(t, res.Schema, attributes)
	dInit.SetId("nacl-test:attachment")
	diff, err := newInstanceDiff("alicloud_network_acl_attachment", attributes, map[string]interface{}{
		"resources": []interface{}{map[string]interface{}{
			"resource_id": newResourceId, "resource_type": "VSwitch",
		}},
	}, dInit.State())
	if err != nil {
		t.Fatal(err)
	}
	d, err := schema.InternalMap(res.Schema).Data(dInit.State(), diff)
	if err != nil {
		t.Fatal(err)
	}
	return res, d
}

func testNetworkAclAttachmentDeleteData(t *testing.T, resourceId string) (*schema.Resource, *schema.ResourceData) {
	res := vpcprovider.ResourceAliyunNetworkAclAttachment()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"network_acl_id": "nacl-test",
		"resources": []interface{}{map[string]interface{}{
			"resource_id": resourceId, "resource_type": "VSwitch",
		}},
	})
	d.SetId("nacl-test:attachment")
	return res, d
}

func TestUnitAlicloudNetworkAclAttachmentUpdateWriteRetry(t *testing.T) {
	type call struct{ action, token string }
	var calls []call
	unassociateAttempts := 0
	associateAttempts := 0
	describeCalls := 0
	client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		action := r.Form.Get("Action")
		calls = append(calls, call{action: action, token: r.Form.Get("ClientToken")})
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "UnassociateNetworkAcl":
			unassociateAttempts++
			if unassociateAttempts == 1 {
				testNetworkAclAttachmentWriteError(w, "OperationDenied.NetworkAclAttachmentInMiddleStatus")
				return
			}
		case "AssociateNetworkAcl":
			associateAttempts++
			if associateAttempts == 1 {
				testNetworkAclAttachmentWriteError(w, "NetworkStatus.Modifying")
				return
			}
		case "DescribeNetworkAclAttributes":
			describeCalls++
			if describeCalls == 1 {
				_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse()))
				return
			}
			_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse("vsw-new")))
			return
		}
		_, _ = w.Write([]byte(`{"RequestId":"fixture"}`))
	}))
	res, d := testNetworkAclAttachmentUpdateData(t, "vsw-old", "vsw-new")
	timeout := 15 * time.Second
	res.Timeouts.Update = &timeout
	if diags := res.UpdateContext(context.Background(), d, client); diags.HasError() {
		t.Fatalf("UpdateContext returned an error: %v", diags)
	}
	actions := make([]string, 0, len(calls))
	for _, c := range calls {
		actions = append(actions, c.action)
	}
	want := "UnassociateNetworkAcl,UnassociateNetworkAcl,DescribeNetworkAclAttributes,AssociateNetworkAcl,AssociateNetworkAcl,DescribeNetworkAclAttributes,DescribeNetworkAclAttributes"
	if got := strings.Join(actions, ","); got != want {
		t.Fatalf("request order=%s, want %s", got, want)
	}
	if calls[0].token == "" || calls[0].token != calls[1].token {
		t.Fatalf("Unassociate ClientToken changed across retries: %q, %q", calls[0].token, calls[1].token)
	}
	if calls[3].token == "" || calls[3].token != calls[4].token {
		t.Fatalf("Associate ClientToken changed across retries: %q, %q", calls[3].token, calls[4].token)
	}
}

func TestUnitAlicloudNetworkAclAttachmentUpdateIdempotentUnbind(t *testing.T) {
	for _, code := range []string{"InvalidResource.NotBinding", "InvalidResource.NotFound", "Mismatch.NetworkAclAndVSwitch"} {
		t.Run(code, func(t *testing.T) {
			describeCalls := 0
			associateCalls := 0
			client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "UnassociateNetworkAcl":
					testNetworkAclAttachmentWriteError(w, code)
				case "AssociateNetworkAcl":
					associateCalls++
					_, _ = w.Write([]byte(`{"RequestId":"fixture"}`))
				case "DescribeNetworkAclAttributes":
					describeCalls++
					if describeCalls == 1 {
						_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse()))
						return
					}
					_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse("vsw-new")))
				}
			}))
			res, d := testNetworkAclAttachmentUpdateData(t, "vsw-old", "vsw-new")
			if diags := res.UpdateContext(context.Background(), d, client); diags.HasError() {
				t.Fatalf("UpdateContext rejected an already-absent binding: %v", diags)
			}
			if describeCalls != 3 || associateCalls != 1 {
				t.Fatalf("idempotent unbind skipped convergence: describes=%d associates=%d", describeCalls, associateCalls)
			}
		})
	}
}

func TestUnitAlicloudNetworkAclAttachmentUpdateSharedOperationBudget(t *testing.T) {
	describeCalls := 0
	associateCalls := 0
	client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "UnassociateNetworkAcl":
			time.Sleep(40 * time.Millisecond)
			_, _ = w.Write([]byte(`{"RequestId":"fixture"}`))
		case "AssociateNetworkAcl":
			associateCalls++
			_, _ = w.Write([]byte(`{"RequestId":"fixture"}`))
		case "DescribeNetworkAclAttributes":
			describeCalls++
			if describeCalls == 1 {
				_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse()))
				return
			}
			_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponseWithStatus("BINDING", "vsw-new")))
		}
	}))
	res, d := testNetworkAclAttachmentUpdateData(t, "vsw-old", "vsw-new")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	diags := res.UpdateContext(ctx, d, client)
	if !diags.HasError() {
		t.Fatal("UpdateContext succeeded while the desired binding remained BINDING")
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("UpdateContext did not honor its shared 100ms operation deadline: %s", elapsed)
	}
	if associateCalls != 1 || describeCalls < 2 {
		t.Fatalf("test did not cross removed and Available phases: describes=%d associates=%d", describeCalls, associateCalls)
	}
}

func TestUnitAlicloudNetworkAclAttachmentDeleteEmptyBindingResponse(t *testing.T) {
	responses := []struct {
		name     string
		response string
	}{
		{name: "missing resources", response: `{"NetworkAclAttribute":{"Status":"Available"}}`},
		{name: "null resources", response: `{"NetworkAclAttribute":{"Status":"Available","Resources":null}}`},
		{name: "missing resource list", response: `{"NetworkAclAttribute":{"Status":"Available","Resources":{}}}`},
		{name: "null resource list", response: `{"NetworkAclAttribute":{"Status":"Available","Resources":{"Resource":null}}}`},
		{name: "empty resource list", response: `{"NetworkAclAttribute":{"Status":"Available","Resources":{"Resource":[]}}}`},
	}
	for _, tc := range responses {
		t.Run(tc.name, func(t *testing.T) {
			actions := make([]string, 0, 2)
			client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				actions = append(actions, r.Form.Get("Action"))
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.response))
			}))
			res, d := testNetworkAclAttachmentDeleteData(t, "vsw-old")
			if diags := res.DeleteContext(context.Background(), d, client); diags.HasError() {
				t.Fatalf("DeleteContext rejected an empty binding response: %v", diags)
			}
			if got := strings.Join(actions, ","); got != "DescribeNetworkAclAttributes,DescribeNetworkAclAttributes" {
				t.Fatalf("DeleteContext actions=%s, want two Describes and no write", got)
			}
		})
	}
}

func TestUnitAlicloudNetworkAclAttachmentDeleteIdempotentRace(t *testing.T) {
	for _, code := range []string{"InvalidResource.NotFound", "Mismatch.NetworkAclAndVSwitch"} {
		t.Run(code, func(t *testing.T) {
			var actions []string
			describeCalls := 0
			client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					t.Error(err)
				}
				action := r.Form.Get("Action")
				actions = append(actions, action)
				w.Header().Set("Content-Type", "application/json")
				switch action {
				case "UnassociateNetworkAcl":
					testNetworkAclAttachmentWriteError(w, code)
				case "DescribeNetworkAclAttributes":
					describeCalls++
					if describeCalls == 1 {
						_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse("vsw-old")))
						return
					}
					_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse()))
				}
			}))
			res, d := testNetworkAclAttachmentDeleteData(t, "vsw-old")
			if diags := res.DeleteContext(context.Background(), d, client); diags.HasError() {
				t.Fatalf("DeleteContext rejected an already-absent binding: %v", diags)
			}
			if got := strings.Join(actions, ","); got != "DescribeNetworkAclAttributes,UnassociateNetworkAcl,DescribeNetworkAclAttributes" {
				t.Fatalf("DeleteContext did not prove Deleted after the idempotent error: %s", got)
			}
		})
	}
}

func TestUnitAlicloudNetworkAclAttachmentDeleteIdempotentStillWaits(t *testing.T) {
	describeCalls := 0
	client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.Form.Get("Action") {
		case "UnassociateNetworkAcl":
			testNetworkAclAttachmentWriteError(w, "Mismatch.NetworkAclAndVSwitch")
		case "DescribeNetworkAclAttributes":
			describeCalls++
			_, _ = w.Write([]byte(testNetworkAclAttachmentDescribeResponse("vsw-old")))
		}
	}))
	res, d := testNetworkAclAttachmentDeleteData(t, "vsw-old")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	diags := res.DeleteContext(ctx, d, client)
	if !diags.HasError() {
		t.Fatal("DeleteContext succeeded although the binding record remained after the idempotent API error")
	}
	if describeCalls < 2 {
		t.Fatalf("DeleteContext returned before the Deleted wait: describes=%d", describeCalls)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("DeleteContext did not honor its 50ms operation deadline: %s", elapsed)
	}
}

func TestUnitAlicloudNetworkAclAttachmentMalformedAttribute(t *testing.T) {
	for _, response := range []string{
		`{"NetworkAclAttribute":null}`,
		`{"NetworkAclAttribute":[]}`,
		`{"NetworkAclAttribute":"invalid"}`,
	} {
		t.Run(response, func(t *testing.T) {
			client := testNetworkAclAttachmentClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(response))
			}))
			res, d := testNetworkAclAttachmentDeleteData(t, "vsw-old")
			if diags := res.ReadContext(context.Background(), d, client); !diags.HasError() {
				t.Fatalf("ReadContext accepted a malformed NetworkAclAttribute: %s", response)
			}
		})
	}
}

func TestUnitAlicloudNetworkAclAttachmentOperationBudget(t *testing.T) {
	now := time.Unix(100, 0)
	deadline := now.Add(15 * time.Minute)
	firstPhase := networkAclAttachmentRemaining(deadline, now.Add(4*time.Minute))
	secondPhase := networkAclAttachmentRemaining(deadline, now.Add(8*time.Minute))
	if firstPhase != 11*time.Minute || secondPhase != 7*time.Minute {
		t.Fatalf("sequential phases reset the operation budget: first=%s second=%s", firstPhase, secondPhase)
	}
	if got := networkAclAttachmentRemaining(now.Add(20*time.Second), now); got != 20*time.Second {
		t.Fatalf("custom short timeout=%s, want 20s", got)
	}
	if got := networkAclAttachmentRemaining(deadline, now.Add(16*time.Minute)); got != -time.Minute {
		t.Fatalf("expired deadline=%s, want -1m", got)
	}
}

func TestUnitAlicloudNetworkAclAttachmentWriteRetry(t *testing.T) {
	for _, code := range []string{"NetworkStatus.Modifying", "ResourceStatus.Error", "OperationDenied.NetworkAclAttachmentInMiddleStatus"} {
		t.Run(code, func(t *testing.T) {
			attempts := 0
			err := networkAclAttachmentRetry(context.Background(), time.Now().Add(5*time.Second), func() error {
				attempts++
				if attempts == 1 {
					return sdkerrors.NewServerError(400, fmt.Sprintf(`{"Code":%q,"Message":"transient"}`, code), "")
				}
				return nil
			})
			if err != nil || attempts != 2 {
				t.Fatalf("transient write error: attempts=%d error=%v, want two attempts and success", attempts, err)
			}
		})
	}

	attempts := 0
	permanent := sdkerrors.NewServerError(400, `{"Code":"InvalidParameter","Message":"terminal"}`, "")
	err := networkAclAttachmentRetry(context.Background(), time.Now().Add(5*time.Second), func() error {
		attempts++
		return permanent
	})
	if err == nil || attempts != 1 {
		t.Fatalf("terminal write error was replayed: attempts=%d error=%v", attempts, err)
	}
}

func TestUnitAlicloudNetworkAclAttachmentDeleteIdempotentErrors(t *testing.T) {
	for _, code := range []string{"InvalidNetworkAcl.NotFound", "InvalidResource.NotBinding", "InvalidResource.NotFound", "Mismatch.NetworkAclAndVSwitch"} {
		err := sdkerrors.NewServerError(400, fmt.Sprintf(`{"Code":%q,"Message":"binding already absent"}`, code), "")
		if !networkAclAttachmentDeleteIdempotent(err) {
			t.Fatalf("%s after the Describe snapshot must continue to the Deleted wait", code)
		}
	}
	err := sdkerrors.NewServerError(403, `{"Code":"Forbidden","Message":"access denied"}`, "")
	if networkAclAttachmentDeleteIdempotent(err) {
		t.Fatal("a permission error must not be treated as an absent binding")
	}
}

func testAccNetworkAclAttachment_create(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

func testAccNetworkAclAttachment_associate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

func testAccNetworkAclAttachment_unassociate(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

func testAccNetworkAclAttachment_multiInstanceDependence(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}multi%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.2.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.3.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default3" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.4.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default2" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  depends_on     = [alicloud_network_acl_attachment.default]
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

// testAccNetworkAclAttachment_multiInstanceDropSibling keeps the sibling attachment bound to the
// shared network acl while the attachment that moved to the third vswitch is removed, so its
// destroy must not release the sibling's binding.
func testAccNetworkAclAttachment_multiInstanceDropSibling(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}multi%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.2.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.3.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default3" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.4.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default2" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}

func testAccNetworkAclAttachment_multiInstance(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc_network_acl"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_network_acl" "default" {
  vpc_id           = "${alicloud_vpc.default.id}"
  network_acl_name = "${var.name}multi%d"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.2.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_vswitch" "default2" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.3.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_network_acl_attachment" "default" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  resources {
    resource_id   = "${alicloud_vswitch.default.id}"
    resource_type = "VSwitch"
  }
}

resource "alicloud_network_acl_attachment" "default2" {
  network_acl_id = "${alicloud_network_acl.default.id}"
  depends_on     = [alicloud_network_acl_attachment.default]
  resources {
    resource_id   = "${alicloud_vswitch.default2.id}"
    resource_type = "VSwitch"
  }
}
`, randInt)
}
