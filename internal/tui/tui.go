package tui

import (
	"bytes"
	"encoding/json"

	"github.com/r-hermanto/leqman/internal/leq"
	"github.com/rivo/tview"
)

type Application struct {
	App         *tview.Application
	Collections *tview.TreeView
	Request     *tview.TextView
	Response    *tview.TextView
}

func (a *Application) Run() {
	if err := a.App.Run(); err != nil {
		panic(err)
	}
}

func NewApplication() *Application {
	apps := &Application{
		App:         tview.NewApplication(),
		Collections: tview.NewTreeView(),
		Request:     tview.NewTextView(),
		Response:    tview.NewTextView(),
	}

	collections := leq.GetCollections()
	collectionRoot := buildTreeNode(collections)

	apps.Collections.
		SetRoot(collectionRoot).
		SetCurrentNode(collectionRoot).
		SetTopLevel(1)

	apps.Collections.
		SetBorder(true).
		SetTitle("Collections Section")

	apps.Collections.SetSelectedFunc(func(node *tview.TreeNode) {
		c := node.GetReference().(*leq.Collection)

		if c.IsDir {
			return
		}

		lc := leq.GetRequest(c.Path)
		b, err := json.Marshal(lc)
		if err != nil {
			panic(err)
		}
		var req bytes.Buffer
		json.Indent(&req, b, "", "    ")

		apps.Request.SetText(req.String())
	})

	apps.Request.
		SetBorder(true).
		SetTitle("Request Section")

	lc := leq.GetRequest(collections[0].Path)
	b, err := json.Marshal(lc)
	if err != nil {
		panic(err)
	}
	var req bytes.Buffer
	json.Indent(&req, b, "", "    ")

	apps.Request.SetText(req.String())

	resp := lc.Execute()
	apps.Response.SetText(resp)

	flex := tview.NewFlex().
		AddItem(apps.Collections, 0, 1, true).
		AddItem(apps.Request, 0, 3, false).
		AddItem(apps.Response, 0, 3, false)

		// TODO:
		// 4. update design to handle
		//      a. url
		//      b. header
		//      c. body
		//      d. options

	apps.App.SetRoot(flex, true).EnableMouse(true)

	return apps
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
