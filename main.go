// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/math"
	"github.com/google/gxui/samples/flags"
	"strconv"
	//"github.com/google/gxui/gxfont"
	"fmt"
)

var image gxui.Image
var label3 gxui.Label
var customDriver *gxui.Driver

var size1, size2 uint64

type Vertice struct {
	number int
	typev int
	posX int
	posY int
	radius int
	selected bool
}

type Relation struct {
	first *Vertice
	second *Vertice
	selected bool
}

var width = 407
var height = 640

var vertices []Vertice
var relations []Relation

var lastSelected *Vertice

var solutions *[][]int

var solutionCount = 0
var currentSolution = 0

func init() {
	vertices = make([]Vertice, 0)
	relations = make([]Relation, 0)
	lastSelected = nil
}

func countChanged() {
	vertices = make([]Vertice, size1+size2)
	relations = relations[:0]
	lastSelected = nil

	padding1 := int(uint64(height)/(size1+1))
	padding2 := int(uint64(height)/(size2+1))
	blockSize1 := padding1/5
	for i:=0; i < int(size1); i++ {
		vertices[i].posX = 35
		vertices[i].posY = padding1*(i+1)
		vertices[i].typev = 0
		vertices[i].number = i
		vertices[i].radius = blockSize1
	}
	blockSize2 := padding2/5
	for i:=0; i < int(size2); i++ {
		vertices[int(size1)+i].posX = width-35
		vertices[int(size1)+i].posY = padding2*(i+1)
		vertices[int(size1)+i].typev = 1
		vertices[int(size1)+i].number = i
		vertices[int(size1)+i].radius = blockSize2
	}
	RedrawCanvas()
}

func RedrawCanvas() {
	if size1 == 0 || size2 == 0 {
		return
	}
	driver := *customDriver
	canvas := driver.CreateCanvas(math.Size { width, height})
	image.SetSize(math.Size{width, height})
	brush := gxui.CreateBrush(gxui.Black)
	brushRed := gxui.CreateBrush(gxui.Red)
	pen := gxui.CreatePen(1, gxui.White)
	penRelation := gxui.CreatePen(2, gxui.Blue)
	penRelationSelected := gxui.CreatePen(4, gxui.Red)

	for _, v := range vertices {
		if v.selected {
			canvas.DrawRoundedRect(math.CreateRect(v.posX-v.radius,v.posY-v.radius,v.posX+v.radius,v.posY+v.radius), float32(v.radius), float32(v.radius), float32(v.radius), float32(v.radius),pen,brushRed)
		} else {
			canvas.DrawRoundedRect(math.CreateRect(v.posX-v.radius,v.posY-v.radius,v.posX+v.radius,v.posY+v.radius), float32(v.radius), float32(v.radius), float32(v.radius), float32(v.radius),pen,brush)
		}
	}

	for _, v := range relations {
		pol := gxui.Polygon{ gxui.PolygonVertex{Position: math.Point{X: v.first.posX, Y: v.first.posY}},gxui.PolygonVertex{Position: math.Point{X: v.second.posX, Y: v.second.posY}}}
		if v.selected {
			canvas.DrawLines(pol, penRelationSelected)
		} else {
			canvas.DrawLines(pol, penRelation)
		}
	}
	canvas.Complete()
	image.SetCanvas(canvas)
}

func Compute(me gxui.MouseEvent) {
	data := make([][]int,size1)
	for _, v := range relations {
		data[v.first.number] = append(data[v.first.number], v.second.number)
	}
	solutions = GetSolutions(&data)
	solutionCount = len(*solutions)
	currentSolution = 0
	ShowSolution()
	updateCountLabel()
}

func ShowSolution() {
	if solutionCount == 0 {
		return
	}
	ClearSelections()
	for i := 0; i < len((*solutions)[currentSolution]); i += 2 {
		v1 := (*solutions)[currentSolution][i]
		v2 := (*solutions)[currentSolution][i+1]
		for itd, v := range relations {
			if v.first.number == v1 && v.second.number == v2 {
				relations[itd].selected = true
				break
			}
		}
	}
	RedrawCanvas()
}

func PreviousSolution(me gxui.MouseEvent) {
	if solutionCount == 0 {
		currentSolution = 0
	} else {
		currentSolution--
		if currentSolution < 0 {
			currentSolution = solutionCount - 1
		}
		ShowSolution()
	}
	updateCountLabel()
}

func NextSolution(me gxui.MouseEvent) {
	if solutionCount == 0 {
		currentSolution = 0
	} else {
		currentSolution++
		if currentSolution >= solutionCount {
			currentSolution = 0
		}
		ShowSolution()
	}
	updateCountLabel()
}

func updateCountLabel() {
	var str string
	if solutionCount == 0 {
		str = "00000/00000"
	} else {
		str = fmt.Sprintf("%05d/%05d", currentSolution+1, solutionCount)
	}
	label3.SetText(str)
}


