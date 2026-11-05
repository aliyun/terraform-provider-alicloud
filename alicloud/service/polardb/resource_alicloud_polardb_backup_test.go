package polardb_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// polardbBackupTestRegion returns the region the provider actually runs in. acctest.PreCheck
// resolves an unset ALICLOUD_REGION to cn-beijing before the steps run, so reading the variable
// at test-function entry can disagree with the provider and make the Advanced policy quartet
// invalid ("400 DBS.ParamIsInValid: srcRegion should be same with db region"). Config builders
// call this after PreCheck, so the value matches the provider.
func polardbBackupTestRegion() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return acctest.DefaultRegionToTest
}

// TestAccAliCloudPolarDBBackup_basic exercises the property_completeness [M-RT-0040]
// mappings that apply at the default (Normal) backup-policy level: CreateBackup.Comment
// and ModifyBackupPolicy.BackupPolicyLevel/PreferredBackupPeriod/PreferredBackupTime. It
// also covers the collection-identity requirements: two backups of the same cluster must
// each own their own backup set, an in-place policy update must not move the identity, and
// destroying one backup must leave its sibling untouched. RestoreTable [M-RT-0039] is a
// separate model-level operation and is deliberately not part of this update contract.
func TestAccAliCloudPolarDBBackup_basic(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	siblingId := "alicloud_polardb_backup.sibling"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolardbBackup%s%d", polardbBackupTestRegion(), rand)

	var primaryID string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBBackupDestroy,
		Steps: []resource.TestStep{
			{
				// Two backups in one cluster: each resource must resolve and keep its own
				// backup set instead of the newest (or comment-matched) one in the list.
				Config: testAccAliCloudPolarDBBackupPair(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(siblingId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "backup_job_id"),
					resource.TestCheckResourceAttrSet(siblingId, "backup_job_id"),
					resource.TestCheckResourceAttr(resourceId, "comment", name),
					resource.TestCheckResourceAttr(siblingId, "comment", name+"sibling"),
					testAccCapturePolarDBBackupID(resourceId, &primaryID),
					testAccCheckPolarDBBackupDistinct(resourceId, siblingId),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Normal"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_period", "Wednesday,Tuesday"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_time", "10:00Z-11:00Z"),
				),
			},
			{
				// In-place policy update: the backup set identity must not move, and the
				// sibling backup must stay readable under its own identity.
				Config: testAccAliCloudPolarDBBackupPairUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolarDBBackupIDUnchanged(resourceId, &primaryID),
					testAccCheckPolarDBBackupExists(resourceId),
					testAccCheckPolarDBBackupExists(siblingId),
					resource.TestCheckResourceAttr(resourceId, "comment", name),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Normal"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_period", "Monday,Wednesday,Saturday"),
					resource.TestCheckResourceAttr(resourceId, "preferred_backup_time", "15:00Z-16:00Z"),
				),
			},
			{
				// An explicitly empty comment must stay "" in state: pinning it to null made it
				// differ from the configuration forever and, with RequiresReplaceIfConfigured,
				// replaced the backup set on every plan (P2). The step itself only plans
				// cleanly if the empty comment converges; the implicit post-apply plan check
				// fails otherwise.
				Config: testAccAliCloudPolarDBBackupEmptyComment(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "comment", ""),
					testAccCapturePolarDBBackupID(resourceId, &primaryID),
					testAccCheckPolarDBBackupExists(resourceId),
					testAccCheckPolarDBBackupExists(siblingId),
					testAccCheckPolarDBBackupDistinct(resourceId, siblingId),
				),
			},
			{
				// Import must reproduce the API-readable fields of the owned backup set.
				// Only the CreateBackup task fields (not part of any backup-set response)
				// are ignored, including backup_job_id: an imported backup has no create task.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_job_id", "backup_progress_status", "process", "create_time", "job_mode", "task_action"},
			},
			{
				// A missing backup identity must not fall back to the sibling backup that
				// does exist in the same cluster: importing an unknown BackupId has to fail.
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateIdFunc: testAccPolarDBBackupMissingIDFunc(siblingId),
				ExpectError:       regexp.MustCompile(`(?i)(cannot import|not exist|not found)`),
			},
			{
				// advanced_data_policies is only valid under BackupPolicyLevel=Advanced. Under
				// Normal the block is ignored locally instead of being sent together with the
				// Normal policy fields, which the server would reject (P3). The configured
				// entry stays in state, and the cluster policy must not gain that chain.
				Config: testAccAliCloudPolarDBBackupNormalWithAdvancedBlock(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Normal"),
					resource.TestCheckResourceAttr(resourceId, "advanced_data_policies.#", "1"),
					testAccCheckPolarDBBackupNoAdvancedPolicy(resourceId),
				),
			},
			{
				// The primary is removed from the configuration while the sibling stays:
				// the destroyed backup set must be gone and the sibling must survive. A
				// "first/last entry wins" read or a mis-targeted Delete fails here.
				Config: testAccAliCloudPolarDBBackupSiblingOnly(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolarDBBackupIDGone(&primaryID),
					resource.TestCheckResourceAttrSet(siblingId, "id"),
					testAccCheckPolarDBBackupExists(siblingId),
				),
			},
		},
	})
}

// TestAccAliCloudPolarDBBackup_advancedDataPolicies exercises the AdvancedDataPolicies
// mapping (ModifyBackupPolicy.AdvancedDataPolicies, 17 sub-fields) under the Advanced
// backup-policy level. AdvancedDataPolicies is only meaningful with BackupPolicyLevel =
// Advanced, which requires the cluster to have advanced backup enabled (申请使用, per
// the BackupPolicyLevel API contract). If the test account/cluster lacks that capability
// the case is reported as a capability gate rather than an implementation defect.
func TestAccAliCloudPolarDBBackup_advancedDataPolicies(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolardbBackupAdv%s%d", polardbBackupTestRegion(), rand)

	// G3: the Advanced in-place update must keep the backup-set identity. Capture the full
	// composite id before the 7 -> 14 retention change and compare it afterwards; the Normal
	// identity assertion in the basic case cannot stand in for this path.
	var advancedID string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBBackupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBBackupAdvanced(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "backup_id"),
					resource.TestCheckResourceAttrSet(resourceId, "backup_job_id"),
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Advanced"),
					testAccCapturePolarDBBackupID(resourceId, &advancedID),
					// QA 45954: Advanced policies require the full src_type/dest_type/
					// src_region/dest_region quartet; assert the configured db->level1
					// chain instead of a single src_type entry. The check matches by
					// chain so it is independent of API ordering and of the
					// service-generated policies.
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, polardbBackupTestRegion(), "7"),
				),
			},
			{
				// G3: same backup set, retention changed in place through ModifyBackupPolicy.
				Config: testAccAliCloudPolarDBBackupAdvancedUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "backup_policy_level", "Advanced"),
					testAccCheckPolarDBBackupIDUnchanged(resourceId, &advancedID),
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, polardbBackupTestRegion(), "14"),
				),
			},
			{
				// G3: a plan after apply must be empty. A PlanOnly step fails on any diff, so
				// this asserts the read-back of the updated policy converges instead of merely
				// accepting whichever value the API returns.
				Config:   testAccAliCloudPolarDBBackupAdvancedUpdate(name),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Only values the API cannot reproduce are ignored:
				// - the CreateBackup task fields are not part of any backup-set response;
				// - preferred_backup_period has no read-back under BackupPolicyLevel=Advanced
				//   (doc 2319233: 开启高级备份后，该参数不再生效; the API reports all seven days);
				// - action_type is the write-side operation of one configured entry and is never
				//   returned by DescribeBackupPolicy. Every other advanced_data_policies field,
				//   including the 16 API-readable ones, is compared (PRR-05).
				ImportStateVerifyIgnore: []string{"backup_job_id", "backup_progress_status", "process", "create_time", "job_mode", "task_action",
					"preferred_backup_period",
					"advanced_data_policies.0.action_type"},
			},
		},
	})
}

