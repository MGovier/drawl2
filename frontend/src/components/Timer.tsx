interface Props {
  remaining: number;
  total: number;
}

export function Timer({ remaining, total }: Props) {
  const pct = (remaining / total) * 100;
  const urgent = remaining <= 10;

  return (
    <div className="timer">
      <div className="timer-bar">
        <div
          className={`timer-fill ${urgent ? 'urgent' : ''}`}
          style={{ width: `${pct}%` }}
        />
      </div>
      <span className={`timer-text ${urgent ? 'urgent' : ''}`}>{remaining}s</span>
    </div>
  );
}
