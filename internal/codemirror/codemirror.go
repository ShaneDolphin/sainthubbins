// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live console editor stub — CodeMirror 6 is browser-only; Go stub provides theme/flash API.

package codemirror

// Editor is stub for the live console editor.
type Editor struct {
	Content string
	Theme   string
}

func New(content string) *Editor { return &Editor{Content: content, Theme: "hubbinsTheme"} }

func (e *Editor) SetContent(s string) { e.Content = s }
func (e *Editor) Flash(from, to int) {}
