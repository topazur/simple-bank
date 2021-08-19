package api

import (
	"fmt"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/topaz-h/go-simple-bank/util"
)

/**
 * 因为Gin的HTML渲染就是基于html/template实现的
 */
func (server *Server) htmlTemplate(ctx *gin.Context) {
	ctx.Status(200)
	const templateText = `微信公众号: {{printf "%s" .}}`
	tmpl, err := template.New("htmlTest").Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	tmpl.Execute(ctx.Writer, "hs")
}

func (server *Server) uploadFile(ctx *gin.Context) {
	//获取普通文本
	name := ctx.PostForm("name")
	var err error
	var file *multipart.FileHeader
	// 获取文件(注意这个地方的file要和html模板中的name一致)
	if file, err = ctx.FormFile("file"); err != nil {
		fmt.Println("接收的数据", name, file.Filename)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "读取文件失败",
		})
		return
	}

	// 获取文件的后缀名
	extstring := path.Ext(file.Filename)
	// 根据当前时间鹾生成一个新的文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	// 新的文件名
	fileName := fileNameStr + extstring
	// 保存上传文件
	filePath := filepath.Join(util.MkdirFolder("tmp"), "/", fileName)
	ctx.SaveUploadedFile(file, filePath)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

func (server *Server) uploadFiles(ctx *gin.Context) {
	var form *multipart.Form
	var err error
	if form, err = ctx.MultipartForm(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "读取文件失败",
		})
		return
	}

	//1.获取文件
	files := form.File["file"]

	//2.循环全部的文件
	for _, file := range files {
		// 获取文件的后缀名
		// extstring := path.Ext(file.Filename)
		// 3.根据时间戳生成文件名
		fileNameInt := time.Now().Unix()
		fileNameStr := strconv.FormatInt(fileNameInt, 10)
		//4.新的文件名(如果是同时上传多张图片的时候就会同名，因此这里使用时间鹾加文件名方式)
		fileName := fileNameStr + file.Filename
		//5.保存上传文件
		filePath := filepath.Join(util.MkdirFolder("tmp"), "/", fileName)
		ctx.SaveUploadedFile(file, filePath)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// 微信内部 浏览器打开
// iphone 微信内 诱导 长按图片保存，其他的都不好使
// 用UC，QQ浏览器下载，safari不支持下载（除了图片）。

// func (server *Server) download(ctx *gin.Context) {
// 	name, ok := ctx.GetPostForm("name")
// 	fmt.Println(name, ok)
// 	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "filename"))
// 	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
// 	ctx.File(download_file)
// }

func (server *Server) downPng(c *gin.Context) {
	// https://cn.vuejs.org/images/logo.svg
	response, err := http.Get("https://cn.vuejs.org/images/dcloud2.png")
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

func (server *Server) downZip(c *gin.Context) {
	response, err := http.Get("http://nginx.org/download/nginx-1.20.1.zip")
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="nginx-1.20.1.zip"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
