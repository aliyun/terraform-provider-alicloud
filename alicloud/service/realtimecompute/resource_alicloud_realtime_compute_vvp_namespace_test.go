package realtimecompute_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/realtimecompute"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
	// The Flink workspace and namespace names, and the OSS bucket derived from them, only accept
	// lowercase letters, digits and hyphens (foasconsole CreateInstance/CreateNamespace and OSS
	// bucket naming rules), so these tests keep the lowercase tf-testacc prefix instead of
	// emitting names the product rejects.
	name := fmt.Sprintf("tf-testacc%srtcns%d", strings.ReplaceAll(region, "-", ""), rand)
	namespace := fmt.Sprintf("tf-testacc-ns-%d", rand)
	namespaceHa := fmt.Sprintf("tf-testacc-ns-ha-%d", rand)
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
				// The minimal configuration omits ha and uses explicit timeouts: the created namespace
				// must still come back with the service default (false on a single-zone workspace), and
				// the local timeouts block must be accepted without a per-apply change.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceMinimal0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "ha", "false"),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "gmt_create"),
					resource.TestCheckResourceAttrSet(resourceId, "gmt_modified"),
					resource.TestCheckResourceAttr(resourceId, "resource_used.#", "1"),
					resource.TestCheckResourceAttrSet(resourceId, "resource_used.0.cpu"),
					resource.TestCheckResourceAttrSet(resourceId, "resource_used.0.cu"),
					resource.TestCheckResourceAttrSet(resourceId, "resource_used.0.memory_gb"),
					testAccRealtimeComputeVvpCheckNoTags(resourceId),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Import right after the create must reproduce the created state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAliCloudRealtimeComputeVvpNamespaceMinimal0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Import right after re-applying the same configuration must still match.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Growing the namespace is the supported update. A successful shrink is not
				// covered: ModifyNamespaceSpec rejects a target below the sum of the namespace's
				// resource queues (903099), and that rejected shrink is asserted later in this test
				// as the negative case.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceMinimal0(name, namespace, zone, 8, 32, 3, 12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "3"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "12"),
				),
			},
			{
				// Import right after the grow must reproduce the updated state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Shrinking below the capacity the namespace's resource queues already use is
				// rejected by the product, so the update is asserted as a negative case rather
				// than dropped: 903099 is the expected, identifiable error.
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceMinimal0(name, namespace, zone, 8, 32, 2, 8),
				ExpectError: regexp.MustCompile("903099"),
			},
			{
				// An HA namespace needs HA capacity on the parent workspace. The workspace this fixture
				// creates is non-HA (the provider's workspace resource cannot provision HA capacity), so
				// the product refuses the create with 903100. Asserting that refusal keeps the ha=true
				// path exercised and pins the prerequisite instead of pretending the lifecycle ran.
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceHa0(name, namespace, namespaceHa, zone, 8, 32),
				ExpectError: regexp.MustCompile("903100"),
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
				// Import right after the create must reproduce the created state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Second apply of the same configuration: the plan must be empty.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Import right after re-applying the same configuration must still match.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudRealtimeComputeVvpNamespace_intlAccount repeats the ap-southeast-1 lifecycle,
// import and cleanup with the dedicated international test account supplied through
// ALICLOUD_ACCESS_KEY_1 / ALICLOUD_SECRET_KEY_1, so the international-site evidence does not
// depend on the ambient ACC account. It skips when either variable is missing; a SKIP is not
// international-site evidence. The case runs alone: it rewrites the provider credential
// environment for the process and restores every variable when the test ends.
func TestAccAliCloudRealtimeComputeVvpNamespace_intlAccount(t *testing.T) {
	// resource.Test enforces the TF_ACC opt-in only once it runs, which is after the account
	// preflight below. Gate first so a plain `go test` never issues signed cloud requests just
	// because the dedicated credentials happen to be present.
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	ak := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_1"))
	sk := strings.TrimSpace(os.Getenv("ALICLOUD_SECRET_KEY_1"))
	if ak == "" || sk == "" {
		t.Skip("Skipping: set ALICLOUD_ACCESS_KEY_1 and ALICLOUD_SECRET_KEY_1 to run the international-site lifecycle on the dedicated international test account")
	}
	region := "ap-southeast-1"
	// These have to be in place before the provider and SDK clients are built: the provider reads
	// them in Configure, and the read/import/CheckDestroy helpers reuse that configured client.
	// Long-term keys are not paired with an ambient STS token or profile, and the region stays in the
	// environment so no credential-bearing provider block is rendered.
	t.Setenv("ALICLOUD_ACCESS_KEY", ak)
	t.Setenv("ALICLOUD_SECRET_KEY", sk)
	t.Setenv("ALICLOUD_SECURITY_TOKEN", "")
	t.Setenv("ALICLOUD_PROFILE", "")
	t.Setenv("ALIBABA_CLOUD_PROFILE", "")
	t.Setenv("ALICLOUD_REGION", region)

	accountId := testAccRealtimeComputeVvpIntlAccountIdentity(t, region, ak, sk)

	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srtcnsintlacct%d", strings.ReplaceAll(region, "-", ""), rand)
	namespace := fmt.Sprintf("tf-testacc-ns-intl-acct-%d", rand)
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
				// Plan-only first: the provider is configured without creating anything, so the
				// identity assertion proves the configured client is the dedicated account before
				// a single paid resource exists.
				Config:             testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, "", zone, 4, 16, 0, 0),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccRealtimeComputeVvpProviderIdentity(accountId),
				),
			},
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
				// Import right after the create must reproduce the created state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Second apply of the same configuration: the plan must be empty.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceBasic0(name, namespace, zone, 8, 32, 2, 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Import right after re-applying the same configuration must still match.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
				// Import right after the grow must reproduce the updated state; every attribute this
				// resource tracks is returned by DescribeNamespaces, so the comparison needs no
				// ignore list.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
				// Import right after re-applying the same configuration must still match.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				// Import right after the final successful configuration must still match. The side
				// namespace only exists outside Terraform, so the comparison covers the managed
				// default namespace exactly as the previous step left it.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAliCloudRealtimeComputeVvpNamespace_validation0 covers the plan-time input rules: an
