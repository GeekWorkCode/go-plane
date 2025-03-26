package markdown

import (
	"testing"
)

func TestToHTML(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "empty string",
			markdown: "",
			want:     "",
		},
		{
			name:     "simple paragraph",
			markdown: "This is a simple paragraph.",
			want:     "This is a simple paragraph.",
		},
		{
			name:     "bold text",
			markdown: "This is **bold** text.",
			want:     "This is <strong>bold</strong> text.",
		},
		{
			name:     "italic text with asterisks",
			markdown: "This is *italic* text.",
			want:     "This is <em>italic</em> text.",
		},
		{
			name:     "italic text with underscores",
			markdown: "This is _italic_ text.",
			want:     "This is <em>italic</em> text.",
		},
		{
			name:     "link",
			markdown: "Check out [this link](https://example.com).",
			want:     "Check out <a href=\"https://example.com\">this link</a>.",
		},
		{
			name:     "inline code",
			markdown: "Use the `fmt.Println()` function.",
			want:     "Use the <code>fmt.Println()</code> function.",
		},
		{
			name:     "code block with language",
			markdown: "```go\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}\n```",
			want:     "<pre><code class=\"language-go\">func main() {\n\tfmt.Println(\"Hello, World!\")\n}\n</code></pre>",
		},
		{
			name:     "code block without language",
			markdown: "```\nplain text code block\n```",
			want:     "<pre><code>plain text code block\n</code></pre>",
		},
		{
			name:     "h1 header",
			markdown: "# Header 1",
			want:     "<h1>Header 1</h1>",
		},
		{
			name:     "h2 header",
			markdown: "## Header 2",
			want:     "<h2>Header 2</h2>",
		},
		{
			name:     "h3 header",
			markdown: "### Header 3",
			want:     "<h3>Header 3</h3>",
		},
		{
			name:     "unordered list",
			markdown: "- Item 1\n- Item 2\n- Item 3",
			want:     "<ul>\n  <li>Item 1</li>\n  <li>Item 2</li>\n  <li>Item 3</li>\n</ul>\n",
		},
		{
			name:     "mentions",
			markdown: "Hello @user!",
			want:     "Hello <span class=\"mention\">@user</span>!",
		},
		{
			name:     "multiple paragraphs",
			markdown: "Paragraph 1.\n\nParagraph 2.",
			want:     "Paragraph 1.<br><br>Paragraph 2.",
		},
		{
			name:     "mixed formatting",
			markdown: "# Title\n\nThis is a **bold** statement with a [link](https://example.com) and some `code`.\n\n- List item 1\n- List item 2",
			want:     "<h1>Title</h1><br><br>This is a <strong>bold</strong> statement with a <a href=\"https://example.com\">link</a> and some <code>code</code>.<br><br><ul>\n  <li>List item 1</li>\n  <li>List item 2</li>\n</ul>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToHTML(tt.markdown)
			if got != tt.want {
				t.Errorf("ToHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
