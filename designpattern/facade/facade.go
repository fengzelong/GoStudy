package facade

import "fmt"

// apiImpl facade实现
type apiImpl struct {
	a AModuleAPI
	b BModuleAPI
}

// FaceApi 暴露给外部的接口
type FaceApi interface {
	FinalCall() string
}

func NewFaceApi() FaceApi {
	return &apiImpl{
		a: &aModuleImpl{},
		b: &bModuleImpl{},
	}
}

func (a *apiImpl) FinalCall() string {
	aRes, _ := a.a.FuncA()
	bRes := a.b.FuncB()
	return fmt.Sprintf("%s\n%s", aRes, bRes)
}

// aModuleImpl a模块标记结构体
type aModuleImpl struct{}

// AModuleAPI a模块接口签名
type AModuleAPI interface {
	FuncA() (string, bool)
}

// FuncA a模块接口实现
func (*aModuleImpl) FuncA() (string, bool) {
	return "A module running", true
}

type bModuleImpl struct{}

type BModuleAPI interface {
	FuncB() string
}

func (*bModuleImpl) FuncB() string {
	return "B module running"
}
