# Getting Started with Go Gantt Chart Application

This guide will help you get the Gantt Chart application up and running in just a few minutes.

## Quick Start (Windows)

The fastest way to get started on Windows:

```powershell
.\start.ps1
```

This script will:
1. Check for Docker and Go installations
2. Give you options to run with Docker or locally
3. Build and start the application automatically
4. Open your browser to the application

## Option 1: Docker (Recommended)

### Prerequisites
- Docker Desktop installed and running
- No other software required!

### Steps

1. **Build the container:**
   ```powershell
   docker build -t go-ghant .
   ```

2. **Run the container:**
   ```powershell
   docker run -p 8080:8080 -v ${PWD}/data:/root go-ghant
   ```

3. **Open your browser:**
   Navigate to http://localhost:8080

4. **Start creating charts!**

### Using Docker Compose (Even Easier)

```powershell
docker-compose up -d
```

To stop:
```powershell
docker-compose down
```

## Option 2: Local Go Installation

### Prerequisites
- Go 1.21 or higher installed
- Git (to clone the repository)

### Steps

1. **Download dependencies:**
   ```powershell
   go mod download
   ```

2. **Run the application:**
   ```powershell
   go run .
   ```

3. **Open your browser:**
   Navigate to http://localhost:8080

## First Steps with the Application

### 1. Create Your First Chart

1. Click **"New Chart"** button in the header
2. Enter a title for your chart (e.g., "Project Roadmap 2024")
3. Set your timeline:
   - Start Year: 2024
   - Start Quarter: Q1
   - End Year: 2025
   - End Quarter: Q4

### 2. Add Categories

Categories help organize your tasks into logical groups.

1. Click **"+ Add Category"** in the sidebar
2. Enter a name (e.g., "Development", "Marketing", "Research")
3. Choose a color for the category
4. Click **"Save"**

### 3. Add Tasks to Categories

1. Click the **"+ Task"** button next to a category
2. Fill in the task details:
   - **Title**: Brief name for the task
   - **Description**: Additional details (optional)
   - **Start Quarter**: When the task begins
   - **End Quarter**: When the task ends
   - **Color**: Custom color (optional, uses category color by default)
3. Click **"Save"**

### 4. View Your Chart

As you add categories and tasks, the chart preview automatically updates on the right side of the screen showing your Gantt chart with:
- Quarter-based timeline at the top
- Categories with colored headers
- Tasks as horizontal bars spanning their timeframes

### 5. Save Your Chart

Click **"Save Chart"** in the header to persist your work. The chart will be saved to a JSON file.

### 6. Export Your Chart

Once saved, you can export your chart in multiple formats:

1. **SVG** - Scalable vector format, best for editing
   - Click "Export as SVG"
   - Opens in vector graphics software
   
2. **PNG** - Raster image format, best for presentations
   - Click "Export as PNG"
   - Ready to insert into documents
   
3. **PDF** - Document format, best for sharing
   - Click "Export as PDF"
   - Professional format for stakeholders

## Example Chart

The application includes an example chart to help you get started:

```powershell
# Copy the example to use it
cp charts.example.json charts.json
```

Then restart the application to see the example chart.

## Common Tasks

### Edit a Category
1. Click **"Edit"** next to the category name
2. Modify the name or color
3. Click **"Save"**

### Edit a Task
1. Click **"Edit"** next to the task
2. Modify any field
3. Click **"Save"**

### Delete Items
Click the **"Delete"** button next to any category or task. Categories will delete all their tasks.

### Adjust Timeline
Change the start/end year and quarters in the sidebar to adjust the overall chart timeframe. Tasks outside the range will still exist but won't be visible in the chart.

### Change Chart Title
Simply edit the "Chart Title" field at the top of the sidebar.

## Tips for Best Results

### Planning Your Chart
1. **Start with categories** - Think about major phases or departments
2. **Use meaningful colors** - Different colors for different types of work
3. **Keep tasks focused** - One clear objective per task
4. **Add descriptions** - Future you will thank you for the details

### Timeline Management
1. **Start broad** - You can always narrow the view later
2. **Quarterly planning** - Quarters are ideal for strategic planning
3. **Buffer time** - Leave gaps between dependent tasks

### Visual Design
1. **Limit categories** - 4-6 categories maximum for readability
2. **Color contrast** - Choose colors that are easy to distinguish
3. **Task length** - Very short tasks (1 quarter) might be hard to see

### Export Strategy
1. **SVG for editing** - Keep a copy you can modify later
2. **PDF for stakeholders** - Professional format for meetings
3. **PNG for presentations** - Easy to insert into slides

## Troubleshooting

### Application won't start
- Check if port 8080 is available
- Try a different port: `docker run -p 3000:3000 -e PORT=3000 go-ghant`

### Can't save chart
- Check file permissions in the application directory
- With Docker, ensure volume is mounted: `-v ${PWD}/data:/root`

### Chart doesn't display
- Ensure you have at least one category with one task
- Check that task dates fall within the chart timeframe
- Try refreshing the browser

### Export fails
- Make sure you've saved the chart first
- Check browser console for errors (F12)
- Verify the server is running

## Next Steps

### Learn More
- Read the [README.md](README.md) for detailed documentation
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for technical details
- View [CHANGELOG.md](CHANGELOG.md) for version history

### Customize
- Modify `static/style.css` to change the look and feel
- Edit `static/app.js` to add client-side features
- Extend the Go backend to add new API endpoints

### Deploy
- Use Docker for production deployments
- Set up a reverse proxy (nginx, Traefik) for HTTPS
- Configure persistent storage with volume mounts
- Consider adding authentication for multi-user scenarios

## Need Help?

- Check the troubleshooting section above
- Review the documentation files
- Open an issue on GitHub
- Check existing issues for solutions

## Enjoy Creating Gantt Charts! 🎉

You're now ready to create professional Gantt charts for your projects. Happy planning!
