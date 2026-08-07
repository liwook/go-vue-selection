package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// DSN 仅用于代码生成阶段连库反向生成模型，运行期配置见 config.yaml
const dsn = "host=127.0.0.1 port=5432 user=postgres password=postgres123! dbname=vue_admin sslmode=disable search_path=app,public"

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		// ModelPkgPath 用不含分隔符的名字 "model"：getModelOutputPath() 在 Windows 上
		// 用 os.PathSeparator(\) 检测分隔符，写 "dal/model"(用/) 会误判，走 else 分支拼成
		// dal/dal/model。用 "model" 则 filepath.Join(filepath.Dir(OutPath), "model")=dal/model。
		// import 路径由 fillModelPkgPath 从实际目录反推为 vue_admin/dal/model。
		OutPath:           "dal/query",
		ModelPkgPath:      "model",
		Mode:              gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery,
		FieldNullable:     true,
		FieldCoverable:    false, //generate pointer when field has default value
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  false,
	})

	g.UseDB(db)

	all := g.GenerateAllTable()
	g.ApplyBasic(all...)
	g.Execute()

	log.Println("gen done: models -> dal/model, queries -> dal/query")
}
