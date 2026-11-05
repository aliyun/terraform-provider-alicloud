package slb_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	slbsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	frameworkschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckListenerAclAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccListenerAclAttachmentConfigHttp(name, "first", "white", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckListenerAclAttachmentAttrs("alicloud_slb_acl.first", "white"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentResourceId, "listener_id", listenerAclAttachmentListenerId, "id"),
					resource.TestCheckResourceAttrPair(listenerAclAttachmentDsId, "slb_listeners.0.id", listenerAclAttachmentListenerId, "id"),
					testAccCheckListenerNonAclAttrsRemote(listenerAclAttachmentListenerId, name),
					resource.ComposeTestCheckFunc(testAccListenerAttrChecks(name)...),
					testAccCaptureListenerNonAclAttrs(listenerAclAttachmentListenerId, &createdNonAclAttrs),
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
		for _, id := range slbsvc.ListenerAclIds(object) {
			if id == acl.Primary.ID {
				return nil
			}
		}
		return fmt.Errorf("listener %s does not report policy group %s, reported: %v", listener.Primary.ID, acl.Primary.ID, slbsvc.ListenerAclIds(object))
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
		ids := slbsvc.ListenerAclIds(object)
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
				ExpectError: regexp.MustCompile(`(?i)(aclnotexist|not exist|invalidparameter)`),
			},
			{
				Config:      testAccListenerAclAttachmentConfigIpVersionConflict(name),
				ExpectError: regexp.MustCompile(`(?i)(ipversionconflict|ip version|invalidparameter)`),
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
// Unit coverage for the resource lifecycle. The acceptance scenarios above use
// a live account; these cases drive the same lifecycle functions with a fixed
// fake API so each failure mode is asserted on the exact call sequence.
// ---------------------------------------------------------------------------

// fakeListenerAclAttachmentApi implements slbsvc.ListenerAclAttachmentApi.
type fakeListenerAclAttachmentApi struct {
	describe  func(ctx context.Context, listenerId string) (map[string]interface{}, error)
	set       func(ctx context.Context, listenerId, aclId, aclType string, enable bool) error
	describes int
	sets      int
	setAclIds []string
}

func (f *fakeListenerAclAttachmentApi) Describe(ctx context.Context, listenerId string) (map[string]interface{}, error) {
	f.describes++
	if f.describe == nil {
		return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
	}
	return f.describe(ctx, listenerId)
}

func (f *fakeListenerAclAttachmentApi) SetAcl(ctx context.Context, listenerId, aclId, aclType string, enable bool) error {
	f.sets++
	f.setAclIds = append(f.setAclIds, aclId)
	if f.set == nil {
		return nil
	}
	return f.set(ctx, listenerId, aclId, aclType, enable)
}

func listenerAclAttachmentTestPlan(listenerId, aclId, aclType string) slbsvc.ListenerAclAttachmentModel {
	return slbsvc.ListenerAclAttachmentModel{
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue(aclId),
		AclType:    types.StringValue(aclType),
	}
}

func TestUnitListenerAclAttachmentParseId(t *testing.T) {
	cases := []struct {
		name       string
		id         string
		listenerId string
		aclId      string
		wantErr    bool
	}{
		{name: "three segments", id: "lb-123:http:80", listenerId: "lb-123:http:80"},
		{name: "four segments", id: "lb-123:https:443:acl-456", listenerId: "lb-123:https:443", aclId: "acl-456"},
		{name: "two segments", id: "lb-123:80", wantErr: true},
		{name: "five segments", id: "lb-123:http:80:acl-1:extra", wantErr: true},
		{name: "empty protocol", id: "lb-123::80", wantErr: true},
		{name: "unsupported protocol", id: "lb-123:smtp:80", wantErr: true},
		{name: "port zero", id: "lb-123:http:0", wantErr: true},
		{name: "port above range", id: "lb-123:http:65536", wantErr: true},
		{name: "non numeric port", id: "lb-123:http:http", wantErr: true},
		{name: "empty acl suffix", id: "lb-123:http:80:", wantErr: true},
		{name: "combined acl suffix", id: "lb-123:http:80:acl-1,acl-2", wantErr: true},
		{name: "port at upper bound", id: "lb-123:udp:65535", listenerId: "lb-123:udp:65535"},
		{name: "port at lower bound", id: "lb-123:tcp:1", listenerId: "lb-123:tcp:1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			listenerId, aclId, err := slbsvc.ParseListenerAclAttachmentId(tc.id)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %q", tc.id)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %s", tc.id, err)
			}
			if listenerId != tc.listenerId || aclId != tc.aclId {
				t.Fatalf("got (%q, %q), want (%q, %q)", listenerId, aclId, tc.listenerId, tc.aclId)
			}
		})
	}
}

