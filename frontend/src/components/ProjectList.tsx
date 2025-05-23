import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  Grid,
  IconButton,
  Menu,
  MenuItem,
  Typography,
} from "@mui/material";
// src/components/ProjectList.tsx
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { setSelectedProject } from "../features/projects/slice";
import { deleteProject, updateProject } from "../features/projects/slice";
import type { Project } from "../features/projects/types";
import type { RootState } from "../store";
import AddButton from "./AddButton";
import MoreMenu, { type MoreMenuOption } from "./MoreMenu";
import ProjectForm from "./ProjectForm";

const ProjectList: React.FC = () => {
  const projects = useSelector((state: RootState) => state.projects.projects);
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [editOpen, setEditOpen] = React.useState(false);
  const [editTarget, setEditTarget] = React.useState<Project | null>(null);

  const openEdit = (id: string) => {
    const projectToEdit = projects.find((p) => p.id === id);
    if (projectToEdit) {
      setEditTarget(projectToEdit);
      setEditOpen(true);
    }
  };

  const closeEdit = () => {
    setEditOpen(false);
    setEditTarget(null);
  };

  const handleSave = (updated: Project) => {
    dispatch(updateProject(updated));
    closeEdit();
  };

  const options: MoreMenuOption<string>[] = [
    { label: "編集", onClick: (id) => navigate(`/projects/edit/${id}`) },
    {
      label: "削除",
      onClick: (id) => {
        if (window.confirm("本当に削除しますか？")) {
          dispatch(deleteProject(id));
        }
      },
    },
  ];

  const handleClickCard = (projectId: string) => {
    dispatch(setSelectedProject(projectId));
    navigate("/boards");
  };

  const onAdd = () => {
    navigate("/projects/new");
  };

  return (
    <Box sx={{ p: 2 }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          padding: "8px",
        }}
      >
        <Typography variant="h4" gutterBottom>
          Projects
        </Typography>
        <AddButton label="Projectを追加" onClick={onAdd} />
      </div>
      <Grid container spacing={2}>
        {projects.map((proj) => (
          <Grid item xs={12} sm={6} md={4} key={proj.id}>
            <Card
              sx={{
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
              }}
            >
              <CardActionArea onClick={() => handleClickCard(proj.id)}>
                <CardContent>
                  <Box>
                    <Typography variant="h5">{proj.name}</Typography>
                    {proj.description && (
                      <Typography variant="body2" color="textSecondary">
                        {proj.description}
                      </Typography>
                    )}
                  </Box>
                </CardContent>
              </CardActionArea>
              <MoreMenu id={proj.id} options={options} />
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default ProjectList;
