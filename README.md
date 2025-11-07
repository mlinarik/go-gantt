# Go Gantt Chart Application

A containerized Go application for creating custom Gantt charts with an interactive web UI. Export your charts as PDF, PNG, or SVG.

## Features

- 📊 **Interactive Web UI** - Create and manage Gantt charts through a modern web interface
- 📅 **Quarter-based Timeline** - Select timeframes by year and quarter
- 🎨 **Category Grouping** - Organize tasks into color-coded categories
- ✏️ **Drag & Drop** - Adjust task timelines with intuitive drag-and-drop (UI ready)
- 📝 **Custom Titles & Notes** - Add titles and descriptions to individual tasks
- 💾 **Multiple Export Formats** - Export charts as PDF, PNG, or SVG
- 🐳 **Containerized** - Easy deployment with Docker

## Quick Start

### Using Docker

1. **Build the container:**
   ```bash
   docker build -t go-ghant .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 go-ghant
   ```

3. **Access the application:**
   Open your browser to `http://localhost:8080`

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Run the application:**
   ```bash
   go run .
   ```

3. **Access the application:**
   Open your browser to `http://localhost:8080`

## Usage

### Creating a Chart

1. Click **"New Chart"** to start
2. Set your chart title and timeframe (start/end year and quarter)
3. Add categories with the **"+ Add Category"** button
4. Add tasks to each category with the **"+ Task"** button
5. Fill in task details including:
   - Task title
   - Description (optional)
   - Start and end quarters
   - Custom color (optional, defaults to category color)

### Exporting Charts

1. Click **"Save Chart"** to persist your chart
2. Use the export buttons in the sidebar:
   - **Export as SVG** - Vector format, ideal for further editing
   - **Export as PNG** - Raster image format
   - **Export as PDF** - Document format for sharing

### Managing Tasks

- **Edit** - Click the "Edit" button on any task or category
- **Delete** - Click the "Delete" button to remove items
- **Reorder** - Drag tasks to adjust their position (visual feedback in UI)

## API Endpoints

The application provides a REST API:

- `GET /api/charts` - List all charts
- `POST /api/charts` - Create a new chart
- `GET /api/charts/{id}` - Get a specific chart
- `PUT /api/charts/{id}` - Update a chart
- `DELETE /api/charts/{id}` - Delete a chart
- `GET /api/charts/{id}/export/svg` - Export as SVG
- `GET /api/charts/{id}/export/png` - Export as PNG
- `GET /api/charts/{id}/export/pdf` - Export as PDF

## Configuration

Set the `PORT` environment variable to change the default port (8080):

```bash
docker run -p 3000:3000 -e PORT=3000 go-ghant
```

## Data Persistence

Charts are saved to `charts.json` in the application directory. To persist data:

```bash
docker run -p 8080:8080 -v $(pwd)/data:/root go-ghant
```

## Technology Stack

- **Backend:** Go with Gorilla Mux router
- **Frontend:** Vanilla JavaScript, HTML, CSS
- **Export:** gofpdf for PDF generation, native image/png for PNG
- **Container:** Docker with multi-stage builds

## Project Structure

```
go-ghant/
├── main.go           # HTTP server and API handlers
├── models.go         # Data structures and storage
├── renderer.go       # SVG chart generation
├── export.go         # PDF and PNG export functionality
├── go.mod            # Go module dependencies
├── Dockerfile        # Container configuration
└── static/           # Frontend files
    ├── index.html    # Web UI
    ├── style.css     # Styles
    └── app.js        # JavaScript application
```

## License

MIT License - feel free to use this project for your needs.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## Future Enhancements

- [ ] Database storage (PostgreSQL, SQLite)
- [ ] User authentication and multiple users
- [ ] Real-time collaboration
- [ ] More chart customization options
- [ ] Milestone markers
- [ ] Resource allocation tracking
- [ ] Import/export to other formats (Excel, MS Project)
