package pflagx

import "github.com/spf13/pflag"

type funcValue func(string) error

// String implements the pflag.Value interface.
func (f funcValue) String() string { return "" }

// Set implements the pflag.Value interface.
func (f funcValue) Set(s string) error { return f(s) }

// Type implements the pflag.Value interface.
func (f funcValue) Type() string { return "string" }

// Func defines a flag with the specified name and usage string.
// Each time the flag is seen, fn is called with the value of the flag.
// If fn returns a non-nil error, it will be treated as a flag value parsing error.
func Func(fs *pflag.FlagSet, name, usage string, fn func(string) error) {
	fs.Var(funcValue(fn), name, usage)
}

// FuncP is like Func, but accepts a shorthand letter that can be used after a
// single dash.
func FuncP(fs *pflag.FlagSet, name, shorthand, usage string, fn func(string) error) {
	fs.VarP(funcValue(fn), name, shorthand, usage)
}
