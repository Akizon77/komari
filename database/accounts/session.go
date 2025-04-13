// session.go
package database

import (
	"errors"
	"time"
)

// GetAllSessions 获取所有会话
func GetAllSessions() (sessions []Session, err error) {
	db := GetDBInstance()
	err = db.Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// CreateSession 创建新会话
func CreateSession(uuid string, expires int) (string, error) {
	db := GetDBInstance()
	session := generateRandomString(32)

	sessionRecord := Session{
		UUID:    uuid,
		Session: session,
		Expires: time.Now().Add(time.Duration(expires) * time.Second),
	}

	err := db.Create(&sessionRecord).Error
	if err != nil {
		return "", err
	}
	return session, nil
}

// GetSession 根据会话 ID 获取 UUID
func GetSession(session string) (uuid string, err error) {
	db := GetDBInstance()
	var sessionRecord Session
	err = db.Where("session = ?", session).First(&sessionRecord).Error
	if err != nil {
		return "", err
	}

	if time.Now().After(sessionRecord.Expires) {
		// 会话已过期，删除它
		_ = DeleteSession(session)
		return "", errors.New("session expired")
	}

	return sessionRecord.UUID, nil
}

// DeleteSession 删除指定会话
func DeleteSession(session string) (err error) {
	db := GetDBInstance()
	result := db.Where("session = ?", session).Delete(&Session{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
