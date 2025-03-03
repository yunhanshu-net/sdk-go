package contenttype

import "strings"

var (
	// TextPlainCharsetUtf8 文本类型
	TextPlainCharsetUtf8      = "text/plain; charset=utf-8"      // 普通文本
	TextHtmlCharsetUtf8       = "text/html; charset=utf-8"       // HTML 文档
	TextCssCharsetUtf8        = "text/css; charset=utf-8"        // CSS 样式表
	TextCsvCharsetUtf8        = "text/csv; charset=utf-8"        // 逗号分隔值
	TextJavascriptCharsetUtf8 = "text/javascript; charset=utf-8" // JavaScript 文件

	// ApplicationJsonCharsetUtf8 应用类型
	ApplicationJsonCharsetUtf8               = "application/json; charset=utf-8"                  // JSON 格式
	ApplicationXmlCharsetUtf8                = "application/xml; charset=utf-8"                   // XML 格式
	ApplicationZip                           = "application/zip"                                  // ZIP 压缩文件
	ApplicationPdf                           = "application/pdf"                                  // PDF 文档
	ApplicationOctetStream                   = "application/octet-stream"                         // 二进制流
	ApplicationXWwwFormUrlencodedCharsetUtf8 = "application/x-www-form-urlencoded; charset=utf-8" // 表单数据
	ApplicationJavascriptCharsetUtf8         = "application/javascript; charset=utf-8"            // JavaScript 文件

	// ImageJpeg 图像类型
	ImageJpeg   = "image/jpeg"    // JPEG 图像
	ImagePng    = "image/png"     // PNG 图像
	ImageGif    = "image/gif"     // GIF 动画
	ImageBmp    = "image/bmp"     // BMP 图像
	ImageSvgXml = "image/svg+xml" // SVG 矢量图

	// AudioMpeg 音频类型
	AudioMpeg = "audio/mpeg" // MPEG 音频
	AudioWav  = "audio/wav"  // WAV 音频
	AudioOgg  = "audio/ogg"  // OGG 音频

	// VideoMp4 视频类型
	VideoMp4  = "video/mp4"  // MP4 视频
	VideoMpeg = "video/mpeg" // MPEG 视频
	VideoOgg  = "video/ogg"  // OGG 视频
	VideoWebm = "video/webm" // WebM 视频
)

type HttpContentType struct {
	Type        string // MIME 类型
	SubType     string // 子类型
	Description string // 描述
	ContentType string // 全称
	IsFile      bool   //是否是文件
}

// GetContentType 获取contentType 类型
func GetContentType(contentType string) string { //application/json; charset=utf-8
	s := strings.Split(contentType, ";")[0] //application/json

	contentType = strings.ToLower(s) //application/json
	for _, httpContentType := range ContentTypes {
		if strings.Split(httpContentType.ContentType, ";")[0] == contentType { //application/json
			return httpContentType.ContentType //application/json; charset=utf-8
		}
	}
	return "text/plain; charset=utf-8"

}

// ContentTypes 定义常见的 Content-Type
var ContentTypes = []HttpContentType{
	// 文本类型
	{Type: "text", SubType: "plain", Description: "普通文本", ContentType: "text/plain; charset=utf-8"},
	{Type: "text", SubType: "html", Description: "HTML 文档", ContentType: "text/html; charset=utf-8"},
	{Type: "text", SubType: "css", Description: "CSS 样式表", ContentType: "text/css; charset=utf-8"},
	{Type: "text", SubType: "csv", Description: "逗号分隔值", ContentType: "text/csv; charset=utf-8"},
	{Type: "text", SubType: "javascript", Description: "JavaScript 文件", ContentType: "text/javascript; charset=utf-8"},

	// 应用类型
	{Type: "application", SubType: "json", Description: "JSON 格式", ContentType: "application/json; charset=utf-8"},
	{Type: "application", SubType: "xml", Description: "XML 格式", ContentType: "application/xml; charset=utf-8"},
	{Type: "application", SubType: "zip", Description: "ZIP 压缩文件", ContentType: "application/zip", IsFile: true},
	{Type: "application", SubType: "pdf", Description: "PDF 文档", ContentType: "application/pdf", IsFile: true},
	{Type: "application", SubType: "octet-stream", Description: "二进制流", ContentType: "application/octet-stream", IsFile: true},
	{Type: "application", SubType: "x-www-form-urlencoded", Description: "表单数据", ContentType: "application/x-www-form-urlencoded; charset=utf-8"},
	{Type: "application", SubType: "javascript", Description: "JavaScript 文件", ContentType: "application/javascript; charset=utf-8"},

	// 图像类型
	{Type: "image", SubType: "jpeg", Description: "JPEG 图像", ContentType: "image/jpeg"},
	{Type: "image", SubType: "png", Description: "PNG 图像", ContentType: "image/png"},
	{Type: "image", SubType: "gif", Description: "GIF 动画", ContentType: "image/gif"},
	{Type: "image", SubType: "bmp", Description: "BMP 图像", ContentType: "image/bmp"},
	{Type: "image", SubType: "svg+xml", Description: "SVG 矢量图", ContentType: "image/svg+xml"},

	// 音频类型
	{Type: "audio", SubType: "mpeg", Description: "MPEG 音频", ContentType: "audio/mpeg"},
	{Type: "audio", SubType: "wav", Description: "WAV 音频", ContentType: "audio/wav"},
	{Type: "audio", SubType: "ogg", Description: "OGG 音频", ContentType: "audio/ogg"},

	// 视频类型
	{Type: "video", SubType: "mp4", Description: "MP4 视频", ContentType: "video/mp4"},
	{Type: "video", SubType: "mpeg", Description: "MPEG 视频", ContentType: "video/mpeg"},
	{Type: "video", SubType: "ogg", Description: "OGG 视频", ContentType: "video/ogg"},
	{Type: "video", SubType: "webm", Description: "WebM 视频", ContentType: "video/webm"},
}
