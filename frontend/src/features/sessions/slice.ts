import { type PayloadAction, createSlice } from "@reduxjs/toolkit";

interface sessionState {
  isLoggedIn: boolean;
  token: string | null;
}

// リロードしてもログインを保つため sessionStorage に置く。使えない環境（プライベートモードなど）では毎回ログインする
export const tokenKey = "accessToken";
const saved = (() => {
  try {
    return sessionStorage.getItem(tokenKey);
  } catch {
    return null;
  }
})();

const initialState: sessionState = {
  isLoggedIn: saved !== null,
  token: saved,
};

const sessionSlice = createSlice({
  name: "session",
  initialState,
  reducers: {
    login(state, action: PayloadAction<{ token: string }>) {
      state.isLoggedIn = true;
      state.token = action.payload.token;
    },
    logout(state) {
      state.isLoggedIn = false;
      state.token = null;
    },
  },
});

export const { login, logout } = sessionSlice.actions;
export default sessionSlice.reducer;
