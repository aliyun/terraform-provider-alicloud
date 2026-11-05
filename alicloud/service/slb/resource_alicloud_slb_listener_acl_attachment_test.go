package slb_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	slbsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	listenerAclAttachmentResourceId = "alicloud_slb_listener_acl_attachment.default"
	listenerAclAttachmentListenerId = "alicloud_slb_listener.default"
	listenerAclAttachmentDsId       = "data.alicloud_slb_listeners.default"
)

// TestAccAliCloudSLBListenerAclAttachment_basic0 walks the whole HTTP
// lifecycle: create, in place single field updates, import after every
// successful step, detach (access control off, listener untouched) and
// re-attach.
func TestAccAliCloudSLBListenerAclAttachment_basic0(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAcl%d", tfacctest.RandIntRange(10000, 99999))
	var createdNonAclAttrs map[string]string
	var listenerId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				// The remote non-ACL baseline and the listener identity are
				// captured before any ACL binding exists.
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "off"),
					testAccCaptureResourceId(listenerAclAttachmentListenerId, &listenerId),
					testAccCaptureListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			{
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentResourceId, "listener_id", listenerAclAttachmentListenerId, "id"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentDsId, "slb_listeners.0.id", listenerAclAttachmentListenerId, "id"),
					testAccCheckListenerNonAclAttrsRemote(listenerAclAttachmentListenerId, name),
					resource.ComposeTestCheckFunc(testAccListenerAttrChecks(name)...),
					testAccCompareListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
			{
				Config:   testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				PlanOnly: true,
			},
			{
				// Only acl_id changes: the listener and every other attribute
				// must stay untouched.
				Config: testAccListenerAclAttachmentConfigHttp(name, "second", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.second", "white"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentDsId, "slb_listeners.0.id", listenerAclAttachmentListenerId, "id"),
					testAccCheckListenerNonAclAttrsRemote(listenerAclAttachmentListenerId, name),
					resource.ComposeTestCheckFunc(testAccListenerAttrChecks(name)...),
					testAccCompareListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.second"),
			{
				Config:   testAccListenerAclAttachmentConfigHttp(name, "second", "white", true),
				PlanOnly: true,
			},
			{
				// Only acl_type changes. The full non ACL attribute set and the
				// remote read are checked again, not just the state values.
				Config: testAccListenerAclAttachmentConfigHttp(name, "second", "black", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.second", "black"),
					testAccCheckListenerNonAclAttrsRemote(listenerAclAttachmentListenerId, name),
					resource.ComposeTestCheckFunc(testAccListenerAttrChecks(name)...),
					testAccCompareListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.second"),
			{
				Config:   testAccListenerAclAttachmentConfigHttp(name, "second", "black", true),
				PlanOnly: true,
			},
			{
				// Detach: access control is disabled and the listener keeps
				// running with all non ACL settings unchanged.
				Config: testAccListenerAclAttachmentConfigHttp(name, "second", "black", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listenerAclAttachmentDsId, "slb_listeners.#", "1"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentDsId, "slb_listeners.0.id", listenerAclAttachmentListenerId, "id"),
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "off"),
					testAccCheckListenerNonAclAttrsRemote(listenerAclAttachmentListenerId, name),
					resource.ComposeTestCheckFunc(testAccListenerAttrChecks(name)...),
					testAccCompareListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			{
				Config:   testAccListenerAclAttachmentConfigHttp(name, "second", "black", false),
				PlanOnly: true,
			},
			{
				// Re-attach after the detach; the successful step is followed by
				// an import like every other successful step.
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "on"),
					testAccCompareListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_multiAcl covers the policy group
// identity when the listener exposes several ACLs: the three segment import is
// rejected instead of guessing, and the four segment import selects the exact
// requested group.
func TestAccAliCloudSLBListenerAclAttachment_multiAcl(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclMulti%d", tfacctest.RandIntRange(10000, 99999))
	var listenerId, secondAclId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					// The first binding must be verified on the listener before
					// any out-of-band change, otherwise an incorrect apply could
					// be masked by the later import steps.
					testAccCheckListenerRemoteAclContains(listenerAclAttachmentListenerId, "alicloud_slb_acl.first"),
					testAccCaptureListenerIdentity(&listenerId, &secondAclId),
				),
			},
			// The first successful apply is imported before the second group is
			// bound outside Terraform.
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
			{
				// Bind a second group outside Terraform so the listener exposes
				// two ACLs while the resource still tracks the first.
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccBindExtraListenerAcl("alicloud_slb_acl.first", "alicloud_slb_acl.second"),
					testAccCheckListenerRemoteAclContains(listenerAclAttachmentListenerId, "alicloud_slb_acl.first"),
					testAccCheckListenerRemoteAclContains(listenerAclAttachmentListenerId, "alicloud_slb_acl.second"),
				),
			},
			{
				Config:   testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				PlanOnly: true,
			},
			{
				ResourceName:      listenerAclAttachmentResourceId,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					listener := s.RootModule().Resources[listenerAclAttachmentListenerId]
					if listener == nil {
						return "", fmt.Errorf("cannot find listener state")
					}
					return listener.Primary.ID, nil
				},
				ExpectError: regexp.MustCompile(`(?i)multiple`),
			},
			{
				ResourceName:      listenerAclAttachmentResourceId,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					if listenerId == "" || secondAclId == "" {
						return "", fmt.Errorf("listener or ACL identity was not captured")
					}
					return listenerId + ":" + secondAclId, nil
				},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected one imported state for the second policy group, got %d", len(states))
					}
					state := states[0]
					want := map[string]string{
						"id":          listenerId,
						"listener_id": listenerId,
						"acl_id":      secondAclId,
						"acl_type":    "white",
						"acl_status":  "on",
					}
					for key, value := range want {
						if got := state.Attributes[key]; got != value {
							return fmt.Errorf("imported state %s = %q, want %q", key, got, value)
						}
					}
					return nil
				},
			},
		},
	})
}

// testAccCaptureListenerIdentity records the listener and the second policy
// group IDs so the later import steps can assert the exact identity instead of
// reusing a single state value.
func testAccCaptureListenerIdentity(listenerId, secondAclId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerAclAttachmentListenerId]
		second := s.RootModule().Resources["alicloud_slb_acl.second"]
		if listener == nil || second == nil {
			return fmt.Errorf("cannot find listener or second ACL state")
		}
		*listenerId = listener.Primary.ID
		*secondAclId = second.Primary.ID
		return nil
	}
}

// testAccCheckListenerRemoteAclContains asserts that the policy group is part of
// the ACLs the listener actually reports, not only of the attachment state.
func testAccCheckListenerRemoteAclContains(listenerResourceId, aclResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		acl := s.RootModule().Resources[aclResource]
		if listener == nil || acl == nil {
			return fmt.Errorf("cannot find listener or ACL state")
		}
		object, err := testAccListenerDescribe(listener.Primary.ID)
		if err != nil {
			return err
		}
		for _, id := range slbsvc.ListenerAclIdsForTest(object) {
			if id == acl.Primary.ID {
				return nil
			}
		}
		return fmt.Errorf("listener %s does not report policy group %s, reported: %v", listener.Primary.ID, acl.Primary.ID, slbsvc.ListenerAclIdsForTest(object))
	}
}

// TestAccAliCloudSLBListenerAclAttachment_tcp covers the TCP dispatch path and
// imports the freshly created resource.
func TestAccAliCloudSLBListenerAclAttachment_tcp(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclTcp%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigTcp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentResourceId, "listener_id", listenerAclAttachmentListenerId, "id"),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "on"),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
			{
				Config:   testAccListenerAclAttachmentConfigTcp(name),
				PlanOnly: true,
			},
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_https covers the HTTPS dispatch path.
func TestAccAliCloudSLBListenerAclAttachment_https(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclHttps%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigHttps(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.acl.0", "white"),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "on"),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.acl.0"),
			{
				Config:   testAccListenerAclAttachmentConfigHttps(name),
				PlanOnly: true,
			},
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_udp covers the UDP dispatch path.
func TestAccAliCloudSLBListenerAclAttachment_udp(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclUdp%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigUdp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "on"),
				),
			},
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_externalDelete deletes the listener
// outside Terraform and verifies that a refresh removes the binding and the
// next apply converges without an error loop.
func TestAccAliCloudSLBListenerAclAttachment_externalDelete(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclExt%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigTcp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					testAccDeleteListenerExternally(listenerAclAttachmentListenerId),
				),
				// The listener was deleted outside Terraform, so the next plan
				// is expected to recreate the listener and the binding.
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccListenerAclAttachmentConfigTcp(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "on"),
				),
			},
			// The successful recreate is followed by an import as well.
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
		},
	})
}

// TestAccAliCloudSLBListenerAclAttachment_negative covers the terminal error
// paths: a listener that does not exist, an ACL that does not exist and an ACL
// that conflicts with the listener IP version.
//
// The listener, load balancer and policy groups are created by a successful step
// first, so every rejected apply has an existing target. A rejected apply never
// runs its Check, so each rejection is followed by a PlanOnly step that plans the
// surviving configuration with -refresh=false: the plan can only be empty when
// the rejection left no attachment object in Terraform state, and a refresh
// cannot hide such an object by dropping it. The same step's PreConfig reads the
// complete remote ACL of the same listener before the next write, so a corrective
// apply can never mask a side effect.
func TestAccAliCloudSLBListenerAclAttachment_negative(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclNeg%d", tfacctest.RandIntRange(10000, 99999))
	var listenerId, loadBalancerId string
	var aclIds []string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// The state assertions below must not be masked by a refresh that would
		// drop a leftover attachment object, so every plan skips the refresh.
		AdditionalCLIOptions: &resource.AdditionalCLIOptions{
			Plan: resource.PlanOptions{NoRefresh: true},
		},
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckListenerAclAttachmentDestroy,
			testAccCheckListenerAclAttachmentScopedResourcesDestroy(&listenerId, &loadBalancerId, &aclIds),
		),
		Steps: []resource.TestStep{
			// A combined ACL value would bind several groups while the resource
			// tracks one, so it is rejected during validation.
			{
				Config:      testAccListenerAclAttachmentConfigCombinedAcl(),
				ExpectError: regexp.MustCompile(`Invalid acl_id`),
			},
			// The listener identity is validated locally for both the
			// configuration and the import path. The patterns use tokens that
			// Terraform does not wrap across lines.
			{
				Config:      testAccListenerAclAttachmentConfigInvalidListenerPort(),
				ExpectError: regexp.MustCompile(`Invalid listener_id`),
			},
			{
				Config:      testAccListenerAclAttachmentConfigInvalidListenerProtocol(),
				ExpectError: regexp.MustCompile(`Invalid listener_id`),
			},
			// The surviving configuration creates the SLB, the policy groups and
			// the TCP listener used by every rejected step. Its Check records the
			// identities that the destroy check has to prove released.
			{
				Config: testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNoListenerAclAttachmentState,
					testAccCaptureResourceId("alicloud_slb_load_balancer.default", &loadBalancerId),
					testAccCaptureResourceId(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerRemoteAcl(listenerAclAttachmentListenerId, "off"),
				),
			},
			{
				Config:      testAccListenerAclAttachmentConfigMissingListener(name),
				ExpectError: regexp.MustCompile(`(?i)(not exist|notfound|invalidloadbalancerid)`),
			},
			// Before any correction: the rejection left no attachment in state and
			// did not touch the listener ACL.
			{
				Config: testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name),
				PreConfig: func() {
					testAccAssertListenerAcl(t, listenerId, "off")
				},
				PlanOnly: true,
			},
			{
				Config:      testAccListenerAclAttachmentConfigMissingAcl(name),
				ExpectError: regexp.MustCompile(`(?i)(aclnotexist|acl.?not.?exist|specified acl)`),
			},
			{
				Config: testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name),
				PreConfig: func() {
					testAccAssertListenerAcl(t, listenerId, "off")
				},
				PlanOnly: true,
			},
			{
				Config:      testAccListenerAclAttachmentConfigIpVersionConflict(name),
				ExpectError: regexp.MustCompile(`(?i)ipversionconflict`),
			},
			{
				Config: testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name),
				PreConfig: func() {
					testAccAssertListenerAcl(t, listenerId, "off")
				},
				PlanOnly: true,
			},
			// The corrective apply runs only after the same listener ACL was proven
			// off and the state was proven free of an attachment. Its Check asserts
			// the same facts again and records every policy group for the destroy
			// check, so the rejected requests left no side effect behind.
			{
				Config: testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNoListenerAclAttachmentState,
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
					testAccCheckListenerRemoteAcl(listenerAclAttachmentListenerId, "off"),
					testAccCaptureListenerAclIds(&aclIds),
				),
			},
		},
	})
}

