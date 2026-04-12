import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const ProductBrowsing = () => {
  const [products, setProducts] = useState([]);
  const [favorites, setFavorites] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchProducts();
    fetchFavorites();
  }, []);

  const fetchProducts = async () => {
    setLoading(true);
    try {
      const response = await api.request('/products');
      if (response.plantings) {
        setProducts(response.plantings);
      } else if (response.error) {
        setError(response.error);
      }
    } catch (err) {
      setError('Failed to fetch products');
    } finally {
      setLoading(false);
    }
  };

  const fetchFavorites = async () => {
    try {
      const response = await api.request('/user/favorites');
      if (response.favorites) {
        setFavorites(response.favorites);
      }
    } catch (err) {
      console.error('Failed to fetch favorites', err);
    }
  };

  const toggleFavorite = async (plantingId) => {
    try {
      const isFavorited = favorites.some(fav => fav.id === plantingId);
      if (isFavorited) {
        await api.request(`/user/favorite/${plantingId}`, {
          method: 'DELETE',
        });
        setFavorites(prev => prev.filter(fav => fav.id !== plantingId));
      } else {
        await api.request(`/user/favorite/${plantingId}`, {
          method: 'POST',
        });
        const product = products.find(p => p.id === plantingId);
        if (product) {
          setFavorites(prev => [product, ...prev]);
        }
      }
    } catch (err) {
      setError('Failed to toggle favorite');
    }
  };

  const toggleLike = async (plantingId) => {
    try {
      await api.request(`/user/like/${plantingId}`, {
        method: 'POST',
      });
      // 这里可以添加点赞状态的更新逻辑
    } catch (err) {
      setError('Failed to toggle like');
    }
  };

  const shareProduct = (plantingId) => {
    // 实现分享功能
    alert(`Sharing product ${plantingId}`);
  };

  const isFavorited = (plantingId) => {
    return favorites.some(fav => fav.id === plantingId);
  };

  return (
    <div className="product-browsing">
      <h3>Product Browsing</h3>
      {error && <div className="error-message">{error}</div>}

      <div className="product-list">
        {loading ? (
          <p>Loading...</p>
        ) : products.length === 0 ? (
          <p>No products found</p>
        ) : (
          <div className="product-grid">
            {products.map(product => (
              <div key={product.id} className="product-card">
                <div className="product-header">
                  <h4>Product #{product.id}</h4>
                  <div className="product-actions">
                    <button
                      onClick={() => toggleFavorite(product.id)}
                      className={`action-button ${isFavorited(product.id) ? 'favorited' : ''}`}
                      title={isFavorited(product.id) ? 'Remove from favorites' : 'Add to favorites'}
                    >
                      {isFavorited(product.id) ? '★' : '☆'}
                    </button>
                    <button
                      onClick={() => toggleLike(product.id)}
                      className="action-button"
                      title="Like"
                    >
                      ❤
                    </button>
                    <button
                      onClick={() => shareProduct(product.id)}
                      className="action-button"
                      title="Share"
                    >
                      ↗
                    </button>
                  </div>
                </div>
                <div className="product-details">
                  <p><strong>Planting Date:</strong> {product.planting_date}</p>
                  {product.transplanting_date && (
                    <p><strong>Transplanting Date:</strong> {product.transplanting_date}</p>
                  )}
                  {product.location && (
                    <p><strong>Location:</strong> {product.location}</p>
                  )}
                  {product.notes && (
                    <p><strong>Notes:</strong> {product.notes}</p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="favorites-section">
        <h4>My Favorites</h4>
        {favorites.length === 0 ? (
          <p>No favorites yet</p>
        ) : (
          <div className="favorite-grid">
            {favorites.map(product => (
              <div key={product.id} className="favorite-card">
                <h5>Product #{product.id}</h5>
                <p>Planting Date: {product.planting_date}</p>
                <button
                  onClick={() => toggleFavorite(product.id)}
                  className="remove-favorite-button"
                >
                  Remove
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ProductBrowsing;