// TestAccAliCloudPolarDBBackup_releasedCluster covers P1: after the cluster is released,
// DescribeBackupPolicy reports it as absent, so Read must decide the fate of the backup set from
// the backup lists instead of clearing state on the policy NotFound alone.
//
// The cluster is released outside Terraform (DeleteDBCluster with
// BackupRetentionPolicyOnClusterDeletion=ALL); a refresh-only step then checks that the state
// follows the API: while DescribeBackups or DescribeDetachedBackups still report the exact
// backup set the resource must stay tracked under the same identity (and, when the set is only
// reported among a released cluster's backups, from the detached read-back fields), and when
// neither source reports it the resource must be removed instead of staying tracked.
func TestAccAliCloudPolarDBBackup_releasedCluster(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	clusterId := "alicloud_polardb_cluster.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolardbBackupReleased%s%d", polardbBackupTestRegion(), rand)

	var backupID, dbClusterID string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBBackupReleasedDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPolarDBBackupRetained(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "backup_retention_policy_on_cluster_deletion", "ALL"),
					testAccCapturePolarDBBackupID(resourceId, &backupID),
					testAccCapturePolarDBClusterID(clusterId, &dbClusterID),
				),
			},
			{
				PreConfig: func() {
					testAccReleasePolarDBCluster(t, dbClusterID)
					// G1: the positive detached branch is the whole point of this case. The
					// released-cluster backup list is populated asynchronously (retained backups
					// have been observed 10-25 minutes after the release), so wait for the exact
					// backup set to appear instead of refreshing during the population lag and
					// silently exercising only the "source did not answer" branch.
					// The captured identity is the composite DBClusterId:BackupId, while the
					// detached list reports the bare BackupId.
					_, waitBackupID := testAccSplitPolarDBBackupID(backupID)
					testAccWaitPolarDBDetachedBackupVisible(t, dbClusterID, waitBackupID, 30*time.Minute)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolarDBBackupReleasedState(resourceId, &backupID),
				),
			},
		},
	})
}

// testAccCheckPolarDBBackupDestroy fails while any backup set tracked by Terraform still
// exists. The lookup is by the exact composite identity, so a Delete that silently
// no-op'd or removed a different backup set is caught.
func testAccCheckPolarDBBackupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_backup" {
			continue
		}
		dbClusterID, backupID := testAccSplitPolarDBBackupID(rs.Primary.ID)
		if dbClusterID == "" || backupID == "" {
			return fmt.Errorf("polardb backup destroy check: malformed composite id %q", rs.Primary.ID)
		}
		if err := testAccWaitPolarDBBackupGone(dbClusterID, backupID, 3*time.Minute); err != nil {
			return err
		}
	}
	return nil
}

// testAccCheckPolarDBBackupReleasedDestroy is the destroy check of the released-cluster
// scenario. It fails while the backup set is still listed for a running cluster. A set that is
// only reported among a released cluster's backups cannot be deleted through DeleteBackup (the
// API requires a running cluster), so destroy only stops tracking it in that case.
func testAccCheckPolarDBBackupReleasedDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_backup" {
			continue
		}
		dbClusterID, backupID := testAccSplitPolarDBBackupID(rs.Primary.ID)
		if dbClusterID == "" || backupID == "" {
			return fmt.Errorf("polardb backup destroy check: malformed composite id %q", rs.Primary.ID)
		}
		live, err := testAccPolarDBBackupExists(dbClusterID, backupID)
		if err != nil {
			return fmt.Errorf("describing backups for cluster %s backup %s failed: %s", dbClusterID, backupID, err)
		}
		if live {
			return fmt.Errorf("polardb backup set %s for cluster %s is still listed by DescribeBackups after destroy", backupID, dbClusterID)
		}
	}
	return nil
}

// testAccWaitPolarDBBackupGone polls the exact backup identity until it disappears, so the
// check tolerates the short list-consistency delay after DeleteBackup.
func testAccWaitPolarDBBackupGone(dbClusterID, backupID string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		found, err := testAccPolarDBBackupExists(dbClusterID, backupID)
		if err != nil {
			return fmt.Errorf("describing backups for cluster %s backup %s failed: %s", dbClusterID, backupID, err)
		}
		if !found {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("polardb backup set %s for cluster %s still exists after destroy", backupID, dbClusterID)
		}
		time.Sleep(5 * time.Second)
	}
}

func testAccSplitPolarDBBackupID(id string) (dbClusterID, backupID string) {
	if idx := strings.Index(id, ":"); idx >= 0 {
		return id[:idx], id[idx+1:]
	}
	return id, ""
}

// testAccPolarDBBackupExists reports whether the exact backup set is currently listed.
//
// DescribeBackups answers very wide StartTime/EndTime ranges with an empty list (this package
// relies on that fact and walks 7-day windows), and its server-side BackupId filter returned an
// empty list for Advanced-policy clusters in the recorded observation. The check therefore queries a bounded
// 7-day window without the BackupId filter and matches the exact backup id in Go.
func testAccPolarDBBackupExists(dbClusterID, backupID string) (bool, error) {
	return testAccPolarDBBackupExistsWithClient(acctest.Provider.Meta().(*connectivity.AliyunClient), dbClusterID, backupID, 50)
}

// testAccPolarDBBackupExistsWithClient pages DescribeBackups to the end before reporting "not
// listed". An unfinished scan is never a NotFound: a target in a later page would otherwise make
// the existence check and CheckDestroy pass spuriously (PRR-04).
func testAccPolarDBBackupExistsWithClient(client *connectivity.AliyunClient, dbClusterID, backupID string, maxPages int) (bool, error) {
	const pageSize = 100
	now := time.Now().UTC()
	for page := 1; page <= maxPages; page++ {
		request := map[string]interface{}{
			"DBClusterId": dbClusterID,
			"StartTime":   now.Add(-7 * 24 * time.Hour).Format("2006-01-02T15:04Z"),
			"EndTime":     now.Add(time.Hour).Format("2006-01-02T15:04Z"),
			"PageNumber":  page,
			"PageSize":    pageSize,
			"RegionId":    client.RegionId,
		}
		response, err := client.RpcPost("polardb", "2017-08-01", "DescribeBackups", nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidDBCluster.NotFound", "InvalidBackupSetID.NotFound"}) {
				return false, nil
			}
			return false, err
		}
		items, _ := response["Items"].(map[string]interface{})
		backups, _ := items["Backup"].([]interface{})
		for _, item := range backups {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(entry["BackupId"]) != backupID {
				continue
			}
			if entryClusterID := fmt.Sprint(entry["DBClusterId"]); entryClusterID != "" && entryClusterID != "<nil>" && entryClusterID != dbClusterID {
				continue
			}
			return true, nil
		}
		if len(backups) < pageSize {
			return false, nil
		}
	}
	return false, fmt.Errorf("DescribeBackups did not finish paging within %d pages while looking for backup set %s of cluster %s; the scan is incomplete, so existence is unknown", maxPages, backupID, dbClusterID)
}

// testAccPolarDBBackupMissingIDFunc builds an import id whose BackupId does not exist while
// the cluster (and a sibling backup) does, so import has to fail instead of binding another
// backup set of the same cluster.
func testAccPolarDBBackupMissingIDFunc(resourceId string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return "", fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		dbClusterID, _ := testAccSplitPolarDBBackupID(rs.Primary.ID)
		return dbClusterID + ":999999999999", nil
	}
}

func testAccCapturePolarDBBackupID(resourceId string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		*target = rs.Primary.ID
		return nil
	}
}

func testAccCheckPolarDBBackupIDUnchanged(resourceId string, expected *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		if *expected == "" {
			return fmt.Errorf("polardb backup %s: no backup identity was captured", resourceId)
		}
		if rs.Primary.ID != *expected {
			return fmt.Errorf("polardb backup %s identity changed from %s to %s", resourceId, *expected, rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckPolarDBBackupIDGone(expected *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *expected == "" {
			return fmt.Errorf("no backup identity was captured for the destroyed resource")
		}
		dbClusterID, backupID := testAccSplitPolarDBBackupID(*expected)
		return testAccWaitPolarDBBackupGone(dbClusterID, backupID, 3*time.Minute)
	}
}

// testAccCheckPolarDBBackupExists fails when the resource's own backup set is not listed.
func testAccCheckPolarDBBackupExists(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		dbClusterID, backupID := testAccSplitPolarDBBackupID(rs.Primary.ID)
		if dbClusterID == "" || backupID == "" {
			return fmt.Errorf("polardb backup %s has malformed composite id %q", resourceId, rs.Primary.ID)
		}
		found, err := testAccPolarDBBackupExists(dbClusterID, backupID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("polardb backup set %s for cluster %s is not listed", backupID, dbClusterID)
		}
		return nil
	}
}

// testAccCheckPolarDBBackupDistinct asserts the two resources own different backup sets.
func testAccCheckPolarDBBackupDistinct(leftId, rightId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		left, ok := s.RootModule().Resources[leftId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", leftId)
		}
		right, ok := s.RootModule().Resources[rightId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", rightId)
		}
		_, leftBackupID := testAccSplitPolarDBBackupID(left.Primary.ID)
		_, rightBackupID := testAccSplitPolarDBBackupID(right.Primary.ID)
		if leftBackupID == "" || rightBackupID == "" {
			return fmt.Errorf("missing backup set id (left %q, right %q)", left.Primary.ID, right.Primary.ID)
		}
		if leftBackupID == rightBackupID {
			return fmt.Errorf("both resources track the same backup set %s", leftBackupID)
		}
		return nil
	}
}

