package miscellaneous_test

import (
	"testing"

	"github.com/kiln-mid/pkg/miscellaneous"
	"github.com/stretchr/testify/require"
)

func Test_Contain(t *testing.T) {
	t.Run("splitToString", func(t *testing.T) {
		var inputs = []int{1, 2, 3}

		res := miscellaneous.SplitToString(inputs, ",")

		require.Equal(t, res, "1,2,3")
	})
}
