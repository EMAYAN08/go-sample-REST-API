import React from 'react';
import ImageCard from './ImageCard';

export default function ImageGrid({ images, onDelete, onView, searchQuery }) {
  if (images.length === 0) {
    if (searchQuery) {
      return (
        <div className="empty-state">
          <p>No images match '{searchQuery}'</p>
        </div>
      );
    }
    return (
      <div className="empty-state">
        <p>No images yet. Upload your first image!</p>
      </div>
    );
  }

  return (
    <div className="image-grid">
      {images.map((image) => (
        <ImageCard
          key={image.id}
          image={image}
          onDelete={onDelete}
          onView={onView}
        />
      ))}
    </div>
  );
}
