package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"

	"github.com/emmansun/gmsm/sm2"
	"github.com/emmansun/gmsm/sm3"
)

var (
	ErrInvalidP1      = errors.New("invalid P1 format")
	ErrInvalidQ1      = errors.New("invalid Q1 format")
	ErrInvalidE       = errors.New("invalid E format")
	ErrInvalidT1      = errors.New("invalid T1 format")
	ErrKeyGenFailed   = errors.New("key generation failed")
	ErrSignFailed     = errors.New("sign failed")
	ErrDecryptFailed  = errors.New("decrypt failed")
)

// SM2Curve SM2曲线参数
var SM2Curve = sm2.P256()

// N SM2曲线阶
var N = SM2Curve.Params().N

// SM2CoopKeyGenResult 协同密钥生成结果
type SM2CoopKeyGenResult struct {
	D2     []byte // 服务端私钥分量
	D2Inv  []byte // D2的逆
	P2     []byte // 服务端公钥分量
	Pa     []byte // 协同公钥
}

// SM2CoopSignResult 协同签名结果
type SM2CoopSignResult struct {
	R  []byte // 签名分量 r
	S2 []byte // 签名分量 s2
	S3 []byte // 签名分量 s3
}

// GenerateKeyPair 生成完整的 SM2 密钥对（用于注册时自动生成）
func GenerateKeyPair() (*sm2.PrivateKey, error) {
	return sm2.GenerateKey(rand.Reader)
}

// GenerateRandom 生成随机数
func GenerateRandom(length int) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	return b, err
}

// SM3Hash 计算 SM3 哈希
func SM3Hash(data []byte) []byte {
	h := sm3.New()
	h.Write(data)
	return h.Sum(nil)
}

// SM3HashWithPassword 计算密码哈希（加盐）
func SM3HashWithPassword(password, salt []byte) []byte {
	h := sm3.New()
	h.Write(salt)
	h.Write(password)
	return h.Sum(nil)
}

// CoopKeyGenInit 协同密钥生成初始化
// 输入: P1 - 客户端公钥分量 (64字节, 未压缩格式, 无前缀04)
// 输出: D2, D2Inv, P2, Pa
func CoopKeyGenInit(p1 []byte) (*SM2CoopKeyGenResult, error) {
	if len(p1) != 64 {
		return nil, ErrInvalidP1
	}

	// 解析 P1 为椭圆曲线点
	p1X := new(big.Int).SetBytes(p1[:32])
	p1Y := new(big.Int).SetBytes(p1[32:64])

	// 生成随机 d2 (服务端私钥分量)
	d2, err := rand.Int(rand.Reader, N)
	if err != nil {
		return nil, ErrKeyGenFailed
	}

	// 计算 d2Inv = d2^(-1) mod n
	d2Inv := new(big.Int).ModInverse(d2, N)

	// 计算 P2 = d2Inv * G
	p2X, p2Y := SM2Curve.ScalarBaseMult(d2Inv.Bytes())
	p2 := make([]byte, 64)
	copy(p2[:32], p2X.Bytes())
	copy(p2[32:64], p2Y.Bytes())

	// 计算 Pa = d2Inv * P1 + (n-1) * G
	// 先计算 d2Inv * P1
	paX, paY := SM2Curve.ScalarMult(p1X, p1Y, d2Inv.Bytes())
	// 再计算 (n-1) * G = -G (因为 (n-1) ≡ -1 mod n)
	minusOne := new(big.Int).Sub(N, big.NewInt(1))
	gX, gY := SM2Curve.ScalarBaseMult(minusOne.Bytes())
	// 最后相加
	paX, paY = SM2Curve.Add(paX, paY, gX, gY)
	pa := make([]byte, 64)
	copy(pa[:32], paX.Bytes())
	copy(pa[32:64], paY.Bytes())

	return &SM2CoopKeyGenResult{
		D2:    d2.Bytes(),
		D2Inv: d2Inv.Bytes(),
		P2:    p2,
		Pa:    pa,
	}, nil
}

// CoopSign 协同签名
// 输入: d2Inv - D2的逆, q1 - 客户端盲化因子 (64字节), e - 消息哈希 (32字节)
// 输出: r, s2, s3
func CoopSign(d2Inv, q1, e []byte) (*SM2CoopSignResult, error) {
	if len(q1) != 64 {
		return nil, ErrInvalidQ1
	}
	if len(e) != 32 {
		return nil, ErrInvalidE
	}

	// 解析 Q1 为椭圆曲线点
	q1X := new(big.Int).SetBytes(q1[:32])
	q1Y := new(big.Int).SetBytes(q1[32:64])

	// 生成随机 k2, k3
	k2, err := rand.Int(rand.Reader, N)
	if err != nil {
		return nil, ErrSignFailed
	}
	k3, err := rand.Int(rand.Reader, N)
	if err != nil {
		return nil, ErrSignFailed
	}

	// 计算 Q2 = k2 * G
	q2X, q2Y := SM2Curve.ScalarBaseMult(k2.Bytes())

	// 计算 x1 = k3 * Q1 + Q2
	x1X, x1Y := SM2Curve.ScalarMult(q1X, q1Y, k3.Bytes())
	x1X, x1Y = SM2Curve.Add(x1X, x1Y, q2X, q2Y)

	// 计算 r = (e + x1) mod n
	r := new(big.Int).SetBytes(e)
	r.Add(r, x1X)
	r.Mod(r, N)

	// 计算 s2 = d2Inv * k3 mod n
	d2InvBig := new(big.Int).SetBytes(d2Inv)
	s2 := new(big.Int).Mul(d2InvBig, k3)
	s2.Mod(s2, N)

	// 计算 s3 = d2Inv * (r + k2) mod n
	rPlusK2 := new(big.Int).Add(r, k2)
	s3 := new(big.Int).Mul(d2InvBig, rPlusK2)
	s3.Mod(s3, N)

	return &SM2CoopSignResult{
		R:  r.Bytes(),
		S2: s2.Bytes(),
		S3: s3.Bytes(),
	}, nil
}

// CoopDecrypt 协同解密
// 输入: d2Inv - D2的逆, t1 - 客户端密文变换 (64字节)
// 输出: t2
func CoopDecrypt(d2Inv, t1 []byte) ([]byte, error) {
	if len(t1) != 64 {
		return nil, ErrInvalidT1
	}

	// 解析 T1 为椭圆曲线点
	t1X := new(big.Int).SetBytes(t1[:32])
	t1Y := new(big.Int).SetBytes(t1[32:64])

	// 计算 T2 = d2Inv * T1
	t2X, t2Y := SM2Curve.ScalarMult(t1X, t1Y, d2Inv)
	t2 := make([]byte, 64)
	copy(t2[:32], t2X.Bytes())
	copy(t2[32:64], t2Y.Bytes())

	return t2, nil
}

// EncodeToBase64 将字节切片编码为 Base64 字符串
func EncodeToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeFromBase64 从 Base64 字符串解码为字节切片
func DecodeFromBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
