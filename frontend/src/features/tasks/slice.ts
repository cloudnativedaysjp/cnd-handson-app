import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { v4 as uuidv4 } from 'uuid';
import { Task } from './types';
function generateId(): string {
    return Math.random().toString(36).substr(2, 9);
  }
  
type TasksState = {
  tasks: Task[];
};

const initialState: TasksState = {
    tasks: [
      {
        id: generateId(),
        title: 'サンプルタスク1',
        description: '説明1',
        status: 'open',
        columnId: 'column-1',
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: 'user-1',
      },
      {
        id: generateId(),
        title: 'サンプルタスク2',
        description: '説明2',
        status: 'in-progress',
        columnId: 'column-2',
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: 'user-2',
      },
    ],
  };

const tasksSlice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {
    addTask: {
      reducer(state, action: PayloadAction<Task>) {
        state.tasks.push(action.payload);
      },
      prepare(
        title: string,
        status: string,
        columnId: string,
        assigneeId: string,
        description: string = '',
        startTimeAt?: string,
        endTimeAt?: string
      ): { payload: Task } {
        const timestamp = new Date().toISOString();
        return {
          payload: {
            id: generateId(),
            title,      
            description,
            status,
            startTimeAt,
            endTimeAt,
            createdAt: timestamp,
            updatedAt: timestamp,
            columnId,
            assigneeId,
          },
        };
      },
    },
    updateTaskColumn(state, action: PayloadAction<{ taskId: string; columnId: string }>) {
      const { taskId, columnId } = action.payload;
      const task = state.tasks.find((task) => task.id === taskId);
      if (task) {
        task.columnId = columnId;
        task.updatedAt = new Date().toISOString();
      }
    },
  },
});

export const { addTask } = tasksSlice.actions;
export const { updateTaskColumn } = tasksSlice.actions;
export default tasksSlice.reducer;