// invalid cpu/memory ratio, an explicitly empty resource_spec and a null resource_spec element must
// be rejected while planning (no cloud resource is created), while omitting resource_spec
// altogether or leaving its values unknown must still plan.
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
				// [null] satisfies SizeAtLeast(1) but carries no cpu/memory_gb; the serializer would
				// skip the element and send the request without a spec, so it is rejected at plan time.
				Config:      testAccAliCloudRealtimeComputeVvpNamespaceSpecNullElement(),
				ExpectError: regexp.MustCompile(`resource_spec\[0\] is null`),
			},
			{
				// Values that are unknown at plan time are deferred to apply time, not rejected.
				Config:             testAccAliCloudRealtimeComputeVvpNamespaceSpecUnknown(),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
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

// testAccRealtimeComputeVvpIntlAccountIdentity verifies in this process that the dedicated
// credentials are a working account for the international region and that Realtime Compute answers
// there, before the lifecycle spends anything. Only the account tail and the workspace count are
// logged; the credentials themselves are never printed.
func testAccRealtimeComputeVvpIntlAccountIdentity(t *testing.T, region, ak, sk string) string {
	t.Helper()
	client, err := sdk.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		t.Fatalf("building the STS client for the international test account failed: %s", err)
	}
	identity, err := testAccRealtimeComputeVvpIntlCall(client, requests.POST, "Sts", "2015-04-01", "sts."+region+".aliyuncs.com", "GetCallerIdentity", nil)
	if err != nil {
		t.Fatalf("Sts GetCallerIdentity with the international test account failed: %s", err)
	}
	accountId := strings.TrimSpace(fmt.Sprint(identity["AccountId"]))
	if accountId == "" || accountId == "<nil>" {
		t.Fatalf("Sts GetCallerIdentity returned no AccountId for the international test account")
	}
	tail := accountId
	if len(tail) > 4 {
		tail = tail[len(tail)-4:]
	}
	t.Logf("international test account identity verified via Sts GetCallerIdentity (account tail %s)", tail)

	// Read-only product probe for the region the lifecycle runs in: reachability plus how many
	// workspaces the account already holds there.
	instances, err := testAccRealtimeComputeVvpIntlCall(client, requests.GET, "foasconsole", "2021-10-28", "foasconsole.aliyuncs.com", "DescribeInstances", map[string]string{"Region": region})
	if err != nil {
		t.Fatalf("DescribeInstances in %s with the international test account failed: %s", region, err)
	}
	count := 0
	if items, ok := instances["Instances"].([]interface{}); ok {
		count = len(items)
	}
	t.Logf("Realtime Compute answered in %s; the account currently holds %d workspace(s) there", region, count)
	return accountId
}

// testAccRealtimeComputeVvpIntlCall issues one signed common request with the dedicated
// credentials, so the preflight never falls back to the ambient account. The HTTP method follows
// the action: GetCallerIdentity accepts POST while foasconsole DescribeInstances is GET-only.
func testAccRealtimeComputeVvpIntlCall(client *sdk.Client, method, product, version, domain, action string, params map[string]string) (map[string]interface{}, error) {
	request := requests.NewCommonRequest()
	request.Method = method
	request.Scheme = "https"
	request.Domain = domain
	request.Version = version
	request.Product = product
	request.ApiName = action
	for key, value := range params {
		request.QueryParams[key] = value
	}
	request.TransToAcsRequest()

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal([]byte(response.GetHttpContentString()), &result)
	if !response.IsSuccess() {
		return nil, fmt.Errorf("%s failed with %v %v", action, result["Code"], result["Message"])
	}
	return result, nil
}

