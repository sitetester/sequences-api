package service

import (
	"github.com/sitetester/sequence-api/api"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SequenceService struct {
	Db *gorm.DB
}

func (ss *SequenceService) GetByID(id uint) *api.Sequence {
	var foundSequence api.Sequence
	ss.Db.Where("id = ?", id).First(&foundSequence)
	return &foundSequence
}

func (ss *SequenceService) GetByName(name string) *api.Sequence {
	var foundSequence api.Sequence
	ss.Db.Where("name = ?", name).First(&foundSequence)
	return &foundSequence
}

// GetOtherSequenceWithSameName https://gorm.io/docs/query.html#String-Conditions
func (ss *SequenceService) GetOtherSequenceWithSameName(name string, id uint) *api.Sequence {
	var otherSequence api.Sequence
	ss.Db.Where("name = ? AND id != ?", name, id).First(&otherSequence)
	return &otherSequence
}

func (ss *SequenceService) GetWithSteps(id uint64) *api.Sequence {
	var foundSequence api.Sequence
	ss.Db.Preload(clause.Associations).Where("id = ?", id).Find(&foundSequence)
	return &foundSequence
}

func (ss *SequenceService) Update(foundSequence *api.Sequence, sequence api.Sequence) {
	foundSequence.Name = sequence.Name
	foundSequence.OpenTrackingEnabled = sequence.OpenTrackingEnabled
	foundSequence.ClickTrackingEnabled = sequence.ClickTrackingEnabled

	ss.Db.Omit("SequenceStep").Save(&foundSequence)
}

func (sss *SequenceService) Create(sequence *api.Sequence) {
	sss.Db.Omit("SequenceStep").Create(&sequence)
}
