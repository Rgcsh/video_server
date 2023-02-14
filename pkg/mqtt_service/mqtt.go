// 
// All rights reserved
// create time '2022/12/8 14:49'
//
// Usage:
//

package mqtt_service

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"video_server/conf"
	"video_server/pkg/glog"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
var MqttClient *mqtt.Client

//
//  @Description: mqtt根据参数初始化,并连接到mqtt服务器
//  @param port:
//  @param broker: ip地址
//  @param username:
//  @param password:
//
func MqttSetUp() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", conf.Conf.MqttConf.Host, conf.Conf.MqttConf.Port))
	opts.SetClientID("GoCamera")
	opts.SetUsername(conf.Conf.MqttConf.UserName)
	opts.SetPassword(conf.Conf.MqttConf.Password)
	opts.SetCleanSession(false)
	opts.SetDefaultPublishHandler(messagePubHandler)
	//opts.OnConnect = func(client mqtt.Client) {
	//	glog.Log.Info("连接mqtt服务成功")
	//}
	//opts.OnConnectionLost = func(client mqtt.Client, err error) {
	//	glog.Log.Info("失去mqtt连接")
	//}
	tmpClient := mqtt.NewClient(opts)
	MqttClient = &tmpClient
	if token := (*MqttClient).Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("连接mqtt服务失败,错误为:%v", token.Error()))
	}
	return
}

//
//  @Description: 发布消息到topic
//  @param client:
//  @param topic:
//  @param message:
//
func MqttPublish(MqttMessagesChannel <-chan []string) {
	for {
		messages, ok := <-MqttMessagesChannel
		if !ok {
			glog.Log.Info("MqttMessagesChannel channel关闭,此mqtt不再发布消息")
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

//
//  @Description: 订阅某个topic
//  @param client:
//
func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s \n", topic)
}
