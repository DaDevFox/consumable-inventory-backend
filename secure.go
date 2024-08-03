package main

// https://github.com/Dentrax/obscure-go/blob/master/types/int.go
// https://medium.com/trendyol-tech/secure-types-memory-safety-with-go-d3a20aa1e727

// const KEY int = 5920
//
// type SecureInt struct {
// 	key         int
// 	realValue   int
// 	fakeValue   int
// 	initialized bool
// }
//
// type ISecureInt interface {
// 	Apply() ISecureInt
// 	SetKey(int)
// 	Inc() ISecureInt
// 	Dec() ISecureInt
// 	Set(int) ISecureInt
// 	Get() int
// 	GetSelf() *SecureInt
// 	Decrypt() int
// 	RandomizeKey()
// 	IsEquals(ISecureInt) bool
// }
//
// func (i *SecureInt) Apply() ISecureInt {
// 	if !i.initialized {
// 		i.realValue = i.XOR(i.realValue, i.key)
// 		i.initialized = true
// 	}
//
// 	return i
// }
//
// func NewInt(value int) ISecureInt {
// 	s := &SecureInt{
// 		key:         KEY,
// 		realValue:   value,
// 		fakeValue:   value,
// 		initialized: false,
// 	}
// 	s.Apply()
// 	return s
// }
//
// func (i *SecureInt) XOR(value int, key int) int {
// 	return value ^ key
// }
// func (i *SecureInt) SetKey(key int) {
// 	i.key = key
// }
//
// func (i *SecureInt) RandomizeKey() {
// 	rand.Seed(time.Now().UnixNano())
// 	i.realValue = i.Decrypt()
// 	i.key = rand.Intn(int(^uint(0) >> 1))
// 	i.realValue = i.XOR(i.realValue, i.key)
// }