// testAccAssertListenerAcl asserts the complete remote access control state of
// the listener before the next write. PreConfig has no error return and runs
// before the plan, which is the only hook that can observe the listener before
// the following apply.
func testAccAssertListenerAcl(t *testing.T, listenerId, wantStatus string) {
	t.Helper()
	if listenerId == "" {
		t.Fatal("the listener identity was not captured before the corrective apply")
	}
	object, err := testAccListenerDescribe(listenerId)
	if err != nil {
		t.Fatalf("reading listener %s before the next write: %s", listenerId, err)
	}
	status, _ := object["AclStatus"].(string)
	if !strings.EqualFold(status, wantStatus) {
		t.Fatalf("listener %s access control status before the next write = %q, want %q", listenerId, status, wantStatus)
	}
	if !strings.EqualFold(wantStatus, "off") {
		return
	}
	// Access control never succeeded in this test, so disabling it must leave no
	// bound policy group: the member list is part of the ACL contract, not only
	// the status switch.
	if ids := slbsvc.ListenerAclIdsForTest(object); len(ids) != 0 {
		t.Fatalf("listener %s still reports bound policy groups %v before the next write", listenerId, ids)
	}
}

// testAccListenerAclAttachmentImportStep imports the attachment with the four
// segment ID right after a successful configuration step.
func testAccListenerAclAttachmentImportStep(aclResource string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      listenerAclAttachmentResourceId,
		ImportState:       true,
		ImportStateVerify: true,
		ImportStateIdFunc: func(s *terraform.State) (string, error) {
			listener := s.RootModule().Resources[listenerAclAttachmentListenerId]
			acl := s.RootModule().Resources[aclResource]
			if listener == nil || acl == nil {
				return "", fmt.Errorf("cannot find listener or ACL state for import")
			}
			return listener.Primary.ID + ":" + acl.Primary.ID, nil
		},
	}
}

func testAccCheckListenerAclAttachmentAttrs(aclResource, aclType string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(listenerAclAttachmentResourceId, "acl_id", aclResource, "id"),
		resource.TestCheckResourceAttr(listenerAclAttachmentResourceId, "acl_type", aclType),
		resource.TestCheckResourceAttr(listenerAclAttachmentResourceId, "acl_status", "on"),
		resource.TestCheckResourceAttrSet(listenerAclAttachmentResourceId, "id"),
	)
}

func testAccCheckListenerAclAttachmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_listener_acl_attachment" {
			continue
		}
		object, err := testAccListenerDescribe(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		if status, _ := object["AclStatus"].(string); strings.EqualFold(status, "on") {
			return fmt.Errorf("access control of listener %s is still enabled after destroy", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckNoListenerAclAttachmentState(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "alicloud_slb_listener_acl_attachment" {
			return fmt.Errorf("unexpected attachment state for %s after a rejected request", rs.Primary.ID)
		}
	}
	return nil
}

// testAccCheckListenerRemoteAcl proves the complete remote access control state
// of the listener: the status and, when access control is off, that no policy
// group is still bound.
func testAccCheckListenerRemoteAcl(listenerResourceId, wantStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		if listener == nil {
			return fmt.Errorf("cannot find %s state", listenerResourceId)
		}
		object, err := testAccListenerDescribe(listener.Primary.ID)
		if err != nil {
			return err
		}
		status, _ := object["AclStatus"].(string)
		if !strings.EqualFold(status, wantStatus) {
			return fmt.Errorf("listener %s access control status = %q, want %q", listener.Primary.ID, status, wantStatus)
		}
		if strings.EqualFold(wantStatus, "off") {
			if ids := slbsvc.ListenerAclIdsForTest(object); len(ids) != 0 {
				return fmt.Errorf("listener %s still reports bound policy groups %v while access control is off", listener.Primary.ID, ids)
			}
		}
		return nil
	}
}

func testAccCheckListenerRemoteAclStatus(listenerResourceId, wantStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		if listener == nil {
			return fmt.Errorf("cannot find %s state", listenerResourceId)
		}
		object, err := testAccListenerDescribe(listener.Primary.ID)
		if err != nil {
			return err
		}
		status, _ := object["AclStatus"].(string)
		if !strings.EqualFold(status, wantStatus) {
			return fmt.Errorf("listener %s access control status = %q, want %q", listener.Primary.ID, status, wantStatus)
		}
		return nil
	}
}

// listenerAclAttachmentExpectedNonAclAttrs is the exact non ACL attribute set
// configured by the HTTP test listener. Both the remote read and the data
// source assertions compare against it, so an update can never reset a value
// that only looks non empty.
func listenerAclAttachmentExpectedNonAclAttrs(description string) map[string]string {
	return map[string]string{
		"Bandwidth":              "10",
		"Scheduler":              "wrr",
		"HealthCheck":            "on",
		"HealthCheckURI":         "/cons",
		"HealthCheckConnectPort": "20",
		"HealthyThreshold":       "8",
		"UnhealthyThreshold":     "8",
		"HealthCheckTimeout":     "8",
		"HealthCheckInterval":    "5",
		"HealthCheckHttpCode":    "http_2xx,http_3xx",
		"StickySession":          "on",
		"StickySessionType":      "insert",
		"CookieTimeout":          "86400",
		"RequestTimeout":         "80",
		"IdleTimeout":            "30",
		"Description":            description,
	}
}

func testAccListenerNonAclAttrs(s *terraform.State, listenerResourceId string) (map[string]string, error) {
	listener := s.RootModule().Resources[listenerResourceId]
	if listener == nil {
		return nil, fmt.Errorf("cannot find %s state", listenerResourceId)
	}
	object, err := testAccListenerDescribe(listener.Primary.ID)
	if err != nil {
		return nil, err
	}
	attrs := make(map[string]string, len(listenerAclAttachmentExpectedNonAclAttrs("")))
	for key := range listenerAclAttachmentExpectedNonAclAttrs("") {
		attrs[key] = strings.TrimSpace(fmt.Sprint(object[key]))
	}
	return attrs, nil
}

// testAccCheckListenerNonAclAttrsRemote reads the listener from the API and
// proves that binding or unbinding the ACL never resets any configured non ACL
// setting to a different value.
func testAccCheckListenerNonAclAttrsRemote(listenerResourceId, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		got, err := testAccListenerNonAclAttrs(s, listenerResourceId)
		if err != nil {
			return err
		}
		for key, want := range listenerAclAttachmentExpectedNonAclAttrs(description) {
			if got[key] != want {
				return fmt.Errorf("listener %s %s = %q, want %q", listenerResourceId, key, got[key], want)
			}
		}
		return nil
	}
}

func testAccCaptureResourceId(resourceId string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource := s.RootModule().Resources[resourceId]
		if resource == nil {
			return fmt.Errorf("cannot find %s state", resourceId)
		}
		*id = resource.Primary.ID
		return nil
	}
}

// testAccCaptureListenerAclIds records the policy groups created by the
// negative test so the destroy check can prove each of them is released.
func testAccCaptureListenerAclIds(ids *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var captured []string
		for _, resourceId := range []string{"alicloud_slb_acl.first", "alicloud_slb_acl.second", "alicloud_slb_acl.ipv6"} {
			resource := s.RootModule().Resources[resourceId]
			if resource == nil {
				return fmt.Errorf("cannot find %s state", resourceId)
			}
			captured = append(captured, resource.Primary.ID)
		}
		*ids = captured
		return nil
	}
}

// testAccCheckListenerAclAttachmentScopedResourcesDestroy proves that the
// listener, the load balancer and the policy groups created by this test case
// are actually released. The attachment destroy check alone only proves the
// binding is disabled.
func testAccCheckListenerAclAttachmentScopedResourcesDestroy(listenerId, loadBalancerId *string, aclIds *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		service := slbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
		if *listenerId != "" {
			if _, err := service.DescribeSlbListener(*listenerId); err != nil {
				if !errs.NotFoundError(err) {
					return fmt.Errorf("reading listener %s after destroy: %s", *listenerId, err)
				}
			} else {
				return fmt.Errorf("listener %s still exists after destroy", *listenerId)
			}
		}
		if *loadBalancerId != "" {
			if _, err := service.DescribeSlbLoadBalancer(*loadBalancerId); err != nil {
				if !errs.NotFoundError(err) {
					return fmt.Errorf("reading load balancer %s after destroy: %s", *loadBalancerId, err)
				}
			} else {
				return fmt.Errorf("load balancer %s still exists after destroy", *loadBalancerId)
			}
		}
		for _, aclId := range *aclIds {
			if _, err := service.DescribeSlbAcl(aclId); err != nil {
				if !errs.NotFoundError(err) {
					return fmt.Errorf("reading ACL %s after destroy: %s", aclId, err)
				}
			} else {
				return fmt.Errorf("ACL %s still exists after destroy", aclId)
			}
		}
		return nil
	}
}

func testAccCheckListenerIdUnchanged(listenerResourceId string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		if listener == nil {
			return fmt.Errorf("cannot find %s state", listenerResourceId)
		}
		if *id == "" {
			return fmt.Errorf("listener identity baseline was not captured")
		}
		if listener.Primary.ID != *id {
			return fmt.Errorf("listener identity changed from %q to %q", *id, listener.Primary.ID)
		}
		return nil
	}
}

