import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import { Board } from './types';
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
        name: 'サンプルボード1',
        description: '説明1',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
      {
        id: generateId(),
        name: 'サンプルボード2',
        description: '説明2',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      },
    ],
};

const boardsSlice = createSlice({
  name: 'boards',
  initialState,
  reducers: {
    addBoard: (state, action: PayloadAction<Board>) => {
      state.boards.push(action.payload);
    },
  },
});
export default boardsSlice.reducer;