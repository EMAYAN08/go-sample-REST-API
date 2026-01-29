// GET /api/images
export async function fetchImages() {
  const response = await fetch('/api/images');
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to fetch images');
  }
  return await response.json();
}

// POST /api/images/upload
export async function uploadImage(file) {
  const formData = new FormData();
  formData.append('image', file);

  const response = await fetch('/api/images/upload', {
    method: 'POST',
    body: formData
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload image');
  }

  return await response.json();
}

// DELETE /api/images/{id}
export async function deleteImage(id) {
  const response = await fetch(`/api/images/${id}`, {
    method: 'DELETE'
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete image');
  }

  return await response.json();
}
