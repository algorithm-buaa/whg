package controllers

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/revel/modules/db/app"
	r "github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"wuhuaguo.com/whgv01/app/models"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTableWithName(models.User{}, "user").SetKeys(true, "UserId")

	t.ColMap("Password").Transient = true

	t = Dbm.AddTableWithName(models.Item{}, "item").SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"Name":       64,
		"Intro":      1024,
		"Spec":       2048,
		"SmImage":    512,
		"FirstImage": 512,
		"Images":     4096,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTablesIfNotExists()

	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("123456"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "游客", "游客", "123456", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	r.INFO.Println("gorp's begin")
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
