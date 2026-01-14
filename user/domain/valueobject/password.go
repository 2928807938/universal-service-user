package valueobject

import (
	"github.com/2928807938/universal-service-user/user/domain/errors"
	"golang.org/x/crypto/bcrypt"
)

// Password 密码值对象
type Password struct {
	hash string
}

// NewPassword 从明文创建密码值对象（自动加密）
// 使用 cost=12 提高安全性（默认是10，每增加1计算时间翻倍）
func NewPassword(plainPassword string) (*Password, error) {
	if plainPassword == "" {
		return nil, errors.ErrPasswordStrengthWeak
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return nil, err
	}
	return &Password{hash: string(hash)}, nil
}

// NewPasswordFromHash 从哈希值创建密码值对象
func NewPasswordFromHash(hash string) (*Password, error) {
	if hash == "" {
		return nil, errors.ErrPasswordStrengthWeak
	}
	return &Password{hash: hash}, nil
}

// String 返回哈希值（注意：不返回明文密码）
func (p *Password) String() string {
	return p.hash
}

// Hash 返回密码哈希
func (p *Password) Hash() string {
	return p.hash
}

// Verify 验证密码是否匹配
func (p *Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
	return err == nil
}
