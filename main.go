package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

var e *Encryptor
var mongo Mongo

func main() {

	InitConfig()
	mongo = NewMongo()

	e = NewEncryptor(os.Getenv("PASSPHRASE"), func(e error) { logrus.Errorf("Encryptor error: %s", e.Error()) })

	ProcessTelegramBot()

	// temp
	//time.Sleep(time.Hour)
}
