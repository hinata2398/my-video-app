"use client";

import { useEffect, useRef } from "react";

type Props = {
  src: string;         // .m3u8 URL
  poster?: string;
};

export default function HlsPlayer({ src, poster }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    // ネイティブHLS対応ブラウザ（Safari）はそのまま再生
    if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = src;
      return;
    }

    // その他（Chrome/Firefox等）はhls.jsで再生
    import("hls.js").then(({ default: Hls }) => {
      if (!Hls.isSupported()) return;
      const hls = new Hls();
      hls.loadSource(src);
      hls.attachMedia(video);
      return () => hls.destroy();
    });
  }, [src]);

  return (
    <video
      ref={videoRef}
      controls
      poster={poster}
      style={{ width: "100%", maxHeight: 450, display: "block" }}
    />
  );
}
