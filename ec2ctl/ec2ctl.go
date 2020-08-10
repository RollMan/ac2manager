package main

import (
  "time"
  "fmt"
)

func main(){
  time.Local = time.FixedZone("GMT", 0)
  prev := time.Now()
  fmt.Println(prev)
  for {
    now := time.Now()
    sleep_by := now.Add(time.Minute * time.Duration(1))
    sleepUntilNextMinute(sleep_by)


    prev = now
  }
}

func sleepUntilNextMinute(target time.Time){
}

func roundDownIntoMinute(t time.Time) time.Time{
  tSecond := t.Second()
  subtracted := t.Add(time.Second * time.Duration(-tSecond))
  rounded := subtracted.Round(time.Minute)
  return rounded
}
