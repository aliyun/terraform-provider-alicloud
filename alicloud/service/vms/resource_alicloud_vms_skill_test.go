package vms_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// CloudCreateSkill/CloudDeleteSkill/CloudListSkill target a Voice Messaging Service
// enterprise (call center instance), addressed by EnterpriseId. The shared test
// account owns one; ALICLOUD_VMS_ENTERPRISE_ID selects it. The test is skipped when
// the value is absent, because an enterprise id is a runtime test-account condition
// rather than a resource-integration gap.
const vmsSkillEnterpriseIdEnv = "ALICLOUD_VMS_ENTERPRISE_ID"

func vmsSkillEnterpriseId() string {
	return strings.TrimSpace(os.Getenv(vmsSkillEnterpriseIdEnv))
}

func TestAccAliCloudVMSSkill_basic0(t *testing.T) {
	enterpriseId := vmsSkillEnterpriseId()
	if enterpriseId == "" {
		t.Skipf("skipping TestAccAliCloudVMSSkill_basic0: %s is not set; provide a valid Voice Messaging Service enterprise (call center instance) ID", vmsSkillEnterpriseIdEnv)
	}
	if _, err := strconv.ParseInt(enterpriseId, 10, 64); err != nil {
		t.Fatalf("%s %q is not a valid integer: %v", vmsSkillEnterpriseIdEnv, enterpriseId, err)
	}

	resourceId := "alicloud_vms_skill.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vms-skill-%d", rand)
	comment := "skill for terraform acc test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckVMSSkillDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudVMSSkillBasic0(enterpriseId, name, comment),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("enterprise_id"), knownvalue.StringExact(enterpriseId)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("comment"), knownvalue.StringExact(comment)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("create_time"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudVMSSkillBasic0(enterpriseId, name, comment string) string {
	return fmt.Sprintf(`
resource "alicloud_vms_skill" "default" {
  enterprise_id = "%[1]s"
  name          = "%[2]s"
  comment       = "%[3]s"
}
`, enterpriseId, name, comment)
}

// testAccCheckVMSSkillDestroy verifies the skill is no longer listable: it pages
// CloudListSkill by EnterpriseId and fails if the destroyed skill Id is still
// present. CloudListSkill has no single-object Get, so the read path is the list.
func testAccCheckVMSSkillDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_vms_skill" {
			continue
		}
		enterpriseId, skillId, ok := parseSkillId(rs.Primary.ID)
		if !ok {
			return fmt.Errorf("destroy check: bad resource id %q", rs.Primary.ID)
		}
		found, err := vmsSkillListContainsId(client, enterpriseId, skillId)
		if err != nil {
			return fmt.Errorf("destroy check: calling CloudListSkill for %s: %v", rs.Primary.ID, err)
		}
		if found {
			return fmt.Errorf("the skill %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func vmsSkillListContainsId(client *connectivity.AliyunClient, enterpriseId, skillId int64) (bool, error) {
	const pageSize = 100
	for offset := 0; ; offset += pageSize {
		request := map[string]interface{}{
			"EnterpriseId": enterpriseId,
			"Limit":        pageSize,
			"Offset":       offset,
		}
		response, err := client.RpcPost("Dyvmsapi", "2017-05-25", "CloudListSkill", nil, request, true)
		if err != nil {
			return false, err
		}
		data, _ := response["Data"].(map[string]interface{})
		list, _ := data["List"].([]interface{})
		for _, item := range list {
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			id, idErr := toInt64(obj["Id"])
			if idErr == nil && id == skillId {
				return true, nil
			}
		}
		if hasNext, _ := data["HasNextPage"].(bool); !hasNext {
			return false, nil
		}
	}
}

// parseSkillId mirrors the resource id format "<enterprise_id>:<skill_id>" used
// by the vms package; it is duplicated here because the test is an external
// package and the resource id is the contract under test.
func parseSkillId(id string) (enterpriseId, skillId int64, ok bool) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	eid, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
	if err != nil {
		return 0, 0, false
	}
	sid, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return eid, sid, true
}

// toInt64 coerces a CloudListSkill Id value (float64 / json.Number / numeric
// string) to int64, mirroring the vms package helper.
func toInt64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case float64:
		return int64(n), nil
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(n), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported integer type %T", v)
	}
}
