package cobrax

import (
	"github.com/martinohmann/exp/pflagx"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// HookViper returns a hook func which executes pflagx.BindViper on
// *cobra.Command's flags with the passed in *viper.Viper instance. Should be
// set as the PersistentPreRunE hook on the root command. Any non-nil func in
// chain will be called after the Intertwine hook succeeds.
//
// See pflagx.BindViper for more information.
func HookViper(v *viper.Viper, chain ...func(*cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := pflagx.BindViper(cmd.Flags(), v); err != nil {
			return err
		}

		for _, fn := range chain {
			if fn == nil {
				continue
			}

			if err := fn(cmd, args); err != nil {
				return err
			}
		}

		return nil
	}
}

// Execute executes cmd after binding v to its flags. This should usually be
// called on the root command as as replacement for cmd.Execute().
//
// See pflagx.BindViper for more information.
func Execute(cmd *cobra.Command, v *viper.Viper) error {
	if cmd.PersistentPreRunE == nil && cmd.PersistentPreRun != nil {
		persistentPreRun := cmd.PersistentPreRun
		// Convert an existing cmd.PersistentPreRun func into
		// cmd.PersistentPreRunE to not lose it if we install the hook as cobra
		// will ignore cmd.PersistentPreRun if cmd.PersistentPreRunE is set.
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			persistentPreRun(cmd, args)
			return nil
		}
	}

	// Register the Intertwine hook as the first one in the PersistentPreRunE
	// chain.
	cmd.PersistentPreRunE = HookViper(v, cmd.PersistentPreRunE)

	return cmd.Execute()
}
