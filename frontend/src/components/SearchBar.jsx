import React from 'react';

export default function SearchBar({ searchQuery, onSearchChange }) {
  return (
    <div className="search-bar">
      <span className="search-icon">🔍</span>
      <input
        type="text"
        placeholder="Search images..."
        value={searchQuery}
        onChange={(e) => onSearchChange(e.target.value)}
      />
      {searchQuery && (
        <button className="clear-btn" onClick={() => onSearchChange('')}>
          ✕
        </button>
      )}
    </div>
  );
}
