package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	stravaLoginURL = "https://www.strava.com/login"
	flagPOSTURL    = "https://www.strava.com/flags"
)

func activityURL(activityID string) string {
	return "https://www.strava.com/activities/" + activityID
}

func newFlagURL(activityID string) string {
	return activityURL(activityID) + "/flags/new"
}

type Args struct {
	SessionID   string `arg:"required,env:STRAVA_SESSION"`
	ActivityID  string `arg:"required"`
	FlagComment string `arg:"required"`
}

func main() {
	var args Args
	arg.MustParse(&args)
	if err := run(args.SessionID, args.ActivityID, args.FlagComment); err != nil {
		log.Fatal(err)
	}
}

func isActivityFlagged(activityID string, nodes *[]*cdp.Node) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(activityURL(activityID)),
		chromedp.WaitVisible(`.activity-name`),
		chromedp.Nodes(`.flagged`, nodes, chromedp.AtLeast(0)),
	}
}

func flagThatShit(activityID, flagComment string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(newFlagURL(activityID)),
		chromedp.WaitVisible(`#flag_comment`),
		chromedp.SendKeys(`#flag_comment`, flagComment),
		chromedp.Submit(`#flag_comment`),
		chromedp.WaitVisible(`.activity-name`),
	}
}

// copied from https://github.com/chromedp/examples/blob/da6357a0fc35/cookie/main.go
func setStravaCookies(cookies ...string) chromedp.Tasks {
	if len(cookies)%2 != 0 {
		panic("length of cookies must be divisible by 2")
	}
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			for i := 0; i < len(cookies); i += 2 {
				err := network.SetCookie(cookies[i], cookies[i+1]).
					WithExpires(&expr).
					WithDomain(".strava.com").
					WithHTTPOnly(true).
					Do(ctx)
				if err != nil {
					return err
				}
			}
			return nil
		}),
	}
}

func run(sessionID, activityID, flagComment string) error {
	dir, err := os.MkdirTemp("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	if err := chromedp.Run(taskCtx, setStravaCookies("_strava4_session", sessionID)); err != nil {
		return err
	}

	flaggedNodes := []*cdp.Node{}
	if err := chromedp.Run(taskCtx, isActivityFlagged(activityID, &flaggedNodes)); err != nil {
		return err
	}
	if len(flaggedNodes) > 0 {
		log.Print("activity flagged, done")
		return nil
	}

	log.Print("activity not flagged, flagging it")
	if err := chromedp.Run(taskCtx, flagThatShit(activityID, flagComment)); err != nil {
		return err
	}

	return nil
}
