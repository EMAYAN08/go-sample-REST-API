import React from 'react';

export default function ImageCard({ image, onDelete, onView }) {
  // Format size in MB
  const formatSize = (bytes) => {
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  };

  // Format date
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  return (
    <div className="glass-card">
      <div className="card-image">
        <img
          src={`/api/images/file/${image.id}`}
          alt={image.originalName}
          onError={(e) => {
            e.target.style.display = 'none';
            e.target.nextSibling.style.display = 'flex';
          }}
        />
        <div className="image-placeholder" style={{ display: 'none' }}>
          <span>📷</span>
          <p>{image.originalName}</p>
        </div>
      </div>
      <div className="card-info">
        <div className="filename" title={image.originalName}>
          {image.originalName}
        </div>
        <div className="metadata">
          {formatSize(image.size)} • {formatDate(image.uploadedAt)}
        </div>
      </div>
      <div className="card-actions">
        <button className="btn-view" onClick={() => onView(image)}>
          View
        </button>
        <button className="btn-delete" onClick={() => onDelete(image.id)}>
          Delete
        </button>
      </div>
    </div>
  );
}
