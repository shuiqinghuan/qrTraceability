import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const GrowthMediaManagement = () => {
  const [plantings, setPlantings] = useState([]);
  const [media, setMedia] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState({
    planting_id: '',
    media_type: 'image',
    file_path: '',
    description: '',
  });

  useEffect(() => {
    fetchPlantings();
  }, []);

  const fetchPlantings = async () => {
    setLoading(true);
    try {
      const response = await api.request('/planting');
      if (response.plantings) {
        setPlantings(response.plantings);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to fetch plantings');
    } finally {
      setLoading(false);
    }
  };

  const fetchMedia = async (plantingId) => {
    try {
      const response = await api.request(`/growth/planting/${plantingId}`);
      if (response.media) {
        setMedia(response.media);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to fetch media');
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      const response = await api.request('/growth', {
        method: 'POST',
        body: JSON.stringify(formData),
      });
      if (response.media) {
        setMedia(prev => [response.media, ...prev]);
        setFormData({ planting_id: '', media_type: 'image', file_path: '', description: '' });
        setShowAddForm(false);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to create media');
    }
  };

  return (
    <div className="growth-media-management">
      <h3>Growth Media Management</h3>
      {error && <div className="error-message">{error}</div>}
      
      <button 
        onClick={() => setShowAddForm(!showAddForm)}
        className="add-button"
      >
        {showAddForm ? 'Cancel' : 'Add Media'}
      </button>

      {showAddForm && (
        <form onSubmit={handleSubmit} className="media-form">
          <div className="form-group">
            <label htmlFor="planting_id">Planting</label>
            <select
              id="planting_id"
              name="planting_id"
              value={formData.planting_id}
              onChange={handleChange}
              required
            >
              <option value="">Select a planting</option>
              {plantings.map(planting => (
                <option key={planting.id} value={planting.id}>
                  Planting #{planting.id}
                </option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label htmlFor="media_type">Media Type</label>
            <select
              id="media_type"
              name="media_type"
              value={formData.media_type}
              onChange={handleChange}
              required
            >
              <option value="image">Image</option>
              <option value="video">Video</option>
            </select>
          </div>
          <div className="form-group">
            <label htmlFor="file_path">File Path</label>
            <input
              type="text"
              id="file_path"
              name="file_path"
              value={formData.file_path}
              onChange={handleChange}
              required
              placeholder="Enter file path or URL"
            />
          </div>
          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={3}
            />
          </div>
          <button type="submit" className="submit-button">
            Create Media
          </button>
        </form>
      )}

      <div className="media-list">
        <h4>Media List</h4>
        <div className="planting-selector">
          <label htmlFor="select-planting">Select Planting:</label>
          <select
            id="select-planting"
            onChange={(e) => {
              if (e.target.value) {
                fetchMedia(e.target.value);
              } else {
                setMedia([]);
              }
            }}
          >
            <option value="">All Plantings</option>
            {plantings.map(planting => (
              <option key={planting.id} value={planting.id}>
                Planting #{planting.id}
              </option>
            ))}
          </select>
        </div>
        {loading ? (
          <p>Loading...</p>
        ) : media.length === 0 ? (
          <p>No media found</p>
        ) : (
          <ul>
            {media.map(item => (
              <li key={item.id} className="media-item">
                <div>
                  <strong>{item.media_type === 'image' ? 'Image' : 'Video'}</strong>
                  <p>File: {item.file_path}</p>
                  {item.description && <p>Description: {item.description}</p>}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default GrowthMediaManagement;