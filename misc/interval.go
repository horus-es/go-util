package misc

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Añade un pgtype.Interval a un pgtype.Timestamp
func AddInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if !t.Valid || !i.Valid {
		return
	}
	result.Time = t.Time.AddDate(0, int(i.Months), int(i.Days)).Add(time.Duration(i.Microseconds) * time.Microsecond)
	result.Valid = true
	return
}

// Sustrae un pgtype.Interval de un pgtype.Timestamp
func SubInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if !t.Valid || !i.Valid {
		return
	}
	result.Time = t.Time.AddDate(0, -int(i.Months), -int(i.Days)).Add(-time.Duration(i.Microseconds) * time.Microsecond)
	result.Valid = true
	return
}
