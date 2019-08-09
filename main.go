package main

import(
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"mine/jwt"
	"os"
)

var jwt_token string

type Configuration struct {
    Nexmo_url string `json:"nexmo_url"`
	Whatsapp_from string `json:"whatsapp_from"`
	Priv_key_path string `json:"priv_key_path"`
	Nexmo_app_id string `json:"nexmo_app_id"`
}

var configuration Configuration

type nexmo_message_request struct {
	To struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"to"`
	From struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"from"`
	Message struct{
		Content struct{
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"message"`

}

type nexmo_status_message struct {
	Message_uuid string `json:"message_uuid"`
	To struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"to"`
	From struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"from"`
	Timestamp string `json:"timestamp"`
	Status string `json:"status"`
}

type nexmo_inbound_message struct {
	Message_uuid string `json:"message_uuid"`
	To struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"to"`
	From struct{
		Number string `json:"number"`
		Type string `json:"type"`
	} `json:"from"`
	Timestamp string `json:"timestamp"`
	Direction string `json:"direction"`
	Message struct{
		Content struct{
			Text string `json:"text"`
			Type string `json:"type"`
		} `json:"content"`
	} `json:"message"`

}

func load_conf(){
    if _, err := os.Stat("conf.json"); err != nil {
        if os.IsNotExist(err) {
			fmt.Println(err)
            panic("conf.json doesn't exists")
        }
	}
	
	conf, _ := ioutil.ReadFile("conf.json")
	err := json.Unmarshal(conf, &configuration)
	
	if err != nil {
		fmt.Println(err)
		panic("config error")
	}
}

func jwt_gen(){
	jwt.Set_priv_key_location(configuration.Priv_key_path)
	jwt.Set_app_id(configuration.Nexmo_app_id)
	jwt_token = jwt.Gen()
}

func main(){

	load_conf()

	jwt_gen()

	http.HandleFunc("/get_jwt", get_jwt)
	http.HandleFunc("/send", send_message)
	http.HandleFunc("/status", receive_status)
	http.HandleFunc("/inbound", handle_inbound)
	http.ListenAndServe(":8008", nil)
}

func get_jwt(w http.ResponseWriter, r *http.Request){
	fmt.Println(jwt_token)
}

func send_message(w http.ResponseWriter, r *http.Request){
	fmt.Println(time.Now().Format("2 Jan 2006 15:04:05"))
	fmt.Printf("IP: %s\n", r.RemoteAddr)

	var request nexmo_message_request

	message, ok_message := r.URL.Query()["message"]
	to, ok_to := r.URL.Query()["to"]
    
    if !ok_message ||len(message[0]) < 1 {
        fmt.Fprintf(w, "Url Param 'message' is missing")
        return
    }
    
    if !ok_to || len(to[0]) < 1 {
        fmt.Fprintf(w, "Url Param 'to' is missing")
        return
    }

	request.To.Type = "whatsapp"
	request.To.Number = string(to[0])
	request.From.Type = "whatsapp"
	request.From.Number = configuration.Whatsapp_from
	request.Message.Content.Type = "text"
	request.Message.Content.Text = string(message[0])

	req_json, err_json := json.Marshal(request)

	if err_json != nil {
		fmt.Println("Error marshal json")
		panic(err_json)
	}

	fmt.Printf("Outgoing Message Request\n")
	fmt.Printf("To: %s\n", request.To.Number)
	fmt.Printf("Text: %s\n", request.Message.Content.Text)

	//send request
	req, _ := http.NewRequest("POST", configuration.Nexmo_url, bytes.NewBuffer(req_json))
	req.Header.Add("Authorization", "Bearer "+jwt_token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err_body := ioutil.ReadAll(res.Body)

	if err_body != nil {
		fmt.Println("Error sending msg")
		panic(err_json)
	}

	fmt.Println(string(body))
	fmt.Println("")
	

}

func receive_status(w http.ResponseWriter, r *http.Request){
	fmt.Println(time.Now().Format("2 Jan 2006 15:04:05"))
	fmt.Printf("IP: %s\n", r.RemoteAddr)

	var m nexmo_status_message
	body_json, err_body := ioutil.ReadAll(r.Body)

	if err_body != nil {
		fmt.Println("Error read body")
		panic(err_body)
	}

	err_json := json.Unmarshal(body_json, &m)

	if err_json != nil {
		fmt.Println("Error parse pson")
		fmt.Printf("Body: %s\n\n", body_json)
		panic(err_json)
	}

	fmt.Printf("UUID %s: %s\n", m.Message_uuid, m.Status)
	fmt.Fprintf(w, "status")
}

func handle_inbound(w http.ResponseWriter, r *http.Request){
	fmt.Println(time.Now().Format("2 Jan 2006 15:04:05"))
	fmt.Printf("IP: %s\n", r.RemoteAddr)
	
	var m nexmo_inbound_message
	body_json, err_body := ioutil.ReadAll(r.Body)

	if err_body != nil {
		fmt.Println("Error read body")
		panic(err_body)
	}

	err_json := json.Unmarshal(body_json, &m)

	if err_json != nil {
		fmt.Println("Error parse json")
		fmt.Printf("Body: %s\n\n", body_json)
		panic(err_json)
	}

	fmt.Printf("Incoming Inbound Message\n")
	fmt.Printf("UUID: %s\n", m.Message_uuid)
	fmt.Printf("From: %s\n", m.From.Number)
	fmt.Printf("Text: %s\n\n", m.Message.Content.Text)
	fmt.Fprintf(w, "inbound")
}