func testAccCaptureListenerNonAclAttrs(listenerResourceId string, snapshot *map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		got, err := testAccListenerNonAclAttrs(s, listenerResourceId)
		if err != nil {
			return err
		}
		*snapshot = got
		return nil
	}
}

// testAccCompareListenerNonAclAttrs proves the full non ACL attribute set did
// not drift between the first apply and the later update, detach and re-attach
// steps.
func testAccCompareListenerNonAclAttrs(listenerResourceId string, snapshot *map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snapshot == nil || len(*snapshot) == 0 {
			return fmt.Errorf("non ACL attribute snapshot for %s was not captured", listenerResourceId)
		}
		got, err := testAccListenerNonAclAttrs(s, listenerResourceId)
		if err != nil {
			return err
		}
		for key, want := range *snapshot {
			if got[key] != want {
				return fmt.Errorf("listener %s %s changed from %q to %q", listenerResourceId, key, want, got[key])
			}
		}
		return nil
	}
}

func testAccListenerDescribe(listenerId string) (map[string]interface{}, error) {
	service := slbsvc.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	return service.DescribeSlbListener(listenerId)
}

func testAccDeleteListenerExternally(listenerResourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		if listener == nil {
			return fmt.Errorf("cannot find %s state", listenerResourceId)
		}
		parts := strings.Split(listener.Primary.ID, ":")
		if len(parts) != 3 {
			return fmt.Errorf("unexpected listener id %q", listener.Primary.ID)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"LoadBalancerId": parts[0],
			"ListenerPort":   parts[2],
		}
		if _, err := client.RpcPost("Slb", "2014-05-15", "DeleteLoadBalancerListener", nil, request, true); err != nil {
			return fmt.Errorf("deleting listener %s externally: %s", listener.Primary.ID, err)
		}
		return nil
	}
}

// testAccBindExtraListenerAcl binds a second policy group outside Terraform so
// the listener exposes more than one group while the resource tracks the first.
func testAccBindExtraListenerAcl(firstAcl, secondAcl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerAclAttachmentListenerId]
		first := s.RootModule().Resources[firstAcl]
		second := s.RootModule().Resources[secondAcl]
		if listener == nil || first == nil || second == nil {
			return fmt.Errorf("cannot find listener or ACL state")
		}
		parts := strings.Split(listener.Primary.ID, ":")
		if len(parts) != 3 {
			return fmt.Errorf("unexpected listener id %q", listener.Primary.ID)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		action := fmt.Sprintf("SetLoadBalancer%sListenerAttribute", strings.ToUpper(parts[1]))
		request := map[string]interface{}{
			"RegionId":       client.RegionId,
			"LoadBalancerId": parts[0],
			"ListenerPort":   parts[2],
			"AclStatus":      "on",
			"AclId":          first.Primary.ID + "," + second.Primary.ID,
			"AclType":        "white",
		}
		if _, err := client.RpcPost("Slb", "2014-05-15", action, nil, request, true); err != nil {
			return fmt.Errorf("binding extra ACL to listener %s: %s", listener.Primary.ID, err)
		}
		return nil
	}
}

func testAccListenerAttrChecks(description string) []resource.TestCheckFunc {
	attrs := map[string]string{
		"bandwidth":                 "10",
		"scheduler":                 "wrr",
		"health_check":              "on",
		"health_check_uri":          "/cons",
		"health_check_connect_port": "20",
		"healthy_threshold":         "8",
		"unhealthy_threshold":       "8",
		"health_check_timeout":      "8",
		"health_check_interval":     "5",
		"health_check_http_code":    "http_2xx,http_3xx",
		"sticky_session":            "on",
		"sticky_session_type":       "insert",
		"cookie_timeout":            "86400",
		"request_timeout":           "80",
		"idle_timeout":              "30",
		"description":               description,
	}
	checks := make([]resource.TestCheckFunc, 0, len(attrs))
	for key, want := range attrs {
		checks = append(checks, resource.TestCheckResourceAttr(listenerAclAttachmentDsId, "slb_listeners.0."+key, want))
	}
	return checks
}

func testAccListenerAclAttachmentCommonConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name   = var.name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  load_balancer_spec   = "slb.s1.small"
}

resource "alicloud_slb_acl" "first" {
  name       = "${var.name}-first"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.10.0/24"
  }
}

resource "alicloud_slb_acl" "second" {
  name       = "${var.name}-second"
  ip_version = "ipv4"
  entry_list {
    entry = "10.10.11.0/24"
  }
}
`, name)
}

func testAccListenerAclAttachmentConfigHttp(name, aclRef, aclType string, withAttachment bool) string {
	attachment := ""
	if withAttachment {
		attachment = fmt.Sprintf(`
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.%s.id
  acl_type    = "%s"
}
`, aclRef, aclType)
	}
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  scheduler                 = "wrr"
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"
  sticky_session            = "on"
  sticky_session_type       = "insert"
  cookie_timeout            = 86400
  request_timeout           = 80
  idle_timeout              = 30
  description               = var.name

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }

  # Keep the policy groups alive until the listener is gone: SLB refuses to
  # delete an ACL that is still associated with a listener.
  depends_on = [alicloud_slb_acl.first, alicloud_slb_acl.second]
}
%s
data "alicloud_slb_listeners" "default" {
  load_balancer_id = alicloud_slb_listener.default.load_balancer_id
}
`, testAccListenerAclAttachmentCommonConfig(name), attachment)
}

func testAccListenerAclAttachmentConfigTcp(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 22
  frontend_port             = 22
  protocol                  = "tcp"
  bandwidth                 = 10
  health_check_type         = "tcp"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.first.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentCommonConfig(name))
}

func testAccListenerAclAttachmentConfigUdp(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 11
  frontend_port             = 11
  protocol                  = "udp"
  bandwidth                 = 10
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.first.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentCommonConfig(name))
}

func testAccListenerAclAttachmentConfigHttps(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener" "default" {
  load_balancer_id          = alicloud_slb_load_balancer.default.id
  backend_port              = 80
  frontend_port             = 443
  protocol                  = "https"
  server_certificate_id     = alicloud_slb_server_certificate.default.id
  bandwidth                 = 10
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.acl.0.id
  acl_type    = "white"
}
`, resourceSlbHTTPSListenerConfigDependence(name))
}

// testAccListenerAclAttachmentConfigNegativeBase is the state every rejected step
// must leave untouched: the load balancer, three policy groups and the TCP
// listener. Each rejected config adds only the attachment under test, so a
// rejection can neither create nor destroy any other object and the following
// plan can attribute every difference to the attachment.
func testAccListenerAclAttachmentConfigNegativeBase(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_acl" "ipv6" {
  name       = "${var.name}-ipv6"
  ip_version = "ipv6"
  entry_list {
    entry = "2001:db8::/32"
  }
}

resource "alicloud_slb_listener" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  backend_port     = 22
  frontend_port    = 22
  protocol         = "tcp"
  bandwidth        = 10
  health_check_type = "tcp"

  lifecycle {
    ignore_changes = [acl_status, acl_type, acl_id, acl_ids]
  }
}
`, testAccListenerAclAttachmentCommonConfig(name))
}

func testAccListenerAclAttachmentConfigMissingListener(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = "lb-mock:http:80"
  acl_id      = alicloud_slb_acl.first.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentConfigNegativeBase(name))
}