// testAccRealtimeComputeVvpProviderIdentity proves the configured provider client is the same
// account the preflight verified, so the HCL resources, reads, import and CheckDestroy all run on
// that account instead of only the environment having been rewritten.
func testAccRealtimeComputeVvpProviderIdentity(expectedAccountId string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return fmt.Errorf("reading the configured provider caller identity failed: %s", err)
		}
		if identity.AccountId != expectedAccountId {
			return fmt.Errorf("the configured provider account does not match the dedicated international test account")
		}
		return nil
	}
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

func testAccAliCloudRealtimeComputeVvpNamespaceSpecNullElement() string {
	return `
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"

  resource_spec = [null]
}`
}

// testAccAliCloudRealtimeComputeVvpNamespaceSpecUnknown keeps the cpu/memory_gb pair unknown at
// plan time by deriving it from a resource that does not exist yet, so the plan-time validator has
// to defer the ratio check to apply time instead of rejecting the configuration.
func testAccAliCloudRealtimeComputeVvpNamespaceSpecUnknown() string {
	return `
resource "alicloud_oss_bucket" "default" {
  bucket = "tf-example-vvp-ns-unknown"
}

resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"

  resource_spec = [
    {
      cpu       = length(alicloud_oss_bucket.default.id)
      memory_gb = 4 * length(alicloud_oss_bucket.default.id)
    }
  ]
}`
}

func testAccAliCloudRealtimeComputeVvpNamespaceSpecOmitted() string {
	return `
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = "f-example-instance"
  namespace   = "tf-example-ns"
}`
}

// testAccAliCloudRealtimeComputeVvpNamespaceMinimal0 is the shared-VPC fixture with a minimal
// namespace block: ha is omitted so the service default is exercised, and the timeouts block uses
// non-default values so the local operation budget is actually configured.
func testAccAliCloudRealtimeComputeVvpNamespaceMinimal0(name, namespace, zone string, instanceCpu, instanceMemoryGb, cpu, memoryGb int) string {
	namespaceBlock := ""
	if namespace != "" {
		namespaceBlock = fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q

  resource_spec = [
    {
      cpu       = %d
      memory_gb = %d
    }
  ]

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}`, namespace, cpu, memoryGb)
	}
	return testAccAliCloudRealtimeComputeVvpNamespaceSharedVpcWorkspace(name, zone, instanceCpu, instanceMemoryGb) + namespaceBlock
}

// testAccAliCloudRealtimeComputeVvpNamespaceHa0 keeps the main namespace at its grown
// specification (3/12) and adds an HA namespace next to it. The non-HA workspace refuses that
// create with 903100, which is what the caller asserts.
func testAccAliCloudRealtimeComputeVvpNamespaceHa0(name, namespace, namespaceHa, zone string, instanceCpu, instanceMemoryGb int) string {
	return testAccAliCloudRealtimeComputeVvpNamespaceSharedVpcWorkspace(name, zone, instanceCpu, instanceMemoryGb) + fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q

  resource_spec = [
    {
      cpu       = 3
      memory_gb = 12
    }
  ]
}

resource "alicloud_realtime_compute_vvp_namespace" "ha" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q
  ha          = true

  resource_spec = [
    {
      cpu       = 1
      memory_gb = 4
    }
  ]
}`, namespace, namespaceHa)
}

