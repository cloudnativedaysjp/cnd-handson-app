import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  Grid,
  Typography,
} from "@mui/material";
import type React from "react";
import { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { useNavigate } from "react-router-dom";
import { setSelectedProject } from "../features/projects/slice";
import type { Project } from "../features/projects/types";
import { api } from "../utils/api";

const ProjectList: React.FC = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [failed, setFailed] = useState(false);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await api<{ projects: Project[] }>("/api/projects");
        setProjects(res.projects);
      } catch {
        setFailed(true);
      }
    };
    load();
  }, []);

  const handleClickCard = (projectId: string) => {
    dispatch(setSelectedProject(projectId));
    navigate("/boards");
  };

  return (
    <Box sx={{ p: 2 }}>
      <Typography variant="h4" gutterBottom sx={{ p: 1 }}>
        Projects
      </Typography>
      {failed && (
        <Typography color="error">
          プロジェクトを読み込めませんでした
        </Typography>
      )}
      <Grid container spacing={2} data-testid="project-list">
        {projects.map((proj) => (
          <Grid item xs={12} sm={6} md={4} key={proj.id}>
            <Card data-testid="project-item">
              <CardActionArea onClick={() => handleClickCard(proj.id)}>
                <CardContent>
                  <Typography variant="h5">{proj.name}</Typography>
                  {proj.description && (
                    <Typography variant="body2" color="textSecondary">
                      {proj.description}
                    </Typography>
                  )}
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default ProjectList;
