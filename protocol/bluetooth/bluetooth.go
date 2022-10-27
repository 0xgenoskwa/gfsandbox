package bluetooth

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/0xgenoskwa/gfsandbox/config"
	"github.com/0xgenoskwa/gfsandbox/domain"
	"github.com/0xgenoskwa/gfsandbox/handler"
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

	Handler *handler.Handler
}

func ProvideBluetooth(c *config.Config, h *handler.Handler) *Bluetooth {
	return &Bluetooth{
		Config:  c,
		Handler: h,
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
					cmdType, resp, err := b.Handler.OnData(value)
					if err != nil {
						respErr := domain.ErrorResponse{
							Data:  value,
							Error: err.Error(),
						}
						bytes, _ := json.Marshal(respErr)
						b.TxChar.Write(bytes)
					}
					if resp != nil {
						cmdTypeStr := strconv.Itoa(int(cmdType))
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

func (b *Bluetooth) Stop() {
	fmt.Println("stop bluetooth")
	err := b.Advertisement.Stop()
	fmt.Println("stop bluetooth error", err)
}
