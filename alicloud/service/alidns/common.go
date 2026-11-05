package alidns

// DefaultTimeout and DefaultTimeoutMedium are local 1:1 copies of the general
// waiter timeout constants defined in alicloud/common.go. The root copies are
// retained for their many existing root consumers; alidns cannot import the
// root package. Values are unchanged and remain consumed as seconds by the
// package waiters (e.g. WaitForAlidnsDomainAttachment timeout int).
const DefaultTimeout = 120
const DefaultTimeoutMedium = 500
