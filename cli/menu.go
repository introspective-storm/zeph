package main

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type item struct {
	title  string
	desc   string
	action string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func (i item) GetAction() string   { return i.action }

type listModel struct {
	list list.Model
}

func mainMenu() listModel {
	items := []list.Item{
		item{
			title:  "Load Existing",
			desc:   "Choose a config to load a previous project",
			action: "load",
		},
		item{
			title:  "New Project",
			desc:   "Create a project with the interactive wizard",
			action: "wizard",
		},
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	return listModel{list: l}
}

func (l listModel) WithProject(projectName string) listModel {
	existingItems := l.list.Items()
	projectItem := item{
		title:  fmt.Sprintf("Project: %s", projectName),
		desc:   "Continue your last project",
		action: "cont",
	}
	newItems := append([]list.Item{projectItem}, existingItems...)
	l.list.SetItems(newItems)
	return l
}

func (l listModel) SelectedItem() list.Item {
	return l.list.SelectedItem()
}

func (l listModel) SetSize(width int, height int) listModel {
	l.list.SetSize(width, height)
	return l
}

func (l listModel) Update(msg tea.Msg) (listModel, tea.Cmd) {
	var cmd tea.Cmd
	l.list, cmd = l.list.Update(msg)
	return l, cmd
}

func (l listModel) View() string {
	return l.list.View()
}
