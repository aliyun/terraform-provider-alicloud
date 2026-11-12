package expressconnectrouter_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnectrouter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudExpressConnectRouterVpcAssociationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_vpc_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterVpcAssociationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}_fake"},
		}),
	}

	associationIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}_fake",
		}),
	}

	vpcIdConfConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}_fake",
		}),
	}

	associationRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "ACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
			"association_id":        "${alicloud_express_connect_router_vpc_association.default.association_id}",
			"vpc_id":                "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
			"status":                "ACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_vpc_association.default.id}_fake"},
			"association_id":        "${alicloud_express_connect_router_vpc_association.default.association_id}_fake",
			"vpc_id":                "${alicloud_express_connect_router_vpc_association.default.vpc_id}_fake",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}_fake",
			"status":                "CREATING",
		}),
	}

	var existAliCloudExpressConnectRouterVpcAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"associations.#":                       "1",
			"associations.0.id":                    acctest.CHECKSET,
			"associations.0.ecr_id":                acctest.CHECKSET,
			"associations.0.association_id":        acctest.CHECKSET,
			"associations.0.vpc_id":                acctest.CHECKSET,
			"associations.0.association_node_type": acctest.CHECKSET,
			"associations.0.vpc_owner_id":          acctest.CHECKSET,
			"associations.0.allowed_prefixes_mode": acctest.CHECKSET,
			"associations.0.description":           "tf-testacc-ds-description",
			"associations.0.status":                acctest.CHECKSET,
			"associations.0.create_time":           acctest.CHECKSET,
			"associations.0.modify_time":           acctest.CHECKSET,
			"associations.0.allowed_prefixes.#":    acctest.CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterVpcAssociationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_express_connect_router_vpc_associations.default",
		ExistMapFunc: existAliCloudExpressConnectRouterVpcAssociationsMapFunc,
		FakeMapFunc:  fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterVpcAssociationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, associationIdConf, vpcIdConfConf, associationRegionIdConf, statusConf, allConf)
}

// TestAccAliCloudExpressConnectRouterVpcAssociationsDataSource_identity0 pins the
// queried association to the resource under test so an unrelated association
// with the same description cannot satisfy the check.
func TestAccAliCloudExpressConnectRouterVpcAssociationsDataSource_identity0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_vpc_associations.default"
	associationResourceId := "alicloud_express_connect_router_vpc_association.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dataSourceVpcAssociationsIdentityConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "associations.0.id", associationResourceId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "associations.0.ecr_id", associationResourceId, "ecr_id"),
					resource.TestCheckResourceAttrPair(resourceId, "associations.0.association_id", associationResourceId, "association_id"),
					resource.TestCheckResourceAttrPair(resourceId, "associations.0.vpc_id", associationResourceId, "vpc_id"),
					resource.TestCheckResourceAttr(resourceId, "associations.0.description", "tf-testacc-ds-description"),
				),
			},
		},
	})
}

// dataSourceVpcAssociationsIdentityConfig scopes the query to the association
// created in the same configuration so the attribute pairs compare against the
// object that was actually created.
func dataSourceVpcAssociationsIdentityConfig(name string) string {
	return dataSourceExpressConnectRouterVpcAssociationsConfig(name) + `
data "alicloud_express_connect_router_vpc_associations" "default" {
  ecr_id = alicloud_express_connect_router_vpc_association.default.ecr_id
  ids    = [alicloud_express_connect_router_vpc_association.default.id]
}
`
}

func dataSourceExpressConnectRouterVpcAssociationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65532"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_express_connect_router_vpc_association" "default" {
  ecr_id                = alicloud_express_connect_router_express_connect_router.default.id
  vpc_id                = alicloud_vpc.default.id
  association_region_id = data.alicloud_regions.default.regions.0.id
  vpc_owner_id          = data.alicloud_account.default.id
  allowed_prefixes      = ["172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"]
  description           = "tf-testacc-ds-description"
}
`, name)
}

func TestUnitVpcAssociationsReadPaginatesToTargetID(t *testing.T) {
	var calls int32
	client := newVpcAssociationUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse request: %v", err)
		}
		if got := r.Form.Get("Action"); got != "DescribeExpressConnectRouterAssociation" {
			t.Errorf("unexpected action %q", got)
		}
		page := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		switch page {
		case 1:
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "unit-test",
				"NextToken": "page-2",
				"AssociationList": []map[string]interface{}{{
					"EcrId":               "ecr-1",
					"AssociationId":       "assoc-a",
					"VpcId":               "vpc-a",
					"AssociationNodeType": "VPC",
					"Status":              "ACTIVE",
					"Description":         "desc-a",
				}},
			}); err != nil {
				t.Errorf("encode page 1: %v", err)
			}
		case 2:
			if got := r.Form.Get("NextToken"); got != "page-2" {
				t.Errorf("second page NextToken = %q, want page-2", got)
			}
			if err := json.NewEncoder(w).Encode(map[string]interface{}{
				"RequestId": "unit-test",
				"AssociationList": []map[string]interface{}{{
					"EcrId":               "ecr-1",
					"AssociationId":       "assoc-b",
					"VpcId":               "vpc-b",
					"AssociationNodeType": "VPC",
					"Status":              "ACTIVE",
					"Description":         "desc-b",
				}},
			}); err != nil {
				t.Errorf("encode page 2: %v", err)
			}
		default:
			t.Errorf("unexpected extra request page %d", page)
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"Code": "InvalidParameter", "Message": "unexpected extra page"})
		}
	})

	dataSource := expressconnectrouter.DataSourceVpcAssociations()
	d := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{
		"ecr_id": "ecr-1",
		"ids":    []interface{}{"ecr-1:assoc-b:vpc-b"},
	})
	if err := dataSource.Read(d, client); err != nil {
		t.Fatalf("read data source: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("DescribeExpressConnectRouterAssociation calls = %d, want 2", got)
	}
	associations, ok := d.Get("associations").([]interface{})
	if !ok || len(associations) != 1 {
		t.Fatalf("associations = %#v, want exactly the target association", d.Get("associations"))
	}
	association, ok := associations[0].(map[string]interface{})
	if !ok {
		t.Fatalf("association entry = %#v, want a map", associations[0])
	}
	if got := association["id"]; got != "ecr-1:assoc-b:vpc-b" {
		t.Errorf("association id = %v, want ecr-1:assoc-b:vpc-b", got)
	}
	if got := association["description"]; got != "desc-b" {
		t.Errorf("association description = %v, want desc-b", got)
	}
}
