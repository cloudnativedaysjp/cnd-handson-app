// src/components/Column.tsx
import React from 'react';
import { Box, Typography } from '@mui/material';
import { Column as ColumnType } from '../features/columns/types';
import { Task } from '../features/tasks/types';
import TaskCard from './TaskCrad';

interface ColumnProps {
  column: ColumnType;
  tasks: Task[];
}

const Column: React.FC<ColumnProps> = ({ column, tasks }) => {
  return (
    <Box
      sx={{
        border: '1px solid #ccc',
        borderRadius: 2,
        p: 2,
        flex: 1,
        minWidth: 250,
        backgroundColor: '#fafafa',
      }}
    >
      <Typography variant="h5" gutterBottom>
        {column.name}
      </Typography>
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} />
      ))}
    </Box>
  );
};

export default Column;
