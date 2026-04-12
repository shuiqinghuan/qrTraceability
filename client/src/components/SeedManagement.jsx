import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const SeedManagement = () => {
  const [seeds, setSeeds] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showAddForm, setShowAddForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    variety: '',
    description: '',
  });

  useEffect(() => {
    fetchSeeds();
  }, []);

  const fetchSeeds = async () => {
    setLoading(true);
    try {
      const response = await api.request('/seed');
      if (response.seeds) {
        setSeeds(response.seeds);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to fetch seeds');
    } finally {
      setLoading(false);
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
      const response = await api.request('/seed', {
        method: 'POST',
        body: JSON.stringify(formData),
      });
      if (response.seed) {
        setSeeds(prev => [response.seed, ...prev]);
        setFormData({ name: '', variety: '', description: '' });
        setShowAddForm(false);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to create seed');
    }
  };

  return (
    <div className="seed-management">
      <h3>Seed Management</h3>
      {error && <div className="error-message">{error}</div>}
      
      <button 
        onClick={() => setShowAddForm(!showAddForm)}
        className="add-button"
      >
        {showAddForm ? 'Cancel' : 'Add Seed'}
      </button>

      {showAddForm && (
        <form onSubmit={handleSubmit} className="seed-form">
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="variety">Variety</label>
            <input
              type="text"
              id="variety"
              name="variety"
              value={formData.variety}
              onChange={handleChange}
              required
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
            Create Seed
          </button>
        </form>
      )}

      <div className="seed-list">
        <h4>Seed List</h4>
        {loading ? (
          <p>Loading...</p>
        ) : seeds.length === 0 ? (
          <p>No seeds found</p>
        ) : (
          <ul>
            {seeds.map(seed => (
              <li key={seed.id} className="seed-item">
                <div>
                  <strong>{seed.name}</strong> ({seed.variety})
                  {seed.description && <p>{seed.description}</p>}
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default SeedManagement;