func testAccListenerAclAttachmentConfigMissingAcl(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = "acl-notexist00000"
  acl_type    = "white"
}
`, testAccListenerAclAttachmentConfigNegativeBase(name))
}

func testAccListenerAclAttachmentConfigIpVersionConflictNoAttachment(name string) string {
	return testAccListenerAclAttachmentConfigNegativeBase(name)
}

func testAccListenerAclAttachmentConfigIpVersionConflict(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.ipv6.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentConfigNegativeBase(name))
}

// testAccListenerAclAttachmentConfigCombinedAcl configures a comma separated
// ACL value, which the schema rejects before any request is sent.
func testAccListenerAclAttachmentConfigCombinedAcl() string {
	return `
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = "lb-mock:http:80"
  acl_id      = "acl-notexist1,acl-notexist2"
  acl_type    = "white"
}
`
}

// testAccListenerAclAttachmentConfigInvalidListenerPort uses a port outside the
// contract range so the listener_id validator rejects it locally.
func testAccListenerAclAttachmentConfigInvalidListenerPort() string {
	return `
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = "lb-mock:http:65536"
  acl_id      = "acl-notexist1"
  acl_type    = "white"
}
`
}

// testAccListenerAclAttachmentConfigInvalidListenerProtocol uses an unsupported
// protocol so the listener_id validator rejects it locally.
func testAccListenerAclAttachmentConfigInvalidListenerProtocol() string {
	return `
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = "lb-mock:smtp:80"
  acl_id      = "acl-notexist1"
  acl_type    = "white"
}
`
}

// ---------------------------------------------------------------------------
// Lifecycle scenarios through the public resource factory and the provider
// protocol. The transport is replaced through a test-only hook, so these tests
// exercise plan/state decoding, diagnostics and the response state without
// exporting production symbols.
// ---------------------------------------------------------------------------

const listenerAclAttachmentTestTypeName = "alicloud_slb_listener_acl_attachment"

var listenerAclAttachmentTestType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"id":          tftypes.String,
	"listener_id": tftypes.String,
	"acl_id":      tftypes.String,
	"acl_type":    tftypes.String,
	"acl_status":  tftypes.String,
}}

type listenerAclAttachmentTestProvider struct {
	client *connectivity.AliyunClient
}

func (p *listenerAclAttachmentTestProvider) Metadata(ctx context.Context, req fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *listenerAclAttachmentTestProvider) Schema(ctx context.Context, req fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *listenerAclAttachmentTestProvider) Configure(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *listenerAclAttachmentTestProvider) Resources(ctx context.Context) []func() fwresource.Resource {
	return []func() fwresource.Resource{slbsvc.NewListenerAclAttachmentResource}
}

func (p *listenerAclAttachmentTestProvider) DataSources(ctx context.Context) []func() fwdatasource.DataSource {
	return nil
}

func newListenerAclAttachmentServer(t *testing.T) tfprotov5.ProviderServer {
	t.Helper()
	return newListenerAclAttachmentServerWith(t, &connectivity.AliyunClient{RegionId: "cn-beijing"})
}

func newListenerAclAttachmentServerWith(t *testing.T, client *connectivity.AliyunClient) tfprotov5.ProviderServer {
	t.Helper()
	server := providerserver.NewProtocol5(&listenerAclAttachmentTestProvider{
		client: client,
	})()
	providerType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	config, err := tfprotov5.NewDynamicValue(providerType, tftypes.NewValue(providerType, map[string]tftypes.Value{}))
	if err != nil {
		t.Fatalf("encoding the provider config: %s", err)
	}
	resp, err := server.ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{Config: &config})
	if err != nil {
		t.Fatalf("ConfigureProvider: %s", err)
	}
	listenerAclTestRequireNoError(t, "ConfigureProvider", resp.Diagnostics)
	return server
}

// newListenerAclAttachmentTransportClient builds a real client whose SLB
// endpoint is the local test server, so the lifecycle touches the actual
// request construction, signing and response decoding instead of a stub.
func newListenerAclAttachmentTransportClient(t *testing.T, serverURL string) *connectivity.AliyunClient {
	t.Helper()
	for _, name := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy"} {
		t.Setenv(name, "")
	}
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("unit-test-access-key").
		SetAccessKeySecret("unit-test-access-key-secret"))
	if err != nil {
		t.Fatalf("building the credential: %s", err)
	}
	var endpoints sync.Map
	endpoints.Store("slb", strings.TrimPrefix(serverURL, "http://"))
	client, err := (&connectivity.Config{
		SkipRegionValidation: true,
		RegionId:             "cn-beijing",
		Protocol:             "HTTP",
		AccountType:          "Normal",
		AccessKey:            "unit-test-access-key",
		SecretKey:            "unit-test-access-key-secret",
		Credential:           credential,
		Endpoints:            &endpoints,
		SignVersion:          &sync.Map{},
	}).Client()
	if err != nil {
		t.Fatalf("building the local transport client: %s", err)
	}
	return client
}

// listenerAclRpcStub replaces the SLB transport for one test.
type listenerAclRpcStub struct {
	set       func(request map[string]interface{}) error
	describe  func() (map[string]interface{}, error)
	sets      int
	describes int
	requests  []map[string]interface{}
}

func (s *listenerAclRpcStub) install(t *testing.T) {
	t.Helper()
	original := *slbsvc.ListenerAclAttachmentRPCPostForTest
	*slbsvc.ListenerAclAttachmentRPCPostForTest = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
		s.requests = append(s.requests, request)
		if strings.HasPrefix(action, "Set") {
			s.sets++
			if s.set != nil {
				return nil, s.set(request)
			}
			return map[string]interface{}{}, nil
		}
		s.describes++
		if s.describe != nil {
			return s.describe()
		}
		return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
	}
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentRPCPostForTest = original })
}

func listenerAclTestObject(t *testing.T, attrs map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	full := make(map[string]tftypes.Value, len(listenerAclAttachmentTestType.AttributeTypes))
	for name, attrType := range listenerAclAttachmentTestType.AttributeTypes {
		if v, ok := attrs[name]; ok {
			full[name] = v
			continue
		}
		full[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(listenerAclAttachmentTestType, full)
}

func listenerAclTestDynamic(t *testing.T, v tftypes.Value) *tfprotov5.DynamicValue {
	t.Helper()
	dv, err := tfprotov5.NewDynamicValue(listenerAclAttachmentTestType, v)
	if err != nil {
		t.Fatalf("encoding the attachment value: %s", err)
	}
	return &dv
}

func listenerAclTestNull(t *testing.T) *tfprotov5.DynamicValue {
	t.Helper()
	return listenerAclTestDynamic(t, tftypes.NewValue(listenerAclAttachmentTestType, nil))
}

func listenerAclTestAttrs(t *testing.T, dv *tfprotov5.DynamicValue) map[string]tftypes.Value {
	t.Helper()
	if dv == nil {
		t.Fatal("the response carried no state")
	}
	val, err := dv.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the attachment state: %s", err)
	}
	attrs := map[string]tftypes.Value{}
	if err := val.As(&attrs); err != nil {
		t.Fatalf("reading the attachment state: %s", err)
	}
	return attrs
}

func listenerAclTestString(t *testing.T, v tftypes.Value) string {
	t.Helper()
	var out string
	if err := v.As(&out); err != nil {
		t.Fatalf("reading a string attribute: %s", err)
	}
	return out
}

func listenerAclTestHasError(diags []*tfprotov5.Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			return true
		}
	}
	return false
}

func listenerAclTestRequireNoError(t *testing.T, step string, diags []*tfprotov5.Diagnostic) {
	t.Helper()
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("%s: %s: %s", step, d.Summary, d.Detail)
		}
	}
}

func listenerAclTestPlanCreate(t *testing.T, server tfprotov5.ProviderServer, listenerId, aclId, aclType string) (*tfprotov5.PlanResourceChangeResponse, *tfprotov5.DynamicValue) {
	t.Helper()
	config := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"listener_id": tftypes.NewValue(tftypes.String, listenerId),
		"acl_id":      tftypes.NewValue(tftypes.String, aclId),
		"acl_type":    tftypes.NewValue(tftypes.String, aclType),
	}))
	plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         listenerAclAttachmentTestTypeName,
		PriorState:       listenerAclTestNull(t),
		ProposedNewState: config,
		Config:           config,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	listenerAclTestRequireNoError(t, "PlanResourceChange (create)", plan.Diagnostics)
	return plan, config
}

func listenerAclTestApplyCreate(t *testing.T, server tfprotov5.ProviderServer, plan *tfprotov5.PlanResourceChangeResponse, config *tfprotov5.DynamicValue) *tfprotov5.ApplyResourceChangeResponse {
	t.Helper()
	applied, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       listenerAclAttachmentTestTypeName,
		PriorState:     listenerAclTestNull(t),
		PlannedState:   plan.PlannedState,
		Config:         config,
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	return applied
}

func TestUnitListenerAclAttachmentProtocolCreate(t *testing.T) {
	stub := &listenerAclRpcStub{}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
	applied := listenerAclTestApplyCreate(t, server, plan, config)
	listenerAclTestRequireNoError(t, "ApplyResourceChange (create)", applied.Diagnostics)

	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["id"]); got != "lb-123:http:80" {
		t.Fatalf("id = %q", got)
	}
	if got := listenerAclTestString(t, attrs["listener_id"]); got != "lb-123:http:80" {
		t.Fatalf("listener_id = %q", got)
	}
	if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
		t.Fatalf("acl_id = %q", got)
	}
	if got := listenerAclTestString(t, attrs["acl_type"]); got != "white" {
		t.Fatalf("acl_type = %q", got)
	}
	if got := listenerAclTestString(t, attrs["acl_status"]); got != "on" {
		t.Fatalf("acl_status = %q", got)
	}
	if stub.sets != 1 {
		t.Fatalf("expected one write, got %d", stub.sets)
	}
}

func TestUnitListenerAclAttachmentProtocolCreateFirstReadFailure(t *testing.T) {
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			return nil, errors.New("read exploded")
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
	applied := listenerAclTestApplyCreate(t, server, plan, config)
	listenerAclTestRequireNoError(t, "ApplyResourceChange (create, read failure)", applied.Diagnostics)

	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["id"]); got != "lb-123:http:80" {
		t.Fatalf("the identity must survive a failed first read, got id %q", got)
	}
	if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
		t.Fatalf("acl_id = %q", got)
	}
	warned := false
	for _, d := range applied.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityWarning {
			warned = true
		}
	}
	if !warned {
		t.Fatal("expected a warning about the failed first read")
	}
}

func TestUnitListenerAclAttachmentProtocolUnknownWrite(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 200 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })

	transportErr := errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)

	t.Run("confirmed readback keeps the binding", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			set: func(map[string]interface{}) error { return transportErr },
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
		applied := listenerAclTestApplyCreate(t, server, plan, config)
		listenerAclTestRequireNoError(t, "ApplyResourceChange (unknown write)", applied.Diagnostics)
		attrs := listenerAclTestAttrs(t, applied.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
			t.Fatalf("confirmed binding was not kept, acl_id = %q", got)
		}
	})

	t.Run("missing type does not confirm", func(t *testing.T) {
		previous := *slbsvc.ListenerAclAttachmentConfirmTimeoutForTest
		*slbsvc.ListenerAclAttachmentConfirmTimeoutForTest = 100 * time.Millisecond
		t.Cleanup(func() { *slbsvc.ListenerAclAttachmentConfirmTimeoutForTest = previous })

		stub := &listenerAclRpcStub{
			set: func(map[string]interface{}) error { return transportErr },
			describe: func() (map[string]interface{}, error) {
				// The Describe contract requires AclType while access control is
				// on, so an absent type cannot prove the planned value.
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclId": "acl-1"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
		applied := listenerAclTestApplyCreate(t, server, plan, config)
		if !listenerAclTestHasError(applied.Diagnostics) {
			t.Fatal("expected an error when the readback cannot prove the type")
		}
		val, err := applied.NewState.Unmarshal(listenerAclAttachmentTestType)
		if err != nil {
			t.Fatalf("decoding the state: %s", err)
		}
		if !val.IsNull() {
			t.Fatalf("no state may be written for an unconfirmed write: %s", val)
		}
	})
}

func TestUnitListenerAclAttachmentProtocolImportNonMember(t *testing.T) {
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			// The requested policy group is not part of the listener ACLs.
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-other"}, nil
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
		TypeName: listenerAclAttachmentTestTypeName,
		ID:       "lb-123:http:80:acl-missing",
	})
	if err != nil {
		t.Fatalf("ImportResourceState: %s", err)
	}
	listenerAclTestRequireNoError(t, "ImportResourceState", imported.Diagnostics)
	if len(imported.ImportedResources) != 1 {
		t.Fatalf("expected one imported resource, got %d", len(imported.ImportedResources))
	}
	read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     listenerAclAttachmentTestTypeName,
		CurrentState: imported.ImportedResources[0].State,
	})
	if err != nil {
		t.Fatalf("ReadResource after import: %s", err)
	}
	listenerAclTestRequireNoError(t, "ReadResource after import", read.Diagnostics)
	val, err := read.NewState.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the read state: %s", err)
	}
	if !val.IsNull() {
		t.Fatalf("a non-member import must clear the state, got %s", val)
	}
}

func TestUnitListenerAclAttachmentProtocolUpdateStaleType(t *testing.T) {
	previousInterval := *slbsvc.ListenerAclAttachmentPollIntervalForTest
	*slbsvc.ListenerAclAttachmentPollIntervalForTest = 10 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentPollIntervalForTest = previousInterval })

	describes := 0
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			describes++
			if describes <= 2 {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
			}
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "black", "AclId": "acl-1"}, nil
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	prior := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "white"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))
	config := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "black"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))
	plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         listenerAclAttachmentTestTypeName,
		PriorState:       prior,
		ProposedNewState: config,
		Config:           config,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (update): %s", err)
	}
	listenerAclTestRequireNoError(t, "PlanResourceChange (update)", plan.Diagnostics)

	applied, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       listenerAclAttachmentTestTypeName,
		PriorState:     prior,
		PlannedState:   plan.PlannedState,
		Config:         config,
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (update): %s", err)
	}
	listenerAclTestRequireNoError(t, "ApplyResourceChange (update)", applied.Diagnostics)
	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["acl_type"]); got != "black" {
		t.Fatalf("the stale type was stored: %q", got)
	}
	if stub.sets != 1 {
		t.Fatalf("expected one update write, got %d", stub.sets)
	}
}

// The Read and Delete not-found paths, the import forms and the single field
// updates are exercised through the same public factory and protocol entry as
// the rest of the lifecycle.
func TestUnitListenerAclAttachmentProtocolReadAndDestroyMissingListener(t *testing.T) {
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			return nil, &tea.SDKError{StatusCode: tea.Int(404), Code: tea.String("InvalidLoadBalancerId.NotFound"), Message: tea.String("The specified resource does not exist."), Data: tea.String("{}")}
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	prior := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "white"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))

	read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     listenerAclAttachmentTestTypeName,
		CurrentState: prior,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	listenerAclTestRequireNoError(t, "ReadResource", read.Diagnostics)
	val, err := read.NewState.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the read state: %s", err)
	}
	if !val.IsNull() {
		t.Fatalf("a missing listener must clear the state, got %s", val)
	}

	plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         listenerAclAttachmentTestTypeName,
		PriorState:       prior,
		ProposedNewState: listenerAclTestNull(t),
		Config:           listenerAclTestNull(t),
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (destroy): %s", err)
	}
	listenerAclTestRequireNoError(t, "PlanResourceChange (destroy)", plan.Diagnostics)
	destroyed, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       listenerAclAttachmentTestTypeName,
		PriorState:     prior,
		PlannedState:   plan.PlannedState,
		Config:         listenerAclTestNull(t),
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (destroy): %s", err)
	}
	listenerAclTestRequireNoError(t, "ApplyResourceChange (destroy)", destroyed.Diagnostics)
	val, err = destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the destroyed state: %s", err)
	}
	if !val.IsNull() {
		t.Fatalf("destroy of a missing listener must not fail, got %s", val)
	}
}

// Only a decodable SLB business error proves that the listener is gone. A
// gateway 404 without a business body, an empty business code and an enabled
// listener that does not report its members must keep the state on Read and must
// fail Delete instead of reporting a successful unbind. The decodable 404 is the
// control case that still releases the state.
func TestUnitListenerAclAttachmentProtocolGoneClassification(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 300 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })
	previousInterval := *slbsvc.ListenerAclAttachmentPollIntervalForTest
	*slbsvc.ListenerAclAttachmentPollIntervalForTest = 10 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentPollIntervalForTest = previousInterval })

	prior := func(t *testing.T) *tfprotov5.DynamicValue {
		t.Helper()
		return listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
			"acl_type":    tftypes.NewValue(tftypes.String, "white"),
			"acl_status":  tftypes.NewValue(tftypes.String, "on"),
		}))
	}

	cases := []struct {
		name   string
		derive func() (map[string]interface{}, error)
		gone   bool
	}{
		{
			name: "gateway 404 without a business body",
			derive: func() (map[string]interface{}, error) {
				return nil, fmt.Errorf("%w: status 404: unexpected EOF", connectivity.ErrRpcRejectedResponse)
			},
		},
		{
			name: "empty business code on a 403",
			derive: func() (map[string]interface{}, error) {
				return nil, &tea.SDKError{StatusCode: tea.Int(403), Code: tea.String(""), Message: tea.String("Forbidden"), Data: tea.String("{}")}
			},
		},
		{
			name: "enabled listener without a reported member list",
			derive: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white"}, nil
			},
		},
		{
			name: "missing listener proven by the business code",
			derive: func() (map[string]interface{}, error) {
				return nil, &tea.SDKError{StatusCode: tea.Int(404), Code: tea.String("InvalidLoadBalancerId.NotFound"), Message: tea.String("The specified resource does not exist."), Data: tea.String("{}")}
			},
			gone: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run("read", func(t *testing.T) {
				stub := &listenerAclRpcStub{describe: tc.derive}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
					TypeName:     listenerAclAttachmentTestTypeName,
					CurrentState: prior(t),
				})
				if err != nil {
					t.Fatalf("ReadResource: %s", err)
				}
				listenerAclTestRequireNoError(t, "ReadResource", read.Diagnostics)
				value, err := read.NewState.Unmarshal(listenerAclAttachmentTestType)
				if err != nil {
					t.Fatalf("decoding the read state: %s", err)
				}
				if tc.gone {
					if !value.IsNull() {
						t.Fatalf("a proven missing listener must clear the state, got %s", value)
					}
					return
				}
				if value.IsNull() {
					t.Fatal("an unproven missing listener must keep the state")
				}
				if got := listenerAclTestString(t, listenerAclTestAttrs(t, read.NewState)["acl_id"]); got != "acl-1" {
					t.Fatalf("acl_id = %q, want acl-1", got)
				}
			})

			t.Run("destroy", func(t *testing.T) {
				stub := &listenerAclRpcStub{describe: tc.derive}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				state := prior(t)
				plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
					TypeName:         listenerAclAttachmentTestTypeName,
					PriorState:       state,
					ProposedNewState: listenerAclTestNull(t),
					Config:           listenerAclTestNull(t),
				})
				if err != nil {
					t.Fatalf("PlanResourceChange (destroy): %s", err)
				}
				listenerAclTestRequireNoError(t, "PlanResourceChange (destroy)", plan.Diagnostics)
				destroyed, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
					TypeName:       listenerAclAttachmentTestTypeName,
					PriorState:     state,
					PlannedState:   plan.PlannedState,
					Config:         listenerAclTestNull(t),
					PlannedPrivate: plan.PlannedPrivate,
				})
				if err != nil {
					t.Fatalf("ApplyResourceChange (destroy): %s", err)
				}
				if tc.gone {
					listenerAclTestRequireNoError(t, "ApplyResourceChange (destroy of a missing listener)", destroyed.Diagnostics)
					value, err := destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
					if err != nil {
						t.Fatalf("decoding the destroyed state: %s", err)
					}
					if !value.IsNull() {
						t.Fatalf("destroy of a proven missing listener must not fail, got %s", value)
					}
					return
				}
				if !listenerAclTestHasError(destroyed.Diagnostics) {
					t.Fatal("an unproven missing listener must fail the destroy instead of reporting success")
				}
				value, err := destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
				if err != nil {
					t.Fatalf("decoding the destroyed state: %s", err)
				}
				if value.IsNull() {
					t.Fatal("an unconfirmed unbinding must keep the state")
				}
			})
		})
	}
}

func TestUnitListenerAclAttachmentProtocolImportForms(t *testing.T) {
	importAndRead := func(t *testing.T, describe func() (map[string]interface{}, error), id string) (*tfprotov5.DynamicValue, []*tfprotov5.Diagnostic) {
		t.Helper()
		stub := &listenerAclRpcStub{describe: describe}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
			TypeName: listenerAclAttachmentTestTypeName,
			ID:       id,
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		listenerAclTestRequireNoError(t, "ImportResourceState", imported.Diagnostics)
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     listenerAclAttachmentTestTypeName,
			CurrentState: imported.ImportedResources[0].State,
		})
		if err != nil {
			t.Fatalf("ReadResource: %s", err)
		}
		return read.NewState, read.Diagnostics
	}

	t.Run("three segments uses the only bound group", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-9"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
			TypeName: listenerAclAttachmentTestTypeName,
			ID:       "lb-123:http:80",
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		listenerAclTestRequireNoError(t, "ImportResourceState", imported.Diagnostics)
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     listenerAclAttachmentTestTypeName,
			CurrentState: imported.ImportedResources[0].State,
		})
		if err != nil {
			t.Fatalf("ReadResource: %s", err)
		}
		listenerAclTestRequireNoError(t, "ReadResource", read.Diagnostics)
		attrs := listenerAclTestAttrs(t, read.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-9" {
			t.Fatalf("acl_id = %q, want acl-9", got)
		}
	})

	// An enabled listener always reports its members. A response that omits them
	// has not reported the membership, so an import must not adopt it and must
	// not conclude that nothing is bound.
	t.Run("three segments without reported members cannot confirm", func(t *testing.T) {
		state, diagnostics := importAndRead(t, func() (map[string]interface{}, error) {
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white"}, nil
		}, "lb-123:http:80")
		if !listenerAclTestHasError(diagnostics) {
			t.Fatalf("an unreported member list must not confirm the import, state: %s", state)
		}
		if !listenerAclTestDiagnosticsContain(diagnostics, "cannot confirm the imported binding") {
			t.Fatalf("the diagnostic did not report the pending confirmation: %v", diagnostics)
		}
	})

	// The control case: an explicitly empty member list is a reported membership
	// that proves nothing is bound, so the import clears the state.
	t.Run("three segments with an empty reported member list clears the state", func(t *testing.T) {
		state, diagnostics := importAndRead(t, func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"ListenerPort": 80,
				"AclStatus":    "on",
				"AclType":      "white",
				"AclIds":       map[string]interface{}{"AclId": []interface{}{}},
			}, nil
		}, "lb-123:http:80")
		listenerAclTestRequireNoError(t, "ReadResource (empty member list)", diagnostics)
		value, err := state.Unmarshal(listenerAclAttachmentTestType)
		if err != nil {
			t.Fatalf("decoding the read state: %s", err)
		}
		if !value.IsNull() {
			t.Fatalf("a listener without a bound group must clear the state, got %s", value)
		}
	})

	t.Run("three segments with several groups asks for the explicit selection", func(t *testing.T) {
		_, diagnostics := importAndRead(t, func() (map[string]interface{}, error) {
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-a,acl-b"}, nil
		}, "lb-123:http:80")
		if !listenerAclTestHasError(diagnostics) {
			t.Fatal("a three segment import with several groups must fail instead of guessing")
		}
		if !listenerAclTestDiagnosticsContain(diagnostics, "multiple") {
			t.Fatalf("the diagnostic did not ask for the explicit selection: %v", diagnostics)
		}
	})

	t.Run("four segments keeps a non-first group", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-a,acl-b,acl-c"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
			TypeName: listenerAclAttachmentTestTypeName,
			ID:       "lb-123:http:80:acl-b",
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		listenerAclTestRequireNoError(t, "ImportResourceState", imported.Diagnostics)
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     listenerAclAttachmentTestTypeName,
			CurrentState: imported.ImportedResources[0].State,
		})
		if err != nil {
			t.Fatalf("ReadResource: %s", err)
		}
		listenerAclTestRequireNoError(t, "ReadResource", read.Diagnostics)
		attrs := listenerAclTestAttrs(t, read.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-b" {
			t.Fatalf("acl_id = %q, want acl-b", got)
		}
	})

	// The same group may be reported in both response fields and repeated inside
	// the list. The identity extraction has to deduplicate it, otherwise a single
	// binding would look like an ambiguous multi group response.
	t.Run("duplicate response fields keep a single group", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{
					"ListenerPort": 80,
					"AclStatus":    "on",
					"AclType":      "white",
					"AclId":        "acl-a",
					"AclIds":       map[string]interface{}{"AclId": []interface{}{"acl-a", "acl-a"}},
				}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
			TypeName: listenerAclAttachmentTestTypeName,
			ID:       "lb-123:http:80",
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     listenerAclAttachmentTestTypeName,
			CurrentState: imported.ImportedResources[0].State,
		})
		if err != nil {
			t.Fatalf("ReadResource: %s", err)
		}
		listenerAclTestRequireNoError(t, "ReadResource", read.Diagnostics)
		attrs := listenerAclTestAttrs(t, read.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-a" {
			t.Fatalf("acl_id = %q, want acl-a", got)
		}
	})

	t.Run("duplicate response fields still report several groups", func(t *testing.T) {
		_, diagnostics := importAndRead(t, func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"ListenerPort": 80,
				"AclStatus":    "on",
				"AclType":      "white",
				"AclId":        "acl-a,acl-b",
				"AclIds":       map[string]interface{}{"AclId": []interface{}{"acl-b", "acl-a", "acl-a"}},
			}, nil
		}, "lb-123:http:80")
		if !listenerAclTestHasError(diagnostics) {
			t.Fatal("two distinct groups must still require the explicit selection after deduplication")
		}
	})

	// A gateway may serialise an unreported field as an explicit JSON null. That
	// is the same shape as an absent key: it must not turn a readable member list
	// into a malformed one.
	t.Run("explicit null members behave as absent", func(t *testing.T) {
		state, diagnostics := importAndRead(t, func() (map[string]interface{}, error) {
			return map[string]interface{}{
				"ListenerPort": 80,
				"AclStatus":    "on",
				"AclType":      "white",
				"AclId":        "acl-9",
				"AclIds":       nil,
			}, nil
		}, "lb-123:http:80")
		listenerAclTestRequireNoError(t, "ReadResource (null members)", diagnostics)
		attrs := listenerAclTestAttrs(t, state)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-9" {
			t.Fatalf("acl_id = %q, want acl-9", got)
		}
	})

	// The import ID is the public entry point, so each malformed form must be
	// rejected locally instead of reaching Describe with a truncated identity.
	t.Run("invalid import IDs are rejected", func(t *testing.T) {
		for _, id := range []string{
			"lb-123",
			"lb-123:http",
			"lb-123:http:80:acl-1:extra",
			"lb-123:http:80:",
			"lb-123:http:0",
			"lb-123:http:70000",
			"lb-123:smtp:80",
			":http:80",
		} {
			t.Run(id, func(t *testing.T) {
				stub := &listenerAclRpcStub{}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
					TypeName: listenerAclAttachmentTestTypeName,
					ID:       id,
				})
				if err != nil {
					return // a protocol-level rejection is an acceptable failure
				}
				if !listenerAclTestHasError(imported.Diagnostics) {
					t.Fatalf("import ID %q must be rejected", id)
				}
				if len(imported.ImportedResources) != 0 {
					t.Fatalf("a rejected import ID must not produce a resource: %d", len(imported.ImportedResources))
				}
				if stub.describes != 0 {
					t.Fatalf("a rejected import ID must not reach Describe, got %d calls", stub.describes)
				}
			})
		}
	})
}

// An import must not claim a binding that the describe response cannot prove.
// The same incomplete responses are still kept for an ordinary refresh, where an
// existing resource must not be dropped.
func TestUnitListenerAclAttachmentProtocolImportUnconfirmed(t *testing.T) {
	responses := []struct {
		name     string
		response map[string]interface{}
	}{
		{name: "missing AclStatus", response: map[string]interface{}{"ListenerPort": 80, "AclType": "white", "AclId": "acl-1"}},
		{name: "wrong AclStatus type", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": 1, "AclType": "white", "AclId": "acl-1"}},
		{name: "missing AclType", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclId": "acl-1"}},
		{name: "wrong AclType type", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": 1, "AclId": "acl-1"}},
		{name: "malformed member", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclIds": map[string]interface{}{"AclId": []interface{}{"acl-1", 42}}}},
		{name: "non string AclId", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": []interface{}{"acl-1"}}},
	}

	for _, tc := range responses {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, id := range []string{"lb-123:http:80", "lb-123:http:80:acl-1"} {
				id := id
				t.Run(id, func(t *testing.T) {
					stub := &listenerAclRpcStub{
						describe: func() (map[string]interface{}, error) { return tc.response, nil },
					}
					stub.install(t)
					server := newListenerAclAttachmentServer(t)
					imported, err := server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
						TypeName: listenerAclAttachmentTestTypeName,
						ID:       id,
					})
					if err != nil {
						t.Fatalf("ImportResourceState: %s", err)
					}
					listenerAclTestRequireNoError(t, "ImportResourceState", imported.Diagnostics)
					read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
						TypeName:     listenerAclAttachmentTestTypeName,
						CurrentState: imported.ImportedResources[0].State,
					})
					if err != nil {
						t.Fatalf("ReadResource: %s", err)
					}
					if !listenerAclTestHasError(read.Diagnostics) {
						t.Fatal("an incomplete response must not confirm the imported binding")
					}
					if !listenerAclTestDiagnosticsContain(read.Diagnostics, "cannot confirm the imported binding") {
						t.Fatalf("the diagnostic did not report the pending confirmation: %v", read.Diagnostics)
					}
				})
			}

			t.Run("refresh keeps the existing resource", func(t *testing.T) {
				stub := &listenerAclRpcStub{
					describe: func() (map[string]interface{}, error) { return tc.response, nil },
				}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				prior := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
					"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
					"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
					"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
					"acl_type":    tftypes.NewValue(tftypes.String, "white"),
					"acl_status":  tftypes.NewValue(tftypes.String, "on"),
				}))
				read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
					TypeName:     listenerAclAttachmentTestTypeName,
					CurrentState: prior,
				})
				if err != nil {
					t.Fatalf("ReadResource: %s", err)
				}
				listenerAclTestRequireNoError(t, "ReadResource (refresh)", read.Diagnostics)
				attrs := listenerAclTestAttrs(t, read.NewState)
				if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
					t.Fatalf("an incomplete response must keep the existing resource, acl_id = %q", got)
				}
			})
		})
	}
}

// listenerAclTestDiagnosticsContain reports whether any diagnostic mentions the
// token; it keeps the assertions readable for protocol level errors.
func listenerAclTestDiagnosticsContain(diagnostics []*tfprotov5.Diagnostic, token string) bool {
	for _, d := range diagnostics {
		if strings.Contains(strings.ToLower(d.Summary+" "+d.Detail), strings.ToLower(token)) {
			return true
		}
	}
	return false
}

// Each single field update changes exactly one configured value, so the write
// it issues and the state it stores can be attributed to that change alone. A
// fixture that moves both fields at once cannot tell an acl_id only update from
// an acl_type only update.
func TestUnitListenerAclAttachmentProtocolSingleFieldUpdates(t *testing.T) {
	priorState := func(t *testing.T, aclId, aclType string) *tfprotov5.DynamicValue {
		t.Helper()
		return listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"acl_id":      tftypes.NewValue(tftypes.String, aclId),
			"acl_type":    tftypes.NewValue(tftypes.String, aclType),
			"acl_status":  tftypes.NewValue(tftypes.String, "on"),
		}))
	}
	apply := func(t *testing.T, stub *listenerAclRpcStub, prior *tfprotov5.DynamicValue, aclId, aclType string) map[string]tftypes.Value {
		t.Helper()
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		config := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"acl_id":      tftypes.NewValue(tftypes.String, aclId),
			"acl_type":    tftypes.NewValue(tftypes.String, aclType),
			"acl_status":  tftypes.NewValue(tftypes.String, "on"),
		}))
		plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
			TypeName:         listenerAclAttachmentTestTypeName,
			PriorState:       prior,
			ProposedNewState: config,
			Config:           config,
		})
		if err != nil {
			t.Fatalf("PlanResourceChange (update): %s", err)
		}
		listenerAclTestRequireNoError(t, "PlanResourceChange (update)", plan.Diagnostics)
		applied, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
			TypeName:       listenerAclAttachmentTestTypeName,
			PriorState:     prior,
			PlannedState:   plan.PlannedState,
			Config:         config,
			PlannedPrivate: plan.PlannedPrivate,
		})
		if err != nil {
			t.Fatalf("ApplyResourceChange (update): %s", err)
		}
		listenerAclTestRequireNoError(t, "ApplyResourceChange (update)", applied.Diagnostics)
		return listenerAclTestAttrs(t, applied.NewState)
	}

	t.Run("acl_id only", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-2"}, nil
			},
		}
		attrs := apply(t, stub, priorState(t, "acl-1", "white"), "acl-2", "white")
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-2" {
			t.Fatalf("acl_id = %q, want acl-2", got)
		}
		if got := listenerAclTestString(t, attrs["acl_type"]); got != "white" {
			t.Fatalf("acl_type = %q, want white", got)
		}
		if stub.sets != 1 {
			t.Fatalf("expected exactly one write for the acl_id change, got %d", stub.sets)
		}
		request := stub.requests[0]
		if request["AclId"] != "acl-2" || request["AclType"] != "white" || request["AclStatus"] != "on" {
			t.Fatalf("the acl_id only write carried the wrong parameters: %v", request)
		}
	})

	t.Run("acl_type only", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "black", "AclId": "acl-1"}, nil
			},
		}
		attrs := apply(t, stub, priorState(t, "acl-1", "white"), "acl-1", "black")
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
			t.Fatalf("acl_id = %q, want acl-1", got)
		}
		if got := listenerAclTestString(t, attrs["acl_type"]); got != "black" {
			t.Fatalf("acl_type = %q, want black", got)
		}
		if stub.sets != 1 {
			t.Fatalf("expected exactly one write for the acl_type change, got %d", stub.sets)
		}
		request := stub.requests[0]
		if request["AclId"] != "acl-1" || request["AclType"] != "black" || request["AclStatus"] != "on" {
			t.Fatalf("the acl_type only write carried the wrong parameters: %v", request)
		}
	})

	t.Run("no acl field change issues no write", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
			},
		}
		attrs := apply(t, stub, priorState(t, "acl-1", "white"), "acl-1", "white")
		if stub.sets != 0 {
			t.Fatalf("an unchanged configuration must not issue a write, got %d", stub.sets)
		}
		if got := listenerAclTestString(t, attrs["acl_status"]); got != "on" {
			t.Fatalf("acl_status = %q, want on", got)
		}
	})
}

// A response that reached the service but could not be decoded is reconciled
// against the listener instead of being reported as a plain failure.
func TestUnitListenerAclAttachmentProtocolUnreadableResponse(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 200 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })

	stub := &listenerAclRpcStub{
		set: func(map[string]interface{}) error {
			return fmt.Errorf("%w: unexpected EOF", connectivity.ErrUnreadableResponse)
		},
		describe: func() (map[string]interface{}, error) {
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)
	plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
	applied := listenerAclTestApplyCreate(t, server, plan, config)
	listenerAclTestRequireNoError(t, "ApplyResourceChange (unreadable response)", applied.Diagnostics)
	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
		t.Fatalf("reconciled binding was not kept, acl_id = %q", got)
	}
}

// A write whose response is truncated after the service received it leaves the
// outcome unknown. The transport has to classify it as unreadable and the
// resource has to reconcile it against the following readback.
func TestUnitListenerAclAttachmentTransportTruncatedWrite(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 3 * time.Second
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })

	var sets, describes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		switch r.Form.Get("Action") {
		case "SetLoadBalancerHTTPListenerAttribute":
			sets.Add(1)
			// Announce a larger body than is written and drop the connection so
			// the client observes a truncated, undecodable response.
			hijacker, ok := w.(http.Hijacker)
			if !ok {
				t.Errorf("the test server does not support hijacking")
				return
			}
			conn, buffer, err := hijacker.Hijack()
			if err != nil {
				t.Errorf("hijacking the response: %s", err)
				return
			}
			_, _ = buffer.WriteString("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 64\r\n\r\n{\"RequestId\":\"r-1\"")
			_ = buffer.Flush()
			_ = conn.Close()
		case "DescribeLoadBalancerHTTPListenerAttribute":
			describes.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"RequestId":"r-2","ListenerPort":80,"AclStatus":"on","AclType":"white","AclId":"acl-1"}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"InvalidAction","Message":"unexpected action","RequestId":"r-0"}`))
		}
	}))
	defer server.Close()

	client := newListenerAclAttachmentTransportClient(t, server.URL)
	providerServer := newListenerAclAttachmentServerWith(t, client)
	plan, config := listenerAclTestPlanCreate(t, providerServer, "lb-123:http:80", "acl-1", "white")
	applied := listenerAclTestApplyCreate(t, providerServer, plan, config)
	listenerAclTestRequireNoError(t, "ApplyResourceChange (truncated write)", applied.Diagnostics)

	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
		t.Fatalf("the reconciled binding was not kept, acl_id = %q", got)
	}
	if got := sets.Load(); got != 1 {
		t.Fatalf("expected exactly one write, got %d", got)
	}
	if got := describes.Load(); got == 0 {
		t.Fatal("the unreadable write outcome must be reconciled with a readback")
	}
	warned := false
	for _, d := range applied.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityWarning {
			warned = true
		}
	}
	if !warned {
		t.Fatal("expected a warning for the reconciled write")
	}
}

