// src/components/TaskModal.tsx
import React from 'react';
import { Dialog, DialogTitle, DialogContent, DialogContentText, DialogActions, Button } from '@mui/material';
import { Task } from '../features/tasks/types';

interface TaskModalProps {
  open: boolean;
  task: Task | null;
  onClose: () => void;
}

const TaskModal: React.FC<TaskModalProps> = ({ open, task, onClose }) => {
  if (!task) return null;

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{task.title}</DialogTitle>
      <DialogContent>
        <DialogContentText>{task.description || '説明なし'}</DialogContentText>
        <DialogContentText>Status: {task.status}</DialogContentText>
        <DialogContentText>作成日時: {task.createdAt}</DialogContentText>
        {/* 他の必要なフィールドを追加 */}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>閉じる</Button>
      </DialogActions>
    </Dialog>
  );
};

export default TaskModal;
