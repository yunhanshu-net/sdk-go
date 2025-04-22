package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// 需要排除的目录列表
var excludeDirs = map[string]bool{
	".git":         true,
	"soft":         true,
	"vendor":       true,
	"node_modules": true,
}

// 支持的代码文件扩展名及其对应语言标记
var codeLang = map[string]string{
	".go":   "go",
	".mod":  "mod",
	".sum":  "sum",
	".toml": "toml",
	".yaml": "yaml",
	".yml":  "yaml",
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <project_root> <output.md>")
		return
	}

	root := os.Args[1]
	output := os.Args[2]

	var buf bytes.Buffer
	buf.WriteString("# Project Code Structure\n\n")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过排除目录
		if info.IsDir() && excludeDirs[info.Name()] {
			return filepath.SkipDir
		}

		// 只处理普通文件
		if !info.Mode().IsRegular() {
			return nil
		}

		// 获取相对路径
		relPath, _ := filepath.Rel(root, path)
		dir, file := filepath.Split(relPath)

		// 写入文件头
		buf.WriteString(fmt.Sprintf("## %s\n", relPath))
		buf.WriteString(fmt.Sprintf("**Path:** `%s`  \n", dir))
		buf.WriteString(fmt.Sprintf("**File:** `%s`\n\n", file))

		// 读取文件内容
		content, err := ioutil.ReadFile(path)
		if err != nil {
			buf.WriteString("```\nError reading file\n```\n\n")
			return nil
		}

		// 获取代码语言标记
		ext := filepath.Ext(file)
		lang, ok := codeLang[ext]
		if !ok {
			lang = "text"
		}

		// 写入代码块
		buf.WriteString(fmt.Sprintf("```%s\n", lang))
		buf.Write(content)
		if !bytes.HasSuffix(content, []byte("\n")) {
			buf.WriteByte('\n')
		}
		buf.WriteString("```\n\n")

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
		return
	}

	// 写入输出文件
	if err := ioutil.WriteFile(output, buf.Bytes(), 0644); err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated documentation: %s\n", output)
}
