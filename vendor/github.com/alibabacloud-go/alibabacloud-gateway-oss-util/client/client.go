// This file is auto-generated, don't edit it. Thanks.
/**
 *
 * @param bodyStr the response body
 * @param apiName the api name
 * @return the parsed result
 */
package client

import (
	"reflect"

	xml "github.com/alibabacloud-go/tea-xml/service"
)

var typeRegistry = make(map[string]reflect.Type)

func init() {
	// for oss
	typeRegistry["CompleteMultipartUpload"] = reflect.TypeOf(CompleteMultipartUploadResponseBody{})
	typeRegistry["CopyObject"] = reflect.TypeOf(CopyObjectResponseBody{})
	typeRegistry["CreateAccessPoint"] = reflect.TypeOf(CreateAccessPointResponseBody{})
	typeRegistry["CreateAccessPointForObjectProcess"] = reflect.TypeOf(CreateAccessPointForObjectProcessResponseBody{})
	typeRegistry["CreateBucketDataRedundancyTransition"] = reflect.TypeOf(CreateBucketDataRedundancyTransitionResponseBody{})
	typeRegistry["CreateCnameToken"] = reflect.TypeOf(CreateCnameTokenResponseBody{})
	typeRegistry["DeleteMultipleObjects"] = reflect.TypeOf(DeleteMultipleObjectsResponseBody{})
	typeRegistry["DescribeRegions"] = reflect.TypeOf(DescribeRegionsResponseBody{})
	typeRegistry["DoMetaQuery"] = reflect.TypeOf(DoMetaQueryResponseBody{})
	typeRegistry["GetAccessPoint"] = reflect.TypeOf(GetAccessPointResponseBody{})
	typeRegistry["GetAccessPointConfigForObjectProcess"] = reflect.TypeOf(GetAccessPointConfigForObjectProcessResponseBody{})
	typeRegistry["GetAccessPointForObjectProcess"] = reflect.TypeOf(GetAccessPointForObjectProcessResponseBody{})
	typeRegistry["GetAccessPointPublicAccessBlock"] = reflect.TypeOf(GetAccessPointPublicAccessBlockResponseBody{})
	typeRegistry["GetBucketAccessMonitor"] = reflect.TypeOf(GetBucketAccessMonitorResponseBody{})
	typeRegistry["GetBucketAcl"] = reflect.TypeOf(GetBucketAclResponseBody{})
	typeRegistry["GetBucketArchiveDirectRead"] = reflect.TypeOf(GetBucketArchiveDirectReadResponseBody{})
	typeRegistry["GetBucketCallbackPolicy"] = reflect.TypeOf(GetBucketCallbackPolicyResponseBody{})
	typeRegistry["GetBucketCors"] = reflect.TypeOf(GetBucketCorsResponseBody{})
	typeRegistry["GetBucketDataRedundancyTransition"] = reflect.TypeOf(GetBucketDataRedundancyTransitionResponseBody{})
	typeRegistry["GetBucketEncryption"] = reflect.TypeOf(GetBucketEncryptionResponseBody{})
	typeRegistry["GetBucketHttpsConfig"] = reflect.TypeOf(GetBucketHttpsConfigResponseBody{})
	typeRegistry["GetBucketInfo"] = reflect.TypeOf(GetBucketInfoResponseBody{})
	typeRegistry["GetBucketInventory"] = reflect.TypeOf(GetBucketInventoryResponseBody{})
	typeRegistry["GetBucketLifecycle"] = reflect.TypeOf(GetBucketLifecycleResponseBody{})
	typeRegistry["GetBucketLocation"] = reflect.TypeOf(GetBucketLocationResponseBody{})
	typeRegistry["GetBucketLogging"] = reflect.TypeOf(GetBucketLoggingResponseBody{})
	typeRegistry["GetBucketPolicyStatus"] = reflect.TypeOf(GetBucketPolicyStatusResponseBody{})
	typeRegistry["GetBucketPublicAccessBlock"] = reflect.TypeOf(GetBucketPublicAccessBlockResponseBody{})
	typeRegistry["GetBucketReferer"] = reflect.TypeOf(GetBucketRefererResponseBody{})
	typeRegistry["GetBucketReplication"] = reflect.TypeOf(GetBucketReplicationResponseBody{})
	typeRegistry["GetBucketReplicationLocation"] = reflect.TypeOf(GetBucketReplicationLocationResponseBody{})
	typeRegistry["GetBucketReplicationProgress"] = reflect.TypeOf(GetBucketReplicationProgressResponseBody{})
	typeRegistry["GetBucketRequestPayment"] = reflect.TypeOf(GetBucketRequestPaymentResponseBody{})
	typeRegistry["GetBucketResourceGroup"] = reflect.TypeOf(GetBucketResourceGroupResponseBody{})
	typeRegistry["GetBucketResponseHeader"] = reflect.TypeOf(GetBucketResponseHeaderResponseBody{})
	typeRegistry["GetBucketStat"] = reflect.TypeOf(GetBucketStatResponseBody{})
	typeRegistry["GetBucketTags"] = reflect.TypeOf(GetBucketTagsResponseBody{})
	typeRegistry["GetBucketTransferAcceleration"] = reflect.TypeOf(GetBucketTransferAccelerationResponseBody{})
	typeRegistry["GetBucketVersioning"] = reflect.TypeOf(GetBucketVersioningResponseBody{})
	typeRegistry["GetBucketWebsite"] = reflect.TypeOf(GetBucketWebsiteResponseBody{})
	typeRegistry["GetBucketWorm"] = reflect.TypeOf(GetBucketWormResponseBody{})
	typeRegistry["GetCnameToken"] = reflect.TypeOf(GetCnameTokenResponseBody{})
	typeRegistry["GetLiveChannelHistory"] = reflect.TypeOf(GetLiveChannelHistoryResponseBody{})
	typeRegistry["GetLiveChannelInfo"] = reflect.TypeOf(GetLiveChannelInfoResponseBody{})
	typeRegistry["GetLiveChannelStat"] = reflect.TypeOf(GetLiveChannelStatResponseBody{})
	typeRegistry["GetMetaQueryStatus"] = reflect.TypeOf(GetMetaQueryStatusResponseBody{})
	typeRegistry["GetObjectAcl"] = reflect.TypeOf(GetObjectAclResponseBody{})
	typeRegistry["GetObjectTagging"] = reflect.TypeOf(GetObjectTaggingResponseBody{})
	typeRegistry["GetPublicAccessBlock"] = reflect.TypeOf(GetPublicAccessBlockResponseBody{})
	typeRegistry["GetStyle"] = reflect.TypeOf(GetStyleResponseBody{})
	typeRegistry["GetUserAntiDDosInfo"] = reflect.TypeOf(GetUserAntiDDosInfoResponseBody{})
	typeRegistry["GetUserDefinedLogFieldsConfig"] = reflect.TypeOf(GetUserDefinedLogFieldsConfigResponseBody{})
	typeRegistry["InitiateMultipartUpload"] = reflect.TypeOf(InitiateMultipartUploadResponseBody{})
	typeRegistry["ListAccessPoints"] = reflect.TypeOf(ListAccessPointsResponseBody{})
	typeRegistry["ListAccessPointsForObjectProcess"] = reflect.TypeOf(ListAccessPointsForObjectProcessResponseBody{})
	typeRegistry["ListBucketAntiDDosInfo"] = reflect.TypeOf(ListBucketAntiDDosInfoResponseBody{})
	typeRegistry["ListBucketDataRedundancyTransition"] = reflect.TypeOf(ListBucketDataRedundancyTransitionResponseBody{})
	typeRegistry["ListBucketInventory"] = reflect.TypeOf(ListBucketInventoryResponseBody{})
	typeRegistry["ListBuckets"] = reflect.TypeOf(ListBucketsResponseBody{})
	typeRegistry["ListCname"] = reflect.TypeOf(ListCnameResponseBody{})
	typeRegistry["ListLiveChannel"] = reflect.TypeOf(ListLiveChannelResponseBody{})
	typeRegistry["ListMultipartUploads"] = reflect.TypeOf(ListMultipartUploadsResponseBody{})
	typeRegistry["ListObjectVersions"] = reflect.TypeOf(ListObjectVersionsResponseBody{})
	typeRegistry["ListObjects"] = reflect.TypeOf(ListObjectsResponseBody{})
	typeRegistry["ListObjectsV2"] = reflect.TypeOf(ListObjectsV2ResponseBody{})
	typeRegistry["ListParts"] = reflect.TypeOf(ListPartsResponseBody{})
	typeRegistry["ListStyle"] = reflect.TypeOf(ListStyleResponseBody{})
	typeRegistry["PutLiveChannel"] = reflect.TypeOf(PutLiveChannelResponseBody{})
	typeRegistry["UploadPartCopy"] = reflect.TypeOf(UploadPartCopyResponseBody{})

	// for hcs-mgw
	typeRegistry["GetAddress"] = reflect.TypeOf(GetAddressResponseBody{})
	typeRegistry["GetAgent"] = reflect.TypeOf(GetAgentResponseBody{})
	typeRegistry["GetAgentStatus"] = reflect.TypeOf(GetAgentStatusResponseBody{})
	typeRegistry["GetJob"] = reflect.TypeOf(GetJobResponseBody{})
	typeRegistry["GetJobResult"] = reflect.TypeOf(GetJobResultResponseBody{})
	typeRegistry["GetReport"] = reflect.TypeOf(GetReportResponseBody{})
	typeRegistry["GetTunnel"] = reflect.TypeOf(GetTunnelResponseBody{})
	typeRegistry["ListAddress"] = reflect.TypeOf(ListAddressResponseBody{})
	typeRegistry["ListAgent"] = reflect.TypeOf(ListAgentResponseBody{})
	typeRegistry["ListJob"] = reflect.TypeOf(ListJobResponseBody{})
	typeRegistry["ListJobHistory"] = reflect.TypeOf(ListJobHistoryResponseBody{})
	typeRegistry["ListTunnel"] = reflect.TypeOf(ListTunnelResponseBody{})
	typeRegistry["VerifyAddress"] = reflect.TypeOf(VerifyAddressResponseBody{})
}

func ParseXml(bodyStr *string, apiName *string) (_result interface{}, _err error) {
	var bodyStruct interface{} = nil
	if typ, ok := typeRegistry[*apiName]; ok {
		bodyStruct = reflect.New(typ).Interface()
	}
	return xml.ParseXml(bodyStr, bodyStruct), nil
}
