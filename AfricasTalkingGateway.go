//
// Copyright 2017
//By: davidmutia47@gmail.com. 
//All rights reserved.
// Use of this source code is governed by a BSD-style
//

/*
AfricasTalking Golang gateway implementation
*/

package AfricasTalkingGateway

import(
	"errors"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"strconv"
	"net/url"
	"strings"
	"fmt"
	
)

type AfricasTalking struct{
	UserName string
	ApiKey string
	Environment string
	ResponseCode int
}


func (a AfricasTalking) AfricasTalkingGatewayException(err string) (error) {
	return errors.New(err)
}

func AfricasTalkingGateway(username,apikey string,params ...string) (*AfricasTalking){
	a:=&AfricasTalking{}
	a.ApiKey,a.UserName,a.ResponseCode= apikey,username,0
	if len(params) >0{
		a.Environment=params[0]
	}else{
		a.Environment = "production"
	}
	return a			
}

//Parameter for sending message
type MessageParameters struct{
	From string
	BulkSMSMode int
	Enqueue int
	Keyword string
	LinkId string
	RetryDurationInHours string
}



//Function to handle sending messages
func (a AfricasTalking) SendMessage(to,message string,params ...MessageParameters) ([]interface{},error) {
	var recipients []interface{}
	parameters :=map[string]interface{}{}
	parameters["username"],parameters["to"],parameters["message"] =  a.UserName,to,message
	if len(params)>0{
		if params[0].From != ""{
			parameters["from"] =params[0].From
		}
		parameters["bulkSMSMode"] = params[0].BulkSMSMode		
		parameters["enqueue"] = params[0].Enqueue
		if params[0].LinkId !=""{
			if linkId,err :=strconv.Atoi(params[0].LinkId);err==nil{
				parameters["linkId"] = linkId	
			}
		}
		if params[0].RetryDurationInHours !=""{
			if retryDurationInHours,err :=strconv.Atoi(params[0].RetryDurationInHours);err==nil{
				parameters["retryDurationInHours"] = retryDurationInHours	
			}
		}
	}else{
		//set defults for bulksmsmode and enqueue parameters
		parameters["bulkSMSMode"] = 1
		parameters["enqueue"] = 0
	}
	responseBody,err := a.sendRequest(a.getSmsUrl(),parameters)
	if err != nil{
		return recipients,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode ==201 {
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{

			return recipients,a.AfricasTalkingGatewayException(err.Error())
		}
		data :=jsonData.(map[string]interface{})
		sMSMessageData := data["SMSMessageData"].(map[string]interface{})
		recipients = sMSMessageData["Recipients"].([]interface{})
		if len(recipients)>0{
		 	return recipients,nil
		}
	}
	return recipients,a.AfricasTalkingGatewayException(string(responseBody))

}



func (a AfricasTalking) FetchMessages(lastReceivedId int) ([]interface{},error){
	var messages []interface{}
	//use package url later
	urlString := a.getSmsUrl()+"?username="+a.UserName+"&lastReceivedId="+strconv.Itoa(lastReceivedId)
	responseBody,err := a.sendRequestGet(urlString)
	if err !=nil{
		return messages,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode ==200 {
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return messages,a.AfricasTalkingGatewayException(err.Error())
		}
		data :=jsonData.(map[string]interface{})
		sMSMessageData := data["SMSMessageData"].(map[string]interface{})
		messages = sMSMessageData["Messages"].([]interface{})
		if len(messages)>0{
		 	return messages,nil
		}
	}
	return messages,a.AfricasTalkingGatewayException(string(responseBody))
}

//Subscription methods - creating subscription
func (a AfricasTalking) CreateSubscription(phoneNumber,shortCode,keyword string) (map[string]interface{},error){
	urlString := a.getSmsSubscriptionUrl()+"/create"
	return a.subscription(urlString,phoneNumber,shortCode,keyword)
}

//Subscription methods - deleting subscription
func (a AfricasTalking) DeleteSubscription(phoneNumber,shortCode,keyword string) (map[string]interface{},error){
	urlString := a.getSmsSubscriptionUrl()+"/delete"
	return a.subscription(urlString,phoneNumber,shortCode,keyword)
}

//Subscription methods - fetching  subscriptions
func (a AfricasTalking) FetchPremiumSubscriptions(shortCode,keyword string, lastReceivedId ...int) ([]interface{},error) {
	response :=[]interface{}{}
	if shortCode=="" || keyword == "" {
		return response,a.AfricasTalkingGatewayException("Supply short code and Keyword")
	}
	//use package url later
	urlString :=a.getSmsSubscriptionUrl()+"?username="+a.UserName+"&shortCode="+shortCode+"&keyword="+keyword+"&lastReceivedId="
	if len(lastReceivedId) >0{
		urlString += strconv.Itoa(lastReceivedId[0])
	}else{
		urlString += strconv.Itoa(0)
	}
	responseBody,err := a.sendRequestGet(urlString)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 200{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		data := jsonData.(map[string]interface{})
		subscriptions := data["Subscriptions"].([]interface{})
		if len(subscriptions) >0{
			return subscriptions,nil
		}
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
}


func (a AfricasTalking) subscription(urlString,phoneNumber,shortCode,keyword string) (map[string]interface{},error){
	response :=map[string]interface{}{}
	if phoneNumber =="" || shortCode=="" || keyword == "" {
		return response,a.AfricasTalkingGatewayException("Supply phone number, short code and Keyword")
	}
	parameters := map[string]interface{}{}
	parameters["username"],parameters["phoneNumber"],parameters["shortCode"],parameters["keyword"] = a.UserName,phoneNumber,shortCode,keyword
	responseBody,err := a.sendRequest(urlString,parameters)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		response= jsonData.(map[string]interface{})
			return response,nil
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
}

//Voice methods - calling 
func (a AfricasTalking) Call(from,to string) (map[string]interface{},error) {
	response :=map[string]interface{}{}
	parameters := map[string]interface{}{}
	parameters["from"],parameters["to"],parameters["username"] = from,to,a.UserName
	urlString := a.getVoiceUrl()+"/call"
	responseBody,err := a.sendRequest(urlString,parameters)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 200{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error()+"############# json marshal")
		}
		response= jsonData.(map[string]interface{})
		return response,nil
	}else{
		return response,a.AfricasTalkingGatewayException(string(responseBody)+"########### "+strconv.Itoa(a.ResponseCode))	
	}
	
}

