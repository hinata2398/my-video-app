"use client";

import { useEffect, useRef, useState } from "react";

type Level = {
  index: number;
  label: string;
  height: number;
};

type Props = {
  src: string;
  poster?: string;
};

export default function HlsPlayer({ src, poster }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<import("hls.js").default | null>(null);
  const [levels, setLevels] = useState<Level[]>([]);
  const [selectedLevel, setSelectedLevel] = useState<number>(-1);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    // Safari はネイティブHLS対応
    if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = src;
      return;
    }

    import("hls.js").then(({ default: Hls }) => {
      if (!Hls.isSupported()) return;

      const hls = new Hls({ startLevel: -1 });
      hlsRef.current = hls;
      hls.loadSource(src);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, (_, data) => {
        const lvls: Level[] = data.levels
          .map((l, i) => ({
            index: i,
            height: l.height,
            label: l.height ? `${l.height}p` : `品質 ${i + 1}`,
          }))
          .sort((a, b) => b.height - a.height);
        setLevels(lvls);
        setSelectedLevel(-1);
      });
    });

    return () => {
      hlsRef.current?.destroy();
      hlsRef.current = null;
      setLevels([]);
    };
  }, [src]);

  const handleQualityChange = (levelIndex: number) => {
    const hls = hlsRef.current;
    if (!hls) return;
    if (levelIndex === -1) {
      hls.currentLevel = -1;
      hls.loadLevel = -1;
    } else {
      hls.loadLevel = levelIndex;
      hls.currentLevel = levelIndex;
    }
    setSelectedLevel(levelIndex);
  };

  return (
    <div>
      {/* 動画プレイヤー */}
      <video
        ref={videoRef}
        controls
        poster={poster}
        style={{ width: "100%", maxHeight: 450, display: "block" }}
      />

      {/* 画質セレクター：プレイヤーの下に配置してコントロールと重ならないようにする */}
      {levels.length > 0 && (
        <div style={{
          display: "flex", alignItems: "center", gap: "0.5rem",
          padding: "0.5rem 0.75rem",
          background: "#1a1a1a",
          borderTop: "1px solid #333",
        }}>
          <span style={{ color: "#aaa", fontSize: "0.8rem" }}>画質</span>
          <div style={{ display: "flex", gap: "0.35rem" }}>
            <button
              onClick={() => handleQualityChange(-1)}
              style={{
                padding: "0.2rem 0.6rem", borderRadius: 4, fontSize: "0.78rem",
                border: "1px solid",
                borderColor: selectedLevel === -1 ? "#e00" : "#555",
                background: selectedLevel === -1 ? "#e00" : "transparent",
                color: selectedLevel === -1 ? "#fff" : "#aaa",
                cursor: "pointer",
              }}
            >
              自動
            </button>
            {levels.map((l) => (
              <button
                key={l.index}
                onClick={() => handleQualityChange(l.index)}
                style={{
                  padding: "0.2rem 0.6rem", borderRadius: 4, fontSize: "0.78rem",
                  border: "1px solid",
                  borderColor: selectedLevel === l.index ? "#e00" : "#555",
                  background: selectedLevel === l.index ? "#e00" : "transparent",
                  color: selectedLevel === l.index ? "#fff" : "#aaa",
                  cursor: "pointer",
                }}
              >
                {l.label}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
