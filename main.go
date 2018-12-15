// Command timeto prints the time difference between the wall time and the time given on the command
// line. Times at the zero date (0000-01-01) are assumed to be relative to current time, such as
// '8pm' or '12' or '23:45' and will be evaluated using the wall time's date.
//
// This tries to parse a handful of time formats, but may not be comprehensive because it's
// a one-off tool I use just for setting PagerDuty overrides in minutes.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var timeFormats = []string{
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.999999999Z07:00",
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
		fmt.Fprintf(os.Stderr, "Usage: %s <time>...\n", prog)
	}
	flag.Parse()

	for _, ts := range flag.Args() {
		if t, err := parseTime(ts); err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			t = t.UTC()
			if t.Year() == 0 && t.Day() == 1 && t.Month() == 1 {
				t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), now.Location())
				if t.Before(now) {
					t = t.AddDate(0, 0, 1)
				}
			}
			d := t.Sub(now)
			fmt.Printf("%d\t%d\t%d\t%v\n", int64(d/time.Minute), int64(d/time.Second), int64(d), d)
		}
	}
}

func parseTime(ts string) (time.Time, error) {
	for _, tf := range timeFormats {
		t, err := time.Parse(tf, ts)
		if err != nil {
			continue
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("cannot parse time %q", ts)
}
