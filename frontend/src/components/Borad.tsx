import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import TaskCard from './TaskCrad';

const Board: React.FC = () => {

  const tasks = useSelector((state: RootState) => state.tasks.tasks);

  const todoTasks = tasks.filter((task: { status: string; }) => task.status === 'todo');
  const inProgressTasks = tasks.filter((task: { status: string; }) => task.status === 'in-progress');
  const doneTasks = tasks.filter((task: { status: string; }) => task.status === 'done');

  return (
    <div style={{ display: 'flex', gap: '16px' }}>
      <div style={{ flex: 1 }}>
        <h2>Todo</h2>
        {todoTasks.map((task: { id: any; title: string; description?: string | undefined; status: string; startTimeAt?: string | undefined; endTimeAt?: string | undefined; createdAt: string; updatedAt: string; assigneeId: string; }) => (
          <TaskCard key={task.id} task={task} />
        ))}
      </div>
      <div style={{ flex: 1 }}>
        <h2>In Progress</h2>
        {inProgressTasks.map((task: { id: any; title: string; description?: string | undefined; status: string; startTimeAt?: string | undefined; endTimeAt?: string | undefined; createdAt: string; updatedAt: string; assigneeId: string; }) => (
          <TaskCard key={task.id} task={task} />
        ))}
      </div>
      <div style={{ flex: 1 }}>
        <h2>Done</h2>
        {doneTasks.map((task: { id: any; title: string; description?: string | undefined; status: string; startTimeAt?: string | undefined; endTimeAt?: string | undefined; createdAt: string; updatedAt: string; assigneeId: string; }) => (
          <TaskCard key={task.id} task={task} />
        ))}
      </div>
    </div>
  );
};

export default Board;
