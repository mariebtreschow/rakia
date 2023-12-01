package main

import (
	"testing"
)

func TestDecode(t *testing.T) {

	decoder := NewDecoder()

	t.Run("single digit", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("1")
		if err != nil {
			t.Error(err)
		}
		if *result != 1 {
			t.Errorf("expected 1 result, got %d", result)
		}

	})

	t.Run("two digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("12")
		if err != nil {
			t.Error(err)
		}
		if *result != 2 {
			t.Errorf("expected 2 results, got %d", result)
		}
	})

	t.Run("three digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("123")
		if err != nil {
			t.Error(err)
		}
		if *result != 3 {
			t.Errorf("expected 3 results, got %d", result)
		}
	})

	t.Run("four digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("1234")
		if err != nil {
			t.Error(err)
		}
		if *result != 3 {
			t.Errorf("expected 3 results, got %d", result)
		}
	})
}
