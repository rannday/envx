package envx

import (
	"testing"
	"time"
)

type testConfig struct {
	Port  int    `env:"PORT" default:"8080"`
	Host  string `env:"HOST" required:"true"`
	Debug bool   `env:"DEBUG" default:"false"`
}

func TestLoad_Success(t *testing.T) {
	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "9090")

	var cfg testConfig

	err := Load(&cfg, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 9090 {
		t.Fatalf("expected Port=9090, got %d", cfg.Port)
	}

	if cfg.Host != "localhost" {
		t.Fatalf("expected Host=localhost, got %s", cfg.Host)
	}

	if cfg.Debug != false {
		t.Fatalf("expected Debug=false, got %v", cfg.Debug)
	}
}

func TestLoad_DefaultValue(t *testing.T) {
	t.Setenv("HOST", "localhost")

	var cfg testConfig

	err := Load(&cfg, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Port != 8080 {
		t.Fatalf("expected default Port=8080, got %d", cfg.Port)
	}
}

func TestLoad_RequiredMissing(t *testing.T) {
	var cfg testConfig

	err := Load(&cfg, Options{})
	if err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestLoad_InvalidInt(t *testing.T) {
	t.Setenv("HOST", "localhost")
	t.Setenv("PORT", "notanumber")

	var cfg testConfig

	err := Load(&cfg, Options{})
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestLoad_NotPointer(t *testing.T) {
	var cfg testConfig

	err := Load(cfg, Options{})
	if err == nil {
		t.Fatal("expected error when passing non-pointer")
	}
}

func TestLoad_Int64(t *testing.T) {
	type cfg struct {
		Timeout int64 `env:"TIMEOUT" default:"60"`
	}

	var c cfg

	err := Load(&c, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Timeout != 60 {
		t.Fatalf("expected 60, got %d", c.Timeout)
	}
}

func TestLoad_Duration(t *testing.T) {
	type cfg struct {
		Timeout time.Duration `env:"TIMEOUT" default:"5s"`
	}

	var c cfg

	err := Load(&c, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Timeout != 5*time.Second {
		t.Fatalf("expected 5s, got %v", c.Timeout)
	}
}

func TestLoad_StringSlice(t *testing.T) {
	type cfg struct {
		Hosts []string `env:"HOSTS" default:"a.com,b.com"`
	}

	var c cfg

	err := Load(&c, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(c.Hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(c.Hosts))
	}

	if c.Hosts[0] != "a.com" || c.Hosts[1] != "b.com" {
		t.Fatalf("unexpected slice values: %v", c.Hosts)
	}
}

func TestLoad_OptionalEmptyInt(t *testing.T) {
	type cfg struct {
		Port int `env:"PORT"`
	}

	var c cfg

	err := Load(&c, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.Port != 0 {
		t.Fatalf("expected zero value, got %d", c.Port)
	}
}

func TestLoad_UnexportedIgnored(t *testing.T) {
	type cfg struct {
		host string `env:"HOST"`
	}

	var c cfg

	err := Load(&c, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
