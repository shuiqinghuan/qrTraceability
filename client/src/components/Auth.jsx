import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

const Auth = ({ onAuthSuccess }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    phone_number: '',
    password: '',
    role: 'clentcustomer',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState('');
  const [passwordStrengthClass, setPasswordStrengthClass] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));

    // 密码强度检查
    if (name === 'password' && value) {
      checkPasswordStrength(value);
    }
  };

  const checkPasswordStrength = (password) => {
    let strength = 0;
    let strengthText = '';
    let strengthClass = '';

    if (password.length >= 8) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[a-z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[^A-Za-z0-9]/.test(password)) strength++;

    switch (strength) {
      case 0:
      case 1:
        strengthText = 'Weak';
        strengthClass = 'password-weak';
        break;
      case 2:
      case 3:
        strengthText = 'Medium';
        strengthClass = 'password-medium';
        break;
      case 4:
      case 5:
        strengthText = 'Strong';
        strengthClass = 'password-strong';
        break;
      default:
        strengthText = '';
        strengthClass = '';
    }

    setPasswordStrength(strengthText);
    setPasswordStrengthClass(strengthClass);
  };

  const validateForm = () => {
    if (!formData.phone_number) {
      setError('Please enter your phone number');
      return false;
    }

    if (!/^1[3-9]\d{9}$/.test(formData.phone_number)) {
      setError('Please enter a valid phone number');
      return false;
    }

    if (!formData.password) {
      setError('Please enter your password');
      return false;
    }

    if (formData.password.length < 6) {
      setError('Password must be at least 6 characters');
      return false;
    }

    if (!isLogin && !formData.role) {
      setError('Please select a role');
      return false;
    }

    return true;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!validateForm()) {
      return;
    }

    setLoading(true);

    try {
      let response;
      if (isLogin) {
        response = await api.auth.login(formData);
      } else {
        response = await api.auth.register(formData);
      }

      if (response.token) {
        localStorage.setItem('token', response.token);
        localStorage.setItem('user', JSON.stringify(response.user));
        onAuthSuccess(response.user);
      } else if (response.error) {
        setError(response.error);
      } else {
        setError('An unexpected error occurred');
      }
    } catch (err) {
      console.error('Auth error:', err);
      setError('Failed to connect to server. Please check your network.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <h2>{isLogin ? 'Login' : 'Register'}</h2>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit} className="auth-form">
        <div className="form-group">
          <label htmlFor="phone_number">Phone Number</label>
          <input
            type="text"
            id="phone_number"
            name="phone_number"
            value={formData.phone_number}
            onChange={handleChange}
            required
            placeholder="Enter your phone number"
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
            minLength={6}
            placeholder="Enter your password"
          />
          {!isLogin && formData.password && (
            <div className={`password-strength ${passwordStrengthClass}`}>
              Password strength: {passwordStrength}
            </div>
          )}
        </div>
        {!isLogin && (
          <div className="form-group">
            <label htmlFor="role">Role</label>
            <select
              id="role"
              name="role"
              value={formData.role}
              onChange={handleChange}
              required
            >
              <option value="serverseed">Server Seed (种子管理)</option>
              <option value="servergrow">Server Grow (生长管理)</option>
              <option value="servermanager">Server Manager (品质管理)</option>
              <option value="clentcustomer">Client Customer (客户)</option>
            </select>
          </div>
        )}
        <button type="submit" className="auth-button" disabled={loading}>
          {loading ? 'Processing...' : isLogin ? 'Login' : 'Register'}
        </button>
      </form>
      <p className="auth-toggle">
        {isLogin ? "Don't have an account?" : "Already have an account?"}
        <button
          onClick={() => setIsLogin(!isLogin)}
          className="toggle-button"
        >
          {isLogin ? 'Register' : 'Login'}
        </button>
      </p>
    </div>
  );
};

export default Auth;