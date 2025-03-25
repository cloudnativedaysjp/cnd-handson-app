// src/components/TaskCard.tsx
import React, { useState } from 'react';
import { Card, CardContent, Typography } from '@mui/material';
import { Task } from '../features/tasks/types';
import TaskModal from './TaskModal';

interface TaskCardProps {
  task: Task;
}

const TaskCard: React.FC<TaskCardProps> = ({ task }) => {
  const [open, setOpen] = useState(false);

  const handleClick = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <>
      <Card style={{ margin: '8px', cursor: 'pointer' }} onClick={handleClick}>
        <CardContent>
          <Typography variant="h6">{task.title}</Typography>
          <Typography variant="body2" color="textSecondary">
            {task.description || '説明なし'}
          </Typography>
          <Typography variant="caption">Status: {task.status}</Typography>
        </CardContent>
      </Card>
      <TaskModal task={task} open={open} onClose={handleClose} />
    </>
  );
};

export default TaskCard;
