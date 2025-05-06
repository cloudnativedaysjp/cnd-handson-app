
import React, { useEffect } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Box,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { RootState } from '../store';
import { setProjects, setSelectedProject } from '../features/projects/slice';
import { Project } from '../features/projects/types';

const Header: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const projects = useSelector((state: RootState) => state.projects.list);
  const selectedId = useSelector((state: RootState) => state.projects.selectedId);


  useEffect(() => {
    const dummy: Project[] = [
      { id: 'proj-1', name: 'Project Alpha' },
      { id: 'proj-2', name: 'Project Beta' },
      { id: 'proj-3', name: 'Project Gamma' },
    ];
    dispatch(setProjects(dummy));
  }, [dispatch]);

  const handleProjectChange = (e: React.ChangeEvent<{ value: unknown }>) => {
    const id = e.target.value as string;
    dispatch(setSelectedProject(id));
    navigate('/boards');
  };

  return (
    <AppBar position="static">
      <Toolbar sx={{ display: 'flex', justifyContent: 'space-between' }}>
        <Box
          sx={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}
          onClick={() => navigate('/')}
        >
          <Typography variant="h6">My Kanban App</Typography>

            <FormControl variant="standard" sx={{ minWidth: 200, marginLeft: 4 }}>
            <InputLabel id="project-select-label" sx={{ color: '#fff' }}>
                Project
            </InputLabel>
            <Select
                labelId="project-select-label"
                value={selectedId}
                onChange={handleProjectChange}
                sx={{
                color: '#fff',
                '& .MuiSelect-icon': { color: '#fff' },
                '&:before, &:after': { borderColor: '#fff' },
                }}
                label="Project"
            >
                {projects.map((proj: { id: any; name: any; }) => (
                <MenuItem key={proj.id} value={proj.id}>
                    {proj.name}
                </MenuItem>
                ))}
            </Select>
            </FormControl>
        </Box>
        <Box>
          <Button color="inherit" onClick={() => navigate('/')}>
            Boards
          </Button>
          <Button color="inherit" onClick={() => navigate('/settings')}>
            Settings
          </Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;
