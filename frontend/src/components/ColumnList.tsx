// src/components/ColumnList.tsx
import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from '../store';
import { addColumn, removeColumn } from '../features/columns/slice';
import {
  Button,
  TextField,
  List,
  ListItem,
  ListItemText,
  IconButton,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';

const ColumnList: React.FC = () => {
  const columns = useSelector((state: RootState) => state.columns.columns);
  const dispatch = useDispatch();
  const [newColumnName, setNewColumnName] = useState('');

  const handleAddColumn = () => {
    if (newColumnName.trim()) {
      dispatch(addColumn(newColumnName));
      setNewColumnName('');
    }
  };

  const handleRemoveColumn = (id: string) => {
    dispatch(removeColumn(id));
  };

  return (
    <div style={{ padding: '16px' }}>
      <h2>Columns</h2>
      <TextField
        value={newColumnName}
        onChange={(e) => setNewColumnName(e.target.value)}
        placeholder="New Column Name"
        size="small"
      />
      <Button variant="contained" onClick={handleAddColumn} style={{ marginLeft: '8px' }}>
        Add Column
      </Button>
      <List>
        {columns.map((column) => (
          <ListItem key={column.id} divider>
            <ListItemText primary={column.name} />
            <IconButton edge="end" onClick={() => handleRemoveColumn(column.id)}>
              <DeleteIcon />
            </IconButton>
          </ListItem>
        ))}
      </List>
    </div>
  );
};

export default ColumnList;
