package validator

import "testing"

func TestValid(t *testing.T) {
	v := New()
	t.Run("should be valid", func(t *testing.T) {
		if !v.Valid() {
			t.Fatalf("got invalid validator without error, expected to be valid")
		}
	})

	t.Run("should be invalid", func(t *testing.T) {
		v.AddError("test", "this is an error")

		if v.Valid() {
			t.Fatalf("got valid validator with and error, expected to be invalid")
		}
	})
}

func TestCheck(t *testing.T) {
	t.Run("should append an error", func(t *testing.T) {
		v := New()
		v.Check(false, "test", "invalid")

		if _, ok := v.Errors["test"]; !ok {
			t.Fatal("got no error in the list, expected an error")
		}
	})

	t.Run("should not append an error", func(t *testing.T) {
		v := New()
		v.Check(true, "test", "valid")

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

func TestEmailRX(t *testing.T) {
	badstr := "bad-email"
	goodstr := "test@test.com"

	t.Run("shouldn't match", func(t *testing.T) {
		if EmailRX.MatchString(badstr) {
			t.Fatalf("regex shouldn't match: %s", badstr)
		}
	})

	t.Run("should match", func(t *testing.T) {
		if !EmailRX.MatchString(goodstr) {
			t.Fatalf("regex should match: %s", badstr)
		}
	})
}

func TestSpecialCharRX(t *testing.T) {
	badstr := "not special"
	goodstr := `.!@#$%^&:;<>,./\?()[]{}*~-_+=`
	le := len(goodstr)

	t.Run("shouldn't match one", func(t *testing.T) {
		if SpecialCharRX(1, 1).MatchString(badstr) {
			t.Fatalf("regex shouldn't match: %s", badstr)
		}
	})

	t.Run("shouldn't match min", func(t *testing.T) {
		if SpecialCharRX(le+1, 0).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match min of %d: %s", le+1, goodstr)
		}
	})

	t.Run("shouldn't match max", func(t *testing.T) {
		if SpecialCharRX(0, le-1).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match max of %d: %s", le-1, goodstr)
		}
	})

	t.Run("should match", func(t *testing.T) {
		if !SpecialCharRX(le, le).MatchString(goodstr) {
			t.Fatalf("regex should match: %s", goodstr)
		}
	})
}

func TestDigitRX(t *testing.T) {
	badstr := "not digit"
	goodstr := "12345"
	le := len(goodstr)

	t.Run("shouldn't match one", func(t *testing.T) {
		if DigitRX(1, 1).MatchString(badstr) {
			t.Fatalf("regex shouldn't match: %s", badstr)
		}
	})

	t.Run("shouldn't match min", func(t *testing.T) {
		if DigitRX(le+1, 0).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match min of %d: %s", le+1, goodstr)
		}
	})

	t.Run("shouldn't match max", func(t *testing.T) {
		if DigitRX(0, le-1).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match max of %d: %s", le-1, goodstr)
		}
	})

	t.Run("should match", func(t *testing.T) {
		if !DigitRX(le, le).MatchString(goodstr) {
			t.Fatalf("regex should match: %s", goodstr)
		}
	})
}

func TestLowercaseRX(t *testing.T) {
	badstr := "NOT LOWER"
	goodstr := "lower"
	le := len(goodstr)

	t.Run("shouldn't match one", func(t *testing.T) {
		if LowercaseRX(1, 1).MatchString(badstr) {
			t.Fatalf("regex shouldn't match: %s", badstr)
		}
	})

	t.Run("shouldn't match min", func(t *testing.T) {
		if LowercaseRX(le+1, 0).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match min of %d: %s", le+1, goodstr)
		}
	})

	t.Run("shouldn't match max", func(t *testing.T) {
		if LowercaseRX(0, le-1).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match max of %d: %s", le-1, goodstr)
		}
	})

	t.Run("should match", func(t *testing.T) {
		if !LowercaseRX(le, le).MatchString(goodstr) {
			t.Fatalf("regex should match: %s", goodstr)
		}
	})
}

func TestUppercaseRX(t *testing.T) {
	badstr := "not upper"
	goodstr := "UPPER"
	le := len(goodstr)

	t.Run("shouldn't match one", func(t *testing.T) {
		if UppercaseRX(1, 1).MatchString(badstr) {
			t.Fatalf("regex shouldn't match: %s", badstr)
		}
	})

	t.Run("shouldn't match min", func(t *testing.T) {
		if UppercaseRX(le+1, 0).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match min of %d: %s", le+1, goodstr)
		}
	})

	t.Run("shouldn't match max", func(t *testing.T) {
		if UppercaseRX(0, le-1).MatchString(goodstr) {
			t.Fatalf("regex shouldn't match max of %d: %s", le-1, goodstr)
		}
	})

	t.Run("should match", func(t *testing.T) {
		if !UppercaseRX(le, le).MatchString(goodstr) {
			t.Fatalf("regex should match: %s", goodstr)
		}
	})
}
