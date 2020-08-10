package main

import (
  "testing"
  "time"
)

func TestRoundDownIntoMinute(test *testing.T){
  t := time.Date(2020, time.August, 10, 8, 10, 0, 0, time.UTC)
  for d := 0; d < 60; d++ {
    s := t.Add(time.Duration(d) * time.Second)
    rounded := roundDownIntoMinute(s)
    if !rounded.Equal(t) {
      test.Fatalf("Test failed: t=%v, s=%v, rounded=%v", t, s, rounded)
    }
  }

  expected := t.Add(1 * time.Minute)
  for d := 60; d < 120; d++{
    s := t.Add(time.Duration(d) * time.Second)
    rounded := roundDownIntoMinute(s)
    if !rounded.Equal(expected) {
      test.Fatalf("Test failed: t=%v, s=%v, rounded=%v", t, s, rounded)
    }
  }
}
