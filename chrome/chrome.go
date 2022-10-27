package chrome

import (
	"context"

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

func (c *Chrome) OpenUrl(url string) error {
	err := chromedp.Run(c.Context, chromedp.Navigate(url))
	if err != nil {
		return err
	}
	return nil
}
