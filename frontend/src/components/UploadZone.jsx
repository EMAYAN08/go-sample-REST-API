import React, { useRef, useState } from 'react';

export default function UploadZone({ onUpload, isUploading }) {
  const fileInputRef = useRef(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    if (files.length > 0) {
      handleFiles(files);
    }
  };

  const handleFileSelect = (e) => {
    const files = Array.from(e.target.files);
    if (files.length > 0) {
      handleFiles(files);
    }
    // Reset input so same file can be selected again
    e.target.value = '';
  };

  const handleFiles = async (files) => {
    for (const file of files) {
      if (file.type.startsWith('image/')) {
        await onUpload(file);
      }
    }
  };

  const triggerFileInput = () => {
    fileInputRef.current.click();
  };

  return (
    <div
      className={`upload-zone ${isDragging ? 'dragging' : ''} ${isUploading ? 'uploading' : ''}`}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
    >
      {isUploading ? (
        <div className="upload-message">
          <div className="spinner"></div>
          <p>Uploading...</p>
        </div>
      ) : (
        <>
          <div className="upload-icon">📁</div>
          <p>Drag & drop images here</p>
          <span>or</span>
          <button onClick={triggerFileInput}>Choose Files</button>
        </>
      )}
      <input
        ref={fileInputRef}
        type="file"
        multiple
        accept="image/*"
        onChange={handleFileSelect}
        style={{ display: 'none' }}
      />
    </div>
  );
}
