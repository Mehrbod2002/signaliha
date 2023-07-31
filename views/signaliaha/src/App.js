import React, { useState, useEffect } from 'react';
import AdminLogin from './AdminLogin';
import Home from './Home';
import axios from 'axios';
import { ToastContainer, toast } from 'react-toastify';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import DataList from './DataList';

const App = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const handleLogin = (token) => {
    localStorage.setItem('token', token);
    setIsAuthenticated(true);
  };

  const handleLogout = async () => {
      try {
        const token = localStorage.getItem('token');
        const formData = new URLSearchParams();
        formData.append('token', token);
        const response = await axios.post(`http://signaliha.com/logout`, formData, {
          headers: {
            "Authorization": localStorage.getItem('token'),
          },
          withCredentials: true,
        });

        if (response.status === 200) {
          toast.success('Logged out successfully');
          localStorage.removeItem('token');
          setIsAuthenticated(false);
        } else {
          toast.error('Failed to log out:', response.status);
        }
      } catch (error) {
        toast.error('An error occurred while logging out:', error);
      }
  };

  return (
    <div>
      <Router>
        <Routes>
          <Route path="/datas" element={<DataList />} />
          <Route path="/" element={isAuthenticated ? <Home onLogout={handleLogout} /> : <AdminLogin onLogin={handleLogin} />} />
        </Routes>
      </Router>
      <ToastContainer />
    </div>
  );  
};

export default App;
