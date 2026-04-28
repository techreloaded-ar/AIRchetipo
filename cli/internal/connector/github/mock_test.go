package github

import (
	"context"
	"strings"
	"testing"
)

// mockRunner is a minimal recordable Runner. Tests pre-program responses
// keyed by a substring match on the full args string and capture every call
// for later assertions.
type mockRunner struct {
	t        *testing.T
	calls    []call
	matchers []responder
}

type call struct {
	args  []string
	stdin []byte
}

type responder struct {
	matchPrefix string // matches when the joined args start with this
	stdout      string
	stderr      string
	err         error
}

func newMock(t *testing.T) *mockRunner {
	return &mockRunner{t: t}
}

func (m *mockRunner) on(prefix string, stdout string) *mockRunner {
	m.matchers = append(m.matchers, responder{matchPrefix: prefix, stdout: stdout})
	return m
}

func (m *mockRunner) Run(_ context.Context, stdin []byte, args ...string) ([]byte, []byte, error) {
	m.calls = append(m.calls, call{args: append([]string(nil), args...), stdin: append([]byte(nil), stdin...)})
	full := strings.Join(args, " ")
	for _, r := range m.matchers {
		if strings.HasPrefix(full, r.matchPrefix) {
			return []byte(r.stdout), []byte(r.stderr), r.err
		}
	}
	m.t.Helper()
	m.t.Logf("unexpected gh call: %v", args)
	return nil, []byte("no matcher"), nil
}

func (m *mockRunner) calledWithPrefix(prefix string) bool {
	for _, c := range m.calls {
		if strings.HasPrefix(strings.Join(c.args, " "), prefix) {
			return true
		}
	}
	return false
}
