import React, { useState, useEffect } from 'react';
import {
  Table,
  TableContainer,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  Typography,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { Link } from 'react-router-dom';
import { toast } from 'react-toastify';
import axios from 'axios';

const DataList = () => {
  const [data, setData] = useState([]);
  const [selectedItem, setSelectedItem] = useState(null);

  const handleItemClick = (item) => {
    setSelectedItem(item);
  };

  const handleCloseDialog = () => {
    setSelectedItem(null);
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp * 1000);
    date.setHours(date.getHours() - 3);
    date.setMinutes(date.getMinutes() - 30);

    const options = {
      year: 'numeric',
      month: 'numeric',
      day: 'numeric',
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
      timeZone: 'Asia/Tehran',
    };

    const formattedDate = date.toLocaleDateString("fa-IR", options);
    return formattedDate;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://signaliha.com/admin/coins', {
          headers: {
            Authorization: localStorage.getItem('token'),
          },
          withCredentials: true,
        });

        if (response.status === 200) {
          const { data } = await response.data;
          setData([...data].reverse());
        } else {
          toast.error('Failed to fetch token data:', response.status);
        }
      } catch (error) {
        toast.error('An error occurred while fetching token data:', error);
      }
    };

    fetchData();
  }, []);

  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        background: '#f5f5f5',
        padding: '2rem',
      }}
    >
      <Typography variant="h4" gutterBottom>
        Data List
      </Typography>

      <Button variant="outlined" component={Link} to="/">
        Return Home
      </Button>

      <TableContainer
        style={{
          marginTop: '2rem',
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
          background: 'white',
          borderRadius: '8px',
        }}
      >
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Coin</TableCell>
              <TableCell>Base Currency</TableCell>
              <TableCell>Platform</TableCell>
              <TableCell>Timestamp</TableCell>
              <TableCell>Entries</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.map((item) => (
              <TableRow
                key={item.MessageID}
                onClick={() => handleItemClick(item)}
                style={{ cursor: 'pointer' }}
              >
                <TableCell>{item.MessageID}</TableCell>
                <TableCell>{item.Coin}</TableCell>
                <TableCell>{item.BaseCurrency}</TableCell>
                <TableCell>{item.Platform}</TableCell>
                <TableCell>{formatDate(item.Timestamp)}</TableCell>
                <TableCell>{item.Entries}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog open={Boolean(selectedItem)} onClose={handleCloseDialog}>
        <DialogTitle>Data Details - ID: {selectedItem?.MessageID}</DialogTitle>
        <DialogContent style={{ width: '600px' }}>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Message ID</TableCell>
                  <TableCell>Coin</TableCell>
                  <TableCell>Base Currency</TableCell>
                  <TableCell>Platform</TableCell>
                  <TableCell>Leverage</TableCell>
                  <TableCell>Side</TableCell>
                  <TableCell>Entries</TableCell>
                  <TableCell>Margin</TableCell>
                  <TableCell>TP</TableCell>
                  <TableCell>SL</TableCell>
                  <TableCell>Timestamp</TableCell>
                  <TableCell>Exit</TableCell>
                  <TableCell>Risk</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {selectedItem && (
                  <TableRow>
                    <TableCell>{selectedItem.MessageID}</TableCell>
                    <TableCell>{selectedItem.Coin}</TableCell>
                    <TableCell>{selectedItem.BaseCurrency}</TableCell>
                    <TableCell>{selectedItem.Platform}</TableCell>
                    <TableCell>{selectedItem.Leverage}</TableCell>
                    <TableCell>{selectedItem.Side}</TableCell>
                    <TableCell>{selectedItem.Entries}</TableCell>
                    <TableCell>{selectedItem.Margin}</TableCell>
                    <TableCell>{selectedItem.Tp}</TableCell>
                    <TableCell>{selectedItem.SL}</TableCell>
                    <TableCell>{formatDate(selectedItem.Timestamp)}</TableCell>
                    <TableCell>{String(selectedItem.Exit)}</TableCell>
                    <TableCell>{String(selectedItem.Risk)}</TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Close</Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default DataList;
