package memcache

import (
	"fmt"
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	for _, tt := range []struct {
		name       string
		cache      *cache
		key        string
		val        interface{}
		validUntil time.Time
	}{
		{
			name:  "test empty key",
			cache: newMemCache(),
			key:   "",
		},
		{
			name:  "test non-empty key",
			cache: newMemCache(),
			key:   "foo",
		},
		{
			name:  "test non-empty val",
			cache: newMemCache(),
			key:   "foo",
			val:   "bar",
		},
		{
			name:       "test non-empty validity lifetime",
			cache:      newMemCache(),
			key:        "foo",
			val:        "bar",
			validUntil: time.Unix(0, 1),
		},
		{
			name: "test overwrite existing value",
			cache: &cache{
				m: map[string]*expiringVal{
					"foo": {
						v:          "bar",
						validUntil: time.Unix(0, 0),
					},
				},
			},
			key:        "foo",
			val:        "baz",
			validUntil: time.Unix(0, 1),
		},
	} {
		tt.cache.Put(tt.key, tt.val, tt.validUntil)
		val, ok := tt.cache.m[tt.key]
		if !ok {
			t.Errorf(fmt.Sprintf("expected value in cache not found: %s", tt.val))
		}
		if val.v != tt.val {
			t.Errorf(fmt.Sprintf("expected value in cache did not match actual value: %s, %s", tt.val, val.v))
		}
		if val.validUntil != tt.validUntil {
			t.Errorf(fmt.Sprintf("expected validity lifetime in cache did not match actual validity lifetime: %s, %s", tt.validUntil, val.validUntil))
		}
	}
}

func TestGet(t *testing.T) {
	type expected struct {
		val   interface{}
		valid bool
		err   string
	}
	for _, tt := range []struct {
		name     string
		cache    *cache
		key      string
		expected expected
	}{
		{
			name:  "test key does not exist",
			key:   "foo",
			cache: newMemCache(),
			expected: expected{
				err: "key does not exist in cache: foo",
			},
		},
		{
			name: "test value within validity lifetime",
			key:  "foo",
			cache: &cache{
				m: map[string]*expiringVal{
					"foo": {
						v:          "bar",
						validUntil: time.Now().Add(time.Minute),
					},
				},
			},
			expected: expected{
				val:   "bar",
				valid: true,
			},
		},
		{
			name: "test value outside validity lifetime",
			key:  "foo",
			cache: &cache{
				m: map[string]*expiringVal{
					"foo": {
						v:          "bar",
						validUntil: time.Now().Add(-time.Minute),
					},
				},
			},
			expected: expected{
				valid: false,
			},
		},
	} {
		val, valid, err := tt.cache.Get(tt.key)
		if tt.expected.err != "" && err == nil {
			t.Errorf(fmt.Sprintf("expected error but did not get one"))
		}
		if err != nil && tt.expected.err != err.Error() {
			t.Errorf(fmt.Sprintf("expected error did not match actual error: %s, %s", tt.expected.err, err.Error()))
		}
		if tt.expected.valid != valid {
			t.Errorf(fmt.Sprintf("expected validity did not match actual validity: %t, %t", tt.expected.valid, valid))
		}
		if tt.expected.valid && tt.expected.val != val {
			t.Errorf(fmt.Sprintf("expected value did not match actual value: %s, %s", tt.expected.val, val))
		}
	}
}

func newMemCache() *cache {
	return &cache{
		m: make(map[string]*expiringVal),
	}
}