// testAccCapturePolarDBClusterID records the cluster id of a dependency resource.
func testAccCapturePolarDBClusterID(resourceId string, target *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb cluster %s not found in state", resourceId)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("polardb cluster %s has no id", resourceId)
		}
		*target = rs.Primary.ID
		return nil
	}
}

// testAccCheckPolarDBBackupReleasedState asserts that the refreshed state follows the API after
// the cluster was released. The cluster was released with
// backup_retention_policy_on_cluster_deletion=ALL, so the backup set survives its cluster by
// contract and must stay tracked under the same identity: an empty detached listing is a
// visibility delay, never an absence proof. When the detached list reports the set, the
// backup-set fields (including store_status, G2) must be readable.
func testAccCheckPolarDBBackupReleasedState(resourceId string, expectedID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *expectedID == "" {
			return fmt.Errorf("polardb backup %s: no backup identity was captured", resourceId)
		}
		dbClusterID, backupID := testAccSplitPolarDBBackupID(*expectedID)
		var live, detached, inconclusive, tracked bool
		var rs *terraform.ResourceState
		for attempt := 0; ; attempt++ {
			var err error
			live, err = testAccPolarDBBackupExists(dbClusterID, backupID)
			if err != nil {
				return fmt.Errorf("describing backups for cluster %s backup %s failed: %s", dbClusterID, backupID, err)
			}
			detached, err = testAccPolarDBDetachedBackupExists(dbClusterID, backupID)
			inconclusive = err != nil && helper.NeedRetry(err)
			if err != nil && !inconclusive {
				return fmt.Errorf("listing detached backups for cluster %s failed: %s", dbClusterID, err)
			}
			rs, tracked = s.RootModule().Resources[resourceId]
			if tracked == (live || detached) || attempt >= 6 {
				break
			}
			time.Sleep(10 * time.Second)
		}
		if !tracked {
			return fmt.Errorf("polardb backup %s was removed from state although backup_retention_policy_on_cluster_deletion=ALL keeps backup set %s (live=%t detached=%t inconclusive=%t)", resourceId, backupID, live, detached, inconclusive)
		}
		if rs.Primary.ID != *expectedID {
			return fmt.Errorf("polardb backup identity changed from %s to %s", *expectedID, rs.Primary.ID)
		}
		attrs := rs.Primary.Attributes
		if got := attrs["backup_id"]; got != backupID {
			return fmt.Errorf("polardb backup %s reports backup_id %q, expected %q", resourceId, got, backupID)
		}
		if got := attrs["status"]; got != "Success" {
			return fmt.Errorf("polardb backup %s reports status %q, expected Success", resourceId, got)
		}
		if detached {
			// G2: the released-cluster list is the only source that reports StoreStatus, so the
			// positive branch must carry it into state instead of leaving it empty.
			if got := attrs["store_status"]; got == "" {
				return fmt.Errorf("polardb backup %s has an empty store_status although DescribeDetachedBackups reports the set", resourceId)
			}
		}
		if got := attrs["backup_retention_policy_on_cluster_deletion"]; got != "ALL" {
			return fmt.Errorf("polardb backup %s lost its cluster policy fields: backup_retention_policy_on_cluster_deletion is %q", resourceId, got)
		}
		return nil
	}
}

// testAccWaitPolarDBDetachedBackupVisible waits until DescribeDetachedBackups reports the exact
// backup set of a released cluster, so the positive detached read branch can be exercised. The
// released-cluster list is populated asynchronously after a release; failing here means the
// environment could not exercise the branch within the documented observation window.
func testAccWaitPolarDBDetachedBackupVisible(t *testing.T, dbClusterID, backupID string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		found, err := testAccPolarDBDetachedBackupExists(dbClusterID, backupID)
		if err == nil && found {
			return
		}
		if err != nil && !helper.NeedRetry(err) {
			t.Fatalf("listing detached backups for cluster %s failed: %s", dbClusterID, err)
		}
		if time.Now().After(deadline) {
			t.Fatalf("DescribeDetachedBackups did not report backup set %s of released cluster %s within %s, so the released-cluster read branch could not be exercised", backupID, dbClusterID, timeout)
		}
		time.Sleep(30 * time.Second)
	}
}

// testAccCheckPolarDBBackupNoAdvancedPolicy asserts that the cluster's backup policy does not
// contain the configured db->level1 advanced chain. It proves that an advanced_data_policies
// block under BackupPolicyLevel=Normal was ignored instead of being applied.
func testAccCheckPolarDBBackupNoAdvancedPolicy(resourceId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		dbClusterID := rs.Primary.Attributes["db_cluster_id"]
		if dbClusterID == "" {
			return fmt.Errorf("polardb backup %s has no db_cluster_id", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		request := map[string]interface{}{"DBClusterId": dbClusterID, "RegionId": client.RegionId}
		response, err := client.RpcPost("polardb", "2017-08-01", "DescribeBackupPolicy", nil, request, true)
		if err != nil {
			return fmt.Errorf("describing the backup policy of cluster %s failed: %s", dbClusterID, err)
		}
		container, _ := response["AdvancedDataPolicies"].(map[string]interface{})
		policies, _ := container["AdvancedDataPolicy"].([]interface{})
		for _, item := range policies {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(entry["SrcType"]) == "db" && fmt.Sprint(entry["DestType"]) == "level1" {
				return fmt.Errorf("cluster %s reports an advanced db->level1 policy although backup_policy_level is Normal", dbClusterID)
			}
		}
		return nil
	}
}

// testAccReleasePolarDBCluster releases a cluster outside Terraform while retaining all of its
// backups, then waits until the cluster is gone from the cluster API.
func testAccReleasePolarDBCluster(t *testing.T, dbClusterID string) {
	t.Helper()
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"DBClusterId":                            dbClusterID,
		"BackupRetentionPolicyOnClusterDeletion": "ALL",
		"RegionId":                               client.RegionId,
	}
	if _, err := client.RpcPost("polardb", "2017-08-01", "DeleteDBCluster", nil, request, true); err != nil {
		t.Fatalf("releasing cluster %s failed: %s", dbClusterID, err)
	}
	deadline := time.Now().Add(20 * time.Minute)
	for {
		_, err := client.RpcPost("polardb", "2017-08-01", "DescribeDBClusterAttribute", nil, map[string]interface{}{
			"DBClusterId": dbClusterID,
			"RegionId":    client.RegionId,
		}, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
				return
			}
			if !helper.NeedRetry(err) {
				t.Fatalf("describing cluster %s while waiting for its release failed: %s", dbClusterID, err)
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("cluster %s was not released within 20 minutes", dbClusterID)
		}
		time.Sleep(15 * time.Second)
	}
}

// testAccPolarDBDetachedBackupExists reports whether the exact backup set is listed among the
// detached backups of a released cluster. Transient 5xx answers (the backend can fail while the
// release is still settling) are retried; the window end is now because a future EndTime is
// rejected with 500 InternalFailure by DescribeDetachedBackups.
func testAccPolarDBDetachedBackupExists(dbClusterID, backupID string) (bool, error) {
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		found, err := testAccPolarDBDetachedBackupExistsOnce(dbClusterID, backupID)
		if err == nil {
			return found, nil
		}
		if !helper.NeedRetry(err) {
			return false, err
		}
		lastErr = err
		time.Sleep(5 * time.Second)
	}
	return false, lastErr
}

func testAccPolarDBDetachedBackupExistsOnce(dbClusterID, backupID string) (bool, error) {
	const pageSize = 100
	const maxPages = 50
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	now := time.Now().UTC()
	for page := 1; page <= maxPages; page++ {
		request := map[string]interface{}{
			"DBClusterId": dbClusterID,
			"StartTime":   now.Add(-24 * time.Hour).Format("2006-01-02T15:04Z"),
			"EndTime":     now.Format("2006-01-02T15:04Z"),
			"PageNumber":  page,
			"PageSize":    pageSize,
			"RegionId":    client.RegionId,
		}
		response, err := client.RpcPost("polardb", "2017-08-01", "DescribeDetachedBackups", nil, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) {
				return false, nil
			}
			return false, err
		}
		items, _ := response["Items"].(map[string]interface{})
		backups, _ := items["Backup"].([]interface{})
		for _, item := range backups {
			entry, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(entry["BackupId"]) == backupID {
				return true, nil
			}
		}
		if len(backups) < pageSize {
			return false, nil
		}
	}
	return false, fmt.Errorf("DescribeDetachedBackups did not finish paging within %d pages while looking for backup set %s of cluster %s", maxPages, backupID, dbClusterID)
}

