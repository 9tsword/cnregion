// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package cnregion

import (
	"os"
	"strings"
	"testing"

	"github.com/issue9/assert/v4"

	"github.com/issue9/cnregion/v2/id"
)

func TestDB_Search(t *testing.T) {
	a := assert.New(t, false)

	rs := obj.Search(&Options{Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].name, "合肥")

	rs = obj.Search(&Options{Parent: "340000000000", Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].name, "合肥")

	rs = obj.Search(&Options{Parent: "000000000000", Text: "合肥"})
	a.Equal(1, len(rs)).
		Equal(rs[0].name, "合肥")

	// 限定 level 只能是省以及 parent 为 34 开头
	rs = obj.Search(&Options{Parent: "340000000000", Level: id.Province, Text: "合肥"})
	a.Equal(0, len(rs))

	// 未限定 parent 且 level 正确
	rs = obj.Search(&Options{Level: id.City, Text: "合肥"})
	a.Equal(1, len(rs))

	rs = obj.Search(&Options{Level: id.City, Text: "湖"})
	a.Equal(2, len(rs))

	rs = obj.Search(&Options{Level: id.City, Parent: "340000000000", Text: "湖"})
	a.Equal(2, len(rs))

	// parent = 浙江
	rs = obj.Search(&Options{Parent: "330000000000", Text: "合肥"})
	a.Equal(0, len(rs))

	// parent 不存在
	rs = obj.Search(&Options{Parent: "110000000000", Text: "合肥"})
	a.Equal(0, len(rs))

	// 只有 Level
	rs = obj.Search(&Options{Level: id.City})
	a.Equal(4, len(rs))
	for _, r := range rs {
		a.True(strings.HasSuffix(r.fullID, "00000000"))
	}

	// 只有 Level
	rs = obj.Search(&Options{Level: id.City + id.Town})
	a.Equal(4, len(rs))

	// 只有 Level
	rs = obj.Search(&Options{Level: id.City + id.Province})
	a.Equal(6, len(rs))
}

func TestDB_SearchWithData(t *testing.T) {
	a := assert.New(t, false)

	obj, err := LoadFS(os.DirFS("./data"), "regions.db", "-", true)
	a.NotError(err).NotNil(obj)
	got := obj.Search(&Options{Text: "温州"})
	a.NotEmpty(got)

	// Level 不匹配
	got = obj.Search(&Options{Text: "温州", Level: id.Province})
	a.Empty(got)
}
