package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

func screenshot(urlStr string) ([]byte, error) {
	var (
		buf []byte
		ctx context.Context
	)
	if remoteUrl != "" {
		wsAddr, err := addrFromRemoteUrl()
		if err != nil {
			return nil, err
		}
		ctx, _ = context.WithTimeout(context.Background(), timeout)
		ctx, _ = chromedp.NewRemoteAllocator(ctx, wsAddr)
	} else {
		if debug {
			chromedp.DefaultExecAllocatorOptions = []chromedp.ExecAllocatorOption{
				chromedp.NoFirstRun,
				chromedp.NoDefaultBrowserCheck,
			}
		}
		ctx, _ = context.WithTimeout(context.Background(), timeout)
	}
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	if err := chromedp.Run(ctx, screenshotAction(urlStr, 90, &buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

func addrFromRemoteUrl() (string, error) {
	resp, err := http.Get(remoteUrl + "/json/version")
	if err != nil {
		return "", err
	}
	var result map[string]string
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(b, &result)
	if v, ok := result["webSocketDebuggerUrl"]; ok {
		return v, nil
	}
	return "", errors.New("无法获取webSocketDebuggerUrl")
}

func screenshotAction(url string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}
			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
