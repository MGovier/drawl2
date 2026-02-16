interface Props {
  message: string;
  onRestart: () => void;
}

export function AIError({ message, onRestart }: Props) {
  return (
    <div className="ai-error">
      <div className="ai-error-robot">
        <svg width="120" height="120" viewBox="0 0 120 120" fill="none">
          <rect x="25" y="30" width="70" height="60" rx="12" fill="#2d2d4a" stroke="#6c5ce7" strokeWidth="3"/>
          <rect x="15" y="50" width="12" height="20" rx="4" fill="#2d2d4a" stroke="#6c5ce7" strokeWidth="2"/>
          <rect x="93" y="50" width="12" height="20" rx="4" fill="#2d2d4a" stroke="#6c5ce7" strokeWidth="2"/>
          <circle cx="45" cy="55" r="8" fill="#1a1a2e"/>
          <circle cx="45" cy="58" r="4" fill="#e74c3c"/>
          <circle cx="75" cy="55" r="8" fill="#1a1a2e"/>
          <circle cx="75" cy="58" r="4" fill="#e74c3c"/>
          <path d="M42 75 Q60 68 78 75" stroke="#e74c3c" strokeWidth="3" strokeLinecap="round" fill="none"/>
          <line x1="60" y1="15" x2="60" y2="30" stroke="#6c5ce7" strokeWidth="3"/>
          <circle cx="60" cy="12" r="5" fill="#e74c3c"/>
          <rect x="35" y="90" width="15" height="18" rx="4" fill="#2d2d4a" stroke="#6c5ce7" strokeWidth="2"/>
          <rect x="70" y="90" width="15" height="18" rx="4" fill="#2d2d4a" stroke="#6c5ce7" strokeWidth="2"/>
        </svg>
      </div>
      <h2 className="ai-error-title">AI Player Failed</h2>
      <p className="ai-error-message">{message}</p>
      <button className="btn btn-primary" onClick={onRestart}>
        Back to Home
      </button>
    </div>
  );
}