func ClearSelections() {
	for i := range relations {
		relations[i].selected = false
	}
}

func FullClearSelections() {
	ClearSelections()
	solutions = nil
	currentSolution = 0
	solutionCount = 0
	updateCountLabel()
}

func ProceedClick(me gxui.MouseEvent) {
	for i, v := range vertices {
		if me.Point.X >= v.posX-v.radius*3/2 && me.Point.X <= v.posX+v.radius*3/2 && me.Point.Y >= v.posY-v.radius*3/2 && me.Point.Y <= v.posY+v.radius*3/2 {
			if lastSelected != nil {
				if lastSelected.typev == vertices[i].typev {
					lastSelected.selected = false
				} else {
					relation := new(Relation)
					if lastSelected.typev == 0 {
						relation.first = lastSelected
						relation.second = &vertices[i]
					} else {
						relation.second = lastSelected
						relation.first = &vertices[i]
					}
					FullClearSelections()
					found := false
					for ind, relv := range relations {
						if relv.first == relation.first && relv.second == relation.second {
							relations[ind] = relations[len(relations)-1]
							relations = relations[:len(relations)-1]
							found = true
							RedrawCanvas()
							break
						}
					}
					if found == false {
						relations = append(relations, *relation)
					}
					lastSelected.selected = false
					lastSelected = nil
					RedrawCanvas()
					break
				}
			}
			lastSelected = &vertices[i]
			vertices[i].selected = true
			RedrawCanvas()
			break
		}
	}
}


func appMain(driver gxui.Driver) {
	customDriver = &driver
	theme := flags.CreateTheme(driver)

	/*font, err := driver.CreateFont(gxfont.Default, 20)
	if err != nil {
		panic(err)
	}*/

	window := theme.CreateWindow(407, 700, "biGraphSolver")
	window.SetTitle("BiGraph Solver")
	brush := gxui.CreateBrush(gxui.Black)
	window.SetBackgroundBrush(brush)

	size1Button := theme.CreateTextBox()
	//size1Button.SetSize(math.Size{10, 10})
	//size1Button.SetFont(font)
	size1Button.OnTextChanged(func(param []gxui.TextBoxEdit) {
		temp, ok := strconv.ParseUint(size1Button.Text(), 10, 64)
		if ok == nil {
			size1 = temp
			countChanged()
		}
	})


	size2Button := theme.CreateTextBox()
	//size2Button.SetFont(font)

	size2Button.OnTextChanged(func(param []gxui.TextBoxEdit) {
		temp, ok := strconv.ParseUint(size2Button.Text(), 10, 64)
		if ok == nil {
			size2 = temp
			countChanged()
		}
	})


	label := theme.CreateLabel()
	label.SetText("Count left:")
	label2 := theme.CreateLabel()
	label2.SetText("Count right:")

	button := theme.CreateButton()
	button.SetText("Compute")
	button.SetPadding(math.Spacing{L: 17, R: 17, T: 17, B: 17})
	button2 := theme.CreateButton()
	button2.SetText("<")
	button2.SetPadding(math.Spacing{L: 6, R: 6, T: 17, B: 17})
	button3 := theme.CreateButton()
	button3.SetText(">")
	button3.SetPadding(math.Spacing{L: 6, R: 6, T: 17, B: 17})
	label3 = theme.CreateLabel()
	label3.SetVerticalAlignment(gxui.AlignMiddle)
	label3.SetText("00000/00000")

	layoutMain := theme.CreateLinearLayout()
	layoutMain.SetHorizontalAlignment(gxui.AlignRight)
	layout1 := theme.CreateLinearLayout()
	layout1.SetDirection(gxui.LeftToRight)
	layout1.SetVerticalAlignment(gxui.AlignMiddle)
	layout1.AddChild(label)
	layout1.AddChild(size1Button)
	layout2 := theme.CreateLinearLayout()
	layout2.SetDirection(gxui.LeftToRight)
	layout2.SetVerticalAlignment(gxui.AlignMiddle)
	layout2.AddChild(label2)
	layout2.AddChild(size2Button)
	layoutMain.AddChild(layout1)
	layoutMain.AddChild(layout2)

	layoutMaMain := theme.CreateLinearLayout()
	layoutMaMain.SetDirection(gxui.LeftToRight)
	layoutMaMain.SetVerticalAlignment(gxui.AlignMiddle)
	layoutMaMain.SetPadding(math.Spacing{L: 2})
	layoutMaMain.AddChild(layoutMain)
	layoutMaMain.AddChild(button)
	layoutMaMain.AddChild(button2)
	layoutMaMain.AddChild(label3)
	layoutMaMain.AddChild(button3)

	button.OnClick(Compute)
	button2.OnClick(PreviousSolution)
	button3.OnClick(NextSolution)

	image = theme.CreateImage()
	image.OnClick(ProceedClick)
	layout := theme.CreateLinearLayout()
	layout.AddChild(layoutMaMain)
	layout.AddChild(image)
	window.AddChild(layout)

	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}