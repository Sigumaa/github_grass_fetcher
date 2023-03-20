package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

func main() {
	username := flag.String("user", "", "GitHub username")
	theme := flag.String("theme", "light", "Theme for the GitHub page")
	out := flag.String("out", "contributions.png", "Output file name")
	flag.Parse()

	if *username == "" {
		fmt.Fprintln(os.Stderr, "GitHub username is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*username, *theme, *out); err != nil {
		log.Fatal(err)
	}
}

func run(username, theme, out string) error {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	selector := ".js-yearly-contributions"
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(fmt.Sprintf("https://github.com/%s", username), selector, &buf, theme)); err != nil {
		return err
	}

	if err := os.WriteFile(out, buf, 0o644); err != nil {
		return err
	}

	return nil
}
func elementScreenshot(urlstr, sel string, res *[]byte, theme string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EvaluateAsDevTools(`document.querySelector("html").setAttribute("data-color-mode", "`+theme+`");`, nil),
		chromedp.WaitVisible(sel),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
