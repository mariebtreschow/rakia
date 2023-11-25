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
		if len(result) != 1 {
			t.Errorf("expected 1 result, got %d", len(result))
		}
		if result[0] != "a" {
			t.Errorf("expected a, got %s", result[0])
		}
	})

	t.Run("two digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("12")
		if err != nil {
			t.Error(err)
		}
		if len(result) != 2 {
			t.Errorf("expected 2 results, got %d", len(result))
		}
		if result[0] != "ab" {
			t.Errorf("expected ab, got %s", result[0])
		}
		if result[1] != "l" {
			t.Errorf("expected l, got %s", result[1])
		}
	})

	t.Run("three digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("123")
		if err != nil {
			t.Error(err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 results, got %d", len(result))
		}
		if result[0] != "abc" {
			t.Errorf("expected abc, got %s", result[0])
		}
		if result[1] != "aw" {
			t.Errorf("expected aw, got %s", result[1])
		}
		if result[2] != "lc" {
			t.Errorf("expected lc, got %s", result[2])
		}
	})

	t.Run("four digits", func(t *testing.T) {
		result, err := decoder.FindAllCombinations("1234")
		if err != nil {
			t.Error(err)
		}
		if len(result) != 3 {
			t.Errorf("expected 3 results, got %d", len(result))
		}
		if result[0] != "abcd" {
			t.Errorf("expected abcd, got %s", result[0])
		}
		if result[1] != "awd" {
			t.Errorf("expected awd, got %s", result[1])
		}

		if result[2] != "lcd" {
			t.Errorf("expected lcd, got %s", result[2])
		}
	})

}
