import React, { useState } from 'react';
import Chat from './Chat';

function App() {
  const [token, setToken] = useState(localStorage.getItem('jwt_token') || '');
  const [boardId, setBoardId] = useState(1);
  const [userId, setUserId] = useState(1);
  const [username, setUsername] = useState('User1');

  if (!token) {
    return (
      <main style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        backgroundColor: '#f5f5f5'
      }}>
        <section style={{
          backgroundColor: 'white',
          padding: '32px',
          borderRadius: '8px',
          boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
          width: '400px'
        }}>
          <h1 style={{
            marginBottom: '24px',
            color: '#2c3e50',
            textAlign: 'center'
          }}>
            TaskManager Chat Login
          </h1>
          <form onSubmit={(e) => {
            e.preventDefault();
            localStorage.setItem('jwt_token', 'demo-jwt-token-12345');
            window.location.reload();
          }} style={{
            display: 'flex',
            flexDirection: 'column',
            gap: '16px'
          }}>
            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              style={{
                padding: '12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
                fontSize: '16px'
              }}
            />
            <input
              type="number"
              placeholder="Board ID"
              value={boardId}
              onChange={(e) => setBoardId(Number(e.target.value))}
              style={{
                padding: '12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
                fontSize: '16px'
              }}
            />
            <input
              type="number"
              placeholder="User ID"
              value={userId}
              onChange={(e) => setUserId(Number(e.target.value))}
              style={{
                padding: '12px',
                border: '1px solid #ced4da',
                borderRadius: '4px',
                fontSize: '16px'
              }}
            />
            <button type="submit" style={{
              padding: '12px',
              backgroundColor: '#007bff',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              fontSize: '16px',
              fontWeight: '600',
              cursor: 'pointer'
            }}>
              Enter Chat
            </button>
          </form>
        </section>
      </main>
    );
  }

  return (
    <Chat 
      token={token}
      boardId={boardId}
      userId={userId}
      username={username}
    />
  );
}

export default App;