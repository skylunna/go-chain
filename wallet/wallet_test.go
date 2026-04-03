package wallet

import (
	"testing"
)

func TestGenerateWallet(t *testing.T) {
	wallet, err := GenerateWallet()
	if err != nil {
		t.Fatalf("生成钱包失败: %v", err)
	}

	if wallet.PrivateKey == nil {
		t.Error("私钥不应为空")
	}

	if wallet.PublicKey == nil {
		t.Error("公钥不应为空")
	}

	addr := wallet.Address()
	if len(addr) == 0 {
		t.Error("地址不应为空")
	}

	t.Logf("生成钱包成功! 地址: %s", addr)
}

func TestSignAndVerify(t *testing.T) {
	// 生成钱包
	wallet, _ := GenerateWallet()

	// 准备数据
	data := []byte("Hello, Blockchain! ")

	// 签名
	signature, err := wallet.Sign(data)
	if err != nil {
		t.Fatalf("签名失败: %v", err)
	}

	// 验证 使用正确的公钥
	if !Verify(wallet.PublicKey, data, signature) {
		t.Error("签名验证应该成功")
	}

	// 验证 使用错误的数据
	wrongData := []byte("Hacked data!")
	if Verify(wallet.PublicKey, wrongData, signature) {
		t.Error("使用错误数据验证应该失败")
	}

	// 验证 使用错误的公钥
	wrongWallet, _ := GenerateWallet()
	if Verify(wrongWallet.PublicKey, data, signature) {
		t.Error("使用错误公钥验证应该失败")
	}

	t.Log("签名和验证测试通过!")
}

func TestPrivateKeyExportImport(t *testing.T) {
	// 生成钱包
	wallet1, _ := GenerateWallet()
	addr1 := wallet1.Address()

	// 导出私钥
	privateKeyHex := wallet1.PrivateKeyToHex()
	t.Logf("导出私钥: %s", privateKeyHex)

	// 导入私钥
	importedKey, err := HexToPrivateKey(privateKeyHex)
	if err != nil {
		t.Fatalf("导入私钥失败: %v", err)
	}

	// 用导入的私钥创建钱包
	wallet2 := &Wallet{
		PrivateKey: importedKey,
		PublicKey:  &importedKey.PublicKey,
	}

	// 验证地址是否相同
	addr2 := wallet2.Address()
	if addr1 != addr2 {
		t.Errorf("地址不匹配: %s != %s", addr1, addr2)
	}

	t.Log("私钥导入导出测试通过! ")
}

func TestPublicKeyBytes(t *testing.T) {
	wallet, _ := GenerateWallet()

	pubBytes := wallet.PublicKeyBytes()
	if len(pubBytes) != 64 {
		t.Errorf("公钥字节长度应为 64, 实际: %d", len(pubBytes))
	}

	// 还原公钥
	pubKey, err := BytesToPublicKey(pubBytes)
	if err != nil {
		t.Fatalf("还原公钥失败: %v", err)
	}

	// 验证公钥是否相同
	if wallet.PublicKey.X.Cmp(pubKey.X) != 0 || wallet.PrivateKey.Y.Cmp(pubKey.Y) != 0 {
		t.Error("还原的公钥与原公钥不匹配")
	}

	t.Log("公钥字节转换测试通过! ")
}
