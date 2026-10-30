package alicloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeArmsNotifyTimeRangeDefaultsEmptyValues(t *testing.T) {
	start, end := normalizeArmsNotifyTimeRange("", "")
	require.Equal(t, "00:00", start)
	require.Equal(t, "23:59", end)

	start, end = normalizeArmsNotifyTimeRange(nil, nil)
	require.Equal(t, "00:00", start)
	require.Equal(t, "23:59", end)
}

func TestNormalizeArmsNotifyTimeRangeKeepsExplicitValues(t *testing.T) {
	start, end := normalizeArmsNotifyTimeRange("08:30", "20:15")
	require.Equal(t, "08:30", start)
	require.Equal(t, "20:15", end)
}
