"use client";

import { useEffect, useRef, useState } from "react";

type Level = {
  index: number;
  label: string;  // "1080p", "720p" など
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
  const [currentLevel, setCurrentLevel] = useState<number>(-1); // -1 = 自動

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    // Safari はネイティブHLS対応のためhls.jsは不要
    if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = src;
      return;
    }

    import("hls.js").then(({ default: Hls }) => {
      if (!Hls.isSupported()) return;

      const hls = new Hls({ startLevel: -1 }); // 自動品質選択で開始
      hlsRef.current = hls;
      hls.loadSource(src);
      hls.attachMedia(video);

      // 解像度リストが確定したらセレクターに反映
      hls.on(Hls.Events.MANIFEST_PARSED, (_, data) => {
        const lvls: Level[] = data.levels.map((l, i) => ({
          index: i,
          height: l.height,
          label: l.height ? `${l.height}p` : `品質 ${i + 1}`,
        }));
        // 高解像度順に並べる
        lvls.sort((a, b) => b.height - a.height);
        setLevels(lvls);
        setCurrentLevel(-1);
      });

      // 実際に再生中のレベルが変わったら表示を更新
      hls.on(Hls.Events.LEVEL_SWITCHED, (_, data) => {
        // 手動選択中でなければ「自動」のまま
        if (hls.autoLevelEnabled) setCurrentLevel(-1);
      });
    });

    return () => {
      hlsRef.current?.destroy();
      hlsRef.current = null;
    };
  }, [src]);

  const handleQualityChange = (levelIndex: number) => {
    const hls = hlsRef.current;
    if (!hls) return;

    if (levelIndex === -1) {
      // 自動
      hls.currentLevel = -1;
      hls.loadLevel = -1;
    } else {
      hls.loadLevel = levelIndex;
      hls.currentLevel = levelIndex;
    }
    setCurrentLevel(levelIndex);
  };

  return (
    <div style={{ position: "relative", background: "#000" }}>
      <video
        ref={videoRef}
        controls
        poster={poster}
        style={{ width: "100%", maxHeight: 450, display: "block" }}
      />

      {/* 解像度セレクター（hls.jsが使える場合のみ表示） */}
      {levels.length > 0 && (
        <div style={{
          position: "absolute", bottom: 52, right: 12,
          display: "flex", alignItems: "center", gap: "0.4rem",
        }}>
          <span style={{ color: "#fff", fontSize: "0.75rem", opacity: 0.8 }}>画質</span>
          <select
            value={currentLevel}
            onChange={(e) => handleQualityChange(Number(e.target.value))}
            style={{
              background: "rgba(0,0,0,0.7)", color: "#fff",
              border: "1px solid rgba(255,255,255,0.3)", borderRadius: 4,
              padding: "0.2rem 0.4rem", fontSize: "0.8rem", cursor: "pointer",
            }}
          >
            <option value={-1}>自動</option>
            {levels.map((l) => (
              <option key={l.index} value={l.index}>
                {l.label}
              </option>
            ))}
          </select>
        </div>
      )}
    </div>
  );
}
