// client.go
package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UpdateClientByUUID 更新指定 UUID 的客户端配置
func UpdateClientByUUID(config Client) error {
	db := GetDBInstance()
	result := db.Model(&Client{}).Where("uuid = ?", config.UUID).Updates(config)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetClientUUIDByToken 根据 Token 获取客户端 UUID
func GetClientUUIDByToken(token string) (uuid string, err error) {
	db := GetDBInstance()
	var client Client
	err = db.Where("token = ?", token).First(&client).Error
	if err != nil {
		return "", err
	}
	return client.UUID, nil
}

// CreateClient 创建新客户端
func CreateClient(config Client) (clientUUID, token string, err error) {
	db := GetDBInstance()
	token = generateToken()
	clientUUID = uuid.New().String()

	client := Client{
		UUID:        clientUUID,
		Token:       token,
		CPU:         config.CPU,
		GPU:         config.GPU,
		RAM:         config.RAM,
		SWAP:        config.SWAP,
		LOAD:        config.LOAD,
		UPTIME:      config.UPTIME,
		TEMP:        config.TEMP,
		OS:          config.OS,
		DISK:        config.DISK,
		NET:         config.NET,
		PROCESS:     config.PROCESS,
		Connections: config.Connections,
		Interval:    config.Interval,
	}

	err = db.Create(&client).Error
	if err != nil {
		return "", "", err
	}
	return clientUUID, token, nil
}

// GetAllClients 获取所有客户端配置
func GetAllClients() (clients []Client, err error) {
	db := GetDBInstance()
	err = db.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// GetClientConfig 获取指定 UUID 的客户端配置
func GetClientConfig(uuid string) (client Client, err error) {
	db := GetDBInstance()
	err = db.Where("uuid = ?", uuid).First(&client).Error
	if err != nil {
		return client, err
	}
	return client, nil
}

// ClientBasicInfo 客户端基本信息（假设的结构体，需根据实际定义调整）
type ClientBasicInfo struct {
	CPU       CPUReport       `json:"cpu"`
	GPU       GPUReport       `json:"gpu"`
	IpAddress IPAddressReport `json:"ip"`
	OS        string          `json:"os"`
}

// GetClientBasicInfo 获取指定 UUID 的客户端基本信息
func GetClientBasicInfo(uuid string) (client ClientBasicInfo, err error) {
	db := GetDBInstance()
	var clientInfo ClientInfo
	err = db.Where("client_uuid = ?", uuid).First(&clientInfo).Error
	if err != nil {
		return client, err
	}

	client = ClientBasicInfo{
		CPU: CPUReport{
			Name:  clientInfo.CPUNAME,
			Arch:  clientInfo.CPUARCH,
			Cores: clientInfo.CPUCORES,
		},
		GPU: GPUReport{
			Name: clientInfo.GPUNAME,
		},
		OS: clientInfo.OS,
		// IpAddress: 未在数据库中找到对应字段，需确认
	}

	return client, nil
}
