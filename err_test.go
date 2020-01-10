package std

import (
	"errors"
	"testing"
)

func TestErrorf(t *testing.T) {
	t.Log(Errorf("sample %s %s", "1", "2"))
}

func TestErrorWrapf(t *testing.T) {
	t.Log(ErrorWrapf(errors.New("raw error"), "sample %d/%d", 1, 2).SetData(1))
}

func TestErrorWrap(t *testing.T) {
	t.Log(ErrorWrap(errors.New("raw error"), "what fuck?"))
}

func TestError(t *testing.T) {
	t.Log(Error("what fuck?"))
}
