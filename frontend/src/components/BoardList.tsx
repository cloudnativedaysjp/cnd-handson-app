
import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../store';
import { Board } from '../features/board/types';
import { Grid, Card, CardContent, Typography, CardActionArea } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const BoardList: React.FC = () => {
  const navigate = useNavigate();
  const boards = useSelector((state: RootState) => state.boards.boards);

  const handleBoardClick = (board: Board) => {
    navigate(`/boards/${board.id}`);
  };

  return (
    <div style={{ padding: '16px' }}>
      <Typography variant="h4" gutterBottom>
        Boards
      </Typography>
      <Grid container spacing={2}>
        {boards.map((board: Board) => (
          <Grid item xs={12} sm={6} md={4} key={board.id}>
            <Card>
              <CardActionArea onClick={() => handleBoardClick(board)}>
                <CardContent>
                  <Typography variant="h5">{board.name}</Typography>
                  <Typography variant="body2" color="textSecondary">
                    {board.description || 'No description available.'}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </div>
  );
};

export default BoardList;
