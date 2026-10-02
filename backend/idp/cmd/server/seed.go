package main

import (
	"context"
	"os"
	"time"

	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/model"
	"github.com/cloudnativedaysjp/cnd-handson-app/backend/idp/internal/idp/service"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 契約テストの USER_ID（e2e/tests/contract/lib/grpc.ts）。UI テストはこの ID のプロジェクトを作る
var demoUserID = uuid.MustParse("00000000-0000-4000-8000-000000000071")

func seed(ctx context.Context, conn *gorm.DB) error {
	db := conn.WithContext(ctx)
	var member model.Role
	for _, name := range []string{"admin", service.DefaultRole} {
		role := model.Role{ID: uuid.New(), Name: name}
		if err := db.Where(model.Role{Name: name}).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		if name == service.DefaultRole {
			member = role
		}
	}

	email, password := os.Getenv("IDP_DEMO_EMAIL"), os.Getenv("IDP_DEMO_PASSWORD")
	if email == "" || password == "" {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now()
	demo := model.User{
		ID:           demoUserID,
		Name:         "demo",
		Email:        email,
		PasswordHash: string(hash),
		RoleID:       member.ID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&demo).Error
}
