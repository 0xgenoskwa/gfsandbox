package bluetooth

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("start bluetooth")
	adapter := tinybl.DefaultAdapter
	// set connect handler
	adapter.SetConnectHandler(func(device tinybl.Addresser, connected bool) {
		if connected {
			fmt.Println("connected, not advertising...")
			b.AdvertisementState = false
		} else {
			fmt.Println("disconnected, advertising...")
			b.AdvertisementState = true
		}
	})

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
					resp, err := b.Handler.OnData(value)
					if err != nil {
						respErr := domain.ErrorResponse{
							Data:  value,
							Error: err.Error(),
						}
						bytes, _ := json.Marshal(respErr)
						b.TxChar.Write(bytes)
					}
					if resp != nil {
						b.TxChar.Write(resp)
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
