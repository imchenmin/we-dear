package main

import (
	"fmt"
	"log"
	"strings"
	"we-dear/config"
	"we-dear/models"
	"we-dear/utils"

	"github.com/mozillazg/go-pinyin"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置和数据库连接
	config.Init()
	config.InitDB()

	db := config.DB

	// 1. 首先添加新字段（允许为空）
	log.Println("Adding new columns...")
	err := db.Exec(`
		DO $$ 
		BEGIN 
			BEGIN
				ALTER TABLE patients ADD COLUMN username text NULL;
			EXCEPTION 
				WHEN duplicate_column THEN NULL;
			END;
			BEGIN
				ALTER TABLE patients ADD COLUMN password text NULL;
			EXCEPTION 
				WHEN duplicate_column THEN NULL;
			END;
			BEGIN
				ALTER TABLE patients ADD COLUMN salt text NULL;
			EXCEPTION 
				WHEN duplicate_column THEN NULL;
			END;
		END $$;
	`).Error
	if err != nil {
		log.Fatalf("Failed to add columns: %v", err)
	}

	// 2. 获取所有患者
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		log.Fatalf("Failed to fetch patients: %v", err)
	}

	// 用于检查用户名是否已存在
	usernameMap := make(map[string]bool)

	// 3. 遍历每个患者并更新信息
	log.Println("Updating patient records...")
	for _, patient := range patients {
		updates := make(map[string]interface{})

		// 处理用户名
		if patient.Username == "" {
			// 将姓名转换为拼音
			args := pinyin.NewArgs()
			args.Style = pinyin.FirstLetter
			pys := pinyin.Pinyin(patient.Name, args)

			// 合并拼音数组为字符串
			username := ""
			for _, py := range pys {
				username += py[0]
			}
			username = strings.ToLower(username)

			// 如果用户名已存在，添加数字后缀
			baseUsername := username
			suffix := 1
			for usernameMap[username] {
				username = fmt.Sprintf("%s%d", baseUsername, suffix)
				suffix++
			}

			updates["username"] = username
			usernameMap[username] = true
		}

		// 处理密码和盐
		if patient.Password == "" || patient.Salt == "" {
			salt, err := utils.GenerateSalt()
			if err != nil {
				log.Printf("Failed to generate salt for patient %s: %v", patient.ID, err)
				continue
			}

			// 使用默认密码 123456
			defaultPassword := "123456"
			hashedPassword := utils.HashPassword(defaultPassword, salt)

			updates["password"] = hashedPassword
			updates["salt"] = salt
		}

		// 如果有需要更新的字段
		if len(updates) > 0 {
			if err := db.Model(&patient).Updates(updates).Error; err != nil {
				log.Printf("Failed to update patient %s: %v", patient.ID, err)
				continue
			}
			log.Printf("Updated patient %s: name=%s, username=%s", patient.ID, patient.Name, updates["username"])
		}
	}

	// 4. 验证所有记录都有必要的值
	log.Println("Verifying all records have required values...")
	var count int64
	if err := db.Model(&models.Patient{}).Where("username IS NULL OR password IS NULL OR salt IS NULL").Count(&count).Error; err != nil {
		log.Fatalf("Failed to verify records: %v", err)
	}
	if count > 0 {
		log.Fatalf("Found %d records with missing values. Please check the data.", count)
	}

	// 5. 添加唯一索引和非空约束
	log.Println("Adding constraints...")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 添加唯一索引
		if err := tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_patients_username ON patients(username)").Error; err != nil {
			return fmt.Errorf("failed to create unique index: %v", err)
		}

		// 添加非空约束
		if err := tx.Exec("ALTER TABLE patients ALTER COLUMN username SET NOT NULL").Error; err != nil {
			return fmt.Errorf("failed to set username not null: %v", err)
		}
		if err := tx.Exec("ALTER TABLE patients ALTER COLUMN password SET NOT NULL").Error; err != nil {
			return fmt.Errorf("failed to set password not null: %v", err)
		}
		if err := tx.Exec("ALTER TABLE patients ALTER COLUMN salt SET NOT NULL").Error; err != nil {
			return fmt.Errorf("failed to set salt not null: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to add constraints: %v", err)
	}

	log.Println("Migration completed successfully")
	log.Println("Note: All patients without a password have been set to default password: 123456")
}