// testAccCheckPolarDBBackupAdvancedPolicy verifies that the configured db->level1
// advanced data policy is present with the expected retention. It matches by chain
// instead of list index so the assertion is independent of API ordering and of
// service-generated policies.
func testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, retentionValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("polardb backup %s not found in state", resourceId)
		}
		attrs := rs.Primary.Attributes
		count, err := strconv.Atoi(attrs["advanced_data_policies.#"])
		if err != nil {
			return fmt.Errorf("polardb backup %s has no readable advanced_data_policies count: %v", resourceId, err)
		}
		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("advanced_data_policies.%d.", i)
			if attrs[prefix+"src_type"] != "db" || attrs[prefix+"dest_type"] != "level1" ||
				attrs[prefix+"src_region"] != region || attrs[prefix+"dest_region"] != region {
				continue
			}
			if got := attrs[prefix+"bak_type"]; got != "F" {
				return fmt.Errorf("advanced policy db->level1 has bak_type %q, expected F", got)
			}
			if got := attrs[prefix+"retention_type"]; got != "delay" {
				return fmt.Errorf("advanced policy db->level1 has retention_type %q, expected delay", got)
			}
			if got := attrs[prefix+"retention_value"]; got != retentionValue {
				return fmt.Errorf("advanced policy db->level1 has retention_value %q, expected %s", got, retentionValue)
			}
			return nil
		}
		return fmt.Errorf("polardb backup %s has no db->level1 advanced policy for region %s among %d entries", resourceId, region, count)
	}
}

func testAccAliCloudPolarDBBackupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_polardb_node_classes" "default" {
  pay_type   = "PostPaid"
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
}

locals {
  // Use a class the queried zone actually offers instead of a hard-coded one. The first
  // entry is not always orderable (COMMODITY.INVALID_COMPONENT), so prefer the known
  // PostPaid-able class whenever the zone lists it and otherwise take the first offer.
  polardb_classes   = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.*.db_node_class
  polardb_preferred = [for c in local.polardb_classes : c if c == "polar.mysql.x4.large"]
  db_node_class     = length(local.polardb_preferred) > 0 ? local.polardb_preferred[0] : local.polardb_classes[0]
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  // Take the available class selected from the queried zone (local.db_node_class).
  db_node_class = local.db_node_class
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
  // encrypt_new_tables is Optional (not Computed) on the cluster resource and the
  // API returns "OFF" for a fresh cluster; set it explicitly so the post-apply plan
  // is empty (otherwise the fixture drifts: encrypt_new_tables = "OFF" -> null).
  encrypt_new_tables = "OFF"
}
`, name)
}

func testAccAliCloudPolarDBBackupPair(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id       = alicloud_polardb_cluster.default.id
  comment             = var.name
  backup_policy_level = "Normal"
  // The API echoes the configured weekday order (verified on cn-beijing), so the
  // assertion below follows it; the resource additionally compares weekday sets.
  preferred_backup_period = "Wednesday,Tuesday"
  preferred_backup_time   = "10:00Z-11:00Z"
}

resource "alicloud_polardb_backup" "sibling" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  comment       = "${var.name}sibling"

  // Serialize the two create calls: PolarDB runs one manual backup per cluster.
  depends_on = [alicloud_polardb_backup.default]
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupPairUpdate(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  comment                 = var.name
  backup_policy_level     = "Normal"
  preferred_backup_period = "Monday,Wednesday,Saturday"
  preferred_backup_time   = "15:00Z-16:00Z"
}

resource "alicloud_polardb_backup" "sibling" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  comment       = "${var.name}sibling"

  depends_on = [alicloud_polardb_backup.default]
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupSiblingOnly(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "sibling" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  comment       = "${var.name}sibling"
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupEmptyComment(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  comment                 = ""
  backup_policy_level     = "Normal"
  preferred_backup_period = "Monday,Wednesday,Saturday"
  preferred_backup_time   = "15:00Z-16:00Z"
}

resource "alicloud_polardb_backup" "sibling" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  comment       = "${var.name}sibling"

  depends_on = [alicloud_polardb_backup.default]
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupNormalWithAdvancedBlock(name string) string {
	region := polardbBackupTestRegion()
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  comment                 = ""
  backup_policy_level     = "Normal"
  preferred_backup_period = "Monday,Wednesday,Saturday"
  preferred_backup_time   = "15:00Z-16:00Z"

  // Only valid under BackupPolicyLevel=Advanced; under Normal the block must be
  // ignored locally instead of being sent with the Normal policy fields.
  advanced_data_policies {
    action_type     = "CREATE"
    bak_type        = "F"
    src_type        = "db"
    dest_type       = "level1"
    src_region      = "%s"
    dest_region     = "%s"
    retention_type  = "delay"
    retention_value = "7"
  }
}

resource "alicloud_polardb_backup" "sibling" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  comment       = "${var.name}sibling"

  depends_on = [alicloud_polardb_backup.default]
}
`, testAccAliCloudPolarDBBackupDependence(name), region, region)
}

func testAccAliCloudPolarDBBackupRetained(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id                               = alicloud_polardb_cluster.default.id
  comment                                     = var.name
  backup_retention_policy_on_cluster_deletion = "ALL"
}
`, testAccAliCloudPolarDBBackupDependence(name))
}

func testAccAliCloudPolarDBBackupAdvanced(name string) string {
	region := polardbBackupTestRegion()
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  backup_policy_level     = "Advanced"
  preferred_backup_period = "Tuesday,Wednesday"
  preferred_backup_time   = "10:00Z-11:00Z"

  // Advanced policies are validated server-side as a full chain: the request must
  // carry srcType, destType, srcRegion and destRegion together. db -> level1 in the
  // test region is the base chain documented for ModifyBackupPolicy.AdvancedDataPolicies,
  // and the scheduling fields (FilterType/FilterKey/FilterValue, OnlyPreserveOneEachDay)
  // are part of the policy definition: omitting them makes DBS reject the save with
  // DBS.Policy.BackupPolicy.SaveFailed.
  advanced_data_policies {
    action_type                = "CREATE"
    bak_type                   = "F"
    src_type                   = "db"
    dest_type                  = "level1"
    src_region                 = "%s"
    dest_region                = "%s"
    filter_type                = "crontab"
    filter_key                 = "dayOfWeek"
    filter_value               = "1,2,3,4,5,6,7"
    only_preserve_one_each_day = true
    retention_type             = "delay"
    retention_value            = "7"
  }
}
`, testAccAliCloudPolarDBBackupDependence(name), region, region)
}

func testAccAliCloudPolarDBBackupAdvancedUpdate(name string) string {
	region := polardbBackupTestRegion()
	return fmt.Sprintf(`%s

resource "alicloud_polardb_backup" "default" {
  db_cluster_id           = alicloud_polardb_cluster.default.id
  backup_policy_level     = "Advanced"
  preferred_backup_period = "Tuesday,Wednesday"
  preferred_backup_time   = "10:00Z-11:00Z"

  advanced_data_policies {
    action_type                = "UPDATE"
    bak_type                   = "F"
    src_type                   = "db"
    dest_type                  = "level1"
    src_region                 = "%s"
    dest_region                = "%s"
    filter_type                = "crontab"
    filter_key                 = "dayOfWeek"
    filter_value               = "1,2,3,4,5,6,7"
    only_preserve_one_each_day = true
    retention_type             = "delay"
    retention_value            = "14"
  }
}
`, testAccAliCloudPolarDBBackupDependence(name), region, region)
}

// The unit tests below drive the resource through its public entry points (Create/Read/
// Update/Delete) against an isolated HTTP endpoint, so they exercise the real request
// serialization, the real SDK decoding and the real state machine. They never touch the cloud
// and never install a client interface seam: the provider client is the production client with
// its endpoint pointed at a local server. Operation budgets are set through the public
// `timeouts` block, so the same code paths run without waiting for production delays.

const (
	backupUnitRegion    = "cn-beijing"
	backupUnitClusterID = "pc-unit-cluster"
	backupUnitBackupID  = "bk-unit-backup"
	backupUnitStartTime = "2026-10-20T00:00:00Z"
)

type backupUnitServer struct {
	t       *testing.T
	mu      sync.Mutex
	actions []string
	handler func(action string, params map[string]string) string
}

