package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("requires GitHub's username on arg(1)")
		os.Exit(1)
	}
	username := os.Args[1]

	if err := run(username); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(username string) error {
	theme := "light"
	if len(os.Args) >= 3 {
		theme = os.Args[2]
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("https://github.com/%s", username)),
		chromedp.EvaluateAsDevTools(`
            const elementList = document.getElementsByClassName("flex-shrink-0");
            let i = 0;
            while (true) {
                const element = elementList.item(i);
                if (!element) {
                    break;
                }
                element.remove();
                i++;
            }

            document.querySelector("html").setAttribute("data-color-mode", "`+theme+`");
        `, &buf),
		chromedp.WaitVisible(".js-yearly-contributions"),
		chromedp.Screenshot(".js-yearly-contributions", &buf, chromedp.NodeVisible),
	)

	if err != nil {
		return err
	}

	file, err := os.Create("output.jpeg")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
