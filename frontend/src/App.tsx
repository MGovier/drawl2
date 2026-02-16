import { useCallback } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { useGameState } from './hooks/useGameState';
import { Home } from './components/Home';
import { Lobby } from './components/Lobby';
import { GamePlay } from './components/GamePlay';
import { Reveal } from './components/Reveal';
import { AIError } from './components/AIError';
import {
  MSG_ADD_AI, MSG_START_GAME, MSG_SUBMIT_DRAWING,
  MSG_SUBMIT_GUESS, MSG_KICK_PLAYER,
  MSG_SUBMIT_VOTES, MSG_PLAY_AGAIN,
} from './lib/protocol';

export default function App() {
  const [state, dispatch] = useGameState();

  const { connect, send, disconnect } = useWebSocket(
    useCallback((msg) => dispatch({ type: 'ws_message', msg }), [dispatch])
  );

  const handleCreateGame = async (name: string, password: string) => {
    const res = await fetch('/api/games', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ playerName: name, password }),
    });
    if (!res.ok) {
      const err = await res.json();
      dispatch({ type: 'ws_message', msg: { type: 'error', data: { message: err.error } } });
      return;
    }
    const data = await res.json();
    dispatch({ type: 'set_connection', code: data.code, playerId: data.playerId });
    connect(data.token, data.code);
  };

  const handleJoinGame = async (name: string, code: string) => {
    const res = await fetch('/api/games/join', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ playerName: name, code }),
    });
    if (!res.ok) {
      const err = await res.json();
      dispatch({ type: 'ws_message', msg: { type: 'error', data: { message: err.error } } });
      return;
    }
    const data = await res.json();
    dispatch({ type: 'set_connection', code, playerId: data.playerId });
    connect(data.token, code);
  };

  const handleGoHome = () => {
    disconnect();
    dispatch({ type: 'set_screen', screen: 'home' });
  };

  return (
    <div className="app">
      {state.screen === 'home' && (
        <Home
          onCreateGame={handleCreateGame}
          onJoinGame={handleJoinGame}
          error={state.error}
        />
      )}

      {state.screen === 'lobby' && (
        <Lobby
          code={state.code}
          players={state.players}
          hostId={state.hostId}
          playerId={state.playerId}
          onStart={() => send(MSG_START_GAME)}
          onAddAI={() => send(MSG_ADD_AI)}
          onKick={(id) => send(MSG_KICK_PLAYER, { playerId: id })}
        />
      )}

      {state.screen === 'playing' && (
        <GamePlay
          round={state.round}
          totalRounds={state.totalRounds}
          turnType={state.turnType}
          prompt={state.prompt}
          timeRemaining={state.timeRemaining}
          timeLimit={state.timeLimit}
          waiting={state.waiting}
          onSubmitDrawing={(drawing) => send(MSG_SUBMIT_DRAWING, { drawing })}
          onSubmitGuess={(guess) => send(MSG_SUBMIT_GUESS, { guess })}
        />
      )}

      {state.screen === 'reveal' && (
        <Reveal
          chains={state.chains}
          players={state.players}
          scores={state.scores}
          playerId={state.playerId}
          hostId={state.hostId}
          favDrawing={state.favDrawing}
          waiting={state.waiting}
          onSubmitVotes={(successChains, favDrawing) => send(MSG_SUBMIT_VOTES, { successChains, favDrawing })}
          onPlayAgain={() => send(MSG_PLAY_AGAIN)}
          onHome={handleGoHome}
        />
      )}

      {state.screen === 'ai_error' && (
        <AIError
          message={state.aiError}
          onRestart={handleGoHome}
        />
      )}
    </div>
  );
}
