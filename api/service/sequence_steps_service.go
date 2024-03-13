package service

import (
	"github.com/sitetester/sequence-api/api"
	"gorm.io/gorm"
)

type SequenceStepsService struct {
	Db *gorm.DB
}

func (sss *SequenceStepsService) GetByID(id uint) *api.SequenceStep {
	var foundSequenceStep api.SequenceStep
	sss.Db.Where("id = ?", id).First(&foundSequenceStep)
	return &foundSequenceStep
}

// SubjectAvailablePerSequence https://gorm.io/docs/query.html#String-Conditions
func (sss *SequenceStepsService) SubjectAvailablePerSequence(subject string, sequenceID uint) bool {
	result := sss.Db.Where("subject = ? AND sequence_id = ?", subject, sequenceID).Find(&api.SequenceStep{})
	return result.RowsAffected == 0
}

func (sss *SequenceStepsService) Update(foundSequenceStep *api.SequenceStep, sequenceStep api.SequenceStep) {
	foundSequenceStep.Subject = sequenceStep.Subject
	foundSequenceStep.Content = sequenceStep.Content
	sss.Db.Save(&foundSequenceStep)
}

func (sss *SequenceStepsService) Create(sequenceStep *api.SequenceStep) {
	sss.Db.Create(&sequenceStep)
}

func (sss *SequenceStepsService) Delete(sequenceStep *api.SequenceStep) {
	sss.Db.Delete(&sequenceStep)
}