// Two listeners that share the same load balancer and frontend port but use
// different protocols must dispatch both the write and the read to their own
// action.
func TestUnitListenerAclAttachmentProtocolDispatch(t *testing.T) {
	wantSet := map[string]string{
		"http":  "SetLoadBalancerHTTPListenerAttribute",
		"https": "SetLoadBalancerHTTPSListenerAttribute",
		"tcp":   "SetLoadBalancerTCPListenerAttribute",
		"udp":   "SetLoadBalancerUDPListenerAttribute",
	}
	wantDescribe := map[string]string{
		"http":  "DescribeLoadBalancerHTTPListenerAttribute",
		"https": "DescribeLoadBalancerHTTPSListenerAttribute",
		"tcp":   "DescribeLoadBalancerTCPListenerAttribute",
		"udp":   "DescribeLoadBalancerUDPListenerAttribute",
	}

	seen := map[string]string{}
	for protocol, setAction := range wantSet {
		listenerId, aclId, err := slbsvc.ParseListenerAclAttachmentId("lb-shared:" + protocol + ":80:acl-1")
		if err != nil {
			t.Fatalf("parse %s listener id: %s", protocol, err)
		}
		if listenerId != "lb-shared:"+protocol+":80" || aclId != "acl-1" {
			t.Fatalf("unexpected identity for %s: (%q, %q)", protocol, listenerId, aclId)
		}
		action, err := slbsvc.ListenerAttributeAction(protocol)
		if err != nil || action != setAction {
			t.Fatalf("set action for %s = %q (%v), want %q", protocol, action, err, setAction)
		}
		describeAction, err := slbsvc.ListenerAttributeDescribeAction(protocol)
		if err != nil || describeAction != wantDescribe[protocol] {
			t.Fatalf("describe action for %s = %q (%v), want %q", protocol, describeAction, err, wantDescribe[protocol])
		}
		for otherProtocol, otherAction := range seen {
			if otherAction == action {
				t.Fatalf("%s and %s share set action %q", protocol, otherProtocol, action)
			}
		}
		seen[protocol] = action
	}

	for _, protocol := range []string{"HTTP", "smtp", "", "http2"} {
		if _, err := slbsvc.ListenerAttributeAction(protocol); err == nil {
			t.Fatalf("expected unsupported protocol %q to fail the set action", protocol)
		}
		if _, err := slbsvc.ListenerAttributeDescribeAction(protocol); err == nil {
			t.Fatalf("expected unsupported protocol %q to fail the describe action", protocol)
		}
	}
}

func TestUnitListenerAclIdsFromDescribe(t *testing.T) {
	object := map[string]interface{}{
		"AclId": "acl-a,acl-b",
		"AclIds": map[string]interface{}{
			"AclId": []interface{}{"acl-b", "acl-c"},
		},
	}
	ids := slbsvc.ListenerAclIds(object)
	if strings.Join(ids, ",") != "acl-a,acl-b,acl-c" {
		t.Fatalf("unexpected acl list: %v", ids)
	}

	single := map[string]interface{}{"AclId": "acl-only"}
	if got := slbsvc.ListenerAclIds(single); len(got) != 1 || got[0] != "acl-only" {
		t.Fatalf("unexpected single acl list: %v", got)
	}
}

