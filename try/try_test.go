package try_test

import (
	"errors"
	"testing"

	"github.com/martinohmann/exp/try"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("must", func(t *testing.T) {
		var a, b, c int

		err := try.Run(func() {
			a = 42
			try.Check(nil)
			b = 43
			try.Check(errors.New("whoops"))
			c = 44
		})

		require.EqualError(t, err, "whoops")
		require.Equal(t, 42, a)
		require.Equal(t, 43, b)
		require.Equal(t, 0, c)
	})

	t.Run("mustf", func(t *testing.T) {
		var a, b, c int

		err := try.Run(func() {
			a = 42
			try.Checkf(nil, "error occurred")
			b = 43
			try.Checkf(errors.New("whoops"), "error occurred: a=%d, b=%d, c=%d", a, b, c)
			c = 44
		})

		require.EqualError(t, err, "error occurred: a=42, b=43, c=0: whoops")
		require.Equal(t, 42, a)
		require.Equal(t, 43, b)
		require.Equal(t, 0, c)
	})

	t.Run("no error", func(t *testing.T) {
		var a int
		err := try.Run(func() {
			a = 42
		})

		require.NoError(t, err)
		require.Equal(t, 42, a)
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("expected panic but got nil")
			}
		}()
		_ = try.Run(func() {
			panic("whoops")
		})
	})

	t.Run("nested", func(t *testing.T) {
		fn := func() {
			try.Check(errors.New("foo"))
		}

		err := try.Run(func() {
			fn()
			try.Check(errors.New("bar"))
		})

		require.EqualError(t, err, "foo")
	})

	t.Run("catch", func(t *testing.T) {
		f := func() (err error) {
			defer try.Catch(&err)
			try.Check(errors.New("foo"))
			return
		}

		require.EqualError(t, f(), "foo")
	})

	t.Run("catch nested", func(t *testing.T) {
		g := func() {
			try.Check(errors.New("foo"))
		}

		f := func() (err error) {
			defer try.Catch(&err)
			g()
			try.Check(errors.New("bar"))
			return
		}

		require.EqualError(t, f(), "foo")
	})
}
