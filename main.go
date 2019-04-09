// Command timeto prints the time difference between the wall time and the time given on the command
// line. Times at the zero date (0000-01-01) are assumed to be relative to current time, such as
// '8pm' or '12' or '23:45' and will be evaluated using the wall time's date.
//
// This tries to parse a handful of time formats, but may not be comprehensive because it's
// a one-off tool I use just for setting PagerDuty overrides in minutes.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var timeFormats = []string{
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05.999999999Z0700",
	"2006-01-02 15:04:05Z0700",
	"2006-01-02 15:04:05.999999999Z0700",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05.999999999Z07:00",
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	"3:04PM",
	"3:04pm",
	"3:04:05PM",
	"3:04:05pm",
	"3:04:05.999999999PM",
	"3:04:05.999999999pm",
	"3PM",
	"3pm",
	"3 PM",
	"3 pm",
	"15",
	"15:04",
	"15:04:05",
	"15:04:05.999999999",
	time.Kitchen,
	strings.ToLower(time.Kitchen),
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func main() {
	prog := filepath.Base(os.Args[0])
	if prog == "" {
		prog = "timeto"
	}

	now := time.Now()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <time|duration>...\n", prog)
	}
	flag.Parse()

	var buf bytes.Buffer
	for _, ts := range flag.Args() {
		t, err := parseTime(ts, now)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		t = t.UTC()
		if t.Year() == 0 && t.Day() == 1 && t.Month() == 1 {
			t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), now.Location())
			if t.Before(now) {
				t = t.AddDate(0, 0, 1)
			}
		}

		d := t.Sub(now)
		fmt.Fprintf(&buf, "%d\t%d\t%d\t%v\t%d\t%d\t%f\t%v\t%v\n",
			int64(d.Round(time.Minute)/time.Minute), // Minutes - rounded
			int64(d.Round(time.Second)/time.Second), // Seconds - rounded
			int64(d),                           // Nanoseconds
			d,                                  // Duration
			t.Unix(),                           // Seconds since Unix epoch
			t.UnixNano(),                       // Nanoseconds since Unix epoch (may not make sense past 2262, so shouldn't be a concern to anyone)
			d.Seconds(),                        // Seconds, floating point (milliseconds and lower units are fractional)
			t.Format(time.RFC3339Nano),         // RFC3339 time, UTC
			t.Local().Format(time.RFC3339Nano), // RFC3339 time, local
		)
	}
	buf.WriteTo(os.Stdout)
}

func parseTime(ts string, now time.Time) (time.Time, error) {
	for _, tf := range timeFormats {
		t, err := time.Parse(tf, ts)
		if err != nil {
			continue
		}
		return t, nil
	}
	if d, err := time.ParseDuration(ts); err == nil {
		return now.Add(d), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q", ts)
}
