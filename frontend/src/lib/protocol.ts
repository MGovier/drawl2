// Client -> Server message types
export const MSG_ADD_AI = 'add_ai';
export const MSG_START_GAME = 'start_game';
export const MSG_SUBMIT_DRAWING = 'submit_drawing';
export const MSG_SUBMIT_GUESS = 'submit_guess';
export const MSG_KICK_PLAYER = 'kick_player';
export const MSG_SUBMIT_VOTES = 'submit_votes';
export const MSG_PLAY_AGAIN = 'play_again';

// Server -> Client message types
export const MSG_GAME_STATE = 'game_state';
export const MSG_PLAYER_JOINED = 'player_joined';
export const MSG_PLAYER_LEFT = 'player_left';
export const MSG_GAME_STARTED = 'game_started';
export const MSG_TURN_START = 'turn_start';
export const MSG_TURN_TICK = 'turn_tick';
export const MSG_WAITING = 'waiting';
export const MSG_ROUND_COMPLETE = 'round_complete';
export const MSG_GAME_OVER = 'game_over';
export const MSG_AI_ERROR = 'ai_error';
export const MSG_ERROR = 'error';
export const MSG_SCORE_UPDATE = 'score_update';
export const MSG_RETURN_TO_LOBBY = 'return_to_lobby';

export const TURN_DRAW = 0;
export const TURN_GUESS = 1;

export const PHASE_LOBBY = 0;
export const PHASE_PLAYING = 1;
export const PHASE_REVEAL = 2;

export interface Player {
  id: string;
  name: string;
  type: number; // 0 = human, 1 = AI
  index: number;
}

export interface ChainEntry {
  playerId: string;
  type: number;
  drawing?: string;
  guess?: string;
}

export interface Chain {
  originalWord: string;
  ownerId: string;
  entries: ChainEntry[];
}

export interface GameStateData {
  code: string;
  phase: number;
  players: Player[];
  round: number;
  totalRounds: number;
  hostId: string;
  playerId?: string;
  scores?: Record<string, number>;
}

export interface TurnStartData {
  round: number;
  totalRounds: number;
  turnType: number;
  prompt: string;
  timeLimit: number;
}

export interface ServerMessage {
  type: string;
  data: any;
}