func (s *backupUnitServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	params := map[string]string{}
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	// The RPC client may send the API parameters in the POST body as JSON or as a form. The body
	// is read above, so it is decoded here explicitly instead of relying on ParseForm.
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if len(raw) > 0 {
		if strings.Contains(contentType, "json") {
			decoded := map[string]interface{}{}
			if json.Unmarshal(raw, &decoded) == nil {
				backupUnitFlattenParams("", decoded, params)
			}
		} else if values, err := url.ParseQuery(string(raw)); err == nil {
			for key := range values {
				params[key] = values.Get(key)
			}
		}
	}
	action := params["Action"]
	// A nested list parameter arrives as one JSON-encoded form value; expand it so the
	// assertions can address its leaves.
	for key, value := range params {
		trimmed := strings.TrimSpace(value)
		if !strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "{") {
			continue
		}
		decoded := interface{}(nil)
		if json.Unmarshal([]byte(trimmed), &decoded) == nil {
			backupUnitFlattenParams(key, decoded, params)
		}
	}
	s.mu.Lock()
	s.actions = append(s.actions, action)
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	body := "{}"
	if s.handler != nil {
		body = s.handler(action, params)
	}
	_, _ = w.Write([]byte(body))
}

// backupUnitFlattenParams flattens a decoded JSON body into dotted keys, so captured parameters
// are addressable regardless of whether the client nests entries or flattens them.
func backupUnitFlattenParams(prefix string, value interface{}, out map[string]string) {
	switch typed := value.(type) {
	case map[string]interface{}:
		for key, item := range typed {
			name := key
			if prefix != "" {
				name = prefix + "." + key
			}
			backupUnitFlattenParams(name, item, out)
		}
	case []interface{}:
		for index, item := range typed {
			backupUnitFlattenParams(fmt.Sprintf("%s.%d", prefix, index), item, out)
		}
	case nil:
		out[prefix] = ""
	default:
		out[prefix] = fmt.Sprint(typed)
	}
}

func (s *backupUnitServer) seen(action string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, seen := range s.actions {
		if seen == action {
			return true
		}
	}
	return false
}

func newBackupUnitClient(t *testing.T, handler func(action string, params map[string]string) string) (*connectivity.AliyunClient, *backupUnitServer) {
	t.Helper()
	server := &backupUnitServer{t: t, handler: handler}
	httpServer := httptest.NewServer(server)
	t.Cleanup(httpServer.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("building the test credential: %v", err)
	}
	endpoint := strings.TrimPrefix(httpServer.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             backupUnitRegion,
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("building the test client: %v", err)
	}
	endpoints.Store("polardb", endpoint)
	return client, server
}

func newBackupUnitResource(t *testing.T, client *connectivity.AliyunClient) (fwresource.Resource, fwschema.Schema) {
	t.Helper()
	r := polardb.NewBackupResource()
	configured, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("the backup resource does not accept the provider client")
	}
	var configureResp fwresource.ConfigureResponse
	configured.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building the schema: %v", schemaResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

func backupUnitObjectType(t *testing.T, sch fwschema.Schema) tftypes.Object {
	t.Helper()
	objectType, ok := sch.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("the resource schema does not describe an object")
	}
	return objectType
}

// backupUnitObject builds a fully typed object value: every attribute not listed is null.
func backupUnitObject(t *testing.T, objectType tftypes.Object, attributes map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	full := map[string]tftypes.Value{}
	for name, attributeType := range objectType.AttributeTypes {
		if value, ok := attributes[name]; ok {
			full[name] = value
			continue
		}
		full[name] = tftypes.NewValue(attributeType, nil)
	}
	return tftypes.NewValue(objectType, full)
}

func backupUnitString(value string) tftypes.Value {
	return tftypes.NewValue(tftypes.String, value)
}

func backupUnitBool(value bool) tftypes.Value {
	return tftypes.NewValue(tftypes.Bool, value)
}

// backupUnitTimeouts builds the public timeouts block for a lifecycle object. The resource
// bounds its create resolution and read grace by these operation timeouts, so the unit tests
// drive the production code paths without waiting for the default operation budgets.
func backupUnitTimeouts(t *testing.T, objectType tftypes.Object, create, read string) tftypes.Value {
	t.Helper()
	timeoutsType, ok := objectType.AttributeTypes["timeouts"].(tftypes.Object)
	if !ok {
		t.Fatal("the resource schema has no timeouts object")
	}
	attributes := map[string]tftypes.Value{}
	for name, attributeType := range timeoutsType.AttributeTypes {
		attributes[name] = tftypes.NewValue(attributeType, nil)
	}
	if create != "" {
		attributes["create"] = backupUnitString(create)
	}
	if read != "" {
		attributes["read"] = backupUnitString(read)
	}
	return tftypes.NewValue(timeoutsType, attributes)
}

// backupUnitPolicyList builds the advanced_data_policies block value from per-entry attributes.
func backupUnitPolicyList(t *testing.T, objectType tftypes.Object, entries []map[string]tftypes.Value) tftypes.Value {
	t.Helper()
	listType, ok := objectType.AttributeTypes["advanced_data_policies"].(tftypes.List)
	if !ok {
		t.Fatal("advanced_data_policies is not a list block")
	}
	elementType, ok := listType.ElementType.(tftypes.Object)
	if !ok {
		t.Fatal("advanced_data_policies entries are not objects")
	}
	values := make([]tftypes.Value, 0, len(entries))
	for _, entry := range entries {
		full := map[string]tftypes.Value{}
		for name, attributeType := range elementType.AttributeTypes {
			if value, ok := entry[name]; ok {
				full[name] = value
				continue
			}
			full[name] = tftypes.NewValue(attributeType, nil)
		}
		values = append(values, tftypes.NewValue(elementType, full))
	}
	return tftypes.NewValue(listType, values)
}

func backupUnitNotFound(code string) string {
	return fmt.Sprintf(`{"Code":%q,"Message":"unit test not found"}`, code)
}

func backupUnitBackupsJSON(entries ...string) string {
	if len(entries) == 0 {
		return `{"Items":{"Backup":[]},"TotalRecordCount":0}`
	}
	return fmt.Sprintf(`{"Items":{"Backup":[%s]},"TotalRecordCount":%d}`, strings.Join(entries, ","), len(entries))
}

// backupUnitBackupEntry builds a DescribeBackups entry with only the fields that action
// documents; StoreStatus belongs to DescribeDetachedBackups and is added by that fixture.
func backupUnitBackupEntry(backupID, clusterID, startTime string) string {
	return fmt.Sprintf(`{"BackupId":%q,"DBClusterId":%q,"BackupStatus":"Success","BackupStartTime":%q,"BackupMode":"Manual"}`,
		backupID, clusterID, startTime)
}

func backupUnitStateAttribute(t *testing.T, state tfsdk.State, name string) (string, bool) {
	t.Helper()
	var value types.String
	if diagnostics := state.GetAttribute(context.Background(), fwpath.Root(name), &value); diagnostics.HasError() {
		return "", false
	}
	if value.IsNull() || value.IsUnknown() {
		return "", false
	}
	return value.ValueString(), true
}

func backupUnitDiagnosticsText(diagnostics interface{ HasError() bool }) string {
	return fmt.Sprint(diagnostics)
}

// TestUnitAliCloudPolarDBBackupCreateRefusesForeignBackup covers PRR-01: the created backup set
// is bound only through the task id returned by CreateBackup, and a set that another actor
// created at a different time is never accepted as this task's backup set.
func TestUnitAliCloudPolarDBBackupCreateRefusesForeignBackup(t *testing.T) {
	foreign := backupUnitBackupEntry("bk-someone-else", backupUnitClusterID, "2026-10-19T23:58:00Z")
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			// Our task's backup set is never listed; another actor's set is.
			return backupUnitBackupsJSON(foreign)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-foreign"),
		"timeouts":      backupUnitTimeouts(t, objectType, "150ms", ""),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Create accepted a create whose own backup set never appeared: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeBackups") {
		t.Fatal("Create did not verify the created backup set against DescribeBackups")
	}
	if detail := backupUnitDiagnosticsText(resp.Diagnostics); !strings.Contains(detail, backupUnitBackupID) {
		t.Errorf("Create error does not name the backup task id it refused to bind: %s", detail)
	}
}

