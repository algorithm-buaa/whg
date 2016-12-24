package models

import (
	"github.com/revel/revel"
	"strings"
)

type Item struct {
	Id   int
	Name string
	//卖方id
	SallerId int
	Price    float32
	// 商品介绍， json格式 描述商品信息
	Intro string
	// 商品规格
	Spec string
	//小图片
	SmImage string
	//首页图片
	FirstImage string
	//介绍图片  多张图片
	Images     string
	DisplayPos int
}

type ItemToRender struct {
	item     *Item
	Images   []string
	SmImages []string
}

func (i *ItemToRender) Name() string {
	return i.item.Name
}

func (i *ItemToRender) Intro() string {
	return i.item.Intro
}

func (i *Item) Validate(v *revel.Validation) {
	v.Check(i.Name,
		revel.Required{},
		revel.MaxSize{64},
		revel.MinSize{2},
	)

	v.Check(i.Intro,
		revel.Required{},
		revel.MinSize{32},
	)

	v.Check(i.Spec,
		revel.Required{},
		revel.MinSize{32},
	)
}

func (i *Item) ItemToRender() *ItemToRender {
	images := strings.Split(i.FirstImage, ",")
	smImage := strings.Split(i.SmImage, ",")
	ir := &ItemToRender{item: i, Images: images, SmImages: smImage}
	return ir
}
