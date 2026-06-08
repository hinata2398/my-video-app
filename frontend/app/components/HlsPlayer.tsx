"use client";

import { useEffect, useRef, useState } from "react";
import type Hls from "hls.js";

type Props = {
  src: string;
  poster?: string;
};

type QualityLevel = {
  index: number;
  label: string;
};

export default function HlsPlayer({ src, poster }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [levels, setLevels] = useState<QualityLevel[]>([]);
  const [currentLevel, setCurrentLevel] = useState<number>(-1);
  const [menuOpen, setMenuOpen] = useState(false);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    let hlsInstance: Hls | null = null;

    import("hls.js").then(({ default: Hls }) => {
      if (Hls.isSupported()) {
        const hls = new Hls();
        hlsInstance = hls;
        hlsRef.current = hls;
        hls.loadSource(src);
        hls.attachMedia(video);

        hls.on(Hls.Events.MANIFEST_PARSED, () => {
          const qualityLevels = hls.levels.map((level, index) => ({
            index,
            label: `${level.height}p`,
          }));
          setLevels(qualityLevels);
          setCurrentLevel(-1);
        });
      } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
        video.src = src;
      }
    });

    return () => {
      hlsInstance?.destroy();
    };
  }, [src]);

  const handleQualityChange = (index: number) => {
    if (!hlsRef.current) return;
    hlsRef.current.currentLevel = index;
    setCurrentLevel(index);
    setMenuOpen(false);
  };

  const currentLabel =
    currentLevel === -1
      ? "自動"
      : (levels.find((l) => l.index === currentLevel)?.label ?? "自動");

  return (
    <div style={{ background: "#000" }}>
      <video
        ref={videoRef}
        controls
        poster={poster}
        style={{ width: "100%", maxHeight: 450, display: "block" }}
      />

      {/* 動画の下のコントロールバー */}
      {levels.length > 0 && (
        <div
          style={{
            position: "relative",
            background: "#111",
            padding: "6px 12px",
            display: "flex",
            alignItems: "center",
            gap: 8,
          }}
        >
          <span style={{ color: "#aaa", fontSize: "0.8rem" }}>⚙️ 画質</span>

          {/* 画質選択ボタン */}
          <div style={{ position: "relative" }}>
            <button
              onClick={() => setMenuOpen((prev) => !prev)}
              style={{
                background: "transparent",
                border: "1px solid #444",
                borderRadius: 4,
                color: "#fff",
                cursor: "pointer",
                padding: "3px 10px",
                fontSize: "0.85rem",
              }}
            >
              {currentLabel} ▲
            </button>

            {/* メニュー（上に展開） */}
            {menuOpen && (
              <div
                style={{
                  position: "absolute",
                  bottom: "calc(100% + 6px)",
                  left: 0,
                  background: "rgba(20,20,20,0.97)",
                  borderRadius: 8,
                  overflow: "hidden",
                  minWidth: 110,
                  boxShadow: "0 4px 16px rgba(0,0,0,0.5)",
                }}
              >
                {[{ index: -1, label: "自動" }, ...levels].map((level) => (
                  <button
                    key={level.index}
                    onClick={() => handleQualityChange(level.index)}
                    style={{
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "space-between",
                      width: "100%",
                      padding: "9px 16px",
                      background: "transparent",
                      border: "none",
                      color: "#fff",
                      fontSize: "0.9rem",
                      cursor: "pointer",
                      gap: 12,
                    }}
                  >
                    <span>{level.label}</span>
                    {currentLevel === level.index && <span>✓</span>}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* メニュー外クリックで閉じる */}
          {menuOpen && (
            <div
              onClick={() => setMenuOpen(false)}
              style={{ position: "fixed", inset: 0, zIndex: -1 }}
            />
          )}
        </div>
      )}
    </div>
  );
}
