package alicloud

import (
	"fmt"
	"log"
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// SLS sub-resource sweepers (project-scoped).
//
// The parent alicloud_log_project sweeper (resource_alicloud_log_project_test.go)
// already handles whole-project deletion. These sweepers clean the data-bearing
// SUB-RESOURCES (logstore / machine_group / alert / dashboard) that live INSIDE
// a project, so that a stale test project's sub-resources are removed even when
// the project itself is protected/managed and cannot be deleted.
//
// 🔴 RED-LINE GUARD (real存 preserve): the prefix check on the PARENT PROJECT
// name is applied UNCONDITIONALLY — even when sweepAll() is true. Managed/infra
// SLS projects (starops-*, jarvis-*, trcn-*, default-cms-*, cloud-resource-*,
// slsaudit-*) do not match the test prefixes, so their sub-resources are NEVER
// touched, regardless of SWEEP_ALL. This prevents silent production data loss.
// NODELETING markers on either the project name or the sub-resource name are
// also honored via protectedSweepVpcName.

// slsSubResourceSweepPrefixes is the canonical test-prefix list for SLS
// sub-resource sweepers. A project must match one of these (case-insensitive)
// for its sub-resources to be swept. This list intentionally mirrors the parent
// alicloud_log_project sweeper's test prefixes (minus the infra-ish k8s/dbaudit/
// kms prefixes, which are preserved at the sub-resource level).
var slsSubResourceSweepPrefixes = []string{
	"tf-testacc",
	"tf_testacc",
	"tf-example",
	"tf_example",
	"terraform-example",
	"tf_test_",
	"tf-test-",
	"tf-test",
	"tf_test",
	"test-",
	"test_",
}

// slsSubResourceProjectShouldSweep decides whether sub-resources inside the
// given SLS project should be swept. The decision is based ONLY on the parent
// project name and is applied even under SWEEP_ALL.
//
// Returns (true, "") when sub-resources inside the project should be swept.
// Returns (false, reason) when the project (and therefore its sub-resources)
// must be skipped.
func slsSubResourceProjectShouldSweep(projectName string) (bool, string) {
	if protectedSweepVpcName(projectName) {
		return false, "NODELETING/do-not-delete marker in project name"
	}
	lname := strings.ToLower(projectName)
	for _, prefix := range slsSubResourceSweepPrefixes {
		if strings.HasPrefix(lname, strings.ToLower(prefix)) {
			return true, ""
		}
	}
	return false, "project name does not match test prefixes (managed/infra project — sub-resources preserved)"
}

// slsSubResourceNameShouldSweep is the per-sub-resource name guard. It only
// blocks sub-resources carrying a NODELETING/do-not-delete marker; within a
// qualifying test project every other sub-resource is swept.
func slsSubResourceNameShouldSweep(name string) (bool, string) {
	if protectedSweepVpcName(name) {
		return false, "NODELETING/do-not-delete marker in sub-resource name"
	}
	return true, ""
}

// slsListSweepableProjects returns the names of SLS projects whose
// sub-resources should be swept, using slsSubResourceProjectShouldSweep.
func slsListSweepableProjects(client *connectivity.AliyunClient) ([]string, error) {
	raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
		return slsClient.ListProject()
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving SLS projects: %s", WrapError(err))
	}
	allProjects, _ := raw.([]string)
	var sweepable []string
	for _, p := range allProjects {
		ok, _ := slsSubResourceProjectShouldSweep(p)
		if ok {
			sweepable = append(sweepable, p)
		} else {
			log.Printf("[INFO] Skipping SLS sub-resources of project %q (managed/non-test)", p)
		}
	}
	return sweepable, nil
}

func init() {
	resource.AddTestSweepers("alicloud_log_store", &resource.Sweeper{
		Name: "alicloud_log_store",
		F:    testSweepLogStores,
	})
	resource.AddTestSweepers("alicloud_log_machine_group", &resource.Sweeper{
		Name: "alicloud_log_machine_group",
		F:    testSweepLogMachineGroups,
	})
	resource.AddTestSweepers("alicloud_log_alert", &resource.Sweeper{
		Name: "alicloud_log_alert",
		F:    testSweepLogAlerts,
	})
	resource.AddTestSweepers("alicloud_log_dashboard", &resource.Sweeper{
		Name: "alicloud_log_dashboard",
		F:    testSweepLogDashboards,
	})
}

