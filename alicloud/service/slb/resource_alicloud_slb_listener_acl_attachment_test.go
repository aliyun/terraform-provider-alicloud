package slb_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	slbsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/slb"
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
					// Bind a second group outside Terraform so the listener
					// exposes two ACLs while the resource still tracks the first.
					testAccBindExtraListenerAcl("alicloud_slb_acl.first", "alicloud_slb_acl.second"),
					testAccCaptureListenerIdentity(&listenerId, &secondAclId),
				),
			},
			// The first successful apply is imported before the second group is
			// used for the explicit four segment import below.
			testAccListenerAclAttachmentImportStep("alicloud_slb_acl.first"),
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
				ExpectError: regexp.MustCompile(`(?i)multiple access control policy groups`),
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
			// A combined ACL value would bind several groups while the resource
			// tracks one, so it is rejected during validation.
			{
				Config:      testAccListenerAclAttachmentConfigCombinedAcl(),
				ExpectError: regexp.MustCompile(`(?i)single access control policy group`),
			},
			// The listener identity is validated locally for both the
			// configuration and the import path.
			{
				Config:      testAccListenerAclAttachmentConfigInvalidListenerPort(),
				ExpectError: regexp.MustCompile(`(?i)between 1 and 65535`),
			},
			{
				Config:      testAccListenerAclAttachmentConfigInvalidListenerProtocol(),
				ExpectError: regexp.MustCompile(`(?i)protocol must be one of`),
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
