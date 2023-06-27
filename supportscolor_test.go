package supportscolor

import (
	"runtime"
	"testing"
)

type testEnvironment struct {
	flags              []string
	env                map[string]string
	isNotTerminal      bool
	winMajorVersion    uint32
	winMinorVersion    uint32
	winBuildNumber     uint32
	colorCantBeEnabled bool
	colorWasEnabled    bool
	goos               string
}

func (test *testEnvironment) LookupEnv(name string) (string, bool) {
	val, present := test.env[name]
	return val, present
}

func (test *testEnvironment) Getenv(name string) string {
	return test.env[name]
}

func (test *testEnvironment) HasFlag(flag string) bool {
	for _, f := range test.flags {
		if f == flag {
			return true
		}
	}
	return false
}

func (test *testEnvironment) IsTerminal(fd int) bool {
	return !test.isNotTerminal
}

func (test *testEnvironment) getWindowsVersion() (majorVersion, minorVersion, buildNumber uint32) {
	return test.winMajorVersion, test.winMinorVersion, test.winBuildNumber
}

func (test *testEnvironment) osEnableColor() bool {
	test.colorWasEnabled = true
	return !test.colorCantBeEnabled
}

func (test *testEnvironment) getGOOS() string {
	if test.goos != "" {
		return test.goos
	}
	return runtime.GOOS
}

func TestReturnBasicIfForceColorAndNotTty(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:           map[string]string{"FORCE_COLOR": "true"},
		isNotTerminal: true,
	}))

	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestForceColorAnd256Flag(t *testing.T) {
	// return true if `FORCE_COLOR` is in env, but honor 256 flag
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{"FORCE_COLOR": "true"},
		flags: []string{"color=256"},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{"FORCE_COLOR": "1"},
		flags: []string{"color=256"},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestForceColorIs0(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "0"},
	}))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestDoNotCacheForceColor(t *testing.T) {
	env := map[string]string{"FORCE_COLOR": "0"}

	result := SupportsColor(0, setEnvironment(&testEnvironment{env: env}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}

	env["FORCE_COLOR"] = "1"
	result = SupportsColor(0, setEnvironment(&testEnvironment{env: env}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestNoColor(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:           map[string]string{"NO_COLOR": ""},
		isNotTerminal: true,
	}))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestReturnNoneIfNotTty(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{isNotTerminal: true}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestReturnNoneIfNoColorFlagIsUsed(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{"TERM": "xterm-256color"},
		flags: []string{"no-color"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestReturnNoneIfNoColorsFlagIsUsed(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{"TERM": "xterm-256color"},
		flags: []string{"no-colors"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestReturnBasicIfColorFlagIsUsed(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestReturnBasicIfColorsFlagIsUsed(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"colors"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestReturnBasicIfColorTermInEnv(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{"COLORTERM": "true"},
		flags: []string{},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestSupportColorTrueFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=true"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestSupportColorAlwaysFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=always"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestSupportColorFalseFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=false"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestSupportColor256Flag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=256"},
	}))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
	if result.Has256 == false {
		t.Errorf("Expected Has256 to be true")
	}
}

func TestSupportColor16mFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=16m"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
	if result.Has256 == false {
		t.Errorf("Expected Has256 to be true")
	}
	if result.Has16m == false {
		t.Errorf("Expected Has16m to be true")
	}
}

func TestSupportColorFullFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=full"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
	if result.Has256 == false {
		t.Errorf("Expected Has256 to be true")
	}
	if result.Has16m == false {
		t.Errorf("Expected Has16m to be true")
	}
}

func TestSupportColorTruecolorFlag(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env:   map[string]string{},
		flags: []string{"color=truecolor"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
	if result.Has256 == false {
		t.Errorf("Expected Has256 to be true")
	}
	if result.Has16m == false {
		t.Errorf("Expected Has16m to be true")
	}
}

func TestReturnNoneIfCIIsInEnv(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"CI": ""},
	}))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestReturnBasicIfTravis(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"CI": "Travis", "TRAVIS": "1"},
	}))

	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestReturnBasicIfCI(t *testing.T) {
	for _, ci := range []string{"CIRCLECI", "APPVEYOR", "GITLAB_CI", "GITHUB_ACTIONS", "GITEA_ACTIONS", "BUILDKITE", "DRONE"} {
		result := SupportsColor(0, setEnvironment(&testEnvironment{
			env: map[string]string{"CI": "true", ci: "true"},
		}))

		if result.Level != Basic {
			t.Errorf("%v: Expected %v, got %v", ci, Basic, result.Level)
		}
	}
}

