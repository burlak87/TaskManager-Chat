import React, { useState, useEffect, useRef, useCallback } from 'react';
import { wsService } from './websocket';
import { Message } from './types';

interface ChatProps {
  token: string;
  boardId: number;
  userId: number;
  username: string;
}

const Chat: React.FC<ChatProps> = ({ token, boardId, userId, username }) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [participants, setParticipants] = useState<string[]>([username]);
  const [connectionStatus, setConnectionStatus] = useState('Connecting');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    wsService.connect(token, boardId);

    wsService.on('connect', () => {
      setConnectionStatus('Connected');
      console.log('Connected to chat');
    });

    wsService.on('disconnect', () => {
      setConnectionStatus('Disconnected');
    });

    wsService.on('error', (error) => {
      console.error('WebSocket error:', error);
      setConnectionStatus('Error');
    });

    wsService.on('message', (message: Message) => {
      setMessages(prev => [...prev, message]);
      
      if (!participants.includes(message.username)) {
        setParticipants(prev => [...prev, message.username]);
      }
      
      if (message.user_id !== userId && 'Notification' in window && Notification.permission === 'granted') {
        new Notification('New message', {
          body: `${message.username}: ${message.content.substring(0, 50)}`
        });
      }
    });

    const loadHistory = async () => {
      try {
        const response = await fetch(`/api/messages?board_id=${boardId}&limit=50`);
        if (response.ok) {
          const data = await response.json();
          setMessages(data);
          
          const users = Array.from(new Set(data.map((msg: Message) => msg.username)));
          setParticipants(users);
        }
      } catch (error) {
        console.error('Failed to load message history:', error);
      }
    };

    loadHistory();

    if ('Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission();
    }

    return () => {
      wsService.disconnect();
    };
  }, [token, boardId, userId]);

 
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = useCallback((e: React.FormEvent) => {
    e.preventDefault();
    
    if (!inputMessage.trim() || !wsService.isConnected()) return;

    wsService.sendMessage({
      board_id: boardId,
      content: inputMessage
    });

    setInputMessage('');
  }, [inputMessage, boardId]);

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      const form = e.currentTarget.closest('form');
      if (form) {
        const submitEvent = new Event('submit', { cancelable: true });
        form.dispatchEvent(submitEvent);
        if (!submitEvent.defaultPrevented) {
          handleSendMessage(e as any);
        }
      }
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
  };

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'Connected': return '#4CAF50';
      case 'Connecting': return '#FF9800';
      case 'Disconnected': return '#F44336';
      default: return '#757575';
    }
  };

  return (
    <main style={{
      display: 'flex',
      flexDirection: 'column',
      height: '100vh',
      maxWidth: '1200px',
      margin: '0 auto',
      backgroundColor: '#f5f5f5'
    }}>
      <header style={{
        backgroundColor: '#2c3e50',
        color: 'white',
        padding: '16px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
      }}>
        <h1 style={{ fontSize: '18px', fontWeight: '500' }}>Chat - Board #{boardId}</h1>
        <span style={{
          padding: '4px 12px',
          borderRadius: '12px',
          fontSize: '14px',
          backgroundColor: getStatusColor(),
          color: 'white'
        }}>
          {connectionStatus}
        </span>
      </header>

      <section style={{
        display: 'flex',
        flex: 1,
        overflow: 'hidden'
      }}>
        <section style={{
          flex: 1,
          padding: '16px',
          overflowY: 'auto',
          backgroundColor: 'white',
          display: 'flex',
          flexDirection: 'column',
          gap: '12px'
        }}>
          {messages.map((msg) => (
            <article key={msg.id} style={{
              maxWidth: '70%',
              padding: '12px',
              borderRadius: '12px',
              alignSelf: msg.user_id === userId ? 'flex-end' : 'flex-start',
              backgroundColor: msg.user_id === userId ? '#007bff' : '#e9ecef',
              color: msg.user_id === userId ? 'white' : '#212529'
            }}>
              <header style={{
                display: 'flex',
                justifyContent: 'space-between',
                marginBottom: '4px',
                fontSize: '14px'
              }}>
                <strong style={{
                  fontWeight: '600',
                  color: msg.user_id === userId ? 'rgba(255,255,255,0.9)' : 'inherit'
                }}>
                  {msg.username}
                </strong>
                <time style={{
                  opacity: '0.7',
                  fontSize: '12px'
                }}>
                  {formatTime(msg.created_at)}
                </time>
              </header>
              <p style={{
                wordWrap: 'break-word',
                lineHeight: '1.4'
              }}>
                {msg.content}
              </p>
            </article>
          ))}
          <div ref={messagesEndRef} />
        </section>

        <aside style={{
          width: '250px',
          backgroundColor: 'white',
          borderLeft: '1px solid #dee2e6',
          padding: '16px',
          overflowY: 'auto'
        }}>
          <h2 style={{
            margin: '0 0 16px 0',
            color: '#2c3e50',
            fontSize: '16px'
          }}>
            Participants ({participants.length})
          </h2>
          <ul style={{
            listStyle: 'none',
            padding: '0',
            margin: '0'
          }}>
            {participants.map((participant) => (
              <li key={participant} style={{
                padding: '8px',
                borderBottom: '1px solid #eee',
                color: '#34495e',
                backgroundColor: participant === username ? '#e8f4fd' : 'transparent',
                fontWeight: participant === username ? '600' : '400',
                color: participant === username ? '#007bff' : 'inherit'
              }}>
                {participant} {participant === username && '(You)'}
              </li>
            ))}
          </ul>
        </aside>
      </section>

      <form onSubmit={handleSendMessage} style={{
        display: 'flex',
        gap: '8px',
        padding: '16px',
        backgroundColor: 'white',
        borderTop: '1px solid #dee2e6'
      }}>
        <textarea
          value={inputMessage}
          onChange={(e) => setInputMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          placeholder="Type your message here..."
          rows={3}
          disabled={!wsService.isConnected()}
          style={{
            flex: 1,
            padding: '12px',
            border: '1px solid #ced4da',
            borderRadius: '8px',
            fontSize: '16px',
            resize: 'none',
            fontFamily: 'inherit'
          }}
        />
        <button 
          type="submit" 
          disabled={!inputMessage.trim() || !wsService.isConnected()}
          style={{
            padding: '12px 24px',
            backgroundColor: !inputMessage.trim() || !wsService.isConnected() ? '#6c757d' : '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            fontWeight: '600',
            cursor: !inputMessage.trim() || !wsService.isConnected() ? 'not-allowed' : 'pointer',
            transition: 'background-color 0.2s'
          }}
          onMouseEnter={(e) => {
            if (inputMessage.trim() && wsService.isConnected()) {
              e.currentTarget.style.backgroundColor = '#0056b3';
            }
          }}
          onMouseLeave={(e) => {
            if (inputMessage.trim() && wsService.isConnected()) {
              e.currentTarget.style.backgroundColor = '#007bff';
            }
          }}
        >
          Send
        </button>
      </form>
    </main>
  );
};

export default Chat;