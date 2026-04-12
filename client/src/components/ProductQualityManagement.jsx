import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const ProductQualityManagement = () => {
  const [plantings, setPlantings] = useState([]);
  const [quality, setQuality] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState({
    planting_id: '',
    sugar_content: '',
    weight: '',
    taste_analysis: '',
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

  const fetchQuality = async (plantingId) => {
    try {
      const response = await api.request(`/quality/planting/${plantingId}`);
      if (response.quality) {
        setQuality(response.quality);
      } else if (response.error) {
        setError(response.error);
        setQuality(null);
      }
    } catch (err) {
      setError('Failed to fetch quality data');
      setQuality(null);
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
      const response = await api.request('/quality', {
        method: 'POST',
        body: JSON.stringify(formData),
      });
      if (response.quality) {
        setQuality(response.quality);
        setFormData({ planting_id: '', sugar_content: '', weight: '', taste_analysis: '' });
        setShowAddForm(false);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to create quality record');
    }
  };

  return (
    <div className="product-quality-management">
      <h3>Product Quality Management</h3>
      {error && <div className="error-message">{error}</div>}
      
      <div className="planting-selector">
        <label htmlFor="select-planting">Select Planting:</label>
        <select
          id="select-planting"
          onChange={(e) => {
            if (e.target.value) {
              fetchQuality(e.target.value);
            } else {
              setQuality(null);
            }
          }}
        >
          <option value="">Select a planting</option>
          {plantings.map(planting => (
            <option key={planting.id} value={planting.id}>
              Planting #{planting.id}
            </option>
          ))}
        </select>
      </div>

      {!quality && (
        <button 
          onClick={() => setShowAddForm(!showAddForm)}
          className="add-button"
        >
          {showAddForm ? 'Cancel' : 'Add Quality Record'}
        </button>
      )}

      {showAddForm && (
        <form onSubmit={handleSubmit} className="quality-form">
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
            <label htmlFor="sugar_content">Sugar Content</label>
            <input
              type="number"
              id="sugar_content"
              name="sugar_content"
              value={formData.sugar_content}
              onChange={handleChange}
              step="0.01"
              placeholder="e.g., 12.5"
            />
          </div>
          <div className="form-group">
            <label htmlFor="weight">Weight</label>
            <input
              type="number"
              id="weight"
              name="weight"
              value={formData.weight}
              onChange={handleChange}
              step="0.01"
              placeholder="e.g., 150.0"
            />
          </div>
          <div className="form-group">
            <label htmlFor="taste_analysis">Taste Analysis</label>
            <textarea
              id="taste_analysis"
              name="taste_analysis"
              value={formData.taste_analysis}
              onChange={handleChange}
              rows={3}
              placeholder="Enter taste analysis"
            />
          </div>
          <button type="submit" className="submit-button">
            Create Quality Record
          </button>
        </form>
      )}

      <div className="quality-details">
        <h4>Quality Details</h4>
        {loading ? (
          <p>Loading...</p>
        ) : quality ? (
          <div className="quality-card">
            <p><strong>Planting ID:</strong> {quality.planting_id}</p>
            {quality.sugar_content && (
              <p><strong>Sugar Content:</strong> {quality.sugar_content}%</p>
            )}
            {quality.weight && (
              <p><strong>Weight:</strong> {quality.weight}g</p>
            )}
            {quality.taste_analysis && (
              <p><strong>Taste Analysis:</strong> {quality.taste_analysis}</p>
            )}
          </div>
        ) : (
          <p>No quality data found for selected planting</p>
        )}
      </div>
    </div>
  );
};

export default ProductQualityManagement;