// TestUnitAliCloudPolarDBBackupCreateRefusesAmbiguousBackupSet covers the other side of the
// task-start-time binding: when two manual backup sets started at the task start time, the
// create must refuse instead of picking one.
func TestUnitAliCloudPolarDBBackupCreateRefusesAmbiguousBackupSet(t *testing.T) {
	first := backupUnitBackupEntry("bk-first", backupUnitClusterID, backupUnitStartTime)
	second := backupUnitBackupEntry("bk-second", backupUnitClusterID, backupUnitStartTime)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			return backupUnitBackupsJSON(first, second)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-ambiguous"),
		"timeouts":      backupUnitTimeouts(t, objectType, "150ms", ""),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Create bound one of two backup sets that started at the task start time: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if detail := backupUnitDiagnosticsText(resp.Diagnostics); !strings.Contains(detail, "backup sets starting at that second") {
		t.Errorf("Create error does not explain the ambiguous backup sets: %s", detail)
	}
}

// TestUnitAliCloudPolarDBBackupCreateBindsTaskBackupSet covers the positive branch: the task id
// from CreateBackup differs from the backup set id, and the set is bound because it is the only
// manual backup starting at the task start time.
func TestUnitAliCloudPolarDBBackupCreateBindsTaskBackupSet(t *testing.T) {
	created := backupUnitBackupEntry("3218271464", backupUnitClusterID, backupUnitStartTime)
	automated := backupUnitBackupEntry("3217277072", backupUnitClusterID, backupUnitStartTime)
	automated = strings.Replace(automated, `"BackupMode":"Manual"`, `"BackupMode":"Automated"`, 1)
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "CreateBackup":
			return fmt.Sprintf(`{"BackupJobId":%q}`, backupUnitBackupID)
		case "DescribeBackupTasks":
			return fmt.Sprintf(`{"Items":{"BackupJob":[{"BackupJobId":%q,"BackupProgressStatus":"Finished","StartTime":%q}]}}`, backupUnitBackupID, backupUnitStartTime)
		case "DescribeBackups":
			return backupUnitBackupsJSON(automated, created)
		case "DescribeBackupPolicy":
			return `{"DBClusterId":"` + backupUnitClusterID + `","BackupPolicyLevel":"Normal"}`
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"comment":       backupUnitString("unit-bind"),
		"timeouts":      backupUnitTimeouts(t, objectType, "2s", ""),
	})
	var resp fwresource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Create failed to bind the set produced by its task: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id != backupUnitClusterID+":3218271464" {
		t.Fatalf("Create did not record the resolved backup set id: %q present=%v", id, ok)
	}
	if jobId, ok := backupUnitStateAttribute(t, resp.State, "backup_job_id"); !ok || jobId != backupUnitBackupID {
		t.Fatalf("Create did not record the authoritative task id: %q present=%v", jobId, ok)
	}
}

// TestUnitAliCloudPolarDBBackupReadDetachedRetained covers G1/G2: after the cluster is released,
// Read resolves the retained backup set from DescribeDetachedBackups, keeps the resource and
// reports the StoreStatus the released-cluster list is the only source for.
func TestUnitAliCloudPolarDBBackupReadDetachedRetained(t *testing.T) {
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			// Only the fields DescribeDetachedBackups documents: the comment is not one of them,
			// so it must be preserved from state rather than taken from the response.
			return backupUnitBackupsJSON(`{"BackupId":"` + backupUnitBackupID + `","DBClusterId":"` + backupUnitClusterID + `","BackupStatus":"Success","StoreStatus":"Success","BackupStartTime":"` + backupUnitStartTime + `","IsAvail":"1"}`)
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("ALL"),
		"preferred_backup_time":                       backupUnitString("10:00Z-11:00Z"),
		"comment":                                     backupUnitString("unit-retained"),
		"timeouts":                                    backupUnitTimeouts(t, objectType, "", "1s"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed on the retained detached branch: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id != backupUnitClusterID+":"+backupUnitBackupID {
		t.Fatalf("Read dropped the retained backup set from state: id=%q present=%v", id, ok)
	}
	if storeStatus, ok := backupUnitStateAttribute(t, resp.State, "store_status"); !ok || storeStatus != "Success" {
		t.Fatalf("Read did not report store_status from DescribeDetachedBackups: %q present=%v", storeStatus, ok)
	}
	if retention, ok := backupUnitStateAttribute(t, resp.State, "backup_retention_policy_on_cluster_deletion"); !ok || retention != "ALL" {
		t.Fatalf("Read lost the cluster policy fields of a released cluster: %q present=%v", retention, ok)
	}
	if comment, ok := backupUnitStateAttribute(t, resp.State, "comment"); !ok || comment != "unit-retained" {
		t.Fatalf("Read did not map the detached backup-set fields: comment=%q present=%v", comment, ok)
	}
}

// TestUnitAliCloudPolarDBBackupReadKeepsRetainedAbsence covers N2: retention ALL keeps the backup
// set by contract, so an empty released-cluster listing must keep the state instead of removing a
// possibly billed backup.
func TestUnitAliCloudPolarDBBackupReadKeepsRetainedAbsence(t *testing.T) {
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy", "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON()
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("ALL"),
		"timeouts": backupUnitTimeouts(t, objectType, "", "1s"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed although the retention contract keeps the backup set: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if len(resp.Diagnostics.Warnings()) == 0 {
		t.Fatal("Read did not warn that the retained backup set is not listed yet")
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id == "" {
		t.Fatal("Read removed a backup set that retention ALL keeps")
	}
}

// TestUnitAliCloudPolarDBBackupReadRemovesUnretainedAbsence covers the other half of N2: when the
// retention policy does not keep the backup set, repeated complete empty answers remove it.
func TestUnitAliCloudPolarDBBackupReadRemovesUnretainedAbsence(t *testing.T) {
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy", "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON()
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("NONE"),
		"timeouts": backupUnitTimeouts(t, objectType, "", "150ms"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed while removing an unretained backup set: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeDetachedBackups") {
		t.Fatal("Read removed the backup set without consulting the released-cluster list")
	}
	if _, ok := backupUnitStateAttribute(t, resp.State, "id"); ok {
		t.Fatal("Read kept a backup set that is gone and that retention does not keep")
	}
}

// TestUnitAliCloudPolarDBBackupReadToleratesExpiredTaskRecord covers the tolerated missing-task
// case: DescribeBackupTasks keeps only running and recent tasks, so an older backup has no record.
// The read must still succeed, keep the tracked task id and report no error.
func TestUnitAliCloudPolarDBBackupReadToleratesExpiredTaskRecord(t *testing.T) {
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy":
			return fmt.Sprintf(`{"DBClusterId":%q,"BackupPolicyLevel":"Normal"}`, backupUnitClusterID)
		case "DescribeBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		case "DescribeBackupTasks":
			return `{"Items":{"BackupJob":[]}}`
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":            backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"backup_id":     backupUnitString(backupUnitBackupID),
		"backup_job_id": backupUnitString(backupUnitBackupID),
		"timeouts":      backupUnitTimeouts(t, objectType, "", "1s"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed although a missing task record is only the task-list retention window: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeBackupTasks") {
		t.Fatal("Read did not query the tracked backup task")
	}
	if jobId, ok := backupUnitStateAttribute(t, resp.State, "backup_job_id"); !ok || jobId != backupUnitBackupID {
		t.Fatalf("Read lost the tracked backup task id when the task record expired: %q present=%v", jobId, ok)
	}
	if warnings := resp.Diagnostics.Warnings(); len(warnings) != 0 {
		t.Fatalf("Read warned about an expired task record: %v", warnings)
	}
}

// TestUnitAliCloudPolarDBBackupReadReportsTaskQueryFailure covers the other side of the task
// query: permission, throttling and service failures must not look like a successful refresh.
// The read has to fail with a diagnostic that names the query instead of returning the state.
func TestUnitAliCloudPolarDBBackupReadReportsTaskQueryFailure(t *testing.T) {
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy":
			return fmt.Sprintf(`{"DBClusterId":%q,"BackupPolicyLevel":"Normal"}`, backupUnitClusterID)
		case "DescribeBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		case "DescribeBackupTasks":
			return `{"Code":"Forbidden.RAM","Message":"the caller is not authorized to perform DescribeBackupTasks"}`
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":            backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"backup_id":     backupUnitString(backupUnitBackupID),
		"backup_job_id": backupUnitString(backupUnitBackupID),
		"timeouts":      backupUnitTimeouts(t, objectType, "", "1s"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("Read swallowed a failed backup task query and reported a successful refresh")
	}
	detail := backupUnitDiagnosticsText(resp.Diagnostics)
	if !strings.Contains(detail, "DescribeBackupTasks") || !strings.Contains(detail, "Forbidden.RAM") {
		t.Fatalf("Read diagnostic does not identify the failed task query: %s", detail)
	}
}

// TestUnitAliCloudPolarDBBackupReadToleratesReleasedClusterTaskQuery covers the released-cluster
// read: after the cluster is released, DescribeBackupTasks answers InvalidDBClusterId.NotFound.
// That is a definitive "no task record", not a failed query, so Read must keep the tracked
// resource and its task id instead of failing the refresh.
func TestUnitAliCloudPolarDBBackupReadToleratesReleasedClusterTaskQuery(t *testing.T) {
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DescribeBackupPolicy":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeBackups":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		case "DescribeBackupTasks":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_job_id":     backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
		"backup_retention_policy_on_cluster_deletion": backupUnitString("ALL"),
		"timeouts": backupUnitTimeouts(t, objectType, "", "1s"),
	})

	var resp fwresource.ReadResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Read failed although a released cluster cannot answer DescribeBackupTasks: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if !server.seen("DescribeBackupTasks") {
		t.Fatal("Read did not query the tracked backup task")
	}
	if id, ok := backupUnitStateAttribute(t, resp.State, "id"); !ok || id != backupUnitClusterID+":"+backupUnitBackupID {
		t.Fatalf("Read dropped the released-cluster backup set: id=%q present=%v", id, ok)
	}
	if jobId, ok := backupUnitStateAttribute(t, resp.State, "backup_job_id"); !ok || jobId != backupUnitBackupID {
		t.Fatalf("Read lost the tracked backup task id: %q present=%v", jobId, ok)
	}
}

// TestUnitAliCloudPolarDBBackupDeletePropagatesNonNotFound covers QAP-01: a DeleteBackup failure
// that is not a NotFound must propagate; it must not be reinterpreted as "retained" success.
func TestUnitAliCloudPolarDBBackupDeletePropagatesNonNotFound(t *testing.T) {
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		if action == "DeleteBackup" {
			return `{"Code":"InvalidParameter","Message":"the backup set can not be deleted"}`
		}
		if action == "DescribeDetachedBackups" {
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":            backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id": backupUnitString(backupUnitClusterID),
		"backup_id":     backupUnitString(backupUnitBackupID),
	})

	var resp fwresource.DeleteResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Delete swallowed a non-NotFound DeleteBackup failure: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if server.seen("DescribeDetachedBackups") {
		t.Fatal("Delete consulted the retention list before classifying the original error")
	}
}

