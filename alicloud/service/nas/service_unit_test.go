package nas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// unitTestProviderResources 提供本产品包已注册的 NAS Resource Schema 映射，供 alicloud/service/nas 内部
// TestUnit 用例（*_unit_test.go）在不依赖根 provider 与 acctest 的情况下访问资源 Schema 与 CRUD 方法。
// 镜像 alicloud/service/ecs/disk_unit_test.go 中 unitTestProviderResources 的本地副本惯例。
var unitTestProviderResources = map[string]*schema.Resource{
	"alicloud_nas_access_group":         resourceAliCloudNasAccessGroup(),
	"alicloud_nas_access_rule":          resourceAliCloudNasAccessRule(),
	"alicloud_nas_auto_snapshot_policy": resourceAliCloudNasAutoSnapshotPolicy(),
	"alicloud_nas_fileset":              resourceAliCloudNasFileset(),
	"alicloud_nas_file_system":          resourceAliCloudNasFileSystem(),
	"alicloud_nas_mount_target":         resourceAliCloudNasMountTarget(),
	"alicloud_nas_data_flow":            resourceAlicloudNasDataFlow(),
	"alicloud_nas_lifecycle_policy":     resourceAlicloudNasLifecyclePolicy(),
	"alicloud_nas_recycle_bin":          resourceAlicloudNasRecycleBin(),
	"alicloud_nas_snapshot":             resourceAlicloudNasSnapshot(),
	"alicloud_nas_smb_acl_attachment":   resourceAlicloudNasSmbAclAttachment(),
}
