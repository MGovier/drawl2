import { useState } from 'react';

interface Props {
  onCreateGame: (name: string, password: string) => void;
  onJoinGame: (name: string, code: string) => void;
  error: string;
}

export function Home({ onCreateGame, onJoinGame, error }: Props) {
  const [name, setName] = useState('');
  const [code, setCode] = useState('');
  const [password, setPassword] = useState('');
  const [mode, setMode] = useState<'menu' | 'create' | 'join'>('menu');

  const handleCreate = () => {
    if (name.trim()) onCreateGame(name.trim(), password);
  };

  const handleJoin = () => {
    if (name.trim() && code.trim()) onJoinGame(name.trim(), code.trim());
  };

  return (
    <div className="home">
      <h1 className="title">Drawl</h1>
      <p className="subtitle">A multiplayer drawing game</p>

      {error && <div className="error-msg">{error}</div>}

      {mode === 'menu' && (
        <div className="menu-buttons">
          <button className="btn btn-primary" onClick={() => setMode('create')}>
            Create Game
          </button>
          <button className="btn btn-secondary" onClick={() => setMode('join')}>
            Join Game
          </button>
        </div>
      )}

      {mode === 'create' && (
        <div className="form">
          <input
            type="text"
            placeholder="Your name"
            value={name}
            onChange={e => setName(e.target.value)}
            maxLength={20}
            autoFocus
          />
          <input
            type="password"
            placeholder="Game password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && handleCreate()}
          />
          <div className="form-buttons">
            <button className="btn btn-primary" onClick={handleCreate} disabled={!name.trim()}>
              Create
            </button>
            <button className="btn btn-ghost" onClick={() => setMode('menu')}>
              Back
            </button>
          </div>
        </div>
      )}

      {mode === 'join' && (
        <div className="form">
          <input
            type="text"
            placeholder="Your name"
            value={name}
            onChange={e => setName(e.target.value)}
            maxLength={20}
            autoFocus
          />
          <input
            type="text"
            placeholder="Game code"
            value={code}
            onChange={e => setCode(e.target.value.toUpperCase())}
            maxLength={5}
            onKeyDown={e => e.key === 'Enter' && handleJoin()}
          />
          <div className="form-buttons">
            <button
              className="btn btn-primary"
              onClick={handleJoin}
              disabled={!name.trim() || !code.trim()}
            >
              Join
            </button>
            <button className="btn btn-ghost" onClick={() => setMode('menu')}>
              Back
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
