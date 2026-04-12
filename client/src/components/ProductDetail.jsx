import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const ProductDetail = ({ productId }) => {
  const [product, setProduct] = useState(null);
  const [seed, setSeed] = useState(null);
  const [media, setMedia] = useState([]);
  const [quality, setQuality] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [liked, setLiked] = useState(false);
  const [favorite, setFavorite] = useState(false);

  useEffect(() => {
    fetchProductDetails();
  }, [productId]);

  const fetchProductDetails = async () => {
    setLoading(true);
    try {
      // 获取产品基本信息
      const productResponse = await api.request(`/products/${productId}`);
      if (productResponse.planting) {
        setProduct(productResponse.planting);
        // 获取种子信息
        if (productResponse.planting.seed_id) {
          fetchSeedInfo(productResponse.planting.seed_id);
        }
        // 获取生长媒体
        fetchMedia(productResponse.planting.id);
        // 获取品质信息
        fetchQuality(productResponse.planting.id);
      } else if (productResponse.error) {
        setError(productResponse.error);
      }
    } catch (err) {
      setError('Failed to fetch product details');
    } finally {
      setLoading(false);
    }
  };

  const fetchSeedInfo = async (seedId) => {
    try {
      const response = await api.request(`/seed/${seedId}`);
      if (response.seed) {
        setSeed(response.seed);
      }
    } catch (err) {
      console.error('Failed to fetch seed info', err);
    }
  };

  const fetchMedia = async (plantingId) => {
    try {
      const response = await api.request(`/growth/planting/${plantingId}`);
      if (response.media) {
        setMedia(response.media);
      }
    } catch (err) {
      console.error('Failed to fetch media', err);
    }
  };

  const fetchQuality = async (plantingId) => {
    try {
      const response = await api.request(`/quality/planting/${plantingId}`);
      if (response.quality) {
        setQuality(response.quality);
      }
    } catch (err) {
      console.error('Failed to fetch quality info', err);
    }
  };

  const handleLike = async () => {
    try {
      // 实现防刷机制，这里简化处理
      if (localStorage.getItem(`liked_${productId}`)) {
        alert('You have already liked this product');
        return;
      }
      
      await api.request(`/user/like/${productId}`, {
        method: 'POST',
      });
      setLiked(true);
      localStorage.setItem(`liked_${productId}`, 'true');
    } catch (err) {
      setError('Failed to like product');
    }
  };

  const handleFavorite = async () => {
    try {
      await api.request(`/user/favorite/${productId}`, {
        method: 'POST',
      });
      setFavorite(true);
    } catch (err) {
      setError('Failed to favorite product');
    }
  };

  const handleShare = () => {
    // 实现分享功能
    if (navigator.share) {
      navigator.share({
        title: `Product #${productId}`,
        text: 'Check out this product',
        url: window.location.href,
      });
    } else {
      // 复制链接到剪贴板
      navigator.clipboard.writeText(window.location.href);
      alert('Link copied to clipboard');
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  if (error) {
    return <div className="error-message">{error}</div>;
  }

  if (!product) {
    return <div className="error-message">Product not found</div>;
  }

  return (
    <div className="product-detail">
      <h2>Product Details</h2>
      
      {/* 第一部分：产品基本信息 */}
      <section className="product-section">
        <h3>Product Information</h3>
        <div className="product-info">
          <p><strong>Variety:</strong> {seed ? seed.name : 'N/A'} ({seed ? seed.variety : 'N/A'})</p>
          <p><strong>Planting Location:</strong> {product.location || 'N/A'}</p>
          <p><strong>Planting Date:</strong> {product.planting_date}</p>
          {product.transplanting_date && (
            <p><strong>Transplanting Date:</strong> {product.transplanting_date}</p>
          )}
        </div>
      </section>

      {/* 第二部分：产品图片和视频 */}
      <section className="product-section">
        <h3>Media</h3>
        {media.length === 0 ? (
          <p>No media available</p>
        ) : (
          <div className="media-grid">
            {media.map(item => (
              <div key={item.id} className="media-item">
                {item.media_type === 'image' ? (
                  <img src={item.file_path} alt={`Media ${item.id}`} className="media-image" />
                ) : (
                  <video src={item.file_path} controls className="media-video" />
                )}
                {item.description && <p>{item.description}</p>}
              </div>
            ))}
          </div>
        )}
      </section>

      {/* 第三部分：采收质量信息 */}
      <section className="product-section">
        <h3>Quality Information</h3>
        {quality ? (
          <div className="quality-info">
            <p><strong>Sugar Content:</strong> {quality.sugar_content || 'N/A'}%</p>
            <p><strong>Weight:</strong> {quality.weight || 'N/A'}g</p>
            <p><strong>Taste Analysis:</strong> {quality.taste_analysis || 'N/A'}</p>
          </div>
        ) : (
          <p>No quality information available</p>
        )}
      </section>

      {/* 第四部分：用户交互反馈 */}
      <section className="product-section">
        <h3>Feedback</h3>
        <div className="feedback-actions">
          <button 
            onClick={handleLike} 
            className={`action-button ${liked ? 'liked' : ''}`}
            disabled={liked}
          >
            ❤ {liked ? 'Liked' : 'Like'}
          </button>
          <button 
            onClick={handleFavorite} 
            className={`action-button ${favorite ? 'favorited' : ''}`}
            disabled={favorite}
          >
            ★ {favorite ? 'Favorited' : 'Favorite'}
          </button>
          <button onClick={handleShare} className="action-button">
            ↗ Share
          </button>
        </div>
      </section>
    </div>
  );
};

export default ProductDetail;