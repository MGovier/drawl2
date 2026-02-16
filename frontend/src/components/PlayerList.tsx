import { Player } from '../lib/protocol';

interface Props {
  players: Player[];
  hostId: string;
  playerId: string;
  isHost: boolean;
  onKick?: (id: string) => void;
}

export function PlayerList({ players, hostId, playerId, isHost, onKick }: Props) {
  return (
    <div className="player-list">
      {players.map(p => (
        <div key={p.id} className={`player-item ${p.id === playerId ? 'you' : ''}`}>
          <span className="player-name">
            {p.name}
            {p.id === hostId && <span className="badge host">Host</span>}
            {p.type === 1 && <span className="badge ai">AI</span>}
            {p.id === playerId && <span className="badge you">You</span>}
          </span>
          {isHost && p.id !== playerId && onKick && (
            <button className="btn-kick" onClick={() => onKick(p.id)} title="Kick player">
              &times;
            </button>
          )}
        </div>
      ))}
    </div>
  );
}
