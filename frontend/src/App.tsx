import React from 'react';
import { Provider } from 'react-redux';
import { store } from './store';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import BoardList from './components/BoardList';
import Borad from './components/Borad';
import TaskModal from './components/TaskModal';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Routes>
          <Route path="/" element={<BoardList />} />
          <Route path="/boards/:id" element={<Borad />} />
          <Route path="/tasks/:id" element={<TaskModal open={false} task={null} onClose={function (): void {
            throw new Error('Function not implemented.');
          } } />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;