// testAccAliCloudRealtimeComputeVvpNamespaceSharedVpcWorkspace renders only the workspace part of
// the shared-VPC fixture, so the fixture variants can append their own namespace blocks.
func testAccAliCloudRealtimeComputeVvpNamespaceSharedVpcWorkspace(name, zone string, instanceCpu, instanceMemoryGb int) string {
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
`, name, zone, zone, instanceCpu, instanceMemoryGb)
}

// testAccRealtimeComputeVvpCheckNoTags accepts the two ways an untagged namespace can be rendered:
// the service omits Tags (the attribute is null and absent from the flat state) or returns an empty
// array (tags.% = 0). A namespace that reports tags fails the check.
func testAccRealtimeComputeVvpCheckNoTags(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("cannot find %s in state", resourceId)
		}
		if count, ok := rs.Primary.Attributes["tags.%"]; ok && count != "0" {
			return fmt.Errorf("%s reports %s tags, expected none", resourceId, count)
		}
		return nil
	}
}

// vvpNamespaceUpdateMock serves the foasconsole actions the update path can call and records the
// ordered actions plus the specification the write request carried. The DescribeNamespaces response
// follows the specification the last ModifyNamespaceSpec wrote, so the refreshed state proves that
// the provider used the service answer rather than the plan.
type vvpNamespaceUpdateMock struct {
	mu           sync.Mutex
	actions      []string
	modifyCpu    string
	modifyMemory string
	currentCpu   int64
	currentMem   int64
}

func (m *vvpNamespaceUpdateMock) snapshot() ([]string, string, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.actions...), m.modifyCpu, m.modifyMemory
}

func (m *vvpNamespaceUpdateMock) spec() (int64, int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentCpu, m.currentMem
}

// vvpNamespaceUpdateTestServer starts a foasconsole stub. delay slows every DescribeNamespaces so a
// case can let the operation budget run out after the write started; onModify runs right after the
// write is accepted so a case can cancel the lifecycle in flight.
func vvpNamespaceUpdateTestServer(t *testing.T, failReads bool, delay time.Duration, onModify func()) (*httptest.Server, *vvpNamespaceUpdateMock) {
	t.Helper()
	mock := &vvpNamespaceUpdateMock{currentCpu: 1, currentMem: 4}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		action := r.FormValue("Action")
		mock.mu.Lock()
		mock.actions = append(mock.actions, action)
		mock.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch action {
		case "ModifyNamespaceSpec":
			mock.mu.Lock()
			mock.modifyCpu = r.FormValue("ModifyNamespaceSpecRequest.ResourceSpec.Cpu")
			mock.modifyMemory = r.FormValue("ModifyNamespaceSpecRequest.ResourceSpec.MemoryGB")
			if cpu, err := strconv.ParseInt(mock.modifyCpu, 10, 64); err == nil {
				mock.currentCpu = cpu
			}
			if memory, err := strconv.ParseInt(mock.modifyMemory, 10, 64); err == nil {
				mock.currentMem = memory
			}
			mock.mu.Unlock()
			if onModify != nil {
				onModify()
			}
			fmt.Fprint(w, `{"RequestId":"mock"}`)
		default:
			if delay > 0 {
				time.Sleep(delay)
			}
			if failReads {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"RequestId":"mock","Code":"InternalError","Message":"mock read failure"}`)
				return
			}
			cpu, memory := mock.spec()
			fmt.Fprintf(w, `{"RequestId":"mock","Namespaces":[{"InstanceId":"inst-unit","Namespace":"demo","Status":"SUCCESS","Ha":false,"GmtCreate":1700000000,"GmtModified":1700000000,"ResourceSpec":{"Cpu":%d,"MemoryGB":%d},"ResourceUsed":{"Cpu":1,"Cu":1,"MemoryGB":4},"Tags":[]}]}`, cpu, memory)
		}
	}))
	t.Cleanup(server.Close)
	return server, mock
}

// vvpNamespaceUpdateTestResource configures a fresh resource against a local foasconsole endpoint.
func vvpNamespaceUpdateTestResource(t *testing.T, endpoint string) (fwresource.Resource, fwschema.Schema) {
	t.Helper()
	t.Setenv("NO_PROXY", "127.0.0.1")
	t.Setenv("no_proxy", "127.0.0.1")
	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").
		SetAccessKeyId("test-access-key-id").
		SetAccessKeySecret("test-access-key-secret"))
	if err != nil {
		t.Fatalf("build credential failed: %v", err)
	}
	endpoints := new(sync.Map)
	config := connectivity.Config{
		AccessKey:            "test-access-key-id",
		SecretKey:            "test-access-key-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build client failed: %v", err)
	}
	// Both the tea RPC layer and the legacy SDK prepend the scheme themselves.
	endpoints.Store("foasconsole", strings.TrimPrefix(endpoint, "http://"))

	res := realtimecompute.NewVvpNamespaceResource()
	configurable, ok := res.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("the resource does not implement the Configure interface")
	}
	var configureResp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	var schemaResp fwresource.SchemaResponse
	res.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	return res, schemaResp.Schema
}

// vvpNamespaceUpdateTestValue builds one state or plan value at the given specification with the
// given timeouts.update.
func vvpNamespaceUpdateTestValue(t *testing.T, sch fwschema.Schema, updateTimeout string, cpu, memoryGb int64) tftypes.Value {
	t.Helper()
	objectType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("resource schema does not describe an object")
	}
	specListType, ok := objectType.AttributeTypes["resource_spec"].(tftypes.List)
	if !ok {
		t.Fatalf("resource_spec is not a list: %T", objectType.AttributeTypes["resource_spec"])
	}
	spec := tftypes.NewValue(specListType, []tftypes.Value{
		tftypes.NewValue(specListType.ElementType, map[string]tftypes.Value{
			"cpu":       tftypes.NewValue(tftypes.Number, cpu),
			"memory_gb": tftypes.NewValue(tftypes.Number, memoryGb),
		}),
	})
	timeouts := tftypes.NewValue(objectType.AttributeTypes["timeouts"], map[string]tftypes.Value{
		"create": tftypes.NewValue(tftypes.String, nil),
		"update": tftypes.NewValue(tftypes.String, updateTimeout),
		"delete": tftypes.NewValue(tftypes.String, nil),
	})
	return tftypes.NewValue(objectType, map[string]tftypes.Value{
		"id":            tftypes.NewValue(tftypes.String, "inst-unit:demo"),
		"gmt_create":    tftypes.NewValue(tftypes.Number, 1700000000),
		"gmt_modified":  tftypes.NewValue(tftypes.Number, 1700000000),
		"ha":            tftypes.NewValue(tftypes.Bool, false),
		"instance_id":   tftypes.NewValue(tftypes.String, "inst-unit"),
		"namespace":     tftypes.NewValue(tftypes.String, "demo"),
		"resource_spec": spec,
		"resource_used": tftypes.NewValue(objectType.AttributeTypes["resource_used"], nil),
		"status":        tftypes.NewValue(tftypes.String, "SUCCESS"),
		"tags":          tftypes.NewValue(objectType.AttributeTypes["tags"], nil),
		"timeouts":      timeouts,
	})
}

