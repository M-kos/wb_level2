package domains

import "time"

type Event struct {
	ID          int
	UserID      int
	Title       string
	Description string
	Date        time.Time
}
