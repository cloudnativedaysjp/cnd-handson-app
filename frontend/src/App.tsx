import { Provider } from "react-redux";
import {
  Navigate,
  Route,
  BrowserRouter as Router,
  Routes,
} from "react-router-dom";
import { legacy } from "./board";
import {
  LegacyColumns,
  LegacyProject,
  LegacyProjectList,
  LegacyTask,
} from "./components/Legacy";
import Login from "./components/Login";
import Logout from "./components/Logout";
import {
  ModernBoard,
  ModernColumns,
  ModernProjectList,
} from "./components/Modern";
import PrivateLayout from "./components/PrivateLayout";
import { store } from "./store";

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/logout" element={<Logout />} />
          {/* データと URL は 2 つの版で共通。違うのは見せ方だけ */}
          <Route element={<PrivateLayout />}>
            <Route
              path="/"
              element={legacy ? <LegacyProjectList /> : <ModernProjectList />}
            />
            <Route
              path="/projects/:id"
              element={legacy ? <LegacyProject /> : <ModernBoard />}
            />
            {/* modern はボードの上に詳細のパネルを開く */}
            <Route
              path="/projects/:id/tasks/:taskId"
              element={legacy ? <LegacyTask /> : <ModernBoard />}
            />
            <Route
              path="/projects/:id/columns"
              element={legacy ? <LegacyColumns /> : <ModernColumns />}
            />
          </Route>
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;
