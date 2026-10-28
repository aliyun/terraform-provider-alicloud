package realtimecompute_test

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudRealtimeComputeVvpNamespace_basic0(t *testing.T) {
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	// Pin Hangzhou instead of inheriting the worker default region: it is the region the
	// product's own VvpInstance test targets and where this lifecycle is known to provision.
	region := "cn-hangzhou"
	if err := os.Setenv("ALICLOUD_REGION", region); err != nil {
		t.Fatalf("setting ALICLOUD_REGION failed: %s", err)
	}
	if zone := os.Getenv("ALICLOUD_ZONE"); zone != "" {
		t.Logf("using the zone from ALICLOUD_ZONE: %s", zone)
	}
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srtcns%d", strings.ReplaceAll(region, "-", ""), rand)
	namespace := fmt.Sprintf("tf-testacc-ns-%d", rand)
	zone := realtimeComputeTestZone(region)
	var instanceId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRealtimeComputeVvpNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				// A pay-as-you-go workspace only finishes provisioning when it is created with a
				// resource spec. Its default namespace then takes that whole spec, so the extra
				// capacity is bought in the next step's PreConfig, before the namespace exists.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, "", zone, 4, 16, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_realtime_compute_vvp_instance.default", "id"),
					testAccRealtimeComputeVvpCaptureInstanceId(&instanceId),
					testAccRealtimeComputeVvpCheckDefaultNamespace(name, 4, 16),
				),
			},
			{
				PreConfig: func() {
					if err := testAccRealtimeComputeVvpPurchaseCapacity(instanceId, 8, 32); err != nil {
						t.Fatalf("purchasing capacity for %s failed: %s", instanceId, err)
					}
				},
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "ha", "false"),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "gmt_create"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Growing the namespace is the supported update. Shrinking is not covered here:
				// ModifyNamespaceSpec rejects a target below the sum of the namespace's resource
				// queues (903099) and this resource does not manage those queues.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, namespace, zone, 8, 32, 3, 12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "3"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "12"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Shrinking below the capacity the namespace's resource queues already use is
				// rejected by the product, so the update is asserted as a negative case rather
				// than dropped: 903099 is the expected, identifiable error.
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, namespace, zone, 8, 32, 2, 8),
				ExpectError: regexp.MustCompile("903099"),
			},
		},
	})
}

// TestAccAliCloudRealtimeComputeVvpNamespace_intl0 covers the international site (ap-southeast-1):
// it prepares the workspace, creates the namespace 2/8, re-applies the same configuration for an
// empty plan and destroys. The CU grow, import and 903099 negative case stay in basic0.
func TestAccAliCloudRealtimeComputeVvpNamespace_intl0(t *testing.T) {
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	region := "ap-southeast-1"
	if err := os.Setenv("ALICLOUD_REGION", region); err != nil {
		t.Fatalf("setting ALICLOUD_REGION failed: %s", err)
	}
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srtcnsintl%d", strings.ReplaceAll(region, "-", ""), rand)
	namespace := fmt.Sprintf("tf-testacc-ns-intl-%d", rand)
	zone := realtimeComputeTestZone(region)
	var instanceId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRealtimeComputeVvpNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, "", zone, 4, 16, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_realtime_compute_vvp_instance.default", "id"),
					testAccRealtimeComputeVvpCaptureInstanceId(&instanceId),
				),
			},
			{
				PreConfig: func() {
					if err := testAccRealtimeComputeVvpPurchaseCapacity(instanceId, 8, 32); err != nil {
						t.Fatalf("purchasing capacity for %s failed: %s", instanceId, err)
					}
				},
				Config: testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "ha", "false"),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Second apply of the same configuration: the plan must be empty.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
		},
	})
}

