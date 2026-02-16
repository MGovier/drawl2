import { Chain, Player, TURN_DRAW } from '../lib/protocol';

interface Props {
  chain: Chain;
  chainIndex: number;
  players: Player[];
  voteMode?: boolean;
  voted?: boolean;
  onVoteSuccess?: () => void;
  favouriteMode?: boolean;
  favouriteKey?: string;
  onFavourite?: (entryKey: string) => void;
}

export function ChainCard({ chain, chainIndex, players, voteMode, voted, onVoteSuccess, favouriteMode, favouriteKey, onFavourite }: Props) {
  const getPlayerName = (id: string) => {
    return players.find(p => p.id === id)?.name || 'Unknown';
  };

  const owner = getPlayerName(chain.ownerId);

  return (
    <div className={`chain-card ${voted ? 'chain-voted' : ''}`}>
      <div className="chain-header">
        <strong>{owner}'s chain</strong>
        <div className="chain-header-right">
          <span className="chain-word">Original: "{chain.originalWord}"</span>
          {voteMode && (
            <button
              className={`btn-vote ${voted ? 'btn-vote-active' : ''}`}
              onClick={onVoteSuccess}
              title="This chain survived!"
            >
              <span className="vote-icon">{voted ? '\u{1F44D}' : '\u{1F44D}'}</span>
            </button>
          )}
        </div>
      </div>
      <div className="chain-entries">
        {chain.entries.map((entry, i) => {
          const entryKey = `${chainIndex}:${i}`;
          const isDrawing = entry.type === TURN_DRAW;
          const isFavSelected = favouriteKey === entryKey;

          return (
            <div key={i} className="chain-entry">
              <div className="chain-entry-header">
                <span className="chain-player">{getPlayerName(entry.playerId)}</span>
                {favouriteMode && isDrawing && entry.drawing && onFavourite && (
                  <button
                    className={`btn-fav ${isFavSelected ? 'btn-fav-active' : ''}`}
                    onClick={() => onFavourite(entryKey)}
                    title="Pick as favourite drawing"
                  >
                    {isFavSelected ? '\u2605' : '\u2606'}
                  </button>
                )}
              </div>
              {isDrawing ? (
                entry.drawing ? (
                  <img src={entry.drawing} alt="drawing" className="chain-drawing" />
                ) : (
                  <div className="chain-blank">No drawing</div>
                )
              ) : (
                <div className="chain-guess">"{entry.guess || '???'}"</div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
