package scale

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// startPostgres runs a real PostgreSQL 16 in a container for the test and removes it at the end.
// The test is skipped, not failed, when Docker is not running.
func startPostgres(t *testing.T) *sql.DB {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("cine"), postgres.WithUsername("cine"), postgres.WithPassword("test-only"),
		postgres.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, container)
	if err != nil {
		t.Fatal(err)
	}
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestIntegrationSeatStoreAgainstARealPostgres(t *testing.T) {
	ctx := context.Background()
	db := startPostgres(t)
	store := NewSeatStore(db)
	if err := store.Migrate(ctx); err != nil {
		t.Fatal(err)
	}

	t.Run("hold, confirm and release move the seats", func(t *testing.T) {
		for _, ref := range []string{"a", "b", "c"} {
			if err := store.Hold(ctx, ref); err != nil {
				t.Fatal(err)
			}
		}
		if err := store.Confirm(ctx, "a"); err != nil {
			t.Fatal(err)
		}
		if err := store.Release(ctx, "b"); err != nil {
			t.Fatal(err)
		}
		counts, err := store.Counts(ctx)
		if err != nil || counts["HELD"] != 1 || counts["CONFIRMED"] != 1 || counts["RELEASED"] != 1 {
			t.Fatalf("counts %v, error %v", counts, err)
		}
		t.Logf("counts after hold a, b, c, confirm a and release b: %v", counts)
	})

	// These are the cases that a mock of the store accepts and a real database refuses.
	t.Run("the primary key refuses a second hold of the same order", func(t *testing.T) {
		if err := store.Hold(ctx, "a"); !errors.Is(err, ErrAlreadyHeld) {
			t.Fatalf("got %v, want ErrAlreadyHeld", err)
		}
		t.Log("second hold of order a: ErrAlreadyHeld (SQLSTATE 23505, pk_seat_hold)")
	})

	t.Run("the named CHECK refuses a status outside the list", func(t *testing.T) {
		_, err := db.ExecContext(ctx, `INSERT INTO seat_hold (order_ref, status) VALUES ('x', 'BOGUS')`)
		if err == nil || !strings.Contains(err.Error(), "ck_seat_hold_status") {
			t.Fatalf("got %v, want a violation of ck_seat_hold_status", err)
		}
		t.Logf("insert with status BOGUS: %v", err)
	})

	t.Run("a seat that is no longer held cannot be moved", func(t *testing.T) {
		for _, move := range []func(context.Context, string) error{store.Confirm, store.Release} {
			if err := move(ctx, "b"); !errors.Is(err, ErrNotFound) {
				t.Fatalf("a released seat: got %v, want ErrNotFound", err)
			}
		}
		if err := store.Confirm(ctx, "unknown"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("an unknown order: got %v, want ErrNotFound", err)
		}
		t.Log("confirm and release of the released order b, and confirm of an unknown order: ErrNotFound")
	})
}
