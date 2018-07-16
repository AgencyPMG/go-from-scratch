package sqlrepo

import (
	"database/sql/driver"
	"fmt"
	"time"
)

//UTCTime is a helper type that forces value-ing and scanning Time values
//to and from the UTC Timezone.
type UTCTime struct {
	Time time.Time
}

//Scan attempts to scan src into t.Time.
//If scanning is successful, t.Time will be in UTC Timezone via time.Time.UTC().
func (t *UTCTime) Scan(src interface{}) error {
	srcTime, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("sqlrepo: invalid src type %T for scanning into UTCTime", src)
	}

	t.Time = srcTime.UTC()

	return nil
}

//Value returns t.Time.In(time.UTC), nil.
func (t UTCTime) Value() (driver.Value, error) {
	return t.Time.In(time.UTC), nil
}
