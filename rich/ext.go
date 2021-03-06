package rich

import (
	"net/http"

	"github.com/labstack/echo"
)

// Ext 扩展定义
type Ext struct {
	Key   string `json:"key"`   // 键
	Value string `json:"value"` // 值
}

// 扩展路由
func (w *Web) extsRoute(c *echo.Group) {
	c.GET("", w.extsGet)    // 全部扩展
	c.GET("/:id", w.extGet) // 获取扩展
	c.PUT("/:id", w.extPut) // 修改扩展
}

var extKeys = []string{
	"E-C", // 客户扩展数据
	"E-I", // 商品扩展数据
}

// 获取全部扩展定义
func (w *Web) extsGet(c echo.Context) error {
	ret := make(map[string][]Ext)
	for _, key := range extKeys {
		ret[key] = w.getExtByID(key)
	}
	return c.JSON(http.StatusOK, ret)
}

// 根据id获取扩展定义
func (w *Web) getExtByID(id string) []Ext {
	ret := []Ext{}
	w.Get([]byte(id), &ret)
	return ret
}

// 获取扩展定义
func (w *Web) extGet(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, w.getExtByID(id))
}

// 扩展定义修改
func (w *Web) extPut(c echo.Context) error {
	id := c.Param("id")
	ext := w.getExtByID(id)
	idBs := []byte(id)
	w.Get(idBs, &ext)
	if err := c.Bind(&ext); err != nil {
		return err
	}
	w.Put(idBs, ext)
	return c.JSON(http.StatusOK, ext)
}
