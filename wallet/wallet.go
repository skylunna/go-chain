package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Wallet 包含公私钥对 代表一个区块链钱包
type Wallet struct {
	// 不直接存结构体，只存其引用地址
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// 由公钥生成 (简化版: 公钥哈希)，交易地址
func (w *Wallet) Address() string {
	// 将公钥的 X 和 Y 坐标拼接，将 Y 的所有字节，一个个加到 X 中
	pubKey := append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)
	// 进行 SHA256 哈希
	hash := sha256.Sum256(pubKey)
	// 返回 十六进制 字符串
	return hex.EncodeToString(hash[:])
}

// 获取公钥字的字节表示 (用于交易中)
func (w *Wallet) PublicKeyBytes() []byte {
	return append(w.PublicKey.X.Bytes(), w.PublicKey.Y.Bytes()...)
}

// GenerateWallet 生成新钱包
func GenerateWallet() (*Wallet, error) {
	// 生成 椭圆曲线私钥（带随机数）
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// 用私钥创建钱包
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// 对数据签名
func (w *Wallet) Sign(data []byte) ([]byte, error) {
	// 先对数据进行哈希
	hash := sha256.Sum256(data)
	// 使用私钥签名
	signature, err := ecdsa.SignASN1(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// 验证签名 (静态方法，使用公钥验证)
func Verify(pubKey *ecdsa.PublicKey, data []byte, signature []byte) bool {
	// 对数据进行哈希
	hash := sha256.Sum256(data)
	// 使用公钥验证签名
	return ecdsa.VerifyASN1(pubKey, hash[:], signature)
}

// BytesToPublicKey 将字节数组还原为公钥
func BytesToPublicKey(pubBytes []byte) (*ecdsa.PublicKey, error) {
	// P256 曲线的公钥坐标长度是 32 字节
	// X 32 + Y 32 = 64
	if len(pubBytes) != 64 {
		return nil, fmt.Errorf("无效的公钥长度: %d", len(pubBytes))
	}

	x := new(big.Int).SetBytes(pubBytes[:32])
	y := new(big.Int).SetBytes(pubBytes[32:])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(), // 曲线类型
		X:     x,
		Y:     y,
	}, nil
}

// 将私钥导出为十六进制(用于备份)
func (w *Wallet) PrivateKeyToHex() string {
	return hex.EncodeToString(w.PrivateKey.D.Bytes())
}

// 从十六进制导入私钥
func HexToPrivateKey(hexStr string) (*ecdsa.PrivateKey, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}

	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = new(big.Int).SetBytes(bytes)

	// 根据私钥推导出公钥
	privateKey.PublicKey.X, privateKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(bytes)

	return privateKey, nil
}

// // 将16进制字符串转回公钥（用于验证交易）
// func HexToPublicKey(hexStr string) (*ecdsa.PublicKey, error) {
// 	// 简化处理，实际需要从 X, Y 坐标解析
// 	// 为了方便，直接在交易里存公钥字节
// 	return nil, fmt.Errorf("not implemented in this snippet")
// }
