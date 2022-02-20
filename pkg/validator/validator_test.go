package validator

import "testing"

func TestValid(t *testing.T) {
	v := New()
	v.KeyError = "test"

	t.Run("should be valid", func(t *testing.T) {
		if !v.Valid() {
			t.Fatalf("got invalid validator without error, expected to be valid")
		}
	})

	t.Run("should be invalid", func(t *testing.T) {
		v.AddError("this is an error")

		if v.Valid() {
			t.Fatalf("got valid validator with and error, expected to be invalid")
		}
	})
}

func TestCheck(t *testing.T) {
	t.Run("should append an error", func(t *testing.T) {
		v := New()
		v.KeyError = "test"

		v.Check(false, "invalid")

		if _, ok := v.Errors["test"]; !ok {
			t.Fatal("got no error in the list, expected an error")
		}
	})

	t.Run("should not append an error", func(t *testing.T) {
		v := New()
		v.KeyError = "test"

		v.Check(true, "valid")

		if _, ok := v.Errors["test"]; ok {
			t.Fatal("got an error in the list, expected no error")
		}
	})
}

func TestIn(t *testing.T) {
	t.Run("should detect value in args", func(t *testing.T) {
		if ok := In("test", "arg1", "other", "test"); !ok {
			t.Fatal("got no list value detection, expected to be detected")
		}
	})

	t.Run("should not value in args", func(t *testing.T) {
		if ok := In("test", "arg1", "args2"); ok {
			t.Fatal("got list value detection, expected nothing")
		}
	})
}
