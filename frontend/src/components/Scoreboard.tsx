import { Player } from '../lib/protocol';

interface Props {
  players: Player[];
  scores: Record<string, number>;
  compact?: boolean;
}

export function Scoreboard({ players, scores, compact }: Props) {
  const sorted = [...players].sort((a, b) => (scores[b.id] || 0) - (scores[a.id] || 0));

  return (
    <div className={`scoreboard ${compact ? 'scoreboard-compact' : ''}`}>
      {!compact && <h3 className="scoreboard-title">Scoreboard</h3>}
      <div className="scoreboard-list">
        {sorted.map((p, i) => (
          <div key={p.id} className={`scoreboard-row ${i === 0 && (scores[p.id] || 0) > 0 ? 'scoreboard-leader' : ''}`}>
            <span className="scoreboard-rank">{i + 1}</span>
            <span className="scoreboard-name">{p.name}</span>
            <span className="scoreboard-score">{scores[p.id] || 0}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
