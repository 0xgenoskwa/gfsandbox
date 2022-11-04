package mqtt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.genframe.xyz/config"
	"go.genframe.xyz/domain"
	"go.genframe.xyz/internal/genframe/usecase"
	"go.genframe.xyz/pkg/chrome"
)

type MqttMessage struct {
	Payload map[string]string
	Type    string
	Status  string
	Topic   string
}

type Mqtt struct {
	Config *config.Config

	Usecase *usecase.Usecase
	Chrome  *chrome.Chrome
	Client  mqtt.Client

	notify chan error
}

func ProvideMQTT(c *config.Config, u *usecase.Usecase, chr *chrome.Chrome) *Mqtt {
	m := Mqtt{}

	fmt.Printf("\n\n%s:%d\n\n", c.MqttUrl, c.MqttPort)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", c.MqttUrl, c.MqttPort))
	opts.SetClientID(c.DeviceName)
	opts.SetDefaultPublishHandler(m.onReceiveMessage)
	opts.OnConnect = m.ConnectHandler
	opts.OnConnectionLost = m.ConnectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	m.Client = client

	return &m
}

func (m *Mqtt) onData(data []byte) (domain.CommandType, []byte, error) {
	parts := bytes.Split(data, []byte(":"))
	firstNumber, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return -1, nil, err
	}
	cmdType := domain.CommandType(firstNumber)
	var msg []byte
	if len(parts) > 1 {
		msg = bytes.Join(parts[1:], []byte(":"))
	}

	switch cmdType {
	case domain.CommandTypeInformation:
		resp, err := m.Usecase.GetInformation()
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSetup:
		resp, err := m.Usecase.Setup(msg)
		if err != nil {
			m.Chrome.Toast(err.Error())
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeOpenUrl:
		resp, err := m.Usecase.OpenUrl(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSendTouchEvent:
		resp, err := m.Usecase.SendTouchEvent(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	default:
		return -1, nil, errors.New("unknown cmd")
	}
}

func (m Mqtt) onReceiveMessage(client mqtt.Client, msg mqtt.Message) {
	b := []byte(msg.Payload())
	cmdType, resp, err := m.onData(b)
	cmdTypeStr := strconv.Itoa(int(cmdType))
	if err != nil {
		respErr := domain.ErrorResponse{
			Data:  msg.Payload(),
			Error: err.Error(),
		}
		bytes, _ := json.Marshal(respErr)
		msg := []byte(cmdTypeStr)
		msg = append(msg, []byte(":")...)
		msg = append(msg, bytes...)
		fmt.Println("err data", string(msg))
	}
	if resp != nil {
		cmdTypeStr := strconv.Itoa(int(cmdType))
		msg := []byte(cmdTypeStr)
		msg = append(msg, []byte(":")...)
		msg = append(msg, resp...)
		fmt.Println("err data", string(msg))
	}
}

func (m *Mqtt) ConnectHandler(client mqtt.Client) {
	fmt.Println("Connected")

	if token := m.Client.Subscribe(m.Config.FdChannel, 0, m.onReceiveMessage); token.Wait() && token.Error() != nil {
		m.notify <- token.Error()
	}
}

func (m *Mqtt) ConnectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (m *Mqtt) Connect() error {
	token := m.Client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		return err
	}
	return nil
}

func (m *Mqtt) Sub(topic string) {
	token := m.Client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

// Notify -.
func (m *Mqtt) Notify() <-chan error {
	return m.notify
}

// Shutdown -.
func (m *Mqtt) Shutdown() error {
	m.Client.Disconnect(1000)
	return nil
}
