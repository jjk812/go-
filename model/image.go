package model

type Image struct {
	// 使用 `gorm:"primaryKey"` 标签来指明这个字段是主键
	Username string `gorm:"primaryKey"`
	ImgURL   string `gorm:"column:imgurl"` // 指定数据库中对应的列名
}
