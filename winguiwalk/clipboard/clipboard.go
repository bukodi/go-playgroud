// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var te *walk.TextEdit

	if _, err := (MainWindow{
		Title:   "Margareta Desktop Tool",
		MinSize: Size{Width: 300, Height: 200},
		Layout:  VBox{},
		Children: []Widget{
			TextEdit{
				AssignTo:   &te,
				VScroll:    true,
				ReadOnly:   true,
				Background: SolidColorBrush{Color: walk.RGB(255, 255, 127)},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "Copy",
						OnClicked: func() {
							if err := walk.Clipboard().SetText(te.Text()); err != nil {
								log.Print("Copy: ", err)
							}
						},
					},
					PushButton{
						Text: "Paste",
						OnClicked: func() {
							if text, err := walk.Clipboard().Text(); err != nil {
								log.Print("Paste: ", err)
							} else {
								te.SetText(text)
							}
						},
					},
				},
			},
		},
	}).Run(); err != nil {
		log.Fatal(err)
	}
}
