import { useRef, useState, useCallback, useEffect, useImperativeHandle, forwardRef } from 'react';
import { Stroke, Point, redrawCanvas, drawStroke } from '../lib/canvas';

const COLORS = ['#000000', '#ff0000', '#0000ff', '#00aa00', '#ff8800', '#8800ff', '#00aaaa', '#888888'];
const SIZES = [3, 6, 12, 20];
const CANVAS_W = 600;
const CANVAS_H = 600;

export interface DrawingCanvasHandle {
  submit: () => void;
}

interface Props {
  prompt: string;
  onSubmit: (dataUrl: string) => void;
}

export const DrawingCanvas = forwardRef<DrawingCanvasHandle, Props>(function DrawingCanvas({ prompt, onSubmit }, ref) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [strokes, setStrokes] = useState<Stroke[]>([]);
  const [currentStroke, setCurrentStroke] = useState<Point[]>([]);
  const [color, setColor] = useState('#000000');
  const [brushSize, setBrushSize] = useState(6);
  const [isEraser, setIsEraser] = useState(false);
  const drawingRef = useRef(false);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d')!;
    redrawCanvas(ctx, strokes, CANVAS_W, CANVAS_H);
  }, [strokes]);

  const getPoint = useCallback((e: React.PointerEvent): Point => {
    const canvas = canvasRef.current!;
    const rect = canvas.getBoundingClientRect();
    return {
      x: ((e.clientX - rect.left) / rect.width) * CANVAS_W,
      y: ((e.clientY - rect.top) / rect.height) * CANVAS_H,
    };
  }, []);

  const handlePointerDown = useCallback((e: React.PointerEvent) => {
    e.preventDefault();
    (e.target as HTMLElement).setPointerCapture(e.pointerId);
    drawingRef.current = true;
    setCurrentStroke([getPoint(e)]);
  }, [getPoint]);

  const handlePointerMove = useCallback((e: React.PointerEvent) => {
    if (!drawingRef.current) return;
    e.preventDefault();
    const pt = getPoint(e);
    setCurrentStroke(prev => {
      const next = [...prev, pt];
      // Live preview
      const canvas = canvasRef.current;
      if (canvas) {
        const ctx = canvas.getContext('2d')!;
        drawStroke(ctx, {
          points: next,
          color: isEraser ? '#ffffff' : color,
          width: isEraser ? 20 : brushSize,
        });
      }
      return next;
    });
  }, [getPoint, color, brushSize, isEraser]);

  const handlePointerUp = useCallback(() => {
    if (!drawingRef.current) return;
    drawingRef.current = false;
    setCurrentStroke(prev => {
      if (prev.length > 0) {
        const stroke: Stroke = {
          points: prev,
          color: isEraser ? '#ffffff' : color,
          width: isEraser ? 20 : brushSize,
        };
        setStrokes(s => [...s, stroke]);
      }
      return [];
    });
  }, [color, brushSize, isEraser]);

  const handleUndo = () => {
    setStrokes(s => s.slice(0, -1));
  };

  const handleClear = () => {
    setStrokes([]);
  };

  const handleSubmit = useCallback(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    onSubmit(canvas.toDataURL('image/png'));
  }, [onSubmit]);

  useImperativeHandle(ref, () => ({ submit: handleSubmit }), [handleSubmit]);

  return (
    <div className="drawing-area">
      <div className="drawing-prompt">
        Draw: <strong>{prompt}</strong>
      </div>

      <canvas
        ref={canvasRef}
        width={CANVAS_W}
        height={CANVAS_H}
        className="drawing-canvas"
        onPointerDown={handlePointerDown}
        onPointerMove={handlePointerMove}
        onPointerUp={handlePointerUp}
        onPointerLeave={handlePointerUp}
        style={{ touchAction: 'none' }}
      />

      <div className="toolbar">
        <div className="color-picker">
          {COLORS.map(c => (
            <button
              key={c}
              className={`color-btn ${c === color && !isEraser ? 'active' : ''}`}
              style={{ backgroundColor: c }}
              onClick={() => { setColor(c); setIsEraser(false); }}
            />
          ))}
        </div>
        <div className="size-picker">
          {SIZES.map(s => (
            <button
              key={s}
              className={`size-btn ${s === brushSize && !isEraser ? 'active' : ''}`}
              onClick={() => { setBrushSize(s); setIsEraser(false); }}
            >
              <span className="size-dot" style={{ width: s, height: s }} />
            </button>
          ))}
        </div>
        <div className="tool-buttons">
          <button
            className={`btn btn-tool ${isEraser ? 'active' : ''}`}
            onClick={() => setIsEraser(!isEraser)}
          >
            Eraser
          </button>
          <button className="btn btn-tool" onClick={handleUndo} disabled={strokes.length === 0}>
            Undo
          </button>
          <button className="btn btn-tool" onClick={handleClear} disabled={strokes.length === 0}>
            Clear
          </button>
        </div>
      </div>

      <button className="btn btn-primary btn-submit" onClick={handleSubmit}>
        Submit Drawing
      </button>
    </div>
  );
});
