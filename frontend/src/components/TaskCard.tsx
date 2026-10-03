import { useDraggable } from "@dnd-kit/core";
import { Box, Card, CardContent, Typography } from "@mui/material";
// src/components/TaskCard.tsx
import type React from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { deleteTask } from "../features/tasks/slice";
import type { Task } from "../features/tasks/types";
import MoreMenu, { type MoreMenuOption } from "./MoreMenu";

interface TaskCardProps {
  task: Task;
}

const TaskCard: React.FC<TaskCardProps> = ({ task }) => {
  const navigate = useNavigate();
  const dispatch = useDispatch();

  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: task.id,
    data: { task },
  });

  const style = {
    transform: transform
      ? `translate3d(${transform.x}px, ${transform.y}px, 0)`
      : undefined,
    transition: "transform 200ms ease",
    margin: "8px 0",
    cursor: "grab",
  };

  const options: MoreMenuOption<string>[] = [
    { label: "編集", onClick: (id) => navigate(`/tasks/edit/${id}`) },
    {
      label: "削除",
      onClick: (id) => {
        if (window.confirm("本当に削除しますか？")) {
          dispatch(deleteTask(id));
        }
      },
    },
  ];

  return (
    <div ref={setNodeRef} style={style} {...attributes} {...listeners}>
      <Card
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
        }}
      >
        <CardContent>
          <Box>
            <Typography variant="h6">{task.title}</Typography>
            <Typography variant="body2" color="textSecondary">
              {task.description || "説明なし"}
            </Typography>
            <Typography variant="caption">Status: {task.status}</Typography>
          </Box>
        </CardContent>
        <MoreMenu id={task.id} options={options} />
      </Card>
    </div>
  );
};

export default TaskCard;
