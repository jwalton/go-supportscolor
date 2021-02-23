package hasflag

import (
	"testing"
)

func checkHasFlag(t *testing.T, expected bool, flag string, argv []string) {
	if hasFlag(flag, argv) != expected {
		t.Errorf("Expected %v, got %v for flag %v in %v", expected, !expected, flag, argv)
	}
}

func TestHasFlag(t *testing.T) {
	checkHasFlag(t, true, "unicorn", []string{"--foo", "--unicorn", "--bar"})
	// Optional prefix.
	checkHasFlag(t, true, "--unicorn", []string{"--foo", "--unicorn", "--bar"})
	checkHasFlag(t, true, "unicorn=rainbow", []string{"--foo", "--unicorn=rainbow", "--bar"})
	checkHasFlag(t, true, "unicorn", []string{"--unicorn", "--", "--foo"})
	// Don't match flags after terminator.
	checkHasFlag(t, false, "unicorn", []string{"--foo", "--", "--unicorn"})
	checkHasFlag(t, false, "unicorn", []string{"--foo"})
	checkHasFlag(t, true, "-u", []string{"-f", "-u", "-b"})
	checkHasFlag(t, true, "-u", []string{"-u", "--", "-f"})
	checkHasFlag(t, true, "u", []string{"-f", "-u", "-b"})
	checkHasFlag(t, true, "u", []string{"-u", "--", "-f"})
	checkHasFlag(t, false, "f", []string{"-u", "--", "-f"})
}
