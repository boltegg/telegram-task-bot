package main

import "time"

func main() {

	InitConfig()

	go ProcessTelegramBot()


	// temp
	time.Sleep(time.Hour)
}