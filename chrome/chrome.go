package chrome

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type Chrome struct {
	Context context.Context
}

func ProvideChrome() *Chrome {
	return &Chrome{}
}

func (c *Chrome) Init(ctx context.Context) (error, func()) {
	// chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // cause default option is headless
		chromedp.Flag("start-fullscreen", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("mute-audio", false), // cause default option is muted
		chromedp.Flag("no-first-run", true),
	)
	ctx, _ = chromedp.NewExecAllocator(ctx, opts...)
	ctx, cancel := chromedp.NewContext(ctx)
	c.Context = ctx

	return nil, cancel
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
