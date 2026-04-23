package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("DATABASE_MAX_CONNS", "50")
	os.Setenv("DATABASE_CONN_MAX_LIFETIME", "10m")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected 9090, got %s", cfg.Port)
	}
	if cfg.DatabaseMaxConns != 50 {
		t.Errorf("expected 50, got %d", cfg.DatabaseMaxConns)
	}
	if cfg.DatabaseConnLifetime != 10*time.Minute {
		t.Errorf("expected 10m, got %v", cfg.DatabaseConnLifetime)
	}

	os.Clearenv()
	_, err = Load()
	if err == nil {
		t.Error("expected validation error for missing DATABASE_URL")
	}
}

func TestGetEnvBool(t *testing.T) {
	if getEnvBool("NON_EXISTENT", true) != true {
		t.Error("expected default true")
	}
	os.Setenv("TEST_BOOL", "true")
	if getEnvBool("TEST_BOOL", false) != true {
		t.Error("expected true")
	}
	os.Setenv("TEST_BOOL", "0")
	if getEnvBool("TEST_BOOL", true) != false {
		t.Error("expected false")
	}
	os.Clearenv()
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "invalid")
	if getEnvInt("TEST_INT", 10) != 10 {
		t.Error("expected fallback to 10")
	}
	os.Clearenv()
}

func TestGetEnvDuration(t *testing.T) {
	os.Setenv("TEST_DUR", "invalid")
	if getEnvDuration("TEST_DUR", time.Minute) != time.Minute {
		t.Error("expected fallback to 1m")
	}
	os.Clearenv()
}
