package vegenerelib

import "testing"

func TestEncrypt(t *testing.T) {
	expected_result := "LIUKMGQ SLW XYF XABIC\n"
	if Encrypt("enc_test.txt", "secretkey") != expected_result {
		t.Error(`Encrypt("enc_test.txt", "secretkey") = false`)
	}
}

func TestDecrypt(t *testing.T) {
	expected_result := "TESTING ONE TWO THREE\n"
	if Decrypt("dec_test.txt", "secretkey") != expected_result {
		t.Error(`Decrypt("dec_test.txt", "secretkey") = false`)
	}
}

func TestDecryptKey(t *testing.T) {
	expected_result := "NEWWORLD"
	if DecryptKey("key_test.txt") != expected_result {
		t.Error(`DecryptKey("key_test.txt") = false`)
	}
}