// TestUnitAliCloudPolarDBBackupDeleteRetainedWarning covers the G1 delete branch: DeleteBackup
// reports the released cluster as absent, the detached list still reports the exact set, so the
// destroy stops tracking it with a warning instead of failing or claiming a deletion.
func TestUnitAliCloudPolarDBBackupDeleteRetainedWarning(t *testing.T) {
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "DeleteBackup":
			return backupUnitNotFound("InvalidDBClusterId.NotFound")
		case "DescribeDetachedBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":     backupUnitString(backupUnitClusterID),
		"backup_id":         backupUnitString(backupUnitBackupID),
		"backup_start_time": backupUnitString(backupUnitStartTime),
	})

	var resp fwresource.DeleteResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: tfsdk.State{Schema: sch, Raw: state}}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Delete failed although the backup set is retained by a released cluster: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if len(resp.Diagnostics.Warnings()) == 0 {
		t.Fatal("Delete did not warn that the retained backup set stays in the cloud")
	}
}

// TestUnitAliCloudPolarDBBackupUpdateRejectsSameCountReplacement covers QAP-02: replacing one
// advanced policy with a different one keeps the list length but removes an entry, so Update must
// reject it instead of leaving the old cloud policy unmanaged.
func TestUnitAliCloudPolarDBBackupUpdateRejectsSameCountReplacement(t *testing.T) {
	client, server := newBackupUnitClient(t, func(action string, params map[string]string) string {
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	policy := func(srcType, destType, retention string) map[string]tftypes.Value {
		return map[string]tftypes.Value{
			"action_type":     backupUnitString("UPDATE"),
			"bak_type":        backupUnitString("F"),
			"src_type":        backupUnitString(srcType),
			"dest_type":       backupUnitString(destType),
			"src_region":      backupUnitString(backupUnitRegion),
			"dest_region":     backupUnitString(backupUnitRegion),
			"retention_type":  backupUnitString("delay"),
			"retention_value": backupUnitString(retention),
		}
	}
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{policy("db", "level1", "7")}),
	})
	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{policy("db", "level2", "7")}),
	})

	var resp fwresource.UpdateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Update(context.Background(),
		fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}, State: tfsdk.State{Schema: sch, Raw: state}},
		&resp)
	if !resp.Diagnostics.HasError() {
		t.Fatalf("Update accepted a same-sized advanced-data-policy replacement: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	if server.seen("ModifyBackupPolicy") {
		t.Fatal("Update sent ModifyBackupPolicy although an advanced entry was implicitly removed")
	}
}

// TestUnitAliCloudPolarDBBackupUpdateAdvancedFieldsSerialized covers G4: every advanced-policy
// field the configuration can carry must reach ModifyBackupPolicy with its own value, and every
// readable field must come back into state. Fields that a cloud run cannot safely construct
// (auto_created, policy_id are server-assigned) are covered here with fixed values.
func TestUnitAliCloudPolarDBBackupUpdateAdvancedFieldsSerialized(t *testing.T) {
	region := backupUnitRegion
	entry := map[string]tftypes.Value{
		"action_type":                 backupUnitString("UPDATE"),
		"only_preserve_one_each_hour": backupUnitBool(true),
		"storage_class":               backupUnitString("standard"),
		"retention_value":             backupUnitString("14"),
		"bak_type":                    backupUnitString("F"),
		"filter_key":                  backupUnitString("dayOfWeek"),
		"filter_value":                backupUnitString("1,2,3,4,5,6,7"),
		"dest_type":                   backupUnitString("level1"),
		"filter_type":                 backupUnitString("crontab"),
		"src_region":                  backupUnitString(region),
		"auto_created":                backupUnitBool(false),
		"only_preserve_one_each_day":  backupUnitBool(true),
		"src_type":                    backupUnitString("db"),
		"dump_action":                 backupUnitString("copy"),
		"dest_region":                 backupUnitString(region),
		"policy_id":                   backupUnitString("pid-unit"),
		"retention_type":              backupUnitString("delay"),
	}
	var modifyMu sync.Mutex
	var modifyRequest map[string]string
	client, _ := newBackupUnitClient(t, func(action string, params map[string]string) string {
		switch action {
		case "ModifyBackupPolicy":
			modifyMu.Lock()
			modifyRequest = params
			modifyMu.Unlock()
			return "{}"
		case "DescribeBackupPolicy":
			return fmt.Sprintf(`{"DBClusterId":%q,"BackupPolicyLevel":"Advanced","BackupRetentionPolicyOnClusterDeletion":"ALL",`+
				`"AdvancedDataPolicies":{"AdvancedDataPolicy":[{"AutoCreated":false,"BakType":"F","DestRegion":%q,"DestType":"level1",`+
				`"DumpAction":"copy","FilterKey":"dayOfWeek","FilterType":"crontab","FilterValue":"1,2,3,4,5,6,7",`+
				`"OnlyPreserveOneEachDay":true,"OnlyPreserveOneEachHour":true,"PolicyId":"pid-unit","RetentionType":"delay",`+
				`"RetentionValue":"14","SrcRegion":%q,"SrcType":"db","StorageClass":"standard"}]}}`,
				backupUnitClusterID, region, region)
		case "DescribeBackups":
			return backupUnitBackupsJSON(backupUnitBackupEntry(backupUnitBackupID, backupUnitClusterID, backupUnitStartTime))
		}
		return "{}"
	})
	r, sch := newBackupUnitResource(t, client)
	objectType := backupUnitObjectType(t, sch)

	prior := map[string]tftypes.Value{}
	for key, value := range entry {
		prior[key] = value
	}
	prior["retention_value"] = backupUnitString("7")
	state := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_start_time":      backupUnitString(backupUnitStartTime),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{prior}),
	})
	plan := backupUnitObject(t, objectType, map[string]tftypes.Value{
		"id":                     backupUnitString(backupUnitClusterID + ":" + backupUnitBackupID),
		"db_cluster_id":          backupUnitString(backupUnitClusterID),
		"backup_id":              backupUnitString(backupUnitBackupID),
		"backup_start_time":      backupUnitString(backupUnitStartTime),
		"backup_policy_level":    backupUnitString("Advanced"),
		"advanced_data_policies": backupUnitPolicyList(t, objectType, []map[string]tftypes.Value{entry}),
	})

	var resp fwresource.UpdateResponse
	resp.State = tfsdk.State{Schema: sch, Raw: state}
	r.Update(context.Background(),
		fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: sch, Raw: plan}, State: tfsdk.State{Schema: sch, Raw: state}},
		&resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Update failed on the full advanced policy: %v", backupUnitDiagnosticsText(resp.Diagnostics))
	}
	modifyMu.Lock()
	captured := modifyRequest
	modifyMu.Unlock()
	if captured == nil {
		t.Fatal("Update did not call ModifyBackupPolicy for the changed advanced policy")
	}
	// RPC requests carry nested entries either as AdvancedDataPolicies.N.<Field> form keys or
	// inside one JSON body; match on the leaf name so the assertion does not depend on which.
	leaf := func(name, want string) {
		for key, value := range captured {
			if strings.EqualFold(key, name) || strings.HasSuffix(strings.ToLower(key), "."+strings.ToLower(name)) {
				if strings.EqualFold(value, want) {
					return
				}
			}
		}
		t.Errorf("ModifyBackupPolicy request does not carry %s=%q (params: %v)", name, want, captured)
	}
	leaf("DBClusterId", backupUnitClusterID)
	leaf("ActionType", "UPDATE")
	leaf("BakType", "F")
	leaf("DestType", "level1")
	leaf("SrcType", "db")
	leaf("SrcRegion", region)
	leaf("DestRegion", region)
	leaf("FilterType", "crontab")
	leaf("FilterKey", "dayOfWeek")
	leaf("FilterValue", "1,2,3,4,5,6,7")
	leaf("OnlyPreserveOneEachDay", "true")
	leaf("OnlyPreserveOneEachHour", "true")
	leaf("RetentionType", "delay")
	leaf("RetentionValue", "14")
	leaf("StorageClass", "standard")
	leaf("DumpAction", "copy")
	leaf("PolicyId", "pid-unit")
	leaf("AutoCreated", "false")
	payloadPresent := false
	for key := range captured {
		if strings.Contains(key, "AdvancedDataPolicies.") {
			payloadPresent = true
			break
		}
	}
	if !payloadPresent {
		t.Errorf("ModifyBackupPolicy request carries no AdvancedDataPolicies payload (params: %v)", captured)
	}
	// G4 read-back: the API-readable advanced policy fields are mapped back into state.
	if value, ok := backupUnitPolicyAttribute(t, resp.State, "retention_value"); !ok || value != "14" {
		t.Errorf("advanced_data_policies.retention_value in state = %q present=%v, want 14", value, ok)
	}
	if value, ok := backupUnitPolicyAttribute(t, resp.State, "policy_id"); !ok || value != "pid-unit" {
		t.Errorf("advanced_data_policies.policy_id in state = %q present=%v, want pid-unit", value, ok)
	}
}

