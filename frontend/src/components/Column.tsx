// src/components/Column.tsx
import React from 'react';
import { Box, Typography } from '@mui/material';
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
    border: isOver ? '2px dashed #333' : '1px solid #ccc',
    borderRadius: 2,
    padding: '16px',
    minWidth: '250px',
    backgroundColor: '#fafafa',
  };

  return (
    <Box
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
    </Box>
  );
};

export default Column;
