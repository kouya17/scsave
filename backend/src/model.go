package scsave

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type JobState string

const (
  JobStateWaiting JobState = "waiting"
  JobStateRunning JobState = "running"
  JobStateCompleted JobState = "completed"
  JobStateFailed JobState = "failed"
)

func (s *JobState) Scan(value interface{}) error {
	*s = JobState(value.(string))
	return nil
}

func (s JobState) Value() (driver.Value, error) {
	return string(s), nil
}

type ScrapePropertiesJob struct {
  gorm.Model
  Url string
  Args string
  Type string
  State JobState `sql:"type:state" gorm:"default:waiting"`
  Progress uint `gorm:"default:0"` //0(waiting)-100(completed)
  Message string
  Tag string
}

type Property struct {
  gorm.Model
  Url string
  Price uint32
  LandArea float32
  BuildingArea float32
  Station string
  City string
  Layout string
  BuildYear uint16
  Access string
  Road string
  OtherCost string
  CoverageRatio string
  Timing string
  Rights string
  Structure string
  BuildCompany string
  Reform string
  LandKind string
  AreaPurpose string
  OtherRestriction string
  OtherNotice string
  JobId uint
  AreaPrice int64
  EstimatedBuildingPrice int64
  ClickCount uint64
  Revision uint16
}

type ResponseProperies struct {
  Count int64
  Properties []Property
}

type AreaPrice struct {
  gorm.Model
  Name string
  Price int64
}