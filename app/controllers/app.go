package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"wuhuaguo.com/whgv01/app/models"
	"wuhuaguo.com/whgv01/app/routes"
)

type Application struct {
	GorpController
}

func (c Application) CheckUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		revel.INFO.Println("check User  user: nil")
		return c.Redirect("https://open.weixin.qq.com/connect/oauth2/authorize?appid=wx520c15f417810387&redirect_uri=https%3A%2F%2Fchong.qq.com%2Fphp%2Findex.php%3Fd%3D%26c%3DwxAdapter%26m%3DmobileDeal%26showwxpaytitle%3D1%26vb2ctag%3D4_2030_5_1194_60&response_type=code&scope=snsapi_base&state=123#wechat_redirect")
	}
	return nil
}

func (c Application) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.RenderArgs["user"] = user
		revel.INFO.Println("add user " + user.Name)
	}
	return nil
}

func (c Application) wxId() *models.User {
	if userid, ok := c.Session["wxid"]; ok {
		return c.getUser(userid)
	}
	return nil
}

func (c Application) connected() *models.User {
	if c.RenderArgs["user"] != nil {
		return c.RenderArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	//返回游客账号
	//判断浏览器
	ua := c.Request.UserAgent()
	revel.INFO.Println(ua)
	if strings.Contains(c.Request.UserAgent(), "micromessenger") {
		return nil
	} else {
		return c.getUser("游客")
	}
}

func (c Application) getUser(username string) *models.User {
	users, err := c.Txn.Select(models.User{}, `select * from user where Username = ?`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c Application) Index() revel.Result {
	user := c.connected()

	revel.INFO.Println("username: %q", user.Username)
	c.Session["user"] = user.Username

	//获取boxlist 商品
	irs, err := c.boxListItems()
	if err != nil {
		revel.ERROR.Println(err)
	}
	return c.Render(user, irs)
}

func (c Application) Register() revel.Result {
	return c.Render()
}

func (c Application) SaveUser(user models.User, verifyPassword string) revel.Result {
	c.Validation.Required(verifyPassword)
	c.Validation.Required(verifyPassword == user.Password).
		Message("Password does not match")
	user.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Application.Register())
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword(
		[]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		panic(err)
	}

	c.Session["user"] = user.Username
	c.Flash.Success("Welcome, " + user.Name)
	return c.Redirect(routes.Application.Index())
}

func (c *Application) LoginIndex() revel.Result {
	return c.Render()
}

func (c Application) Login(username, password string, remember bool) revel.Result {
	user := c.getUser(username)
	if user != nil {
		err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		if err == nil {
			c.Session["user"] = username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + username)
			return c.Redirect(routes.Application.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("登录失败，账号或密码错误！")
	return c.Redirect(routes.Application.LoginIndex())
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Index())
}

//在首页列表展示区域的商品
func (c Application) boxListItems() ([]*models.ItemToRender, error) {
	results, err := c.Txn.Select(models.Item{}, `select * from item where DisplayPos = ?`, 1)
	if err != err {
		return nil, err
	}
	var irs []*models.ItemToRender
	for _, r := range results {
		i := r.(*models.Item).ItemToRender()
		irs = append(irs, i)
	}
	return irs, nil
}
