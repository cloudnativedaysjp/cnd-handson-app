import { useDroppable } from "@dnd-kit/core";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import { Box, IconButton, TextField, Typography } from "@mui/material";
import type React from "react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import type { Column as ColumnType } from "../features/columns/types";
import type { Task } from "../features/tasks/types";
import TaskCard from "./TaskCard";

interface ColumnProps {
  column: ColumnType;
  tasks: Task[];
  onUpdateColumnName?: (id: string, name: string) => void;
  onDeleteColumn?: (id: string) => void;
  onAdd?: () => void;
  initiallyEditing?: boolean;
}

const Column: React.FC<ColumnProps> = ({
  column,
  tasks,
  onUpdateColumnName,
  onDeleteColumn,
  initiallyEditing = false,
}) => {
  const { setNodeRef, isOver } = useDroppable({ id: column.id });

  const [isEditing, setIsEditing] = useState(initiallyEditing);
  const [name, setName] = useState(column.name);
  const navigate = useNavigate();

  const handleSubmit = (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setIsEditing(false);
    if (onUpdateColumnName) {
      onUpdateColumnName(column.id, name);
    }
  };

  const startEditing = () => {
    setIsEditing(true);
  };

  const handleDelete = () => {
    if (onDeleteColumn && window.confirm("このカラムを削除しますか？")) {
      onDeleteColumn(column.id);
    }
  };

  return (
    <Box
      ref={setNodeRef}
      sx={{
        border: isOver ? "2px dashed #4CAF50" : "1px solid #ccc",
        borderRadius: 2,
        p: 2,
        minWidth: 280,
        backgroundColor: "#fafafa",
        boxShadow: 1,
        position: "relative",
      }}
    >
      <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
        {isEditing ? (
          <form onSubmit={handleSubmit} style={{ flexGrow: 1 }}>
            <TextField
              value={name}
              onChange={(e) => setName(e.target.value)}
              size="small"
              autoFocus
              fullWidth
              onBlur={() => handleSubmit()}
            />
          </form>
        ) : (
          <Typography
            variant="h6"
            sx={{ flexGrow: 1, cursor: "pointer" }}
            onClick={startEditing}
          >
            {column.name}
          </Typography>
        )}
        {!isEditing && (
          <>
            <IconButton size="small" onClick={startEditing}>
              <EditIcon fontSize="small" />
            </IconButton>
            <IconButton size="small" onClick={handleDelete}>
              <DeleteIcon fontSize="small" />
            </IconButton>
          </>
        )}
      </Box>

      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} />
      ))}

      <IconButton
        size="small"
        color="primary"
        sx={{ position: "absolute", bottom: 8, right: 8 }}
        onClick={() => navigate(`/tasks/new?columnId=${column.id}`)}
      >
        <AddIcon />
      </IconButton>
    </Box>
  );
};

export default Column;
