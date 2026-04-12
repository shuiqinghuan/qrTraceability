import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom';
import './App.css';
import Auth from './components/Auth';
import SeedManagement from './components/SeedManagement';
import PlantingManagement from './components/PlantingManagement';
import GrowthMediaManagement from './components/GrowthMediaManagement';
import ProductQualityManagement from './components/ProductQualityManagement';
import ProductBrowsing from './components/ProductBrowsing';
import ProductDetail from './components/ProductDetail';

function App() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // 检查本地存储中是否有用户信息
    const storedUser = localStorage.getItem('user');
    const storedToken = localStorage.getItem('token');
    if (storedUser && storedToken) {
      setUser(JSON.parse(storedUser));
    }
    setIsLoading(false);
  }, []);

  const handleAuthSuccess = (userData) => {
    setUser(userData);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <Router>
      <div className="app">
        <header className="app-header">
          <h1>Plantation Management System</h1>
          {user && (
            <div className="user-info">
              <span>Welcome, {user.phone_number}</span>
              <button onClick={handleLogout} className="logout-button">
                Logout
              </button>
            </div>
          )}
        </header>
        <main className="app-main">
          <Routes>
            {/* 产品详情页路由，不需要登录 */}
            <Route path="/product/:id" element={<ProductDetailWrapper />} />
            
            {/* 登录路由 */}
            <Route path="/login" element={!user ? <Auth onAuthSuccess={handleAuthSuccess} /> : <Navigate to="/" />} />
            
            {/* 主页路由 */}
            <Route path="/" element={user ? <Dashboard user={user} /> : <Navigate to="/login" />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

// 产品详情页包装组件
const ProductDetailWrapper = () => {
  const { id } = useParams();
  return <ProductDetail productId={parseInt(id)} />;
};

// 仪表板组件
const Dashboard = ({ user }) => {
  return (
    <div className="dashboard">
      <h2>Dashboard</h2>
      <p>Welcome to the dashboard, {user.role}!</p>
      <div className="role-content">
        {user.role === 'serverseed' && (
          <div>
            <SeedManagement />
            <PlantingManagement />
          </div>
        )}
        {user.role === 'servergrow' && (
          <div>
            <GrowthMediaManagement />
          </div>
        )}
        {user.role === 'servermanager' && (
          <div>
            <ProductQualityManagement />
          </div>
        )}
        {user.role === 'clentcustomer' && (
          <div>
            <ProductBrowsing />
          </div>
        )}
      </div>
    </div>
  );
};

export default App;