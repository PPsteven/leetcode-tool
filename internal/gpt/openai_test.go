package gpt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenai(t *testing.T) {
	o := NewOpenai("sk-xxxxxx", "gpt-3.5-turbo")

	t.Run("chat", func(t *testing.T) {
		ret, err := o.Chat("写一个冒泡排序")
		assert.NoError(t, err)

		t.Log(ret)
	})
}
