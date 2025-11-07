# Go Gantt Chart Application - Technical Overview

## Architecture

### Backend (Go)
- **HTTP Server**: Gorilla Mux router handling REST API requests
- **Data Storage**: JSON file-based persistence (charts.json)
- **Chart Rendering**: SVG generation with quarter-based timeline
- **Export Engine**: Multi-format export (SVG, PNG, PDF)

### Frontend (Vanilla JavaScript)
- **SPA Architecture**: Single-page application with dynamic rendering
- **State Management**: Client-side state with reactive updates
- **UI Components**: Modal dialogs, drag-and-drop interface
- **Real-time Preview**: Client-side SVG generation for instant feedback

## Data Model

```go
Chart {
    ID: string (UUID)
    Title: string
    StartYear: int
    StartQuarter: int (1-4)
    EndYear: int
    EndQuarter: int (1-4)
    Categories: []Category
    CreatedAt: timestamp
    UpdatedAt: timestamp
}

Category {
    ID: string (UUID)
    Name: string
    Color: string (hex color)
    Tasks: []Task
}

Task {
    ID: string (UUID)
    Title: string
    Description: string
    StartYear: int
    StartQuarter: int (1-4)
    EndYear: int
    EndQuarter: int (1-4)
    Color: string (optional, hex color)
}
```

## API Endpoints

### Chart Management
- `GET /api/charts` - List all charts
- `POST /api/charts` - Create new chart
- `GET /api/charts/{id}` - Get specific chart
- `PUT /api/charts/{id}` - Update chart
- `DELETE /api/charts/{id}` - Delete chart

### Export
- `GET /api/charts/{id}/export/svg` - Download SVG
- `GET /api/charts/{id}/export/png` - Download PNG
- `GET /api/charts/{id}/export/pdf` - Download PDF

### Static Files
- `GET /` - Serve index.html and static assets

## File Structure

```
go-ghant/
├── main.go              # HTTP server and request handlers
├── models.go            # Data structures and storage logic
├── renderer.go          # SVG chart generation
├── export.go            # PDF and PNG export functionality
├── go.mod               # Go dependencies
├── go.sum               # Dependency checksums
├── Dockerfile           # Container configuration
├── docker-compose.yml   # Docker Compose setup
├── Makefile             # Build automation
├── start.ps1            # Windows quick start script
├── .gitignore           # Git ignore rules
├── README.md            # User documentation
├── ARCHITECTURE.md      # This file
├── charts.example.json  # Example data
└── static/              # Frontend files
    ├── index.html       # Main HTML page
    ├── style.css        # Styling
    └── app.js           # JavaScript application logic
```

## Chart Rendering

### SVG Generation Process
1. Calculate quarter range from start/end parameters
2. Determine chart dimensions based on:
   - Number of quarters (width)
   - Number of categories and tasks (height)
3. Render components in layers:
   - Background and title
   - Quarter headers with alternating colors
   - Grid lines for visual separation
   - Category headers with transparent color overlay
   - Task labels and descriptions
   - Task bars with rounded corners and opacity

### Layout Constants
- Header height: 80px
- Row height: 40px
- Quarter width: 120px
- Label width: 200px
- Category header height: 35px
- Padding: 20px

## Export Formats

### SVG
- Native format, directly from renderer
- Scalable vector graphics
- Editable in vector graphics software
- Smallest file size

### PNG
- Raster image format
- Generated from SVG-like rendering to image.RGBA
- Fixed resolution based on chart dimensions
- Good for presentations and documentation

### PDF
- Document format using gofpdf library
- Landscape A4 orientation
- Scaled dimensions for print compatibility
- Best for formal documentation and sharing

## Frontend Architecture

### State Management
- `currentChart`: Active chart being edited
- `currentCategoryId`: Category being modified
- `currentTaskId`: Task being modified
- `editingCategory`: Category in edit mode
- `editingTask`: Task in edit mode

### Key Functions
- `createNewChart()`: Initialize empty chart
- `updateUI()`: Refresh all UI components
- `renderCategories()`: Display category list
- `renderTasks()`: Display task list for category
- `updatePreview()`: Generate and display chart SVG
- `saveChart()`: Persist chart via API
- `exportChart(format)`: Download chart in specified format

