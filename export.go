package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePNG creates a PNG image of the Gantt chart
func GeneratePNG(chart *Chart) ([]byte, error) {
	// Calculate dimensions
	quarters := calculateQuarters(chart.StartYear, chart.StartQ, chart.EndYear, chart.EndQ)
	totalQuarters := len(quarters)

	// Layout constants
	headerHeight := 80
	rowHeight := 40
	quarterWidth := 120
	labelWidth := 200
	padding := 20
	categoryHeaderHeight := 35

	// Count total rows
	totalRows := 0
	for _, cat := range chart.Categories {
		totalRows++ // Category header
		totalRows += len(cat.Tasks)
	}

	width := labelWidth + totalQuarters*quarterWidth + padding*2
	height := headerHeight + totalRows*rowHeight + padding*2

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{250, 250, 250, 255}}, image.Point{}, draw.Src)

	// Draw grid and quarters
	for i := range quarters {
		x := padding + labelWidth + i*quarterWidth
		y := headerHeight

		// Quarter background
		bgColor := color.RGBA{232, 232, 232, 255}
		if i%2 == 1 {
			bgColor = color.RGBA{245, 245, 245, 255}
		}
		drawRect(img, x, y-30, quarterWidth, 30, bgColor)
		drawRectBorder(img, x, y-30, quarterWidth, 30, color.RGBA{204, 204, 204, 255})

		// Vertical grid lines
		drawVerticalLine(img, x, y, height-padding, color.RGBA{221, 221, 221, 255})
	}

	// Draw categories and tasks
	currentY := headerHeight
	for _, cat := range chart.Categories {
		catColor := parseColor(cat.Color)

		// Category header background
		catColorAlpha := color.RGBA{catColor.R, catColor.G, catColor.B, 76} // 0.3 opacity
		drawRect(img, padding, currentY, labelWidth, categoryHeaderHeight, catColorAlpha)

		// Category background span
		catColorLight := color.RGBA{catColor.R, catColor.G, catColor.B, 13} // 0.05 opacity
		drawRect(img, padding+labelWidth, currentY, totalQuarters*quarterWidth, categoryHeaderHeight, catColorLight)

		currentY += categoryHeaderHeight

		// Tasks
		for _, task := range cat.Tasks {
			// Task label background
			drawRect(img, padding, currentY, labelWidth, rowHeight, color.RGBA{255, 255, 255, 255})
			drawRectBorder(img, padding, currentY, labelWidth, rowHeight, color.RGBA{221, 221, 221, 255})

			// Draw task bar
			startIdx := findQuarterIndex(quarters, task.StartYear, task.StartQ)
			endIdx := findQuarterIndex(quarters, task.EndYear, task.EndQ)

			if startIdx >= 0 && endIdx >= 0 {
				barX := padding + labelWidth + startIdx*quarterWidth
				barWidth := (endIdx - startIdx + 1) * quarterWidth
				barY := currentY + 8
				barHeight := rowHeight - 16

				taskColor := parseColor(task.Color)
				if task.Color == "" {
					taskColor = catColor
				}

				taskColorAlpha := color.RGBA{taskColor.R, taskColor.G, taskColor.B, 204} // 0.8 opacity
				drawRoundedRect(img, barX+2, barY, barWidth-4, barHeight, 4, taskColorAlpha)
			}

			currentY += rowHeight
		}
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GeneratePDF creates a PDF document of the Gantt chart
func GeneratePDF(chart *Chart) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// Calculate dimensions (scaled for PDF)
	quarters := calculateQuarters(chart.StartYear, chart.StartQ, chart.EndYear, chart.EndQ)
	_ = len(quarters) // totalQuarters not used in this simplified version

	// Layout constants (in mm)
	headerHeight := 20.0
	rowHeight := 10.0
	quarterWidth := 25.0
	labelWidth := 50.0
	padding := 10.0
	categoryHeaderHeight := 8.0

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.SetXY(padding, padding)
	pdf.Cell(0, 10, chart.Title)

	// Draw quarter headers
	pdf.SetFont("Arial", "B", 9)
	for i, q := range quarters {
		x := padding + labelWidth + float64(i)*quarterWidth
		y := headerHeight

		// Quarter background
		if i%2 == 0 {
			pdf.SetFillColor(232, 232, 232)
		} else {
			pdf.SetFillColor(245, 245, 245)
		}
		pdf.Rect(x, y-8, quarterWidth, 8, "F")

		// Quarter text
		pdf.SetXY(x, y-7)
		pdf.Cell(quarterWidth, 6, fmt.Sprintf("Q%d %d", q.quarter, q.year))
	}

	// Draw categories and tasks
	currentY := headerHeight
	pdf.SetFont("Arial", "", 8)

	for _, cat := range chart.Categories {
		// Category header
		catR, catG, catB := parseColorRGB(cat.Color)
		pdf.SetFillColor(catR, catG, catB)
		pdf.SetAlpha(0.3, "Normal")
		pdf.Rect(padding, currentY, labelWidth, categoryHeaderHeight, "F")
		pdf.SetAlpha(1.0, "Normal")

		pdf.SetFont("Arial", "B", 9)
		pdf.SetXY(padding+2, currentY+2)
		pdf.Cell(labelWidth-4, categoryHeaderHeight-4, cat.Name)
		pdf.SetFont("Arial", "", 8)

		currentY += categoryHeaderHeight

		// Tasks
		for _, task := range cat.Tasks {
			// Task label background
			pdf.SetDrawColor(221, 221, 221)
			pdf.Rect(padding, currentY, labelWidth, rowHeight, "D")

			pdf.SetXY(padding+2, currentY+2)
			pdf.Cell(labelWidth-4, 4, truncate(task.Title, 20))

			// Draw task bar
			startIdx := findQuarterIndex(quarters, task.StartYear, task.StartQ)
			endIdx := findQuarterIndex(quarters, task.EndYear, task.EndQ)

			if startIdx >= 0 && endIdx >= 0 {
				barX := padding + labelWidth + float64(startIdx)*quarterWidth
				barWidth := float64(endIdx-startIdx+1) * quarterWidth
				barY := currentY + 2
				barHeight := rowHeight - 4

				taskColor := task.Color
				if taskColor == "" {
					taskColor = cat.Color
				}

				taskR, taskG, taskB := parseColorRGB(taskColor)
				pdf.SetFillColor(taskR, taskG, taskB)
				pdf.SetAlpha(0.8, "Normal")
				pdf.Rect(barX+1, barY, barWidth-2, barHeight, "F")
				pdf.SetAlpha(1.0, "Normal")
			}

			currentY += rowHeight
		}
	}

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper functions for image drawing
func drawRect(img *image.RGBA, x, y, width, height int, col color.RGBA) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if i >= 0 && i < img.Bounds().Dx() && j >= 0 && j < img.Bounds().Dy() {
				img.Set(i, j, col)
			}
		}
	}
}

