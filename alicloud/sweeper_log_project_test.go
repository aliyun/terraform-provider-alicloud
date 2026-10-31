package alicloud

import (
	"fmt"
	"log"
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This root-package copy of the alicloud_log_project sweeper and its
// testSweepLogProjects / testSweepLogProjectsWithPrefixAndSuffix helpers is
// retained per the migration stock-method-disposition (point 4): the sls
// product subpackage (alicloud/service/sls, external sls_test) now owns the
// canonical copy for the migrated SLS resources, but cross-product root test
// callers — notably alicloud/resource_alicloud_cs_kubernetes_test.go, which
// sweeps the log projects created by the K8s Service — still reference
// testSweepLogProjectsWithPrefixAndSuffix directly in the root alicloud test
// package. The root copy uses the root-package bare helpers
// (sharedClientForRegion, sweepAll, WrapError, IsExpectedErrors) and does not
// import the sls product package, so it introduces no cross-product
// dependency. The behavior is identical to the migrated sls_test copy.
func init() {
	resource.AddTestSweepers("alicloud_log_project", &resource.Sweeper{
		Name: "alicloud_log_project",
		F:    testSweepLogProjects,
	})
}

func testSweepLogProjects(region string) error {
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf-example",
		"tf_example",
		"terraform-example",
		"tf_test_",
		"tf-test-",
		"k8s-log-",
		"dbaudit",
		"kms-log-",
	}
	return testSweepLogProjectsWithPrefixAndSuffix(region, prefixes, []string{})
}

func testSweepLogProjectsWithPrefixAndSuffix(region string, prefixes, suffixes []string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Log Projects: %s", WrapError(err))
	}
	names, _ := raw.([]string)

	for _, v := range names {
		name := v
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				for _, suffix := range suffixes {
					if strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffix)) {
						skip = false
						break
					}
				}
			}
			// Sweep the project which from the k8s cluster
			if skip && strings.HasPrefix(name, "k8s-log-") {
				k8sId := strings.TrimPrefix(name, "k8s-log-")
				raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
					return csClient.DescribeCluster(k8sId)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"ErrorClusterNotFound"}) {
						skip = false
					} else {
						log.Printf("[ERROR] DescribeCluster got an error: %#v", err)
					}
				} else {
					cluster, _ := raw.(cs.ClusterType)
					if strings.HasPrefix(strings.ToLower(cluster.Name), "tf-testacc") || strings.HasPrefix(strings.ToLower(cluster.Name), "tf_testacc") {
						skip = false
					}
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Log Project: %s", name)
				continue
			}
		}
		log.Printf("[INFO] Deleting Log Project: %s", name)
		_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return nil, slsClient.DeleteProject(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Log Project (%s): %s", name, err)
		}
	}
	return nil
}
