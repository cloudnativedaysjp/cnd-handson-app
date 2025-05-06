
import React, { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
} from '@mui/material';
import { useDispatch } from 'react-redux';
import { login } from '../features/sessions/slice';
import { useNavigate } from 'react-router-dom';

const Login: React.FC = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim() || !password) return;

    dispatch(login({ email }));
    navigate('/');
  };

  return (
    <Box
      sx={{
        width: 400,
        margin: '80px auto',
      }}
    >
      <Paper sx={{ p: 4 }}>
        <Typography variant="h5" gutterBottom align="center">
          ログイン
        </Typography>
        <Box
          component="form"
          onSubmit={handleSubmit}
          sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}
        >
          <TextField
            label="メールアドレス"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <TextField
            label="パスワード"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button type="submit" variant="contained" fullWidth>
            ログイン
          </Button>
        </Box>
      </Paper>
    </Box>
  );
};

export default Login;
