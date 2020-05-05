package main

import (
  "time"
  "testing"
  "context"
  "os"
  "log"
  "strconv"
  "github.com/chromedp/chromedp"
  "github.com/RollMan/ac2manager/app/models"
  "github.com/chromedp/cdproto/network"
  // "github.com/chromedp/chromedp/kb"
)

type InputValuePair struct {
  Input string
  Value interface{}
}

func TestAddForm(t *testing.T){
  var optsA = [...]chromedp.ExecAllocatorOption{
    chromedp.NoFirstRun,
    chromedp.NoDefaultBrowserCheck,
    chromedp.DisableGPU,

    chromedp.Flag("disable-background-networking", true),
    chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
    chromedp.Flag("disable-background-timer-throttling", true),
    chromedp.Flag("disable-backgrounding-occluded-windows", true),
    chromedp.Flag("disable-breakpad", true),
    chromedp.Flag("disable-client-side-phishing-detection", true),
    chromedp.Flag("disable-default-apps", true),
    chromedp.Flag("disable-dev-shm-usage", true),
    chromedp.Flag("disable-extensions", true),
    chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
    chromedp.Flag("disable-hang-monitor", true),
    chromedp.Flag("disable-ipc-flooding-protection", true),
    chromedp.Flag("disable-popup-blocking", true),
    chromedp.Flag("disable-prompt-on-repost", true),
    chromedp.Flag("disable-renderer-backgrounding", true),
    chromedp.Flag("disable-sync", true),
    chromedp.Flag("force-color-profile", "srgb"),
    chromedp.Flag("metrics-recording-only", true),
    chromedp.Flag("safebrowsing-disable-auto-update", true),
    chromedp.Flag("enable-automation", true),
    chromedp.Flag("password-store", "basic"),
    chromedp.Flag("use-mock-keychain", true),
}
  opts := optsA[:]
  allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
  defer cancel()
  ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
  defer cancel()

  // prepare for http response retrieve
  var statuscode *int64 = new(int64)
  err := chromedp.Run(ctx, network.Enable())
  chromedp.ListenTarget(ctx, func(event interface{}) {
    switch responseReceivedEvent := event.(type) {
    case *network.EventResponseReceived:
      response := responseReceivedEvent.Response
      *statuscode = response.Status
    }
  })


  // login
  forms_login := make([]InputValuePair, 0)
  forms_login = append(forms_login, InputValuePair{`//input[@name="userid"]`, os.Getenv("AC2_APP_ADMINUSERNAME")}, InputValuePair{`//input[@name="pw"]`, os.Getenv("AC2_APP_ADMINPASSWORD")})

  var res string
  err = chromedp.Run(ctx, send(`http://localhost:8000/login`, forms_login, &res))
  if err != nil {
    log.Fatal(err)
  }
  if *statuscode != 200 {
    log.Fatalf("Bad status code when login: %d", *statuscode)
  }

  // add
  event := models.Event{
    Startdate: time.Now(),
    Track: `monza`,
    WeatherRandomness                    : 2,
    P_hourOfDay                          : 14,
    P_timeMultiplier                     : 1,
    P_sessionDurationMinute              : 30,
    Q_hourOfDay                          : 14,
    Q_timeMultiplier                     : 1,
    Q_sessionDurationMinute              : 30,
    R_hourOfDay                          : 14,
    R_timeMultiplier                     : 1,
    R_sessionDurationMinute              : 30,
    PitWindowLengthSec                   : 600,
    IsRefuellingAllowedInRace            : true,
    MandatoryPitstopCount                : 1,
    IsMandatoryPitstopRefuellingRequired : false,
    IsMandatoryPitstopTyreChangeRequired : true,
    IsMandatoryPitstopSwapDriverRequired : true,
    TyreSetCount                         : 3,
  }

  startdate_day := event.Startdate.Format("2006-01-02")
  startdate_time := event.Startdate.Format("15:04")


  forms_add := make([]InputValuePair, 0)
  forms_add = append(forms_add,
  InputValuePair{`//input[@name="tyreSetCount"]`, event.TyreSetCount},
  InputValuePair{`//input[@name="startdate_day"`, startdate_day},
  InputValuePair{`//input[@name="startdate_time"`, startdate_time},
  InputValuePair{`//input[@name="track"]`, event.Track},
  InputValuePair{`//input[@name="weatherRandomness"]`, event.WeatherRandomness},
  InputValuePair{`//input[@name="P_hourOfDay"]`, event.P_hourOfDay},
  InputValuePair{`//input[@name="P_timeMultiplier"]`, event.P_timeMultiplier},
  InputValuePair{`//input[@name="P_sessionDurationMinute"]`, event.P_sessionDurationMinute},
  InputValuePair{`//input[@name="Q_hourOfDay"]`, event.Q_hourOfDay},
  InputValuePair{`//input[@name="Q_timeMultiplier"]`, event.Q_timeMultiplier},
  InputValuePair{`//input[@name="Q_sessionDurationMinute"]`, event.Q_sessionDurationMinute},
  InputValuePair{`//input[@name="R_hourOfDay"]`, event.R_hourOfDay},
  InputValuePair{`//input[@name="R_timeMultiplier"]`, event.R_timeMultiplier},
  InputValuePair{`//input[@name="R_sessionDurationMinute"]`, event.R_sessionDurationMinute},
  InputValuePair{`//input[@name="pitWindowLengthSec"]`, event.PitWindowLengthSec},
  InputValuePair{`//input[@name="isRefuellingAllowedInRace"]`, event.IsRefuellingAllowedInRace},
  InputValuePair{`//input[@name="mandatoryPitstopCount"]`, event.MandatoryPitstopCount},
  InputValuePair{`//input[@name="isMandatoryPitstopRefuellingRequired"]`, event.IsMandatoryPitstopRefuellingRequired},
  InputValuePair{`//input[@name="isMandatoryPitstopTyreChangeRequired"]`, event.IsMandatoryPitstopTyreChangeRequired},
  InputValuePair{`//input[@name="isMandatoryPitstopSwapDriverRequired"]`, event.IsMandatoryPitstopSwapDriverRequired},
)

  // Have to parse startdate and time.
  err = chromedp.Run(ctx, send(`http://localhost:8000/add_event`, forms_add, &res))

  log.Printf(res)

  if err != nil {
    log.Fatal(err)
  }

  if *statuscode != 200 {
    log.Fatalf("Bad status code when login: %d", *statuscode)
  }
}


func send(urlstr string, forms []InputValuePair, res *string) chromedp.Tasks{
  tasks := make(chromedp.Tasks, 0)
  tasks = append(tasks, chromedp.Navigate(urlstr))
  for _, f := range(forms) {
    tasks = append(tasks, chromedp.WaitVisible(f.Input))
    var value string
    switch t := f.Value.(type) {
    case int:
      value = strconv.Itoa(t)
    case bool:
      value = strconv.FormatBool(t)
    case string:
      value = t
    }
    tasks = append(tasks, chromedp.SendKeys(f.Input, value))
  }
  tasks = append(tasks, chromedp.Submit(forms[0].Input))
  tasks = append(tasks, chromedp.WaitReady(`/html/body`))
  tasks = append(tasks, chromedp.Text(`/html`, res))

  return tasks
}
