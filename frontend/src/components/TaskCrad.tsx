// src/components/TaskCard.tsx
import React, { useState } from 'react';
import { Card, CardContent, Typography } from '@mui/material';
import { Task } from '../features/tasks/types';
import { useDraggable } from '@dnd-kit/core';
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

  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: task.id,
    data: { task }
  });

  const style = {
    transform: transform ? `translate3d(${transform.x}px, ${transform.y}px, 0)` : undefined,
    transition: 'transform 200ms ease',
    margin: '8px 0',
    cursor: 'grab',
  };
    
  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
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
    </div>
  );
};

export default TaskCard;