// A terminal HTTP rejection is an explicit failure: it must not be reconciled,
// must not invent state and must not trigger another request.
func TestUnitListenerAclAttachmentTransportTerminalRejection(t *testing.T) {
	var sets, describes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		if r.Form.Get("Action") == "SetLoadBalancerHTTPListenerAttribute" {
			sets.Add(1)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"Code":"IpVersionConflict","Message":"the ACL IP version conflicts with the listener","RequestId":"r-9"}`))
			return
		}
		describes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"RequestId":"r-2","ListenerPort":80,"AclStatus":"off"}`))
	}))
	defer server.Close()

	client := newListenerAclAttachmentTransportClient(t, server.URL)
	providerServer := newListenerAclAttachmentServerWith(t, client)
	plan, config := listenerAclTestPlanCreate(t, providerServer, "lb-123:http:80", "acl-1", "white")
	applied := listenerAclTestApplyCreate(t, providerServer, plan, config)

	if !listenerAclTestHasError(applied.Diagnostics) {
		t.Fatal("a terminal rejection must fail the create")
	}
	if !listenerAclTestDiagnosticsContain(applied.Diagnostics, "IpVersionConflict") {
		t.Fatalf("the terminal error code was lost: %v", applied.Diagnostics)
	}
	value, err := applied.NewState.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the state: %s", err)
	}
	if !value.IsNull() {
		t.Fatalf("a terminal rejection must not write state, got %s", value)
	}
	if got := sets.Load(); got != 1 {
		t.Fatalf("a terminal rejection must not be retried, got %d writes", got)
	}
	if got := describes.Load(); got != 0 {
		t.Fatalf("a terminal rejection must not be reconciled, got %d reads", got)
	}
}

