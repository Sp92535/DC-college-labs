import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import NameInputForm from './pages/NameInputForm';
import RoomPage from './pages/RoomPage';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<NameInputForm />} />
        <Route path="/room/:roomCode" element={<RoomPage />} />
      </Routes>
    </Router>
  );
};

export default App;
