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

			c.Set(tt.setKey, tt.value)

			v, ok := c.Get(tt.getKey)
			require.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantValue, v)
		})
	}
}

func TestCache_Set(t *testing.T) {
	t.Run("overwrites_value", func(t *testing.T) {
		c := New(context.Background(), time.Minute, NoopExpire)
		c.Set(1, "value1")
		v, ok := c.Get(1)
		require.True(t, ok)
		assert.Equal(t, "value1", v)

		c.Set(1, "value2")
		v2, ok := c.Get(1)
		require.True(t, ok)
		assert.Equal(t, "value2", v2)
	})
}

func TestCache_Delete(t *testing.T) {
	const (
		setKey = 1
		value  = "test"
	)

	tests := map[string]struct {
		key    interface{}
		wantOk bool
	}{
		"found": {
			key:    setKey,
			wantOk: false,
		},
		"not_found": {
			key:    100500,
			wantOk: false,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := New(ctx, time.Minute, NoopExpire)

	t.Run("prepare", func(t *testing.T) {
		c.Set(setKey, value)
		_, ok := c.Get(setKey)
		require.True(t, ok)
	})

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c.Delete(tt.key)
			_, ok := c.Get(tt.key)
			require.Equal(t, tt.wantOk, ok)
		})
	}

	t.Run("vaults_are_closed", func(t *testing.T) {
		cancel()
		<-c.Done()
	})
}
