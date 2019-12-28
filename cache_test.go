package wcache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_Get(t *testing.T) {
	tests := map[string]struct {
		setKey    interface{}
		getKey    interface{}
		value     interface{}
		wantValue interface{}
		wantOk    bool
	}{
		"found": {
			setKey:    1,
			getKey:    1,
			value:     "test",
			wantValue: "test",
			wantOk:    true,
		},
		"not_found": {
			setKey:    1,
			getKey:    100500,
			value:     "test",
			wantValue: nil,
			wantOk:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := New(context.Background(), time.Minute, NoopExpire)

			err := c.Set(tt.setKey, tt.value)
			require.NoError(t, err)

			v, ok := c.Get(tt.getKey)
			require.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}
