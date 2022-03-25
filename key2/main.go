package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

var db *gorm.DB

func init() {
	//创建一个数据库的连接
	var err error
	db, err = gorm.Open("mysql", "root:root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	//迁移the schema
	db.AutoMigrate(&kodo{})
}

type (
	kodo struct {
		gorm.Model
		Task  string `json:"task"`
		Judge int    `json:"judge"`
	}

	kodos struct {
		ID    uint   `json:"id"`
		Task  string `json:"task "`
		Judge bool   `json:"judge"`
	}
)

func main() {

	r := gin.Default()
	//添加一条任务
	r.POST("/", func(c *gin.Context) {

		judge, _ := strconv.Atoi(c.PostForm("judge"))
		k := kodo{Task: c.PostForm("task"), Judge: judge}
		db.Save(&k)
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "struct": k})
	})
	//查询全部任务
	r.GET("/", func(c *gin.Context) {
		var a []kodo
		var b kodos
		db.Find(&a)
		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if item.Judge == 1 {
				b.Judge = true
			} else {
				b.Judge = false
			}
			if e := c.ShouldBind(&b); e == nil {
				c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			}
		}

	})
	//查询全部已完成任务
	r.GET("/t", func(c *gin.Context) {
		var a []kodo
		var b kodos
		db.Find(&a)
		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if item.Judge == 1 {
				b.Judge = true
			} else {
				b.Judge = false
			}
			if b.Judge == true {

				if e := c.ShouldBind(&b); e == nil {
					c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				}
			}
		}

	})
	//查询全部未完成任务
	r.GET("/f", func(c *gin.Context) {
		var a []kodo
		var b kodos
		db.Find(&a)
		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if item.Judge == 1 {
				b.Judge = true
			} else {
				b.Judge = false
			}
			if b.Judge == false {

				if e := c.ShouldBind(&b); e == nil {
					c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				}
			}
		}

	})
	//关键词查询任务，任务为完成
	r.GET("/:name/t", func(c *gin.Context) {
		key := c.Param("name")
		var a []kodo
		var b kodos

		db.Where("task LIKE ? AND judge = ?", "%"+key+"%", "1").Find(&a)
		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if item.Judge == 1 {
				b.Judge = true
			} else {
				b.Judge = false
			}
			if b.Judge == true {

				if e := c.ShouldBind(&b); e == nil {
					c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				}
			}
		}

	})
	//关键词查询任务，任务为未完成
	r.GET("/:name/f", func(c *gin.Context) {
		key := c.Param("name")
		var a []kodo
		var b kodos

		db.Where("task LIKE ? AND judge = ?", "%"+key+"%", "false").Find(&a)

		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if item.Judge == 1 {
				b.Judge = true
			} else {
				b.Judge = false
			}
			if b.Judge == false {

				if e := c.ShouldBind(&b); e == nil {
					c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				}
			}
		}

	})
	//单个任务修改完成状态
	r.PUT("/:id", func(c *gin.Context) {
		key := c.Param("id")
		var b kodos
		var a kodo
		var d uint
		db.First(&a, key)
		b.Task = a.Task
		b.ID = a.ID
		key2 := a.Judge
		if key2 != 1 {
			b.Judge = true
			d = 1
		} else {
			b.Judge = false
			d = 0
		}
		db.Model(&a).Update("judge", d)
		if e := c.ShouldBind(&b); e == nil {
			c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		}
	})
	//修改全部任务状态
	r.PUT("/put/:judge", func(c *gin.Context) {
		key, _ := strconv.Atoi(c.Param("judge"))
		var b kodos
		var a []kodo
		var d uint
		db.Find(&a)
		for _, item := range a {
			b.ID = item.ID
			b.Task = item.Task

			if key == 1 {
				b.Judge = true
				d = 1
			} else {
				b.Judge = false
				d = 0
			}
			db.Model(&a).Update("judge", d)
			if e := c.ShouldBind(&b); e == nil {
				c.JSON(201, gin.H{"id": b.ID, "task": b.Task, "judge": b.Judge})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
			}
		}

	})
	//删除
	r.DELETE("/:id", func(c *gin.Context) {
		key := c.Param("id")
		var a kodo
		db.First(&a, key)
		db.Delete(&a)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
	})

	r.Run()
}
