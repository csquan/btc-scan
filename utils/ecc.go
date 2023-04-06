package utils

import (
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/suiguo/hwlib/ecies"
)

func GenKey() (pub string, pri string, generror error) {
	prv1, err := ecies.GenerateKey(rand.Reader, elliptic.P256(), nil)
	if err != nil {
		return "", "", fmt.Errorf(err.Msg())
	}
	return prv1.PublicKey.String(), prv1.String(), nil
}

type Encrypt interface {
	ECCEncrypt(msg []byte) ([]byte, error)
}
type Decrypt interface {
	ECCDecrypt(msg []byte) ([]byte, error)
}
type en_tool struct {
	pub *ecies.PublicKey
}

func (e *en_tool) ECCEncrypt(msg []byte) ([]byte, error) {
	if e.pub == nil {
		return nil, fmt.Errorf("pub key is nil")
	}
	out, err := ecies.Encrypt(rand.Reader, e.pub, msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf(err.Msg())
	}
	return out, nil
}

// 获取一个加密
func EnTool(pubstr string) (Encrypt, error) {
	pub, err := ecies.PublicFromString(pubstr)
	if err != nil {
		return nil, fmt.Errorf(err.Msg())
	}
	if pub == nil {
		return nil, fmt.Errorf("pub is nil")
	}
	return &en_tool{pub: pub}, nil
}

// 获取一个解密
func DeTool(pristr string) (Decrypt, error) {
	pri, err := ecies.PrivateFromString(pristr)
	if err != nil {
		return nil, fmt.Errorf(err.Msg())
	}
	if pri == nil {
		return nil, fmt.Errorf("pri is nil")
	}
	return &de_tool{pri: pri}, nil
}

type de_tool struct {
	pri *ecies.PrivateKey
}

func (d *de_tool) ECCDecrypt(msg []byte) ([]byte, error) {
	if d.pri == nil {
		return nil, fmt.Errorf("pri is nil")
	}
	pt, err := d.pri.Decrypt(msg, nil, nil)
	if err != nil {
		return nil, fmt.Errorf(err.Msg())
	}
	return pt, nil
}
