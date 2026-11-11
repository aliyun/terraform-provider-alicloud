package sweep_test

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"testing"

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
	var matched []slbsdk.ServerCertificate
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
		matched = append(matched, serverCertificate)
	}
	return deleteSlbServerCertificates(matched, func(id string) error {
		deleteReq := slbsdk.CreateDeleteServerCertificateRequest()
		deleteReq.ServerCertificateId = id
		_, err := client.WithSlbClient(func(slbClient *slbsdk.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(deleteReq)
		})
		return err
	})
}

// deleteSlbServerCertificates keeps deleting the remaining certificates but
// returns every failure, so a sweeper cannot report success while one remains.
func deleteSlbServerCertificates(certificates []slbsdk.ServerCertificate, deleteCertificate func(id string) error) error {
	var failures []error
	for _, certificate := range certificates {
		if err := deleteCertificate(certificate.ServerCertificateId); err != nil {
			log.Printf("[ERROR] Failed to delete Slb Server Certificate (%s (%s)): %s", certificate.ServerCertificateName, certificate.ServerCertificateId, err)
			failures = append(failures, fmt.Errorf("deleting Slb Server Certificate %s (%s): %w", certificate.ServerCertificateName, certificate.ServerCertificateId, err))
		}
	}
	return errors.Join(failures...)
}

func TestUnitSweepSlbServerCertificateDocumentationExampleDeleteError(t *testing.T) {
	certificates := []slbsdk.ServerCertificate{
		{ServerCertificateName: "terraform-example-slb-domain-extension", ServerCertificateId: "cert-1"},
		{ServerCertificateName: "terraform-example-slb-domain-extension-2", ServerCertificateId: "cert-2"},
	}
	deleteFailure := errors.New("InvalidServerCertificate.NotFound")

	err := deleteSlbServerCertificates(certificates, func(id string) error {
		if id == "cert-2" {
			return deleteFailure
		}
		return nil
	})
	if err == nil {
		t.Fatal("want the delete failure to be reported, got nil")
	}
	if !errors.Is(err, deleteFailure) {
		t.Fatalf("want wrapped delete failure in %v", err)
	}
	if !strings.Contains(err.Error(), "cert-2") {
		t.Fatalf("want failing certificate identified in %v", err)
	}

	if err := deleteSlbServerCertificates(certificates, func(string) error { return nil }); err != nil {
		t.Fatalf("want nil when every deletion succeeds, got %v", err)
	}
}