// TestUnitRealtimeComputeVvpNamespaceUpdateConfirmationReadBudget drives the public Update request
// path. It pins the single operation budget across the stages that follow the write (the SUCCESS
// wait and the confirmation read), the in-flight cancellation, the exact ModifyNamespaceSpec
// request and the refreshed state: the mock answers with the specification the write carried, so a
// successful case proves the provider stored the service answer rather than the plan.
func TestUnitRealtimeComputeVvpNamespaceUpdateConfirmationReadBudget(t *testing.T) {
	cases := []struct {
		name             string
		updateTimeout    string
		cancelBefore     bool
		cancelAfterWrite bool
		failReads        bool
		readDelay        time.Duration
		planCpu          int64
		planMemoryGb     int64
		wantErrSubstring string
		wantModifies     int
		wantReads        int
		wantModifyCpu    string
		wantModifyMem    string
		wantStateCpu     int64
		wantStateMem     int64
	}{
		{
			name:             "exhausted budget keeps state without a read",
			updateTimeout:    "1ns",
			planCpu:          1,
			planMemoryGb:     4,
			wantErrSubstring: "timeouts.update",
			wantStateCpu:     1,
			wantStateMem:     4,
		},
		{
			name:             "exhausted budget on a cancelled lifecycle keeps state without a read",
			updateTimeout:    "1ns",
			cancelBefore:     true,
			planCpu:          1,
			planMemoryGb:     4,
			wantErrSubstring: "timeouts.update",
			wantStateCpu:     1,
			wantStateMem:     4,
		},
		{
			name:             "confirmation read after the write runs out of the same budget",
			updateTimeout:    "800ms",
			readDelay:        500 * time.Millisecond,
			planCpu:          2,
			planMemoryGb:     8,
			wantErrSubstring: "confirmation read interrupted",
			wantModifies:     1,
			wantReads:        2,
			wantModifyCpu:    "2",
			wantModifyMem:    "8",
			wantStateCpu:     1,
			wantStateMem:     4,
		},
		{
			name:             "cancelled lifecycle interrupts the wait after the write",
			updateTimeout:    "1m",
			cancelAfterWrite: true,
			planCpu:          2,
			planMemoryGb:     8,
			wantErrSubstring: "waiting for SUCCESS",
			wantModifies:     1,
			wantModifyCpu:    "2",
			wantModifyMem:    "8",
			wantStateCpu:     1,
			wantStateMem:     4,
		},
		{
			name:          "remaining budget updates and refreshes the state",
			updateTimeout: "1m",
			planCpu:       2,
			planMemoryGb:  8,
			wantModifies:  1,
			wantReads:     2,
			wantModifyCpu: "2",
			wantModifyMem: "8",
			wantStateCpu:  2,
			wantStateMem:  8,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			var cancel context.CancelFunc
			if tc.cancelBefore || tc.cancelAfterWrite {
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
			}
			if tc.cancelBefore {
				cancel()
			}
			server, mock := vvpNamespaceUpdateTestServer(t, tc.failReads, tc.readDelay, func() {
				if tc.cancelAfterWrite {
					cancel()
				}
			})
			res, sch := vvpNamespaceUpdateTestResource(t, server.URL)
			stateValue := vvpNamespaceUpdateTestValue(t, sch, tc.updateTimeout, 1, 4)
			planValue := vvpNamespaceUpdateTestValue(t, sch, tc.updateTimeout, tc.planCpu, tc.planMemoryGb)

			req := fwresource.UpdateRequest{
				Plan:  tfsdk.Plan{Schema: sch, Raw: planValue},
				State: tfsdk.State{Schema: sch, Raw: stateValue},
			}
			resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: sch, Raw: stateValue}}
			res.Update(ctx, req, resp)

			actions, modifyCpu, modifyMem := mock.snapshot()
			modifies, reads := 0, 0
			for _, action := range actions {
				switch action {
				case "ModifyNamespaceSpec":
					modifies++
				case "DescribeNamespaces":
					reads++
				}
			}
			if modifies != tc.wantModifies || reads != tc.wantReads {
				t.Fatalf("ModifyNamespaceSpec = %d, DescribeNamespaces = %d, want %d/%d (actions: %v)",
					modifies, reads, tc.wantModifies, tc.wantReads, actions)
			}
			if tc.wantModifyCpu != "" && (modifyCpu != tc.wantModifyCpu || modifyMem != tc.wantModifyMem) {
				t.Fatalf("ModifyNamespaceSpec carried cpu=%q memory_gb=%q, want %q/%q", modifyCpu, modifyMem, tc.wantModifyCpu, tc.wantModifyMem)
			}
			if tc.wantErrSubstring == "" {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
				}
				for i, action := range actions {
					want := []string{"ModifyNamespaceSpec", "DescribeNamespaces", "DescribeNamespaces"}
					if i < len(want) && action != want[i] {
						t.Fatalf("action order = %v, want %v", actions, want)
					}
				}
			} else {
				if !resp.Diagnostics.HasError() {
					t.Fatalf("expected an error containing %q, got none", tc.wantErrSubstring)
				}
				detail := resp.Diagnostics.Errors()[0].Summary() + " " + resp.Diagnostics.Errors()[0].Detail()
				if !strings.Contains(detail, tc.wantErrSubstring) {
					t.Fatalf("expected an error containing %q, got %q", tc.wantErrSubstring, detail)
				}
			}

			if resp.State.Raw.IsNull() {
				t.Fatalf("the update dropped the namespace from state")
			}
			var id types.String
			resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
			if id.ValueString() != "inst-unit:demo" {
				t.Fatalf("state id = %q, want inst-unit:demo", id.ValueString())
			}
			var spec types.List
			resp.State.GetAttribute(context.Background(), path.Root("resource_spec"), &spec)
			if spec.IsNull() || len(spec.Elements()) != 1 {
				t.Fatalf("state resource_spec = %v, want one element", spec)
			}
			attrs := spec.Elements()[0].(types.Object).Attributes()
			cpu := attrs["cpu"].(types.Int64).ValueInt64()
			memory := attrs["memory_gb"].(types.Int64).ValueInt64()
			if cpu != tc.wantStateCpu || memory != tc.wantStateMem {
				t.Fatalf("state resource_spec = %d/%d, want %d/%d", cpu, memory, tc.wantStateCpu, tc.wantStateMem)
			}
		})
	}
}

