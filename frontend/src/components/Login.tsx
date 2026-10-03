import { Box, Button, Paper, TextField, Typography } from "@mui/material";
import type React from "react";
import { useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { login } from "../features/sessions/slice";
import { api } from "../utils/api";

const Login: React.FC = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [failed, setFailed] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim() || !password) return;

    try {
      const { accessToken } = await api<{ accessToken: string }>("/api/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      dispatch(login({ token: accessToken }));
      navigate("/");
    } catch {
      setFailed(true);
    }
  };

  return (
    <Box
      sx={{
        width: 400,
        margin: "80px auto",
      }}
    >
      <Paper sx={{ p: 4 }}>
        <Typography variant="h5" gutterBottom align="center">
          ログイン
        </Typography>
        <Box
          component="form"
          onSubmit={handleSubmit}
          sx={{ display: "flex", flexDirection: "column", gap: 2 }}
        >
          <TextField
            label="メールアドレス"
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            slotProps={{ htmlInput: { "data-testid": "login-email" } }}
          />
          <TextField
            label="パスワード"
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            slotProps={{ htmlInput: { "data-testid": "login-password" } }}
          />
          {failed && (
            <Typography color="error">
              メールアドレスかパスワードが違います
            </Typography>
          )}
          <Button
            type="submit"
            variant="contained"
            fullWidth
            data-testid="login-submit"
          >
            ログイン
          </Button>
        </Box>
      </Paper>
    </Box>
  );
};

export default Login;
