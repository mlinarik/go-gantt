# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-11-06

### Added
- Initial release of Go Gantt Chart Application
- Interactive web UI for creating Gantt charts
- Quarter-based timeline with year and quarter selection
- Category management with color coding
- Task management with titles, descriptions, and timelines
- Drag-and-drop UI preparation (visual feedback ready)
- Export functionality for SVG, PNG, and PDF formats
- REST API for chart CRUD operations
- Docker containerization support
- Docker Compose configuration
- File-based JSON storage for charts
- Example chart data for quick start
- Comprehensive documentation (README, ARCHITECTURE)
- PowerShell quick start script for Windows
- Makefile for common operations
- Responsive web interface
- Modal dialogs for category and task editing
- Real-time chart preview
- Client-side SVG generation
- Color picker for categories and tasks
- Custom task colors (optional, inherits from category)

### Features
- **Chart Management**: Create, read, update, delete charts
- **Timeline Control**: Select start/end year and quarter (Q1-Q4)
- **Categories**: Group tasks into colored categories
- **Tasks**: Add detailed tasks with titles, descriptions, and timelines
- **Visual Editor**: Interactive UI for building charts
- **Multiple Exports**: Download as SVG, PNG, or PDF
- **Containerized**: Easy deployment with Docker
- **Persistent Storage**: Charts saved to JSON file
- **Example Data**: Pre-loaded example chart for demonstration

### Technical
- Go 1.21+ backend
- Gorilla Mux router for HTTP handling
- gofpdf for PDF generation
- Native Go image/png for PNG export
- Vanilla JavaScript frontend (no frameworks)
- Responsive CSS Grid and Flexbox layout
- SVG-based chart rendering
- UUID-based resource identification
- Multi-stage Docker builds for optimization

### Documentation
- README with quick start guide
- Architecture documentation
- API endpoint documentation
- Deployment instructions
- Docker usage examples
- Contributing guidelines
- Troubleshooting guide

## [Unreleased]

### Planned
- Database backend (PostgreSQL/MongoDB)
- User authentication and authorization
- Multi-user support
- Real-time collaboration
- Full drag-and-drop task repositioning
- Task dependencies visualization
- Milestone markers
- Resource allocation tracking
- Chart templates
- Import/Export to Excel and MS Project formats
- Version history and rollback
- Comments and discussions on tasks
- Email notifications
- Webhook integrations
- Search and filter functionality
- Chart sharing with permissions
- Custom themes and branding
- Mobile-responsive improvements
- Unit and integration tests
- API documentation (OpenAPI/Swagger)
- Monitoring and metrics
- Performance optimizations for large charts
