package aes

import (
	"testing"
)

func TestDecryptJs(t *testing.T) {

	got, err := AesDecrypt([]byte("a8ee9b4cfee2c682f382f1c292a93c6bb100be65cc34ba4246b45d744e8e3abf"), []byte("1595231954955"))
	if err != nil {
		t.Errorf("DecryptJs() error = %v", err)
		return
	}
	t.Log(string(got))
}
func TestDecryptJs2(t *testing.T) {

	got, err := AesDecrypt([]byte("e1d4d9ae99fa43e5d968bae170f6ca25e07cfdf41d738f156e7b076b3074067b"), []byte("1595241274949"))
	if err != nil {
		t.Errorf("DecryptJs() error = %v", err)
		return
	}
	t.Log(string(got))
}

func TestJsAesEncrypt(t *testing.T) {
	raw, key := "Admin@arrian.007", "1595231954955"
	d, err := AesEncrypt(raw, key)
	if err != nil {
		t.Error(err)
		return
	}
	if d != "a8ee9b4cfee2c682f382f1c292a93c6bb100be65cc34ba4246b45d744e8e3abf" {
		t.Error(d)
		return
	}

	got, err := AesDecrypt([]byte(d), []byte(key))
	if err != nil {
		t.Errorf("DecryptJs() error = %v", err)
		return
	}
	if string(got) != raw {
		t.Error("failed")
	}
}
