package rich

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/xuender/go-utils"
)

// 账户路由
func (w *Web) profileRoute(c *echo.Group) {
	c.GET("", w.profileGet)              // 账户信息
	c.PATCH("", w.profilePatch)          // 修改账户信息
	c.PATCH("/pass", w.profilePassPatch) // 修改密码
}

// profile 获取当前账户
func (w *Web) profile(c echo.Context) (u User, err error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := new(utils.ID)
	if err = id.Parse(claims["id"].(string)); err != nil {
		return
	}
	w.Get(id[:], &u)
	return
}

// profileGet 当前账户信息
func (w *Web) profileGet(c echo.Context) error {
	u, err := w.profile(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, u)
}

// profilePassPatch 修改密码
func (w *Web) profilePassPatch(c echo.Context) error {
	old := c.QueryParam("old")
	if old == "" {
		return errors.New("原密码不能为空")
	}
	pass := c.QueryParam("pass")
	if pass == "" {
		return errors.New("新密码不能为空")
	}
	u, err := w.profile(c)
	if err != nil {
		return err
	}
	if !bytes.Equal(u.Pass, Pass(old)) {
		return errors.New("原密码错误")
	}
	w.UpdatePass(&u, pass)
	t, err := w.Signed(u.ID.String(), u.Pass)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

// profilePatch 修改账户信息
func (w *Web) profilePatch(c echo.Context) error {
	u, err := w.profile(c)
	if err != nil {
		return err
	}
	users := w.users()
	users = users.filter(func(o *User) bool { return !o.ID.Equal(u.ID) })
	n := User{}
	err = w.Bind(c, &n)
	if err != nil {
		return err
	}
	if n.Name != "" {
		for _, o := range users {
			if o.Name == n.Name {
				return errors.New("名称重复")
			}
		}
		u.Name = n.Name
	}
	if n.Phone != "" {
		for _, o := range users {
			if o.Phone == n.Phone {
				return errors.New("手机号重复")
			}
		}
		u.Phone = n.Phone
	}
	w.Put(u.ID[:], u)
	return c.JSON(http.StatusOK, u)
}
