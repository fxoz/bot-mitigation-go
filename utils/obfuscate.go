package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func ObfuscateJS(jsCode string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(jsCode))

	chunkSize := 10
	var chunks []string
	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]
		chunks = append(chunks, reverseString(chunk))
	}

	var sb strings.Builder
	sb.WriteString("(function(){\n")
	sb.WriteString("  if(false){var _0xdummy = function(){return 'dummy';};}\n")
	sb.WriteString("  var _0xarr = [")
	for i, chunk := range chunks {
		sb.WriteString(fmt.Sprintf("\"%s\"", chunk))
		if i != len(chunks)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("];\n")
	sb.WriteString("  var _0xstr = \"\";\n")
	sb.WriteString("  for(var i = _0xarr.length - 1; i >= 0; i--){\n")
	sb.WriteString("      _0xstr += _0xarr[i];\n")
	sb.WriteString("  }\n")
	sb.WriteString("  _0xstr = _0xstr.split(\"\").reverse().join(\"\");\n")
	sb.WriteString("  eval(atob(_0xstr));\n")
	sb.WriteString("})();")

	return sb.String()
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
