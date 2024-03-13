package api

// Sequence https://gorm.io/docs/models.html#Conventions
// DB table name will be `sequences` (plural)
type Sequence struct {
	ID                   uint   `gorm:"primaryKey"`
	Name                 string ` valid:"alphanum,required,minstringlength(3),maxstringlength(30)" gorm:"unique"`
	OpenTrackingEnabled  bool
	ClickTrackingEnabled bool
	SequenceSteps        []SequenceStep `json:"-"` // wouldn't show in JSON output
}

// SequenceStep https://gorm.io/docs/has_many.html#Has-Many
type SequenceStep struct {
	ID         uint   `gorm:"primaryKey"`
	Subject    string `valid:"required,minstringlength(3)"`
	Content    string `valid:"required,minstringlength(3)"`
	SequenceID uint
}

type SequenceWithSteps struct {
	Sequence *Sequence
	Steps    *[]SequenceStep
}

type ErrorResponse struct {
	Error string
}
