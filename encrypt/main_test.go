package main

import (
	"bytes"
	"testing"
)

func TestPKCS7PaddingAndUnPadding(t *testing.T) {
	data := []byte("hello")
	padded := PKCS7Padding(data, 8)
	if len(padded)%8 != 0 {
		t.Fatalf("填充后长度 = %d，期望能被 8 整除", len(padded))
	}

	got, err := PKCS7UnPadding(padded)
	if err != nil {
		t.Fatalf("去除填充失败: %v", err)
	}
	if !bytes.Equal(got, data) {
		t.Fatalf("去除填充结果 = %q，期望 %q", got, data)
	}
}

func TestEnPwdCodeAndDePwdCode(t *testing.T) {
	plain := []byte("GoStudy 加解密示例")
	encrypted, err := EnPwdCode(plain)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}
	if encrypted == string(plain) {
		t.Fatal("密文不应该等于明文")
	}

	decrypted, err := DePwdCode(encrypted)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}
	if !bytes.Equal(decrypted, plain) {
		t.Fatalf("解密结果 = %q，期望 %q", decrypted, plain)
	}
}

func TestDePwdCodeInvalidBase64(t *testing.T) {
	if _, err := DePwdCode("不是 base64"); err == nil {
		t.Fatal("非法 base64 应该返回错误")
	}
}