// testAccRealtimeComputeVvpPurchaseHACapacity buys same-city HA capacity for the workspace and waits
// until DescribeInstances reports it. HA namespaces are refused with 903100 on a workspace that has
// no HA capacity, so the HA lifecycle test provisions it first.
func testAccRealtimeComputeVvpPurchaseHACapacity(instanceId string, cpu, memoryGb, haCpu, haMemoryGb int, haVSwitchIds []string) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	resourceSpec, err := json.Marshal(map[string]int{"Cpu": cpu, "MemoryGB": memoryGb})
	if err != nil {
		return fmt.Errorf("encoding the capacity spec %d/%d failed: %s", cpu, memoryGb, err)
	}
	haSpec, err := json.Marshal(map[string]int{"Cpu": haCpu, "MemoryGB": haMemoryGb})
	if err != nil {
		return fmt.Errorf("encoding the HA capacity spec %d/%d failed: %s", haCpu, haMemoryGb, err)
	}
	vswitches, err := json.Marshal(haVSwitchIds)
	if err != nil {
		return fmt.Errorf("encoding the HA vswitch ids failed: %s", err)
	}
	request := map[string]interface{}{
		"Region":         client.RegionId,
		"InstanceId":     instanceId,
		"ResourceSpec":   string(resourceSpec),
		"Ha":             true,
		"HaResourceSpec": string(haSpec),
		"HaVSwitchIds":   string(vswitches),
	}
	if _, err := client.RpcPost("foasconsole", "2021-10-28", "ModifyInstanceSpec", nil, request, true); err != nil {
		return fmt.Errorf("buying HA capacity %d/%d failed: %s", haCpu, haMemoryGb, err)
	}
	deadline := time.Now().Add(15 * time.Minute)
	var lastResponse map[string]interface{}
	var lastErr error
	for {
		response, err := client.RpcGet("foasconsole", "2021-10-28", "DescribeInstances",
			map[string]interface{}{"Region": client.RegionId, "InstanceId": instanceId}, nil)
		lastResponse, lastErr = response, err
		if err == nil {
			if instances, ok := response["Instances"].([]interface{}); ok && len(instances) == 1 {
				instance, _ := instances[0].(map[string]interface{})
				haSpecValue, _ := instance["HaResourceSpec"].(map[string]interface{})
				if haSpecValue != nil &&
					fmt.Sprint(haSpecValue["Cpu"]) == fmt.Sprint(haCpu) &&
					fmt.Sprint(haSpecValue["MemoryGB"]) == fmt.Sprint(haMemoryGb) &&
					fmt.Sprint(instance["ClusterStatus"]) == "RUNNING" {
					return nil
				}
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("waiting for the HA capacity to become %d/%d: last response %v, last error %v",
				haCpu, haMemoryGb, lastResponse, lastErr)
		}
		time.Sleep(15 * time.Second)
	}
}

