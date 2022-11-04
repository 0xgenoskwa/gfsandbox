package mqtt

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/0xgenoskwa/gfsandbox/config"
	"github.com/0xgenoskwa/gfsandbox/domain"
	"github.com/0xgenoskwa/gfsandbox/handler"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttMessage struct {
	Payload map[string]string
	Type    string
	Status  string
	Topic   string
}

type MqttClient struct {
	Port     int
	Broker   string
	ClientID string
	UserName string
	Password string
	CMessage chan *MqttMessage
	CError   chan error
	Client   mqtt.Client
	Handler  *handler.Handler
	Config   *config.Config
}

func NewMqttClient() MqttClient {
	client := MqttClient{}
	client.Broker = "mqtt.dev.generative.xyz"
	client.Port = 1883
	client.ClientID = "go_mqtt_client"
	client.CError = make(chan error)
	client.CMessage = make(chan *MqttMessage)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", client.Broker, client.Port))
	opts.SetClientID(client.ClientID)

	opts.SetDefaultPublishHandler(client.MessagePubHandler)
	opts.OnConnect = client.ConnectHandler
	opts.OnConnectionLost = client.ConnectLostHandler
	client.Client = mqtt.NewClient(opts)
	return client
}

func (c MqttClient) MessagePubHandler(client mqtt.Client, msg mqtt.Message) {
	b := []byte(msg.Payload())
	cmdType, resp, err := c.Handler.OnData(b)
	if err != nil {
		respErr := domain.ErrorResponse{
			Data:  msg.Payload(),
			Error: err.Error(),
		}
		bytes, _ := json.Marshal(respErr)
		fmt.Println("err data", bytes)
	}
	if resp != nil {
		cmdTypeStr := strconv.Itoa(int(cmdType))
		msg := []byte(cmdTypeStr)
		msg = append(msg, []byte(":")...)
		msg = append(msg, resp...)
		fmt.Println("err data", msg)

	}

}

func (c MqttClient) ConnectHandler(client mqtt.Client) {
	fmt.Println("Connected")
}

func (c MqttClient) ConnectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (c MqttClient) Connect() error {
	token := c.Client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	return nil
}

func (c MqttClient) Sub(topic string) {
	token := c.Client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func ProvideMQTT(c *config.Config, h *handler.Handler) *MqttClient {
	return &MqttClient{
		Config:  c,
		Handler: h,
	}
}
