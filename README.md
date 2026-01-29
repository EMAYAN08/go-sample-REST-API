# Image Gallery - Full Stack Application

A beautiful full-stack image gallery application with glassomorphic UI design. Upload, view, search, and manage your images with a modern, responsive interface.

## Features

- **Upload Images**: Drag & drop or click to upload images (supports all image formats, no size limit)
- **Glassomorphic Grid View**: Beautiful grid layout with glassmorphism design
- **Search**: Real-time search by filename
- **Sort & Filter**: Sort images by date or size, ascending or descending
- **Full-Size Preview**: Click to view images in a full-screen modal
- **Delete Images**: Remove unwanted images
- **Persistent Storage**: Images and metadata survive server restarts
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile

## Tech Stack

**Backend:**
- Go 1.22.3
- gorilla/mux for routing
- google/uuid for ID generation
- JSON file-based metadata storage
- Local filesystem for image storage

**Frontend:**
- React 18 (built with Vite)
- Modern CSS with glassomorphism effects
- Responsive grid layout
- Native Fetch API for HTTP requests

## Project Structure

```
go-sample-REST-API/
├── main.go                 # Go backend server
├── go.mod                  # Go dependencies
├── go.sum                  # Go dependency checksums
├── uploads/                # Image files storage (auto-created)
├── data/                   # Metadata storage (auto-created)
│   └── metadata.json       # Image metadata (auto-created)
├── frontend/               # React frontend
│   ├── src/
│   │   ├── App.jsx         # Main app component
│   │   ├── components/     # React components
│   │   ├── styles/         # CSS styles
│   │   └── utils/          # API client
│   ├── build/              # Production build (after npm run build)
│   ├── package.json        # Node dependencies
│   └── vite.config.js      # Vite configuration
└── README.md               # This file
```

## Getting Started

### Prerequisites

- Go 1.22+ installed
- Node.js 16+ and npm installed

### Installation & Running

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-sample-REST-API
   ```

2. **Install frontend dependencies and build**
   ```bash
   cd frontend
   npm install
   npm run build
   cd ..
   ```

3. **Run the server**
   ```bash
   go run main.go
   ```

4. **Open your browser**
   Navigate to: `http://localhost:6000`

The server will:
- Create `uploads/` and `data/` directories automatically
- Serve the React frontend from `frontend/build/`
- Provide REST API endpoints at `/api/*`
- Run on port 6000

## API Endpoints

All API endpoints use the `/api` prefix:

### GET /api/images
Get all images with metadata, sorted by upload date (newest first)

**Response:** `200 OK`
```json
[
  {
    "id": "uuid-string",
    "filename": "uuid_originalname.jpg",
    "originalName": "photo.jpg",
    "size": 1048576,
    "mimeType": "image/jpeg",
    "uploadedAt": "2024-01-01T12:00:00Z",
    "path": "./uploads/uuid_originalname.jpg"
  }
]
```

### GET /api/images/{id}
Get single image metadata by ID

**Response:** `200 OK` or `404 Not Found`

### POST /api/images/upload
Upload a new image

**Request:** `multipart/form-data` with "image" field

**Response:** `201 Created`
```json
{
  "id": "uuid-string",
  "filename": "uuid_originalname.jpg",
  "originalName": "photo.jpg",
  "size": 1048576,
  "mimeType": "image/jpeg",
  "uploadedAt": "2024-01-01T12:00:00Z",
  "path": "./uploads/uuid_originalname.jpg"
}
```

**Errors:**
- `400 Bad Request`: No file provided or failed to read file
- `500 Internal Server Error`: Failed to save image or metadata

### DELETE /api/images/{id}
Delete an image and its file

**Response:** `200 OK`
```json
{
  "message": "Image deleted successfully"
}
```

**Errors:**
- `404 Not Found`: Image not found
- `500 Internal Server Error`: Failed to delete or update metadata

### GET /api/images/file/{id}
Serve the actual image file for viewing

**Response:** `200 OK` with image data and appropriate Content-Type header

**Errors:**
- `404 Not Found`: Image or file not found

## Development

### Running in Development Mode

For hot-reload during frontend development:

**Terminal 1 - Backend:**
```bash
go run main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Frontend dev server runs on `http://localhost:5173` (or similar) with API proxy to backend.

### Building for Production

```bash
cd frontend
npm run build
cd ..
go run main.go
```

The Go server will serve both the API and the built React app from a single port (6000).

## Data Persistence

- **Images**: Stored in `./uploads/` directory
- **Metadata**: Stored in `./data/metadata.json`
- Both persist across server restarts
- On first run, directories and metadata file are created automatically
- If metadata.json is corrupted, the server starts with an empty array

## Design Features

### Glassomorphism UI
- Semi-transparent backgrounds with backdrop blur
- Subtle borders and shadows
- Gradient background (purple to violet)
- Smooth animations and transitions

### Responsive Breakpoints
- **Desktop (1200px+)**: 4-column grid
- **Tablet (768-1199px)**: 3-column grid
- **Mobile (480-767px)**: 2-column grid
- **Small (<480px)**: 1-column grid

## Error Handling

**Backend:**
- Graceful handling of missing directories (auto-created)
- Corrupted metadata recovery (starts fresh)
- File operation error logging
- Proper HTTP status codes

**Frontend:**
- Error messages displayed to user
- Failed image loads show placeholder
- Network error handling
- Loading states for async operations

## Browser Support

Works in all modern browsers:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+

## License

MIT

## Contributing

Feel free to submit issues and pull requests!