// backupUnitPolicyAttribute returns the first known value of a string attribute inside the
// configured advanced_data_policies block of the read-back state.
func backupUnitPolicyAttribute(t *testing.T, state tfsdk.State, name string) (string, bool) {
	t.Helper()
	var list types.List
	if diagnostics := state.GetAttribute(context.Background(), fwpath.Root("advanced_data_policies"), &list); diagnostics.HasError() {
		return "", false
	}
	if list.IsNull() || list.IsUnknown() {
		return "", false
	}
	for _, element := range list.Elements() {
		policy, ok := element.(types.Object)
		if !ok {
			continue
		}
		value, ok := policy.Attributes()[name].(types.String)
		if ok && !value.IsNull() && !value.IsUnknown() {
			return value.ValueString(), true
		}
	}
	return "", false
}

// TestAccAliCloudPolarDBBackup_advancedDataPoliciesMemberOrder exercises the reorder contract of
// the configurable advanced_data_policies list. Two policies that are legal together (db->level1
// and level1->level2) are applied in one order and then in the opposite order; the SDK requires
// the post-apply plan of the second apply to be empty, so a canonicalizing read or an
// order-sensitive change filter fails here.
func TestAccAliCloudPolarDBBackup_advancedDataPoliciesMemberOrder(t *testing.T) {
	resourceId := "alicloud_polardb_backup.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolardbBackupOrder%s%d", polardbBackupTestRegion(), rand)
	region := polardbBackupTestRegion()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccAliCloudPolarDBBackupDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBBackupDestroy,
		Steps: []resource.TestStep{
			{
				// Both directions of the chain are legal together and are accepted in this order.
				Config: testAccConfig(map[string]interface{}{
					"db_cluster_id":       "${alicloud_polardb_cluster.default.id}",
					"backup_policy_level": "Advanced",
					"advanced_data_policies": []interface{}{
						map[string]interface{}{
							"action_type":                "CREATE",
							"bak_type":                   "F",
							"src_type":                   "db",
							"dest_type":                  "level1",
							"src_region":                 region,
							"dest_region":                region,
							"filter_type":                "crontab",
							"filter_key":                 "dayOfWeek",
							"filter_value":               "1,2,3,4,5,6,7",
							"only_preserve_one_each_day": true,
							"retention_type":             "delay",
							"retention_value":            "7",
						},
						map[string]interface{}{
							"action_type":                "CREATE",
							"bak_type":                   "F",
							"src_type":                   "level1",
							"dest_type":                  "level2",
							"src_region":                 region,
							"dest_region":                region,
							"filter_type":                "crontab",
							"filter_key":                 "dayOfWeek",
							"filter_value":               "1,2,3,4,5,6,7",
							"only_preserve_one_each_day": true,
							"retention_type":             "delay",
							"retention_value":            "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "advanced_data_policies.#", "2"),
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, "7"),
				),
			},
			{
				// Only the member order changes; the post-apply plan must still be empty.
				Config: testAccConfig(map[string]interface{}{
					"advanced_data_policies": []interface{}{
						map[string]interface{}{
							"action_type":                "CREATE",
							"bak_type":                   "F",
							"src_type":                   "level1",
							"dest_type":                  "level2",
							"src_region":                 region,
							"dest_region":                region,
							"filter_type":                "crontab",
							"filter_key":                 "dayOfWeek",
							"filter_value":               "1,2,3,4,5,6,7",
							"only_preserve_one_each_day": true,
							"retention_type":             "delay",
							"retention_value":            "7",
						},
						map[string]interface{}{
							"action_type":                "CREATE",
							"bak_type":                   "F",
							"src_type":                   "db",
							"dest_type":                  "level1",
							"src_region":                 region,
							"dest_region":                region,
							"filter_type":                "crontab",
							"filter_key":                 "dayOfWeek",
							"filter_value":               "1,2,3,4,5,6,7",
							"only_preserve_one_each_day": true,
							"retention_type":             "delay",
							"retention_value":            "7",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "advanced_data_policies.#", "2"),
					testAccCheckPolarDBBackupAdvancedPolicy(resourceId, region, "7"),
				),
			},
		},
	})
}

// resourceTestAccConfigFunc returns a config builder that accumulates the primary resource's
// attributes across steps, so a step can change only the attribute under test while the rest of
// the configuration stays applied.
func resourceTestAccConfigFunc(resourceId, name string, dependence func(string) string) func(map[string]interface{}) string {
	attributes := map[string]interface{}{}
	return func(change map[string]interface{}) string {
		for key, value := range change {
			attributes[key] = value
		}
		parts := strings.Split(resourceId, ".")
		return dependence(name) + fmt.Sprintf("\n\nresource %q %q {\n%s\n}\n", parts[0], parts[1], polardbBackupHCLAttributes(2, attributes))
	}
}

func polardbBackupHCLAttributes(indent int, attributes map[string]interface{}) string {
	keys := make([]string, 0, len(attributes))
	for key := range attributes {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pad := strings.Repeat(" ", indent)
	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		switch value := attributes[key].(type) {
		case []interface{}:
			if len(value) > 0 {
				if _, ok := value[0].(map[string]interface{}); ok {
					for _, item := range value {
						lines = append(lines, pad+key+" "+polardbBackupHCLObject(indent, item.(map[string]interface{})))
					}
					continue
				}
			}
			lines = append(lines, pad+key+" = []")
		case map[string]interface{}:
			lines = append(lines, pad+key+" "+polardbBackupHCLObject(indent, value))
		case string:
			lines = append(lines, pad+key+" = "+strconv.Quote(value))
		case bool:
			lines = append(lines, pad+key+" = "+strconv.FormatBool(value))
		case int:
			lines = append(lines, pad+key+" = "+strconv.Itoa(value))
		}
	}
	return strings.Join(lines, "\n")
}

func polardbBackupHCLObject(indent int, attributes map[string]interface{}) string {
	return "{\n" + polardbBackupHCLAttributes(indent+2, attributes) + "\n" + strings.Repeat(" ", indent) + "}"
}
