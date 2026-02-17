import { useRef, useEffect } from 'react';
import { TURN_DRAW } from '../lib/protocol';
import { DrawingCanvas, DrawingCanvasHandle } from './DrawingCanvas';
import { GuessInput } from './GuessInput';
import { Timer } from './Timer';

interface Props {
  round: number;
  totalRounds: number;
  turnType: number;
  prompt: string;
  timeRemaining: number;
  timeLimit: number;
  waiting: boolean;
  onSubmitDrawing: (dataUrl: string) => void;
  onSubmitGuess: (guess: string) => void;
}

export function GamePlay({
  round, totalRounds, turnType, prompt,
  timeRemaining, timeLimit, waiting,
  onSubmitDrawing, onSubmitGuess,
}: Props) {
  const canvasRef = useRef<DrawingCanvasHandle>(null);

  useEffect(() => {
    if (timeRemaining <= 0 && !waiting && turnType === TURN_DRAW) {
      canvasRef.current?.submit();
    }
  }, [timeRemaining, waiting, turnType]);

  return (
    <div className="gameplay">
      <div className="round-info">
        Round {round + 1} of {totalRounds}
      </div>

      <Timer remaining={timeRemaining} total={timeLimit} />

      {waiting ? (
        <div className="waiting-screen">
          <div className="waiting-spinner" />
          <p>Waiting for other players...</p>
        </div>
      ) : turnType === TURN_DRAW ? (
        <DrawingCanvas ref={canvasRef} prompt={prompt} onSubmit={onSubmitDrawing} />
      ) : (
        <GuessInput drawingDataUrl={prompt} onSubmit={onSubmitGuess} />
      )}
    </div>
  );
}
