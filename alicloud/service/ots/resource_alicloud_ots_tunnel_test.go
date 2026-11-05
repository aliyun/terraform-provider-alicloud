package ots_test

import (
	"fmt"
	"regexp"
	"strconv"
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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

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
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				// channel_status is server-side computed and transitions (e.g. WAIT -> OPEN)
				// once a BaseAndStream tunnel begins processing its base-data full table
				// scan, so it drifts between the post-create and import refreshes. The
				// tunnel attributes exercised here (instance_name, table_name, tunnel_name,
				// tunnel_type) are still verified.
				ImportStateVerifyIgnore: []string{
					"channels",
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
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Stream tunnel with start_offset_flag only (StartOffset not set), which
				// makes the SDK send the EARLIEST flag.
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": "EARLIEST",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tunnel_type":       "Stream",
						"start_offset_flag": "EARLIEST",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Explicit start_offset_flag = "LATEST" (the server default). Read does not
				// backfill LATEST, so the schema DiffSuppressFunc must suppress the empty-state
				// vs "LATEST"-config diff; otherwise every plan would force a replace. The
				// following ImportStateVerify re-runs the plan and asserts no diff, confirming
				// no permadiff.
				Config: testAccConfig(map[string]interface{}{
					"instance_name":     "${alicloud_ots_instance.default.name}",
					"table_name":        "${alicloud_ots_table.default.table_name}",
					"tunnel_name":       "${var.name}",
					"tunnel_type":       "Stream",
					"start_offset_flag": "LATEST",
				}),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
		CheckDestroy:             rac.CheckResourceDestroy(),
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
				// apply) are well inside the valid window. The post-apply refresh
				// verifies that Read backfills the offsets from DescribeTunnel.
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
				// channel_status is a server-side computed field that transitions
				// (e.g. WAIT -> OPEN) once the explicit-offset Stream tunnel begins
				// processing its offset range, so it is not stable between the
				// post-create refresh and the import refresh. The flag-based Stream
				// tunnel in TestAccAliCloudOtsTunnel_stream keeps channel_status stable
				// (no offset range to process) and therefore does not need this ignore.
				// The offset and flag attributes, which are the point of this case,
				// are still verified.
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
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceOtsTunnelConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// start_offset is not applicable to BaseAndStream tunnels and must be
				// rejected at create time rather than silently dropped.
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"table_name":    "${alicloud_ots_table.default.table_name}",
					"tunnel_name":   "${var.name}",
					"tunnel_type":   "BaseAndStream",
					"start_offset":  "1600000000000",
				}),
				ExpectError: regexp.MustCompile(`only supported when tunnel_type is Stream`),
			},
		},
	})
}
