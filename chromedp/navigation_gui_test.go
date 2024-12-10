package main_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"gotest.tools/v3/assert"
)

// test-101
func TestGuiNavgation(t *testing.T) {
	// Create a new context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Disable headless mode
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("window-size", "1280,1220"), // Set window size to 1280x720

	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a ##timeout## for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	waitTime := 500 * time.Millisecond
	var buf []byte
	// Navigate to a website and take a screenshot
	t.Log("Start Simulation")
	t.Log("Choose Test 1")
	testTestSelectButtonIsVisible := false
	// Selectors see https://www.w3schools.com/cssref/css_selectors.php
	mainUrl := "https://localhost:8181"
	selectorButton1 := "#button-test-1"
	if err := chromedp.Run(ctx,
		chromedp.Navigate(mainUrl),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(selectorButton1, chromedp.ByID),
		chromedp.Click(selectorButton1, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("screen-1.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	testTestSelectButtonIsVisible = true
	assert.Equal(t, testTestSelectButtonIsVisible, true, "Test 1 Button should be visible")

	// Next Page
	t.Log("Good boy test")
	testNextPageIsVisible := false
	selectorInput1 := "#goodboy"
	input1 := "me"
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(selectorInput1, chromedp.ByID),
		chromedp.Focus(selectorInput1),
		// chromedp.SetValue(questionAreaSelector, targetQuestion, chromedp.ByQuery),
		chromedp.SendKeys(selectorInput1, input1, chromedp.ByQuery),
		chromedp.Sleep(2000*time.Millisecond),
		chromedp.Click(selectorButton1, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		log.Fatal(err)
	} else {
		testNextPageIsVisible = true
	}
	if err := os.WriteFile("screen-2.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, testNextPageIsVisible, true, "Next Page should be visible")

	// Previous Page
	t.Log("Previous page button")
	testPrevPageIsVisible := false
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#button-prev-page`, chromedp.ByID),
		chromedp.Click(`#button-prev-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
	); err != nil {
		log.Fatal(err)
	}
	testPrevPageIsVisible = true
	assert.Equal(t, testPrevPageIsVisible, true, "Prev Page should be visible")

	t.Log("Navigate to the last page")
	// Navgitage to the last page (5 times)
	testSubmitButtonIsVisible := false
	if err := chromedp.Run(ctx,
		chromedp.WaitVisible(`#button-next-page`, chromedp.ByID),
		chromedp.Click(`#button-next-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(`#button-next-page`, chromedp.ByID),
		chromedp.Click(`#button-next-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(`#button-next-page`, chromedp.ByID),
		chromedp.Click(`#button-next-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(`#button-next-page`, chromedp.ByID),
		chromedp.Click(`#button-next-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(`#button-next-page`, chromedp.ByID),
		chromedp.Click(`#button-next-page`, chromedp.ByID),
		chromedp.Sleep(waitTime),
		chromedp.WaitVisible(`#button-submit-test`, chromedp.ByID),
		chromedp.CaptureScreenshot(&buf),
	); err != nil {
		t.Error("Navigation to the last page failed")
		log.Fatal(err)
	}
	if err := os.WriteFile("screen-100.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}
	testSubmitButtonIsVisible = true
	assert.Equal(t, testSubmitButtonIsVisible, true, "Submit Button should be visible on the last page")
	// Look at last Page
	if err := chromedp.Run(ctx,
		chromedp.Sleep(2*waitTime),
	); err != nil {
		log.Fatal(err)
	}

	// Save the screenshot
	if err := os.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
