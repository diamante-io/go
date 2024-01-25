package time

import (
	"strconv"
	goTime "time"
)

//Millis represents time as milliseconds since epoch without any timezone adjustments
type Millis int64

//MillisFromString generates a Millis struct from a string representing an int64
func MillisFromString(s string) (Millis, error) {
	millis, err := strconv.ParseInt(s, 10, 64)
	return Millis(int64(millis)), err
}

//MillisFromInt64 generates a Millis struct from given millis int64
func MillisFromInt64(millis int64) Millis {
	return Millis(millis)
}

//MillisFromTime generates a Millis struct from given go time
func MillisFromTime(t goTime.Time) Millis {
	return Millis(t.UTC().UnixNano() / int64(goTime.Millisecond))
}

//MillisFromSeconds generates a Millis struct from given seconds int64
func MillisFromSeconds(seconds int64) Millis {
	return Millis(seconds * 1000)
}

func (t Millis) increment(millisToAdd int64) Millis {
	return Millis(int64(t) + millisToAdd)
}

//IsNil returns true if the timeMillis has not been initialized to a date other then 0 from epoch
func (t Millis) IsNil() bool {
	return t == 0
}

//RoundUp returns a new Millis instance with a rounded up to d millis
func (t Millis) RoundUp(d int64) Millis {
	if d == 0 {
		return t
	}
	if int64(t)%d != 0 {
		return t.RoundDown(d).increment(d)
	}
	return t
}

//RoundUp returns a new ToInt64 instance with a down to d millis
func (t Millis) RoundDown(d int64) Millis {
	//round down to the nearest d
	return Millis(int64(int64(t)/d) * d)
}

//ToInt64 returns the actual int64 millis since epoch
func (t Millis) ToInt64() int64 {
	return int64(t)
}

//ToTime returns a go time.Time timestamp, UTC adjusted
func (t Millis) ToTime() goTime.Time {
	// Milliseconds   1510831636149
	// Nanoseconds    1510831636149000000
	// Unix sec arg   1510831636
	// Unix nsec arg            149000000
	return goTime.Unix(int64(t)/1000, int64(t)%1000*int64(goTime.Millisecond)).UTC()
}

//Now returns current time in millis
func Now() Millis {
	return MillisFromTime(goTime.Now())
}

func (t Millis) String() string {
	return strconv.FormatInt(t.ToInt64(), 10)
}
