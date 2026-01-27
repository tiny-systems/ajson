package ajson

import (
	"testing"
)

// TestSingleArgStringFunctions tests single-argument string functions (upper, lower, trim, reverse)
func TestSingleArgStringFunctions(t *testing.T) {
	tests := []struct {
		name   string
		fname  string
		value  *Node
		result *Node
		fail   bool
	}{
		// upper
		{name: "upper: lowercase", fname: "upper", value: StringNode("", "hello"), result: StringNode("", "HELLO")},
		{name: "upper: uppercase", fname: "upper", value: StringNode("", "HELLO"), result: StringNode("", "HELLO")},
		{name: "upper: mixed", fname: "upper", value: StringNode("", "HeLLo WoRLd"), result: StringNode("", "HELLO WORLD")},
		{name: "upper: empty", fname: "upper", value: StringNode("", ""), result: StringNode("", "")},
		{name: "upper: unicode", fname: "upper", value: StringNode("", "héllo"), result: StringNode("", "HÉLLO")},
		{name: "upper: nil", fname: "upper", value: nil, result: NullNode("")},
		{name: "upper: numeric error", fname: "upper", value: NumericNode("", 123), fail: true},
		{name: "upper: bool error", fname: "upper", value: BoolNode("", true), fail: true},

		// lower
		{name: "lower: uppercase", fname: "lower", value: StringNode("", "HELLO"), result: StringNode("", "hello")},
		{name: "lower: lowercase", fname: "lower", value: StringNode("", "hello"), result: StringNode("", "hello")},
		{name: "lower: mixed", fname: "lower", value: StringNode("", "HeLLo WoRLd"), result: StringNode("", "hello world")},
		{name: "lower: empty", fname: "lower", value: StringNode("", ""), result: StringNode("", "")},
		{name: "lower: unicode", fname: "lower", value: StringNode("", "HÉLLO"), result: StringNode("", "héllo")},
		{name: "lower: nil", fname: "lower", value: nil, result: NullNode("")},
		{name: "lower: numeric error", fname: "lower", value: NumericNode("", 123), fail: true},
		{name: "lower: array error", fname: "lower", value: ArrayNode("", []*Node{}), fail: true},

		// trim
		{name: "trim: spaces", fname: "trim", value: StringNode("", "  hello  "), result: StringNode("", "hello")},
		{name: "trim: tabs", fname: "trim", value: StringNode("", "\thello\t"), result: StringNode("", "hello")},
		{name: "trim: newlines", fname: "trim", value: StringNode("", "\nhello\n"), result: StringNode("", "hello")},
		{name: "trim: mixed whitespace", fname: "trim", value: StringNode("", " \t\nhello \t\n"), result: StringNode("", "hello")},
		{name: "trim: no whitespace", fname: "trim", value: StringNode("", "hello"), result: StringNode("", "hello")},
		{name: "trim: empty", fname: "trim", value: StringNode("", ""), result: StringNode("", "")},
		{name: "trim: only whitespace", fname: "trim", value: StringNode("", "   \t\n  "), result: StringNode("", "")},
		{name: "trim: nil", fname: "trim", value: nil, result: NullNode("")},
		{name: "trim: numeric error", fname: "trim", value: NumericNode("", 123), fail: true},

		// reverse
		{name: "reverse: simple", fname: "reverse", value: StringNode("", "hello"), result: StringNode("", "olleh")},
		{name: "reverse: palindrome", fname: "reverse", value: StringNode("", "racecar"), result: StringNode("", "racecar")},
		{name: "reverse: empty", fname: "reverse", value: StringNode("", ""), result: StringNode("", "")},
		{name: "reverse: single char", fname: "reverse", value: StringNode("", "a"), result: StringNode("", "a")},
		{name: "reverse: unicode", fname: "reverse", value: StringNode("", "héllo"), result: StringNode("", "olléh")},
		{name: "reverse: with spaces", fname: "reverse", value: StringNode("", "hello world"), result: StringNode("", "dlrow olleh")},
		{name: "reverse: nil", fname: "reverse", value: nil, result: NullNode("")},
		{name: "reverse: numeric error", fname: "reverse", value: NumericNode("", 123), fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := functions[test.fname](test.value)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestMultiArgFunctionHelpers tests the helper functions for multi-arg functions
func TestMultiArgFunctionHelpers(t *testing.T) {
	t.Run("IsMultiArgFunction", func(t *testing.T) {
		// Should return true for known multi-arg functions
		knownFuncs := []string{"split", "join", "contains", "hasprefix", "hassuffix", "replace", "substr", "index"}
		for _, fname := range knownFuncs {
			if !IsMultiArgFunction(fname) {
				t.Errorf("IsMultiArgFunction(%s) should return true", fname)
			}
			// Test case insensitivity
			if !IsMultiArgFunction(fname[:1]+fname[1:]) {
				t.Errorf("IsMultiArgFunction should be case insensitive for %s", fname)
			}
		}

		// Should return false for unknown functions
		if IsMultiArgFunction("unknownfunc") {
			t.Error("IsMultiArgFunction(unknownfunc) should return false")
		}

		// Should return false for single-arg functions
		if IsMultiArgFunction("upper") {
			t.Error("IsMultiArgFunction(upper) should return false")
		}
	})

	t.Run("GetMultiArgFunction", func(t *testing.T) {
		fn, ok := GetMultiArgFunction("split")
		if !ok || fn == nil {
			t.Error("GetMultiArgFunction(split) should return a function")
		}

		_, ok = GetMultiArgFunction("unknownfunc")
		if ok {
			t.Error("GetMultiArgFunction(unknownfunc) should return false")
		}
	})

	t.Run("AddMultiArgFunction", func(t *testing.T) {
		name := "testmultiargfunc"
		if IsMultiArgFunction(name) {
			t.Error("test function should not exist yet")
		}
		AddMultiArgFunction(name, func(args []*Node) (*Node, error) {
			return NumericNode("", 42), nil
		})
		if !IsMultiArgFunction(name) {
			t.Error("test function should exist after adding")
		}
	})
}

// TestSplitFunction tests the split multi-arg function
func TestSplitFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{
			name:   "split: basic",
			args:   []*Node{StringNode("", "a,b,c"), StringNode("", ",")},
			result: ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b"), StringNode("", "c")}),
		},
		{
			name:   "split: single element",
			args:   []*Node{StringNode("", "hello"), StringNode("", ",")},
			result: ArrayNode("", []*Node{StringNode("", "hello")}),
		},
		{
			name:   "split: empty separator",
			args:   []*Node{StringNode("", "abc"), StringNode("", "")},
			result: ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b"), StringNode("", "c")}),
		},
		{
			name:   "split: empty string",
			args:   []*Node{StringNode("", ""), StringNode("", ",")},
			result: ArrayNode("", []*Node{StringNode("", "")}),
		},
		{
			name:   "split: space separator",
			args:   []*Node{StringNode("", "hello world foo"), StringNode("", " ")},
			result: ArrayNode("", []*Node{StringNode("", "hello"), StringNode("", "world"), StringNode("", "foo")}),
		},
		{
			name:   "split: multi-char separator",
			args:   []*Node{StringNode("", "a::b::c"), StringNode("", "::")},
			result: ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b"), StringNode("", "c")}),
		},
		{
			name: "split: too few args",
			args: []*Node{StringNode("", "hello")},
			fail: true,
		},
		{
			name: "split: first arg not string",
			args: []*Node{NumericNode("", 123), StringNode("", ",")},
			fail: true,
		},
		{
			name: "split: second arg not string",
			args: []*Node{StringNode("", "hello"), NumericNode("", 123)},
			fail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("split")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else {
				// Compare arrays element by element
				if !result.IsArray() {
					t.Error("Result should be an array")
					return
				}
				resultArr := result.Inheritors()
				expectedArr := test.result.Inheritors()
				if len(resultArr) != len(expectedArr) {
					t.Errorf("Array length mismatch: got %d, expected %d", len(resultArr), len(expectedArr))
					return
				}
				for i := range resultArr {
					if ok, _ := resultArr[i].Eq(expectedArr[i]); !ok {
						t.Errorf("Element %d mismatch: got %v, expected %v", i, resultArr[i].value.Load(), expectedArr[i].value.Load())
					}
				}
			}
		})
	}
}

