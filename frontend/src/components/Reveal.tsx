import { useState } from 'react';
import { Chain, Player } from '../lib/protocol';
import { ChainCard } from './ChainCard';
import { Scoreboard } from './Scoreboard';

interface Props {
  chains: Chain[];
  players: Player[];
  scores: Record<string, number>;
  playerId: string;
  hostId: string;
  favDrawing: string;
  votingDone: boolean;
  waiting: boolean;
  onSubmitVotes: (successChains: number[], favDrawing: string) => void;
  onPlayAgain: () => void;
  onHome: () => void;
}

export function Reveal({ chains, players, scores, playerId, hostId, favDrawing, votingDone, waiting, onSubmitVotes, onPlayAgain, onHome }: Props) {
  const [selectedChains, setSelectedChains] = useState<Set<number>>(new Set());
  const [selectedFav, setSelectedFav] = useState('');
  const [submitted, setSubmitted] = useState(false);
  const isHost = playerId === hostId;
  const isWaiting = submitted || waiting;

  const toggleChain = (idx: number) => {
    setSelectedChains(prev => {
      const next = new Set(prev);
      if (next.has(idx)) next.delete(idx);
      else next.add(idx);
      return next;
    });
  };

  const handleSubmit = () => {
    setSubmitted(true);
    onSubmitVotes(Array.from(selectedChains), selectedFav);
  };

  return (
    <div className="reveal">
      <h2>Game Over!</h2>
      <p className="reveal-subtitle">
        {votingDone
          ? 'Final results are in!'
          : isWaiting
            ? 'Waiting for others to finish voting...'
            : 'Review the chains, vote on successes, and pick your favourite drawing!'}
      </p>

      <div className="review-layout">
        <div className="review-chains">
          <div className="chains-list">
            {chains.map((chain, i) => (
              <ChainCard
                key={i}
                chain={chain}
                chainIndex={i}
                players={players}
                voteMode={!isWaiting && !votingDone}
                voted={selectedChains.has(i)}
                onVoteSuccess={() => toggleChain(i)}
                favouriteMode={!isWaiting && !votingDone}
                favouriteKey={selectedFav}
                onFavourite={setSelectedFav}
              />
            ))}
          </div>

          {!isWaiting && !votingDone && (
            <button
              className="btn btn-primary btn-submit-votes"
              onClick={handleSubmit}
            >
              Submit Votes
            </button>
          )}

          {isWaiting && !votingDone && (
            <div className="waiting-screen">
              <div className="waiting-spinner" />
              <p className="waiting-text">Waiting for others to finish voting...</p>
            </div>
          )}

          <div className="reveal-actions">
            {isHost ? (
              <button className="btn btn-primary" onClick={onPlayAgain}>
                Play Again
              </button>
            ) : (
              <p className="waiting-text">Waiting for host to start a new game...</p>
            )}
            <button className="btn btn-secondary" onClick={onHome}>
              Home
            </button>
          </div>
        </div>

        <Scoreboard players={players} scores={scores} />
      </div>
    </div>
  );
}
