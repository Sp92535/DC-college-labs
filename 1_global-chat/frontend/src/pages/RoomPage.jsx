import { useState, useEffect } from 'react';
import { useParams, useLocation } from 'react-router-dom';
import axios from 'axios';
import './RoomPage.css';
import { io } from 'socket.io-client';
import { SERVER } from '../utils/config';

const RoomPage = () => {
  const { roomCode } = useParams();

  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState('');
  const [socket, setSocket] = useState(null);

  const { state } = useLocation();
  const username = state?.username;

  // Send message to backend and update the state
  const sendMessage = () => {
    try {
      console.log(username);
      console.log(roomCode + " " + username + " " + message);

      socket.emit("newMessage", { room_id: roomCode, username: username, msg: message });
      setMessage("");
    } catch (error) {

    }
  };

  // Fetch users and messages from backend
  useEffect(() => {
    const fetchData = async () => {
      try {

        const response = await axios.get(`${SERVER}/messages/${roomCode}`);
        setMessages(response.data.messages);
        console.log(response.data);

      } catch (error) {

      }
    };

    const socketInstance = io(SERVER);
    setSocket(socketInstance);

    socketInstance.emit("joinRoom", roomCode);

    socketInstance.on("message", (msg) => {
      setMessages((prevMessages) => [...prevMessages, msg]);
    });

    fetchData();
    return () => {
      socketInstance.disconnect();  // Clean up the socket connection when the component unmounts
    };
  }, [roomCode]);

  return (
    <div className="room-page">
      <div className="chat-container">
        <h2 className="room-code">Room Code: {roomCode}</h2>


        <div className="chat-window">
          {messages.length > 0 ? (
            messages.map((msg, index) => (
              <div key={index} className="message">
                <strong>{msg.username}:</strong> {msg.msg} <span className="timestamp">({msg.timestamp})</span>
              </div>
            ))
          ) : (
            <p>No messages yet.</p>
          )}
        </div>

        <div className="message-input">
          <input
            type="text"
            placeholder="Type a message..."
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
          <button onClick={sendMessage}>Send</button>
        </div>
      </div>
    </div>
  );
};

export default RoomPage;