// A Describe response that is valid JSON but carries no usable AclStatus cannot
// prove that the binding is gone. Read must keep the state and Delete must fail
// instead of reporting a successful unbinding.
func TestUnitListenerAclAttachmentProtocolIncompleteAclStatus(t *testing.T) {
	prior := func(t *testing.T) *tfprotov5.DynamicValue {
		t.Helper()
		return listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
			"acl_type":    tftypes.NewValue(tftypes.String, "white"),
			"acl_status":  tftypes.NewValue(tftypes.String, "on"),
		}))
	}

	for _, tc := range []struct {
		name     string
		response map[string]interface{}
	}{
		{name: "missing AclStatus", response: map[string]interface{}{"ListenerPort": 80, "AclType": "white", "AclId": "acl-1"}},
		{name: "wrong AclStatus type", response: map[string]interface{}{"ListenerPort": 80, "AclStatus": 1, "AclType": "white", "AclId": "acl-1"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("read keeps the resource", func(t *testing.T) {
				stub := &listenerAclRpcStub{
					describe: func() (map[string]interface{}, error) { return tc.response, nil },
				}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
					TypeName:     listenerAclAttachmentTestTypeName,
					CurrentState: prior(t),
				})
				if err != nil {
					t.Fatalf("ReadResource: %s", err)
				}
				listenerAclTestRequireNoError(t, "ReadResource (incomplete AclStatus)", read.Diagnostics)
				attrs := listenerAclTestAttrs(t, read.NewState)
				if got := listenerAclTestString(t, attrs["id"]); got != "lb-123:http:80" {
					t.Fatalf("an incomplete response must keep the resource in state, id = %q", got)
				}
				if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
					t.Fatalf("acl_id = %q, want acl-1", got)
				}
			})

			t.Run("delete is not confirmed", func(t *testing.T) {
				stub := &listenerAclRpcStub{
					describe: func() (map[string]interface{}, error) { return tc.response, nil },
				}
				stub.install(t)
				server := newListenerAclAttachmentServer(t)
				state := prior(t)
				plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
					TypeName:         listenerAclAttachmentTestTypeName,
					PriorState:       state,
					ProposedNewState: listenerAclTestNull(t),
					Config:           listenerAclTestNull(t),
				})
				if err != nil {
					t.Fatalf("PlanResourceChange (destroy): %s", err)
				}
				listenerAclTestRequireNoError(t, "PlanResourceChange (destroy)", plan.Diagnostics)
				destroyed, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
					TypeName:       listenerAclAttachmentTestTypeName,
					PriorState:     state,
					PlannedState:   plan.PlannedState,
					Config:         listenerAclTestNull(t),
					PlannedPrivate: plan.PlannedPrivate,
				})
				if err != nil {
					t.Fatalf("ApplyResourceChange (destroy): %s", err)
				}
				if !listenerAclTestHasError(destroyed.Diagnostics) {
					t.Fatal("an incomplete AclStatus must not confirm the unbinding")
				}
				if stub.sets != 0 {
					t.Fatalf("no disabling write may be issued before the status is known, got %d", stub.sets)
				}
			})
		})
	}
}

