
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
import { setSelectedProject } from '../features/projects/slice';

const Header: React.FC = () => {

  const navigate = useNavigate();
  const dispatch = useDispatch();
  const projects = useSelector((state: RootState) => state.projects.projects);
  const selectedId = useSelector((state: RootState) => state.projects.selectedId);
  const isLoggedIn = useSelector((state: RootState) => state.session.isLoggedIn);

  if (!isLoggedIn) return null;

  const handleProjectChange = (e: React.ChangeEvent<{ value: unknown }>) => {
    const id = e.target.value as string;
    dispatch(setSelectedProject(id));
    navigate('/boards');
  };

  return (
    <AppBar position="static" sx={{ backgroundColor: '#111111' }}>
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
            Projects
          </Button>
          <Button color="inherit" onClick={() => navigate('/boards')}>
            Boards
          </Button>
          <Button
            color="inherit"
            onClick={() => {
              navigate('/login');
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
