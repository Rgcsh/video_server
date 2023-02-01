// (C) Guangcai Ren <rgc@bvrft.com>
// All rights reserved
// create time '2022/12/8 14:49'
//
// Usage:
//

package mqtt_service

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
var MqttClient *mqtt.Client

func MqttSetUp(port int, broker, username, password string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("GoCamera")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("连接mqtt服务成功")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Println("失去mqtt连接")
	}
	tmpClient := mqtt.NewClient(opts)
	MqttClient = &tmpClient
	if token := (*MqttClient).Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("连接mqtt服务失败,错误为:%v", token.Error()))
	}
	return
}

//
//  @Description: 发布消息
//  @param client:
//  @param topic:
//  @param message:
//
func MqttPublish(MqttMessagesChannel <-chan []string) {
	for {
		messages, ok := <-MqttMessagesChannel
		if !ok {
			fmt.Println("MqttMessagesChannel channel关闭,此mqtt不再发布消息")
			break
		}
		topic := messages[0]
		message := messages[1]
		fmt.Printf("开始发布消息:topic:%v,message:%v\n", topic, message)
		token := (*MqttClient).Publish(topic, 2, true, message)
		if token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("发布消息失败:topic:%v,message:%v,err:%v", topic, message, token.Error()))
		}
		fmt.Printf("发布消息成功:topic:%v,message:%v\n", topic, message)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", topic)
}
