import { Player } from '../lib/protocol';
import { PlayerList } from './PlayerList';

interface Props {
  code: string;
  players: Player[];
  hostId: string;
  playerId: string;
  onStart: () => void;
  onAddAI: () => void;
  onKick: (id: string) => void;
}

export function Lobby({ code, players, hostId, playerId, onStart, onAddAI, onKick }: Props) {
  const isHost = playerId === hostId;

  return (
    <div className="lobby">
      <h2>Game Lobby</h2>
      <div className="game-code">
        <span className="code-label">Code:</span>
        <span className="code-value">{code}</span>
      </div>
      <p className="player-count">{players.length} / 8 players</p>

      <PlayerList
        players={players}
        hostId={hostId}
        playerId={playerId}
        isHost={isHost}
        onKick={onKick}
      />

      {isHost && (
        <div className="lobby-actions">
          <button
            className="btn btn-primary"
            onClick={onStart}
            disabled={players.length < 2}
          >
            Start Game
          </button>
          <button
            className="btn btn-secondary"
            onClick={onAddAI}
            disabled={players.length >= 8}
          >
            Add AI Player
          </button>
        </div>
      )}
      {!isHost && <p className="waiting-text">Waiting for host to start...</p>}
    </div>
  );
}
