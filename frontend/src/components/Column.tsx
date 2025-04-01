// src/components/Column.tsx
import React from 'react';
import { Typography, Paper } from '@mui/material';
import { Column as ColumnType } from '../features/columns/types';
import { useDroppable } from '@dnd-kit/core';
import { Task } from '../features/tasks/types';
import TaskCard from './TaskCrad';

interface ColumnProps {
  column: ColumnType;
  tasks: Task[];
}

const Column: React.FC<ColumnProps> = ({ column, tasks }) => {
  const { isOver, setNodeRef } = useDroppable({
    id: column.id,
    data: { column },
  });

  const style = {
    border: isOver ? '2px dashed #4CAF50' : '1px solid #ccc',
    borderRadius: 16,
    padding: '24px',
    minWidth: '300px',
    backgroundColor: isOver ? '#f0fff0' : '#ffffff',
    boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
    transition: 'background-color 0.2s ease, box-shadow 0.2s ease',
  };

  return (
    <Paper
      ref={setNodeRef}
      sx={{
        style,
      }}
    >
      <Typography variant="h5" gutterBottom>
        {column.name}
      </Typography>
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} />
      ))}
    </Paper>
  );
};

export default Column;
