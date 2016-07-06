package config_test

import (
	"fmt"
	"os"
	"testing"

	config "github.com/UltimateSoftware/ultipkg/config"
	"github.com/nbio/st"
)

func TestMain(m *testing.M) {
	os.Setenv("FOO", "foo")
	os.Setenv("FOO_BAR", "foobar")
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	var examples = []struct {
		key   string
		value string
	}{
		{"foo", "foo"},
		{"FOO", "foo"},
		{"foo_bar", "foobar"},
		{"foo.bar", "foobar"},
		{"nope", ""},
	}

	for _, example := range examples {
		if config.Get(example.key) != example.value {
			t.Errorf("Expected %s == %s", example.key, example.value)
		}
	}
}

func TestGetFallback(t *testing.T) {
	val := config.Get("bad", "fallback1", "fallback2")

	st.Expect(t, val, "fallback1")
}

func TestGetDoesNotUseFallback(t *testing.T) {
	val := config.Get("foo", "fallback")

	st.Expect(t, val, "foo")
}

func TestMustGet(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Previous Panic was expected")
		}
	}()

	var examples = []struct {
		key         string
		shouldPanic bool
	}{
		{"foo", false},
		{"FOO", false},
		{"foo_bar", false},
		{"foo.bar", false},
		{"nope", true},
	}

	for _, example := range examples {
		config.MustGet(example.key)
		if example.shouldPanic {
			t.Errorf("MustGet should have paniced for missing key: %s", example.key)
		}
	}
}

func BenchmarkGetSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		config.Get("foo")
	}
}

func BenchmarkGetDotted(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		config.Get("foo.bar")
	}
}

func BenchmarkMustGetPanic(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Previous Panic was expected")
		}
	}()

	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		config.MustGet("missing")
	}
}

func BenchmarkMustGet(b *testing.B) {
	defer func() {
		if r := recover(); r != nil {
			b.Error("Should not have paniced for key: foo")
		}
	}()

	for i := 0; i < b.N; i++ {
		b.ReportAllocs()
		config.MustGet("foo")
	}
}
