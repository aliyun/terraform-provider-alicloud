package slb_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

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
					testAccCaptureListenerId(listenerAclAttachmentListenerId, &listenerId),
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
					testAccCheckListenerIdUnchanged(listenerAclAttachmentListenerId, &listenerId),
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
					testAccCheckListenerRemoteAclNotFirst(listenerAclAttachmentListenerId, "alicloud_slb_acl.second"),
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

// testAccCheckListenerRemoteAclNotFirst proves that the requested group is not
// the first element of the response, so an importer that always takes the first
// group cannot pass the following import step by accident.
func testAccCheckListenerRemoteAclNotFirst(listenerResourceId, aclResource string) resource.TestCheckFunc {
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
		ids := slbsvc.ListenerAclIdsForTest(object)
		if len(ids) < 2 {
			return fmt.Errorf("listener %s reports %d policy groups, expected at least two", listener.Primary.ID, len(ids))
		}
		if ids[0] == acl.Primary.ID {
			return fmt.Errorf("the requested policy group %s is the first element of the response %v; the import case needs it at another position", acl.Primary.ID, ids)
		}
		return nil
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
func TestAccAliCloudSLBListenerAclAttachment_negative(t *testing.T) {
	name := fmt.Sprintf("tf-testAccSlbListenerAclNeg%d", tfacctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
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
			// The cloud backed negative steps stay last: the final config must
			// remain valid so the post-test destroy can clean up.
			{
				Config:      testAccListenerAclAttachmentConfigMissingListener(name),
				ExpectError: regexp.MustCompile(`(?i)(not exist|notfound|invalidloadbalancerid)`),
			},
			{
				Config:      testAccListenerAclAttachmentConfigMissingAcl(name),
				ExpectError: regexp.MustCompile(`(?i)(aclnotexist|acl.?not.?exist|specified acl)`),
			},
			{
				Config:      testAccListenerAclAttachmentConfigIpVersionConflict(name),
				ExpectError: regexp.MustCompile(`(?i)ipversionconflict`),
			},
			{
				// A rejected request must not leave attachment state or write
				// the listener ACLs.
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNoListenerAclAttachmentState,
					testAccCheckListenerRemoteAclStatus(listenerAclAttachmentListenerId, "off"),
				),
			},
		},
	})
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

func testAccCaptureListenerId(listenerResourceId string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		listener := s.RootModule().Resources[listenerResourceId]
		if listener == nil {
			return fmt.Errorf("cannot find %s state", listenerResourceId)
		}
		*id = listener.Primary.ID
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

func testAccListenerAclAttachmentConfigMissingListener(name string) string {
	return fmt.Sprintf(`%s
resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = "lb-mock:http:80"
  acl_id      = alicloud_slb_acl.first.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentCommonConfig(name))
}

func testAccListenerAclAttachmentConfigMissingAcl(name string) string {
	return fmt.Sprintf(`%s
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

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = "acl-notexist00000"
  acl_type    = "white"
}
`, testAccListenerAclAttachmentCommonConfig(name))
}

func testAccListenerAclAttachmentConfigIpVersionConflict(name string) string {
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

resource "alicloud_slb_listener_acl_attachment" "default" {
  listener_id = alicloud_slb_listener.default.id
  acl_id      = alicloud_slb_acl.ipv6.id
  acl_type    = "white"
}
`, testAccListenerAclAttachmentCommonConfig(name))
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
	server := providerserver.NewProtocol5(&listenerAclAttachmentTestProvider{
		client: &connectivity.AliyunClient{RegionId: "cn-beijing"},
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

// listenerAclRpcStub replaces the SLB transport for one test.
type listenerAclRpcStub struct {
	set       func(request map[string]interface{}) error
	describe  func() (map[string]interface{}, error)
	sets      int
	describes int
}

func (s *listenerAclRpcStub) install(t *testing.T) {
	t.Helper()
	original := *slbsvc.ListenerAclAttachmentRPCPostForTest
	*slbsvc.ListenerAclAttachmentRPCPostForTest = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}) (map[string]interface{}, error) {
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
