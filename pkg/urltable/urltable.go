package urltable

import (
	"net/http"
	"strings"

	"github.com/xinliangnote/go-gin-api/pkg/errors"
)

const (
	empty      = ""
	fuzzy      = "*"
	omitted    = "**"
	delimiter  = "/"
	methodView = "VIEW"
)

// parse and validate pattern
func parse(pattern string) ([]string, error) {
	const format = "[get, post, put, patch, delete, view]/{a-Z}+/{*}+/{**}"

	// 去除首尾空格和开头的/
	if pattern = strings.TrimLeft(strings.TrimSpace(pattern), delimiter); pattern == "" {
		return nil, errors.Errorf("pattern illegal, should in format of %s", format)
	}

	// 按 / 分割路径
	// 至少需要 2 个部分（方法 + 至少一个路径段）
	paths := strings.Split(pattern, delimiter)
	if len(paths) < 2 {
		return nil, errors.Errorf("pattern illegal, should in format of %s", format)
	}

	// 去除每个路径段的首尾空格
	for i := range paths {
		paths[i] = strings.TrimSpace(paths[i])
	}

	// 防止无效的通配符用法
	// likes get/ get/* get/**
	if len(paths) == 2 && (paths[1] == empty || paths[1] == fuzzy || paths[1] == omitted) {
		return nil, errors.New("illegal wildcard")
	}

	// 将方法名转为大写
	// 验证是否为支持的 HTTP 方法
	// 支持自定义的 VIEW 方法
	switch paths[0] = strings.ToUpper(paths[0]); paths[0] {
	case http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		methodView:
	default:
		return nil, errors.Errorf("only supports [%s %s %s %s %s %s]",
			http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, methodView)
	}

	// 空路径验证：防止中间出现空路径段（如 GET/api//user）
	// 省略符位置验证：** 必须是最后一个路径段
	for k := 1; k < len(paths); k++ {
		if paths[k] == empty && k+1 != len(paths) {
			return nil, errors.New("pattern contains illegal empty path")
		}

		if paths[k] == omitted && k+1 != len(paths) {
			return nil, errors.New("pattern contains illegal omitted path")
		}
	}

	// 返回解析后的路径段数组，如 ["GET", "api", "v1", "user"]
	return paths, nil
}

// Format pattern
func Format(pattern string) (string, error) {
	paths, err := parse(pattern)
	if err != nil {
		return "", err
	}

	return strings.Join(paths, delimiter), nil
}

// 前缀树
/* 根节点 (section)
├── "api" (section)
│   ├── "v1" (section)
│   │   ├── "user" (section, leaf=true)  # GET/api/v1/user
│   │   └── "order" (section, leaf=true) # GET/api/v1/order
│   └── "v2" (section)
│       └── "user" (section, leaf=true)  # GET/api/v2/user
├── "*" (section, leaf=true)             # GET/* (单段通配符,如/api/"*"/info)
└── "**" (section, leaf=true)            # GET/** (多段省略符，如/api/"**")
*/
type section struct {
	leaf    bool
	mapping map[string]*section
}

func newSection() *section {
	return &section{mapping: make(map[string]*section)}
}

// Table a (thread unsafe) table to store urls
type Table struct {
	size int
	root *section
}

// NewTable create a table
func NewTable() *Table {
	return &Table{root: newSection()}
}

// Size contains how many urls
func (t *Table) Size() int {
	return t.size
}

// Append pattern
// // 插入路由模式
// table.Append("GET/api/v1/user")     // leaf=true 在 "user" 节点
// table.Append("GET/api/v1/order")    // leaf=true 在 "order" 节点
// table.Append("GET/api/*/info")      // leaf=true 在 "*" 节点
// table.Append("GET/api/**")          // leaf=true 在 "**" 节点
func (t *Table) Append(pattern string) error {
	paths, err := parse(pattern)
	if err != nil {
		return err
	}

	insert := false
	root := t.root
	for i, path := range paths {
		if (path == fuzzy && root.mapping[omitted] != nil) ||
			(path == omitted && root.mapping[fuzzy] != nil) ||
			(path != omitted && root.mapping[omitted] != nil) {
			return errors.Errorf("conflict at %s", strings.Join(paths[:i], delimiter))
		}

		next := root.mapping[path]
		if next == nil {
			next = newSection()
			root.mapping[path] = next
			insert = true
		}
		root = next
	}

	if insert {
		t.size++
	}

	root.leaf = true
	return nil
}

// Mapping url to pattern
// // 查询匹配
// pattern, _ := table.Mapping("/api/v1/user")     // 返回 "GET/api/v1/user"
// pattern, _ := table.Mapping("/api/v2/info")     // 返回 "GET/api/*/info"
// pattern, _ := table.Mapping("/api/v1/user/info") // 返回 "GET/api/**"
func (t *Table) Mapping(url string) (string, error) {
	paths, err := parse(url)
	if err != nil {
		return "", err
	}

	pattern := make([]string, 0, len(paths))

	root := t.root
	for _, path := range paths {
		next := root.mapping[path]
		if next == nil {
			nextFuzzy := root.mapping[fuzzy]
			nextOmitted := root.mapping[omitted]
			if nextFuzzy == nil && nextOmitted == nil {
				return "", nil
			}

			if nextOmitted != nil {
				pattern = append(pattern, omitted)
				return strings.Join(pattern, delimiter), nil
			}

			next = nextFuzzy
			path = fuzzy
		}

		root = next
		pattern = append(pattern, path)
	}

	if root.leaf {
		return strings.Join(pattern, delimiter), nil
	}
	return "", nil
}
