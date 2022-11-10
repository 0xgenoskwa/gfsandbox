package bluetooth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.genframe.xyz/config"
	"go.genframe.xyz/domain"
	"go.genframe.xyz/internal/genframe/usecase"
	"go.genframe.xyz/pkg/chrome"
	tinybl "tinygo.org/x/bluetooth"
)

var (
	serviceUUID = tinybl.ServiceUUIDNordicUART
	rxUUID      = tinybl.CharacteristicUUIDUARTRX
	txUUID      = tinybl.CharacteristicUUIDUARTTX
)

type Bluetooth struct {
	Config *config.Config

	Adapter            *tinybl.Adapter
	Advertisement      *tinybl.Advertisement
	AdvertisementState bool
	AdvertisementUntil time.Time

	RxChar tinybl.Characteristic
	TxChar tinybl.Characteristic

	Usecase *usecase.Usecase
	Chrome  *chrome.Chrome

	notify chan error
}

func ProvideBluetooth(c *config.Config, u *usecase.Usecase, chr *chrome.Chrome) *Bluetooth {
	b := Bluetooth{
		Config:  c,
		Usecase: u,
		Chrome:  chr,
		notify:  make(chan error, 1),
	}

	return &b
}

func (b *Bluetooth) onData(data []byte) (domain.CommandType, []byte, error) {
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
		resp, err := b.Usecase.GetInformation()
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeGetMacAddress:
		resp, err := b.Usecase.GetMacAddress()
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSetup:
		resp, err := b.Usecase.Setup(msg)
		if err != nil {
			b.Chrome.Toast(err.Error())
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeOpenUrl:
		resp, err := b.Usecase.OpenUrl(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSendTouchEvent:
		resp, err := b.Usecase.SendTouchEvent(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	case domain.CommandTypeSendKeyEvent:
		resp, err := b.Usecase.SendKeyEvent(msg)
		if err != nil {
			return cmdType, nil, err
		}
		return cmdType, resp, nil
	default:
		return -1, nil, errors.New("unknown cmd")
	}
}

func (b *Bluetooth) Start() error {
	adapter := tinybl.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}
	b.Adapter = adapter
	adv := adapter.DefaultAdvertisement()
	if err := adv.Configure(tinybl.AdvertisementOptions{
		LocalName:    b.Config.DeviceName, // Nordic UART Service
		ServiceUUIDs: []tinybl.UUID{serviceUUID},
	}); err != nil {
		return err
	}
	b.Advertisement = adv
	if err := adv.Start(); err != nil {
		return err
	}
	if err := adapter.AddService(&tinybl.Service{
		UUID: serviceUUID,
		Characteristics: []tinybl.CharacteristicConfig{
			{
				Handle: &b.RxChar,
				UUID:   rxUUID,
				Flags:  tinybl.CharacteristicWritePermission | tinybl.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client tinybl.Connection, offset int, value []byte) {
					fmt.Println("onreceive data", string(value))
					cmdType, resp, err := b.onData(value)
					cmdTypeStr := strconv.Itoa(int(cmdType))
					if err != nil {
						respErr := domain.ErrorResponse{
							Data:  value,
							Error: err.Error(),
						}
						bytes, _ := json.Marshal(respErr)
						msg := []byte(cmdTypeStr)
						msg = append(msg, []byte(":")...)
						msg = append(msg, bytes...)
						b.TxChar.Write(msg)
					}
					if resp != nil {
						msg := []byte(cmdTypeStr)
						msg = append(msg, []byte(":")...)
						msg = append(msg, resp...)
						b.TxChar.Write(msg)
					}
				},
			},
			{
				Handle: &b.TxChar,
				UUID:   txUUID,
				Flags:  tinybl.CharacteristicNotifyPermission | tinybl.CharacteristicReadPermission,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// Notify -.
func (b *Bluetooth) Notify() <-chan error {
	return b.notify
}

// Shutdown -.
func (b *Bluetooth) Shutdown() error {
	return b.Advertisement.Stop()
}
