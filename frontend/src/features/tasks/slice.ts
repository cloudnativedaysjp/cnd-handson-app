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
        title: "サンプルタスク1",
        description: "これはサンプルタスクの説明です。",
        status: "To Do",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-1",
      },
      {
        id: generateId(),
        title: "サンプルタスク1",
        description: "これはサンプルタスクの説明です。",
        status: "To Do",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-1",
      },
      {
        id: generateId(),
        title: "サンプルタスク1",
        description: "これはサンプルタスクの説明です。",
        status: "To Do",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-1",
      },
      {
        id: generateId(),
        title: "サンプルタスク1",
        description: "これはサンプルタスクの説明です。",
        status: "To Do",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-1",
      },
      {
        id: generateId(),
        title: "サンプルタスク1",
        description: "これはサンプルタスクの説明です。",
        status: "To Do",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-1",
      },
      {
        id: generateId(),
        title: "サンプルタスク2",
        description: "別のサンプルタスクです。",
        status: "In Progress",
        startTimeAt: undefined,
        endTimeAt: undefined,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        assigneeId: "user-2",
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
            assigneeId,
          },
        };
      },
    },

  },
});

export const { addTask } = tasksSlice.actions;
export default tasksSlice.reducer;
