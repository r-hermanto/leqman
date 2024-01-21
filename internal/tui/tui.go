package tui

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/r-hermanto/leqman/internal/leq"
	"github.com/rivo/tview"
)

var httpMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
	http.MethodDelete,
	http.MethodHead,
	http.MethodOptions,
}

type RequestUI struct {
	layout  *tview.Flex
	method  *tview.DropDown
	url     *tview.InputField
	header  *tview.TextArea
	body    *tview.TextArea
	config  *leq.LeqConfig
	changed bool
}

func (u *RequestUI) updateScreen() {
	config := u.config

	u.method.SetCurrentOption(slices.Index(httpMethods, config.Method))
	u.url.SetText(config.URL)

	header := ""
	if config.Header != nil {
		b, err := json.MarshalIndent(config.Header, "", " ")
		if err != nil {
			panic(err)
		}

		header = string(b)
	}
	u.header.SetText(header, false)

	body := ""
	if !isJsonRawMessageNil(config.Body) {
		b, err := json.MarshalIndent(config.Body, "", " ")
		if err != nil {
			panic(err)
		}
		body = string(b)
	}
	u.body.SetText(body, false)
}

type Application struct {
	app         *tview.Application
	layout      *tview.Flex
	collections *tview.TreeView
	requetsUI   RequestUI
	response    *tview.TextView
}

func (a *Application) Run() {
	a.fetchCollections()
	a.app.SetRoot(a.layout, true).
		EnableMouse(true)

	if err := a.app.Run(); err != nil {
		panic(err)
	}
}

func (a *Application) fetchCollections() {
	collections := leq.GetCollections()
	collectionRoot := buildTreeNode(collections)

	a.collections.
		SetRoot(collectionRoot).
		SetCurrentNode(collectionRoot)
}

func newRequestUI() RequestUI {
	ui := RequestUI{
		layout: tview.NewFlex(),
		method: tview.NewDropDown(),
		url:    tview.NewInputField(),
		header: tview.NewTextArea(),
		body:   tview.NewTextArea(),
	}

	ui.method.
		SetOptions(httpMethods, nil).
		SetCurrentOption(0)

	ui.method.SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	dropdownIcon := tview.NewTextView().
		SetText("▼  ").
		SetTextAlign(tview.AlignRight).
		SetDynamicColors(true)
	ui.url.SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	urlGroup := tview.NewFlex().
		AddItem(ui.method, 6, 1, false).
		AddItem(dropdownIcon, 3, 1, false).
		AddItem(ui.url, 0, 5, false)
	urlGroup.SetBorder(true).SetTitle("url").SetTitleAlign(tview.AlignLeft)

	ui.header.SetBorder(true).SetTitle("header").SetTitleAlign(tview.AlignLeft)
	ui.body.SetBorder(true).SetTitle("body").SetTitleAlign(tview.AlignLeft)

	ui.layout.
		SetDirection(tview.FlexRow).
		AddItem(urlGroup, 3, 1, false).
		AddItem(ui.header, 0, 2, false).
		AddItem(ui.body, 0, 3, false)

	return ui
}

func NewApplication() *Application {
	app := &Application{
		app:         tview.NewApplication(),
		layout:      tview.NewFlex(),
		collections: tview.NewTreeView(),
		requetsUI:   newRequestUI(),
		response:    tview.NewTextView(),
	}

	app.collections.
		SetTopLevel(1).
		SetBorder(true).
		SetTitle("Collection Section")

	app.collections.SetSelectedFunc(func(node *tview.TreeNode) {
		c := node.GetReference().(*leq.Collection)

		if c.IsDir {
			return
		}

		lc := leq.GetRequest(c.Path)
		app.requetsUI.config = &lc
		app.requetsUI.updateScreen()
	})

	app.response.
		SetBorder(true).
		SetTitle("Response Section")

	app.layout.
		AddItem(app.collections, 0, 1, true).
		AddItem(app.requetsUI.layout, 0, 3, true).
		AddItem(app.response, 0, 3, true)

	app.requetsUI.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlS {
			ui := app.requetsUI

			url := ui.url.GetText()
			_, method := ui.method.GetCurrentOption()

			var header map[string]string
			if ui.header.GetText() != "" {
				err := json.Unmarshal([]byte(ui.header.GetText()), &header)
				if err != nil {
					panic(err)
				}
			}

			var body []byte
			if ui.body.GetText() != "" {
				body = []byte(ui.body.GetText())
			}

			config := ui.config
			config.Method = method
			config.URL = url
			config.Header = header
			config.Body = body

			leq.UpdateRequest(*config)
		}

		return event
	})

	return app
}

func buildTreeNode(c []*leq.Collection) *tview.TreeNode {
	root := tview.NewTreeNode(".")

	var addChild func(c *leq.Collection) *tview.TreeNode
	addChild = func(c *leq.Collection) *tview.TreeNode {
		node := tview.
			NewTreeNode(c.Title).
			SetReference(c)

		if !c.IsDir {
			return node
		}

		for _, collection := range c.Children {
			node.AddChild(addChild(collection))
		}

		return node
	}

	for _, collection := range c {
		root.AddChild(addChild(collection))
	}

	return root
}

func isJsonRawMessageNil(d json.RawMessage) bool {
	return len(d) == 4 && string(d) == "null"
}
