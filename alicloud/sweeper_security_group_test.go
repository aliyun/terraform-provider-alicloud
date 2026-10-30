package alicloud

// The alicloud_security_group sweeper stays in the root test binary: the
// make sweep target runs only ./alicloud, and the SDKv2 sweeper helpers
// (sharedClientForRegion, PageSizeLarge, getNextpageNumber, VpcService,
// sweepAll) live in the root alicloud test package. The resource
// implementation itself lives in alicloud/service/ecs; the exported
// SweepSecurityGroup method on *ecs.Service is reached via
// ecssvc.NewService(client) so the unexported client field is not
// constructed through a struct literal.

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	ecssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_security_group", &resource.Sweeper{
		Name: "alicloud_security_group",
		F:    testSweepSecurityGroups,
		//When implemented, these should be removed firstly
		Dependencies: []string{
			"alicloud_instance",
			"alicloud_ecs_network_interface",
			"alicloud_cs_kubernetes",
		},
	})
}

func testSweepSecurityGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	var groups []ecs.SecurityGroup
	req := ecs.CreateDescribeSecurityGroupsRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroups(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Security Groups: %s", err)
		}
		resp, _ := raw.(*ecs.DescribeSecurityGroupsResponse)
		if resp == nil || len(resp.SecurityGroups.SecurityGroup) < 1 {
			break
		}
		groups = append(groups, resp.SecurityGroups.SecurityGroup...)

		if len(resp.SecurityGroups.SecurityGroup) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(req.PageNumber)
		if err != nil {
			return err
		}
		req.PageNumber = page
	}

	vpcService := VpcService{client}
	ecsService := ecssvc.NewService(client)
	for _, v := range groups {
		name := v.SecurityGroupName
		id := v.SecurityGroupId
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			// If a Security Group created by other service, it should be fetched by vpc name and deleted.
			if skip {
				if need, err := vpcService.needSweepVpc(v.VpcId, ""); err == nil {
					skip = !need
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Security Group: %s (%s)", name, id)
				continue
			}
		}
		log.Printf("[INFO] Deleting Security Group: %s (%s)", name, id)
		if err := ecsService.SweepSecurityGroup(id); err != nil {
			log.Printf("[ERROR] Failed to delete Security Group (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}
