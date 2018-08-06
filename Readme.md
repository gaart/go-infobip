# Infobip API Go client

[![Go Report Card](https://goreportcard.com/badge/github.com/gaart/go-infobip)](https://goreportcard.com/report/github.com/gaart/go-infobip)
[![license](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go client library for interacting with Infobip's API.

## Install

```sh
go get github.com/gaart/go-infobip
```

## Usage

To send SMS:

```go
import "github.com/gaart/go-infobip"
```

```go
username := os.Getenv("INFOBIP_USERNAME")
password := os.Getenv("INFOBIP_PASSWORD")

client, err := infobip.NewClient(username, password)
if err != nil {
    fmt.Println(err.Error())
    return
}

sms := infobip.SMS{
    From: "sender name",
    To:   []string{"+12125557890"},  // list of phone numbers
    Text: "message text here",
}

res, err := client.SendSMS(&sms)
if err != nil {
    fmt.Println(err.Error())
    return
}

fmt.Printf("%+v\n", res)

```

To get SMS delivery report:

```go
messageId := "123123123"  // SMS ID to check

report, err := client.GetDeliveryReport(messageId)
if err != nil {
    fmt.Println(err.Error())
}
fmt.Printf("%+v\n", report)
```
