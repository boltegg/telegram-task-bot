package main

import (
	"testing"
)

func TestEncryptor_EncryptString(t *testing.T) {
	enc := NewEncryptor("watafak")
	r, err := enc.EncryptString("pipika")
	if err != nil {
		t.Error(err)
	}
	if len(r) == 0 {
		t.Error("somenting wrong")
	}
	//fmt.Println(r)
}

func TestEncryptor_DecryptString(t *testing.T) {

	enc := NewEncryptor("watafak")
	r, err := enc.DecryptString("2/pM+KtolFCdNcm6k63MOU4hdGSaRdar1TO5x9m6+pXX/g==")
	if err != nil {
		t.Error(err)
	}
	if r != "pipika" {
		t.Error("decrypt error")
	}
}