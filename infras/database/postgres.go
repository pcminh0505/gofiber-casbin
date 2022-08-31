package database

import (
	"fmt"

	"github.com/pcminh0505/gofiber-casbin/api/models"
	"github.com/pcminh0505/gofiber-casbin/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	admin = "ADMIN"
)

var adminDB *gorm.DB

// Connect to databases
func Connect() {
	adminDsn := getDataSourceName(admin)

	var err error
	adminDB, err = gorm.Open(postgres.Open(adminDsn), &gorm.Config{})

	if err != nil {
		panic("Cannot connect to database")
	}

	// Database migration
	adminDB.AutoMigrate(
		&models.User{},
	)

	// Auto create admin at first load
	if result := adminDB.First(&models.User{}).RowsAffected; result == 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte(config.GetEnv("ROOT_ADMIN_PASSWORD")), bcrypt.DefaultCost)
		data := models.User{
			Username: config.GetEnv("ROOT_ADMIN_USERNAME"),
			Password: string(password),
			Role:     config.GetEnv("ROOT_ADMIN_ROLE"),
		}
		adminDB.Create(&data)
		Casbin().AddGroupingPolicy(fmt.Sprint(data.ID), data.Role)
	}

}

func GetAdminDB() *gorm.DB {
	return adminDB
}

func getDataSourceName(db string) string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok",
		config.GetEnv("POSTGRES_USER_"+db),
		config.GetEnv("POSTGRES_PASSWORD_"+db),
		config.GetEnv("POSTGRES_HOST_"+db),
		config.GetEnv("POSTGRES_PORT_"+db),
		config.GetEnv("POSTGRES_DBNAME_"+db))
}