// A 4xx status is an explicit rejection even when its body cannot be decoded: it
// must fail without a reconciliation read. A 2xx response with the same
// undecodable body leaves the outcome unknown and is reconciled instead.
func TestUnitListenerAclAttachmentTransportBadBodyClassification(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 500 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })

	newServer := func(t *testing.T, status int) (*httptest.Server, *atomic.Int32, *atomic.Int32) {
		t.Helper()
		var sets, describes atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = r.ParseForm()
			if r.Form.Get("Action") == "SetLoadBalancerHTTPListenerAttribute" {
				sets.Add(1)
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`<html><body>the body is not JSON</body></html>`))
				return
			}
			describes.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"RequestId":"r-2","ListenerPort":80,"AclStatus":"on","AclType":"white","AclId":"acl-1"}`))
		}))
		return server, &sets, &describes
	}

	t.Run("undecodable 4xx is an explicit rejection", func(t *testing.T) {
		server, sets, describes := newServer(t, http.StatusBadRequest)
		defer server.Close()
		client := newListenerAclAttachmentTransportClient(t, server.URL)
		providerServer := newListenerAclAttachmentServerWith(t, client)
		plan, config := listenerAclTestPlanCreate(t, providerServer, "lb-123:http:80", "acl-1", "white")
		applied := listenerAclTestApplyCreate(t, providerServer, plan, config)
		if !listenerAclTestHasError(applied.Diagnostics) {
			t.Fatal("an undecodable 4xx rejection must fail the create")
		}
		value, err := applied.NewState.Unmarshal(listenerAclAttachmentTestType)
		if err != nil {
			t.Fatalf("decoding the state: %s", err)
		}
		if !value.IsNull() {
			t.Fatalf("an explicit rejection must not write state, got %s", value)
		}
		if got := sets.Load(); got != 1 {
			t.Fatalf("expected exactly one write, got %d", got)
		}
		if got := describes.Load(); got != 0 {
			t.Fatalf("an explicit rejection must not be reconciled, got %d reads", got)
		}
	})

	t.Run("undecodable 2xx is an unknown write and is reconciled", func(t *testing.T) {
		server, sets, describes := newServer(t, http.StatusOK)
		defer server.Close()
		client := newListenerAclAttachmentTransportClient(t, server.URL)
		providerServer := newListenerAclAttachmentServerWith(t, client)
		plan, config := listenerAclTestPlanCreate(t, providerServer, "lb-123:http:80", "acl-1", "white")
		applied := listenerAclTestApplyCreate(t, providerServer, plan, config)
		listenerAclTestRequireNoError(t, "ApplyResourceChange (undecodable 2xx)", applied.Diagnostics)
		attrs := listenerAclTestAttrs(t, applied.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
			t.Fatalf("the reconciled binding was not kept, acl_id = %q", got)
		}
		if got := sets.Load(); got != 1 {
			t.Fatalf("expected exactly one write, got %d", got)
		}
		if got := describes.Load(); got == 0 {
			t.Fatal("an unknown 2xx outcome must be reconciled with a readback")
		}
	})
}

// A completed readback that still contradicts the plan is a failure, while an
// incomplete readback is only a pending confirmation. Neither may be reported as
// a successful apply that stores a binding the listener never confirmed.
func TestUnitListenerAclAttachmentProtocolReadbackOutcomes(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 300 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })
	previousInterval := *slbsvc.ListenerAclAttachmentPollIntervalForTest
	*slbsvc.ListenerAclAttachmentPollIntervalForTest = 10 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentPollIntervalForTest = previousInterval })

	t.Run("a complete but different binding keeps the planned identity", func(t *testing.T) {
		// The listener keeps reporting the previous type during the whole read
		// budget of the apply and converges only for the following refresh.
		stale := true
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				if stale {
					return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
				}
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "black", "AclId": "acl-1"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "black")
		applied := listenerAclTestApplyCreate(t, server, plan, config)
		if !listenerAclTestHasError(applied.Diagnostics) {
			t.Fatal("a complete readback that contradicts the plan must not be reported as success")
		}
		if !listenerAclTestDiagnosticsContain(applied.Diagnostics, "different access control binding") {
			t.Fatalf("the contradiction was not reported: %v", applied.Diagnostics)
		}
		attrs := listenerAclTestAttrs(t, applied.NewState)
		if got := listenerAclTestString(t, attrs["id"]); got != "lb-123:http:80" {
			t.Fatalf("the acknowledged identity was dropped, id = %q", got)
		}
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
			t.Fatalf("acl_id = %q, want acl-1", got)
		}
		if got := listenerAclTestString(t, attrs["acl_type"]); got != "black" {
			t.Fatalf("the planned type was dropped, acl_type = %q", got)
		}

		// Recovery: the next refresh sees the converged binding and keeps the
		// identity that the failed apply stored.
		stale = false
		read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
			TypeName:     listenerAclAttachmentTestTypeName,
			CurrentState: applied.NewState,
		})
		if err != nil {
			t.Fatalf("ReadResource after the contradiction: %s", err)
		}
		listenerAclTestRequireNoError(t, "ReadResource after the contradiction", read.Diagnostics)
		recovered := listenerAclTestAttrs(t, read.NewState)
		if got := listenerAclTestString(t, recovered["acl_type"]); got != "black" {
			t.Fatalf("the refreshed binding was lost, acl_type = %q", got)
		}
	})

	t.Run("an incomplete readback keeps the planned identity", func(t *testing.T) {
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				// The Describe contract requires AclType while access control is
				// on, so this response proves neither value.
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclId": "acl-1"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		plan, config := listenerAclTestPlanCreate(t, server, "lb-123:http:80", "acl-1", "white")
		applied := listenerAclTestApplyCreate(t, server, plan, config)
		listenerAclTestRequireNoError(t, "ApplyResourceChange (incomplete readback)", applied.Diagnostics)
		attrs := listenerAclTestAttrs(t, applied.NewState)
		if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
			t.Fatalf("the planned identity was not kept, acl_id = %q", got)
		}
		warned := false
		for _, d := range applied.Diagnostics {
			if d.Severity == tfprotov5.DiagnosticSeverityWarning {
				warned = true
			}
		}
		if !warned {
			t.Fatal("a pending confirmation must be reported as a warning")
		}
	})
}

// An acknowledged single field update whose readback keeps the old value must
// fail the apply but keep the planned identity, so the change can be reconciled
// instead of silently keeping the previous group.
func TestUnitListenerAclAttachmentProtocolUpdateContradictionKeepsPlan(t *testing.T) {
	previousOperation := *slbsvc.ListenerAclAttachmentOperationTimeoutForTest
	*slbsvc.ListenerAclAttachmentOperationTimeoutForTest = 300 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentOperationTimeoutForTest = previousOperation })
	previousInterval := *slbsvc.ListenerAclAttachmentPollIntervalForTest
	*slbsvc.ListenerAclAttachmentPollIntervalForTest = 10 * time.Millisecond
	t.Cleanup(func() { *slbsvc.ListenerAclAttachmentPollIntervalForTest = previousInterval })

	// The listener keeps reporting the previous type during the whole read budget
	// of the apply and converges only for the following refresh.
	stale := true
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			if stale {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
			}
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "black", "AclId": "acl-1"}, nil
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)

	prior := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "white"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))
	config := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "black"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))
	plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         listenerAclAttachmentTestTypeName,
		PriorState:       prior,
		ProposedNewState: config,
		Config:           config,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (update): %s", err)
	}
	listenerAclTestRequireNoError(t, "PlanResourceChange (update)", plan.Diagnostics)
	applied, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       listenerAclAttachmentTestTypeName,
		PriorState:     prior,
		PlannedState:   plan.PlannedState,
		Config:         config,
		PlannedPrivate: plan.PlannedPrivate,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (update): %s", err)
	}
	if !listenerAclTestHasError(applied.Diagnostics) {
		t.Fatal("a stale readback after an acknowledged update must fail the apply")
	}
	if !listenerAclTestDiagnosticsContain(applied.Diagnostics, "different access control binding") {
		t.Fatalf("the contradiction was not reported: %v", applied.Diagnostics)
	}
	attrs := listenerAclTestAttrs(t, applied.NewState)
	if got := listenerAclTestString(t, attrs["acl_type"]); got != "black" {
		t.Fatalf("the previous group was kept instead of the planned identity, acl_type = %q", got)
	}
	if got := listenerAclTestString(t, attrs["acl_id"]); got != "acl-1" {
		t.Fatalf("acl_id = %q, want acl-1", got)
	}
	if stub.sets != 1 {
		t.Fatalf("expected exactly one acknowledged write, got %d", stub.sets)
	}

	// Recovery: once the listener converges, the stored identity is kept.
	stale = false
	read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     listenerAclAttachmentTestTypeName,
		CurrentState: applied.NewState,
	})
	if err != nil {
		t.Fatalf("ReadResource after the update contradiction: %s", err)
	}
	listenerAclTestRequireNoError(t, "ReadResource after the update contradiction", read.Diagnostics)
	if got := listenerAclTestString(t, listenerAclTestAttrs(t, read.NewState)["acl_type"]); got != "black" {
		t.Fatalf("the refreshed update was lost, acl_type = %q", got)
	}
}

// Read has to follow the real on to off transition: the resource is kept while
// access control is enabled and removed once it is explicitly off.
func TestUnitListenerAclAttachmentProtocolReadTransitions(t *testing.T) {
	describes := 0
	stub := &listenerAclRpcStub{
		describe: func() (map[string]interface{}, error) {
			describes++
			if describes == 1 {
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
			}
			return map[string]interface{}{"ListenerPort": 80, "AclStatus": "off", "AclId": "acl-1"}, nil
		},
	}
	stub.install(t)
	server := newListenerAclAttachmentServer(t)
	prior := listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
		"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
		"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
		"acl_type":    tftypes.NewValue(tftypes.String, "white"),
		"acl_status":  tftypes.NewValue(tftypes.String, "on"),
	}))

	read, err := server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     listenerAclAttachmentTestTypeName,
		CurrentState: prior,
	})
	if err != nil {
		t.Fatalf("ReadResource (on): %s", err)
	}
	listenerAclTestRequireNoError(t, "ReadResource (on)", read.Diagnostics)
	if got := listenerAclTestString(t, listenerAclTestAttrs(t, read.NewState)["acl_id"]); got != "acl-1" {
		t.Fatalf("the enabled binding was dropped, acl_id = %q", got)
	}

	read, err = server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
		TypeName:     listenerAclAttachmentTestTypeName,
		CurrentState: prior,
	})
	if err != nil {
		t.Fatalf("ReadResource (off): %s", err)
	}
	listenerAclTestRequireNoError(t, "ReadResource (off)", read.Diagnostics)
	value, err := read.NewState.Unmarshal(listenerAclAttachmentTestType)
	if err != nil {
		t.Fatalf("decoding the read state: %s", err)
	}
	if !value.IsNull() {
		t.Fatalf("an explicitly disabled binding must be removed from state, got %s", value)
	}
}

// Destroy has to distinguish a confirmed unbinding from a rejected or unknown
// one: only the confirmed path may drop the state.
func TestUnitListenerAclAttachmentProtocolDestroyOutcomes(t *testing.T) {
	prior := func(t *testing.T) *tfprotov5.DynamicValue {
		t.Helper()
		return listenerAclTestDynamic(t, listenerAclTestObject(t, map[string]tftypes.Value{
			"id":          tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"listener_id": tftypes.NewValue(tftypes.String, "lb-123:http:80"),
			"acl_id":      tftypes.NewValue(tftypes.String, "acl-1"),
			"acl_type":    tftypes.NewValue(tftypes.String, "white"),
			"acl_status":  tftypes.NewValue(tftypes.String, "on"),
		}))
	}
	destroy := func(t *testing.T, server tfprotov5.ProviderServer, state *tfprotov5.DynamicValue) *tfprotov5.ApplyResourceChangeResponse {
		t.Helper()
		plan, err := server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
			TypeName:         listenerAclAttachmentTestTypeName,
			PriorState:       state,
			ProposedNewState: listenerAclTestNull(t),
			Config:           listenerAclTestNull(t),
		})
		if err != nil {
			t.Fatalf("PlanResourceChange (destroy): %s", err)
		}
		listenerAclTestRequireNoError(t, "PlanResourceChange (destroy)", plan.Diagnostics)
		destroyed, err := server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
			TypeName:       listenerAclAttachmentTestTypeName,
			PriorState:     state,
			PlannedState:   plan.PlannedState,
			Config:         listenerAclTestNull(t),
			PlannedPrivate: plan.PlannedPrivate,
		})
		if err != nil {
			t.Fatalf("ApplyResourceChange (destroy): %s", err)
		}
		return destroyed
	}

	t.Run("on to off is confirmed", func(t *testing.T) {
		describes := 0
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				describes++
				if describes == 1 {
					return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
				}
				return map[string]interface{}{"ListenerPort": 80, "AclStatus": "off"}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		destroyed := destroy(t, server, prior(t))
		listenerAclTestRequireNoError(t, "ApplyResourceChange (destroy on to off)", destroyed.Diagnostics)
		if stub.sets != 1 {
			t.Fatalf("expected exactly one disabling write, got %d", stub.sets)
		}
		value, err := destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
		if err != nil {
			t.Fatalf("decoding the destroyed state: %s", err)
		}
		if !value.IsNull() {
			t.Fatalf("a confirmed unbinding must drop the state, got %s", value)
		}
	})

	t.Run("incomplete response during post-write polling keeps the state", func(t *testing.T) {
		describes := 0
		stub := &listenerAclRpcStub{
			describe: func() (map[string]interface{}, error) {
				describes++
				if describes == 1 {
					return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
				}
				// The disabling write was acknowledged, but the polls that follow
				// it cannot confirm the new status.
				return map[string]interface{}{"ListenerPort": 80}, nil
			},
		}
		stub.install(t)
		server := newListenerAclAttachmentServer(t)
		destroyed := destroy(t, server, prior(t))
		if !listenerAclTestHasError(destroyed.Diagnostics) {
			t.Fatal("an incomplete response during the post-write polling must fail the destroy")
		}
		if stub.sets != 1 {
			t.Fatalf("expected the disabling write to be issued once, got %d", stub.sets)
		}
		value, err := destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
		if err != nil {
			t.Fatalf("decoding the destroyed state: %s", err)
		}
		if value.IsNull() {
			t.Fatal("an unconfirmed unbinding must keep the state")
		}
	})

	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "explicit rejection keeps the state", err: fmt.Errorf("%w: status 400", connectivity.ErrRpcRejectedResponse)},
		{name: "unknown result keeps the state", err: fmt.Errorf("%w: unexpected EOF", connectivity.ErrUnreadableResponse)},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			stub := &listenerAclRpcStub{
				set: func(map[string]interface{}) error { return tc.err },
				describe: func() (map[string]interface{}, error) {
					return map[string]interface{}{"ListenerPort": 80, "AclStatus": "on", "AclType": "white", "AclId": "acl-1"}, nil
				},
			}
			stub.install(t)
			server := newListenerAclAttachmentServer(t)
			destroyed := destroy(t, server, prior(t))
			if !listenerAclTestHasError(destroyed.Diagnostics) {
				t.Fatal("an unconfirmed unbinding must fail the destroy")
			}
			value, err := destroyed.NewState.Unmarshal(listenerAclAttachmentTestType)
			if err != nil {
				t.Fatalf("decoding the destroyed state: %s", err)
			}
			if value.IsNull() {
				t.Fatal("an unconfirmed unbinding must keep the state")
			}
		})
	}
}
