## Golang AfricasTalking gateway implementation

More info [GoDoc](http://godoc.org/github.com/davidmutia47/AfricasTalkingGateway)

##Quick Start
######Download and install

    go get github.com/davidmutia47/AfricasTalkingGateway

##Example Sending SMS
######Create file `sendSms.go`
```go
package main

import (
	"github.com/davidmutia47/AfricasTalkingGateway"
	"fmt"
	)

func main(){
    // Specify your login credentials
    username := "username"
    apiKey := "apikey"

    // Specify the numbers that you want to send to in a comma-separated list
    // Please ensure you include the country code (+254 for Kenya in this case)
    recipients := "+25470+++++++,+25475++++++";
    // And of course we want our recipients to know what we really do
    message := "Hello, world";

    //Create instance of getWay
    getWay := AfricasTalkingGateway.AfricasTalkingGateway(username,apikey)

    //call sendMessage to handle sending the message
    response,err := getWay.sendMessage(recipients,message)

    //handle errors if encountered an error
    if err:=nil{
    	//handle error
    }

    for _,receipient :=range response{
    	//get receipient data
    	//type assert to ensure receipient is a map
    	r:= receipient.(map[string]interface{})
    	fmt.Println("number :",r[number],"status :",r["status"])
    }


}
```


## Features

*SMS [Documentation](http://docs.africastalking.com/sms)
######Sending sms
######Fetching sms 
*Sbscriptions [Documentation](http://docs.africastalking.com/subscriptions)
#####Creating subscription
#####Deleting subscription
#####Fetching  subscriptions 
*Voice [Documentation](http://docs.africastalking.com/voice)
#####Making a call
#####Fetching queued calls
#####Uploading media file
*Airtime - buying airtime for a number [Documentation](http://docs.africastalking.com/airtime)
*Mpesa payments [Documentation](http://docs.africastalking.com/payments)
#####B2C
#####B2B
#####C2B

## Documentation

*AfricasTalking Doc [Documentation](http://docs.africastalking.com)
*Implementation Doc [GoDoc](http://godoc.org/github.com/davidmutia47/AfricasTalkingGateway)

## LICENSE

licensed under BSD-style license

