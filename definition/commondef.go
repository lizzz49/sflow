package definition

type Definition struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:定义Id"`
	Name string `json:"name" gorm:"uniqueIndex;column:name;type:varchar(255);comment:定义名"`
}
