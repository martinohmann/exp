package cobrax

import (
	"errors"
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestHookViper(t *testing.T) {
	t.Run("ignores nil funcs in chain", func(t *testing.T) {
		require := require.New(t)

		defer func() {
			if err := recover(); err != nil {
				t.Fatalf("unexpected panic: %v", err)
			}
		}()

		cmd := &cobra.Command{}
		require.NoError(cmd.Flags().Parse([]string{}))

		hook := HookViper(nil, nil)

		require.NoError(hook(cmd, []string{}))
	})

	t.Run("calls funcs in chain in order", func(t *testing.T) {
		require := require.New(t)

		var results []int

		f1 := func(*cobra.Command, []string) error {
			results = append(results, 1)
			return nil
		}

		f2 := func(*cobra.Command, []string) error {
			results = append(results, 2)
			return nil
		}

		cmd := &cobra.Command{}
		require.NoError(cmd.Flags().Parse([]string{}))

		hook := HookViper(nil, f1, f2)

		require.NoError(hook(cmd, []string{}))
		require.Equal([]int{1, 2}, results)
	})

	t.Run("returns error from func chain", func(t *testing.T) {
		require := require.New(t)

		var results []int

		f1 := func(*cobra.Command, []string) error {
			results = append(results, 1)
			return errors.New("whoops")
		}

		f2 := func(*cobra.Command, []string) error {
			results = append(results, 2)
			return nil
		}

		cmd := &cobra.Command{}
		require.NoError(cmd.Flags().Parse([]string{}))

		hook := HookViper(nil, f1, f2)

		require.EqualError(hook(cmd, []string{}), "whoops")
		require.Equal([]int{1}, results)
	})

	t.Run("intertwines before calling func chain", func(t *testing.T) {
		require := require.New(t)

		fn := func(cmd *cobra.Command, args []string) error {
			val, err := cmd.Flags().GetString("the-flag")
			if err != nil {
				return err
			}

			if val != "the-value" {
				return fmt.Errorf("want: %q, got: %q", "the-value", val)
			}

			return nil
		}

		cmd := &cobra.Command{}
		cmd.Flags().String("the-flag", "", "flag usage")

		require.NoError(cmd.Flags().Parse([]string{"--the-flag", "the-value"}))

		hook := HookViper(nil, fn)

		require.NoError(hook(cmd, []string{}))
	})

	t.Run("returns error from intertwine hook", func(t *testing.T) {
		require := require.New(t)

		fn := func(*cobra.Command, []string) error {
			t.Fatal("unexpected call to chained func")
			return nil
		}

		cmd := &cobra.Command{}
		cmd.Flags().IP("ip", nil, "an ip")

		require.NoError(cmd.Flags().Parse([]string{}))

		v := viper.New()

		// causes Intertwine error because ip has invalid format.
		v.Set("ip", "invalid-ip")

		hook := HookViper(v, fn)

		require.EqualError(
			hook(cmd, []string{}),
			`failed to set flag from env or config: invalid argument "invalid-ip" for "--ip" flag: failed to parse IP: "invalid-ip"`,
		)
	})
}

func TestExecute(t *testing.T) {
	t.Run("preserves cmd.PersistentPreRun", func(t *testing.T) {
		var called bool
		cmd := &cobra.Command{
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				called = true
			},
			Run: func(cmd *cobra.Command, args []string) {},
		}

		require.NoError(t, Execute(cmd, nil))
		require.True(t, called)
	})

	t.Run("preserves cmd.PersistentPreRunE", func(t *testing.T) {
		var called bool
		cmd := &cobra.Command{
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				called = true
				return nil
			},
			Run: func(cmd *cobra.Command, args []string) {},
		}

		require.NoError(t, Execute(cmd, nil))
		require.True(t, called)
	})

	t.Run("sets command flags from viper", func(t *testing.T) {
		var (
			val    string
			called bool
		)

		cmd := &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				called = true
				require.Equal(t, "the-value", val)
			},
		}

		cmd.Flags().StringVar(&val, "the-flag", "the-default", "a string flag")

		v := viper.New()
		v.Set("the-flag", "the-value")

		require.NoError(t, Execute(cmd, v))
		require.True(t, called)
	})

	t.Run("preserves cmd.PersistentPreRun from parent cmd", func(t *testing.T) {
		var called bool
		parent := &cobra.Command{
			Use: "parent",
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				called = true
			},
		}

		child := &cobra.Command{
			Use: "child",
			Run: func(cmd *cobra.Command, args []string) {},
		}

		parent.AddCommand(child)

		parent.SetArgs([]string{"child"})

		require.NoError(t, Execute(parent, nil))
		require.True(t, called)
	})
}
