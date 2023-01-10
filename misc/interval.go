package misc

import (
	"time"

	"github.com/jackc/pgtype"
)

// Añade un pgtype.Interval a un pgtype.Timestamp
func AddInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if t.Status != pgtype.Present || i.Status != pgtype.Present {
		result.Status = pgtype.Null
		return
	}
	result.Time = t.Time.AddDate(0, int(i.Months), int(i.Days)).Add(time.Duration(i.Microseconds) * time.Microsecond)
	result.Status = pgtype.Present
	return
}

// Sustrae un pgtype.Interval de un pgtype.Timestamp
func SubInterval(t pgtype.Timestamp, i pgtype.Interval) (result pgtype.Timestamp) {
	if t.Status != pgtype.Present || i.Status != pgtype.Present {
		result.Status = pgtype.Null
		return
	}
	result.Time = t.Time.AddDate(0, -int(i.Months), -int(i.Days)).Add(-time.Duration(i.Microseconds) * time.Microsecond)
	result.Status = pgtype.Present
	return
}
