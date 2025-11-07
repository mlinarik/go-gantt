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
	baseRowHeight := 40
	quarterWidth := 120
	labelWidth := 200
	padding := 20
	categoryHeaderHeight := 35
	titleLineHeight := 14
	descLineHeight := 12
	verticalPaddingPerTask := 8

	// Helper to wrap text by approx character count
	wrapText := func(s string, maxChars int) []string {
		if s == "" {
			return []string{}
		}
		parts := strings.Fields(s)
		var lines []string
		line := ""
		for _, w := range parts {
			if len((line + " " + w)) <= maxChars {
				if line == "" {
					line = w
				} else {
					line = line + " " + w
				}
			} else {
				if line != "" {
					lines = append(lines, line)
				}
				line = w
			}
		}
		if line != "" {
			lines = append(lines, line)
		}
		return lines
	}

	// Calculate dynamic heights per task and category
	totalCategoryHeadersHeight := 0
	totalTaskHeight := 0
	perTaskHeights := make(map[string]int)
	perCategoryHeights := make(map[string]int)
	for _, cat := range chart.Categories {
		catNameLines := wrapText(cat.Name, 30)
		catH := categoryHeaderHeight
		if len(catNameLines) > 1 {
			catH = 18 + len(catNameLines)*14 // base + lines * lineheight
		}
		perCategoryHeights[cat.ID] = catH
		totalCategoryHeadersHeight += catH

		for _, task := range cat.Tasks {
			titleLines := wrapText(task.Title, 28)
			descLines := wrapText(task.Description, 36)
			h := baseRowHeight
			calc := len(titleLines)*titleLineHeight + len(descLines)*descLineHeight + verticalPaddingPerTask
			if calc > h {
				h = calc
			}
			perTaskHeights[task.ID] = h
			totalTaskHeight += h
		}
	}

	width := labelWidth + totalQuarters*quarterWidth + padding*2
	height := headerHeight + totalCategoryHeadersHeight + totalTaskHeight + padding*2

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
		// Category header - dynamic height
		catH := perCategoryHeights[cat.ID]
		if catH == 0 {
			catH = categoryHeaderHeight
		}

		buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="0.3"/>`,
			padding, currentY, labelWidth, catH, cat.Color))

		// Category name with wrapping
		catLines := wrapText(cat.Name, 30)
		if len(catLines) <= 1 {
			buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="category">%s</text>`,
				padding+10, currentY+22, escapeXML(cat.Name)))
		} else {
			buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="category">`, padding+10, currentY+18))
			for i, ln := range catLines {
				dy := 4
				if i > 0 {
					dy = 14
				}
				buf.WriteString(fmt.Sprintf(`<tspan x="%d" dy="%d">%s</tspan>`, padding+10, dy, escapeXML(ln)))
			}
			buf.WriteString(`</text>`)
		}

		// Category background span
		buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="0.05"/>`,
			padding+labelWidth, currentY, totalQuarters*quarterWidth, catH, cat.Color))

		currentY += catH

		// Tasks
		for _, task := range cat.Tasks {
			// dynamic task height
			h := perTaskHeights[task.ID]
			if h == 0 {
				h = baseRowHeight
			}
			// Task label background
			buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#fff" stroke="#ddd" stroke-width="1"/>`,
				padding, currentY, labelWidth, h))

			// Title and description wrapped
			titleLines := wrapText(task.Title, 28)
			descLines := wrapText(task.Description, 36)
			textY := currentY + 14
			if len(titleLines) > 0 {
				buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="label">`, padding+10, textY))
				for i, ln := range titleLines {
					dy := 0
					if i > 0 {
						dy = titleLineHeight
					}
					buf.WriteString(fmt.Sprintf(`<tspan x="%d" dy="%d">%s</tspan>`, padding+10, dy, escapeXML(ln)))
				}
				buf.WriteString(`</text>`)
				textY += len(titleLines) * titleLineHeight
			}
			if len(descLines) > 0 {
				buf.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="desc">`, padding+10, textY+4))
				for i, ln := range descLines {
					dy := 0
					if i > 0 {
						dy = descLineHeight
					}
					buf.WriteString(fmt.Sprintf(`<tspan x="%d" dy="%d">%s</tspan>`, padding+10, dy, escapeXML(ln)))
				}
				buf.WriteString(`</text>`)
			}

			// Draw task bar
			startIdx := findQuarterIndex(quarters, task.StartYear, task.StartQ)
			endIdx := findQuarterIndex(quarters, task.EndYear, task.EndQ)

			if startIdx >= 0 && endIdx >= 0 {
				barX := padding + labelWidth + startIdx*quarterWidth
				barWidth := (endIdx - startIdx + 1) * quarterWidth
				barY := currentY + 8
				barHeight := h - 16

				taskColor := task.Color
				if taskColor == "" {
					taskColor = cat.Color
				}

				if barHeight < 12 {
					barHeight = 12
				}

				buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="4" opacity="0.8"/>`,
					barX+2, barY, barWidth-4, barHeight, taskColor))
				buf.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="none" stroke="%s" stroke-width="2" rx="4"/>`,
					barX+2, barY, barWidth-4, barHeight, darken(taskColor)))
			}

			currentY += h
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
