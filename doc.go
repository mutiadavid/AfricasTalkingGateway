//
// Copyright 2017
//By: davidmutia47@gmail.com. 
//All rights reserved.
// Use of this source code is governed by a BSD-style
//

/*
Package AfricasTalkingGateway provides AfricasTalking api implementation in golang

Example of sending message
	package main

	import (
		"github.com/dave254/AfricasTalkingGateway"
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

more information: https://github.com/davidmutia47/AfricasTalkingGateway
*/
package AfricasTalkingGateway
