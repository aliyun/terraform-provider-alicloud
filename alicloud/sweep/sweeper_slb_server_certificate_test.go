package sweep_test

import (
	"fmt"
	"log"
	"strings"

	slbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfsweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_slb_server_certificate_documentation_example", &resource.Sweeper{
		Name: "alicloud_slb_server_certificate_documentation_example",
		F:    testSweepSlbServerCertificateDocumentationExample,
	})
}

func testSweepSlbServerCertificateDocumentationExample(region string) error {
	rawClient, err := tfsweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	// Only documentation example certificates are in scope here; test
	// certificates keep their existing sweeper and prefix set.
	prefixes := []string{
		"terraform-example",
	}

	req := slbsdk.CreateDescribeServerCertificatesRequest()
	req.RegionId = client.RegionId
	raw, err := client.WithSlbClient(func(slbClient *slbsdk.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(req)
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(*slbsdk.DescribeServerCertificatesResponse)
	for _, serverCertificate := range resp.ServerCertificates.ServerCertificate {
		name := serverCertificate.ServerCertificateName
		id := serverCertificate.ServerCertificateId

		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Slb Server Certificate: %s (%s)", name, id)
			continue
		}
		log.Printf("[INFO] Deleting Slb Server Certificate : %s (%s)", name, id)

		deleteReq := slbsdk.CreateDeleteServerCertificateRequest()
		deleteReq.ServerCertificateId = id
		if _, err := client.WithSlbClient(func(slbClient *slbsdk.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(deleteReq)
		}); err != nil {
			log.Printf("[ERROR] Failed to delete Slb Server Certificate (%s (%s)): %s", name, id, err)
		}
	}
	return nil
}