// TestAccAliCloudRealtimeComputeVvpNamespace_ha0 covers the HA namespace lifecycle: the workspace
// buys same-city HA capacity first (ModifyInstanceSpec Ha/HaResourceSpec/HaVSwitchIds), then the
// namespace is created with ha = true, read back, imported, replaced when ha flips to false and
// finally destroyed with the workspace.
func TestAccAliCloudRealtimeComputeVvpNamespace_ha0(t *testing.T) {
	resourceId := "alicloud_realtime_compute_vvp_namespace.default"
	region := "cn-hangzhou"
	if err := os.Setenv("ALICLOUD_REGION", region); err != nil {
		t.Fatalf("setting ALICLOUD_REGION failed: %s", err)
	}
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srtcnsha%d", strings.ReplaceAll(region, "-", ""), rand)
	namespace := fmt.Sprintf("tf-testacc-ns-ha-%d", rand)
	zone := realtimeComputeTestZone(region)
	var instanceId, haVSwitchId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckRealtimeComputeVvpNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudRealtimeComputeVvpNamespaceHaCapacity0(name, "", zone, 4, 16),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_realtime_compute_vvp_instance.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_vswitch.ha", "id"),
					testAccRealtimeComputeVvpCaptureInstanceId(&instanceId),
					testAccRealtimeComputeVvpCaptureResourceId("alicloud_vswitch.ha", &haVSwitchId),
					testAccRealtimeComputeVvpCheckDefaultNamespace(name, 4, 16),
				),
			},
			{
				PreConfig: func() {
					if err := testAccRealtimeComputeVvpPurchaseCapacity(instanceId, 8, 32); err != nil {
						t.Fatalf("purchasing capacity for %s failed: %s", instanceId, err)
					}
					if err := testAccRealtimeComputeVvpPurchaseHACapacity(instanceId, 8, 32, 2, 8, []string{haVSwitchId}); err != nil {
						t.Fatalf("purchasing HA capacity for %s failed: %s", instanceId, err)
					}
				},
				Config: testAccAliCloudRealtimeComputeVvpNamespaceHaCapacity0(name, namespace, zone, 8, 32),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "ha", "true"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.cpu", "2"),
					resource.TestCheckResourceAttr(resourceId, "resource_spec.0.memory_gb", "8"),
				),
			},
			{
				// Import right after the HA create must reproduce the created state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Flipping ha is a replacement: the same name is destroyed and created without HA.
				Config: testAccAliCloudRealtimeComputeVvpNamespaceHaCapacity0(name, namespace, zone, 8, 32, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "namespace", namespace),
					resource.TestCheckResourceAttr(resourceId, "ha", "false"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
				),
			},
			{
				// Import right after the replacement must reproduce the replaced state.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccAliCloudRealtimeComputeVvpNamespaceHa0 renders the workspace fixture with a backup-AZ
// vSwitch for HA capacity and, when namespace is not empty, an HA namespace next to it. The
// namespace defaults to ha = true; pass ha=false to render the replacement configuration.
func testAccAliCloudRealtimeComputeVvpNamespaceHaCapacity0(name, namespace, zone string, instanceCpu, instanceMemoryGb int, ha ...bool) string {
	namespaceHa := true
	if len(ha) > 0 {
		namespaceHa = ha[0]
	}
	namespaceBlock := ""
	if namespace != "" {
		namespaceBlock = fmt.Sprintf(`
resource "alicloud_realtime_compute_vvp_namespace" "default" {
  instance_id = alicloud_realtime_compute_vvp_instance.default.id
  namespace   = %q
  ha          = %t

  resource_spec = [
    {
      cpu       = 2
      memory_gb = 8
    }
  ]
}`, namespace, namespaceHa)
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

resource "alicloud_vswitch" "ha" {
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = "10.4.20.0/24"
  zone_id      = [for z in data.alicloud_zones.default.zones : z.id if z.id != %s][0]
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
`, name, zone, zone, zone, instanceCpu, instanceMemoryGb, namespaceBlock)
}

// testAccRealtimeComputeVvpCaptureResourceId stores the id of the given resource in target.
func testAccRealtimeComputeVvpCaptureResourceId(resourceName string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("cannot find %s in state", resourceName)
		}
		*target = rs.Primary.ID
		return nil
	}
}

// vvpNamespaceCreateRecoveryTestServer answers CreateNamespace with a 5xx that does not prove the
// request was rejected, while DescribeNamespaces already lists the namespace the request asked for.
// A create that adopts the object it finds is what the recovery path must do.
func vvpNamespaceCreateRecoveryTestServer(t *testing.T) (*httptest.Server, func() []string) {
	t.Helper()
	var mu sync.Mutex
	var actions []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		action := r.FormValue("Action")
		mu.Lock()
		actions = append(actions, action)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if action == "CreateNamespace" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"RequestId":"mock","Code":"InternalError","Message":"mock create failure"}`)
			return
		}
		fmt.Fprint(w, `{"RequestId":"mock","Namespaces":[{"InstanceId":"inst-unit","Namespace":"demo","Status":"SUCCESS","Ha":false,"GmtCreate":1700000000,"GmtModified":1700000000,"ResourceSpec":{"Cpu":2,"MemoryGB":8},"ResourceUsed":{"Cpu":0,"Cu":0,"MemoryGB":0},"Tags":[]}]}`)
	}))
	t.Cleanup(server.Close)
	return server, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), actions...)
	}
}

