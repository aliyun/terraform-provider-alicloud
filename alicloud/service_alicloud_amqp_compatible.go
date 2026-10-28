package alicloud

import "github.com/aliyun/terraform-provider-alicloud/alicloud/convert"

// This file retains root-package names for the amqp instance renewal duration
// unit converters so unmigrated root consumers (e.g. the rocketmq instance
// resource, which shares the "ons" instance model via BssOpenApi) can keep
// calling them unqualified. The single implementation lives in the shared
// alicloud/convert leaf package; the names below are aliases, not duplicate
// bodies, so the migrated alicloud/service/amqpopen package and the root
// consumers cannot diverge.

// Migrated to alicloud/convert.ConvertAmqpInstanceRenewalDurationUnitResponse.
var convertAmqpInstanceRenewalDurationUnitResponse = convert.ConvertAmqpInstanceRenewalDurationUnitResponse

// Migrated to alicloud/convert.ConvertAmqpInstanceRenewalDurationUnitRequest.
var convertAmqpInstanceRenewalDurationUnitRequest = convert.ConvertAmqpInstanceRenewalDurationUnitRequest
