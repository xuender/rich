package rich

import (
	"errors"
	"sort"

	utils "github.com/xuender/go-utils"
)

// Days 日期
type Days []string

// Add 增加
func (d *Days) Add(day string) bool {
	if d.Includes(day) {
		return false
	}
	*d = append(*d, day)
	sort.Sort(sort.Reverse(sort.StringSlice(*d)))
	return true
}

// Includes 包含
func (d *Days) Includes(day string) bool {
	for _, v := range *d {
		if v == day {
			return true
		}
	}
	return false
}

// 日帐增加
func (w *Web) dayAdd(t Trade) error {
	day := t.Ca.Format("2006-01-02")
	w.days.Add(day)
	w.Put(DaysKey, w.days)
	idsKey := []byte(day)
	ids := []utils.ID{}
	w.Get(idsKey, &ids)
	ids = append(ids, t.ID)
	return w.Put(idsKey, ids)
}

// 日帐减少
func (w *Web) dayDel(t Trade) error {
	day := t.Ca.Format("2006-01-02")
	idsKey := []byte(day)
	ids := []utils.ID{}
	w.Get(idsKey, &ids)
	for i, id := range ids {
		if id.Equal(t.ID) {
			ids = append(ids[:i], ids[i+1:]...)
			w.Put(idsKey, ids)
			if len(ids) == 0 {
				for f, d := range w.days {
					if d == day {
						w.days = append(w.days[:f], w.days[f+1:]...)
						w.Put(DaysKey, w.days)
					}
				}
			}
		}
	}
	return errors.New("帐目未找到")
}
