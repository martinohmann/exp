package pflagx

import (
	"fmt"

	"github.com/spf13/pflag"
)

// TransformerFunc is a func that can be registered via RegisterTransformerFunc
// to transform flag values before they are set.
type TransformerFunc func(val string) string

type transformingValue struct {
	pflag.Value
	fn TransformerFunc
}

func newTransformingValue(value pflag.Value, fn TransformerFunc) *transformingValue {
	return &transformingValue{
		Value: value,
		fn:    fn,
	}
}

// Set implements the pflag.Value interface.
func (f *transformingValue) Set(s string) error {
	return f.Value.Set(f.fn(s))
}

// RegisterTransformerFunc registers a TransformerFunc for the named flag on
// fs. The func is invoked before setting the flag value. Panics if fn is nil,
// fs does not contain a flag with name or if the FlagSet is already parsed.
func RegisterTransformerFunc(fs *pflag.FlagSet, name string, fn TransformerFunc) {
	if fn == nil {
		panic("pflagx.RegisterTransformerFunc: nil transformer func")
	}

	if fs.Parsed() {
		panic("pflagx.RegisterTransformerFunc: must be invoked before fs.Parse()")
	}

	flag := fs.Lookup(name)
	if flag == nil {
		panic(fmt.Sprintf("pflagx.RegisterTransformerFunc: flag %q not defined", name))
	}

	flag.Value = newTransformingValue(flag.Value, fn)
}
