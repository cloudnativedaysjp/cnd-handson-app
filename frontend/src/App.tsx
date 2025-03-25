import React from 'react';
import { Provider } from 'react-redux';
import { store } from './store';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Board from './components/Borad';
import TaskModal from './components/TaskModal';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Routes>
          <Route path="/" element={<Board />} />
          <Route path="/tasks/:id" element={<TaskModal />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;
