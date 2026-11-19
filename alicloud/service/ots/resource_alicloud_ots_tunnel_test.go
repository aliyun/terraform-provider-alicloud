package ots_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// checkOtsTunnelDestroy wraps the target destroy check. Terraform destroys the instance
// and table after the tunnel, so the tunnel endpoint can stop resolving before CheckDestroy
// runs. A DNS failure alone is not proof that the tunnel is gone: it is accepted only when
// an independent parent-instance control-plane query reports NotFound.
func checkOtsTunnelDestroy(destroy resource.TestCheckFunc) resource.TestCheckFunc {
	return checkOtsTunnelDestroyWith(destroy, otsParentInstanceExists)
}

func checkOtsTunnelDestroyWith(destroy resource.TestCheckFunc, instanceExists func(string) (bool, error)) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := destroy(s); err != nil {
			if !strings.Contains(err.Error(), "no such host") {
				return err
			}
			instanceName, ok := otsTunnelInstanceName(s)
			if !ok {
				return err
			}
			exists, lookupErr := instanceExists(instanceName)
			if lookupErr != nil || exists {
				return err
			}
			return nil
		}
		return nil
	}
}

func otsTunnelInstanceName(s *terraform.State) (string, bool) {
	if s == nil {
		return "", false
	}
	for _, mod := range s.Modules {
		if mod == nil {
			continue
		}
		for _, rs := range mod.Resources {
			if rs.Type != "alicloud_ots_tunnel" || rs.Primary == nil {
				continue
			}
			parts := strings.Split(rs.Primary.ID, ":")
			if len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != "" {
				return parts[0], true
			}
		}
	}
	return "", false
}

func otsParentInstanceExists(instanceName string) (bool, error) {
	svc := ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	_, err := svc.DescribeOtsInstance(instanceName)
	if err == nil {
		return true, nil
	}
	if errs.NotFoundError(err) {
		return false, nil
	}
	return false, err
}

func TestUnitOtsTunnelCheckDestroyDNSRequiresInstanceNotFound(t *testing.T) {
	state := &terraform.State{Modules: []*terraform.ModuleState{
		{Path: []string{}, Resources: map[string]*terraform.ResourceState{
			"alicloud_ots_tunnel.default": {
				Type:    "alicloud_ots_tunnel",
				Primary: &terraform.InstanceState{ID: "instance:table:tunnel"},
			},
		}},
	}}
	dnsErr := fmt.Errorf("Post \"https://instance.example/tunnel/describe\": dial tcp: lookup instance: no such host")
	destroy := func(*terraform.State) error { return dnsErr }

	cases := []struct {
		name      string
		exists    bool
		lookupErr error
		wantErr   bool
	}{
		{name: "instance alive", exists: true, wantErr: true},
		{name: "instance not found", exists: false, wantErr: false},
		{name: "instance lookup failed", exists: false, lookupErr: fmt.Errorf("lookup failed"), wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkOtsTunnelDestroyWith(destroy, func(string) (bool, error) {
				return tc.exists, tc.lookupErr
			})(state)
			if tc.wantErr {
				if err != dnsErr {
					t.Fatalf("expected the original DNS error, got %v", err)
				}
			} else if err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}

	missingState := &terraform.State{Modules: []*terraform.ModuleState{
		{Path: []string{}, Resources: map[string]*terraform.ResourceState{}},
	}}
	if err := checkOtsTunnelDestroyWith(destroy, func(string) (bool, error) {
		return false, nil
	})(missingState); err != dnsErr {
		t.Fatalf("a DNS error without a tunnel ID must keep the original error, got %v", err)
	}

	// Malformed IDs with empty components are not a trusted tunnel identity. Even
	// when the parent-instance lookup reports NotFound, the original DNS error
	// must be preserved instead of being excused by a partial ID.
	malformedIDs := []string{"instance::", "instance:table:", "instance::tunnel", "instance:table"}
	for _, id := range malformedIDs {
		t.Run("malformed ID "+id, func(t *testing.T) {
			malformedState := &terraform.State{Modules: []*terraform.ModuleState{
				{Path: []string{}, Resources: map[string]*terraform.ResourceState{
					"alicloud_ots_tunnel.default": {
						Type:    "alicloud_ots_tunnel",
						Primary: &terraform.InstanceState{ID: id},
					},
				}},
			}}
			if err := checkOtsTunnelDestroyWith(destroy, func(string) (bool, error) {
				return false, nil
			})(malformedState); err != dnsErr {
				t.Fatalf("malformed tunnel ID %q must keep the original DNS error, got %v", id, err)
			}
		})
	}

	otherErr := fmt.Errorf("describe failed")
	if err := checkOtsTunnelDestroyWith(func(*terraform.State) error { return otherErr }, func(string) (bool, error) {
		return false, nil
	})(state); err != otherErr {
		t.Fatalf("non-DNS error must propagate, got %v", err)
	}
	if err := checkOtsTunnelDestroyWith(func(*terraform.State) error { return nil }, func(string) (bool, error) {
		return true, nil
	})(state); err != nil {
		t.Fatalf("successful destroy must stay successful, got %v", err)
	}
}

func TestAccAliCloudOtsTunnel_basic(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkOtsTunnelDestroy(rac.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   "BaseAndStream",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": acctest.CHECKSET,
						"table_name":    acctest.CHECKSET,
						"tunnel_name":   acctest.CHECKSET,
						"tunnel_type":   "BaseAndStream",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// channel_status is server-side computed by DescribeTunnel and can
				// transition asynchronously (WAIT -> OPEN) between the post-create and
				// import refreshes. Only the two status leaves are ignored;
				// channel_id/type/client_id/rpo and the list size are still compared so a
				// Read regression cannot pass silently.
				ImportStateVerifyIgnore: []string{
					"channels.0.channel_status",
					"channels.1.channel_status",
				},
			},
		},
	})
}

func resourceOtsTunnelConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}
	resource "alicloud_ots_table" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }
	  time_to_live = -1
	  max_version = 1
	}
	`, name)
}

var otsTunnelBasicMap = map[string]string{}

func TestAccAliCloudOtsTunnel_stream(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)

	// This case covers the flag-based stream configurations only (start_offset_flag
	// EARLIEST and LATEST); explicit start_offset/end_offset is exercised by
	// TestAccAliCloudOtsTunnel_streamOffset, which computes offsets at the
	// narrow valid window required by the OTS server.

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkOtsTunnelDestroy(rac.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				// Omitting the flag on the initial create exercises the server-default
				// Optional+Computed path and its post-apply consistency check.
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   "Stream",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":       "Stream",
						"start_offset_flag": "LATEST",
					}),
				),
			},
			{
				// Import only reads from DescribeTunnel, so ImportStateVerify proves that
				// Read backfills the server default LATEST into the imported state.
				// channel_status is server-side computed and observed to transition from
				// WAIT to OPEN between the post-create and import refreshes.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": "LATEST",
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":       "Stream",
						"start_offset_flag": "LATEST",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
			{
				// Removing explicit LATEST must be a no-op rather than a replacement
				// followed by an empty post-apply plan.
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": acctest.REMOVEKEY,
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionNoop),
					},
				},
			},
			{
				// The post-apply plan only proves the omitted configuration converges.
				// ImportStateVerify after this successful step proves the server default
				// LATEST is backfilled into an imported state; channel_status is the only
				// ignored leaf for the same asynchronous reason as above.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
			{
				// EARLIEST is non-default: it must replace the tunnel (ForceNew) and be
				// backfilled by Read.
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": "EARLIEST",
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":       "Stream",
						"start_offset_flag": "EARLIEST",
					}),
				),
			},
			{
				// channel_status is server-side computed and observed to transition from
				// WAIT to OPEN between the post-create and import refreshes.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": acctest.REMOVEKEY,
				}),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceId, plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":       "Stream",
						"start_offset_flag": "LATEST",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
		},
	})
}

// enableOtsTableStreamForAcc enables Stream on the OTS table created by
// resourceOtsTunnelConfigDependence with the maximum Stream expiration (168h =
// 7 days). The OTS server sets Stream expiration (SE) when Stream is first
// enabled on a table and does not allow it to be modified afterwards (docs
// 108059/102491). The maximum SE keeps the valid StartOffset window
// [serverNow - SE + 5m, serverNow) wide enough that offsets captured at test
// start remain valid when the explicit-offset tunnel is created later. The
// alicloud_ots_table resource does not expose stream_specification, so Stream
// is enabled via a direct UpdateTable call (the Go SDK UpdateTableRequest
// supports StreamSpec, doc 101592).
//
// After UpdateTable returns, the helper polls DescribeTable until the server
// reports EnableStream=true with a positive ExpirationTime. The OTS server may
// take up to a minute to apply the UpdateTable setting; creating the Stream
// tunnel before Stream is fully enabled would auto-enable it with a minimal
// expiration, shrinking the valid offset window and rejecting the explicit
// offsets (this is the failure observed in earlier revisions where the tunnel
// was created immediately after UpdateTable).
func enableOtsTableStreamForAcc(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		instanceName := "tf-" + name
		tableName := name
		request := new(tablestore.UpdateTableRequest)
		request.TableName = tableName
		request.StreamSpec = &tablestore.StreamSpecification{
			EnableStream:   true,
			ExpirationTime: 168,
		}
		if err := retry.Retry(3*time.Minute, func() *retry.RetryError {
			_, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.UpdateTable(request)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, errs.OtsTableIsTemporarilyUnavailable) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		}); err != nil {
			return err
		}
		// Poll DescribeTable until the server reports Stream enabled with a
		// positive expiration, so the subsequent CreateTunnel sees the large SE
		// rather than auto-enabling Stream with a minimal expiration.
		describeRequest := new(tablestore.DescribeTableRequest)
		describeRequest.TableName = tableName
		return retry.Retry(2*time.Minute, func() *retry.RetryError {
			raw, err := client.WithTableStoreClient(instanceName, func(tableStoreClient *tablestore.TableStoreClient) (interface{}, error) {
				return tableStoreClient.DescribeTable(describeRequest)
			})
			if err != nil {
				if errs.IsExpectedErrors(err, errs.OtsTableIsTemporarilyUnavailable) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			resp, ok := raw.(*tablestore.DescribeTableResponse)
			if !ok || resp == nil || resp.StreamDetails == nil {
				return retry.RetryableError(fmt.Errorf("stream details not ready on table %q", tableName))
			}
			if !resp.StreamDetails.EnableStream || resp.StreamDetails.ExpirationTime < 1 {
				return retry.RetryableError(fmt.Errorf("stream not yet enabled with a positive expiration on table %q", tableName))
			}
			return nil
		})
	}
}

func TestAccAliCloudOtsTunnel_streamOffset(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)

	// Explicit start_offset/end_offset require the table's Stream to be enabled
	// with a Stream expiration (SE) large enough that the offsets remain inside
	// the valid window [serverNow - SE + 5m, serverNow) when the tunnel is
	// created (doc 108059). SE is set when Stream is first enabled and cannot be
	// modified afterwards (docs 108059/102491), so step 1 enables Stream with the
	// maximum 168h SE via UpdateTable (and polls DescribeTable until the server
	// confirms Stream is enabled) before any tunnel exists; step 2 then creates
	// the Stream tunnel with explicit offsets. The offsets are captured at test
	// start and are ~1h in the past by the time step 2 applies, well inside the
	// 168h window and safely in the past relative to the server clock (earlier
	// revisions used a 2s offset which is fragile against clock skew). The
	// alicloud_ots_table resource does not expose stream_specification, which is
	// why UpdateTable is called directly in a check function.
	now := time.Now()
	// The schema exposes start_offset/end_offset as millisecond timestamps (Java
	// SDK doc 108059); the provider converts to nanoseconds when calling the
	// tablestore SDK (Go SDK doc 176301, which states the unit is nanoseconds;
	// confirmed empirically: millisecond values sent directly are rejected by
	// CreateTunnel as pointing to 1970). The offsets are ~1h in the past by the
	// time step 2 applies, well inside the 168h Stream expiration window and
	// safely in the past relative to the server clock (earlier revisions used a
	// 2s offset which was fragile against clock skew).
	startOffset := now.Add(-1 * time.Hour).UnixMilli()
	endOffset := now.Add(-30 * time.Minute).UnixMilli()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkOtsTunnelDestroy(rac.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				// Step 1: create the instance and table only (no tunnel), then
				// enable Stream on the table with a 168h expiration. SE is
				// immutable after first enable, so the maximum value must be set
				// now, before any tunnel is created. The helper then polls
				// DescribeTable until the server confirms Stream is enabled.
				Config: resourceOtsTunnelConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					enableOtsTableStreamForAcc(name),
				),
			},
			{
				// Step 2: create the Stream tunnel with explicit
				// start_offset/end_offset. Stream is already enabled with a 168h
				// SE, so the offsets (captured at test start, ~1h old by this
				// apply) are well inside the valid window. The Check below reads the
				// planned state and is therefore not proof of Read backfill; the
				// Import step below is, because an imported state can only come from
				// DescribeTunnel via Read.
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   "Stream",
					"start_offset":  strconv.FormatInt(startOffset, 10),
					"end_offset":    strconv.FormatInt(endOffset, 10),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":  "Stream",
						"start_offset": strconv.FormatInt(startOffset, 10),
						"end_offset":   strconv.FormatInt(endOffset, 10),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// channel_status is server-side computed by DescribeTunnel and can
				// transition asynchronously (WAIT -> OPEN) between the post-create
				// refresh and the import refresh. The flag-based Stream case ignores the
				// same leaf for the same reason. The offset and flag attributes, which
				// are the point of this case, are still verified.
				ImportStateVerifyIgnore: []string{"channels.0.channel_status"},
			},
		},
	})
}

func TestAccAliCloudOtsTunnel_nonStreamWithOffset(t *testing.T) {
	var v *otsTunnel.DescribeTunnelResponse

	resourceId := "alicloud_ots_tunnel.default"
	ra := acctest.ResourceAttrInit(resourceId, otsTunnelBasicMap)
	serviceFunc := func() interface{} {
		return ots.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)

	// Each step builds a fresh config (fresh attribute map) so that an offset
	// configured by an earlier rejected step is not carried into the next one.
	newConfig := func(tunnelType string, extra map[string]interface{}) string {
		attrs := map[string]interface{}{
			"instance_name": "${alicloud_ots_instance.default.name}",
			"table_name":    "${alicloud_ots_table.default.table_name}",
			"tunnel_name":   "${var.name}",
			"tunnel_type":   tunnelType,
		}
		for k, v := range extra {
			attrs[k] = v
		}
		return acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)(attrs)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             checkOtsTunnelDestroy(rac.CheckResourceDestroy()),
		Steps: []resource.TestStep{
			{
				// start_offset is not applicable to BaseAndStream tunnels and must be
				// rejected during planning, before a replacement can delete the old tunnel.
				Config:      newConfig("BaseAndStream", map[string]interface{}{"start_offset": "1600000000000"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only supported when tunnel_type is Stream`),
			},
			{
				// end_offset is likewise rejected for BaseAndStream tunnels.
				Config:      newConfig("BaseAndStream", map[string]interface{}{"end_offset": "1600000000000"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only supported when tunnel_type is Stream`),
			},
			{
				// An explicit LATEST flag is also configuration, so it must be rejected for
				// BaseData tunnels instead of being dropped by the flag DiffSuppressFunc.
				Config:      newConfig("BaseData", map[string]interface{}{"start_offset_flag": "LATEST"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only supported when tunnel_type is Stream`),
			},
			{
				Config:      newConfig("BaseData", map[string]interface{}{"start_offset_flag": "EARLIEST"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only supported when tunnel_type is Stream`),
			},
		},
	})
}

// TestAccAliCloudOtsTunnel_offsetValidation covers the plan-time checks of
// start_offset/end_offset and the offset/flag mutual exclusion. These
// configurations are rejected before any cloud resource is created, so the
// case needs no CheckDestroy.
func TestAccAliCloudOtsTunnel_offsetValidation(t *testing.T) {
	resourceId := "alicloud_ots_tunnel.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("testAcc%d", rand)

	newConfig := func(extra map[string]interface{}) string {
		attrs := map[string]interface{}{
			"instance_name": "${alicloud_ots_instance.default.name}",
			"table_name":    "${alicloud_ots_table.default.table_name}",
			"tunnel_name":   "${var.name}",
			"tunnel_type":   "Stream",
		}
		for k, v := range extra {
			attrs[k] = v
		}
		return acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)(attrs)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      newConfig(map[string]interface{}{"start_offset": "0"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be a positive millisecond timestamp`),
			},
			{
				Config:      newConfig(map[string]interface{}{"start_offset": "01600000000000"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`canonical decimal integer without leading zeros`),
			},
			{
				Config:      newConfig(map[string]interface{}{"end_offset": "not-a-number"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be a non-negative integer millisecond timestamp`),
			},
			{
				Config:      newConfig(map[string]interface{}{"end_offset": "9999999999999999999"}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`too large to convert to nanoseconds`),
			},
			{
				Config:      newConfig(map[string]interface{}{"start_offset": ""}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not be empty`),
			},
			{
				Config:      newConfig(map[string]interface{}{"end_offset": ""}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not be empty`),
			},
			{
				Config:      newConfig(map[string]interface{}{"tunnel_type": "BaseData", "start_offset": ""}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not be empty`),
			},
			{
				Config:      newConfig(map[string]interface{}{"tunnel_type": "BaseAndStream", "end_offset": ""}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must not be empty`),
			},
			{
				Config: newConfig(map[string]interface{}{
					"start_offset":      "1600000000000",
					"start_offset_flag": "EARLIEST",
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
			{
				Config: newConfig(map[string]interface{}{
					"start_offset": "1600000000000",
					"end_offset":   "1600000000000",
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`end_offset must be greater than start_offset`),
			},
			{
				Config: newConfig(map[string]interface{}{
					"start_offset": "1600000000001",
					"end_offset":   "1600000000000",
				}),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`end_offset must be greater than start_offset`),
			},
		},
	})
}
