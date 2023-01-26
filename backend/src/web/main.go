package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"kouya17/scsave"
)

var DB *gorm.DB

func main() {
	var err error
	DB, err = open("host=db user=postgres password=postgres dbname=scsave port=5432 sslmode=disable TimeZone=Asia/Tokyo", 30)
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

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "backend server is running")
	})
	e.GET("/scrape-properties-jobs", getScrapePropertyJobs)
	e.POST("/scrape-properties-jobs", createScrapePropertyJob)
	e.GET("/properties", getProperties)
	e.PATCH("/properties/:id", patchProperties)
	e.GET("/properties/count", getPropertiesCount)
	e.Logger.Fatal(e.Start(":8000"))
}

func getScrapePropertyJobs(c echo.Context) error {
	var jobs []scsave.ScrapePropertiesJob
	DB.Order("id desc").Find(&jobs)
	return c.JSON(http.StatusOK, jobs)
}

func createScrapePropertyJob(c echo.Context) error {
	j := new(scsave.ScrapePropertiesJob)
	if err := c.Bind(j); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	DB.Create(j)

	return c.JSON(http.StatusOK, j)
}

func getProperties(c echo.Context) error {
	city := c.QueryParam("city")
	order := c.QueryParam("order")
	properties := new([]scsave.Property)
	var count int64
	q := DB.Where("city LIKE ?", "%"+city+"%")
	if c.QueryParam("min-job") != "" {
		minJobId, _ := strconv.ParseUint(c.QueryParam("min-job"), 10, 64)
		q = q.Where("job_id >= ?", minJobId)
	}
	if c.QueryParam("max-job") != "" {
		maxJobId, _ := strconv.ParseUint(c.QueryParam("max-job"), 10, 64)
		q = q.Where("job_id <= ?", maxJobId)
	}
	q.Find(&properties).Count(&count)
	q.Scopes(Paginate(c.Request())).Order(order).Find(&properties)
	res := scsave.ResponseProperies{Count: count, Properties: *properties}
	return c.JSON(http.StatusOK, res)
}

func patchProperties(c echo.Context) error {
	id := c.Param("id")
	property := new(scsave.Property)
	DB.Find(&property, "id = ?", id)
	if err := c.Bind(property); err != nil {
		return err
	}
	DB.Save(&property)
	return c.JSON(http.StatusOK, property)
}

func getPropertiesCount(c echo.Context) error {
	city := c.QueryParam("city")
	properties := new([]scsave.Property)
	var count int64
	q := DB.Where("city LIKE ?", "%"+city+"%")
	if c.QueryParam("min-job") != "" {
		minJobId, _ := strconv.ParseUint(c.QueryParam("min-job"), 10, 64)
		q = q.Where("job_id >= ?", minJobId)
	}
	if c.QueryParam("max-job") != "" {
		maxJobId, _ := strconv.ParseUint(c.QueryParam("max-job"), 10, 64)
		q = q.Where("job_id <= ?", maxJobId)
	}
	q.Find(&properties).Count(&count)
	return c.JSON(http.StatusOK, count)
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func open(dsn string, count uint) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		if count == 0 {
			return nil, fmt.Errorf("Retry count over")
		}
		time.Sleep(time.Second)
		count--
		return open(dsn, count)
	}
	return db, nil
}
