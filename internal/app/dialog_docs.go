package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"coilforge/internal/render"
	"coilforge/internal/storage"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type docDialogMode int

const (
	docDialogClosed docDialogMode = iota
	docDialogLoad
	docDialogSave
)

type docDialogState struct {
	mode     docDialogMode
	store    storage.DocStore
	docs     []storage.DocInfo
	selected int
	input    string
	errText  string

	lastClickAt    time.Time
	lastClickIndex int
}

func (a *App) openDocDialog(mode docDialogMode) {
	a.settingsOpen = false
	a.docDialog.mode = mode
	a.docDialog.errText = ""
	a.docDialog.lastClickAt = time.Time{}
	a.docDialog.lastClickIndex = -1
	a.refreshDocDialogList()
	if len(a.docDialog.docs) > 0 {
		a.docDialog.selected = 0
		if mode == docDialogLoad || strings.TrimSpace(a.docDialog.input) == "" {
			a.docDialog.input = a.docDialog.docs[0].Name
		}
		return
	}
	a.docDialog.selected = -1
	if mode == docDialogSave && strings.TrimSpace(a.docDialog.input) == "" {
		a.docDialog.input = DefaultProjectPath
	}
}

func (a *App) closeDocDialog(errText string) {
	a.docDialog.mode = docDialogClosed
	a.docDialog.errText = errText
	a.docDialog.lastClickAt = time.Time{}
	a.docDialog.lastClickIndex = -1
}

func (a *App) handleDocDialogMousePress(mouseX, mouseY int) bool {
	idx, ok := render.DocBrowserRowAtScreenPoint(mouseX, mouseY, len(a.docDialog.docs))
	if !ok {
		a.docDialog.lastClickIndex = -1
		return false
	}
	if idx < 0 || idx >= len(a.docDialog.docs) {
		a.docDialog.lastClickIndex = -1
		return false
	}
	a.docDialog.selected = idx
	a.docDialog.input = a.docDialog.docs[idx].Name
	now := time.Now()
	if a.docDialog.mode == docDialogLoad &&
		a.docDialog.lastClickIndex == idx &&
		now.Sub(a.docDialog.lastClickAt) <= 400*time.Millisecond {
		a.commitDocDialog()
		a.docDialog.lastClickIndex = -1
		a.docDialog.lastClickAt = time.Time{}
		return true
	}
	a.docDialog.lastClickIndex = idx
	a.docDialog.lastClickAt = now
	return true
}

func (a *App) refreshDocDialogList() {
	if a.docDialog.store == nil {
		a.docDialog.docs = nil
		a.docDialog.selected = -1
		return
	}
	docs, err := a.docDialog.store.ListDocs()
	if err != nil {
		a.docDialog.docs = nil
		a.docDialog.selected = -1
		a.docDialog.errText = err.Error()
		return
	}
	a.docDialog.docs = docs
	if len(docs) == 0 {
		a.docDialog.selected = -1
		return
	}
	if a.docDialog.selected < 0 || a.docDialog.selected >= len(docs) {
		a.docDialog.selected = 0
	}
}

func (a *App) handleDocDialogTyping() {
	if a.docDialog.mode == docDialogClosed {
		return
	}
	for _, key := range inpututil.AppendJustPressedKeys(nil) {
		a.handleDocDialogKey(key)
	}
	a.appendDocDialogChars()
}

func (a *App) handleDocDialogKey(key ebiten.Key) {
	switch key {
	case ebiten.KeyArrowUp:
		a.selectPrevDoc()
	case ebiten.KeyArrowDown:
		a.selectNextDoc()
	case ebiten.KeyBackspace:
		a.backspaceDocInput()
	case ebiten.KeyEnter:
		a.commitDocDialog()
	case ebiten.KeyDelete:
		a.deleteSelectedDoc()
	}
}

func (a *App) selectPrevDoc() {
	if a.docDialog.selected <= 0 {
		return
	}
	a.docDialog.selected--
	a.docDialog.input = a.docDialog.docs[a.docDialog.selected].Name
}

func (a *App) selectNextDoc() {
	if a.docDialog.selected < 0 || a.docDialog.selected >= len(a.docDialog.docs)-1 {
		return
	}
	a.docDialog.selected++
	a.docDialog.input = a.docDialog.docs[a.docDialog.selected].Name
}

func (a *App) backspaceDocInput() {
	runes := []rune(a.docDialog.input)
	if len(runes) == 0 {
		return
	}
	a.docDialog.input = string(runes[:len(runes)-1])
}

func (a *App) appendDocDialogChars() {
	for _, ch := range ebiten.AppendInputChars(nil) {
		if ch < 32 || ch == 127 || ch == '/' || ch == '\\' {
			continue
		}
		a.docDialog.input += string(ch)
	}
}

func (a *App) commitDocDialog() {
	name := strings.TrimSpace(a.docDialog.input)
	if name == "" {
		a.docDialog.errText = "Filename is required."
		return
	}
	if a.docDialog.store == nil {
		a.docDialog.errText = "Document storage unavailable."
		return
	}
	switch a.docDialog.mode {
	case docDialogLoad:
		data, err := a.docDialog.store.LoadDoc(name)
		if err != nil {
			a.docDialog.errText = err.Error()
			return
		}
		if err := UnmarshalProject(data); err != nil {
			a.docDialog.errText = err.Error()
			return
		}
		a.closeDocDialog("")
	case docDialogSave:
		data, err := MarshalProject()
		if err != nil {
			a.docDialog.errText = err.Error()
			return
		}
		if err := a.docDialog.store.SaveDoc(name, data); err != nil {
			a.docDialog.errText = err.Error()
			return
		}
		a.closeDocDialog("")
	}
}

func (a *App) deleteSelectedDoc() {
	if a.docDialog.store == nil {
		a.docDialog.errText = "Document storage unavailable."
		return
	}
	if a.docDialog.selected < 0 || a.docDialog.selected >= len(a.docDialog.docs) {
		return
	}
	name := a.docDialog.docs[a.docDialog.selected].Name
	err := a.docDialog.store.DeleteDoc(name)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		a.docDialog.errText = err.Error()
		return
	}
	a.refreshDocDialogList()
	if len(a.docDialog.docs) > 0 {
		a.docDialog.input = a.docDialog.docs[a.docDialog.selected].Name
		return
	}
	a.docDialog.input = ""
}

func (a *App) docDialogTitle() string {
	switch a.docDialog.mode {
	case docDialogLoad:
		return "Load Project"
	case docDialogSave:
		return "Save Project"
	default:
		return ""
	}
}

func (a *App) docDialogRows() []render.DocBrowserRow {
	rows := make([]render.DocBrowserRow, 0, len(a.docDialog.docs))
	for i := range a.docDialog.docs {
		doc := a.docDialog.docs[i]
		text := fmt.Sprintf("%s  (%d B)", doc.Name, doc.SizeBytes)
		rows = append(rows, render.DocBrowserRow{
			Text:     text,
			Selected: i == a.docDialog.selected,
		})
	}
	return rows
}

func (a *App) docDialogFooter() string {
	if a.docDialog.errText != "" {
		return a.docDialog.errText
	}
	return "Enter: confirm  Del: delete selected  Esc: close"
}
