"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";

type Video = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  created_at: string;
};

export default function VideoDetailPage() {
  const { id } = useParams();
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`)
      .then((res) => {
        if (!res.ok) { setNotFound(true); return null; }
        return res.json();
      })
      .then((data) => { if (data) setVideo(data); })
      .catch(() => setNotFound(true))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return (
      <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
        <p style={{ color: "#888" }}>読み込み中...</p>
      </main>
    );
  }

  if (notFound || !video) {
    return (
      <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
        <p>動画が見つかりません。</p>
        <Link href="/" style={{ color: "#e00", textDecoration: "none" }}>← 一覧へ戻る</Link>
      </main>
    );
  }

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
      <Link href="/" style={{ color: "#888", textDecoration: "none" }}>← 一覧へ戻る</Link>
      <div style={{ background: "#111", marginTop: "1rem", borderRadius: 8, height: 400, display: "flex", alignItems: "center", justifyContent: "center" }}>
        {video.thumbnail_url ? (
          <img src={video.thumbnail_url} alt={video.title} style={{ width: "100%", height: "100%", objectFit: "cover", borderRadius: 8 }} />
        ) : (
          <span style={{ color: "#555", fontSize: "4rem" }}>▶</span>
        )}
      </div>
      <h1 style={{ marginTop: "1rem" }}>{video.title}</h1>
      <p style={{ color: "#888", fontSize: "0.9rem" }}>{new Date(video.created_at).toLocaleDateString("ja-JP")}</p>
      {video.description && <p style={{ marginTop: "1rem", lineHeight: 1.7 }}>{video.description}</p>}
    </main>
  );
}
