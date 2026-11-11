package scale

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrAlreadyHeld = errors.New("a seat is already held for this order")
	ErrNotFound    = errors.New("there is no held seat for this order")
)

// schema follows the conventions of the project: a singular snake_case table, text with a named CHECK and explicit constraint names.
const schema = `CREATE TABLE seat_hold (
	order_ref text NOT NULL,
	status    text NOT NULL,
	CONSTRAINT pk_seat_hold PRIMARY KEY (order_ref),
	CONSTRAINT ck_seat_hold_status CHECK (status IN ('HELD', 'CONFIRMED', 'RELEASED'))
)`

// SeatStore keeps the seats of the checkout in PostgreSQL.
type SeatStore struct{ db *sql.DB }

func NewSeatStore(db *sql.DB) *SeatStore { return &SeatStore{db: db} }

func (s *SeatStore) Migrate(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, schema)
	return err
}

// Hold takes a seat for an order. A second hold for the same order is refused by the primary key.
func (s *SeatStore) Hold(ctx context.Context, ref string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO seat_hold (order_ref, status) VALUES ($1, 'HELD')`, ref)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return ErrAlreadyHeld
	}
	return err
}

func (s *SeatStore) Confirm(ctx context.Context, ref string) error {
	return s.move(ctx, ref, "CONFIRMED")
}

func (s *SeatStore) Release(ctx context.Context, ref string) error {
	return s.move(ctx, ref, "RELEASED")
}

// move changes a seat that is still held; a seat that is confirmed, released or unknown is not found.
func (s *SeatStore) move(ctx context.Context, ref, status string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE seat_hold SET status = $2 WHERE order_ref = $1 AND status = 'HELD'`, ref, status)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// Counts returns how many seats are in each status.
func (s *SeatStore) Counts(ctx context.Context) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT status, count(*) FROM seat_hold GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var status string
		var n int
		if err := rows.Scan(&status, &n); err != nil {
			return nil, err
		}
		counts[status] = n
	}
	return counts, rows.Err()
}