// TestAccAliCloudRealtimeComputeVvpNamespace_default0 adopts the namespace a workspace creates on
// its own (the default namespace) and adjusts its CU through this resource: import, grow, read
// back, then require an empty second plan. It answers the actual business question — changing the
// CU of the default namespace — instead of only reading the default's initial spec.
func TestAccAliCloudRealtimeComputeVvpNamespace_default0(t *testing.T) {
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	region := "cn-hangzhou"
	if err := os.Setenv("ALICLOUD_REGION", region); err != nil {
		t.Fatalf("setting ALICLOUD_REGION failed: %s", err)
	}
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svvpnsdefault%d", strings.ReplaceAll(region, "-", ""), rand)
	zone := realtimeComputeTestZone(region)
	defaultNamespace := fmt.Sprintf("%s-default", name)
	var instanceId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRealtimeComputeVvpDefaultWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				// The workspace is created alone first, so its own default namespace is the
				// pre-existing object this test adopts instead of one Terraform created.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, "", zone, 4, 16, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_realtime_compute_vvp_instance.default", "id"),
					testAccRealtimeComputeVvpCaptureInstanceId(&instanceId),
					testAccRealtimeComputeVvpCheckDefaultNamespace(name, 4, 16),
				),
			},
			{
				// The workspace already holds 4/16 in its default namespace, so buying 8/32 first
				// leaves capacity for the default to grow into.
				PreConfig: func() {
					if err := testAccRealtimeComputeVvpPurchaseCapacity(instanceId, 16, 64); err != nil {
						t.Fatalf("purchasing capacity for %s failed: %s", instanceId, err)
					}
				},
				Config:            testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: false,
				// Keep the imported default in state so the next step grows it in place instead of
				// trying to create a namespace that already exists.
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return fmt.Sprintf("%s:%s", instanceId, defaultNamespace), nil
				},
			},
			{
				// Pin the value the product actually holds before growing: the workspace expansion
				// bought 8/32 for the workspace, but the default namespace still owns 4/16. Reading
				// the API here is what makes the later "8/32" assertion prove that
				// ModifyNamespaceSpec really grew this namespace instead of the import arriving
				// already grown.
				Config:             testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccRealtimeComputeVvpCheckNamespaceSpec(&instanceId, defaultNamespace, 4, 16),
				),
			},
			{
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", defaultNamespace),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "8"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "32"),
				),
			},
			{
				// Re-applying the same configuration must produce an empty plan.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "8"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "32"),
				),
			},
			{
				// Removing the namespace while it is the workspace's last one must fail with the
				// product's 903054 refusal and must NOT be reported as a successful delete.
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, "", zone, 16, 64, 0, 0),
				ExpectError: regexp.MustCompile("903054|last namespace"),
			},
			{
				// The failed delete must have kept the resource in state: the same configuration plans
				// empty, which it could not do if the previous step had dropped the state entry.
				Config:   testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				PlanOnly: true,
			},
			{
				// Terraform can only delete a namespace that is not the last one, and the destroy at
				// the end of this test has to succeed. Add one untracked namespace through the API
				// (it goes away with the workspace) so the managed default namespace is no longer the
				// workspace's last valid one when Terraform deletes it.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, defaultNamespace, zone, 16, 64, 8, 32),
				Check: resource.ComposeTestCheckFunc(
					testAccRealtimeComputeVvpCreateSideNamespace(&instanceId, fmt.Sprintf("%s-extra", name), 2, 8),
				),
			},
		},
	})
}