func drawRectBorder(img *image.RGBA, x, y, width, height int, col color.RGBA) {
	// Top and bottom
	for i := x; i < x+width; i++ {
		if i >= 0 && i < img.Bounds().Dx() {
			if y >= 0 && y < img.Bounds().Dy() {
				img.Set(i, y, col)
			}
			if y+height >= 0 && y+height < img.Bounds().Dy() {
				img.Set(i, y+height, col)
			}
		}
	}
	// Left and right
	for j := y; j < y+height; j++ {
		if j >= 0 && j < img.Bounds().Dy() {
			if x >= 0 && x < img.Bounds().Dx() {
				img.Set(x, j, col)
			}
			if x+width >= 0 && x+width < img.Bounds().Dx() {
				img.Set(x+width, j, col)
			}
		}
	}
}

func drawRoundedRect(img *image.RGBA, x, y, width, height, radius int, col color.RGBA) {
	// Simplified rounded rectangle - just draw a regular rectangle for now
	drawRect(img, x, y, width, height, col)
}

func drawVerticalLine(img *image.RGBA, x, y1, y2 int, col color.RGBA) {
	for j := y1; j <= y2; j++ {
		if x >= 0 && x < img.Bounds().Dx() && j >= 0 && j < img.Bounds().Dy() {
			img.Set(x, j, col)
		}
	}
}

func parseColor(hexColor string) color.RGBA {
	if hexColor == "" {
		return color.RGBA{100, 149, 237, 255} // Default blue
	}

	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return color.RGBA{100, 149, 237, 255}
	}

	r, _ := strconv.ParseUint(hexColor[0:2], 16, 8)
	g, _ := strconv.ParseUint(hexColor[2:4], 16, 8)
	b, _ := strconv.ParseUint(hexColor[4:6], 16, 8)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func parseColorRGB(hexColor string) (int, int, int) {
	col := parseColor(hexColor)
	return int(col.R), int(col.G), int(col.B)
}
