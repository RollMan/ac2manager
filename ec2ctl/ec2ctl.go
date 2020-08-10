package main

import (
  "time"
  "log"
)

func main(){
  log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
  time.Local = time.FixedZone("GMT", 0)
  prev := time.Now()
  for {
    now := time.Now()
    sleep_by := now.Add(time.Minute * time.Duration(1)).Truncate(time.Minute)
    sleepUntilNextMinute(sleep_by)

    prev = now
  }
}

func sleepUntilNextMinute(target time.Time){
  t1 := time.Now()
  toWait := target.Sub(t1) + time.Second
  for toWait > 0 {
    time.Sleep(toWait)
    t2 := time.Now()
    toWait = toWait - (t2.Sub(t1))
    t1 = t2
  }
}
