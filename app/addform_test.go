package main

import (
  "testing"
  "context"
  "os"
  "log"
  "strings"
  "github.com/chromedp/chromedp"
  // "github.com/chromedp/chromedp/kb"
  // "github.com/RollMan/ac2manager/app/models"
)

type InputValuePair struct {
  Input string
  Value string
}

func TestAddForm(t *testing.T){
  ctx, cancel := chromedp.NewContext(context.Background())
  defer cancel()

  forms := make([]InputValuePair, 0)
  forms = append(forms, InputValuePair{`//input[@name="userid"]`, os.Getenv("AC2_APP_ADMINUSERNAME")}, InputValuePair{`//input[@name="pw"]`, os.Getenv("AC2_APP_ADMINPASSWORD")})

  var res string
  err := chromedp.Run(ctx, send(`http://localhost:8000/login`, forms, &res))
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("/login response: `%s`", strings.TrimSpace(res))

  // err := chromedp.Run(ctx, 

}

func send(urlstr string, forms []InputValuePair, res *string) chromedp.Tasks{
  tasks := make(chromedp.Tasks, 0)
  tasks = append(tasks, chromedp.Navigate(urlstr))
  for _, f := range(forms) {
    tasks = append(tasks, chromedp.WaitVisible(f.Input))
    tasks = append(tasks, chromedp.SendKeys(f.Input, f.Value))
  }
  tasks = append(tasks, chromedp.Submit(forms[0].Input))
  tasks = append(tasks, chromedp.WaitReady(`/html/body`))
  tasks = append(tasks, chromedp.Text(`/html`, res))

  return tasks
}
