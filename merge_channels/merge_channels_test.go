package merge_channels

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMergeChannelPattern(t *testing.T) {
	mergedArray := MergeChannelPattern()
	require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, mergedArray)
}
