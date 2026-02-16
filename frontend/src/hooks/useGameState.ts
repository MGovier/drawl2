import { useReducer } from 'react';
import {
  ServerMessage, GameStateData, TurnStartData, Player, Chain,
  MSG_GAME_STATE, MSG_PLAYER_JOINED, MSG_PLAYER_LEFT, MSG_GAME_STARTED,
  MSG_TURN_START, MSG_TURN_TICK, MSG_WAITING, MSG_ROUND_COMPLETE,
  MSG_GAME_OVER, MSG_AI_ERROR, MSG_ERROR,
  MSG_SCORE_UPDATE, MSG_RETURN_TO_LOBBY,
  PHASE_PLAYING, PHASE_REVEAL,
} from '../lib/protocol';

export type Screen = 'home' | 'lobby' | 'playing' | 'reveal' | 'ai_error';

export interface GameState {
  screen: Screen;
  code: string;
  playerId: string;
  hostId: string;
  players: Player[];
  round: number;
  totalRounds: number;
  turnType: number;
  prompt: string;
  timeLimit: number;
  timeRemaining: number;
  waiting: boolean;
  chains: Chain[];
  scores: Record<string, number>;
  favDrawing: string;
  error: string;
  aiError: string;
}

const initialState: GameState = {
  screen: 'home',
  code: '',
  playerId: '',
  hostId: '',
  players: [],
  round: 0,
  totalRounds: 0,
  turnType: 0,
  prompt: '',
  timeLimit: 60,
  timeRemaining: 60,
  waiting: false,
  chains: [],
  scores: {},
  favDrawing: '',
  error: '',
  aiError: '',
};

type Action =
  | { type: 'ws_message'; msg: ServerMessage }
  | { type: 'set_screen'; screen: Screen }
  | { type: 'set_connection'; code: string; playerId: string }
  | { type: 'clear_error' };

function reducer(state: GameState, action: Action): GameState {
  switch (action.type) {
    case 'set_screen':
      return { ...state, screen: action.screen };
    case 'set_connection':
      return { ...state, code: action.code, playerId: action.playerId };
    case 'clear_error':
      return { ...state, error: '' };
    case 'ws_message':
      return handleWSMessage(state, action.msg);
    default:
      return state;
  }
}

function handleWSMessage(state: GameState, msg: ServerMessage): GameState {
  switch (msg.type) {
    case MSG_GAME_STATE: {
      const d = msg.data as GameStateData;
      let screen: Screen = 'lobby';
      if (d.phase === PHASE_PLAYING) screen = 'playing';
      if (d.phase === PHASE_REVEAL) screen = 'reveal';
      return {
        ...state,
        screen,
        code: d.code,
        playerId: d.playerId || state.playerId,
        hostId: d.hostId,
        players: d.players,
        round: d.round,
        totalRounds: d.totalRounds,
        scores: d.scores || state.scores,
      };
    }
    case MSG_PLAYER_JOINED:
      return {
        ...state,
        players: [...state.players.filter(p => p.id !== msg.data.player.id), msg.data.player],
      };
    case MSG_PLAYER_LEFT:
      return {
        ...state,
        players: state.players.filter(p => p.id !== msg.data.playerId),
        hostId: msg.data.hostId || state.hostId,
      };
    case MSG_GAME_STARTED:
      return { ...state, screen: 'playing', waiting: false };
    case MSG_TURN_START: {
      const d = msg.data as TurnStartData;
      return {
        ...state,
        screen: 'playing',
        round: d.round,
        totalRounds: d.totalRounds,
        turnType: d.turnType,
        prompt: d.prompt,
        timeLimit: d.timeLimit,
        timeRemaining: d.timeLimit,
        waiting: false,
      };
    }
    case MSG_TURN_TICK:
      return { ...state, timeRemaining: msg.data.remaining };
    case MSG_WAITING:
      return { ...state, waiting: true };
    case MSG_ROUND_COMPLETE:
      return { ...state, round: msg.data.round };
    case MSG_GAME_OVER:
      return { ...state, screen: 'reveal', chains: msg.data.chains, scores: msg.data.scores || state.scores, waiting: false, favDrawing: '' };
    case MSG_SCORE_UPDATE:
      return {
        ...state,
        scores: msg.data.scores || state.scores,
        favDrawing: msg.data.favDrawing || state.favDrawing,
      };
    case MSG_RETURN_TO_LOBBY:
      return {
        ...state,
        screen: 'lobby',
        players: msg.data.players,
        scores: msg.data.scores || state.scores,
        hostId: msg.data.hostId || state.hostId,
        waiting: false,
        chains: [],
        favDrawing: '',
      };
    case MSG_AI_ERROR:
      return { ...state, screen: 'ai_error', aiError: msg.data.message };
    case MSG_ERROR:
      return { ...state, error: msg.data.message };
    default:
      return state;
  }
}

export function useGameState() {
  return useReducer(reducer, initialState);
}
