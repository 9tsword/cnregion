// SPDX-FileCopyrightText: 2021-2024 caixw
//
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/issue9/cnregion/v2"
	"github.com/issue9/cnregion/v2/version"
)

func build(dataDir, output string, years ...int) error {
	d := cnregion.NewDB()

	if len(years) == 0 {
		years = version.All()
	}
	for _, year := range years {
		if err := buildYear(d, dataDir, year); err != nil {
			return err
		}
	}

	return d.Dump(output, true)
}

func buildYear(d *cnregion.DB, dataDir string, year int) error {
	fmt.Printf("\n添加 %d 的数据\n", year)
	if !d.AddVersion(year) {
		fmt.Printf("已经存在该年份 %d 的数据\n\n", year)
		return nil
	}

	y := strconv.Itoa(year)
	dataDir = filepath.Join(dataDir, y)

	return filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name()[0] == '.' { // 忽略隐藏文件
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		s := bufio.NewScanner(bytes.NewBuffer(data))
		s.Split(bufio.ScanLines)
		for s.Scan() {
			txt := s.Text()
			values := strings.Split(txt, "\t")
			if len(values) != 2 {
				return fmt.Errorf("无效的格式，位于 %s:%s", path, txt)
			}
			id, name := values[0], values[1]

			if err := d.AddItem(id, name, year); err != nil {
				return err
			}
		}

		return nil
	})
}
