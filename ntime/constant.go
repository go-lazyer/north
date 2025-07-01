package ntime

import "time"

// layout constants
const (
	AtomLayout     = RFC3339Layout
	ANSICLayout    = time.ANSIC
	CookieLayout   = "Monday, 02-Jan-2006 15:04:05 MST"
	KitchenLayout  = time.Kitchen
	RssLayout      = time.RFC1123Z
	RubyDateLayout = time.RubyDate
	UnixDateLayout = time.UnixDate
	W3cLayout      = RFC3339Layout

	RFC1036Layout      = "Mon, 02 Jan 06 15:04:05 -0700"
	RFC1123Layout      = time.RFC1123
	RFC1123ZLayout     = time.RFC1123Z
	RFC2822Layout      = time.RFC1123Z
	RFC3339Layout      = "2006-01-02T15:04:05Z07:00"
	RFC3339MilliLayout = "2006-01-02T15:04:05.999Z07:00"
	RFC3339MicroLayout = "2006-01-02T15:04:05.999999Z07:00"
	RFC3339NanoLayout  = "2006-01-02T15:04:05.999999999Z07:00"
	RFC7231Layout      = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC822Layout       = time.RFC822
	RFC822ZLayout      = time.RFC822Z
	RFC850Layout       = time.RFC850

	ISO8601Layout      = "2006-01-02T15:04:05-07:00"
	ISO8601MilliLayout = "2006-01-02T15:04:05.999-07:00"
	ISO8601MicroLayout = "2006-01-02T15:04:05.999999-07:00"
	ISO8601NanoLayout  = "2006-01-02T15:04:05.999999999-07:00"

	ISO8601ZuluLayout      = "2006-01-02T15:04:05Z"
	ISO8601ZuluMilliLayout = "2006-01-02T15:04:05.999Z"
	ISO8601ZuluMicroLayout = "2006-01-02T15:04:05.999999Z"
	ISO8601ZuluNanoLayout  = "2006-01-02T15:04:05.999999999Z"

	FormattedDateLayout    = "Jan 2, 2006"
	FormattedDayDateLayout = "Mon, Jan 2, 2006"

	DayDateTimeLayout        = "Mon, Jan 2, 2006 3:04 PM"
	DateTimeLayout           = "2006-01-02 15:04:05"
	DateTimeMilliLayout      = "2006-01-02 15:04:05.999"
	DateTimeMicroLayout      = "2006-01-02 15:04:05.999999"
	DateTimeNanoLayout       = "2006-01-02 15:04:05.999999999"
	ShortDateTimeLayout      = "20060102150405"
	ShortDateTimeMilliLayout = "20060102150405.999"
	ShortDateTimeMicroLayout = "20060102150405.999999"
	ShortDateTimeNanoLayout  = "20060102150405.999999999"

	DateLayout           = "2006-01-02"
	DateMilliLayout      = "2006-01-02.999"
	DateMicroLayout      = "2006-01-02.999999"
	DateNanoLayout       = "2006-01-02.999999999"
	ShortDateLayout      = "20060102"
	ShortDateMilliLayout = "20060102.999"
	ShortDateMicroLayout = "20060102.999999"
	ShortDateNanoLayout  = "20060102.999999999"

	TimeLayout           = "15:04:05"
	TimeMilliLayout      = "15:04:05.999"
	TimeMicroLayout      = "15:04:05.999999"
	TimeNanoLayout       = "15:04:05.999999999"
	ShortTimeLayout      = "150405"
	ShortTimeMilliLayout = "150405.999"
	ShortTimeMicroLayout = "150405.999999"
	ShortTimeNanoLayout  = "150405.999999999"

	TimestampLayout      = "unix"
	TimestampMilliLayout = "unixMilli"
	TimestampMicroLayout = "unixMicro"
	TimestampNanoLayout  = "unixNano"
)