// TestAccAliCloudRealtimeComputeVvpNamespace_validation0 covers the plan-time input rules: an
// invalid cpu/memory ratio and an explicitly empty resource_spec must be rejected while planning
// (no cloud resource is created), and omitting resource_spec altogether must still plan.
func TestAccAliCloudRealtimeComputeVvpNamespace_validation0(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceSpecOnly("2", "9"),
				ExpectError: regexp.MustCompile("memory_gb must be 4 times"),
			},
			{
				Config:             testAccAliCloudRealtimeComputeVvpNamespaceSpecOnly("2", "8"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceSpecEmpty(),
				ExpectError: regexp.MustCompile("at least 1"),
			},
			{
				Config:             testAccAliCloudRealtimeComputeVvpNamespaceSpecOmitted(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// testAccRealtimeComputeVvpCaptureInstanceId records the workspace id so a later step's
// PreConfig can buy capacity for it before the namespace is planned.
func testAccRealtimeComputeVvpCaptureInstanceId(target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_realtime_compute_vvp_instance.default"]
		if !ok {
			return fmt.Errorf("cannot find alicloud_realtime_compute_vvp_instance.default in state")
		}
		*target = rs.Primary.ID
		return nil
	}
}

// testAccRealtimeComputeVvpCheckNamespaceSpec reads one namespace straight from the API and
// asserts its allocated CU, so a test can prove what the product holds independently of the
// Terraform state the provider writes.
func testAccRealtimeComputeVvpCheckNamespaceSpec(instanceId *string, namespace string, cpu, memoryGb int) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if *instanceId == "" {
			return fmt.Errorf("workspace id was not captured before checking namespace %s", namespace)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := realtimecompute.DescribeVvpNamespace(client, fmt.Sprintf("%s:%s", *instanceId, namespace))
		if err != nil {
			return fmt.Errorf("describing namespace %s failed: %s", namespace, err)
		}
		if object == nil {
			return fmt.Errorf("namespace %s does not exist", namespace)
		}
		spec, _ := object["ResourceSpec"].(map[string]interface{})
		if spec == nil {
			return fmt.Errorf("namespace %s has no ResourceSpec: %v", namespace, object)
		}
		if got := fmt.Sprint(spec["Cpu"]); got != fmt.Sprint(cpu) {
			return fmt.Errorf("namespace %s cpu = %s, expected %d", namespace, got, cpu)
		}
		if got := fmt.Sprint(spec["MemoryGB"]); got != fmt.Sprint(memoryGb) {
			return fmt.Errorf("namespace %s memory_gb = %s, expected %d", namespace, got, memoryGb)
		}
		return nil
	}
}

// testAccRealtimeComputeVvpCreateSideNamespace creates one namespace straight through the API,
// outside Terraform state. The product refuses to delete a workspace's last namespace (903054), so
// this untracked namespace keeps the managed default namespace deletable during teardown; it is
// removed together with the workspace.
func testAccRealtimeComputeVvpCreateSideNamespace(instanceId *string, namespace string, cpu, memoryGb int) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if *instanceId == "" {
			return fmt.Errorf("workspace id was not captured before creating namespace %s", namespace)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"ClientToken": fmt.Sprintf("TF-CreateNamespace-%d", time.Now().UnixNano()),
			"RegionId":    client.RegionId,
			"CreateNamespaceRequest": map[string]interface{}{
				"InstanceId": *instanceId,
				"Namespace":  namespace,
				"Region":     client.RegionId,
				"ResourceSpec": map[string]interface{}{
					"Cpu":      cpu,
					"MemoryGB": memoryGb,
				},
			},
		}
		if _, err := client.RpcPost("foasconsole", "2019-06-01", "CreateNamespace", nil, request, true); err != nil {
			return fmt.Errorf("creating side namespace %s failed: %s", namespace, err)
		}
		// The product only counts a namespace as valid once it is usable, so wait for SUCCESS
		// before the teardown starts relying on it.
		deadline := time.Now().Add(5 * time.Minute)
		for {
			object, err := realtimecompute.DescribeVvpNamespace(client, fmt.Sprintf("%s:%s", *instanceId, namespace))
			if err == nil && object != nil && fmt.Sprint(object["Status"]) == "SUCCESS" {
				return nil
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("side namespace %s did not reach SUCCESS: last error %v", namespace, err)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// testAccRealtimeComputeVvpCheckDefaultNamespace asserts the fact behind "the default namespace
// CU equals the workspace ceiling": creating the workspace also creates its default namespace,
// which takes the whole purchased spec — so adjusting that namespace is exactly what this
// resource is for. The default namespace is named after the workspace (not after the instance
// id), and it appears shortly after the workspace turns RUNNING, so the check polls the read
// path this resource uses instead of assuming one request-time name.
func testAccRealtimeComputeVvpCheckDefaultNamespace(workspaceName string, cpu, memoryGb int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["alicloud_realtime_compute_vvp_instance.default"]
		if !ok {
			return fmt.Errorf("cannot find alicloud_realtime_compute_vvp_instance.default in state")
		}
		instanceId := rs.Primary.ID
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		candidates := []string{
			fmt.Sprintf("%s-default", workspaceName),
			fmt.Sprintf("%s-default-ha", workspaceName),
		}
		deadline := time.Now().Add(5 * time.Minute)
		for {
			for _, name := range candidates {
				object, err := realtimecompute.DescribeVvpNamespace(client, fmt.Sprintf("%s:%s", instanceId, name))
				if err != nil {
					return fmt.Errorf("describing the default namespace %s failed: %s", name, err)
				}
				if object == nil {
					continue
				}
				spec, _ := object["ResourceSpec"].(map[string]interface{})
				if spec == nil {
					return fmt.Errorf("default namespace %s has no ResourceSpec: %v", name, object)
				}
				if got := fmt.Sprint(spec["Cpu"]); got != fmt.Sprint(cpu) {
					return fmt.Errorf("default namespace %s cpu = %s, expected %d", name, got, cpu)
				}
				if got := fmt.Sprint(spec["MemoryGB"]); got != fmt.Sprint(memoryGb) {
					return fmt.Errorf("default namespace %s memory_gb = %s, expected %d", name, got, memoryGb)
				}
				return nil
			}
			if time.Now().After(deadline) {
				listed, listErr := testAccRealtimeComputeVvpListNamespaces(client, instanceId)
				return fmt.Errorf("workspace %s (name %s) did not get a default namespace named any of %v; namespaces listed by the API: %v (list error: %v)",
					instanceId, workspaceName, candidates, listed, listErr)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

// testAccRealtimeComputeVvpListNamespaces lists the namespace names of a workspace. It is only
// used to explain a failed default-namespace assertion, so a failing call is reported, not fatal.
func testAccRealtimeComputeVvpListNamespaces(client *connectivity.AliyunClient, instanceId string) ([]string, error) {
	response, err := client.RpcGet("foasconsole", "2019-06-01", "DescribeNamespaces", map[string]interface{}{
		"DescribeNamespacesRequest.Region":     client.RegionId,
		"DescribeNamespacesRequest.InstanceId": instanceId,
	}, nil)
	if err != nil {
		return nil, err
	}
	names := []string{}
	if items, ok := response["Namespaces"].([]interface{}); ok {
		for _, item := range items {
			namespace, _ := item.(map[string]interface{})
			if name, ok := namespace["Namespace"].(string); ok {
				names = append(names, name)
			}
		}
	}
	return names, nil
}

// testAccRealtimeComputeVvpPurchaseCapacity buys additional pay-as-you-go capacity for the
// workspace and waits until DescribeInstances reports it.
func testAccRealtimeComputeVvpPurchaseCapacity(instanceId string, cpu, memoryGb int) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	region := client.RegionId
	// ModifyInstanceSpec declares ResourceSpec with the json style, so it has to be sent as
	// an encoded object instead of the dotted parameters used by the flat RPC calls.
	resourceSpec, err := json.Marshal(map[string]int{"Cpu": cpu, "MemoryGB": memoryGb})
	if err != nil {
		return fmt.Errorf("encoding the capacity spec %d/%d failed: %s", cpu, memoryGb, err)
	}
	request := map[string]interface{}{
		"Region":       region,
		"InstanceId":   instanceId,
		"ResourceSpec": string(resourceSpec),
	}
	if _, err := client.RpcPost("foasconsole", "2021-10-28", "ModifyInstanceSpec", nil, request, true); err != nil {
		return fmt.Errorf("buying %d/%d failed: %s", cpu, memoryGb, err)
	}

	deadline := time.Now().Add(10 * time.Minute)
	var lastResponse map[string]interface{}
	var lastErr error
	for {
		response, err := client.RpcGet("foasconsole", "2021-10-28", "DescribeInstances",
			map[string]interface{}{"Region": region, "InstanceId": instanceId}, nil)
		lastResponse, lastErr = response, err
		if err == nil {
			if instances, ok := response["Instances"].([]interface{}); ok && len(instances) == 1 {
				instance, _ := instances[0].(map[string]interface{})
				spec, _ := instance["ResourceSpec"].(map[string]interface{})
				if spec != nil &&
					fmt.Sprint(spec["Cpu"]) == fmt.Sprint(cpu) &&
					fmt.Sprint(spec["MemoryGB"]) == fmt.Sprint(memoryGb) &&
					fmt.Sprint(instance["ClusterStatus"]) == "RUNNING" {
					return nil
				}
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("waiting for the capacity to become %d/%d: last response %v, last error %v",
				cpu, memoryGb, lastResponse, lastErr)
		}
		time.Sleep(10 * time.Second)
	}
}

// realtimeComputeTestZone returns a zone expression usable in HCL. Realtime Compute accepts
// only a subset of a region's zones and alicloud_zones has no filter for this product, so
// known regions use a verified zone while other regions fall back to the first VSwitch zone.
func realtimeComputeTestZone(region string) string {
	if zone := os.Getenv("ALICLOUD_ZONE"); zone != "" {
		return fmt.Sprintf("%q", zone)
	}
	switch region {
	case "cn-beijing":
		return `"cn-beijing-i"`
	case "cn-hangzhou":
		return `"cn-hangzhou-i"`
	case "ap-southeast-1":
		// DescribeSupportedZones reports ap-southeast-1b/1c; the first one is pinned so
		// the lifecycle does not depend on the zone data source ordering.
		return `"ap-southeast-1b"`
	}
	return "data.alicloud_zones.default.zones.0.id"
}

// testAccCheckRealtimeComputeVvpDefaultWorkspaceDestroy verifies cleanup for the default-namespace
// test. The test first creates one namespace outside Terraform state, so the managed default is no
// longer the workspace's last valid one and Terraform really deletes it (DeleteNamespace + gone
// wait); that side namespace is then removed together with the workspace. This check therefore only
// proves the workspace and its dependants are gone — it is not evidence for deleting a namespace on
// its own.
func testAccCheckRealtimeComputeVvpDefaultWorkspaceDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_realtime_compute_vvp_instance" {
			continue
		}
		instanceId := rs.Primary.ID
		deadline := time.Now().Add(5 * time.Minute)
		for {
			response, err := client.RpcGet("foasconsole", "2021-10-28", "DescribeInstances",
				map[string]interface{}{"Region": client.RegionId, "InstanceId": instanceId}, nil)
			if err == nil {
				if items, ok := response["Instances"].([]interface{}); !ok || len(items) == 0 {
					return nil
				}
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("workspace %s still exists after destroy: last error %v", instanceId, err)
			}
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}

func testAccCheckRealtimeComputeVvpNamespaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_realtime_compute_vvp_namespace" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		object, err := realtimecompute.DescribeVvpNamespace(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("describing namespace %s failed: %s", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("namespace %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0 keeps the namespace fixtures inside the
// account's shared default VPC instead of creating a new one per run. The ACC account sits at its
// VPC quota and this resource has no VPC behaviour to verify, so reusing the shared network keeps
// the test focused and avoids failing on QuotaExceeded.Vpc before any namespace code runs.
func testAccAliCloudRealtimeComputeVvpNamespaceSharedVpc0(name, namespace, zone string, instanceCpu, instanceMemoryGb, cpu, memoryGb int) string {
	namespaceBlock := ""
	if namespace != "" {
		namespaceBlock = fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q
  ha          = false

  resource_spec = [
    {
      cpu       = %d
      memory_gb = %d
    }
  ]
}`, namespace, cpu, memoryGb)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = %s
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.bucket
    }
  }
  vvp_instance_name = var.name
  vpc_id            = data.alicloud_vpcs.default.ids.0
  zone_id           = %s
  vswitch_ids       = [data.alicloud_vswitches.default.ids.0]
  payment_type      = "PayAsYouGo"

  resource_spec {
    cpu       = %d
    memory_gb = %d
  }
}
%s
`, name, zone, zone, instanceCpu, instanceMemoryGb, namespaceBlock)
}

func testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, namespace, zone string, instanceCpu, instanceMemoryGb, cpu, memoryGb int) string {
	namespaceBlock := ""
	if namespace != "" {
		namespaceBlock = fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q
  ha          = false

  resource_spec = [
    {
      cpu       = %d
      memory_gb = %d
    }
  ]
}`, namespace, cpu, memoryGb)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = %s
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.bucket
    }
  }
  vvp_instance_name = var.name
  vpc_id            = alicloud_vpc.default.id
  zone_id           = %s
  vswitch_ids       = [alicloud_vswitch.default.id]
  payment_type      = "PayAsYouGo"

  resource_spec {
    cpu       = %d
    memory_gb = %d
  }
}
%s
`, name, zone, zone, instanceCpu, instanceMemoryGb, namespaceBlock)
}

// The three helpers below only exercise plan-time validation: the namespace never reaches the API,
// so no cloud resource is created for them.
func testAccAliCloudRealtimeComputeVvpNamespaceSpecOnly(cpu, memoryGb string) string {
	return fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"

  resource_spec = [
    {
      cpu       = %s
      memory_gb = %s
    }
  ]
}`, cpu, memoryGb)
}

func testAccAliCloudRealtimeComputeVvpNamespaceSpecEmpty() string {
	return `
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"

  resource_spec = []
}`
}

func testAccAliCloudRealtimeComputeVvpNamespaceSpecOmitted() string {
	return `
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"
}`
}
