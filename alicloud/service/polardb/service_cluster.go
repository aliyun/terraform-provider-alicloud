package polardb

import (
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"time"

	polardbsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribePolardbCluster fetches the PolarDB cluster to read its DBType. It is a
// hand-written supplement to the generated PolardbServiceV2: the generator omits
// cluster DBType lookup, which the baseline root PolarDBService.DescribePolarDBCluster
// (alicloud/service_alicloud_polardb.go) provided so that resourceDatabaseCreate can
// default Collate/Ctype to "C" for PostgreSQL/Oracle clusters. Ported verbatim from
// the baseline using the polardb SDK so the proven DescribeDBClusters call is preserved.
func (s *PolardbServiceV2) DescribePolardbCluster(id string) (instance *polardbsdk.DBCluster, err error) {
	request := polardbsdk.CreateDescribeDBClustersRequest()
	request.RegionId = s.client.RegionId
	request.DBClusterIds = id
	var response *polardbsdk.DescribeDBClustersResponse
	wait := tfhelper.IncrementalWait(3*time.Second, 3*time.Second)
	err = tfretry.Retry(5*time.Minute, func() *tfretry.RetryError {
		raw, err := s.client.WithPolarDBClient(func(polarDBClient *polardbsdk.Client) (interface{}, error) {
			return polarDBClient.DescribeDBClusters(request)
		})
		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		tflogs.AddDebug(request.GetActionName(), raw, request)
		response, _ = raw.(*polardbsdk.DescribeDBClustersResponse)
		return nil
	})
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return nil, tferrs.WrapErrorf(err, tferrs.NotFoundMsg, tferrs.AlibabaCloudSdkGoERROR)
		}
		return nil, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, request.GetActionName(), tferrs.AlibabaCloudSdkGoERROR)
	}
	if len(response.Items.DBCluster) < 1 {
		return nil, tferrs.WrapErrorf(tferrs.NotFoundErr("Cluster", id), tferrs.NotFoundMsg, tferrs.AlibabaCloudSdkGoERROR)
	}

	return &response.Items.DBCluster[0], nil
}
