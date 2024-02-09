package pathresolver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidPath(t *testing.T) {
	pw, err := newPathWalker("field[1].[0].nested[3][4].value")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("node=field", func(t *testing.T) {
		assert.Equal(t, false, pw.IsCurrentNodeIndex())
		assert.Equal(t, "field", pw.StringValue())
		assert.Equal(t, -1, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1]", func(t *testing.T) {
		assert.Equal(t, true, pw.IsCurrentNodeIndex())
		assert.Equal(t, "", pw.StringValue())
		assert.Equal(t, 1, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1].[0]", func(t *testing.T) {
		assert.Equal(t, true, pw.IsCurrentNodeIndex())
		assert.Equal(t, "", pw.StringValue())
		assert.Equal(t, 0, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1].[0].nested", func(t *testing.T) {
		assert.Equal(t, false, pw.IsCurrentNodeIndex())
		assert.Equal(t, "nested", pw.StringValue())
		assert.Equal(t, -1, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1].[0].nested[3]", func(t *testing.T) {
		assert.Equal(t, true, pw.IsCurrentNodeIndex())
		assert.Equal(t, "", pw.StringValue())
		assert.Equal(t, 3, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1].[0].nested[3][4]", func(t *testing.T) {
		assert.Equal(t, true, pw.IsCurrentNodeIndex())
		assert.Equal(t, "", pw.StringValue())
		assert.Equal(t, 4, pw.IndexValue())
		assert.Equal(t, true, pw.MoveToNextNode())
	})

	t.Run("node=field[1].[0].nested[3][4].value", func(t *testing.T) {
		assert.Equal(t, false, pw.IsCurrentNodeIndex())
		assert.Equal(t, "value", pw.StringValue())
		assert.Equal(t, -1, pw.IndexValue())
		assert.Equal(t, false, pw.MoveToNextNode())
	})
}

func Test_InvalidPaths(t *testing.T) {
	t.Run("path=test[1].value[[]4]]", func(t *testing.T) {
		pw, err := newPathWalker("test[1].value[[]4]]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1].value[[]4]]\". invalid part \"value[[]4]]\"")
	})

	t.Run("path=test[1]a.value[4]", func(t *testing.T) {
		pw, err := newPathWalker("test[1]a.value[4]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1]a.value[4]\". invalid part \"test[1]a\"")
	})

	t.Run("path=test[1].value[abc]", func(t *testing.T) {
		pw, err := newPathWalker("test[1].value[abc]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1].value[abc]\". invalid part \"value[abc]\"")
	})

	t.Run("path=test[1].value[]", func(t *testing.T) {
		pw, err := newPathWalker("test[1].value[]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1].value[]\". invalid part \"value[]\"")
	})

	t.Run("path=test[1].value][", func(t *testing.T) {
		pw, err := newPathWalker("test[1].value][")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1].value][\". invalid part \"value][\"")
	})

	t.Run("path=test[1].value[[]]", func(t *testing.T) {
		pw, err := newPathWalker("test[1].value[[]]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"test[1].value[[]]\". invalid part \"value[[]]\"")
	})

	t.Run("path=\"\"", func(t *testing.T) {
		pw, err := newPathWalker("")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "path can't be empty")
	})

	t.Run("path=\" \"", func(t *testing.T) {
		pw, err := newPathWalker(" ")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "path can't be empty")
	})

	t.Run("path= ] [", func(t *testing.T) {
		pw, err := newPathWalker(" ] [")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \" ] [\". invalid part \"] [\"")
	})

	t.Run("path=field.subField[-2]", func(t *testing.T) {
		pw, err := newPathWalker("field.subField[-2]")
		assert.Nil(t, pw)
		assert.ErrorContains(t, err, "invalid path. path \"field.subField[-2]\". invalid part \"subField[-2]\"")
	})
}
