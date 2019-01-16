package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var brokerClients = make(chan bool)
var mqttOpts = mqtt.NewClientOptions()
var mqttClient = mqtt.NewClient(mqttOpts)

func brokerClientsHandler(client mqtt.Client, msg mqtt.Message) {
	brokerClients <- true
	fmt.Printf("BrokerClientsHandler      ")
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func event(t int) {
	fmt.Println("tick", t)
	mqttClient.Publish("test", 1, false, "testBack")
}
func main() {

	//opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("router-sample")
	//opts.SetCleanSession(true)
	mqttOpts.AddBroker("tcp://localhost:1883")
	mqttOpts.SetClientID("router-sample")

	mqttClient = mqtt.NewClient(mqttOpts)
	//c := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := mqttClient.Subscribe("test", 0, brokerClientsHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	ticker := time.Tick(5 * time.Second)
	tick := 0
	for {
		select {
		case <-brokerClients:
		case <-ticker:
			event(tick)
			tick++
		}
	}
}
