package main

import (
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kouya17/scsave"
)

var DB *gorm.DB

func main() {
  dsn := "host=db user=postgres password=postgres dbname=scsave port=5432 sslmode=disable TimeZone=Asia/Tokyo"
  var err error
  DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  err = DB.AutoMigrate(&scsave.Property{})
  if err != nil {
    panic("failed migrate property")
  }
  err = DB.AutoMigrate(&scsave.ScrapePropertiesJob{})
  if err != nil {
    panic("failed migrate job")
  }
  err = DB.AutoMigrate(&scsave.AreaPrice{})
  if err != nil {
    panic("failed migrate job")
  }

  var job scsave.ScrapePropertiesJob
  tx := DB.Begin()
  resultSelect := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&job, "state = ?", string(scsave.JobStateWaiting))
  if resultSelect.Error != nil {
    log.Fatal(resultSelect.Error)
  }
  fmt.Printf("selected job: %#v\n", job)
  resultUpdate := tx.Model(&job).Update("state", string(scsave.JobStateRunning))
  if resultUpdate.Error != nil {
    log.Fatal(resultUpdate.Error)
  }
  fmt.Printf("updated job: %#v\n", job)
  tx.Commit()

  properties, jobType := scsave.GetProperties(job.Url, job.Args)
  for _, v := range properties {
    v.JobId = job.ID
    var cachedAreaPrice scsave.AreaPrice
    var areaPrice int64 = -1
    err := DB.First(&cachedAreaPrice, "name = ?", v.Station + "駅").Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
      areaPrice = scsave.FetchAreaPrice(v.Station + "駅")
      cachingAreaPrice := scsave.AreaPrice{}
      cachingAreaPrice.Name = v.Station + "駅"
      cachingAreaPrice.Price = areaPrice
      DB.FirstOrCreate(&cachingAreaPrice, "name = ?", cachingAreaPrice.Name)
    } else {
      areaPrice = cachedAreaPrice.Price
    }
    if areaPrice == -1 {
      // FIXME: 上の処理と共通化
      cachedAreaPrice = scsave.AreaPrice{}
      err := DB.First(&cachedAreaPrice, "name = ?", v.City).Error
      if errors.Is(err, gorm.ErrRecordNotFound) {
        areaPrice = scsave.FetchAreaPrice(v.City)
        cachingAreaPrice := scsave.AreaPrice{}
        cachingAreaPrice.Name = v.City
        cachingAreaPrice.Price = areaPrice
        DB.FirstOrCreate(&cachingAreaPrice, "name = ?", cachingAreaPrice.Name)
      } else {
        areaPrice = cachedAreaPrice.Price
      }
    }
    v.AreaPrice = areaPrice
    v.EstimatedBuildingPrice = int64(v.Price) * 10000 - int64(v.LandArea) * areaPrice
    var old scsave.Property
    err = DB.First(&old, "url = ?", v.Url).Error
    if errors.Is(err, gorm.ErrRecordNotFound) {
      DB.Create(&v)
    } else {
      if v.Price != old.Price {
        v.Revision = old.Revision + 1
        DB.Create(&v)
        DB.Delete(&old)
      }
    }
  }
  DB.Model(&job).Update("type", jobType)
  DB.Model(&job).Update("state", string(scsave.JobStateCompleted))
  var count int64
  DB.Model(&scsave.Property{}).Count(&count)
  fmt.Printf("count: %d\n", count)
}
