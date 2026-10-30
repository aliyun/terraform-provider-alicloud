package lindorm

import "time"

// Shims that let the external test package (lindorm_test) exercise the
// request-level update/delete paths and shorten their polling intervals without
// widening the production API. This file is only compiled into the test binary
// of this package.
var (
	ResourceLindormInstanceUpdateForTest = resourceLindormInstanceUpdate
	ResourceLindormInstanceDeleteForTest = resourceLindormInstanceDelete
)

// ShortenLindormStateWaitForTest shortens both state wait intervals used by the
// update and delete paths and returns a restore function.
func ShortenLindormStateWaitForTest(d time.Duration) func() {
	previousUpgrade := lindormUpgradeStateWaitDuration
	previousPayType := lindormPayTypeWaitDuration
	lindormUpgradeStateWaitDuration = d
	lindormPayTypeWaitDuration = d
	return func() {
		lindormUpgradeStateWaitDuration = previousUpgrade
		lindormPayTypeWaitDuration = previousPayType
	}
}

// LindormPostCreateReadToleranceForCreateTimeoutForTest exposes the tolerance
// calculation so the external test can assert its bound.
func LindormPostCreateReadToleranceForCreateTimeoutForTest(createTimeout time.Duration) time.Duration {
	return lindormPostCreateReadToleranceFor(createTimeout)
}

// LindormPostCreateReadToleranceLastForTest exposes the recorded tolerated
// error so the external test can assert diagnostics recording.
func LindormPostCreateReadToleranceLastForTest(t *LindormPostCreateReadTolerance) error {
	return t.last()
}

// IsLindormTransientPostCreateReadErrorForTest exposes the transient error
// classifier so the external test can cover every error shape the shared client
// may return.
func IsLindormTransientPostCreateReadErrorForTest(err error) bool {
	return isLindormTransientPostCreateReadError(err)
}
