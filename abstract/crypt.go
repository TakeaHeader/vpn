package abstract

// 抽象 加解密
type Cryptable interface {
	//解密
	Decrypt([]byte) ([]byte, error)

	//加密
	Encrypt([]byte) ([]byte, error)
}
