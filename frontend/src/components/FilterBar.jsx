import React from 'react';

export default function FilterBar({ sortBy, sortOrder, onSortChange }) {
  const handleSortByChange = (e) => {
    onSortChange(e.target.value, sortOrder);
  };

  const toggleSortOrder = () => {
    const newOrder = sortOrder === 'asc' ? 'desc' : 'asc';
    onSortChange(sortBy, newOrder);
  };

  return (
    <div className="filter-bar">
      <label>Sort by:</label>
      <select value={sortBy} onChange={handleSortByChange}>
        <option value="date">Date</option>
        <option value="size">Size</option>
      </select>
      <button onClick={toggleSortOrder}>
        {sortOrder === 'asc' ? '↑ Ascending' : '↓ Descending'}
      </button>
    </div>
  );
}
