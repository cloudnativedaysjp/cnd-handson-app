import { configureStore } from "@reduxjs/toolkit";
import boardsReducer from "../features/boards/slice";
import columnReducer from "../features/columns/slice";
import projectReducer from "../features/projects/slice";
import sessionReducer, { tokenKey } from "../features/sessions/slice";
import tasksReducer from "../features/tasks/slice";

export const store = configureStore({
  reducer: {
    tasks: tasksReducer,
    columns: columnReducer,
    boards: boardsReducer,
    projects: projectReducer,
    session: sessionReducer,
  },
});

store.subscribe(() => {
  const { token } = store.getState().session;
  try {
    if (token) sessionStorage.setItem(tokenKey, token);
    else sessionStorage.removeItem(tokenKey);
  } catch {
    // 保存できなくても、このタブの中ではログインを保てる
  }
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
