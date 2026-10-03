import { Provider } from "react-redux";
import {
  Navigate,
  Route,
  BrowserRouter as Router,
  Routes,
} from "react-router-dom";
import BoardForm from "./components/BoardForm";
import BoardList from "./components/BoardList";
import Borad from "./components/Borad";
import { LegacyProject, LegacyProjectList } from "./components/Legacy";
import Login from "./components/Login";
import Logout from "./components/Logout";
import PrivateLayout from "./components/PrivateLayout";
import ProjectForm from "./components/ProjectForm";
import ProjectList from "./components/ProjectList";
import TaskEditPage from "./components/TaskEditPage";
import NewTaskPage from "./components/TaskNewPage";
import { store } from "./store";

// 入口のイメージはビルド時の VARIANT で legacy / modern を切り替える。bff/Dockerfile が VITE_VARIANT に渡す
const legacy = import.meta.env.VITE_VARIANT === "legacy";

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/logout" element={<Logout />} />
          <Route element={<PrivateLayout />}>
            <Route path="/boards" element={<BoardList />} />
            <Route
              path="/"
              element={legacy ? <LegacyProjectList /> : <ProjectList />}
            />
            <Route path="/projects/:id" element={<LegacyProject />} />
            <Route path="/projects/new" element={<ProjectForm />} />
            <Route path="/projects/edit/:id" element={<ProjectForm />} />
            <Route path="/boards/:id" element={<Borad />} />
            <Route path="/boards/add" element={<BoardForm />} />
            <Route path="/boards/edit/:id" element={<BoardForm />} />
            <Route path="/tasks/new" element={<NewTaskPage />} />
            <Route path="/tasks/edit/:id" element={<TaskEditPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;
