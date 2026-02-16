import { useState } from 'react';

interface Props {
  drawingDataUrl: string;
  onSubmit: (guess: string) => void;
}

export function GuessInput({ drawingDataUrl, onSubmit }: Props) {
  const [guess, setGuess] = useState('');

  const handleSubmit = () => {
    if (guess.trim()) {
      onSubmit(guess.trim());
    }
  };

  return (
    <div className="guess-area">
      <div className="guess-prompt">What is this drawing?</div>

      <div className="guess-drawing">
        {drawingDataUrl ? (
          <img src={drawingDataUrl} alt="Drawing to guess" className="guess-image" />
        ) : (
          <div className="guess-blank">No drawing submitted</div>
        )}
      </div>

      <div className="guess-form">
        <input
          type="text"
          placeholder="Type your guess..."
          value={guess}
          onChange={e => setGuess(e.target.value)}
          maxLength={100}
          autoFocus
          onKeyDown={e => e.key === 'Enter' && handleSubmit()}
        />
        <button
          className="btn btn-primary"
          onClick={handleSubmit}
          disabled={!guess.trim()}
        >
          Submit Guess
        </button>
      </div>
    </div>
  );
}
