import React, { useState } from 'react';
import { TextField, Button } from '@mui/material';
import { ToastContainer, toast } from 'react-toastify';
import axios from 'axios';

const TokenForm = ({ onSubmit }) => {
  const [name, setName] = useState('');
  const [limit, setLimit] = useState('');

  const handleNameChange = (event) => {
    setName(event.target.value);
  };

  const handleLimitChange = (event) => {
    setLimit(event.target.value);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!Number.isNaN(parseFloat(limit))) {
      const tokenData = { name, limit };

      try {
        const formData = new URLSearchParams();
        formData.append('name', name);
        formData.append('limit', limit);
        axios.defaults.withCredentials = true;

        const token = localStorage.getItem('token');

        const response = await axios.post('http://signaliha.com/admin/tokens', formData, {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            "Authorization":token,
          },
        });

        if (response.status == 200) {
          const newToken = await response.data;
          onSubmit(newToken);
          setName('');
          setLimit('');
        } else {
          toast.error('Failed to create token:', response.status);
        }
      } catch (error) {
        toast.error('An error occurred while creating token:', error);
      }
    } else {
      toast.error('Invalid limit input');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <ToastContainer />
      <TextField
        label="Name"
        value={name}
        onChange={handleNameChange}
        required
        fullWidth
        margin="normal"
      />
      <TextField
        label="Limit"
        value={limit}
        onChange={handleLimitChange}
        required
        type="number"
        fullWidth
        margin="normal"
      />
      <Button type="submit" variant="contained">
        Create Token
      </Button>
    </form>
  );
};

export default TokenForm;
