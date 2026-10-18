package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// sslSweeperPrefixes are the name prefixes the SSL Certificates Service sweepers
// treat as test fixtures. A resource whose name does not begin with one of these is
// left untouched unless ALICLOUD_SWEEP_ALL_RESOURCES=true is set; the
// certificate-bearing PCA resources keep the prefix check active even then, so a
// real cert is never swept (red-line: preserve real certs/domains).
var sslSweeperPrefixes = []string{
	"tf-testAcc",
	"tf_testAcc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tftestacc",
}

// sslField returns item[key] as a lower-cased string, tolerating absent or
// non-string values so a sweeper never panics on a shape the API did not return.
func sslField(item map[string]interface{}, key string) string {
	v, ok := item[key]
	if !ok || v == nil {
		return ""
	}
	return strings.ToLower(fmt.Sprint(v))
}

// sslSweepNameOk is the standard fleet-wipe gate for the non-certificate-bearing
// SSL resources (company, contact). A name carrying a NODELETING marker is always
// protected; otherwise the tf-prefix check runs only in non-SWEEP_ALL mode, so a
// fleet wipe (ALICLOUD_SWEEP_ALL_RESOURCES=true) cleans every non-NODELETING item.
func sslSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	if sweepAll() {
		return true
	}
	for _, prefix := range sslSweeperPrefixes {
		if strings.HasPrefix(name, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// sslCertSweepNameOk is the always-active gate for the certificate-bearing PCA
// resources (pca_cert, pca_certificate). The tf-prefix check runs regardless of
// SWEEP_ALL: PCA certs may back real private CAs and issued certificates, so a
// fleet wipe must never touch a cert the tests did not create (red-line #5).
func sslCertSweepNameOk(name string) bool {
	if protectedSweepVpcName(name) {
		return false
	}
	for _, prefix := range sslSweeperPrefixes {
		if strings.HasPrefix(name, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// sslRpcPost wraps an RPC call with the retry/backoff pattern shared by every SSL
// sweeper. The `var err error` is pre-declared so the retry closure assigns to the
// outer scope rather than shadowing it with `:=`, which would swallow the error
// that retry.Retry is supposed to surface.
func sslRpcPost(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	query := make(map[string]interface{})
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost(product, version, action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	return response, err
}

// sslList navigates a jsonpath into an RPC response and returns the resulting
// slice. A missing path or a non-array value yields nil, which ranges cleanly.
func sslList(response map[string]interface{}, path string) []interface{} {
	raw, err := jsonpath.Get(path, response)
	if err != nil {
		return nil
	}
	list, _ := raw.([]interface{})
	return list
}

// sslDelete issues a delete RPC and logs the outcome without ever failing the
// sweep: a NotFound is treated as success (already gone), any other error is
// logged at [ERROR] and the sweep moves on (red-line: ALWAYS return nil).
func sslDelete(client *connectivity.AliyunClient, product, version, action, rtype, id string, request map[string]interface{}) {
	if _, err := sslRpcPost(client, product, version, action, request, 1*time.Minute); err != nil {
		if NotFoundError(err) {
			log.Printf("[INFO] Delete ssl %s already gone: %s", rtype, id)
		} else {
			log.Printf("[ERROR] Failed to delete ssl %s (%s): %s", rtype, id, err)
		}
		return
	}
	log.Printf("[INFO] Delete ssl %s success: %s", rtype, id)
}

func init() {
	resource.AddTestSweepers("alicloud_ssl_certificates_service_company", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_company",
		F:    testSweepSslCertificatesServiceCompany,
	})
	resource.AddTestSweepers("alicloud_ssl_certificates_service_contact", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_contact",
		F:    testSweepSslCertificatesServiceContact,
	})
	resource.AddTestSweepers("alicloud_ssl_certificates_service_pca_cert", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_pca_cert",
		F:    testSweepSslCertificatesServicePcaCert,
	})
	resource.AddTestSweepers("alicloud_ssl_certificates_service_pca_certificate", &resource.Sweeper{
		Name: "alicloud_ssl_certificates_service_pca_certificate",
		F:    testSweepSslCertificatesServicePcaCertificate,
	})
}

// testSweepSslCertificatesServiceCompany sweeps alicloud_ssl_certificates_service_company.
// Company records (issuer organization info) are not cert/domain-bearing, so the
// standard fleet-wipe gate applies: SWEEP_ALL bypasses the tf-prefix and cleans
// every non-NODELETING company.
func testSweepSslCertificatesServiceCompany(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "ListCompanies"
	request := map[string]interface{}{
		"ShowSize":    PageSizeLarge,
		"CurrentPage": 1,
	}
	ids := make([]string, 0)
	for {
		response, err := sslRpcPost(client, "cas", "2020-04-07", action, request, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list ssl certificates service company: %s", err)
			return nil
		}
		result := sslList(response, "$.CompanyList")
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := sslField(item, "CompanyName")
			if !sslSweepNameOk(name) {
				log.Printf("[INFO] Skipping ssl company: %s", name)
				continue
			}
			ids = append(ids, fmt.Sprint(item["CompanyId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, id := range ids {
		log.Printf("[INFO] Deleting ssl company: %s", id)
		sslDelete(client, "cas", "2020-04-07", "DeleteCompany", "company", id, map[string]interface{}{
			"CompanyId": id,
		})
	}
	return nil
}

// testSweepSslCertificatesServiceContact sweeps alicloud_ssl_certificates_service_contact.
// Contact records (DV approver phone/email) are not cert/domain-bearing, so the
// standard fleet-wipe gate applies: SWEEP_ALL bypasses the tf-prefix and cleans
// every non-NODELETING contact.
func testSweepSslCertificatesServiceContact(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "ListContact"
	request := map[string]interface{}{
		"ShowSize":    PageSizeLarge,
		"CurrentPage": 1,
	}
	ids := make([]string, 0)
	for {
		response, err := sslRpcPost(client, "cas", "2020-04-07", action, request, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list ssl certificates service contact: %s", err)
			return nil
		}
		result := sslList(response, "$.ContactList")
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := sslField(item, "Name")
			if !sslSweepNameOk(name) {
				log.Printf("[INFO] Skipping ssl contact: %s", name)
				continue
			}
			ids = append(ids, fmt.Sprint(item["ContactId"]))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, id := range ids {
		log.Printf("[INFO] Deleting ssl contact: %s", id)
		sslDelete(client, "cas", "2020-04-07", "DeleteContact", "contact", id, map[string]interface{}{
			"ContactId": id,
		})
	}
	return nil
}

// testSweepSslCertificatesServicePcaCert sweeps alicloud_ssl_certificates_service_pca_cert.
// ListClientCertificate returns both client and server PCA certs; a tf-prefix match
// (checked always-active, never bypassed by SWEEP_ALL) is the only thing that makes
// a cert eligible — PCA certs back real private PKI, so a fleet wipe must never
// reach a cert the tests did not create (red-line #5).
func testSweepSslCertificatesServicePcaCert(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "ListClientCertificate"
	request := map[string]interface{}{
		"ShowSize":    PageSizeXLarge,
		"CurrentPage": 1,
	}
	ids := make([]string, 0)
	for {
		response, err := sslRpcPost(client, "cas", "2020-06-30", action, request, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list ssl certificates service pca cert: %s", err)
			return nil
		}
		result := sslList(response, "$.CertificateList")
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := sslField(item, "AliasName")
			if name == "" {
				name = sslField(item, "CommonName")
			}
			if !sslCertSweepNameOk(name) {
				log.Printf("[INFO] Skipping ssl pca cert: %s", name)
				continue
			}
			ids = append(ids, fmt.Sprint(item["Identifier"]))
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, id := range ids {
		log.Printf("[INFO] Deleting ssl pca cert: %s", id)
		sslDelete(client, "cas", "2020-06-30", "DeleteClientCertificate", "pca cert", id, map[string]interface{}{
			"Identifier": id,
		})
	}
	return nil
}

// testSweepSslCertificatesServicePcaCertificate sweeps
// alicloud_ssl_certificates_service_pca_certificate. DescribeCACertificateList
// returns root and sub CA certs; the tf-prefix check is always-active (never
// bypassed by SWEEP_ALL) because a CA backs every certificate it issued — deleting
// a real CA is irreversible and cascades onto real certs (red-line #5).
func testSweepSslCertificatesServicePcaCertificate(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "Error getting AliCloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "DescribeCACertificateList"
	request := map[string]interface{}{
		"ShowSize":    PageSizeXLarge,
		"CurrentPage": 1,
	}
	ids := make([]string, 0)
	for {
		response, err := sslRpcPost(client, "cas", "2020-06-30", action, request, 5*time.Minute)
		if err != nil {
			log.Printf("[ERROR] Failed to list ssl certificates service pca certificate: %s", err)
			return nil
		}
		result := sslList(response, "$.CertificateList")
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			name := sslField(item, "Alias")
			if name == "" {
				name = sslField(item, "CommonName")
			}
			if !sslCertSweepNameOk(name) {
				log.Printf("[INFO] Skipping ssl pca certificate: %s", name)
				continue
			}
			ids = append(ids, fmt.Sprint(item["Identifier"]))
		}
		if len(result) < PageSizeXLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	for _, id := range ids {
		log.Printf("[INFO] Deleting ssl pca certificate: %s", id)
		sslDelete(client, "cas", "2020-06-30", "DeleteClientCertificate", "pca certificate", id, map[string]interface{}{
			"Identifier": id,
		})
	}
	return nil
}
