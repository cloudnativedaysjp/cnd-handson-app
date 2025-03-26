import React from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store';
import Column from './Column';
import { Column as ColumnType } from '../features/columns/types';
import { updateTaskColumn } from '../features/tasks/slice';
import { DndContext, DragEndEvent } from '@dnd-kit/core';


const KanbanBoard: React.FC = () => {
  const dispatch = useDispatch();
  const columns = useSelector((state: RootState) => state.columns.columns);
  const tasks = useSelector((state: RootState) => state.tasks.tasks);
  
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over) return;
    const taskId = active.id.toString();
    const newColumnId = over.id.toString();
    dispatch(updateTaskColumn({ taskId, columnId: newColumnId }));
  }

  return (
    <DndContext onDragEnd={handleDragEnd}>
      <div style={{ display: 'flex', gap: '16px', padding: '16px' }}>
        {columns.map((column: ColumnType) => {
          const tasksInColumn = tasks.filter(task => task.columnId === column.id);
          return (
            <Column key={column.id} column={column} tasks={tasksInColumn} />
          );
        })}
      </div>
    </DndContext>
  );
};

export default KanbanBoard;
