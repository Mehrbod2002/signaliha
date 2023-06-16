import React, { useState } from 'react';
import { TextField, Button, Box, Typography } from '@mui/material';
import { toast,ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import axios from 'axios';

const AdminLogin = ({ onLogin }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleUsernameChange = (e) => {
    setUsername(e.target.value);
  };
  const handlePasswordChange = (e) => {
    setPassword(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (username.trim() === '' || password.trim() === '') {
      toast.error('Username and password cannot be empty');
      return;
    }

    try {
        const formData = new URLSearchParams();
        formData.append('username', username);
        formData.append('password', password);
        axios.defaults.withCredentials = true;
        const response = await axios.post('http://signaliha.com:8080/login', formData, {
            headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            },
        });

        if (response.status === 200) {
            const { token } = response.data;
            onLogin(token);
        } else {
          const errorMessage = await response.text();
          toast.error(JSON.parse(errorMessage).message);
        }
      } catch (error) {
        console.error('An error occurred during login:', error);
        toast.error('An error occurred during login');
      }      
  };

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight="100vh"
    >
    <ToastContainer/>
      <form onSubmit={handleSubmit}>
        <Box maxWidth={400} width="100%" padding={2}>
          <Typography variant="h4" align="center" gutterBottom>
            Admin Panel
          </Typography>
          <TextField
            label="Username"
            value={username}
            onChange={handleUsernameChange}
            fullWidth
            margin="normal"
          />
          <TextField
            label="Password"
            type="password"
            value={password}
            onChange={handlePasswordChange}
            fullWidth
            margin="normal"
          />
          <Button type="submit" variant="contained" color="primary" fullWidth>
            Login
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default AdminLogin;
