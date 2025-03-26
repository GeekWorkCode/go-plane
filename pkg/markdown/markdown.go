package markdown

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	// 匹配用户@提及
	// Match user @mentions
	mentionRegex = regexp.MustCompile(`@([a-zA-Z0-9_-]+)`)

	// 匹配加粗文本
	// Match bold text
	boldRegex = regexp.MustCompile(`\*\*(.*?)\*\*`)

	// 匹配斜体文本
	// Match italic text
	italicRegex = regexp.MustCompile(`\*(.*?)\*|_(.*?)_`)

	// 匹配链接
	// Match links
	linkRegex = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)

	// 匹配代码块
	// Match code blocks
	codeBlockRegex = regexp.MustCompile("```([a-zA-Z0-9]*)\n([\\s\\S]*?)```")

	// 匹配内联代码
	// Match inline code
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")

	// 匹配标题
	// Match headers
	headerRegex = regexp.MustCompile(`(?m)^(#{1,6})\s+(.*)$`)

	// 匹配无序列表
	// Match unordered lists
	listItemRegex = regexp.MustCompile(`(?m)^[ \t]*[-*+][ \t]+(.*)$`)
)

// ToHTML 将 Markdown 文本转换为 HTML
// ToHTML converts Markdown text to HTML
func ToHTML(markdown string) string {
	if markdown == "" {
		return ""
	}

	html := markdown

	// 处理代码块
	// Process code blocks
	html = codeBlockRegex.ReplaceAllStringFunc(html, func(match string) string {
		submatches := codeBlockRegex.FindStringSubmatch(match)
		lang := submatches[1]
		code := submatches[2]

		if lang == "" {
			return "<pre><code>" + code + "</code></pre>"
		}
		return "<pre><code class=\"language-" + lang + "\">" + code + "</code></pre>"
	})

	// 处理内联代码
	// Process inline code
	html = inlineCodeRegex.ReplaceAllString(html, "<code>$1</code>")

	// 处理标题
	// Process headers
	html = headerRegex.ReplaceAllStringFunc(html, func(match string) string {
		submatches := headerRegex.FindStringSubmatch(match)
		level := len(submatches[1])
		text := submatches[2]

		return "<h" + strconv.Itoa(level) + ">" + text + "</h" + strconv.Itoa(level) + ">"
	})

	// 处理加粗文本
	// Process bold text
	html = boldRegex.ReplaceAllString(html, "<strong>$1</strong>")

	// 处理斜体文本
	// Process italic text
	html = italicRegex.ReplaceAllStringFunc(html, func(match string) string {
		if strings.HasPrefix(match, "*") {
			return "<em>" + match[1:len(match)-1] + "</em>"
		}
		return "<em>" + match[1:len(match)-1] + "</em>"
	})

	// 处理链接
	// Process links
	html = linkRegex.ReplaceAllString(html, "<a href=\"$2\">$1</a>")

	// 处理无序列表
	// Process unordered lists
	if listItemRegex.MatchString(html) {
		var b bytes.Buffer
		lines := strings.Split(html, "\n")
		inList := false

		for _, line := range lines {
			if listItemRegex.MatchString(line) {
				if !inList {
					b.WriteString("<ul>\n")
					inList = true
				}
				matches := listItemRegex.FindStringSubmatch(line)
				b.WriteString("  <li>" + matches[1] + "</li>\n")
			} else {
				if inList {
					b.WriteString("</ul>\n")
					inList = false
				}
				b.WriteString(line + "\n")
			}
		}

		if inList {
			b.WriteString("</ul>\n")
		}

		html = b.String()
	}

	// 处理用户@提及
	// Process user @mentions
	html = mentionRegex.ReplaceAllString(html, "<span class=\"mention\">@$1</span>")

	// 将换行符转换为<br>标签
	// Convert newlines to <br> tags
	html = strings.ReplaceAll(html, "\n\n", "<br><br>")

	return html
}
