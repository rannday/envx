package envx

import (
	"os"
	"testing"
)

func TestParseDotEnv_ParsesKeyValues(t *testing.T) {
	content := "FOO=bar\nBAR=baz"
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %q", values["FOO"])
	}
	if values["BAR"] != "baz" {
		t.Fatalf("expected BAR=baz, got %q", values["BAR"])
	}
}

func TestParseDotEnv_IgnoresCommentsAndEmptyLines(t *testing.T) {
	content := "\n# comment\nFOO=bar\n\n# another\nBAR=baz\n"
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 2 {
		t.Fatalf("expected 2 values, got %d", len(values))
	}
}

func TestParseDotEnv_StripsInlineCommentOutsideQuotes(t *testing.T) {
	content := "FOO=bar # this is a comment"
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %q", values["FOO"])
	}
}

func TestParseDotEnv_KeepsHashInsideQuotes(t *testing.T) {
	content := "FOO=\"abc #123\""
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values["FOO"] != "abc #123" {
		t.Fatalf("expected quoted hash to be preserved, got %q", values["FOO"])
	}
}

func TestParseDotEnv_SupportsSingleQuotes(t *testing.T) {
	content := "FOO='My App'"
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values["FOO"] != "My App" {
		t.Fatalf("expected single-quoted value, got %q", values["FOO"])
	}
}

func TestParseDotEnv_SupportsExportPrefix(t *testing.T) {
	content := "export FOO=bar"
	path := writeTempDotEnv(t, content)

	values, err := parseDotEnv(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if values["FOO"] != "bar" {
		t.Fatalf("expected export-prefixed key to parse, got %q", values["FOO"])
	}
}

func TestParseDotEnv_InvalidKeyErrors(t *testing.T) {
	content := " =bad"
	path := writeTempDotEnv(t, content)

	if _, err := parseDotEnv(path); err == nil {
		t.Fatal("expected invalid key error")
	}
}

func writeTempDotEnv(t *testing.T, content string) string {
	t.Helper()

	file, err := os.CreateTemp("", ".env")
	if err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(file.Name())
	})

	if _, err := file.WriteString(content); err != nil {
		t.Fatalf("failed writing temp file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("failed closing temp file: %v", err)
	}

	return file.Name()
}
