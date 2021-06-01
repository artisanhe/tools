package pic_file

import (
	"fmt"
	"strings"
)

//文件内容识别
func GetContentTypeAndContentStr(content string) (content_type string, content_str string, err error) {
	if !strings.HasPrefix(content, DATA_PEFIX) {
		return "", "", fmt.Errorf("content has no data prefix")
	}

	var encoding_index int
	if encoding_index = strings.Index(content, ENCODING_DESC); encoding_index == -1 {
		return "", "", fmt.Errorf("content has no base64 encode descx")
	}

	image_type := content[len(DATA_PEFIX):encoding_index]
	content_type = image_type[len("image/"):]

	content_str = content[encoding_index+len(ENCODING_DESC):]
	return
}

const (
	DATA_PEFIX    string = "data:"
	ENCODING_DESC string = ";base64,"
)
