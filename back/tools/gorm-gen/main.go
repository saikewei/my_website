package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const dsn = "backend_user:Hyc65319436@tcp(nas.saikewei.tech:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 使用 log.Fatalf 来打印详细错误并退出程序
		log.Fatalf("failed to connect database: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../internal/model/query",
		ModelPkgPath: "../../internal/model",

		Mode: gen.WithoutContext | gen.WithDefaultQuery,
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("albums"),
		g.GenerateModel("photo_metadata"),
		g.GenerateModel("photo_tags"),
		g.GenerateModel("photos"),
		g.GenerateModel("tags"),
	)

	// 执行生成
	g.Execute()
}
