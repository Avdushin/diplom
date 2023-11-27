package main

import (
	"time"
	"wtd/initializers"
	"wtd/pkg/routes"

	"gorm.io/gorm"
)

// Migration - структура для миграции
type Migration struct {
	ID   string
	Up   func(db *gorm.DB) error
	Down func(db *gorm.DB) error
}

// PasswordResetRequest - структура для модели PasswordResetRequest
type PasswordResetRequest struct {
	gorm.Model
	UserID    uint
	Token     string
	ExpiresAt time.Time
}

var migrations []*Migration

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()

	migrations = append(migrations, &Migration{
		ID: "20231128010816_create_password_reset_requests_table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&PasswordResetRequest{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&PasswordResetRequest{})
		},
	})

	// Применяем миграции
	for _, migration := range migrations {
		err := migration.Up(initializers.DB)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// @SETUP ROUTES
	r := routes.SetupRoutes()

	// @RUN DEV SERVER
	r.Run()

	// @RUN PRODUCTION SERVER WITH SSL
	// r.RunTLS(vars.PORT, vars.Cert, vars.Key)
}
