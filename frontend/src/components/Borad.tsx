import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import Column from './Column';
import { Column as ColumnType } from '../features/columns/types';

const KanbanBoard: React.FC = () => {
  const columns = useSelector((state: RootState) => state.columns.columns);
  const tasks = useSelector((state: RootState) => state.tasks.tasks);

  return (
    <div style={{ display: 'flex', gap: '16px', padding: '16px' }}>
      {columns.map((column: ColumnType) => {
        const tasksInColumn = tasks.filter(task => task.status === column.name);
        return (
          <Column key={column.id} column={column} tasks={tasksInColumn} />
        );
      })}
    </div>
  );
};

export default KanbanBoard;
