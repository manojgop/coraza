package actions

import "testing"

func TestBlockInit(t *testing.T) {
	t.Run("no arguments", func(t *testing.T) {
		a := block()
		if err := a.Init(nil, ""); err != nil {
			t.Error(err)
		}
	})

	t.Run("unexpected arguments", func(t *testing.T) {
		a := block()
		if err := a.Init(nil, "abc"); err == nil || err != ErrUnexpectedArguments {
			t.Error("expected error ErrUnexpectedArguments")
		}
	})
}
