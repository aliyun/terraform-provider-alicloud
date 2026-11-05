package sslcertificateservice

import (
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const sslCertificatesServiceWarehouseTimeout = 5 * time.Minute

// CreateSslCertificatesServiceWarehouse calls cas/2020-04-07 CreateWarehouse and returns the
// allocated WarehouseInstanceId (the resource primary key). Name, Type and Biz are the only
// CreateWarehouse inputs; there is no Update API, so every argument is replacement-only.
func (s *SslCertificatesServiceServiceV2) CreateSslCertificatesServiceWarehouse(name, warehouseType, biz string) (id string, err error) {
	client := s.client
	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("CreateWarehouse"),
		"Name":        name,
		"Type":        warehouseType,
		"Biz":         biz,
	}
	action := "CreateWarehouse"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(sslCertificatesServiceWarehouseTimeout, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, name, action, errs.AlibabaCloudSdkGoERROR)
	}

	idRaw, err := jsonpath.Get("$.WarehouseInstanceId", response)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, name, "$.WarehouseInstanceId", response)
	}
	id, _ = idRaw.(string)
	return id, nil
}

// DescribeSslCertificatesServiceWarehouse reads a single warehouse by id through cas/2020-04-07
// ListWarehouse with the WarehouseInstanceIds filter. There is no dedicated GetWarehouse API, so
// ListWarehouse is the only read path; an empty match is reported as NotFound.
func (s *SslCertificatesServiceServiceV2) DescribeSslCertificatesServiceWarehouse(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"WarehouseInstanceIds": []interface{}{id},
	}
	action := "ListWarehouse"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(sslCertificatesServiceWarehouseTimeout, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	dataRaw, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	data, ok := dataRaw.([]interface{})
	if !ok {
		return object, errs.WrapErrorf(errs.NotFoundErr("Warehouse", id), errs.NotFoundMsg, response)
	}
	for _, item := range data {
		m, _ := item.(map[string]interface{})
		if wid, _ := m["WarehouseInstanceId"].(string); wid == id {
			return m, nil
		}
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("Warehouse", id), errs.NotFoundMsg, response)
}

// DeleteSslCertificatesServiceWarehouse calls cas/2020-04-07 DeleteWarehouse. DeleteWarehouse has
// no return body beyond RequestId; a NotFound code on a missing warehouse is treated as success so
// the destroy stays idempotent.
func (s *SslCertificatesServiceServiceV2) DeleteSslCertificatesServiceWarehouse(id string) (err error) {
	client := s.client
	request := map[string]interface{}{
		"WarehouseInstanceId": id,
	}
	action := "DeleteWarehouse"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(sslCertificatesServiceWarehouseTimeout, func() *retry.RetryError {
		response, err = client.RpcPost("cas", "2020-04-07", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NotFound", "ResourceNotFound", "NotExistsWarehouse"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