// TestJoinFunction tests the join multi-arg function
func TestJoinFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{
			name:   "join: basic",
			args:   []*Node{ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b"), StringNode("", "c")}), StringNode("", ",")},
			result: StringNode("", "a,b,c"),
		},
		{
			name:   "join: single element",
			args:   []*Node{ArrayNode("", []*Node{StringNode("", "hello")}), StringNode("", ",")},
			result: StringNode("", "hello"),
		},
		{
			name:   "join: empty separator",
			args:   []*Node{ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b"), StringNode("", "c")}), StringNode("", "")},
			result: StringNode("", "abc"),
		},
		{
			name:   "join: empty array",
			args:   []*Node{ArrayNode("", []*Node{}), StringNode("", ",")},
			result: StringNode("", ""),
		},
		{
			name:   "join: space separator",
			args:   []*Node{ArrayNode("", []*Node{StringNode("", "hello"), StringNode("", "world")}), StringNode("", " ")},
			result: StringNode("", "hello world"),
		},
		{
			name:   "join: multi-char separator",
			args:   []*Node{ArrayNode("", []*Node{StringNode("", "a"), StringNode("", "b")}), StringNode("", "::")},
			result: StringNode("", "a::b"),
		},
		{
			name: "join: too few args",
			args: []*Node{ArrayNode("", []*Node{StringNode("", "hello")})},
			fail: true,
		},
		{
			name: "join: first arg not array",
			args: []*Node{StringNode("", "hello"), StringNode("", ",")},
			fail: true,
		},
		{
			name: "join: second arg not string",
			args: []*Node{ArrayNode("", []*Node{StringNode("", "a")}), NumericNode("", 123)},
			fail: true,
		},
		{
			name: "join: array element not string",
			args: []*Node{ArrayNode("", []*Node{StringNode("", "a"), NumericNode("", 123)}), StringNode("", ",")},
			fail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("join")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestContainsFunction tests the contains multi-arg function
func TestContainsFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "contains: found", args: []*Node{StringNode("", "hello world"), StringNode("", "world")}, result: BoolNode("", true)},
		{name: "contains: not found", args: []*Node{StringNode("", "hello world"), StringNode("", "foo")}, result: BoolNode("", false)},
		{name: "contains: at start", args: []*Node{StringNode("", "hello world"), StringNode("", "hello")}, result: BoolNode("", true)},
		{name: "contains: at end", args: []*Node{StringNode("", "hello world"), StringNode("", "world")}, result: BoolNode("", true)},
		{name: "contains: empty substring", args: []*Node{StringNode("", "hello"), StringNode("", "")}, result: BoolNode("", true)},
		{name: "contains: empty string", args: []*Node{StringNode("", ""), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "contains: both empty", args: []*Node{StringNode("", ""), StringNode("", "")}, result: BoolNode("", true)},
		{name: "contains: case sensitive", args: []*Node{StringNode("", "Hello"), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "contains: too few args", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "contains: first arg not string", args: []*Node{NumericNode("", 123), StringNode("", "1")}, fail: true},
		{name: "contains: second arg not string", args: []*Node{StringNode("", "hello"), NumericNode("", 123)}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("contains")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestHasPrefixFunction tests the hasprefix multi-arg function
func TestHasPrefixFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "hasprefix: true", args: []*Node{StringNode("", "hello world"), StringNode("", "hello")}, result: BoolNode("", true)},
		{name: "hasprefix: false", args: []*Node{StringNode("", "hello world"), StringNode("", "world")}, result: BoolNode("", false)},
		{name: "hasprefix: empty prefix", args: []*Node{StringNode("", "hello"), StringNode("", "")}, result: BoolNode("", true)},
		{name: "hasprefix: empty string", args: []*Node{StringNode("", ""), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "hasprefix: equal strings", args: []*Node{StringNode("", "hello"), StringNode("", "hello")}, result: BoolNode("", true)},
		{name: "hasprefix: case sensitive", args: []*Node{StringNode("", "Hello"), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "hasprefix: too few args", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "hasprefix: first arg not string", args: []*Node{NumericNode("", 123), StringNode("", "1")}, fail: true},
		{name: "hasprefix: second arg not string", args: []*Node{StringNode("", "hello"), NumericNode("", 123)}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("hasprefix")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestHasSuffixFunction tests the hassuffix multi-arg function
func TestHasSuffixFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "hassuffix: true", args: []*Node{StringNode("", "hello world"), StringNode("", "world")}, result: BoolNode("", true)},
		{name: "hassuffix: false", args: []*Node{StringNode("", "hello world"), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "hassuffix: empty suffix", args: []*Node{StringNode("", "hello"), StringNode("", "")}, result: BoolNode("", true)},
		{name: "hassuffix: empty string", args: []*Node{StringNode("", ""), StringNode("", "hello")}, result: BoolNode("", false)},
		{name: "hassuffix: equal strings", args: []*Node{StringNode("", "hello"), StringNode("", "hello")}, result: BoolNode("", true)},
		{name: "hassuffix: case sensitive", args: []*Node{StringNode("", "helloWorld"), StringNode("", "world")}, result: BoolNode("", false)},
		{name: "hassuffix: too few args", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "hassuffix: first arg not string", args: []*Node{NumericNode("", 123), StringNode("", "1")}, fail: true},
		{name: "hassuffix: second arg not string", args: []*Node{StringNode("", "hello"), NumericNode("", 123)}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("hassuffix")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestReplaceFunction tests the replace multi-arg function
func TestReplaceFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "replace: basic", args: []*Node{StringNode("", "hello world"), StringNode("", "world"), StringNode("", "there")}, result: StringNode("", "hello there")},
		{name: "replace: multiple", args: []*Node{StringNode("", "foo foo foo"), StringNode("", "foo"), StringNode("", "bar")}, result: StringNode("", "bar bar bar")},
		{name: "replace: not found", args: []*Node{StringNode("", "hello"), StringNode("", "xyz"), StringNode("", "abc")}, result: StringNode("", "hello")},
		{name: "replace: empty old", args: []*Node{StringNode("", "abc"), StringNode("", ""), StringNode("", "X")}, result: StringNode("", "XaXbXcX")},
		{name: "replace: empty new", args: []*Node{StringNode("", "hello"), StringNode("", "l"), StringNode("", "")}, result: StringNode("", "heo")},
		{name: "replace: empty string", args: []*Node{StringNode("", ""), StringNode("", "x"), StringNode("", "y")}, result: StringNode("", "")},
		{name: "replace: too few args 1", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "replace: too few args 2", args: []*Node{StringNode("", "hello"), StringNode("", "l")}, fail: true},
		{name: "replace: first arg not string", args: []*Node{NumericNode("", 123), StringNode("", "1"), StringNode("", "2")}, fail: true},
		{name: "replace: second arg not string", args: []*Node{StringNode("", "hello"), NumericNode("", 1), StringNode("", "2")}, fail: true},
		{name: "replace: third arg not string", args: []*Node{StringNode("", "hello"), StringNode("", "l"), NumericNode("", 2)}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("replace")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestSubstrFunction tests the substr multi-arg function
func TestSubstrFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "substr: from start", args: []*Node{StringNode("", "hello"), NumericNode("", 0)}, result: StringNode("", "hello")},
		{name: "substr: from middle", args: []*Node{StringNode("", "hello"), NumericNode("", 2)}, result: StringNode("", "llo")},
		{name: "substr: from end", args: []*Node{StringNode("", "hello"), NumericNode("", 5)}, result: StringNode("", "")},
		{name: "substr: negative index", args: []*Node{StringNode("", "hello"), NumericNode("", -2)}, result: StringNode("", "lo")},
		{name: "substr: negative beyond start", args: []*Node{StringNode("", "hello"), NumericNode("", -10)}, result: StringNode("", "hello")},
		{name: "substr: with length", args: []*Node{StringNode("", "hello"), NumericNode("", 1), NumericNode("", 3)}, result: StringNode("", "ell")},
		{name: "substr: length beyond end", args: []*Node{StringNode("", "hello"), NumericNode("", 3), NumericNode("", 10)}, result: StringNode("", "lo")},
		{name: "substr: zero length", args: []*Node{StringNode("", "hello"), NumericNode("", 2), NumericNode("", 0)}, result: StringNode("", "")},
		{name: "substr: empty string", args: []*Node{StringNode("", ""), NumericNode("", 0)}, result: StringNode("", "")},
		{name: "substr: too few args", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "substr: first arg not string", args: []*Node{NumericNode("", 123), NumericNode("", 0)}, fail: true},
		{name: "substr: second arg not number", args: []*Node{StringNode("", "hello"), StringNode("", "0")}, fail: true},
		{name: "substr: third arg not number", args: []*Node{StringNode("", "hello"), NumericNode("", 0), StringNode("", "3")}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("substr")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestIndexFunction tests the index multi-arg function
func TestIndexFunction(t *testing.T) {
	tests := []struct {
		name   string
		args   []*Node
		result *Node
		fail   bool
	}{
		{name: "index: found at start", args: []*Node{StringNode("", "hello"), StringNode("", "h")}, result: NumericNode("", 0)},
		{name: "index: found at middle", args: []*Node{StringNode("", "hello"), StringNode("", "ll")}, result: NumericNode("", 2)},
		{name: "index: found at end", args: []*Node{StringNode("", "hello"), StringNode("", "o")}, result: NumericNode("", 4)},
		{name: "index: not found", args: []*Node{StringNode("", "hello"), StringNode("", "xyz")}, result: NumericNode("", -1)},
		{name: "index: empty substring", args: []*Node{StringNode("", "hello"), StringNode("", "")}, result: NumericNode("", 0)},
		{name: "index: empty string", args: []*Node{StringNode("", ""), StringNode("", "hello")}, result: NumericNode("", -1)},
		{name: "index: both empty", args: []*Node{StringNode("", ""), StringNode("", "")}, result: NumericNode("", 0)},
		{name: "index: case sensitive", args: []*Node{StringNode("", "Hello"), StringNode("", "h")}, result: NumericNode("", -1)},
		{name: "index: too few args", args: []*Node{StringNode("", "hello")}, fail: true},
		{name: "index: first arg not string", args: []*Node{NumericNode("", 123), StringNode("", "1")}, fail: true},
		{name: "index: second arg not string", args: []*Node{StringNode("", "hello"), NumericNode("", 1)}, fail: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn, _ := GetMultiArgFunction("index")
			result, err := fn(test.args)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			} else if ok, err := result.Eq(test.result); !ok {
				if err != nil {
					t.Errorf("Unexpected error on comparison: %s", err.Error())
				}
				t.Errorf("Wrong value: %v != %v", result.value.Load(), test.result.value.Load())
			}
		})
	}
}

// TestStringFunctionsViaEval tests the string functions through the Eval function
// This ensures the RPN parser and eval correctly handle these functions
func TestStringFunctionsViaEval(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   interface{}
		fail       bool
	}{
		// Single-arg string functions
		{name: "eval upper", json: `{"name":"hello"}`, expression: `upper($.name)`, expected: "HELLO"},
		{name: "eval lower", json: `{"name":"HELLO"}`, expression: `lower($.name)`, expected: "hello"},
		{name: "eval trim", json: `{"name":"  hello  "}`, expression: `trim($.name)`, expected: "hello"},
		{name: "eval reverse", json: `{"name":"hello"}`, expression: `reverse($.name)`, expected: "olleh"},

		// Multi-arg string functions
		{name: "eval split", json: `{"data":"a,b,c"}`, expression: `length(split($.data, ','))`, expected: float64(3)},
		{name: "eval contains true", json: `{"text":"hello world"}`, expression: `contains($.text, 'world')`, expected: true},
		{name: "eval contains false", json: `{"text":"hello world"}`, expression: `contains($.text, 'foo')`, expected: false},
		{name: "eval hasprefix true", json: `{"text":"hello world"}`, expression: `hasprefix($.text, 'hello')`, expected: true},
		{name: "eval hasprefix false", json: `{"text":"hello world"}`, expression: `hasprefix($.text, 'world')`, expected: false},
		{name: "eval hassuffix true", json: `{"text":"hello world"}`, expression: `hassuffix($.text, 'world')`, expected: true},
		{name: "eval hassuffix false", json: `{"text":"hello world"}`, expression: `hassuffix($.text, 'hello')`, expected: false},
		{name: "eval replace", json: `{"text":"hello world"}`, expression: `replace($.text, 'world', 'there')`, expected: "hello there"},
		{name: "eval substr 2 args", json: `{"text":"hello"}`, expression: `substr($.text, 2)`, expected: "llo"},
		{name: "eval substr 3 args", json: `{"text":"hello"}`, expression: `substr($.text, 1, 3)`, expected: "ell"},
		{name: "eval index found", json: `{"text":"hello"}`, expression: `index($.text, 'l')`, expected: float64(2)},
		{name: "eval index not found", json: `{"text":"hello"}`, expression: `index($.text, 'x')`, expected: float64(-1)},

		// Combined operations
		{name: "eval upper + contains", json: `{"name":"Hello"}`, expression: `contains(upper($.name), 'ELLO')`, expected: true},
		{name: "eval split + first", json: `{"data":"a/b/c"}`, expression: `first(split($.data, '/'))`, expected: "a"},
		{name: "eval split + last", json: `{"data":"a/b/c"}`, expression: `last(split($.data, '/'))`, expected: "c"},

		// Ternary with string functions
		{name: "eval ternary contains", json: `{"text":"error: failed"}`, expression: `contains($.text, 'error') ? 'has error' : 'ok'`, expected: "has error"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			switch expected := test.expected.(type) {
			case string:
				val, err := result.GetString()
				if err != nil {
					t.Errorf("Expected string result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %q, expected %q", val, expected)
				}
			case bool:
				val, err := result.GetBool()
				if err != nil {
					t.Errorf("Expected bool result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			case float64:
				val, err := result.GetNumeric()
				if err != nil {
					t.Errorf("Expected numeric result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			default:
				t.Errorf("Unknown expected type: %T", test.expected)
			}
		})
	}
}

// TestSplitJoinRoundtrip tests that split and join are inverse operations
func TestSplitJoinRoundtrip(t *testing.T) {
	tests := []struct {
		original  string
		separator string
	}{
		{"a,b,c", ","},
		{"hello world", " "},
		{"one::two::three", "::"},
		{"single", ","},
		{"", ","},
	}

	splitFn, _ := GetMultiArgFunction("split")
	joinFn, _ := GetMultiArgFunction("join")

	for _, test := range tests {
		t.Run(test.original, func(t *testing.T) {
			// Split
			splitResult, err := splitFn([]*Node{StringNode("", test.original), StringNode("", test.separator)})
			if err != nil {
				t.Fatalf("Split failed: %s", err.Error())
			}

			// Join back
			joinResult, err := joinFn([]*Node{splitResult, StringNode("", test.separator)})
			if err != nil {
				t.Fatalf("Join failed: %s", err.Error())
			}

			// Compare
			resultStr, _ := joinResult.GetString()
			if resultStr != test.original {
				t.Errorf("Roundtrip failed: got %q, expected %q", resultStr, test.original)
			}
		})
	}
}

// TestNestedStringFunctions tests deeply nested function calls
func TestNestedStringFunctions(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   interface{}
	}{
		// Nested single-arg functions
		{
			name:       "upper(lower(upper))",
			json:       `{"text":"Hello"}`,
			expression: `upper(lower(upper($.text)))`,
			expected:   "HELLO",
		},
		{
			name:       "reverse(reverse)",
			json:       `{"text":"hello"}`,
			expression: `reverse(reverse($.text))`,
			expected:   "hello",
		},
		{
			name:       "trim(upper)",
			json:       `{"text":"  hello  "}`,
			expression: `upper(trim($.text))`,
			expected:   "HELLO",
		},

		// Nested multi-arg with single-arg
		{
			name:       "split with upper",
			json:       `{"text":"a,b,c"}`,
			expression: `first(split(upper($.text), ','))`,
			expected:   "A",
		},
		{
			name:       "contains with lower",
			json:       `{"text":"HELLO WORLD"}`,
			expression: `contains(lower($.text), 'hello')`,
			expected:   true,
		},
		{
			name:       "replace with trim",
			json:       `{"text":"  hello world  "}`,
			expression: `replace(trim($.text), 'world', 'there')`,
			expected:   "hello there",
		},

		// Complex expressions
		{
			name:       "split and get second element",
			json:       `{"path":"deployment/myapp"}`,
			expression: `last(split($.path, '/'))`,
			expected:   "myapp",
		},
		{
			name:       "split and get first element",
			json:       `{"path":"deployment/myapp"}`,
			expression: `first(split($.path, '/'))`,
			expected:   "deployment",
		},
		{
			name:       "conditional with hasprefix",
			json:       `{"cmd":"scale deployment/nginx --replicas=3"}`,
			expression: `hasprefix($.cmd, 'scale') ? 'scaling' : 'other'`,
			expected:   "scaling",
		},
		{
			name:       "extract namespace from command",
			json:       `{"cmd":"restart deployment/myapp -n production"}`,
			expression: `contains($.cmd, ' -n ') ? 'has namespace' : 'no namespace'`,
			expected:   "has namespace",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			switch expected := test.expected.(type) {
			case string:
				val, err := result.GetString()
				if err != nil {
					t.Errorf("Expected string result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %q, expected %q", val, expected)
				}
			case bool:
				val, err := result.GetBool()
				if err != nil {
					t.Errorf("Expected bool result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			}
		})
	}
}

// TestStringFunctionsWithJSONPath tests string functions combined with JSONPath queries
func TestStringFunctionsWithJSONPath(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   interface{}
	}{
		{
			name:       "split path from array element",
			json:       `{"items":[{"path":"a/b"},{"path":"c/d"}]}`,
			expression: `first(split($.items[0].path, '/'))`,
			expected:   "a",
		},
		{
			name:       "contains with array access",
			json:       `{"tags":["production","critical"]}`,
			expression: `contains($.tags[0], 'prod')`,
			expected:   true,
		},
		{
			name:       "replace in nested object",
			json:       `{"config":{"endpoint":"http://localhost:8080"}}`,
			expression: `replace($.config.endpoint, 'localhost', '127.0.0.1')`,
			expected:   "http://127.0.0.1:8080",
		},
		{
			name:       "upper with nested access",
			json:       `{"user":{"name":"john doe"}}`,
			expression: `upper($.user.name)`,
			expected:   "JOHN DOE",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			switch expected := test.expected.(type) {
			case string:
				val, err := result.GetString()
				if err != nil {
					t.Errorf("Expected string result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %q, expected %q", val, expected)
				}
			case bool:
				val, err := result.GetBool()
				if err != nil {
					t.Errorf("Expected bool result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			}
		})
	}
}

// TestStringFunctionsWithArithmetic tests string functions combined with arithmetic
func TestStringFunctionsWithArithmetic(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   float64
	}{
		{
			name:       "length of split result",
			json:       `{"csv":"a,b,c,d,e"}`,
			expression: `length(split($.csv, ','))`,
			expected:   5,
		},
		{
			name:       "index + 1",
			json:       `{"text":"hello"}`,
			expression: `index($.text, 'l') + 1`,
			expected:   3,
		},
		{
			name:       "substr length",
			json:       `{"text":"hello world"}`,
			expression: `length(substr($.text, 0, 5))`,
			expected:   5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			val, err := result.GetNumeric()
			if err != nil {
				t.Errorf("Expected numeric result: %s", err.Error())
				return
			}
			if val != test.expected {
				t.Errorf("Wrong value: got %v, expected %v", val, test.expected)
			}
		})
	}
}

// TestStringFunctionsEdgeCases tests edge cases for string functions
func TestStringFunctionsEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   interface{}
		fail       bool
	}{
		// Unicode handling
		{
			name:       "split unicode",
			json:       `{"text":"hello→world→foo"}`,
			expression: `length(split($.text, '→'))`,
			expected:   float64(3),
		},
		{
			name:       "upper unicode",
			json:       `{"text":"café"}`,
			expression: `upper($.text)`,
			expected:   "CAFÉ",
		},
		{
			name:       "reverse unicode",
			json:       `{"text":"hello世界"}`,
			expression: `reverse($.text)`,
			expected:   "界世olleh",
		},

		// Empty and whitespace
		{
			name:       "trim only spaces",
			json:       `{"text":"     "}`,
			expression: `trim($.text)`,
			expected:   "",
		},
		{
			name:       "split empty result elements",
			json:       `{"text":"a,,b"}`,
			expression: `length(split($.text, ','))`,
			expected:   float64(3),
		},

		// String literals in expressions
		{
			name:       "contains with literal",
			json:       `{}`,
			expression: `contains('hello world', 'world')`,
			expected:   true,
		},
		{
			name:       "split with literal",
			json:       `{}`,
			expression: `first(split('a:b:c', ':'))`,
			expected:   "a",
		},
		{
			name:       "replace with literal",
			json:       `{}`,
			expression: `replace('hello', 'l', 'L')`,
			expected:   "heLLo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if test.fail {
				if err == nil {
					t.Error("Expected error: nil given")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			switch expected := test.expected.(type) {
			case string:
				val, err := result.GetString()
				if err != nil {
					t.Errorf("Expected string result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %q, expected %q", val, expected)
				}
			case bool:
				val, err := result.GetBool()
				if err != nil {
					t.Errorf("Expected bool result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			case float64:
				val, err := result.GetNumeric()
				if err != nil {
					t.Errorf("Expected numeric result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			}
		})
	}
}

// TestRealWorldUseCases tests realistic use cases for the string functions
func TestRealWorldUseCases(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		expression string
		expected   interface{}
	}{
		// Kubernetes-like path parsing
		{
			name:       "parse kind from path",
			json:       `{"target":"deployment/nginx"}`,
			expression: `first(split($.target, '/'))`,
			expected:   "deployment",
		},
		{
			name:       "parse name from path",
			json:       `{"target":"deployment/nginx"}`,
			expression: `last(split($.target, '/'))`,
			expected:   "nginx",
		},
		{
			name:       "check if has namespace flag",
			json:       `{"args":"--replicas=3 -n production"}`,
			expression: `contains($.args, '-n ')`,
			expected:   true,
		},

		// URL manipulation
		{
			name:       "extract domain from URL",
			json:       `{"url":"https://example.com/path"}`,
			expression: `first(split(last(split($.url, '//')), '/'))`,
			expected:   "example.com",
		},

		// Label parsing
		{
			name:       "parse label key",
			json:       `{"label":"app=nginx"}`,
			expression: `first(split($.label, '='))`,
			expected:   "app",
		},
		{
			name:       "parse label value",
			json:       `{"label":"app=nginx"}`,
			expression: `last(split($.label, '='))`,
			expected:   "nginx",
		},

		// Error message processing
		{
			name:       "check error type",
			json:       `{"error":"NotFound: resource not found"}`,
			expression: `hasprefix($.error, 'NotFound')`,
			expected:   true,
		},
		{
			name:       "extract error type",
			json:       `{"error":"NotFound: resource not found"}`,
			expression: `first(split($.error, ':'))`,
			expected:   "NotFound",
		},

		// Log processing
		{
			name:       "check log level",
			json:       `{"log":"[ERROR] Failed to connect"}`,
			expression: `contains($.log, '[ERROR]')`,
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := Unmarshal([]byte(test.json))
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %s", err.Error())
			}

			result, err := Eval(root, test.expression)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
				return
			}

			switch expected := test.expected.(type) {
			case string:
				val, err := result.GetString()
				if err != nil {
					t.Errorf("Expected string result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %q, expected %q", val, expected)
				}
			case bool:
				val, err := result.GetBool()
				if err != nil {
					t.Errorf("Expected bool result: %s", err.Error())
					return
				}
				if val != expected {
					t.Errorf("Wrong value: got %v, expected %v", val, expected)
				}
			}
		})
	}
}
