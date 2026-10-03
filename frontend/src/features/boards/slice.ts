import { type PayloadAction, createSlice } from "@reduxjs/toolkit";
import type { Board } from "./types";
function generateId(): string {
  return Math.random().toString(36).substr(2, 9);
}

type BoardsState = {
  boards: Board[];
};

const initialState: BoardsState = {
  boards: [
    {
      id: generateId(),
      name: "サンプルボード1",
      description: "説明1",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
    {
      id: generateId(),
      name: "サンプルボード2",
      description: "説明2",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    },
  ],
};

const boardsSlice = createSlice({
  name: "boards",
  initialState,
  reducers: {
    addBoard: (state, action: PayloadAction<Board>) => {
      state.boards.push(action.payload);
    },
    updateBoard: (state, action: PayloadAction<Board>) => {
      const index = state.boards.findIndex(
        (board) => board.id === action.payload.id,
      );
      if (index !== -1) {
        state.boards[index] = { ...state.boards[index], ...action.payload };
      }
    },
    deleteBoard: (state, action: PayloadAction<string>) => {
      state.boards = state.boards.filter(
        (board) => board.id !== action.payload,
      );
    },
  },
});
export const { addBoard, updateBoard, deleteBoard } = boardsSlice.actions;
export default boardsSlice.reducer;
