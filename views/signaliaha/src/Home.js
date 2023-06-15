import React, { useEffect, useState } from 'react';
import { TableContainer, Paper, Table, TableHead, TableRow, TableCell, TableBody, Typography, Button, Box, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { ToastContainer, toast } from 'react-toastify';
import TokenForm from './CreateToken';
import axios from 'axios';

const Home = ({ onLogout }) => {
    const [tokens, setTokens] = useState([]);
    const [isTokenFormOpen, setTokenFormOpen] = useState(false);
    const [refreshTrigger, setRefreshTrigger] = useState(false);
    const [selectedToken, setSelectedToken] = useState(null);
    const [tokenData, setTokenData] = useState([]);
  
    useEffect(() => {
      fetchTokens();
    }, [refreshTrigger]);
  
    const fetchTokens = async () => {
      try {
        const token = localStorage.getItem('token');
  
        const response = await axios.get('http://localhost:8080/admin/tokens', {
          headers: {
            Authorization: token,
          },
          withCredentials: true,
        });
  
        if (response.status === 200) {
          const data = await response.data;
          setTokens(data);
        } else {
          toast.error('Failed to fetch tokens:', response.status);
        }
      } catch (error) {
        toast.error('An error occurred while fetching tokens:', error);
      }
    };
  
    const handleOpenTokenForm = () => {
      setTokenFormOpen(true);
    };
  
    const handleCloseTokenForm = () => {
      setTokenFormOpen(false);
      setRefreshTrigger(!refreshTrigger);
    };
  
    const handleTokenFormSubmit = (newToken) => {
      setTokens([...tokens, newToken]);
      setTokenFormOpen(false);
      setRefreshTrigger(!refreshTrigger);
    };
  
    const handleDeleteToken = async (token) => {
      try {
        const formData = new URLSearchParams();
        formData.append('token', token);
        const response = await axios.post(`http://localhost:8080/admin/tokens/delete`, formData, {
          headers: {
            "Authorization": localStorage.getItem('token'),
          },
          withCredentials: true,
        });
  
        if (response.status === 200) {
          toast.success('Token deleted successfully');
          setRefreshTrigger(!refreshTrigger);
        } else {
          toast.error('Failed to delete token:', response.status);
        }
      } catch (error) {
        toast.error('An error occurred while deleting token:', error);
      }
    };
    
    const handleTokenClick = async (token) => {
        try {
            const formData = new URLSearchParams();
            formData.append('token', token);
            const response = await axios.get(`http://localhost:8080/admin/tokens/${token}/history`, formData, {
                headers: {
                    "Authorization": localStorage.getItem('token'),
                },
                withCredentials: true,
            });
  
            if (response.status === 200) {
                const data = await response.data;
                setTokenData(data);
                setSelectedToken(token);
            } else {
                toast.error('Failed to fetch token data:', response.status);
            }
        } catch (error) {
            toast.error('An error occurred while fetching token data:', error);
        }
      };
    
      const handleCloseTokenDialog = () => {
        setSelectedToken(null);
        setTokenData([]);
      };
    return (
      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', background: '#f5f5f5' }}>
        <Typography variant="h4" gutterBottom>
          Tokens
        </Typography>
        <ToastContainer />
        <Box component={TableContainer} sx={{ width: '80%', background: 'white', borderRadius: '8px', padding: '1rem' }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell align="center">Token</TableCell>
                <TableCell align="center">Name</TableCell>
                <TableCell align="center">Limit</TableCell>
                <TableCell align="center">Timestamp</TableCell>
                <TableCell align="center">Actions</TableCell> {/* Added Actions column */}
              </TableRow>
            </TableHead>
            <TableBody>
              {tokens.map((token) => (
                <TableRow key={token.token}>
                  <TableCell onClick={() => handleTokenClick(token.token)} align="center">{token.token}</TableCell>
                  <TableCell align="center">{token.name}</TableCell>
                  <TableCell align="center">{token.limit}</TableCell>
                  <TableCell align="center">{token.timestamp}</TableCell>
                  <TableCell align="center">
                    <Button variant="contained" color="secondary" onClick={() => handleDeleteToken(token.token)}>
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>
        <Button variant="contained" onClick={handleOpenTokenForm} style={{ marginTop: '1rem' }}>
          Create Token
        </Button>
        <Button variant="contained" onClick={onLogout} style={{ marginTop: '1rem' }}>
          Logout
        </Button>
  
            <Dialog open={Boolean(selectedToken)} onClose={handleCloseTokenDialog}>
            <DialogTitle>Token Details - {selectedToken}</DialogTitle>
            <DialogContent>
            <Table>
                <TableHead>
                <TableRow>
                    <TableCell align="center">ID</TableCell>
                    <TableCell align="center">Token</TableCell>
                    <TableCell align="center">Time</TableCell>
                    <TableCell align="center">Result</TableCell>
                    <TableCell align="center">Request</TableCell>
                </TableRow>
                </TableHead>
                <TableBody>
                {Array.isArray(tokenData) && tokenData.length > 0 ? (
                    tokenData.map((data) => (
                    <TableRow key={data.id}>
                        <TableCell align="center">{data.id}</TableCell>
                        <TableCell align="center">{data.token}</TableCell>
                        <TableCell align="center">{data.time}</TableCell>
                        <TableCell align="center">{data.result}</TableCell>
                        <TableCell align="center">{data.request}</TableCell>
                    </TableRow>
                    ))
                ) : (
                    <TableRow>
                    <TableCell colSpan={5} align="center">
                        No token data available
                    </TableCell>
                    </TableRow>
                )}
                </TableBody>
            </Table>
            </DialogContent>
            <DialogActions>
            <Button onClick={handleCloseTokenDialog}>Close</Button>
            </DialogActions>
        </Dialog>

        <Dialog open={isTokenFormOpen} onClose={handleCloseTokenForm}>
          <DialogTitle>Create Token</DialogTitle>
          <DialogContent>
            <TokenForm onSubmit={handleTokenFormSubmit} onClose={handleCloseTokenForm} />
          </DialogContent>
        </Dialog>
      </div>
    );
  };
  
  export default Home;
  