func TestReturnBasicIfCodeship(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"CI": "true", "CI_NAME": "codeship"},
	}))

	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestTeamcity(t *testing.T) {
	// < 9.1
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TEAMCITY_VERSION": "9.0.5 (build 32523)"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}

	// >= 9.1
	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TEAMCITY_VERSION": "9.1.0 (build 32523)"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestRxvt(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "rxvt"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestPreferLevel2XtermOverColorTerm(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "xterm-256color", "COLORTERM": "1"},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestScreen256Color(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "screen-256color"},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestPutty256Color(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "putty-256color"},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestITerm(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM_PROGRAM": "iTerm.app", "TERM_PROGRAM_VERSION": "3.0.10"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM_PROGRAM": "iTerm.app", "TERM_PROGRAM_VERSION": "2.9.3"},
	}))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestWindows(t *testing.T) {
	// return level 1 if on Windows earlier than 10 build 10586
	env := &testEnvironment{
		goos:               "windows",
		colorCantBeEnabled: true,
		winMajorVersion:    10,
		winMinorVersion:    0,
		winBuildNumber:     10240,
	}
	result := SupportsColor(0, setEnvironment(env))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}

	// return level 2 if on Windows 10 build 10586 or later
	env = &testEnvironment{
		goos:            "windows",
		winMajorVersion: 10,
		winMinorVersion: 0,
		winBuildNumber:  10586,
	}
	result = SupportsColor(0, setEnvironment(env))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
	if !env.colorWasEnabled {
		t.Errorf("Expected color to be enabled")
	}

	// return level 3 if on Windows 10 build 14931 or later
	env = &testEnvironment{
		goos:            "windows",
		winMajorVersion: 10,
		winMinorVersion: 0,
		winBuildNumber:  14931,
	}
	result = SupportsColor(0, setEnvironment(env))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
	if !env.colorWasEnabled {
		t.Errorf("Expected color to be enabled")
	}

	// return level 2 if on Windows and force color flag
	env = &testEnvironment{
		flags:           []string{"color=256"},
		goos:            "windows",
		winMajorVersion: 10,
		winMinorVersion: 0,
		winBuildNumber:  10586,
	}
	result = SupportsColor(0, setEnvironment(env))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
	if !env.colorWasEnabled {
		t.Errorf("Expected color to be enabled")
	}

	// return level 3 if on Windows and force color flag
	env = &testEnvironment{
		flags:           []string{"color=16m"},
		goos:            "windows",
		winMajorVersion: 10,
		winMinorVersion: 0,
		winBuildNumber:  10586,
	}
	result = SupportsColor(0, setEnvironment(env))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
	if !env.colorWasEnabled {
		t.Errorf("Expected color to be enabled")
	}
}

func TestReturnAnsi256WhenForceColorIsSetWhenNotTTYInXterm256(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		isNotTerminal: true,
		env: map[string]string{
			"FORCE_COLOR": "true",
			"TERM":        "xterm-256color",
		},
	}))

	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}
}

func TestSupportsSettingAColorLevelUsingForceColor(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "0"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "1"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "2"},
	}))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "3"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "4"},
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}
}

func TestForceColorWorksWhenSetViaCommandLine(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "true"},
	}))
	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "true", "TERM": "xterm-256color"},
	}))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}

	result = SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"FORCE_COLOR": "false"},
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}
func TestTermIsDumb(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "dumb"},
	}))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestTermIsDumbAndTermProgramIsSet(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "dumb", "TERM_PROGRAM": "Apple_Terminal"},
	}))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestTermIsDumbAndWindows(t *testing.T) {
	env := &testEnvironment{
		env:             map[string]string{"TERM": "dumb", "TERM_PROGRAM": "Apple_Terminal"},
		goos:            "windows",
		winMajorVersion: 10,
		winMinorVersion: 0,
		winBuildNumber:  14931,
	}
	result := SupportsColor(0, setEnvironment(env))

	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
	if env.colorWasEnabled {
		t.Errorf("Expected color to not be enabled")
	}
}

func TestTermIsDumbAndForceColorIsSet(t *testing.T) {
	result := SupportsColor(0, setEnvironment(&testEnvironment{
		env: map[string]string{"TERM": "dumb", "FORCE_COLOR": "1"},
	}))

	if result.Level != Basic {
		t.Errorf("Expected %v, got %v", Basic, result.Level)
	}
}

func TestIgnoreFlagsWhenSniffFlagsIsFalse(t *testing.T) {
	env := &testEnvironment{
		env:   map[string]string{"TERM": "dumb"},
		flags: []string{"color=256"},
	}

	result := SupportsColor(0, setEnvironment(env))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}

	result = SupportsColor(0, SniffFlagsOption(true), setEnvironment(env))
	if result.Level != Ansi256 {
		t.Errorf("Expected %v, got %v", Ansi256, result.Level)
	}

	result = SupportsColor(0, SniffFlagsOption(false), setEnvironment(env))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}

func TestForceTTY(t *testing.T) {
	// Should be able to force that we are a terminal even when we detect we are not.
	result := SupportsColor(0, IsTTYOption(true), setEnvironment(&testEnvironment{
		env:           map[string]string{"TERM_PROGRAM": "iTerm.app", "TERM_PROGRAM_VERSION": "3.0.10"},
		isNotTerminal: true,
	}))
	if result.Level != Ansi16m {
		t.Errorf("Expected %v, got %v", Ansi16m, result.Level)
	}

	// Should be able to force that we are not a terminal even when we detect we are.
	result = SupportsColor(0, IsTTYOption(false), setEnvironment(&testEnvironment{
		env:           map[string]string{"TERM_PROGRAM": "iTerm.app", "TERM_PROGRAM_VERSION": "3.0.10"},
		isNotTerminal: false,
	}))
	if result.Level != None {
		t.Errorf("Expected %v, got %v", None, result.Level)
	}
}