var defaultLayouts = []string{
	DateTimeLayout, DateLayout, TimeLayout, DayDateTimeLayout,

	"2006-01-02 15:04:05 -0700 MST", "2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05-0700", "2006-01-02T15:04:05",

	ISO8601Layout, RFC1036Layout, RFC822Layout, RFC822ZLayout, RFC850Layout, RFC1123Layout, RFC1123ZLayout, RFC3339Layout, RFC7231Layout,
	KitchenLayout, CookieLayout, ANSICLayout, UnixDateLayout, RubyDateLayout,

	ShortDateTimeLayout, ShortDateLayout, ShortTimeLayout,

	DateTimeMilliLayout, DateTimeMicroLayout, DateTimeNanoLayout,
	DateMilliLayout, DateMicroLayout, DateNanoLayout,
	TimeMilliLayout, TimeMicroLayout, TimeNanoLayout,

	ShortDateTimeMilliLayout, ShortDateTimeMicroLayout, ShortDateTimeNanoLayout,
	ShortDateMilliLayout, ShortDateMicroLayout, ShortDateNanoLayout,
	ShortTimeMilliLayout, ShortTimeMicroLayout, ShortTimeNanoLayout,

	ISO8601MilliLayout, ISO8601MicroLayout, ISO8601NanoLayout,
	RFC3339MilliLayout, RFC3339MicroLayout, RFC3339NanoLayout,

	"15:04:05-07",                          // postgres time with time zone type
	"2006-01-02 15:04:05-07",               // postgres timestamp with time zone type
	"2006-01-02 15:04:05-07:00",            // sqlite text type
	"2006-01-02 15:04:05.999999999 -07:00", // sqlserver datetimeoffset type
	"2006",
	"2006-1-2 15:4:5 -0700 MST", "2006-1-2 3:4:5 -0700 MST",
	"2006-1", "2006-1-2", "2006-1-2 15", "2006-1-2 15:4", "2006-1-2 15:4:5", "2006-1-2 15:4:5.999999999",
	"2006.1", "2006.1.2", "2006.1.2 15", "2006.1.2 15:4", "2006.1.2 15:4:5", "2006.1.2 15:4:5.999999999",
	"2006/1", "2006/1/2", "2006/1/2 15", "2006/1/2 15:4", "2006/1/2 15:4:5", "2006/1/2 15:4:5.999999999",
	"2006-01-02 15:04:05PM MST", "2006-01-02 15:04:05.999999999PM MST", "2006-1-2 15:4:5PM MST", "2006-1-2 15:4:5.999999999PM MST",
	"2006-01-02 15:04:05 PM MST", "2006-01-02 15:04:05.999999999 PM MST", "2006-1-2 15:4:5 PM MST", "2006-1-2 15:4:5.999999999 PM MST",
	"1/2/2006", "1/2/2006 15", "1/2/2006 15:4", "1/2/2006 15:4:5", "1/2/2006 15:4:5.999999999",
	"2006-1-2 15:4:5.999999999 -0700 MST", "2006-1-2 15:04:05 -0700 MST", "2006-1-2 15:04:05.999999999 -0700 MST",
	"2006-01-02T15:04:05.999999999", "2006-1-2T3:4:5", "2006-1-2T3:4:5.999999999",
	"2006-01-02T15:04:05.999999999Z07", "2006-1-2T15:4:5Z07", "2006-1-2T15:4:5.999999999Z07",
	"2006-01-02T15:04:05.999999999Z07:00", "2006-1-2T15:4:5Z07:00", "2006-1-2T15:4:5.999999999Z07:00",
	"2006-01-02T15:04:05.999999999-07:00", "2006-1-2T15:4:5-07:00", "2006-1-2T15:4:5.999999999-07:00",
	"2006-01-02T15:04:05.999999999-0700", "2006-1-2T3:4:5-0700", "2006-1-2T3:4:5.999999999-0700",
	"20060102150405-07:00", "20060102150405.999999999-07:00",
	"20060102150405Z07:00", "20060102150405.999999999Z07:00",
}
