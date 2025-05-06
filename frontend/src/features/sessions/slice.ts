
import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface AuthState {
  isLoggedIn: boolean;
  userEmail: string | null;
}

const initialState: AuthState = {
  isLoggedIn: false,
  userEmail: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    login(state, action: PayloadAction<{ email: string }>) {
      state.isLoggedIn = true;
      state.userEmail = action.payload.email;
    },
    logout(state) {
      state.isLoggedIn = false;
      state.userEmail = null;
    },
  },
});

export const { login, logout } = authSlice.actions;
export default authSlice.reducer;
