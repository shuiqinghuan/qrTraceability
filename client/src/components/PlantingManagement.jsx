import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const PlantingManagement = () => {
  const [plantings, setPlantings] = useState([]);
  const [seeds, setSeeds] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState({
    seed_id: '',
    planting_date: '',
    transplanting_date: '',
    location: '',
    notes: '',
  });

  useEffect(() => {
    fetchPlantings();
    fetchSeeds();
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

  const fetchSeeds = async () => {
    try {
      const response = await api.request('/seed');
      if (response.seeds) {
        setSeeds(response.seeds);
      }
    } catch (err) {
      console.error('Failed to fetch seeds', err);
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
      const response = await api.request('/planting', {
        method: 'POST',
        body: JSON.stringify(formData),
      });
      if (response.planting) {
        setPlantings(prev => [response.planting, ...prev]);
        setFormData({ seed_id: '', planting_date: '', transplanting_date: '', location: '', notes: '' });
        setShowAddForm(false);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to create planting');
    }
  };

  return (
    <div className="planting-management">
      <h3>Planting Management</h3>
      {error && <div className="error-message">{error}</div>}
      
      <button 
        onClick={() => setShowAddForm(!showAddForm)}
        className="add-button"
      >
        {showAddForm ? 'Cancel' : 'Add Planting'}
      </button>

      {showAddForm && (
        <form onSubmit={handleSubmit} className="planting-form">
          <div className="form-group">
            <label htmlFor="seed_id">Seed</label>
            <select
              id="seed_id"
              name="seed_id"
              value={formData.seed_id}
              onChange={handleChange}
              required
            >
              <option value="">Select a seed</option>
              {seeds.map(seed => (
                <option key={seed.id} value={seed.id}>
                  {seed.name} ({seed.variety})
                </option>
              ))}
            </select>
          </div>
          <div className="form-group">
            <label htmlFor="planting_date">Planting Date</label>
            <input
              type="date"
              id="planting_date"
              name="planting_date"
              value={formData.planting_date}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="transplanting_date">Transplanting Date</label>
            <input
              type="date"
              id="transplanting_date"
              name="transplanting_date"
              value={formData.transplanting_date}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="location">Location</label>
            <input
              type="text"
              id="location"
              name="location"
              value={formData.location}
              onChange={handleChange}
            />
          </div>
          <div className="form-group">
            <label htmlFor="notes">Notes</label>
            <textarea
              id="notes"
              name="notes"
              value={formData.notes}
              onChange={handleChange}
              rows={3}
            />
          </div>
          <button type="submit" className="submit-button">
            Create Planting
          </button>
        </form>
      )}

      <div className="planting-list">
        <h4>Planting List</h4>
        {loading ? (
          <p>Loading...</p>
        ) : plantings.length === 0 ? (
          <p>No plantings found</p>
        ) : (
          <ul>
            {plantings.map(planting => {
              const seed = seeds.find(s => s.id === planting.seed_id);
              return (
                <li key={planting.id} className="planting-item">
                  <div>
                    <strong>Planting #{planting.id}</strong>
                    {seed && <p>Seed: {seed.name} ({seed.variety})</p>}
                    <p>Planting Date: {planting.planting_date}</p>
                    {planting.transplanting_date && (
                      <p>Transplanting Date: {planting.transplanting_date}</p>
                    )}
                    {planting.location && <p>Location: {planting.location}</p>}
                    {planting.notes && <p>Notes: {planting.notes}</p>}
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
};

export default PlantingManagement;