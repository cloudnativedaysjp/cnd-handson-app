import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Project } from './types';

interface ProjectsState {
  list: Project[];
  selectedId: string | '';
}

const initialState: ProjectsState = {
  list: [],
  selectedId: '',
};

const projectsSlice = createSlice({
  name: 'projects',
  initialState,
  reducers: {
    setProjects(state, action: PayloadAction<Project[]>) {
      state.list = action.payload;
    },
    setSelectedProject(state, action: PayloadAction<string>) {
      state.selectedId = action.payload;
    },
  },
});

export const { setProjects, setSelectedProject } = projectsSlice.actions;
export default projectsSlice.reducer;
