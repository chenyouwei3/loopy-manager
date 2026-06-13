package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type RunningLog struct {
	Name  string `json:"name"`  // 文件名，如 app-2025-09-27.log
	Path  string `json:"path"`  // 完整路径
	Size  int64  `json:"size"`  // 文件大小（字节）
	MTime string `json:"mtime"` // 修改时间（格式化）
	IsLog bool   `json:"isLog"` // 是否是 .log 文件
}

func (r *RunningLog) GetList(path string, skip, limit int, startTime, endTime string) ([]RunningLog, int, error) {
	var logs []RunningLog

	// 检查路径是否存在
	info, err := os.Stat(path)
	if err != nil {
		return nil, 0, err
	}
	if !info.IsDir() {
		return nil, 0, os.ErrInvalid
	}

	// 解析时间参数
	var start, end time.Time
	if startTime != "" {
		start, err = time.Parse("2006-01-02 15:04:05", startTime)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid startTime format, expect yyyy-MM-dd HH:mm:ss")
		}
	}
	if endTime != "" {
		end, err = time.Parse("2006-01-02 15:04:05", endTime)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid endTime format, expect yyyy-MM-dd HH:mm:ss")
		}
	}

	// 读取目录下所有文件
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, 0, err
	}

	// 过滤出 .log 文件，并收集信息
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只保留 .log 结尾的文件
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".log") {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())

		// 获取文件信息
		fi, err := entry.Info()
		if err != nil {
			continue // 跳过出错的文件
		}

		// 根据时间筛选
		mtime := fi.ModTime()
		if !start.IsZero() && mtime.Before(start) {
			continue
		}
		if !end.IsZero() && mtime.After(end) {
			continue
		}

		logs = append(logs, RunningLog{
			Name:  entry.Name(),
			Path:  fullPath,
			Size:  fi.Size(),
			MTime: mtime.Format("2006-01-02 15:04:05"),
			IsLog: true,
		})
	}

	// 按修改时间倒序排序（最新的在前）
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].MTime > logs[j].MTime
	})

	// 计算总数
	total := len(logs)

	endIdx := skip + limit
	if endIdx > len(logs) {
		endIdx = len(logs)
	}
	if skip > len(logs) { // 防止越界
		return []RunningLog{}, total, nil
	}

	return logs[skip:endIdx], total, nil
}
