package main

import (
	app "github.com/RB-PRO/labexp/internal/tgsecret"
)

func main() {
	a, err := app.New()
	if err != nil {
		return
	}

	err = a.Run()
	if err != nil {
		return
	}
}
