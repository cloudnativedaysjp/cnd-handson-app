// src/store/index.ts
import { configureStore } from "@reduxjs/toolkit";
import tasksReducer from '../features/tasks/slice';
import columnReducer from '../features/columns/slice';

export const store = configureStore({
  reducer: {
    tasks: tasksReducer,
    columns: columnReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
