import { useState, useEffect } from 'react';
import './styles/App.css';
import SearchBar from './components/SearchBar';
import FilterBar from './components/FilterBar';
import UploadZone from './components/UploadZone';
import ImageGrid from './components/ImageGrid';
import ImageModal from './components/ImageModal';
import { fetchImages, uploadImage, deleteImage } from './utils/api';

function App() {
  const [images, setImages] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState('date');
  const [sortOrder, setSortOrder] = useState('desc');
  const [selectedImage, setSelectedImage] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Load images on mount
  useEffect(() => {
    loadImages();
  }, []);

  const loadImages = async () => {
    try {
      setLoading(true);
      const data = await fetchImages();
      setImages(data);
      setError(null);
    } catch (err) {
      setError('Failed to load images: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleUpload = async (file) => {
    try {
      setIsUploading(true);
      const newImage = await uploadImage(file);
      setImages([newImage, ...images]);
      setError(null);
    } catch (err) {
      setError('Failed to upload image: ' + err.message);
    } finally {
      setIsUploading(false);
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteImage(id);
      setImages(images.filter((img) => img.id !== id));
      setError(null);
    } catch (err) {
      setError('Failed to delete image: ' + err.message);
    }
  };

  const handleSortChange = (newSortBy, newSortOrder) => {
    setSortBy(newSortBy);
    setSortOrder(newSortOrder);
  };

  // Filter and sort images
  const getFilteredAndSortedImages = () => {
    let filtered = images;

    // Apply search filter
    if (searchQuery) {
      filtered = filtered.filter((img) =>
        img.originalName.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    // Sort
    const sorted = [...filtered].sort((a, b) => {
      let comparison = 0;

      if (sortBy === 'date') {
        const dateA = new Date(a.uploadedAt);
        const dateB = new Date(b.uploadedAt);
        comparison = dateB - dateA; // Newest first by default
      } else if (sortBy === 'size') {
        comparison = b.size - a.size; // Largest first by default
      }

      // Reverse if ascending
      return sortOrder === 'asc' ? -comparison : comparison;
    });

    return sorted;
  };

  const filteredImages = getFilteredAndSortedImages();

  return (
    <div className="app">
      <header className="header">
        <h1>Image Gallery</h1>
        <SearchBar searchQuery={searchQuery} onSearchChange={setSearchQuery} />
      </header>

      <FilterBar
        sortBy={sortBy}
        sortOrder={sortOrder}
        onSortChange={handleSortChange}
      />

      {error && <div className="error-message">{error}</div>}

      <UploadZone onUpload={handleUpload} isUploading={isUploading} />

      {loading ? (
        <div className="loading">Loading images...</div>
      ) : (
        <ImageGrid
          images={filteredImages}
          onDelete={handleDelete}
          onView={setSelectedImage}
          searchQuery={searchQuery}
        />
      )}

      <ImageModal image={selectedImage} onClose={() => setSelectedImage(null)} />
    </div>
  );
}

export default App;
