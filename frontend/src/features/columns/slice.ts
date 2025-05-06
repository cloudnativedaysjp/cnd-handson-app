import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Column } from './types';
import { v4 as uuidv4 } from 'uuid';

interface ColumnsState {
    columns: Column[];
}

const initialState: ColumnsState = {
    columns: [
        {
            id: 'column-1',
            name: 'To Do',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
        },
        {
            id: 'column-2',
            name: 'In Progress',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
        },
        {
            id: 'column-3',
            name: 'Done',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
        },
    ],
};

const columnsSlice = createSlice({
    name: 'columns',
    initialState,
    reducers: {
        addColumn: {
            reducer(state, action: PayloadAction<Column>) {
                state.columns.push(action.payload);
            },
            prepare(name: string): { payload: Column } {
                return {
                    payload: {
                        id: uuidv4(),
                        name,
                        createdAt: new Date().toISOString(),
                        updatedAt: new Date().toISOString(),
                    },
                };
            },
        },
        removeColumn(state, action: PayloadAction<string>) {
            state.columns = state.columns.filter((column) => column.id !== action.payload);
        }
    },

});

export const { addColumn, removeColumn } = columnsSlice.actions;
export default columnsSlice.reducer;