func (a AfricasTalking) GetNumQueuedCalls(phoneNumber string,queueName ...string) (map[string]interface{},error){
	parameters :=map[string]interface{}{}
	response :=map[string]interface{}{}
	parameters["username"],parameters["phoneNumber"] = a.UserName,phoneNumber
	if len(queueName) >0 {
		parameters["queueName"] = queueName[0]	
	}
	urlString := a.getVoiceUrl()+"/queueStatus"
	responseBody,err := a.sendRequest(urlString,parameters)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		response= jsonData.(map[string]interface{})
			return response,nil
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
}

func (a AfricasTalking) UploadMediaFile(url string) (map[string]interface{},error) {
	response := map[string]interface{}{}
	parameters :=map[string]interface{}{}
	parameters["username"],parameters["url"] =a.UserName,url 
	urlString := a.getVoiceUrl()+"/mediaUpload"
	responseBody,err := a.sendRequest(urlString,parameters)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		response= jsonData.(map[string]interface{})
			return response,nil
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
}

//Airtime methods - buy airtime for a number
func (a AfricasTalking) SendAirtime(params []map[string]interface{}) ([]interface{},error) {
	response := []interface{}{}
	parameters := map[string]interface{}{}
	urlString := a.getAirtimeUrl()+"/send"
	parameters["username"] = a.UserName
	recipientsJson,_ := json.Marshal(params)
	parameters["recipients"] = string(recipientsJson)
	responseBody,err := a.sendRequest(urlString,parameters)
	if err!=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201{
		var jsonData interface{}
		err := json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
			data := jsonData.(map[string]interface{})
			response = data["responses"].([]interface{})
			if len(response) >0{
				return response,nil	
			}
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
}

//Payment methods - C2B
func (a AfricasTalking) InitiateMobilePaymentCheckout(productName,phoneNumber,currencyCode string,amount float64,metadata map[string]string) (map[string]string,error){
	parameters := map[string]interface{}{}
	response := map[string]string{}
	parameters["username"],parameters["productName"],parameters["phoneNumber"],parameters["currencyCode"],parameters["amount"],parameters["metadata"] = a.UserName,productName,phoneNumber,currencyCode,amount,metadata
	urlString := a.getMobilePaymentCheckoutUrl()	
	responseBody,err := a.sendJSONRequest(urlString,parameters)
	if err !=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201 {
		var jsonData interface{}
		err = json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		data,ok :=jsonData.(map[string]interface{})
		if !ok {
			return response,a.AfricasTalkingGatewayException("incorrect response data")
		}
		if status,ok := data["status"].(string);ok{
			response["status"] = status
		}
		if description,ok := data["description"].(string);ok{
			response["description"] = description
		}
		if transactionId,ok := data["transactionId"].(string);ok{
			response["transactionId"] = transactionId
		}
		return response,nil	
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
			
}

//Payment methods - B2C
func (a AfricasTalking) MobilePaymentB2CRequest(productName string, recipients []map[string]interface{}) ([]interface{},error){
	parameters := map[string]interface{}{}
	entries := []interface{}{}
	parameters["username"],parameters["productName"],parameters["recipients"] = a.UserName,productName,recipients
	urlString := a.getMobilePaymentB2CUrl()
	responseBody,err := a.sendJSONRequest(urlString,parameters)
	if err !=nil {
		return entries,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode ==201 {
		var jsonData interface{}
		err = json.Unmarshal(responseBody,&jsonData)
		if err != nil {
			return nil,a.AfricasTalkingGatewayException(err.Error())
		}
		data := jsonData.(map[string]interface{})
		entries = data["entries"].([]interface{})
		if len(entries)>0 {
			return entries,nil
		}	
	}
	return entries,a.AfricasTalkingGatewayException(string(responseBody))
}

//Payment methods - B2B
func (a AfricasTalking) MobilePaymentB2BRequest(productName,provider,transferType,currencyCode,destinationChannel string,amount float64,metadata map[string]string,destinationAccount ...string) (map[string]string,error){
	parameters :=map[string]interface{}{}
	response := map[string]string{}
	parameters["username"],parameters["productName"],parameters["provider"],parameters["destinationChannel"],parameters["transferType"],parameters["currencyCode"],parameters["amount"],parameters["metadata"]=a.UserName,productName,provider,destinationChannel,transferType,currencyCode,amount,metadata
	if len(destinationAccount) >0{
		parameters["destinationChannel"] = destinationChannel[0]
	}
	urlString:=a.getMobilePaymentB2BUrl()
	responseBody,err := a.sendJSONRequest(urlString,parameters)
	if err !=nil {
		return response,a.AfricasTalkingGatewayException(err.Error())
	}
	if a.ResponseCode == 201 {
		var jsonData interface{}
		err = json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		data,ok :=jsonData.(map[string]interface{})
		if !ok {
			return response,a.AfricasTalkingGatewayException("incorrect response data")
		}
		if status,ok := data["status"].(string);ok{
			response["status"] = status
		}
		if transactionId,ok := data["transactionId"].(string);ok{
			response["transactionId"] = transactionId
		}
		if transactionFee,ok := data["transactionFee"].(string);ok{
			response["transactionFee"] = transactionFee
		}
		if providerChannel,ok := data["providerChannel"].(string);ok{
			response["providerChannel"] = providerChannel
		}
		return response,nil	
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))
			
}

//User data - getting user data eg balance in account
func (a AfricasTalking) GetUserData() (map[string]interface{},error){
	response:=map[string]interface{}{}
	urlString:= a.getUserDataUrl()+"?username="+a.UserName
	responseBody,err := a.sendRequestGet(urlString)
	if err!=nil {
		
	}
	if a.ResponseCode==200{
		var jsonData interface{}
		err = json.Unmarshal(responseBody,&jsonData)
		if err != nil{
			return response,a.AfricasTalkingGatewayException(err.Error())
		}
		data :=jsonData.(map[string]interface{})
		response = data["UserData"].(map[string]interface{})
		return response,nil
	}
	return response,a.AfricasTalkingGatewayException(string(responseBody))

}


//http post request 
func (a *AfricasTalking) sendRequest(urlString string,parameters map[string]interface{}) ([]byte,error){
	request:= &http.Request{}
	var err error
	//var data []byte
	form := url.Values{}
	if len(parameters)!=0{
		for key,value := range parameters{
			if v,ok := value.(string); ok{
				form.Add(key,v)
			}else{
				if v,ok := value.(int); ok{
					form.Add(key,strconv.Itoa(v))
				}	
			}
		}
	}
	fmt.Println(form.Encode())
	request,err = http.NewRequest("POST",urlString,strings.NewReader(form.Encode()))
	if err !=nil{
		return nil,err
	}
	request.Header.Set("Accept","application/json")
	request.Header.Set("Content-Type","application/x-www-form-urlencoded")
	request.Header.Set("apikey",a.ApiKey)
	client :=&http.Client{}
	response,err := client.Do(request)
	if err !=nil{
		return nil,err
	}
	defer response.Body.Close()
	a.ResponseCode = response.StatusCode 
	body,err := ioutil.ReadAll(response.Body)
	if err !=nil{
		return nil,err
	}
	return body,nil
}
//http get request
func (a *AfricasTalking) sendRequestGet(urlString string) ([]byte,error){
	request,err := http.NewRequest("GET",urlString,nil)
	if err!=nil{
		return nil,err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Set("apikey",a.ApiKey)
	client :=&http.Client{}
	response,err := client.Do(request)
	if err !=nil{
		return nil,err
	}
	defer response.Body.Close()
	a.ResponseCode = response.StatusCode 
	body,err := ioutil.ReadAll(response.Body)
	if err !=nil{
		return nil,err
	}
	return body,nil

}

//json post request
func (a *AfricasTalking) sendJSONRequest(urlString string,parameters map[string]interface{}) ([]byte,error) {
	request:= &http.Request{}
	var err error
	jsonParams,err :=json.Marshal(parameters)
	if err !=nil{
		return nil,err
	}
	request,err = http.NewRequest("POST",urlString,bytes.NewBuffer(jsonParams))
	if err !=nil{
		return nil,err
	}
	request.Header.Set("Content-Type","application/json")
	request.Header.Set("apikey",a.ApiKey)
	client :=&http.Client{}
	response,err := client.Do(request)
	if err !=nil{
		return nil,err
	}
	defer response.Body.Close() 
	a.ResponseCode = response.StatusCode
	body,err := ioutil.ReadAll(response.Body)
	if err !=nil{
		return nil,err
	}
	return body,nil
}

func (a AfricasTalking) getApiHost() string{
	if a.Environment !="sandbox" {
		return "https://api.africastalking.com"
	}
	return "https://api.sandbox.africastalking.com"
}
func (a AfricasTalking) getPaymentHost() string{
	if a.Environment !="sandbox" {
		return "https://payments.africastalking.com" 
	}
	return "https://payments.sandbox.africastalking.com"
	
}

func (a AfricasTalking) getVoiceHost() string{
	if a.Environment != "sandbox"{
		return "https://voice.africastalking.com"
	}
    return "https://voice.sandbox.africastalking.com"         
}
func (a AfricasTalking) getSmsUrl() (string) {
	return a.getApiHost()+"/version1/messaging"
}

func (a AfricasTalking) getMobilePaymentCheckoutUrl() string{
	return a.getPaymentHost()+"/mobile/checkout/request"
}

func (a AfricasTalking) getMobilePaymentB2CUrl() string{
	return a.getPaymentHost() + "/mobile/b2c/request"
}
func (a AfricasTalking) getMobilePaymentB2BUrl() string{
	return a.getPaymentHost() +"/mobile/b2b/request"
}
func (a AfricasTalking) getSmsSubscriptionUrl() string{
	return a.getApiHost() + "/version1/subscription"
}

func (a AfricasTalking) getVoiceUrl() string{
	return a.getVoiceHost()
}

func (a AfricasTalking) getAirtimeUrl() string{
	return a.getApiHost()+"/version1/airtime"
}

func (a AfricasTalking) getUserDataUrl() string{
	return a.getApiHost() + "/version1/user"
}