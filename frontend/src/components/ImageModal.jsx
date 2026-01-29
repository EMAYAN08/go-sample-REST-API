import React, { useEffect } from 'react';

export default function ImageModal({ image, onClose }) {
  useEffect(() => {
    if (!image) return;

    const handleEscape = (e) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [image, onClose]);

  if (!image) return null;

  // Format size in MB
  const formatSize = (bytes) => {
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  };

  // Format date
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'long',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="close-btn" onClick={onClose}>
          ✕
        </button>
        <img
          src={`/api/images/file/${image.id}`}
          alt={image.originalName}
        />
        <div className="modal-info">
          <h3>{image.originalName}</h3>
          <p>{formatSize(image.size)} • {formatDate(image.uploadedAt)}</p>
        </div>
      </div>
    </div>
  );
}
