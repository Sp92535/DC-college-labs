import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import './NameInputForm.css';
import { SERVER } from '../utils/config';


const NameInputForm = () => {

  const navigate = useNavigate();
  const [username, setUsername] = useState("")
  const [roomId, setRoomId] = useState("");

  const handleCreateRoom = async (e) => {
    e.preventDefault();
    try {
      console.log("YOOO");
      
      const response = await axios.get(`${SERVER}/create-room`);
      setRoomId(response.data.room_id);
    } catch (error) {
      console.error("Error creating room:", error);
    }
  };

  const handleJoinRoom = (e) => {
    e.preventDefault();
    if(!roomId) return;
    navigate(`/room/${roomId}`, { state: { username } })
    console.log("YOOOOOOOOOOO");
    
  };

  return (
    <div className="name-input-form">
      <div className="form-container">
        <h1 className="title">Enter Room</h1>
        <form>
          <div className="input-group">
            <label htmlFor="name">Username</label>
            <input
              id="name"
              type="text"
              placeholder="Enter username"
              name="name"
              value={username}
              required
              onChange={(e) => setUsername(e.target.value)}
            />
            <hr />
            <button onClick={handleCreateRoom} className="submit-btn">
              Create Room
            </button>
          </div>
          <div className="input-group">
            <label htmlFor="roomCode">Room ID</label>
            <input
              id="roomCode"
              type="text"
              placeholder="Enter room ID"
              name="roomCode"
              value={roomId}
              onChange={(e) => setRoomId(e.target.value)}
            />
            <button onClick={handleJoinRoom} className="submit-btn">
              Join Room
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default NameInputForm;