### Modal System
Two modal dialogs for data entry:
1. Category Modal: Name, color picker
2. Task Modal: Title, description, timeline, color

### Event Handling
- Form inputs trigger immediate chart updates
- Modal save buttons validate and update data
- Export buttons trigger API calls with file download
- Drag events prepared for task reordering (visual only)

## Deployment Options

### Docker (Recommended)
```bash
docker build -t go-ghant .
docker run -p 8080:8080 -v ./data:/root go-ghant
```

### Docker Compose
```bash
docker-compose up -d
```

### Local Development
```bash
go mod download
go run .
```

### Windows Quick Start
```powershell
.\start.ps1
```

## Data Persistence

### Storage Mechanism
- Charts stored in `charts.json` in application directory
- Automatic save on create/update/delete operations
- JSON format for human readability and easy debugging

### Volume Mounting
To persist data across container restarts:
```bash
docker run -v $(pwd)/data:/root go-ghant
```

## Security Considerations

### Current State (Development)
- No authentication/authorization
- Single-user mode
- Local file storage only
- No input sanitization on server (client-side only)

### Production Recommendations
1. Add user authentication (JWT, OAuth)
2. Implement role-based access control
3. Use database (PostgreSQL, MongoDB)
4. Add input validation and sanitization
5. Enable HTTPS/TLS
6. Implement rate limiting
7. Add CORS configuration
8. Sanitize file uploads/downloads

## Performance Optimization

### Current Optimizations
- Static file serving via http.FileServer
- In-memory chart storage with file persistence
- Client-side chart preview (no server round-trip)
- SVG format for efficient vector rendering

### Future Improvements
1. Database indexing for large datasets
2. Chart caching layer
3. Lazy loading for large charts
4. WebSocket for real-time updates
5. Worker pool for export operations
6. CDN for static assets

## Browser Compatibility

### Tested Browsers
- Chrome/Edge (Chromium-based)
- Firefox
- Safari

### Required Features
- ES6 JavaScript
- CSS Grid and Flexbox
- Fetch API
- SVG rendering
- HTML5 input types (color, number)

## Future Enhancements

### Planned Features
1. **Multi-user Support**: User accounts and permissions
2. **Database Backend**: PostgreSQL or MongoDB
3. **Real-time Collaboration**: Multiple users editing same chart
4. **Advanced Drag & Drop**: Full task repositioning
5. **Dependencies**: Link tasks with dependencies
6. **Milestones**: Add milestone markers
7. **Resources**: Track resource allocation
8. **Templates**: Predefined chart templates
9. **Import/Export**: Excel, MS Project, CSV formats
10. **Notifications**: Email/webhook notifications
11. **Version History**: Chart versioning and rollback
12. **Comments**: Task-level comments and discussions

### Technical Debt
1. Add comprehensive unit tests
2. Add integration tests
3. Improve error handling
4. Add logging infrastructure
5. Implement proper validation
6. Add API documentation (OpenAPI/Swagger)
7. Optimize chart rendering for large datasets
8. Add monitoring and metrics

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Docker (optional)
- Git

### Local Development Workflow
```bash
# Clone repository
git clone <repository-url>
cd go-ghant

# Install dependencies
go mod download

# Run tests (when implemented)
go test ./...

# Run locally
go run .

# Build binary
go build -o go-ghant

# Run binary
./go-ghant
```

### Docker Development Workflow
```bash
# Build image
docker build -t go-ghant .

# Run container
docker run -p 8080:8080 go-ghant

# View logs
docker logs <container-id>

# Stop container
docker stop <container-id>
```

## Troubleshooting

### Common Issues

**Port already in use**
```bash
# Change port with environment variable
docker run -p 3000:3000 -e PORT=3000 go-ghant
```

**Data not persisting**
```bash
# Mount volume for persistence
docker run -v $(pwd)/data:/root go-ghant
```

**Module import errors**
```bash
# Update dependencies
go mod tidy
go mod download
```

## Contributing

### Code Style
- Follow Go standard formatting (gofmt)
- Use meaningful variable names
- Comment exported functions
- Keep functions focused and small

### Pull Request Process
1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Run `go fmt` and `go vet`
5. Submit pull request with description

## License

MIT License - See LICENSE file for details