// vvpNamespaceCreateTestValue builds a create plan for the unit recovery test.
func vvpNamespaceCreateTestValue(t *testing.T, sch fwschema.Schema, createTimeout string, cpu, memoryGb int64) tftypes.Value {
	t.Helper()
	objectType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatalf("resource schema does not describe an object")
	}
	specListType, ok := objectType.AttributeTypes["resource_spec"].(tftypes.List)
	if !ok {
		t.Fatalf("resource_spec is not a list: %T", objectType.AttributeTypes["resource_spec"])
	}
	spec := tftypes.NewValue(specListType, []tftypes.Value{
		tftypes.NewValue(specListType.ElementType, map[string]tftypes.Value{
			"cpu":       tftypes.NewValue(tftypes.Number, cpu),
			"memory_gb": tftypes.NewValue(tftypes.Number, memoryGb),
		}),
	})
	timeouts := tftypes.NewValue(objectType.AttributeTypes["timeouts"], map[string]tftypes.Value{
		"create": tftypes.NewValue(tftypes.String, createTimeout),
		"update": tftypes.NewValue(tftypes.String, nil),
		"delete": tftypes.NewValue(tftypes.String, nil),
	})
	return tftypes.NewValue(objectType, map[string]tftypes.Value{
		"id":            tftypes.NewValue(tftypes.String, nil),
		"gmt_create":    tftypes.NewValue(tftypes.Number, nil),
		"gmt_modified":  tftypes.NewValue(tftypes.Number, nil),
		"ha":            tftypes.NewValue(tftypes.Bool, false),
		"instance_id":   tftypes.NewValue(tftypes.String, "inst-unit"),
		"namespace":     tftypes.NewValue(tftypes.String, "demo"),
		"resource_spec": spec,
		"resource_used": tftypes.NewValue(objectType.AttributeTypes["resource_used"], nil),
		"status":        tftypes.NewValue(objectType.AttributeTypes["status"], nil),
		"tags":          tftypes.NewValue(objectType.AttributeTypes["tags"], nil),
		"timeouts":      timeouts,
	})
}

// TestUnitRealtimeComputeVvpNamespaceCreateAdoptsUncertainResult drives the public Create path
// against a service whose CreateNamespace answers 5xx while the namespace exists: the provider must
// look it up by <instance_id>:<namespace>, adopt it with a warning and keep it in state instead of
// reporting a failure that would leave an untracked namespace.
func TestUnitRealtimeComputeVvpNamespaceCreateAdoptsUncertainResult(t *testing.T) {
	server, observed := vvpNamespaceCreateRecoveryTestServer(t)
	res, sch := vvpNamespaceUpdateTestResource(t, server.URL)
	planValue := vvpNamespaceCreateTestValue(t, sch, "1s", 2, 8)

	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: sch}}
	res.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("create reported an error instead of adopting the namespace: %v", resp.Diagnostics)
	}
	warned := false
	for _, warning := range resp.Diagnostics.Warnings() {
		if strings.Contains(warning.Summary()+warning.Detail(), "create response lost") {
			warned = true
		}
	}
	if !warned {
		t.Fatalf("missing the create-response-lost warning: %v", resp.Diagnostics)
	}
	actions := observed()
	creates, reads := 0, 0
	for _, action := range actions {
		switch action {
		case "CreateNamespace":
			creates++
		case "DescribeNamespaces":
			reads++
		}
	}
	if creates == 0 || reads == 0 {
		t.Fatalf("recovery did not call CreateNamespace and DescribeNamespaces (actions: %v)", actions)
	}
	var id types.String
	resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	if id.ValueString() != "inst-unit:demo" {
		t.Fatalf("state id = %q, want inst-unit:demo", id.ValueString())
	}
	var spec types.List
	resp.State.GetAttribute(context.Background(), path.Root("resource_spec"), &spec)
	if spec.IsNull() || len(spec.Elements()) != 1 {
		t.Fatalf("state resource_spec = %v, want one element", spec)
	}
}
