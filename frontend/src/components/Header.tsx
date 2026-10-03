import {
  AppBar,
  Box,
  Button,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  type SelectChangeEvent,
  Toolbar,
  Typography,
} from "@mui/material";
import type React from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { setSelectedProject } from "../features/projects/slice";
import type { RootState } from "../store";

const Header: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const projects = useSelector((state: RootState) => state.projects.projects);
  const selectedId = useSelector(
    (state: RootState) => state.projects.selectedId,
  );
  const isLoggedIn = useSelector(
    (state: RootState) => state.session.isLoggedIn,
  );

  if (!isLoggedIn) return null;

  const handleProjectChange = (event: SelectChangeEvent) => {
    const id = event.target.value as string;
    dispatch(setSelectedProject(id));
    navigate("/boards");
  };

  return (
    <AppBar position="static">
      <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
        <Box
          sx={{ display: "flex", alignItems: "center", cursor: "pointer" }}
          onClick={() => navigate("/")}
        >
          <Typography variant="h6">My Kanban App</Typography>

          <FormControl variant="standard" sx={{ minWidth: 200, marginLeft: 4 }}>
            <InputLabel id="project-select-label" sx={{ color: "#fff" }}>
              Project
            </InputLabel>
            <Select
              labelId="project-select-label"
              value={selectedId}
              onChange={handleProjectChange}
              sx={{
                color: "#fff",
                "& .MuiSelect-icon": { color: "#fff" },
                "&:before, &:after": { borderColor: "#fff" },
              }}
              label="Project"
            >
              {projects.map((proj: { id: string; name: string }) => (
                <MenuItem key={proj.id} value={proj.id}>
                  {proj.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
        <Box>
          <Button color="inherit" onClick={() => navigate("/")}>
            Projects
          </Button>
          <Button color="inherit" onClick={() => navigate("/boards")}>
            Boards
          </Button>
          <Button
            color="inherit"
            onClick={() => {
              navigate("/login");
            }}
          >
            Logout
          </Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
