package controllers

import (
	"github.com/revel/revel"
	// "wuhuaguo.com/whgv01/app/routes"
	"wuhuaguo.com/whgv01/app/models"
)

type Product struct {
	Application
}

func (c Product) Detail(id int) revel.Result {
	//查询商品数据
	results, err := c.Txn.Select(models.Item{}, `select * from item where Id = ?`, id)
	if err != nil {
		panic(err)
	}
	if len(results) == 0 {
		revel.INFO.Println("%s item not exists", id)
		return nil
	}

	ir := results[0].(*models.Item).ItemToRender()

	revel.INFO.Println("%t", ir)

	return c.Render(ir)
}
