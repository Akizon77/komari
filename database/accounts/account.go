// database/account.go
package database

import (
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CheckPassword 检查密码是否正确
//
// 如果密码正确，返回用户的 UUID 和 true；否则返回空字符串和 false
func CheckPassword(username, passwd string) (uuid string, success bool) {
	db := GetDBInstance()
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return "", false
	}
	if hashPasswd(passwd) != user.Passwd {
		return "", false
	}
	return user.UUID, true
}

// ForceResetPassword 强制重置用户密码
func ForceResetPassword(username, passwd string) (err error) {
	db := GetDBInstance()
	result := db.Model(&User{}).Where("username = ?", username).Update("passwd", hashPasswd(passwd))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// hashPasswd 对密码进行加盐哈希
func hashPasswd(passwd string) string {
	saltedPassword := passwd + constantSalt
	hash := sha256.New()
	hash.Write([]byte(saltedPassword))
	hashedPassword := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return hashedPassword
}

// CreateDefaultAdminAccount 创建默认管理员账户
func CreateDefaultAdminAccount() (username, passwd string, err error) {
	db := GetDBInstance()

	username = os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	passwd = os.Getenv("ADMIN_PASSWORD")
	if passwd == "" {
		passwd = generatePassword()
	}

	hashedPassword := hashPasswd(passwd)

	user := User{
		UUID:     uuid.New().String(),
		Username: username,
		Passwd:   hashedPassword,
	}

	err = db.Create(&user).Error
	if err != nil {
		return "", "", err
	}

	return username, passwd, nil
}
