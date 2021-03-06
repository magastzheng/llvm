package consts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llir/llvm/types"
)

// TODO: Track the upstream removal of HexIntConstant (ref: discussion with
// Sean on llvm-dev).

// Int represents an integer constant.
//
// Examples:
//    42, -37, true, false, [us]0x[0-9A-Fa-f]+
//
// References:
//    http://llvm.org/docs/LangRef.html#simple-constants
type Int struct {
	typ *types.Int
	x   int64
}

// NewInt returns an integer constant based on the given integer type and string
// representation.
func NewInt(typ types.Type, s string) (*Int, error) {
	// Verify integer type.
	v := new(Int)
	var ok bool
	v.typ, ok = typ.(*types.Int)
	if !ok {
		return nil, fmt.Errorf("invalid type %q for integer constant", typ)
	}
	size := v.typ.Size()
	if size > 64 {
		// TODO: Add support for large integer constants (e.g. above 64-bits).
		err := fmt.Sprintf("not yet implemented; support for %q integer constants (e.g. above 64-bits)", typ)
		panic(err)
	}

	// Parse boolean constant.
	if size == 1 {
		switch s {
		case "1", "true":
			v.x = 1
		case "0", "false":
			v.x = 0
		default:
			return nil, fmt.Errorf("invalid integer constant %q for boolean type", s)
		}
		return v, nil
	} else if s == "true" || s == "false" {
		return nil, fmt.Errorf("integer constant %q type mismatch; expected i1, got %v", s, typ)
	}

	// TODO: Implement support for the HexIntConstant representation:
	//    [us]0x[0-9A-Fa-f]+

	// Parse integer constant.
	var err error
	v.x, err = strconv.ParseInt(s, 10, size)
	if err != nil {
		return nil, fmt.Errorf("unable to parse integer constant %q; %v", s, err)
	}

	return v, nil
}

// Type returns the type of the value.
func (v *Int) Type() types.Type {
	return v.typ
}

// String returns a string representation of the integer, either as a signed
// integer (e.g. 42, -13) or as a boolean (e.g. true, false) depending on the
// type. The integer string representation is preceded by the type of the
// constant, e.g.
//
//    i1 true
//    i32 -13
//    i64 42
func (v *Int) String() string {
	s := ""
	if v.typ.Size() == 1 {
		switch v.x {
		case 1:
			s = "true"
		default:
			s = "false"
		}
	} else {
		s = strconv.FormatInt(v.x, 10)
	}

	return fmt.Sprintf("%s %s", v.Type(), s)
}

// Float represents a floating point constant.
//
// Examples:
//    123.45, 1.2345e2, 0x[KLMH]?[0-9A-Fa-f]+
//
// References:
//    http://llvm.org/docs/LangRef.html#simple-constants
type Float struct {
	typ *types.Float
	x   float64
}

// NewFloat returns a floating point constant based on the given floating point
// type and string representation.
func NewFloat(typ types.Type, s string) (*Float, error) {
	// Verify floating point type.
	v := new(Float)
	var ok bool
	v.typ, ok = typ.(*types.Float)
	if !ok {
		return nil, fmt.Errorf("invalid type %q for floating point constant", typ)
	}
	size := v.typ.Size()
	switch size {
	case 32, 64:
		// supported size
	default:
		// TODO: Add support for half, fp128, x86_fp80 and ppc_fp128.
		err := fmt.Sprintf("not yet implemented; support for %q floating point constants", v.typ)
		panic(err)
	}

	// TODO: Implement support for the following representation:
	//    0x[KLMH]?[0-9A-Fa-f]+

	// Parse floating point constant.
	var err error
	v.x, err = strconv.ParseFloat(s, size)
	if err != nil {
		return nil, fmt.Errorf("unable to parse floating point constant %q; %v", s, err)
	}

	// Verify that there was no precision loss.
	if size != 64 && strconv.FormatFloat(v.x, 'g', -1, size) != strconv.FormatFloat(v.x, 'g', -1, 64) {
		return nil, fmt.Errorf("invalid floating point constant %q for type %q; precision loss", s, v.typ)
	}

	return v, nil
}

// Type returns the type of the value.
func (v *Float) Type() types.Type {
	return v.typ
}

// String returns a string representation of the floating point constant using
// scientific notation (e.g. -2.5e10) for large exponents and regular floating
// point representation otherwise (e.g. 3.14). The floating point string
// representation is preceded by the type of the constant, e.g.
//
//    float 2.0
//    double 3.14
//    double -2.5e10
func (v *Float) String() string {
	size := v.typ.Size()
	switch size {
	case 32, 64:
		// supported size
	default:
		// TODO: Add support for half, fp128, x86_fp80 and ppc_fp128.
		err := fmt.Sprintf("not yet implemented; support for %q floating point constants", v.typ)
		panic(err)
	}

	// Insert decimal point if not present.
	//    3e4 -> 3.0e4
	//    42  -> 42.0
	s := strconv.FormatFloat(v.x, 'g', -1, size)
	if !strings.ContainsRune(s, '.') {
		pos := strings.IndexByte(s, 'e')
		if pos != -1 {
			s = s[:pos] + ".0" + s[pos:]
		} else {
			s = s + ".0"
		}
	}

	// Drop explicit plus sign in exponents.
	//    3.0e+4 -> 3.0e4
	s = strings.Replace(s, "e+", "e", -1)

	return fmt.Sprintf("%s %s", v.Type(), s)
}

// TODO: Check if global names are used for anything except functions and global
// variables. If so, be more specific about @foo in the example below by
// providing a comment.

// Pointer represents a pointer constant.
//
// Examples:
//    null, @foo
//
// References:
//    http://llvm.org/docs/LangRef.html#simple-constants
type Pointer struct {
	typ *types.Pointer
}

// TODO: Figure out how to represent pointer constants. Add the necessary fields
// to the Pointer struct and implement the NewPointer constructor afterwards.

// Type returns the type of the value.
func (v *Pointer) Type() types.Type {
	return v.typ
}

// isConst ensures that only constant values can be assigned to the Constant
// interface.
func (*Int) isConst()     {}
func (*Float) isConst()   {}
func (*Pointer) isConst() {}