// testSweepLogStores sweeps alicloud_log_store (logstore) sub-resources inside
// test-prefix SLS projects. Managed projects' logstores are never touched.
func testSweepLogStores(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	projects, err := slsListSweepableProjects(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, project := range projects {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			return slsClient.ListLogStore(project)
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving LogStores for project %q: %s", project, WrapError(err))
			continue
		}
		storeNames, _ := raw.([]string)
		for _, store := range storeNames {
			if ok, reason := slsSubResourceNameShouldSweep(store); !ok {
				log.Printf("[INFO] Skipping LogStore %q in project %q: %s", store, project, reason)
				continue
			}
			log.Printf("[INFO] Deleting LogStore %q in project %q", store, project)
			_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.DeleteLogStore(project, store)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete LogStore (%s/%s): %s", project, store, err)
			}
		}
	}
	return nil
}

// testSweepLogMachineGroups sweeps alicloud_log_machine_group sub-resources
// inside test-prefix SLS projects.
func testSweepLogMachineGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	projects, err := slsListSweepableProjects(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, project := range projects {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			names, _, gErr := slsClient.ListMachineGroup(project, 0, 100)
			return names, gErr
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving MachineGroups for project %q: %s", project, WrapError(err))
			continue
		}
		groupNames, _ := raw.([]string)
		for _, group := range groupNames {
			if ok, reason := slsSubResourceNameShouldSweep(group); !ok {
				log.Printf("[INFO] Skipping MachineGroup %q in project %q: %s", group, project, reason)
				continue
			}
			log.Printf("[INFO] Deleting MachineGroup %q in project %q", group, project)
			_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.DeleteMachineGroup(project, group)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete MachineGroup (%s/%s): %s", project, group, err)
			}
		}
	}
	return nil
}

// testSweepLogAlerts sweeps alicloud_log_alert sub-resources inside
// test-prefix SLS projects.
func testSweepLogAlerts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	projects, err := slsListSweepableProjects(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, project := range projects {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			alerts, _, _, aErr := slsClient.ListAlert(project, "", "", 0, 100)
			return alerts, aErr
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving Alerts for project %q: %s", project, WrapError(err))
			continue
		}
		alerts, _ := raw.([]*sls.Alert)
		for _, alert := range alerts {
			if alert == nil {
				continue
			}
			name := alert.Name
			if ok, reason := slsSubResourceNameShouldSweep(name); !ok {
				log.Printf("[INFO] Skipping Alert %q in project %q: %s", name, project, reason)
				continue
			}
			log.Printf("[INFO] Deleting Alert %q in project %q", name, project)
			_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.DeleteAlert(project, name)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Alert (%s/%s): %s", project, name, err)
			}
		}
	}
	return nil
}

// testSweepLogDashboards sweeps alicloud_log_dashboard sub-resources inside
// test-prefix SLS projects.
func testSweepLogDashboards(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	projects, err := slsListSweepableProjects(client)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return nil
	}
	for _, project := range projects {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			names, _, _, dErr := slsClient.ListDashboard(project, "", 0, 100)
			return names, dErr
		})
		if err != nil {
			log.Printf("[ERROR] Error retrieving Dashboards for project %q: %s", project, WrapError(err))
			continue
		}
		dashboards, _ := raw.([]string)
		for _, dashboard := range dashboards {
			if ok, reason := slsSubResourceNameShouldSweep(dashboard); !ok {
				log.Printf("[INFO] Skipping Dashboard %q in project %q: %s", dashboard, project, reason)
				continue
			}
			log.Printf("[INFO] Deleting Dashboard %q in project %q", dashboard, project)
			_, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				return nil, slsClient.DeleteDashboard(project, dashboard)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Dashboard (%s/%s): %s", project, dashboard, err)
			}
		}
	}
	return nil
}
