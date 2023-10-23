package gitprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitter/internal/providers/gitprovider/testpg"
)

func TestDrift(t *testing.T) {
	t.Run("master0-feature0", func(t *testing.T) {
		play := must(testpg.New())

		err := play.Master()
		require.NoError(t, err)
		err = play.AddTestFile("common-1")
		require.NoError(t, err)

		err = play.Checkout("feature", true)
		require.NoError(t, err)

		drifter := NewDrifter(play.R)
		a, b, err := drifter.calc("master", "feature")
		require.NoError(t, err)
		assert.Equal(t, 0, a)
		assert.Equal(t, 0, b)
	})
	t.Run("master0-feature1", func(t *testing.T) {
		play := must(testpg.New())

		err := play.Master()
		require.NoError(t, err)
		err = play.AddTestFile("common-1")
		require.NoError(t, err)

		err = play.Checkout("feature", true)
		require.NoError(t, err)
		err = play.AddTestFile("feature-1")
		require.NoError(t, err)

		drifter := NewDrifter(play.R)
		a, b, err := drifter.calc("master", "feature")
		require.NoError(t, err)
		assert.Equal(t, 0, a)
		assert.Equal(t, 1, b)
	})
	t.Run("master1-feature0", func(t *testing.T) {
		play := must(testpg.New())

		err := play.Master()
		require.NoError(t, err)
		err = play.AddTestFile("common-1")
		require.NoError(t, err)

		err = play.Checkout("feature", true)
		require.NoError(t, err)

		err = play.Checkout("master", false)
		require.NoError(t, err)
		err = play.AddTestFile("master-1")
		require.NoError(t, err)

		drifter := NewDrifter(play.R)
		a, b, err := drifter.calc("master", "feature")
		require.NoError(t, err)
		assert.Equal(t, 1, a)
		assert.Equal(t, 0, b)
	})
	t.Run("master1-feature1", func(t *testing.T) {
		play := must(testpg.New())

		err := play.Master()
		require.NoError(t, err)
		err = play.AddTestFile("common-1")
		require.NoError(t, err)

		err = play.Checkout("feature", true)
		require.NoError(t, err)
		err = play.AddTestFile("feature-1")
		require.NoError(t, err)

		err = play.Checkout("master", false)
		require.NoError(t, err)
		err = play.AddTestFile("master-1")
		require.NoError(t, err)

		drifter := NewDrifter(play.R)
		a, b, err := drifter.calc("master", "feature")
		require.NoError(t, err)
		assert.Equal(t, 1, a)
		assert.Equal(t, 1, b)
	})
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
