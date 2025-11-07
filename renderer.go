package main

import (
	"bytes"
	"fmt"
	"strings"
)

// GenerateSVG creates an SVG representation of the Gantt chart
func GenerateSVG(chart *Chart) (string, error) {
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

	var buf bytes.Buffer

	// SVG header
	buf.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, width, height))
	buf.WriteString(`<defs><style>.title{font:bold 20px sans-serif;fill:#333}.header{font:bold 12px sans-serif;fill:#555}.label{font:12px sans-serif;fill:#333}.category{font:bold 14px sans-serif;fill:#222}.desc{font:10px sans-serif;fill:#666}</style></defs>`)

	// Background
	buf.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="#fafafa"/>`, width, height))

	// Title
	buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="title">%s</text>`, padding, padding+20, escapeXML(chart.Title)))

	// Draw quarter headers
	for i, q := range quarters {
		x := padding + labelWidth + i*quarterWidth
		y := headerHeight

		// Quarter background
		color := "#e8e8e8"
		if i%2 == 1 {
			color = "#f5f5f5"
		}
		buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" stroke="#ccc" stroke-width="1"/>`,
			x, y-30, quarterWidth, 30, color))

		// Quarter text
		buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="header" text-anchor="middle">Q%d %d</text>`,
			x+quarterWidth/2, y-10, q.quarter, q.year))

		// Vertical grid lines
		buf.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#ddd" stroke-width="1"/>`,
			x, y, x, height-padding))
	}

	// Draw categories and tasks
	currentY := headerHeight
	for _, cat := range chart.Categories {
		// Category header
		buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="0.3"/>`,
			padding, currentY, labelWidth, categoryHeaderHeight, cat.Color))
		buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="category">%s</text>`,
			padding+10, currentY+22, escapeXML(cat.Name)))

		// Category background span
		buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="0.05"/>`,
			padding+labelWidth, currentY, totalQuarters*quarterWidth, categoryHeaderHeight, cat.Color))

		currentY += categoryHeaderHeight

		// Tasks
		for _, task := range cat.Tasks {
			// Task label background
			buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#fff" stroke="#ddd" stroke-width="1"/>`,
				padding, currentY, labelWidth, rowHeight))

			// Task title
			buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="label">%s</text>`,
				padding+10, currentY+18, escapeXML(truncate(task.Title, 25))))

			// Task description
			if task.Description != "" {
				buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="desc">%s</text>`,
					padding+10, currentY+32, escapeXML(truncate(task.Description, 30))))
			}

			// Draw task bar
			startIdx := findQuarterIndex(quarters, task.StartYear, task.StartQ)
			endIdx := findQuarterIndex(quarters, task.EndYear, task.EndQ)

			if startIdx >= 0 && endIdx >= 0 {
				barX := padding + labelWidth + startIdx*quarterWidth
				barWidth := (endIdx - startIdx + 1) * quarterWidth
				barY := currentY + 8
				barHeight := rowHeight - 16

				taskColor := task.Color
				if taskColor == "" {
					taskColor = cat.Color
				}

				buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="4" opacity="0.8"/>`,
					barX+2, barY, barWidth-4, barHeight, taskColor))
				buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="2" rx="4"/>`,
					barX+2, barY, barWidth-4, barHeight, darken(taskColor)))
			}

			currentY += rowHeight
		}
	}

	buf.WriteString(`</svg>`)
	return buf.String(), nil
}

type quarterInfo struct {
	year    int
	quarter int
}

func calculateQuarters(startYear, startQ, endYear, endQ int) []quarterInfo {
	var quarters []quarterInfo

	for year := startYear; year <= endYear; year++ {
		startQuarter := 1
		endQuarter := 4

		if year == startYear {
			startQuarter = startQ
		}
		if year == endYear {
			endQuarter = endQ
		}

		for q := startQuarter; q <= endQuarter; q++ {
			quarters = append(quarters, quarterInfo{year: year, quarter: q})
		}
	}

	return quarters
}

func findQuarterIndex(quarters []quarterInfo, year, quarter int) int {
	for i, q := range quarters {
		if q.year == year && q.quarter == quarter {
			return i
		}
	}
	return -1
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func darken(color string) string {
	// Simple darkening - in production, you'd want proper color manipulation
	if color == "" {
		return "#333"
	}
	return color
}
