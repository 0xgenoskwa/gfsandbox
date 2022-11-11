package chrome

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"go.genframe.xyz/domain"
)

type Chrome struct {
	Context context.Context
}

func ProvideChrome() *Chrome {
	return &Chrome{}
}

func (c *Chrome) Init(ctx context.Context) (func(), error) {
	// chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // cause default option is headless
		chromedp.Flag("start-fullscreen", true),
		chromedp.Flag("disable-gpu", false),
		// chromedp.Flag("disable-web-security", true),
		chromedp.Flag("no-default-browser-check", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("mute-audio", false), // cause default option is muted
		chromedp.Flag("no-first-run", true),
	)
	ctx, _ = chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(ctx)
	c.Context = ctx

	return cancel, nil
}

func (c *Chrome) OpenHtml() error {
	err := chromedp.Run(c.Context, chromedp.Navigate("http://localhost"))
	if err != nil {
		return err
	}

	return nil
}

func (c *Chrome) OpenUrl(url string) error {
	err := chromedp.Run(c.Context, chromedp.Tasks{
		chromedp.SetAttributeValue("#genframe", "src", url, chromedp.ByID),
		chromedp.SetAttributeValue("#genframe", "styles", "display: block", chromedp.ByID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Chrome) Toast(text string) error {
	err := chromedp.Run(c.Context, chromedp.Tasks{
		// chromedp.Evaluate(fmt.Sprintf(`document.getElementById("notification-content").innerHTML = "%s";`, text), nil),
		chromedp.SetJavascriptAttribute(`#notification-content`, "textContent", text, chromedp.ByQuery),
		chromedp.SetAttributeValue("#notification", "style", "display: flex", chromedp.ByID),
	})
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(8 * time.Second)
		chromedp.Run(c.Context, chromedp.Tasks{
			chromedp.SetJavascriptAttribute(`#notification-content`, "textContent", text, chromedp.ByQuery),
			chromedp.SetAttributeValue("#notification", "style", "display: none", chromedp.ByID),
		})
	}()
	return nil
}

func (c *Chrome) SendTouchEvent(event string, value domain.TouchEvent) error {
	err := chromedp.Run(c.Context, chromedp.Tasks{
		chromedp.QueryAfter("body", func(ctx context.Context, eci runtime.ExecutionContextID, n ...*cdp.Node) error {
			boxes, err := dom.GetContentQuads().WithNodeID(n[0].NodeID).Do(ctx)
			if err != nil {
				return err
			}
			content := boxes[0]

			c := len(content)
			if c%2 != 0 || c < 1 {
				return chromedp.ErrInvalidDimensions
			}

			var x, y float64
			for i := 0; i < c; i += 2 {
				x += content[i]
				y += content[i+1]
			}
			x /= float64(c / 2)
			y /= float64(c / 2)

			touchType := input.TouchStart
			if event == "touchend" {
				touchType = input.TouchEnd
			}
			if event == "touchcancel" {
				touchType = input.TouchCancel
			}
			if event == "touchmove" {
				touchType = input.TouchMove
			}

			touchPoints := []*input.TouchPoint{}
			for _, item := range value.Touches {
				touchPoints = append(touchPoints, &input.TouchPoint{
					X:             item.PageX,
					Y:             item.PageY,
					RadiusX:       item.RadiusX,
					RadiusY:       item.RadiusY,
					RotationAngle: item.RorationAngle,
					Force:         item.Force,
				})
			}

			p := input.DispatchTouchEventParams{
				Type:        touchType,
				TouchPoints: touchPoints,
			}
			if err := p.Do(ctx); err != nil {
				return err
			}
			return nil
		}),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Chrome) SendKeyEvent(value string) error {
	err := chromedp.Run(c.Context, chromedp.Tasks{
		chromedp.SendKeys("#genframe", value, chromedp.ByID),
	})
	if err != nil {
		return err
	}
	return nil
}
