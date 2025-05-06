import React from 'react';
import { Provider, useSelector } from 'react-redux';
import { store } from './store';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import BoardList from './components/BoardList';
import Borad from './components/Borad';
import TaskModal from './components/TaskModal';
import Header from './components/Header';
import Login from './components/Login';

const PrivateRoute: React.FC<{ children: JSX.Element }> = ({ children }) => {
  const isLoggedIn = useSelector((state: RootState) => state.session.isLoggedIn);
  return isLoggedIn ? children : <Navigate to="/login" replace />;
};

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Header />
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/" element={ <PrivateRoute><BoardList /></PrivateRoute>} />
          <Route path="/boards/:id" element={<PrivateRoute><Borad /></PrivateRoute>} />
          <Route path="/tasks/:id" element={<PrivateRoute><TaskModal open={false} task={null} onClose={function (): void {
            throw new Error('Function not implemented.');
          } } /></PrivateRoute>} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;