// The schema rejects a combined ACL value and an invalid listener identity
// locally, before any request is built.
func TestUnitListenerAclAttachmentSchemaValidation(t *testing.T) {
	res := slbsvc.NewListenerAclAttachmentResource()
	var schemaResp frameworkresource.SchemaResponse
	res.Schema(context.Background(), frameworkresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building the schema failed: %v", schemaResp.Diagnostics)
	}

	validate := func(attribute, value string) bool {
		attr, ok := schemaResp.Schema.Attributes[attribute].(frameworkschema.StringAttribute)
		if !ok {
			t.Fatalf("%s is not a string attribute", attribute)
		}
		var resp validator.StringResponse
		for _, v := range attr.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root(attribute),
				ConfigValue: types.StringValue(value),
			}, &resp)
		}
		return resp.Diagnostics.HasError()
	}

	cases := []struct {
		attribute string
		value     string
		wantError bool
	}{
		{attribute: "acl_id", value: "acl-123", wantError: false},
		{attribute: "acl_id", value: "acl-123,acl-456", wantError: true},
		{attribute: "acl_id", value: "acl 123", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:80", wantError: false},
		{attribute: "listener_id", value: "lb-123:https:65535", wantError: false},
		{attribute: "listener_id", value: "lb-123:http:0", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:65536", wantError: true},
		{attribute: "listener_id", value: "lb-123:smtp:80", wantError: true},
		{attribute: "listener_id", value: "lb-123:http:http", wantError: true},
	}
	for _, tc := range cases {
		if got := validate(tc.attribute, tc.value); got != tc.wantError {
			t.Fatalf("validate %s=%q: got error %v, want %v", tc.attribute, tc.value, got, tc.wantError)
		}
	}
}

func TestUnitListenerAclAttachmentCreateLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")

	t.Run("write then read failure keeps identity", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, errors.New("read exploded")
			},
		}
		model, warnings, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(warnings) != 1 {
			t.Fatalf("expected one warning about the failed read, got %v", warnings)
		}
		if model.Id.ValueString() != listenerId || model.ListenerId.ValueString() != listenerId || model.AclId.ValueString() != "acl-1" {
			t.Fatalf("identity was not kept after the failed read: %+v", model)
		}
	})

	t.Run("unknown write outcome confirmed keeps identity", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if model.AclId.ValueString() != "acl-1" || model.Id.ValueString() != listenerId {
			t.Fatalf("confirmed binding was not kept in state: %+v", model)
		}
		if api.sets != 1 {
			t.Fatalf("expected exactly one write, got %d", api.sets)
		}
	})

	t.Run("unknown write outcome with disabled access control returns error", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "off", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err == nil {
			t.Fatal("expected an error when access control is disabled")
		}
		if model.AclId.ValueString() != "" {
			t.Fatalf("state must stay empty: %+v", model)
		}
	})

	t.Run("unknown write outcome with stale type returns error", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "black")
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err == nil {
			t.Fatal("expected an error when the listener still reports the previous type")
		}
		if model.AclId.ValueString() != "" {
			t.Fatalf("state must stay empty: %+v", model)
		}
	})

	t.Run("acknowledged write waits for the planned type", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "black")
		describes := 0
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				describes++
				if describes == 1 {
					return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
				}
				return map[string]interface{}{"AclStatus": "on", "AclType": "black", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, warnings, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err != nil || len(warnings) != 0 {
			t.Fatalf("unexpected result: warnings=%v err=%v", warnings, err)
		}
		if model.AclType.ValueString() != "black" {
			t.Fatalf("stale type was stored: %+v", model)
		}
		if describes < 2 {
			t.Fatalf("the read must wait for the planned type, describes=%d", describes)
		}
	})

	t.Run("terminal write error makes a single call and keeps no state", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return &tea.SDKError{Code: tea.String("QuotaLimitExceeds.AclListenerOverLimit"), Message: tea.String("quota exceeded"), Data: tea.String("{}")}
			},
		}
		model, _, err := slbsvc.CreateListenerAclAttachment(context.Background(), api, plan)
		if err == nil {
			t.Fatal("expected the quota error to be returned")
		}
		if api.sets != 1 || api.describes != 0 {
			t.Fatalf("terminal error must fail fast: writes=%d describes=%d", api.sets, api.describes)
		}
		if model.AclId.ValueString() != "" {
			t.Fatalf("no state may be written for a rejected request: %+v", model)
		}
	})
}

func TestUnitListenerAclAttachmentUpdateLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	state := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")
	state.Id = types.StringValue(listenerId)
	state.AclStatus = types.StringValue("on")

	t.Run("read failure keeps the planned identity", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-2", "black")
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, errors.New("read exploded")
			},
		}
		model, warnings, err := slbsvc.UpdateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(warnings) != 1 {
			t.Fatalf("expected one warning, got %v", warnings)
		}
		if model.AclId.ValueString() != "acl-2" || model.AclType.ValueString() != "black" {
			t.Fatalf("planned identity was not kept: %+v", model)
		}
		if api.sets != 1 {
			t.Fatalf("expected one write, got %d", api.sets)
		}
	})

	t.Run("unchanged identity skips the write", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "white")
		api := &fakeListenerAclAttachmentApi{}
		_, warnings, err := slbsvc.UpdateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil || len(warnings) != 0 {
			t.Fatalf("unexpected result: warnings=%v err=%v", warnings, err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write when only unrelated values change, got %d", api.sets)
		}
	})

	t.Run("unknown update outcome with stale type returns error", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "black")
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := slbsvc.UpdateListenerAclAttachment(context.Background(), api, plan, state)
		if err == nil {
			t.Fatal("expected an error when the listener still reports the previous type")
		}
		if model.AclId.ValueString() != "" {
			t.Fatalf("state must stay empty: %+v", model)
		}
	})

	t.Run("unknown update outcome confirmed keeps the plan", func(t *testing.T) {
		plan := listenerAclAttachmentTestPlan(listenerId, "acl-1", "black")
		api := &fakeListenerAclAttachmentApi{
			set: func(context.Context, string, string, string, bool) error {
				return errors.New(`Post "https://slb.aliyuncs.com": dial tcp 10.0.0.1:443: i/o timeout`)
			},
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "black", "AclId": "acl-1", "ListenerPort": 80}, nil
			},
		}
		model, _, err := slbsvc.UpdateListenerAclAttachment(context.Background(), api, plan, state)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if model.AclId.ValueString() != "acl-1" || model.AclType.ValueString() != "black" {
			t.Fatalf("planned binding was not kept: %+v", model)
		}
	})
}

func TestUnitListenerAclAttachmentReadLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"
	baseState := slbsvc.ListenerAclAttachmentModel{
		Id:         types.StringValue(listenerId),
		ListenerId: types.StringValue(listenerId),
		AclId:      types.StringValue("acl-1"),
		AclType:    types.StringValue("white"),
	}

	cases := []struct {
		name        string
		object      map[string]interface{}
		describeErr error
		wantFound   bool
		wantAclId   string
	}{
		{
			name:      "bound group keeps the resource",
			object:    map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-1,acl-2"},
			wantFound: true,
			wantAclId: "acl-1",
		},
		{
			name:      "target that is not the first group keeps its identity",
			object:    map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-a,acl-b,acl-c"},
			wantFound: true,
			wantAclId: "acl-b",
		},
		{
			name:      "disabled access control removes the resource",
			object:    map[string]interface{}{"AclStatus": "off", "AclId": "acl-1"},
			wantFound: false,
		},
		{
			name:      "missing group removes the resource",
			object:    map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-2"},
			wantFound: false,
		},
		{
			name:        "missing listener removes the resource",
			describeErr: listenerAclAttachmentNotFoundError(),
			wantFound:   false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := baseState
			if tc.name == "target that is not the first group keeps its identity" {
				state.AclId = types.StringValue("acl-b")
			}
			api := &fakeListenerAclAttachmentApi{
				describe: func(context.Context, string) (map[string]interface{}, error) {
					return tc.object, tc.describeErr
				},
			}
			model, found, err := slbsvc.ReadListenerAclAttachment(context.Background(), api, state)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if found != tc.wantFound {
				t.Fatalf("found = %v, want %v", found, tc.wantFound)
			}
			if tc.wantFound && model.AclId.ValueString() != tc.wantAclId {
				t.Fatalf("acl_id = %q, want %q", model.AclId.ValueString(), tc.wantAclId)
			}
		})
	}

	t.Run("import without acl suffix picks the only group", func(t *testing.T) {
		state := slbsvc.ListenerAclAttachmentModel{Id: types.StringValue(listenerId)}
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-9"}, nil
			},
		}
		model, found, err := slbsvc.ReadListenerAclAttachment(context.Background(), api, state)
		if err != nil || !found || model.AclId.ValueString() != "acl-9" {
			t.Fatalf("unexpected result: found=%v model=%+v err=%v", found, model, err)
		}
	})

	t.Run("import without acl suffix rejects several groups", func(t *testing.T) {
		state := slbsvc.ListenerAclAttachmentModel{Id: types.StringValue(listenerId)}
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on", "AclType": "white", "AclId": "acl-a,acl-b"}, nil
			},
		}
		_, _, err := slbsvc.ReadListenerAclAttachment(context.Background(), api, state)
		if err == nil || !strings.Contains(err.Error(), "multiple access control policy groups") {
			t.Fatalf("expected the ambiguity diagnostic, got %v", err)
		}
	})
}

func TestUnitListenerAclAttachmentDeleteLifecycle(t *testing.T) {
	listenerId := "lb-123:http:80"

	t.Run("missing listener is already gone", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return nil, listenerAclAttachmentNotFoundError()
			},
		}
		if err := slbsvc.DeleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write, got %d", api.sets)
		}
	})

	t.Run("disabled access control makes no write", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "off"}, nil
			},
		}
		if err := slbsvc.DeleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 0 {
			t.Fatalf("expected no write, got %d", api.sets)
		}
	})

	t.Run("disable then wait until off", func(t *testing.T) {
		describeCalls := 0
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				describeCalls++
				if describeCalls == 1 {
					return map[string]interface{}{"AclStatus": "on"}, nil
				}
				return map[string]interface{}{"AclStatus": "off"}, nil
			},
		}
		if err := slbsvc.DeleteListenerAclAttachment(context.Background(), api, listenerId); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if api.sets != 1 || api.setAclIds[0] != "" {
			t.Fatalf("expected one disable write, got %d (%v)", api.sets, api.setAclIds)
		}
	})

	t.Run("cancelled deadline stops the wait", func(t *testing.T) {
		api := &fakeListenerAclAttachmentApi{
			describe: func(context.Context, string) (map[string]interface{}, error) {
				return map[string]interface{}{"AclStatus": "on"}, nil
			},
		}
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		start := time.Now()
		if err := slbsvc.DeleteListenerAclAttachment(ctx, api, listenerId); err == nil {
			t.Fatal("expected the deadline to fail the delete")
		}
		if elapsed := time.Since(start); elapsed > 5*time.Second {
			t.Fatalf("delete did not stop promptly: %s", elapsed)
		}
	})
}

// listenerAclAttachmentNotFoundError mirrors the error shape the provider
// detects as not found.
func listenerAclAttachmentNotFoundError() error {
	return &tea.SDKError{
		StatusCode: tea.Int(404),
		Code:       tea.String("InvalidLoadBalancerId.NotFound"),
		Message:    tea.String("The specified resource does not exist."),
		Data:       tea.String("{}"),
	}
}
