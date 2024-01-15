package tui

import (
	"github.com/r-hermanto/leqman/internal/leq"
	"github.com/rivo/tview"
)

type Application struct {
	App         *tview.Application
	Collections *tview.TreeView
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
	}

	collections := leq.GetCollections()
	collectionRoot := buildTreeNode(collections)

	apps.Collections.
		SetRoot(collectionRoot).
		SetCurrentNode(collectionRoot).
		SetTopLevel(1)

	apps.Collections.
		SetBorder(true).
		SetTitle("Collection Section")

	flex := tview.NewFlex().
		AddItem(apps.Collections, 0, 1, true).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Request Section"), 0, 3, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Response Section"), 0, 3, false)

	apps.App.SetRoot(flex, true).EnableMouse(true)

	return apps
}

func buildTreeNode(c []*leq.Collection) *tview.TreeNode {
	root := tview.NewTreeNode(".")

	var addChild func(c *leq.Collection) *tview.TreeNode
	addChild = func(c *leq.Collection) *tview.TreeNode {
		node := tview.NewTreeNode(c.Title)

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
